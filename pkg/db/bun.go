package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
	"github.com/user/normark/internal/config"
	"github.com/user/normark/internal/entity"
)

const (
	defaultConnectionTimeout = 10 * time.Second
	defaultDialTimeout       = 5 * time.Second
	defaultReadTimeout       = 5 * time.Second
	defaultWriteTimeout      = 5 * time.Second

	maxOpenConnections    = 25
	maxIdleConnections    = 5
	connectionMaxLifetime = 30 * time.Minute
	connectionMaxIdleTime = 5 * time.Minute

	healthCheckTimeout = 2 * time.Second
)

type DB struct {
	*bun.DB
}

func NewPostgresConnection(ctx context.Context, cfg *config.Postgres) (*DB, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
		cfg.SSLMode,
	)

	connector := pgdriver.NewConnector(
		pgdriver.WithDSN(dsn),
		pgdriver.WithTimeout(defaultConnectionTimeout),
		pgdriver.WithDialTimeout(defaultDialTimeout),
		pgdriver.WithReadTimeout(defaultReadTimeout),
		pgdriver.WithWriteTimeout(defaultWriteTimeout),
	)

	sqlDB := sql.OpenDB(connector)

	sqlDB.SetMaxOpenConns(maxOpenConnections)
	sqlDB.SetMaxIdleConns(maxIdleConnections)
	sqlDB.SetConnMaxLifetime(connectionMaxLifetime)
	sqlDB.SetConnMaxIdleTime(connectionMaxIdleTime)
	bunDB := bun.NewDB(sqlDB, pgdialect.New())

	bunDB.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithVerbose(true),
		bundebug.FromEnv("BUNDEBUG"),
	))

	if err := bunDB.PingContext(ctx); err != nil {
		sqlDB.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	db := &DB{DB: bunDB}

	if err := db.AutoMigrate(ctx); err != nil {
		sqlDB.Close()
		return nil, fmt.Errorf("failed to run auto-migration: %w", err)
	}

	return db, nil
}

func (db *DB) Close() error {
	return db.DB.Close()
}

func (db *DB) Ping(ctx context.Context) error {
	return db.DB.PingContext(ctx)
}

func (db *DB) HealthCheck(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, healthCheckTimeout)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("database health check failed: %w", err)
	}

	return nil
}

func (db *DB) AutoMigrate(ctx context.Context) error {
	models := []any{
		(*entity.User)(nil),
	}

	if _, err := db.ExecContext(ctx, "CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\""); err != nil {
		return fmt.Errorf("failed to create uuid-ossp extension: %w", err)
	}

	if _, err := db.ExecContext(ctx, "CREATE EXTENSION IF NOT EXISTS pgcrypto"); err != nil {
		return fmt.Errorf("failed to create pgcrypto extension: %w", err)
	}

	for _, model := range models {
		if _, err := db.NewCreateTable().Model(model).IfNotExists().Exec(ctx); err != nil {
			return fmt.Errorf("failed to create table: %w", err)
		}
	}

	return nil
}
