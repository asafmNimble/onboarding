package numbers

import (
	"onboarding/common/data/entities"
)

type Manager struct {
	backend DBBackend
}

func NewManager(b DBBackend) *Manager {
	return &Manager{
		backend: b,
	}
}

func (m *Manager) AddNum(num int64) (string, error) {
	number, err := entities.NewNumber(num)
	if err != nil {
		return "", err
	}
	numID, err := m.backend.AddNum(number)
	if err != nil {
		return "", err
	}
	return numID, err
}

func (m *Manager) QueryNumber(num int64) (int64, *[]entities.GuessType, error) {
	number, numGuesses, err := m.backend.QueryNumber(num)
	if err != nil {
		return number, nil, err
	}
	return number, numGuesses, err
}

func (m *Manager) RemoveNum(num int64) (bool, error) {
	_, err := m.backend.GetNumber(num)
	if err != nil {
		return false, err
	}
	m.backend.RemoveNum(num)
	return true, err
}

func (m *Manager) GetNumber(num int64) (*entities.Number, error) {
	number, err := m.backend.GetNumber(num)
	if err != nil {
		return nil, err
	}
	return number, nil
}

func (m *Manager) UpdateGuessForNumber(num int64, guess *entities.GuessType) (string, error) {
	id, err := m.backend.UpdateGuessForNumber(num, guess)
	if err != nil {
		return "", err
	}
	return id, nil
}
