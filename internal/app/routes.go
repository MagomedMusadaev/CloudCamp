package app

import (
	"CloudCamp/internal/handler"
	"net/http"
)

// setupRoutes настраивает маршруты для сервера
func (s *Server) setupRoutes() {
	// Создаем обработчики
	proxyHandler := handler.NewProxyHandler(s.balancer)

	// Настраиваем маршруты
	mux := http.NewServeMux()

	// Маршрут для прокси
	mux.HandleFunc("/", proxyHandler.ServeHTTP)

	// Устанавливаем обработчик для сервера
	s.httpServer.Handler = mux
}
