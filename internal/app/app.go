package app

import (
	"log/slog"
	grpcapp "sso/internal/app/grpc"
	"time"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(
	log *slog.Logger,
	gRPCPort int,
	dsn string,
	tokenTTL time.Duration,
) *App {
	// todo: connect to storage

	// todo: init auth service

	gRPCApp := grpcapp.New(log, gRPCPort)

	return &App{
		GRPCServer: gRPCApp,
	}
}
