package main

import (
	"fmt"
	"sso/internal/config"
)

func main() {
	// todo: инициализировать модель конфига
	cfg := config.MustLoad()
	fmt.Println(cfg)
	// todo: инициализировать логгер

	// todo: инициализировать приложение (app)

	// todo: запустить gRPC сервер приложения
}