package numbers

import (
	"onboarding/common/data/entities"
)

type DBBackend interface {
	AddNum(n *entities.Number) (string, error)
	RemoveNum(n int64) (bool, error)
	QueryNumber(n int64) (int64, *[]entities.GuessType, error)
	Get(n int64) (*entities.Number, error)
}

type NumbersManager interface {
	AddNum(n *entities.Number) (string, error)
	RemoveNum(n int64) (bool, error)
	QueryNumber(n int64) (int64, *[]entities.GuessType, error)
	Get(n int64) (*entities.Number, error)
}
