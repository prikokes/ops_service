package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Server     ServerConfig
	Database   DatabaseConfig
	JWT        JWTConfig
	GRPC       GRPCConfig
	Prometheus PrometheusConfig
}

type ServerConfig struct {
	Port int
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

type JWTConfig struct {
	Secret string
	TTL    time.Duration
}

type GRPCConfig struct {
	Port int
}

type PrometheusConfig struct {
	Port int
}

func LoadConfig() *Config {
	config := &Config{
		Server: ServerConfig{
			Port: getEnvAsInt("SERVER_PORT", 8080),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvAsInt("DB_PORT", 5432),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "ops_service"),
		},
		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET", "your_secret_key"),
			TTL:    getEnvAsDuration("JWT_TTL", 24*time.Hour),
		},
		GRPC: GRPCConfig{
			Port: getEnvAsInt("GRPC_PORT", 3000),
		},
		Prometheus: PrometheusConfig{
			Port: getEnvAsInt("PROMETHEUS_PORT", 9000),
		},
	}

	return config
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	if value, exists := os.LookupEnv(key); exists {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
