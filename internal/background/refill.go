package background

import (
	"CloudCamp/internal/limiter"
	"context"
	"log/slog"
	"sync"
	"time"
)

// TokenRefill периодически пополняет токены в бакетах
type TokenRefill struct {
	limiter *limiter.MemoryRateLimiter
	ticker  *time.Ticker
	wg      sync.WaitGroup
}

// NewTokenRefill создает новый TokenRefill
func NewTokenRefill(limiter *limiter.MemoryRateLimiter, interval time.Duration) *TokenRefill {
	return &TokenRefill{
		limiter: limiter,
		ticker:  time.NewTicker(interval),
	}
}

// Start запускает процесс пополнения токенов
func (tr *TokenRefill) Start(ctx context.Context) {
	tr.wg.Add(1)
	go func() {
		defer tr.wg.Done()
		for {
			select {
			case <-ctx.Done():
				tr.ticker.Stop()
				return
			case <-tr.ticker.C:
				slog.Info("Checking token...")
				tr.refillTokens()
			}
		}
	}()
}

// Wait ожидает завершения работы
func (tr *TokenRefill) Wait() {
	tr.wg.Wait()
}

// refillTokens пополняет токены во всех бакетах
func (tr *TokenRefill) refillTokens() {
	tr.limiter.RefillAll()
}
