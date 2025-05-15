package utils

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse — структура для стандартного ответа об ошибке
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// SendJSON отправляет JSON-ответ
func SendJSON(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	resp := ErrorResponse{
		Code:    code,
		Message: message,
	}

	_ = json.NewEncoder(w).Encode(resp)
}
