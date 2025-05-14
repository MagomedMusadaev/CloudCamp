package balancer

import "CloudCamp/internal/domain/balancerDomain"

type LeastConnectionsBalancer struct {
	*balancerDomain.BaseBalancer
}

// NewLeastConnectionsBalancer создает новый балансировщик с алгоритмом Least Connections
func NewLeastConnectionsBalancer(backends []*balancerDomain.Backend) *LeastConnectionsBalancer {
	return &LeastConnectionsBalancer{
		BaseBalancer: balancerDomain.NewBaseBalancer(backends),
	}
}

// NextBackend выбирает бэкенд с наименьшим числом активных соединений
func (l *LeastConnectionsBalancer) NextBackend() *balancerDomain.Backend {
	available := l.GetAvailableBackends()
	if len(available) == 0 {
		return nil
	}

	var selected *balancerDomain.Backend
	minConnections := int64(^uint64(0) >> 1)

	for _, b := range available {
		if cons := b.GetActiveConnections(); cons < minConnections {
			minConnections = cons
			selected = b
		}
	}

	if selected == nil {
		return nil
	}

	return selected
}

// MarkBackendDown помечает бэкенд как недоступный
func (l *LeastConnectionsBalancer) MarkBackendDown(backend *balancerDomain.Backend) {
	l.BaseBalancer.MarkBackendDown(backend)
}

// MarkBackendUp помечает бэкенд как доступный
func (l *LeastConnectionsBalancer) MarkBackendUp(backend *balancerDomain.Backend) {
	l.BaseBalancer.MarkBackendUp(backend)
}

// UpdateBackends обновляет список бэкендов
func (l *LeastConnectionsBalancer) UpdateBackends(backends []*balancerDomain.Backend) {
	l.BaseBalancer.UpdateBackends(backends)
}
