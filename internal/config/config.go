package config

import (
	"log"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	HTTP     HTTPConfig
	Database DatabaseConfig
	Logger   LoggerConfig
}

type HTTPConfig struct {
	Port    int           `env:"HTTP_PORT" env-required:"true"`
	Timeout time.Duration `env:"HTTP_TIMEOUT" env-required:"true"`
}

type DatabaseConfig struct {
	Host     string `env:"DATABASE_HOST" env-required:"true"`
	Port     string `env:"DATABASE_PORT" env-required:"true"`
	User     string `env:"DATABASE_USER" env-required:"true"`
	Password string `env:"DATABASE_PASSWORD" env-required:"true"`
	Name     string `env:"DATABASE_NAME" env-required:"true"`
}

type LoggerConfig struct {
	Level string `env:"LOGGER_LEVEL" env-required:"true"`
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Warning: Could not load .env file, falling back to system environment variables")
	}

	var config Config

	err = cleanenv.ReadEnv(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func MustLoadConfig() *Config {
	config, err := LoadConfig()
	if err != nil {
		log.Fatalf("Could not load configuration: %v", err)
	}
	return config
}
