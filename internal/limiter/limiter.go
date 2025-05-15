package limiter

import (
	"sync"
	"time"
)

// TokenBucket реализует алгоритм Token Bucket для rate limiting
type TokenBucket struct {
	rate   int           // количество токенов
	per    time.Duration // интервал пополнения токенов
	tokens int           // текущее количество токенов
	last   time.Time     // время последнего пополнения
}

// MemoryRateLimiter реализует RateLimiter с хранением в памяти
type MemoryRateLimiter struct {
	mu      sync.RWMutex
	buckets map[string]*TokenBucket
	clients *ClientSettings
}

// NewMemoryRateLimiter создает новый лимитер с хранением в памяти
func NewMemoryRateLimiter() *MemoryRateLimiter {
	return &MemoryRateLimiter{
		buckets: make(map[string]*TokenBucket),
		clients: NewClientSettings(0, 0),
	}
}

// Allow проверяет, можно ли пропустить запрос
func (m *MemoryRateLimiter) Allow(key string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()

	// Сначала проверяем глобальный лимит
	globalBucket, existsGlobal := m.buckets["global"]
	if existsGlobal {
		elapsed := now.Sub(globalBucket.last)
		tokensToAdd := int(elapsed.Seconds() * float64(globalBucket.rate) / globalBucket.per.Seconds())

		// Обновляем токены глобального бакета
		globalBucket.tokens = minTwoNum(globalBucket.rate, globalBucket.tokens+tokensToAdd)
		globalBucket.last = now

		// Проверяем глобальный лимит
		if globalBucket.tokens <= 0 {
			return false
		}
	}

	// Затем проверяем лимит конкретного клиента
	bucket, exists := m.buckets[key]
	if !exists {
		// Если у клиента нет бакета и глобальный лимит пройден, разрешаем запрос
		if existsGlobal {
			globalBucket.tokens--
		}
		return true
	}

	// Вычисляем, сколько токенов нужно добавить
	elapsed := now.Sub(bucket.last)
	tokensToAdd := int(elapsed.Seconds() * float64(bucket.rate) / bucket.per.Seconds())

	// Обновляем количество токенов клиента
	bucket.tokens = minTwoNum(bucket.rate, bucket.tokens+tokensToAdd)
	bucket.last = now

	// Проверяем, есть ли доступный токен
	if bucket.tokens > 0 {
		bucket.tokens--
		// Если запрос разрешен для клиента, уменьшаем глобальный счетчик
		if existsGlobal {
			globalBucket.tokens--
		}
		return true
	}

	return false
}

// GetLimit возвращает текущий лимит для ключа
func (m *MemoryRateLimiter) GetLimit(key string) (int, time.Duration) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if bucket, exists := m.buckets[key]; exists {
		return bucket.rate, bucket.per
	}
	return 0, 0
}

// SetLimit устанавливает лимит для ключа
func (m *MemoryRateLimiter) SetLimit(key string, rate int, per time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Если это глобальный ключ, обновляем дефолтные значения в клиентских настройках
	if key == "global" {
		m.clients = NewClientSettings(rate, per)
	}

	// Получаем актуальные настройки для ключа
	actualRate, actualPer := m.clients.GetSettings(key)

	m.buckets[key] = &TokenBucket{
		rate:   actualRate,
		per:    actualPer,
		tokens: actualRate,
		last:   time.Now(),
	}

	return nil
}

// SetClientLimit устанавливает индивидуальный лимит для клиента
func (m *MemoryRateLimiter) SetClientLimit(clientID string, rate int, per time.Duration) error {
	m.clients.SetSettings(clientID, rate, per)
	return m.SetLimit(clientID, rate, per)
}

// RemoveClientLimit удаляет индивидуальный лимит клиента
func (m *MemoryRateLimiter) RemoveClientLimit(clientID string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Удаляем настройки клиента
	m.clients.RemoveSettings(clientID)
	// Удаляем bucket клиента
	delete(m.buckets, clientID)
}

// minTwoNum возвращает минимальное из двух чисел
func minTwoNum(a, b int) int {
	if a < b {
		return a
	}
	return b
}
