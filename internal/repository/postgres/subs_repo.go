package postgres

import (
	"context"
	"fmt"
	"log/slog"
	"time"

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

func (r *SubscriptionRepository) Update(ctx context.Context, sub *models.Subscription) error {
	query := `
        UPDATE subscriptions 
        SET service_name = $1, price = $2, start_date = $3, end_date = $4
        WHERE id = $5
    `

	result, err := r.pool.Exec(
		ctx,
		query,
		sub.ServiceName,
		sub.Price,
		sub.StartDate,
		sub.EndDate,
		sub.ID,
	)

	if err != nil {
		slog.Error("Failed to update subscription", "id", sub.ID, "error", err)
		return fmt.Errorf("failed to update subscription: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("subscription not found")
	}

	slog.Info("Subscription updated successfully", "id", sub.ID)
	return nil
}

// Delete удаляет подписку
func (r *SubscriptionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM subscriptions WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		slog.Error("Failed to delete subscription", "id", id, "error", err)
		return fmt.Errorf("failed to delete subscription: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("subscription not found")
	}

	slog.Info("Subscription deleted successfully", "id", id)
	return nil
}

// GetByUserID возвращает все подписки пользователя
func (r *SubscriptionRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Subscription, error) {
	query := `
        SELECT id, service_name, price, user_id, start_date, end_date
        FROM subscriptions 
        WHERE user_id = $1
        ORDER BY start_date DESC
    `

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		slog.Error("Failed to get subscriptions by user ID", "user_id", userID, "error", err)
		return nil, fmt.Errorf("failed to get subscriptions: %w", err)
	}
	defer rows.Close()

	var subscriptions []*models.Subscription
	for rows.Next() {
		var sub models.Subscription
		err := rows.Scan(
			&sub.ID,
			&sub.ServiceName,
			&sub.Price,
			&sub.UserID,
			&sub.StartDate,
			&sub.EndDate,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan subscription: %w", err)
		}
		subscriptions = append(subscriptions, &sub)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	slog.Info("Retrieved subscriptions by user ID", "user_id", userID, "count", len(subscriptions))
	return subscriptions, nil
}

func (r *SubscriptionRepository) GetTotalSpent(
	ctx context.Context,
	from time.Time,
	to time.Time,
	userID *uuid.UUID,
	serviceName *string,
) (int, error) {
	query := `
        SELECT COALESCE(SUM(price), 0)
        FROM subscriptions 
        WHERE start_date <= $2 
          AND (end_date IS NULL OR end_date >= $1)
    `

	args := []interface{}{from, to}
	argIndex := 3

	if userID != nil {
		query += fmt.Sprintf(" AND user_id = $%d", argIndex)
		args = append(args, *userID)
		argIndex++
	}

	if serviceName != nil {
		query += fmt.Sprintf(" AND service_name = $%d", argIndex)
		args = append(args, *serviceName)
	}

	var total int
	err := r.pool.QueryRow(ctx, query, args...).Scan(&total)
	if err != nil {
		slog.Error("Failed to calculate total spent",
			"from", from, "to", to, "user_id", userID, "service_name", serviceName, "error", err)
		return 0, fmt.Errorf("failed to calculate total spent: %w", err)
	}

	slog.Info("Calculated total spent",
		"from", from, "to", to, "user_id", userID, "service_name", serviceName, "total", total)
	return total, nil
}
