package balancer

import (
	"CloudCamp/internal/config"
	"CloudCamp/internal/domain/balancerDomain"
)

// NewBalancerFactory создает новый экземпляр балансировщика с указанной стратегией
func NewBalancerFactory(cfg *config.Config) balancerDomain.Strategy {

	var backends []*balancerDomain.Backend

	for _, url := range cfg.Balancer.Backends {
		backends = append(backends, balancerDomain.NewBackend(url))
	}

	switch cfg.Balancer.Strategy {
	case "round-robin":
		return NewRoundRobinBalancer(backends)
	case "random":
		return NewRandomBalancer(backends)
	case "least-connections":
		return NewLeastConnectionsBalancer(backends)
	default:
		return nil
	}
}
