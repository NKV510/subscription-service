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

// CreateSubscription создает новую подписку
// @Summary Create subscription
// @Description Create a new subscription
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param input body models.CreateSubscriptionRequest true "Subscription data"
// @Success 201 {object} models.Subscription
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /subscriptions [post]

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

// GetSubscriptionByID получает подписку по ID
// @Summary Get subscription by ID
// @Description Get subscription by its ID
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path string true "Subscription ID"
// @Success 200 {object} models.Subscription
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /subscriptions/{id} [get]
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

// UpdateSubscription обновляет подписку
// @Summary Update subscription
// @Description Update existing subscription
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path string true "Subscription ID"
// @Param input body models.UpdateSubscriptionRequest true "Subscription update data"
// @Success 200 {object} models.Subscription
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /subscriptions/{id} [put]
func (h *SubscriptionHandler) UpdateSubscription(c *gin.Context) {
	idStr := c.Param("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		slog.Warn("Invalid UUID format", "id", idStr, "error", err)
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid subscription ID"})
		return
	}

	var req models.UpdateSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Warn("Invalid request body", "error", err)
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	subscription, err := h.service.UpdateSubscription(c.Request.Context(), id, req)
	if err != nil {
		slog.Error("Failed to update subscription", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, subscription)
}

// DeleteSubscription удаляет подписку
// @Summary Delete subscription
// @Description Delete subscription by ID
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path string true "Subscription ID"
// @Success 204
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /subscriptions/{id} [delete]
func (h *SubscriptionHandler) DeleteSubscription(c *gin.Context) {
	idStr := c.Param("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		slog.Warn("Invalid UUID format", "id", idStr, "error", err)
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid subscription ID"})
		return
	}

	if err := h.service.DeleteSubscription(c.Request.Context(), id); err != nil {
		slog.Error("Failed to delete subscription", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Internal server error"})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetSubscriptionsByUserID возвращает подписки пользователя
// @Summary Get user subscriptions
// @Description Get all subscriptions for a user
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param user_id query string true "User ID"
// @Success 200 {array} models.Subscription
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /subscriptions [get]
func (h *SubscriptionHandler) GetSubscriptionsByUserID(c *gin.Context) {
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "user_id query parameter is required"})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		slog.Warn("Invalid user_id format", "user_id", userIDStr, "error", err)
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid user_id format"})
		return
	}

	subscriptions, err := h.service.GetSubscriptionsByUserID(c.Request.Context(), userID)
	if err != nil {
		slog.Error("Failed to get subscriptions by user ID", "user_id", userID, "error", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, subscriptions)
}

// GetTotalSpent вычисляет суммарные траты за период
// @Summary Calculate total spent
// @Description Calculate total amount spent on subscriptions for a period
// @Tags analytics
// @Accept json
// @Produce json
// @Param from query string true "Start date (MM-YYYY)"
// @Param to query string true "End date (MM-YYYY)"
// @Param user_id query string false "User ID filter"
// @Param service_name query string false "Service name filter"
// @Success 200 {object} models.TotalSpentResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /analytics/total [get]
func (h *SubscriptionHandler) GetTotalSpent(c *gin.Context) {
	var req models.TotalSpentRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		slog.Warn("Invalid query parameters", "error", err)
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	total, err := h.service.GetTotalSpent(c.Request.Context(), req.From, req.To, req.UserID, req.ServiceName)
	if err != nil {
		slog.Error("Failed to calculate total spent", "error", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, models.TotalSpentResponse{Total: total})
}
