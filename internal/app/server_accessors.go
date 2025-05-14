package app

import (
	"CloudCamp/internal/domain/balancerDomain"
	"CloudCamp/internal/limiter"
)

// GetLimiter возвращает экземпляр лимитера
func (s *Server) GetLimiter() *limiter.MemoryRateLimiter {
	return s.limiter
}

// GetBackends возвращает список бэкендов
func (s *Server) GetBackends() []*balancerDomain.Backend {
	return s.balancer.GetBackends()
}
