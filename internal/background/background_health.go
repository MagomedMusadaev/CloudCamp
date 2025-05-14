package background

import (
	"CloudCamp/internal/domain/balancerDomain"
	"context"
	"log/slog"
	"net/http"
	"sync"
	"time"
)

// HealthChecker периодически проверяет доступность бэкендов
type HealthChecker struct {
	backends []*balancerDomain.Backend
	ticker   *time.Ticker
	client   *http.Client
	path     string
	wg       sync.WaitGroup
}

// NewHealthChecker создает новый HealthChecker
func NewHealthChecker(backends []*balancerDomain.Backend, interval time.Duration, path string) *HealthChecker {
	return &HealthChecker{
		backends: backends,
		ticker:   time.NewTicker(interval),
		path:     path,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// Start запускает процесс проверки здоровья
func (hc *HealthChecker) Start(ctx context.Context) {
	hc.wg.Add(1)
	go func() {
		defer hc.wg.Done()
		for {
			select {
			case <-ctx.Done():
				hc.ticker.Stop()
				return
			case <-hc.ticker.C:
				slog.Info("Checking backends...")
				hc.checkBackends()
			}
		}
	}()
}

// Wait ожидает завершения работы
func (hc *HealthChecker) Wait() {
	hc.wg.Wait()
}

// checkBackends проверяет доступность всех бэкендов
func (hc *HealthChecker) checkBackends() {
	for _, backend := range hc.backends {
		go func(b *balancerDomain.Backend) {
			resp, err := hc.client.Get(b.URL + hc.path)
			if err != nil {
				slog.Warn("Backend is unavailable", "backend", b.URL+hc.path)
				b.SetAlive(false)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode >= 200 && resp.StatusCode < 300 {
				b.SetAlive(true)
			} else {
				slog.Warn("Backend returned non-2xx status", "backend", b.URL, "status", resp.StatusCode)
				b.SetAlive(false)
			}
		}(backend)
	}
}
