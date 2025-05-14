package balancerDomain

import "sync/atomic"

// Backend представляет собой отдельный сервер в пуле балансировки
type Backend struct {
	URL               string       // адрес бэкенд сервера
	Alive             atomic.Bool  // показывает, доступен ли сервер
	ActiveConnections atomic.Int64 // текущее количество активных соединений
}

func NewBackend(url string) *Backend {
	b := &Backend{
		URL: url,
	}
	b.Alive.Store(true)

	return b
}

// IsAlive проверяет, доступен ли бэкенд
func (b *Backend) IsAlive() bool {
	return b.Alive.Load()
}

// SetAlive устанавливает статус доступности бэкенда
func (b *Backend) SetAlive(alive bool) {
	b.Alive.Store(alive)
}

// IncrementConnections увеличивает счетчик активных соединений
func (b *Backend) IncrementConnections() {
	b.ActiveConnections.Add(1)
}

// DecrementConnections уменьшает счетчик активных соединений
func (b *Backend) DecrementConnections() {
	b.ActiveConnections.Add(-1)
}

// GetActiveConnections возвращает количество активных соединений
func (b *Backend) GetActiveConnections() int64 {
	return b.ActiveConnections.Load()
}
