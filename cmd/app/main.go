package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/NKV510/subscription-service/internal/config"
	"github.com/NKV510/subscription-service/internal/handlers"
	"github.com/NKV510/subscription-service/internal/repository/postgres"
	"github.com/NKV510/subscription-service/internal/service"
	"github.com/NKV510/subscription-service/pkg/database"
	"github.com/gin-gonic/gin"
)

func main() {
	// Инициализация логгера
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Загрузка конфигурации
	cfg := config.Load()

	// Подключение к БД
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := database.NewDBPool(ctx, cfg)
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	// Инициализация слоев
	repo := postgres.NewSubscriptionRepository(pool)
	subscriptionService := service.NewSubscriptionService(repo)
	subscriptionHandler := handlers.NewSubscriptionHandler(subscriptionService)

	// Настройка роутера
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(handlers.LoggingMiddleware())

	// Маршруты
	api := router.Group("/api/v1")
	{
		subscriptions := api.Group("/subscriptions")
		{
			subscriptions.POST("", subscriptionHandler.CreateSubscription)
			subscriptions.GET("/:id", subscriptionHandler.GetSubscriptionByID)
		}
	}

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Graceful shutdown
	srv := &http.Server{
		Addr:    cfg.Server.Port,
		Handler: router,
	}

	go func() {
		slog.Info("Starting server", "port", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Failed to start server", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("Shutting down server...")

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
	}

	slog.Info("Server exited")
}
