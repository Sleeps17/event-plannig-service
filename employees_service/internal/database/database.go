package database

import (
	"context"
	"fmt"
	"github.com/Sleeps17/events-planning-backend-service/employees_service/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	Pool *pgxpool.Pool
}

func MustInit(ctx context.Context, cfg *config.PostgresConfig) *Database {
	pool, err := pgxpool.New(ctx, cfg.Url())
	if err != nil {
		panic(fmt.Sprintf("failed to create new pool: %v", err))
	}

	conn, err := pool.Acquire(ctx)
	if err != nil {
		panic(fmt.Sprintf("failed to acquire connection: %v", err))
	}

	if err := conn.Conn().Ping(ctx); err != nil {
		panic(fmt.Sprintf("failed to ping connection: %v", err))
	}

	return &Database{Pool: pool}
}

func (d *Database) Close() {
	d.Pool.Close()
}
