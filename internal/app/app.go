package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Georgi-Progger/task-tracker-backend/internal/config"
	"github.com/Georgi-Progger/task-tracker-backend/internal/handler"
	"github.com/Georgi-Progger/task-tracker-backend/internal/repo"
	"github.com/Georgi-Progger/task-tracker-backend/internal/service"
	"github.com/Georgi-Progger/task-tracker-common/kafka/producer"
	"github.com/Georgi-Progger/task-tracker-common/logger"
	"github.com/Georgi-Progger/task-tracker-common/postgres"
	"github.com/Georgi-Progger/task-tracker-common/redis"
	"github.com/Georgi-Progger/task-tracker-rate-limiter/limiter"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	_ "github.com/lib/pq"
)

func Run() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.ContextTimeout(10 * time.Second))
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))

	logger := logger.NewLogger()

	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Error(err, "Failed to load config")
		os.Exit(1)
	}

	db, err := postgres.NewDb(cfg.GetUrlDb(), logger)
	if err != nil {
		slog.Error("db connect failed")
	}
	defer func() {
		if err := db.Close(); err != nil {
			logger.Error(err, "db close failed")
		}
	}()

	redisClient, err := redis.NewRedisClient(context.Background(), cfg.GetUrlRedis(), logger)
	if err != nil {
		slog.Error("redis connect failed")
	}

	rateLimiter := limiter.NewLimiter(redisClient)

	jwtSecret := os.Getenv("JWT_SECRET")

	producer := producer.NewProducer(cfg.GetUrlBroker(), "EMAIL_SENDING_TASKS", logger) // TODO: FIX THIS SHIT

	repo := repo.NewRepository(db)
	service := service.NewService(repo, jwtSecret, &producer, 15*time.Minute)
	handler := handler.NewHandler(service, *rateLimiter, logger)
	handler.SetupRoutes(e)

	server := &http.Server{
		Addr:              fmt.Sprintf(":%s", cfg.GetPort()),
		Handler:           e,
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		logger.Info("Server started")
		err := server.ListenAndServe()
		if err != nil {
			logger.Error(err, "Server started error")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Server work ended")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		logger.Error(err, "Server work ended error")
	}

	logger.Info("Server is ended")
}
