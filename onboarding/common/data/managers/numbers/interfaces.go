package numbers

import (
	"onboarding/common/data/entities"
)

type DBBackend interface {
	AddNum(n *entities.Number) (string, error)
	RemoveNum(n int64) (bool, error)
	QueryNumber(n int64) (int64, *[]entities.GuessType, error)
	GetNumber(n int64) (*entities.Number, error)
	UpdateGuessForNumber(n int64, guess *entities.GuessType) (string, error)
}

type NumbersManager interface {
	AddNum(n *entities.Number) (string, error)
	RemoveNum(n int64) (bool, error)
	QueryNumber(n int64) (int64, *[]entities.GuessType, error)
	GetNumber(n int64) (*entities.Number, error)
	UpdateGuessForNumber(n int64, guess *entities.GuessType) (string, error)
}
