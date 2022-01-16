package guessers_counters

import "onboarding/common/data/entities"

type Manager struct {
	backend DBBackend
}

func NewManager(b DBBackend) *Manager {
	return &Manager{
		backend: b,
	}
}

func (m *Manager) CreateGuessersCounter(guesserID int64) error {
	gcm := entities.NewGuesserCounter(guesserID)
	return m.backend.CreateGuessersCounter(gcm)
}

func (m *Manager) IncreaseGuesserCounter(guesserID int64) error {
	return m.backend.IncreaseGuesserCounter(guesserID)
}

func (m *Manager) GetGuesserCounter(guesserID int64) (int64, error) {
	count, err := m.backend.GetGuesserCounter(guesserID)
	if err != nil {return count, err}
	return count, nil
}
