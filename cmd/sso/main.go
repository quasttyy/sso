package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"sso/internal/app"
	"sso/internal/config"
	"sso/internal/logger"
	"syscall"
)

func main() {
	// Инициализируем конфиг
	cfg := config.MustLoad()
	fmt.Println(cfg)
	// Инициализируем логгер
	logger := logger.Setup(cfg.Env)
	logger.Info("logger successfully initialized")

	application := app.New(logger, cfg.GRPC.Port, cfg.Dsn, cfg.TokenTTL)

	go application.GRPCServer.MustRun()
	// todo: инициализировать приложение (app)
	// todo: запустить gRPC сервер приложения

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	gotSignal := <- stop
	logger.Info("stopping application...", slog.Any("signal:", gotSignal))

	application.GRPCServer.Stop()
	logger.Info("application is stopped")
}
