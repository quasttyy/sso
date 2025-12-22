package main

import (
	"fmt"
	"log/slog"
	"os"
	"sso/internal/app"
	"sso/internal/config"
)

const (
	envLocal = "local"
	envDev = "dev"
	envProd = "prod"
)

func main() {
	// Инициализируем конфиг
	cfg := config.MustLoad()
	fmt.Println(cfg)
	// Инициализируем логгер
	logger := setUpLogger(cfg.Env)
	logger.Info("logger successfully initialized")

	application := app.New(logger, cfg.GRPC.Port, cfg.Dsn, cfg.TokenTTL)

	application.GRPCServer.MustRun()
	// todo: инициализировать приложение (app)

	// todo: запустить gRPC сервер приложения
}

func setUpLogger(env string) *slog.Logger {
	var log *slog.Logger
	
	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}