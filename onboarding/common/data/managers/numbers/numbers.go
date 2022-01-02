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

func (m *Manager) AddNum(num int64) (primitive.ObjectID, *entities.Number, error) {
	number, err := entities.NewNumber(num)
	if err != nil {
		return number.ID, nil, err
	}
	numID, err := m.backend.AddNum(number)
	if err != nil {
		return number.ID, nil, err
	}
	return numID, number, err
}

func (m *Manager) QueryNumber(num int64) (primitive.ObjectID, *entities.Number, error) {
	numID, numDetails, err := m.backend.QueryNumber(num)
	if err != nil {
		return numID, nil, err
	}
	return numID, numDetails, err
}

func (m *Manager) RemoveNum(num int64) (primitive.ObjectID, *entities.Number, error) {
	numID, numDetails, err := m.QueryNumber(num)
	if err != nil {
		return numID, nil, err
	}
	numDetails.Active = false
	return numID, numDetails, err
}
