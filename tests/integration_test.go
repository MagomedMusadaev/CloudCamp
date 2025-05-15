package tests

import (
	"CloudCamp/internal/app"
	"CloudCamp/internal/config"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

type MockBackend struct{}

func (m *MockBackend) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func startMockServer(addr string, handler http.Handler) {
	go http.ListenAndServe(addr, handler)
}

// setupTestServerWithMocksForTest — для обычных тестов
func setupTestServerWithMocksForTest(t *testing.T) *app.Server {
	mockBackend := new(MockBackend)

	// Конфиг для сервера с балансировщиком и rate limiting
	cfg := &config.Config{
		Server: config.ServerConfig{Port: 8080},
		Balancer: config.BalancerConfig{
			Strategy: "round-robin",
			Backends: []string{"http://localhost:8081", "http://localhost:8082"},
		},
		RateLimiter: config.RateLimitConfig{
			Enabled: true,
			Rate:    10,
			Period:  time.Second,
		},
	}

	server, err := app.NewServer(cfg)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	// Запуск мок-серверов
	startMockServer(":8081", mockBackend)
	startMockServer(":8082", mockBackend)

	return server
}

// setupTestServerWithMocksForBenchmark — для бенчмарков
func setupTestServerWithMocksForBenchmark(b *testing.B) *app.Server {
	mockBackend := new(MockBackend)

	// Конфиг для сервера с балансировщиком и rate limiting
	cfg := &config.Config{
		Server: config.ServerConfig{Port: 8080},
		Balancer: config.BalancerConfig{
			Strategy: "round-robin",
			Backends: []string{"http://localhost:8081", "http://localhost:8082"},
		},
		RateLimiter: config.RateLimitConfig{
			Enabled: true,
			Rate:    10,
			Period:  time.Second,
		},
	}

	server, err := app.NewServer(cfg)
	if err != nil {
		b.Fatalf("Failed to create server: %v", err)
	}

	// Запуск мок-серверов
	startMockServer(":8081", mockBackend)
	startMockServer(":8082", mockBackend)

	return server
}

// BenchmarkRoundRobin — тест производительности с использованием round-robin балансировки
func BenchmarkRoundRobin(b *testing.B) {
	server := setupTestServerWithMocksForBenchmark(b)
	go server.Run()
	defer server.Shutdown()

	if err := waitForServerReady("http://localhost:8080", 2*time.Second); err != nil {
		b.Fatalf("Server did not start in time: %v", err)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			resp, err := http.Get("http://localhost:8080")
			if err != nil {
				b.Fatalf("Request failed: %v", err)
			}
			if resp.StatusCode != http.StatusOK {
				b.Fatalf("Unexpected status code: %d", resp.StatusCode)
			}
			resp.Body.Close()
		}
	})
}

// waitForServerReady — проверка, что сервер доступен
func waitForServerReady(url string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		resp, err := http.Get(url)
		if err == nil {
			resp.Body.Close()
			return nil
		}
		time.Sleep(50 * time.Millisecond)
	}
	return fmt.Errorf("server not ready at %s", url)
}

// TestRoundRobin — тест для проверки стратегии round-robin
func TestRoundRobin(t *testing.T) {
	server := setupTestServerWithMocksForTest(t)
	go server.Run()
	defer server.Shutdown()

	if err := waitForServerReady("http://localhost:8080", 2*time.Second); err != nil {
		t.Fatalf("Server did not start in time: %v", err)
	}

	client := &http.Client{}
	for i := 0; i < 5; i++ {
		resp, err := client.Get("http://localhost:8080")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		resp.Body.Close()
	}
}
