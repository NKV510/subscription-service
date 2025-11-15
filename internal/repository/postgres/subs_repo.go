package postgres

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/NKV510/subscription-service/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SubscriptionRepository struct {
	pool *pgxpool.Pool
}

func NewSubscriptionRepository(pool *pgxpool.Pool) *SubscriptionRepository {
	return &SubscriptionRepository{pool: pool}
}

func (r *SubscriptionRepository) Create(ctx context.Context, sub *models.Subscription) error {
	query := `
        INSERT INTO subscriptions (service_name, price, user_id, start_date, end_date)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id, created_at, updated_at
    `

	err := r.pool.QueryRow(
		ctx,
		query,
		sub.ServiceName,
		sub.Price,
		sub.UserID,
		sub.StartDate,
		sub.EndDate,
	).Scan(&sub.ID)

	if err != nil {
		slog.Error("Failed to create subscription", "error", err)
		return fmt.Errorf("failed to create subscription: %w", err)
	}

	slog.Info("Subscription created successfully", "id", sub.ID)
	return nil
}

func (r *SubscriptionRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Subscription, error) {
	query := `
        SELECT id, service_name, price, user_id, start_date, end_date, created_at, updated_at
        FROM subscriptions 
        WHERE id = $1
    `

	var sub models.Subscription
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&sub.ID,
		&sub.ServiceName,
		&sub.Price,
		&sub.UserID,
		&sub.StartDate,
		&sub.EndDate,
	)

	if err != nil {
		slog.Error("Failed to get subscription by ID", "id", id, "error", err)
		return nil, fmt.Errorf("subscription not found: %w", err)
	}

	return &sub, nil
}
