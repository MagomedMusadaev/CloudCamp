package app

import (
	"CloudCamp/internal/handler"
	"net/http"
)

// setupRoutes настраивает маршруты для сервера
func (s *Server) setupRoutes() {
	// Создаем обработчики
	proxyHandler := handler.NewProxyHandler(s.balancer)
	rateLimiterMiddleware := handler.NewRateLimiterMiddleware(s.limiter)

	// Настраиваем маршруты
	mux := http.NewServeMux()

	// Маршрут для прокси
	mux.HandleFunc("/", proxyHandler.ServeHTTP)

	// Оборачиваем все маршруты в middleware для rate limiting
	s.httpServer.Handler = rateLimiterMiddleware.Middleware(mux)
}
