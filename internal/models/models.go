package models

import (
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	ID          uuid.UUID  `json:"id"`
	ServiceName string     `json:"service_name" binding:"required"`
	Price       int        `json:"price" binding:"required,min=1"`
	UserID      uuid.UUID  `json:"user_id" binding:"required"`
	StartDate   time.Time  `json:"start_date" binding:"required"`
	EndDate     *time.Time `json:"end_date,omitempty"`
}

type CreateSubscriptionRequest struct {
	ServiceName string    `json:"service_name" binding:"required"`
	Price       int       `json:"price" binding:"required,min=1"`
	UserID      uuid.UUID `json:"user_id" binding:"required"`
	StartDate   string    `json:"start_date" binding:"required"` // формат "MM-YYYY"
}

type ErrorResponse struct {
	Error string `json:"error"`
}
type UpdateSubscriptionRequest struct {
	ServiceName *string `json:"service_name,omitempty"`
	Price       *int    `json:"price,omitempty"`
	StartDate   *string `json:"start_date,omitempty"` // формат "MM-YYYY"
	EndDate     *string `json:"end_date,omitempty"`   // формат "MM-YYYY"
}

type TotalSpentRequest struct {
	From        string     `form:"from" binding:"required"` // формат "MM-YYYY"
	To          string     `form:"to" binding:"required"`   // формат "MM-YYYY"
	UserID      *uuid.UUID `form:"user_id,omitempty"`
	ServiceName *string    `form:"service_name,omitempty"`
}

type TotalSpentResponse struct {
	Total int `json:"total"`
}
