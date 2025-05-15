package app

import (
	"CloudCamp/internal/handler"
	"net/http"
)

// setupRoutes настраивает маршруты для сервера
func (s *Server) setupRoutes() {
	// Создаем обработчики
	proxyHandler := handler.NewProxyHandler(s.balancer)
	clientHandler := handler.NewClientHandler(s.limiter)
	rateLimiterMiddleware := handler.NewRateLimiterMiddleware(s.limiter)

	// Настраиваем маршруты
	mux := http.NewServeMux()

	// Маршрут для прокси
	mux.HandleFunc("/", proxyHandler.ServeHTTP)

	// Маршруты для управления клиентами
	mux.HandleFunc("/clients", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			clientHandler.CreateClient(w, r)
		case http.MethodDelete:
			clientHandler.DeleteClient(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Оборачиваем все маршруты в middleware для rate limiting
	s.httpServer.Handler = rateLimiterMiddleware.Middleware(mux)
}
