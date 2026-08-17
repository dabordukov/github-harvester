package config

import (
	"repo-stat/platform/env"
	"repo-stat/platform/grpcserver"
	"repo-stat/platform/logger"
)

type App struct {
	AppName string `yaml:"app_name" env:"APP_NAME" env-default:"repo-stat-collector"`
}

type Services struct {
	Subscriber string `yaml:"subscriber" env:"SUBSCRIBER_ADDRESS" env-default:"localhost:8081"`
}

type Kafka struct {
	Brokers      []string `yaml:"brokers" env:"KAFKA_BROKERS" env-default:"kafka:9092"`
	RequestTopic string   `yaml:"request_topic" env:"KAFKA_REQUEST_TOPIC" env-default:"collect-requests"`
	ResultTopic  string   `yaml:"result_topic" env:"KAFKA_RESULT_TOPIC" env-default:"collect-results"`
	GroupID      string   `yaml:"group_id" env:"KAFKA_GROUP_ID" env-default:"collector-group"`
}

type Config struct {
	App      App               `yaml:"app"`
	Services Services          `yaml:"services"`
	Kafka    Kafka             `yaml:"kafka"`
	GRPC     grpcserver.Config `yaml:"grpc"`
	Logger   logger.Config     `yaml:"logger"`
}

func MustLoad(path string) Config {
	var cfg Config
	env.MustLoad(path, &cfg)
	return cfg
}
