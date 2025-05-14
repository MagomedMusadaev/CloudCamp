package limiter

import (
	"CloudCamp/internal/config"
	"sync"
	"time"
)

// ClientSettings хранит индивидуальные настройки лимитов для клиентов
type ClientSettings struct {
	mu       sync.RWMutex
	settings map[string]*config.ClientLimit
	defRate  int
	defPer   time.Duration
}

// NewClientSettings создает новый экземпляр ClientSettings
func NewClientSettings(defaultRate int, defaultPer time.Duration) *ClientSettings {
	return &ClientSettings{
		settings: make(map[string]*config.ClientLimit),
		defRate:  defaultRate,
		defPer:   defaultPer,
	}
}

// GetSettings возвращает настройки для клиента или дефолтные значения
func (cs *ClientSettings) GetSettings(clientID string) (int, time.Duration) {
	cs.mu.RLock()
	defer cs.mu.RUnlock()

	if settings, exists := cs.settings[clientID]; exists {
		return settings.RateLimit, settings.Period
	}
	return cs.defRate, cs.defPer
}

// SetSettings устанавливает индивидуальные настройки для клиента
func (cs *ClientSettings) SetSettings(clientID string, rate int, per time.Duration) {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	cs.settings[clientID] = &config.ClientLimit{
		RateLimit: rate,
		Period:    per,
	}
}

// RemoveSettings удаляет индивидуальные настройки клиента
func (cs *ClientSettings) RemoveSettings(clientID string) {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	delete(cs.settings, clientID)
}

// HasCustomSettings проверяет наличие индивидуальных настроек
func (cs *ClientSettings) HasCustomSettings(clientID string) bool {
	cs.mu.RLock()
	defer cs.mu.RUnlock()

	_, exists := cs.settings[clientID]
	return exists
}
