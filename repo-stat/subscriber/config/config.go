package config

import (
	"repo-stat/platform/env"
	"repo-stat/platform/grpcserver"
	"repo-stat/platform/logger"
)

type App struct {
	AppName string `yaml:"app_name" env:"APP_NAME" env-default:"repo-stat-subscriber"`
}

type Database struct {
	DSN            string `yaml:"dsn" env:"DATABASE_DSN" env-required:"true"`
	MigrationsPath string `yaml:"migrations_path" env:"MIGRATIONS_PATH" env-default:"file://subscriber/migrations"`
}

type Config struct {
	App      App               `yaml:"app"`
	Database Database          `yaml:"database"`
	GRPC     grpcserver.Config `yaml:"grpc"`
	Logger   logger.Config     `yaml:"logger"`
}

func MustLoad(path string) Config {
	var cfg Config
	env.MustLoad(path, &cfg)
	return cfg
}
