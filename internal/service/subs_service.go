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
