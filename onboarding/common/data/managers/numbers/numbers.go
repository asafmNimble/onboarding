package numbers

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (m *Manager) QueryNumber(num int64) (primitive.ObjectID, *entities.Number, error) {
	numID, numDetails, err := m.backend.QueryNumber(num)
	if err != nil {
		return numID, nil, err
	}
	return numID, numDetails, err
}

func (m *Manager) RemoveNum(num int64) (bool, error) {
	_, err := m.backend.Get(num)
	if err != nil {
		return false, err
	}
	m.backend.RemoveNum(num)
	return true, err
}

func (m *Manager) Get(num int64) (*entities.Number, error) {
	number, err := m.backend.Get(num)
	if err != nil {
		return nil, err
	}
	return number, nil
}
