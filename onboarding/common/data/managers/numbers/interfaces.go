package numbers

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"onboarding/common/data/entities"
)

type DBBackend interface {
	AddNum(n *entities.Number) (primitive.ObjectID, error)
	RemoveNum(n *entities.Number) (primitive.ObjectID, error)
	QueryNumber(n int64) (primitive.ObjectID, *entities.Number, error)
}

type NumbersManager interface {
	AddNum(n *entities.Number) (primitive.ObjectID, error)
	RemoveNum(n *entities.Number) (primitive.ObjectID, error)
	QueryNumber(n int64) (primitive.ObjectID, *entities.Number, error)
}