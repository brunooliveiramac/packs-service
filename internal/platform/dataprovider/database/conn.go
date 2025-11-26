package database

import (
	"context"
	"fmt"

	"github.com/brunooliveiramac/packs-service/internal/platform/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	pool *pgxpool.Pool
}

func NewFromEnv(ctx context.Context) (*DB, error) {
	cfg := config.LoadDB()
	conn := cfg.URL
	if conn == "" {
		conn = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database, cfg.SSLMode)
	}
	p, err := pgxpool.New(ctx, conn)
	if err != nil {
		return nil, err
	}
	return &DB{pool: p}, nil
}

func (d *DB) Close() { if d.pool != nil { d.pool.Close() } }
func (d *DB) Pool() *pgxpool.Pool { return d.pool }


