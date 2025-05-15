package limiter

import (
	"log/slog"
	"time"
)

// RefillAll пополняет токены во всех бакетах
func (m *MemoryRateLimiter) RefillAll() {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	for _, bucket := range m.buckets {
		slog.Debug("Refill bucket", bucket) // TODO: для мониторинога работы токенов и refill
		if bucket == nil {
			continue
		}

		elapsed := now.Sub(bucket.last)
		tokensToAdd := int(elapsed.Seconds() * float64(bucket.rate) / bucket.per.Seconds())

		bucket.tokens = minTwoNum(bucket.rate, bucket.tokens+tokensToAdd)
		bucket.last = now
	}
}
