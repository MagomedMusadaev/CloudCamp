package balancer

import (
	"CloudCamp/internal/domain/balancerDomain"
	"sync/atomic"
)

type RoundRobinBalancer struct {
	*balancerDomain.BaseBalancer
	current atomic.Uint64 // текущий индекс для Round Robin
}

// NewRoundRobinBalancer создает новый балансировщик с алгоритмом Round Robin
func NewRoundRobinBalancer(backends []*balancerDomain.Backend) *RoundRobinBalancer {
	return &RoundRobinBalancer{
		BaseBalancer: balancerDomain.NewBaseBalancer(backends),
	}
}

// NextBackend возвращает следующий доступный бэкенд
func (r *RoundRobinBalancer) NextBackend() *balancerDomain.Backend {
	available := r.GetAvailableBackends()
	if len(available) == 0 {
		return nil
	}

	// Атомарно получаем текущее значение и инкрементируем его
	current := r.current.Load()                    // Получаем текущее значение атомарно
	next := (current + 1) % uint64(len(available)) // Вычисляем следующий индекс

	// Устанавливаем новое значение
	r.current.Store(next) // Атомарно сохраняем новый индекс

	return available[next]
}

// MarkBackendDown помечает бэкенд как недоступный
func (r *RoundRobinBalancer) MarkBackendDown(backend *balancerDomain.Backend) {
	r.BaseBalancer.MarkBackendDown(backend)
}

// MarkBackendUp помечает бэкенд как доступный
func (r *RoundRobinBalancer) MarkBackendUp(backend *balancerDomain.Backend) {
	r.BaseBalancer.MarkBackendUp(backend)
}

// UpdateBackends обновляет список доступных бэкендов
func (r *RoundRobinBalancer) UpdateBackends(backends []*balancerDomain.Backend) {
	r.BaseBalancer.UpdateBackends(backends) // Обновляем базовый список бэкендов
}
