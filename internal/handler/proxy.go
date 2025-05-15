package handler

import (
	"CloudCamp/internal/domain/balancerDomain"
	"CloudCamp/pkg/utils"
	"io"
	"log/slog"
	"net/http"
	"net/url"
)

// ProxyHandler обработчик для проксирования запросов
type ProxyHandler struct {
	balancer balancerDomain.Strategy
	client   *http.Client
}

// ErrorResponse структура для ошибок, отправляемых пользователю
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// NewProxyHandler создает новый обработчик прокси
func NewProxyHandler(balancer balancerDomain.Strategy) *ProxyHandler {
	return &ProxyHandler{
		balancer: balancer,
		client:   &http.Client{},
	}
}

// ServeHTTP обрабатывает входящие HTTP-запросы
func (h *ProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	const op = "handler.ProxyHandler.ServeHTTP"

	// Получаем следующий доступный бэкенд
	backend := h.balancer.NextBackend()
	if backend == nil {
		slog.Warn("No backend available")
		utils.SendJSON(w,
			http.StatusServiceUnavailable,
			"No backend available",
		)
		return
		return
	}

	// Инкрементируем количество подключений на выбранном бэкенде
	backend.IncrementConnections()
	defer backend.DecrementConnections() // Декрементируем количество подключений после завершения запроса

	// Создаем URL для бэкенда
	targetURL, err := url.Parse(backend.URL)
	if err != nil {
		slog.Error(op,
			"invalid backend URL",
			slog.String("backend", backend.URL),
			slog.String("error", err.Error()),
		)
		utils.SendJSON(w,
			http.StatusInternalServerError,
			"invalid backend URL",
		)
		return
		return
	}

	// Создаем новый HTTP-запрос для проксирования
	proxyReq, err := http.NewRequest(r.Method, targetURL.String()+r.RequestURI, r.Body)
	if err != nil {
		slog.Error(op,
			"failed to create proxy request",
			slog.String("backend", backend.URL),
			slog.String("error", err.Error()),
		)
		utils.SendJSON(w,
			http.StatusInternalServerError,
			"Failed to create proxy request",
		)
		return
	}
	proxyReq.Header = r.Header.Clone()

	// Отправляем запрос на бэкенд
	proxyResp, err := h.client.Do(proxyReq)
	if err != nil {
		slog.Error(op,
			"failed to proxy request",
			slog.String("backend", backend.URL),
			slog.String("error", err.Error()),
		)
		utils.SendJSON(w,
			http.StatusBadGateway,
			"Backend request failed",
		)
		h.balancer.MarkBackendDown(backend)

		return
	}
	defer proxyResp.Body.Close()

	// Копируем заголовки и тело ответа от бэкенда в клиентский ответ
	for k, values := range proxyResp.Header {
		for _, v := range values {
			w.Header().Add(k, v)
		}
	}
	w.WriteHeader(proxyResp.StatusCode)
	io.Copy(w, proxyResp.Body)

	// Логируем успешный прокси запрос
	slog.Info("proxying request",
		slog.String("method", r.Method),
		slog.String("path", r.URL.Path),
		slog.String("backend", backend.URL),
	)
}
