package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

func NewPostgres(dsn string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database config: %w", err)
	}

	config.MaxConns = 25
	config.MinConns = 5
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = 30 * time.Minute
	config.HealthCheckPeriod = time.Minute

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Info().Msg("Database connection established")

	return pool, nil
}

func RunMigrations(pool *pgxpool.Pool, migrationsPath string) error {
	ctx := context.Background()

	migrations := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			wallet_address VARCHAR(42) UNIQUE NOT NULL,
			nonce VARCHAR(64) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE INDEX IF NOT EXISTS idx_users_wallet ON users(wallet_address)`,
		`CREATE TABLE IF NOT EXISTS transactions (
			id SERIAL PRIMARY KEY,
			tx_hash VARCHAR(66) UNIQUE NOT NULL,
			event_type VARCHAR(20) NOT NULL CHECK (event_type IN ('deposit', 'withdrawal', 'transfer')),
			from_address VARCHAR(42) NOT NULL,
			to_address VARCHAR(42),
			amount VARCHAR(78) NOT NULL,
			block_number BIGINT NOT NULL,
			block_timestamp BIGINT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE INDEX IF NOT EXISTS idx_transactions_from ON transactions(from_address)`,
		`CREATE INDEX IF NOT EXISTS idx_transactions_timestamp ON transactions(block_timestamp DESC)`,
		`CREATE INDEX IF NOT EXISTS idx_transactions_block ON transactions(block_number)`,
	}

	for i, migration := range migrations {
		_, err := pool.Exec(ctx, migration)
		if err != nil {
			return fmt.Errorf("migration %d failed: %w", i+1, err)
		}
	}

	log.Info().Msg("Migrations completed successfully")
	return nil
}

func Ping(pool *pgxpool.Pool) error {
	return pool.Ping(context.Background())
}
