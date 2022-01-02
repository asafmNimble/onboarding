package guessers

import (
	"onboarding/common/data/entities"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DBBackend interface {
	AddGuesser(g *entities.Guesser) (primitive.ObjectID, error)
	RemoveGuesser(g *entities.Guesser) (primitive.ObjectID, error)
	QueryGuesser(g *entities.Guesser) (primitive.ObjectID, error)
}

type GuessersManager interface {
	AddGuesser(g *entities.Guesser) (primitive.ObjectID, error)
	RemoveGuesser(g *entities.Guesser) (primitive.ObjectID, error)
	QueryGuesser(g *entities.Guesser) (primitive.ObjectID, error)
}
