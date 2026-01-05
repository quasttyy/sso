package storage

import (
	"context"
	"fmt"
	"sso/internal/config"
	"time"
	"github.com/jackc/pgx/v5/pgxpool"
)

// NewPostgresPool создаёт пул соединений с PostgreSQL
func NewPostgresPool(
	ctx context.Context,
	cfg *config.Config,
) (*pgxpool.Pool, error) {
	if cfg.Postgres.DSN == "" {
		return nil, fmt.Errorf("postgres dsn is empty")
	}

	poolCfg, err := pgxpool.ParseConfig(cfg.Postgres.DSN)
	if err != nil {
		return nil, fmt.Errorf("err parsing pgx config: %w", err)
	}

	// Настройки пула
	poolCfg.MaxConns = cfg.Postgres.MaxConns
	poolCfg.MinConns = cfg.Postgres.MinConns
	poolCfg.MaxConnLifetime = cfg.Postgres.MaxConnLifeTime
	poolCfg.MaxConnIdleTime = cfg.Postgres.MaxConnIdleTime

	poolCfg.ConnConfig.ConnectTimeout = 5 * time.Second

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("err creating pool: %w", err)
	}

	// Проверяем доступность БД
	ctxPing, cancel := context.WithTimeout(ctx, 3 * time.Second)
	defer cancel()

	if err := pool.Ping(ctxPing); err != nil {
		pool.Close()
		return nil, fmt.Errorf("err pinging postgres: %w", err)
	}

	return pool, nil
}
