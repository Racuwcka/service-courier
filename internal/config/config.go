package config

import (
	"log"

	"go.uber.org/config"
)

type AppConfig struct {
	Name string `yaml:"name"`
	Port int    `yaml:"port"`
}

type KafkaConfig struct {
	Brokers    []string `yaml:"brokers"`
	OrderTopic string   `yaml:"order_topic"`
}

type Config struct {
	App   AppConfig   `yaml:"app"`
	Kafka KafkaConfig `yaml:"kafka"`
}

func MustLoad() *Config {
	provider, err := config.NewYAML(config.File("config.yaml"))
	if err != nil {
		log.Fatalf("config file does not exist: %s", err)
	}

	var cfg Config
	if err = provider.Get("").Populate(&cfg); err != nil {
		log.Fatalf("config file does not exist: %s", err)
	}

	return &cfg
}
