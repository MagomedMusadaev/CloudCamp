package balancerDomain

// Strategy определяет интерфейс для алгоритмов балансировки нагрузки
type Strategy interface {
	NextBackend() *Backend              // возвращает следующий доступный бэкенд
	MarkBackendDown(backend *Backend)   // помечает бэкенд как недоступный
	MarkBackendUp(backend *Backend)     // помечает бэкенд как доступный
	UpdateBackends(backends []*Backend) // обновляет список доступных бэкендов
	GetBackends() []*Backend            // возвращает список всех бэкендов
}
