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
	"github.com/Georgi-Progger/task-tracker-backend/pkg/datasource"
	logger "github.com/Georgi-Progger/task-tracker-backend/pkg/looger"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	_ "github.com/lib/pq"
)

func Run() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.ContextTimeout(10 * time.Second))

	logger := logger.NewLogger()

	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Error(err, "Failed to load config")
		os.Exit(1)
	}

	db, err := datasource.NewDb(cfg)
	if err != nil {
		slog.Error("db connect failed")
	}
	defer func() {
		if err := db.Close(); err != nil {
			logger.Error(err, "db close failed")
		}
	}()

	jwtSecret := os.Getenv("JWT_SECRET")

	repo := repo.NewRepository(db)
	service := service.NewService(repo, jwtSecret, 15*time.Minute, logger)
	handler := handler.NewHandler(service)
	handler.SetupRoutes(e)

	server := &http.Server{
		Addr:              fmt.Sprintf(":%s", cfg.Port),
		Handler:           e,
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		logger.Print("Server started")
		err := server.ListenAndServe()
		if err != nil {
			logger.Error(err, "Server started error")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Print("Server work ended")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		logger.Printf("Server work ended error %v", err)
	}

	logger.Println("Server is ended")
}
