package handler

import (
	"CloudCamp/internal/limiter"
	"CloudCamp/pkg/utils"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"
)

// ClientHandler обработчик для управления клиентами
type ClientHandler struct {
	limiter *limiter.MemoryRateLimiter
}

// ClientRequest структура для запроса на создание/обновление клиента
type ClientRequest struct {
	ClientID string `json:"client_id"`
	Rate     int    `json:"rate"`
	Period   string `json:"period"`
}

// NewClientHandler создает новый обработчик для клиентов
func NewClientHandler(limiter *limiter.MemoryRateLimiter) *ClientHandler {
	return &ClientHandler{limiter: limiter}
}

// CreateClient обработчик для создания нового клиента
func (h *ClientHandler) CreateClient(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		slog.Error("Method not allowed", slog.String("method", r.Method))
		utils.SendJSON(w,
			http.StatusMethodNotAllowed,
			"Method not allowed",
		)
		return
	}

	var req ClientRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("Error decoding request", slog.String("error", err.Error()))
		utils.SendJSON(w,
			http.StatusBadRequest,
			"Invalid request body",
		)
		return
	}

	period, err := time.ParseDuration(req.Period)
	if err != nil {
		slog.Error("Error parsing period", slog.String("period", req.Period), slog.String("error", err.Error()))
		utils.SendJSON(w,
			http.StatusBadRequest,
			"Invalid period format. Example: '1s', '500ms', '2m'",
		)
		return
	}

	if req.ClientID == "" || req.Rate <= 0 || period <= 0 {
		slog.Warn("Invalid request", slog.String("client_id", req.ClientID))
		utils.SendJSON(w,
			http.StatusBadRequest,
			"Invalid client parameters",
		)
		return
	}

	// Устанавливаем лимит для клиента
	err = h.limiter.SetClientLimit(req.ClientID, req.Rate, period)
	if err != nil {
		slog.Error("Error setting client limit", slog.String("error", err.Error()))
		utils.SendJSON(w,
			http.StatusInternalServerError,
			"Failed to set client limit",
		)
		return
	}

	utils.SendJSON(w, http.StatusCreated, "Client created successfully")
}

// DeleteClient обработчик для удаления клиента
func (h *ClientHandler) DeleteClient(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		slog.Error("Method not allowed", slog.String("method", r.Method))
		utils.SendJSON(w,
			http.StatusMethodNotAllowed,
			"Method not allowed",
		)
		return
	}

	clientID := r.URL.Query().Get("client_id")
	if clientID == "" {
		slog.Warn("Missing client_id", slog.String("client_id", clientID))
		utils.SendJSON(w,
			http.StatusBadRequest,
			"Client ID is required",
		)
		return
	}

	// Удаляем настройки клиента
	h.limiter.RemoveClientLimit(clientID)
	utils.SendJSON(w, http.StatusOK, "Client deleted successfully")
}
