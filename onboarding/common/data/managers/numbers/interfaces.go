package numbers

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"onboarding/common/data/entities"
)

type DBBackend interface {
	AddNum(n *entities.Number) (string, error)
	RemoveNum(n int64) (bool, error)
	QueryNumber(n int64) (primitive.ObjectID, *entities.Number, error)
	Get(n int64) (*entities.Number, error)
}

type NumbersManager interface {
	AddNum(n *entities.Number) (string, error)
	RemoveNum(n int64) (bool, error)
	QueryNumber(n int64) (primitive.ObjectID, *entities.Number, error)
	Get(n int64) (*entities.Number, error)
}