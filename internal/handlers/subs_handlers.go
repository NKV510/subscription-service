package handlers

import (
	"log/slog"
	"net/http"

	"github.com/NKV510/subscription-service/internal/models"
	"github.com/NKV510/subscription-service/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SubscriptionHandler struct {
	service *service.SubscriptionService
}

func NewSubscriptionHandler(service *service.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{service: service}
}

func (h *SubscriptionHandler) CreateSubscription(c *gin.Context) {
	var req models.CreateSubscriptionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Warn("Invalid request body", "error", err)
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	subscription, err := h.service.CreateSubscription(c.Request.Context(), req)
	if err != nil {
		slog.Error("Failed to create subscription", "error", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Internal server error"})
		return
	}

	c.JSON(http.StatusCreated, subscription)
}

func (h *SubscriptionHandler) GetSubscriptionByID(c *gin.Context) {
	idStr := c.Param("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		slog.Warn("Invalid UUID format", "id", idStr, "error", err)
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid subscription ID"})
		return
	}

	subscription, err := h.service.GetSubscriptionByID(c.Request.Context(), id)
	if err != nil {
		slog.Error("Failed to get subscription", "id", id, "error", err)
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Subscription not found"})
		return
	}

	c.JSON(http.StatusOK, subscription)
}
