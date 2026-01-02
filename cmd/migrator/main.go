package main

import (
	"log/slog"
	"os"
	"path/filepath"
	"sso/internal/config"
	loggerpkg "sso/internal/logger"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	// Инициализируем конфиг
	cfg := config.MustLoad()

	// Инициализируем логгер
	log := loggerpkg.Setup(cfg.Env)

	// Указываем путь к миграциям
	migrationsPath, err := filepath.Abs("migrations")
	if err != nil {
		log.Error("failed to create migrations path", slog.Any("error", err))
		os.Exit(1)
	}

	sourceURL := "file://" + migrationsPath
	dsn := cfg.Dsn

	log.Info("running migrations...")
	log.Info("migration source url", "value", sourceURL)
	log.Info("dsn", "value", dsn)

	// Создаём мигратор
	m, err := migrate.New(sourceURL, dsn)
	if err != nil {
		log.Error("failed to create migrator", slog.Any("error", err))
		os.Exit(1)
	}

	// Закрываем его после завершения
	defer func() {
		_, _ = m.Close()
	}()

	// Применяем миграции
	err = m.Up()
	if err != nil {
		if err == migrate.ErrNoChange {
			log.Info("no new migrations — DB up to date")
			os.Exit(0)
		}
		log.Error("failed to run migrations", slog.Any("error", err))
		os.Exit(1)
	}

	log.Info("migrations applied successfully")
}
