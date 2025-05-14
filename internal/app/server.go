package app

import (
	balancerDir "CloudCamp/internal/balancer"
	"CloudCamp/internal/config"
	"CloudCamp/internal/domain/balancerDomain"
	"CloudCamp/internal/limiter"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

// Server представляет собой HTTP-сервер с балансировщиком нагрузки
type Server struct {
	cfg        *config.Config
	balancer   balancerDomain.Strategy
	limiter    *limiter.MemoryRateLimiter
	httpServer *http.Server
}

// NewServer создает новый сервер
func NewServer(cfg *config.Config) (*Server, error) {
	// Создаем балансировщик через фабрику
	balancer := balancerDir.NewBalancerFactory(cfg)
	if balancer == nil {
		return nil, fmt.Errorf("unsupported balancing strategy: %s", cfg.Balancer.Strategy)
	}

	// Создаем rate limiter
	rl := limiter.NewMemoryRateLimiter()

	return &Server{
		cfg:      cfg,
		balancer: balancer,
		limiter:  rl,
	}, nil
}

func (s *Server) Run() error {
	// Создаем HTTP-сервер
	addr := fmt.Sprintf(":%d", s.cfg.Server.Port)
	s.httpServer = &http.Server{
		Addr: addr,
	}

	// Если включен rate limiting, устанавливаем лимиты
	if s.cfg.RateLimiter.Enabled {
		// Устанавливаем глобальный лимит
		err := s.limiter.SetLimit("global", s.cfg.RateLimiter.Rate, s.cfg.RateLimiter.Period)
		if err != nil {
			return err
		}

		// Устанавливаем индивидуальные лимиты для клиентов
		for clientID, limit := range s.cfg.RateLimiter.Clients {
			err = s.limiter.SetClientLimit(clientID, limit.RateLimit, limit.Period)
			if err != nil {
				return fmt.Errorf("failed to set client limit for %s: %w", clientID, err)
			}
		}
	}

	// Настраиваем маршруты
	s.setupRoutes()

	slog.Info("starting server", slog.String("addr", addr))
	return s.httpServer.ListenAndServe()
}

// Shutdown выполняет корректное завершение работы сервера
func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Останавливаем HTTP-сервер
	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("error shutting down server: %w", err)
	}

	return nil
}
