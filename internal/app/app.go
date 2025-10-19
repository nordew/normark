package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/user/normark/internal/config"
	v1 "github.com/user/normark/internal/controller/http/v1"
	"github.com/user/normark/internal/service"
	"github.com/user/normark/internal/storage"
	"github.com/user/normark/pkg/auth"
	"github.com/user/normark/pkg/db"
	"go.uber.org/zap"
)

const (
	shutdownTimeout = 10 * time.Second
)

type App struct {
	cfg    *config.Config
	logger *zap.Logger
	db     *db.DB
	server *http.Server
}

func New() (*App, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	logger, err := zap.NewProduction()
	if err != nil {
		return nil, fmt.Errorf("failed to create logger: %w", err)
	}

	return &App{
		cfg:    cfg,
		logger: logger,
	}, nil
}

func (a *App) Run() error {
	ctx := context.Background()

	if err := a.initDatabase(ctx); err != nil {
		return err
	}

	if err := a.initServer(); err != nil {
		return err
	}

	return a.start()
}

func (a *App) initDatabase(ctx context.Context) error {
	database, err := db.NewPostgresConnection(ctx, &a.cfg.Postgres)
	if err != nil {
		a.logger.Error("failed to connect to database", zap.Error(err))
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	a.db = database
	a.logger.Info("database connected successfully")
	return nil
}

func (a *App) initServer() error {
	jwtManager, err := auth.NewJWTManager(
		a.cfg.JWT.Secret,
		a.cfg.JWT.AccessTokenExpiry,
		a.cfg.JWT.RefreshTokenExpiry,
	)
	if err != nil {
		a.logger.Error("failed to create jwt manager", zap.Error(err))
		return fmt.Errorf("failed to create jwt manager: %w", err)
	}

	userStorage := storage.NewUserStorage(a.db.DB)
	userService := service.NewUserService(userStorage, jwtManager, a.logger)

	tradingJournalStorage := storage.NewTradingJournalStorage(a.db.DB)
	tradingJournalService := service.NewTradingJournalService(tradingJournalStorage, a.logger)

	tradingJournalEntryStorage := storage.NewTradingJournalEntryStorage(a.db.DB)
	tradingJournalEntryService := service.NewTradingJournalEntryService(tradingJournalEntryStorage, tradingJournalStorage, a.logger)

	middleware := v1.NewMiddleware(a.logger, jwtManager, &a.cfg.CORS)
	rateLimiter := v1.NewRateLimiter(&a.cfg.RateLimit, a.logger)
	handler := v1.NewHandler(userService, tradingJournalService, tradingJournalEntryService, a.logger, middleware, rateLimiter)

	router := handler.InitRoutes()

	addr := ":" + a.cfg.Server.Port

	a.server = &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	a.logger.Info("server initialized", zap.String("addr", addr))
	return nil
}

func (a *App) start() error {
	errChan := make(chan error, 1)

	go func() {
		a.logger.Info("starting server", zap.String("addr", a.server.Addr))
		if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errChan <- err
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	select {
	case err := <-errChan:
		a.logger.Error("server error", zap.Error(err))
		return err
	case sig := <-quit:
		a.logger.Info("shutdown signal received", zap.String("signal", sig.String()))
		return a.shutdown()
	}
}

func (a *App) shutdown() error {
	a.logger.Info("shutting down gracefully")

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := a.server.Shutdown(ctx); err != nil {
		a.logger.Error("server shutdown error", zap.Error(err))
		return fmt.Errorf("server shutdown error: %w", err)
	}

	if err := a.db.Close(); err != nil {
		a.logger.Error("database close error", zap.Error(err))
		return fmt.Errorf("database close error: %w", err)
	}

	if err := a.logger.Sync(); err != nil {
		return fmt.Errorf("logger sync error: %w", err)
	}

	a.logger.Info("shutdown completed")
	return nil
}
