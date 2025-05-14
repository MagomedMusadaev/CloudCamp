package app

import "CloudCamp/internal/domain/balancerDomain"

// GetBackends возвращает список бэкендов
func (s *Server) GetBackends() []*balancerDomain.Backend {
	return s.balancer.GetBackends()
}
