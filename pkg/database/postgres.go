package database

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/NKV510/subscription-service/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewDBPool(ctx context.Context, cfg *config.Config) (*pgxpool.Pool, error) {
	dbURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
		cfg.Database.SSLMode,
	)

	dbConfig, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database config: %w", err)
	}

	dbConfig.MaxConns = cfg.Database.Max_DB_Conns
	dbConfig.HealthCheckPeriod = 1 * time.Minute
	dbConfig.MaxConnLifetime = 1 * time.Hour
	dbConfig.MaxConnIdleTime = 30 * time.Minute

	var dbPool *pgxpool.Pool
	maxRetries := 10
	retryDelay := 3 * time.Second

	for i := 0; i < maxRetries; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		slog.Info("Attempt %d/%d: Connecting to database at %s:%s", i+1, maxRetries, cfg.Database.Host, cfg.Database.Port)

		dbPool, err = pgxpool.NewWithConfig(ctx, dbConfig)
		if err != nil {
			cancel()
			slog.Debug("Attempt %d/%d: Failed to connect to database: %v", i+1, maxRetries, err)
			if i < maxRetries-1 {
				slog.Debug("Retrying in %v...", retryDelay)
				time.Sleep(retryDelay)
				continue
			}
			return nil, fmt.Errorf("unable to connect to database after %d attempts: %w", maxRetries, err)
		}

		if err := dbPool.Ping(ctx); err != nil {
			cancel()
			slog.Debug("Attempt %d/%d: Database ping failed: %v", i+1, maxRetries, err)
			dbPool.Close()
			if i < maxRetries-1 {
				slog.Debug("Retrying in %v...", retryDelay)
				time.Sleep(retryDelay)
				continue
			}
			return nil, fmt.Errorf("unable to ping database after %d attempts: %w", maxRetries, err)
		}

		cancel()
		slog.Info("Successfully connected to database on attempt %d", i+1)
		break
	}

	return dbPool, nil
}
