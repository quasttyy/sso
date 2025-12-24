package grpcapp

import (
	"fmt"
	"log/slog"
	"net"
	authGRPC "sso/internal/grpc/auth"

	"google.golang.org/grpc"
)

type App struct {
	log *slog.Logger
	gRPCServer *grpc.Server
	port int 
}

// New - конструктор объекта приложения
func New(
	log *slog.Logger,
	port int,
) *App {
	gRPCServer := grpc.NewServer() // Создаём новый gRPC сервер

	authGRPC.Register(gRPCServer) // Регистрируем хендлеры

	return &App{
		log: log,
		gRPCServer: gRPCServer,
		port: port,
	} // Возвращаем объект приложения, содержащий логгер, сервер и порт
}

// MustRun вызывает метод Run() и паникует при ошибке
func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

// Run запускает gRPC сервер
func (a *App) Run() error {
	const op = "grpcapp.Run"

	log := a.log.With(
		slog.String("op:", op),
		slog.Int("port:", a.port),
	)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port)) // Начинаем слушать на переданном порте
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("gRPC server is running", slog.String("addr:", listener.Addr().String()))

	if err := a.gRPCServer.Serve(listener); err != nil { // Принимаем tcp соединения через listener
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil 
}

// Stop graceful останавливает работу gRPC сервера
func (a *App) Stop() {
	const op = "grpcapp.Stop"

	a.log.With(slog.String("op:", op)).
	Info("stopping gRPC server", slog.Int("port:", a.port))

	a.gRPCServer.GracefulStop()
}