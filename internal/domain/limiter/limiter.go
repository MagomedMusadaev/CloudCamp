package limiter

import "time"

// RateLimiter определяет интерфейс для ограничения частоты запросов
type RateLimiter interface {
	Allow(key string) bool                                  // проверяет, можно ли пропустить запрос для данного ключа (true - можно пропустить)
	GetLimit(key string) (int, time.Duration)               //  возвращает текущий лимит для ключа
	SetLimit(key string, rate int, per time.Duration) error // устанавливает лимит для ключа
}
