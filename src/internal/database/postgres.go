package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DbName   string
}

type Pool struct {
	*pgxpool.Pool
}

func NewPool(ctx context.Context, config Config) (*Pool, error) {
	dsn := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s", config.User, config.Password, config.Host, config.Port, config.DbName)

	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("postgres - failed to parse config: %w", err)
	}

	poolConfig.MaxConns = 10

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("postgres - failed to connect to database: %w", err)
	}

	return &Pool{pool}, nil
}

func (p *Pool) Close() {
	p.Pool.Close()
}
