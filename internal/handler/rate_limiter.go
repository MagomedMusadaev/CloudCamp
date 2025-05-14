package handler

import (
	"CloudCamp/internal/domain/limiter"
	"CloudCamp/pkg/utils"
	"log/slog"
	"net/http"
	"strings"
)

// RateLimiterMiddleware middleware для ограничения частоты запросов
type RateLimiterMiddleware struct {
	limiter limiter.RateLimiter
}

// NewRateLimiterMiddleware создает новый middleware для rate limiting
func NewRateLimiterMiddleware(limiter limiter.RateLimiter) *RateLimiterMiddleware {
	return &RateLimiterMiddleware{
		limiter: limiter,
	}
}

// Middleware возвращает HTTP middleware для rate limiting
func (m *RateLimiterMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Получаем IP-адрес клиента
		clientIP := getClientIP(r)

		// Получаем идентификатор клиента из заголовка или используем IP
		clientID := r.Header.Get("X-Client-ID")
		if clientID == "" {
			clientID = clientIP
		}

		// Проверяем, не превышен ли лимит
		if !m.limiter.Allow(clientID) {
			slog.Warn("rate limit exceeded",
				slog.String("client_id", clientID),
				slog.String("client_ip", clientIP),
			)
			utils.SendRateLimitExceeded(w,
				http.StatusTooManyRequests,
				"Too Many Requests",
			)
			return
		}

		// Передаем управление следующему обработчику
		next.ServeHTTP(w, r)
	})
}

// getClientIP извлекает IP-адрес клиента из запроса
func getClientIP(r *http.Request) string {
	// Проверяем заголовок X-Forwarded-For
	ip := r.Header.Get("X-Forwarded-For")
	if ip != "" {
		// Берем первый IP из списка
		if i := strings.Index(ip, ","); i > 0 {
			ip = ip[:i]
		}
		return strings.TrimSpace(ip)
	}

	// Проверяем заголовок X-Real-IP
	ip = r.Header.Get("X-Real-IP")
	if ip != "" {
		return strings.TrimSpace(ip)
	}

	// Получаем IP из RemoteAddr
	ip = r.RemoteAddr
	if i := strings.LastIndex(ip, ":"); i > 0 {
		ip = ip[:i]
	}
	return strings.TrimSpace(ip)
}
