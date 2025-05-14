package app

import (
	balancerDir "CloudCamp/internal/balancer"
	"CloudCamp/internal/config"
	"CloudCamp/internal/domain/balancerDomain"
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
	httpServer *http.Server
}

// NewServer создает новый сервер
func NewServer(cfg *config.Config) (*Server, error) {
	// Создаем балансировщик через фабрику
	balancer := balancerDir.NewBalancerFactory(cfg)
	if balancer == nil {
		return nil, fmt.Errorf("unsupported balancing strategy: %s", cfg.Balancer.Strategy)
	}

	return &Server{
		cfg:      cfg,
		balancer: balancer,
	}, nil
}

func (s *Server) Run() error {
	// Создаем HTTP-сервер
	addr := fmt.Sprintf(":%d", s.cfg.Server.Port)
	s.httpServer = &http.Server{
		Addr: addr,
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
