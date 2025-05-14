package balancer

import (
	"CloudCamp/internal/domain/balancerDomain"
	"math/rand"
	"time"
)

type RandomBalancer struct {
	*balancerDomain.BaseBalancer
	rnd *rand.Rand
}

// NewRandomBalancer создает новый балансировщик с алгоритмом Random)
func NewRandomBalancer(backends []*balancerDomain.Backend) *RandomBalancer {
	return &RandomBalancer{
		BaseBalancer: balancerDomain.NewBaseBalancer(backends),
		rnd:          rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// NextBackend выбирает случайный доступный бэкенд
func (r *RandomBalancer) NextBackend() *balancerDomain.Backend {
	available := r.GetAvailableBackends()
	if len(available) == 0 {
		return nil
	}

	idx := r.rnd.Intn(len(available))
	return available[idx]
}

// MarkBackendDown помечает бэкенд как недоступный
func (r *RandomBalancer) MarkBackendDown(backend *balancerDomain.Backend) {
	r.BaseBalancer.MarkBackendDown(backend)
}

// MarkBackendUp помечает бэкенд как доступный
func (r *RandomBalancer) MarkBackendUp(backend *balancerDomain.Backend) {
	r.BaseBalancer.MarkBackendUp(backend)
}

// UpdateBackends обновляет список доступных бэкендов
func (r *RandomBalancer) UpdateBackends(backends []*balancerDomain.Backend) {
	r.BaseBalancer.UpdateBackends(backends)
}
