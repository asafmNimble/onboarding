package guessers

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

func (m *Manager) AddGuesser(guesserID int64, beginAt int64, incrementBy int64, sleep int64) (string, error) {
	guesser, err := entities.NewGuesser(guesserID, beginAt, incrementBy, sleep)
	if err != nil {
		return "", err
	}
	ID, err := m.backend.AddGuesser(guesser)
	if err != nil {
		return "", err
	}
	return ID, err
}

func (m *Manager) QueryGuesser(guesserID int64) (string, *[]entities.Guess, bool, error) {
	guesser, guesses, active, err := m.backend.QueryGuesser(guesserID)
	if err != nil {
		return guesser, nil, false, err
	}
	return guesser, guesses, active, err
}

func (m *Manager) RemoveGuesser(guesserID int64) (bool, error) {
	_, err := m.backend.GetGuesser(guesserID)
	if err != nil {
		return false, err
	}
	m.backend.RemoveGuesser(guesserID)
	return true, err
}

func (m *Manager) GetGuesser(guesserID int64) (*entities.Guesser, error) {
	guesser, err := m.backend.GetGuesser(guesserID)
	if err != nil {
		return nil, err
	}
	return guesser, nil
}

