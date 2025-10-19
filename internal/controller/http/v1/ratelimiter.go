package v1

import (
	"net"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/user/normark/internal/config"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

const (
	headerXForwardedFor = "X-Forwarded-For"
	headerXRealIP       = "X-Real-IP"
	rateLimitExceeded   = "rate limit exceeded"
)

type RateLimiter struct {
	visitors map[string]*rate.Limiter
	mu       sync.RWMutex
	rps      int
	burst    int
	logger   *zap.Logger
}

func NewRateLimiter(cfg *config.RateLimit, logger *zap.Logger) *RateLimiter {
	return &RateLimiter{
		visitors: make(map[string]*rate.Limiter),
		rps:      cfg.RequestsPerSecond,
		burst:    cfg.Burst,
		logger:   logger,
	}
}

func (rl *RateLimiter) getVisitor(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	limiter, exists := rl.visitors[ip]
	if !exists {
		limiter = rate.NewLimiter(rate.Limit(rl.rps), rl.burst)
		rl.visitors[ip] = limiter
	}

	return limiter
}

func (rl *RateLimiter) cleanupVisitors() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	for ip := range rl.visitors {
		delete(rl.visitors, ip)
	}
}

func (rl *RateLimiter) getIP(c *gin.Context) string {
	forwarded := c.GetHeader(headerXForwardedFor)
	if forwarded != "" {
		ip, _, err := net.SplitHostPort(forwarded)
		if err == nil {
			return ip
		}
		return forwarded
	}

	realIP := c.GetHeader(headerXRealIP)
	if realIP != "" {
		return realIP
	}

	ip, _, err := net.SplitHostPort(c.Request.RemoteAddr)
	if err != nil {
		return c.Request.RemoteAddr
	}

	return ip
}

func (rl *RateLimiter) Limit() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := rl.getIP(c)
		limiter := rl.getVisitor(ip)

		if !limiter.Allow() {
			rl.logger.Error(rateLimitExceeded, zap.String("ip", ip))
			newErrorResponse(c, http.StatusTooManyRequests, rateLimitExceeded)
			return
		}

		c.Next()
	}
}
