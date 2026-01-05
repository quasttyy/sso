package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env string `yaml:"env" env-default:"local"`
	TokenTTL time.Duration `yaml:"token_ttl" env-required:"true"` 
	GRPC GRPCConfig `yaml:"grpc"`
	Postgres PostgresConfig `yaml:"postgres"`
}

type GRPCConfig struct {
	Port int `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

type PostgresConfig struct {
	DSN             string `yaml:"dsn"`
	MaxConns        int32 `yaml:"max_conns"`
	MinConns        int32 `yaml:"min_conns"`
	MaxConnLifeTime time.Duration `yaml:"max_conn_lifetime"`
	MaxConnIdleTime time.Duration `yaml:"max_conn_idletime"`
}

func MustLoad() *Config {
	path := fetchConfigPath()
	if path == "" {
		panic("config path is empty.")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file doesn't exist: " + path)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("failed to parse config file: " + err.Error())
	}

	return &cfg
}

// fetchConfigPath возвращает путь до файла конфига из флага либо переменной окружения.
// Приоритет: flag > env > default
func fetchConfigPath() string {
	var path string

	flag.StringVar(&path, "config", "", "path to config file")
	flag.Parse()

	if path == "" {
		path = os.Getenv("CONFIG_PATH")
	}

	return path
}