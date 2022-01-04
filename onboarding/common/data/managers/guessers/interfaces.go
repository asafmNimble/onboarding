package guessers

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"onboarding/common/data/entities"
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
