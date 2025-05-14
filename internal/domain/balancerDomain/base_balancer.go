package balancerDomain

import "sync"

// BaseBalancer предоставляет базовую функциональность для всех стратегий балансировки
type BaseBalancer struct {
	backends []*Backend
	mu       sync.RWMutex
}

// NewBaseBalancer создает новый базовый балансировщик
func NewBaseBalancer(backends []*Backend) *BaseBalancer {
	return &BaseBalancer{
		backends: backends,
	}
}

// GetBackends возвращает список всех бэкендов
func (b *BaseBalancer) GetBackends() []*Backend {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.backends
}

// UpdateBackends обновляет список бэкендов
func (b *BaseBalancer) UpdateBackends(backends []*Backend) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.backends = backends
}

// GetAvailableBackends возвращает список доступных бэкендов
func (b *BaseBalancer) GetAvailableBackends() []*Backend {
	b.mu.RLock()
	defer b.mu.RUnlock()

	var available []*Backend
	for _, backend := range b.backends {
		if backend.IsAlive() {
			available = append(available, backend)
		}
	}
	return available
}

// MarkBackendDown помечает бэкенд как недоступный
func (b *BaseBalancer) MarkBackendDown(backend *Backend) {
	backend.SetAlive(false)
}

// MarkBackendUp помечает бэкенд как доступный
func (b *BaseBalancer) MarkBackendUp(backend *Backend) {
	backend.SetAlive(true)
}
