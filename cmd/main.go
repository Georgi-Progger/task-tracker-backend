package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.ContextTimeout(10 * time.Second))

	server := &http.Server{
		Addr:              net.JoinHostPort("localhost", "8080"),
		Handler:           e,
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		log.Print("Сервер запущен")
		err := server.ListenAndServe()
		if err != nil {
			log.Print("Ошибка запуска сервера")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Print("Завершение работы сервера")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := server.Shutdown(ctx)
	if err != nil {
		log.Printf("Ошибка при остановке сервера %v", err)
	}

	log.Println("Сервер остановлен")
}
