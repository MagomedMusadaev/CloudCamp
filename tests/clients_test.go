package tests

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func TestClientEndpoints(t *testing.T) {
	// Создаем тестовый сервер
	server := setupTestServerWithMocksForTest(t)
	go server.Run()
	defer server.Shutdown()

	if err := waitForServerReady("http://localhost:8080", 2*time.Second); err != nil {
		t.Fatalf("Server did not start in time: %v", err)
	}

	// Тест создания клиента
	t.Run("Create Client", func(t *testing.T) {
		clientData := map[string]interface{}{
			"client_id": "test_client",
			"rate":      20,
			"period":    "1m",
		}
		body, _ := json.Marshal(clientData)

		resp, err := http.Post("http://localhost:8080/clients", "application/json", bytes.NewBuffer(body))
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		resp.Body.Close()

		// Проверяем, что лимит установлен
		limiter := server.GetLimiter()
		rate, per := limiter.GetLimit("test_client")
		assert.Equal(t, 20, rate)
		assert.Equal(t, time.Minute, per)
	})

	// Тест удаления клиента
	t.Run("Delete Client", func(t *testing.T) {
		client := &http.Client{}
		req, _ := http.NewRequest(http.MethodDelete, "http://localhost:8080/clients?client_id=test_client", nil)
		resp, err := client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		resp.Body.Close()

		// Проверяем, что лимит удален
		limiter := server.GetLimiter()
		rate, _ := limiter.GetLimit("test_client")
		assert.Equal(t, 0, rate)
	})
}
