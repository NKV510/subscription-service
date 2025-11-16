package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/NKV510/subscription-service/internal/models"
	"github.com/NKV510/subscription-service/internal/repository/postgres"
	"github.com/google/uuid"
)

type SubscriptionService struct {
	repo *postgres.SubscriptionRepository
}

func NewSubscriptionService(repo *postgres.SubscriptionRepository) *SubscriptionService {
	return &SubscriptionService{repo: repo}
}

func (s *SubscriptionService) CreateSubscription(
	ctx context.Context,
	req models.CreateSubscriptionRequest,
) (*models.Subscription, error) {
	// Парсим дату из формата "MM-YYYY"
	startDate, err := time.Parse("01-2006", req.StartDate)
	if err != nil {
		slog.Error("Invalid start date format", "date", req.StartDate, "error", err)
		return nil, fmt.Errorf("invalid start date format, expected MM-YYYY: %w", err)
	}

	// Устанавливаем начало месяца
	startDate = time.Date(startDate.Year(), startDate.Month(), 1, 0, 0, 0, 0, time.UTC)

	subscription := &models.Subscription{
		ID:          uuid.New(),
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserID:      req.UserID,
		StartDate:   startDate,
		EndDate:     nil, // По умолчанию подписка бессрочная
	}

	if err := s.repo.Create(ctx, subscription); err != nil {
		return nil, err
	}

	return subscription, nil
}

func (s *SubscriptionService) GetSubscriptionByID(
	ctx context.Context,
	id uuid.UUID,
) (*models.Subscription, error) {
	return s.repo.GetByID(ctx, id)
}

// UpdateSubscription обновляет подписку
func (s *SubscriptionService) UpdateSubscription(
	ctx context.Context,
	id uuid.UUID,
	req models.UpdateSubscriptionRequest,
) (*models.Subscription, error) {
	// Получаем существующую подписку
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Обновляем поля если они предоставлены
	if req.ServiceName != nil {
		existing.ServiceName = *req.ServiceName
	}
	if req.Price != nil {
		existing.Price = *req.Price
	}
	if req.StartDate != nil {
		startDate, err := time.Parse("01-2006", *req.StartDate)
		if err != nil {
			return nil, fmt.Errorf("invalid start date format: %w", err)
		}
		existing.StartDate = time.Date(startDate.Year(), startDate.Month(), 1, 0, 0, 0, 0, time.UTC)
	}
	if req.EndDate != nil {
		if *req.EndDate == "" {
			existing.EndDate = nil
		} else {
			endDate, err := time.Parse("01-2006", *req.EndDate)
			if err != nil {
				return nil, fmt.Errorf("invalid end date format: %w", err)
			}
			// Устанавливаем конец месяца
			endDate = time.Date(endDate.Year(), endDate.Month()+1, 0, 0, 0, 0, 0, time.UTC)
			existing.EndDate = &endDate
		}
	}

	if err := s.repo.Update(ctx, existing); err != nil {
		return nil, err
	}

	return existing, nil
}

// DeleteSubscription удаляет подписку
func (s *SubscriptionService) DeleteSubscription(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

// GetSubscriptionsByUserID возвращает подписки пользователя
func (s *SubscriptionService) GetSubscriptionsByUserID(
	ctx context.Context,
	userID uuid.UUID,
) ([]*models.Subscription, error) {
	return s.repo.GetByUserID(ctx, userID)
}

// GetTotalSpent вычисляет суммарные траты за период
func (s *SubscriptionService) GetTotalSpent(
	ctx context.Context,
	fromStr string,
	toStr string,
	userID *uuid.UUID,
	serviceName *string,
) (int, error) {
	// Парсим даты
	from, err := time.Parse("01-2006", fromStr)
	if err != nil {
		return 0, fmt.Errorf("invalid 'from' date format: %w", err)
	}
	from = time.Date(from.Year(), from.Month(), 1, 0, 0, 0, 0, time.UTC)

	to, err := time.Parse("01-2006", toStr)
	if err != nil {
		return 0, fmt.Errorf("invalid 'to' date format: %w", err)
	}
	// Устанавливаем конец месяца для 'to'
	to = time.Date(to.Year(), to.Month()+1, 0, 23, 59, 59, 0, time.UTC)

	return s.repo.GetTotalSpent(ctx, from, to, userID, serviceName)
}
