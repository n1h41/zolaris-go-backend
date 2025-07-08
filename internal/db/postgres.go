package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/afreedicp/zolaris-backend-app/internal/config"
)

// PostgresDB represents a PostgreSQL database connection
type PostgresDB struct {
	pool *pgxpool.Pool
}

// NewPostgresDB creates a new PostgreSQL database connection
func NewPostgresDB(ctx context.Context, cfg *config.Config) (*PostgresDB, error) {
	// Build the connection string
	connString := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.PostgresHost,
		cfg.Database.PostgresPort,
		cfg.Database.PostgresUser,
		cfg.Database.PostgresPassword,
		cfg.Database.PostgresDBName,
		cfg.Database.PostgresSSLMode,
	)

	// Configure connection pool
	pgConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse postgres config: %w", err)
	}

	// Set connection pool options
	pgConfig.MaxConns = 10
	pgConfig.MinConns = 2
	pgConfig.MaxConnLifetime = 30 * time.Minute
	pgConfig.MaxConnIdleTime = 5 * time.Minute
	pgConfig.HealthCheckPeriod = 1 * time.Minute

	// Create the connection pool
	pool, err := pgxpool.NewWithConfig(ctx, pgConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	// Test the connection
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping postgres: %w", err)
	}

	return &PostgresDB{
		pool: pool,
	}, nil
}

// GetPool returns the underlying pgxpool.Pool
func (db *PostgresDB) GetPool() *pgxpool.Pool {
	return db.pool
}

// Close closes the database connection
func (db *PostgresDB) Close() {
	if db.pool != nil {
		db.pool.Close()
	}
}
