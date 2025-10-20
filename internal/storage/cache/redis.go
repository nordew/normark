package cache

import (
	"context"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/redis/go-redis/v9"
)

type Config struct {
	Addr     string
	Password string
	DB       int
}

type SetOptions struct {
	Key        string
	Value      any
	Expiration time.Duration
}

type Redis struct {
	client *redis.Client
}

func New(cfg Config) *Redis {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	return &Redis{
		client: client,
	}
}

func (r *Redis) Ping(ctx context.Context) error {
	if err := r.client.Ping(ctx).Err(); err != nil {
		return errors.Wrap(err, "failed to ping redis")
	}
	return nil
}

func (r *Redis) Close() error {
	if err := r.client.Close(); err != nil {
		return errors.Wrap(err, "failed to close redis connection")
	}
	return nil
}

func (r *Redis) Get(ctx context.Context, key string) (string, error) {
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", errors.Wrap(err, "key not found")
		}
		return "", errors.Wrap(err, "failed to get value")
	}
	return val, nil
}

func (r *Redis) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	if err := r.client.Set(ctx, key, value, expiration).Err(); err != nil {
		return errors.Wrap(err, "failed to set value")
	}
	return nil
}

func (r *Redis) SetWithOptions(ctx context.Context, opts SetOptions) error {
	if err := r.client.Set(ctx, opts.Key, opts.Value, opts.Expiration).Err(); err != nil {
		return errors.Wrap(err, "failed to set value with options")
	}
	return nil
}

func (r *Redis) Delete(ctx context.Context, keys ...string) error {
	if err := r.client.Del(ctx, keys...).Err(); err != nil {
		return errors.Wrap(err, "failed to delete keys")
	}
	return nil
}

func (r *Redis) Exists(ctx context.Context, keys ...string) (int64, error) {
	count, err := r.client.Exists(ctx, keys...).Result()
	if err != nil {
		return 0, errors.Wrap(err, "failed to check existence")
	}
	return count, nil
}

func (r *Redis) Expire(ctx context.Context, key string, expiration time.Duration) error {
	if err := r.client.Expire(ctx, key, expiration).Err(); err != nil {
		return errors.Wrap(err, "failed to set expiration")
	}
	return nil
}

func (r *Redis) TTL(ctx context.Context, key string) (time.Duration, error) {
	duration, err := r.client.TTL(ctx, key).Result()
	if err != nil {
		return 0, errors.Wrap(err, "failed to get ttl")
	}
	return duration, nil
}

func (r *Redis) Increment(ctx context.Context, key string) (int64, error) {
	val, err := r.client.Incr(ctx, key).Result()
	if err != nil {
		return 0, errors.Wrap(err, "failed to increment")
	}
	return val, nil
}

func (r *Redis) IncrementBy(ctx context.Context, key string, value int64) (int64, error) {
	val, err := r.client.IncrBy(ctx, key, value).Result()
	if err != nil {
		return 0, errors.Wrap(err, "failed to increment by value")
	}
	return val, nil
}

func (r *Redis) Decrement(ctx context.Context, key string) (int64, error) {
	val, err := r.client.Decr(ctx, key).Result()
	if err != nil {
		return 0, errors.Wrap(err, "failed to decrement")
	}
	return val, nil
}

func (r *Redis) DecrementBy(ctx context.Context, key string, value int64) (int64, error) {
	val, err := r.client.DecrBy(ctx, key, value).Result()
	if err != nil {
		return 0, errors.Wrap(err, "failed to decrement by value")
	}
	return val, nil
}

func (r *Redis) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	ok, err := r.client.SetNX(ctx, key, value, expiration).Result()
	if err != nil {
		return false, errors.Wrap(err, "failed to set if not exists")
	}
	return ok, nil
}

func (r *Redis) GetSet(ctx context.Context, key string, value interface{}) (string, error) {
	val, err := r.client.GetSet(ctx, key, value).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", errors.Wrap(err, "key not found")
		}
		return "", errors.Wrap(err, "failed to get and set")
	}
	return val, nil
}

func (r *Redis) MGet(ctx context.Context, keys ...string) ([]interface{}, error) {
	vals, err := r.client.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get multiple values")
	}
	return vals, nil
}

type MSetOptions struct {
	Pairs []any
}

func (r *Redis) MSet(ctx context.Context, opts MSetOptions) error {
	if err := r.client.MSet(ctx, opts.Pairs...).Err(); err != nil {
		return errors.Wrap(err, "failed to set multiple values")
	}
	return nil
}

func (r *Redis) FlushDB(ctx context.Context) error {
	if err := r.client.FlushDB(ctx).Err(); err != nil {
		return errors.Wrap(err, "failed to flush database")
	}
	return nil
}

func (r *Redis) Client() *redis.Client {
	return r.client
}
