package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"time"
)

type Environment string

const (
	EnvDev  Environment = "development"
	EnvProd Environment = "production"
	EnvTest Environment = "test"
)

type Config struct {
	Env           Environment         `yaml:"env"`
	Server        ServerConfig        `yaml:"server"`
	Balancer      BalancerConfig      `yaml:"balancer"`
	RateLimiter   RateLimitConfig     `yaml:"rate_limiter"`
	HealthChecker HealthCheckerConfig `yaml:"health_checker"`
	Log           LogConfig           `yaml:"log"`
}

// ServerConfig — содержит настройки сервера
type ServerConfig struct {
	Port int `yaml:"port"`
}

// ClientLimit содержит настройки лимита для конкретного клиента
type ClientLimit struct {
	ClientID  int           `yaml:"client_id"` // Уникальный идентификатор клиента
	RateLimit int           `yaml:"rate"`      // Персональный лимит: максимальное количество запросов, которые клиент может сделать за период
	Period    time.Duration `yaml:"period"`    // Персональный период (в секундах), с учётом которого будут добавляться токены для клиента
}

// RateLimitConfig конфигурация для rate limiter
type RateLimitConfig struct {
	Enabled  bool                   `yaml:"enabled"`  // Включение или выключение rate limiter
	Interval time.Duration          `yaml:"interval"` // Интервал времени для глобального лимита
	Rate     int                    `yaml:"rate"`     // Глобальный лимит: максимальное количество запросов для всех клиентов за глобальный период
	Period   time.Duration          `yaml:"period"`   // Глобальный период, в течение которого обновляются токены для всех клиентов
	Clients  map[string]ClientLimit `yaml:"clients"`  // Карта индивидуальных лимитов для каждого клиента (ключ — ID клиента)
}

// HealthCheckerConfig — содержит настройки проверки нод
type HealthCheckerConfig struct {
	Enabled  bool          `yaml:"enabled"`
	Interval time.Duration `yaml:"interval"`
	Path     string        `yaml:"path"`
}

// LogConfig - содержит настройки slog
type LogConfig struct {
	FilePath string `yaml:"file_path"`
	Dir      string `yaml:"dir"`
}

// BalancerConfig содержит настройки балансировщика
type BalancerConfig struct {
	Backends []string `yaml:"backends"`
	Strategy string   `yaml:"strategy"` // round-robin, least-connections, random
}

// LoadConfig — читает YAML-файл конфигурации и возвращает заполненную структуру
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err = yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	switch config.Env {
	case EnvDev, EnvProd, EnvTest:
	default:
		return nil, fmt.Errorf("invalid environment: %s", config.Env)
	}

	if config.Balancer.Strategy == "" {
		config.Balancer.Strategy = "round-robin" // Default strategy
	}

	return &config, nil
}

func (l *LogConfig) GetLogDir() string {
	return l.Dir
}
