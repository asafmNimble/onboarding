package guessers

import (
	"onboarding/common/data/entities"
)

type DBBackend interface {
	AddGuesser(g *entities.Guesser) (string, error)
	RemoveGuesser(guesserID int64) (bool, error)
	QueryGuesser(guesserID int64) (string, *[]entities.Guess, bool, error)
	GetGuesser(guesserID int64) (*entities.Guesser, error)
	UpdateGuessedNumForGuesser(guesserID int64, guess *entities.Guess) (string, error)
}

type GuessersManager interface {
	AddGuesser(g *entities.Guesser) (string, error)
	RemoveGuesser(guesserID int64) (bool, error)
	QueryGuesser(guesserID int64) (string, *[]entities.Guess, bool, error)
	GetGuesser(guesserID int64) (*entities.Guesser, error)
	UpdateGuessedNumForGuesser(guesserID int64, guess *entities.Guess) (string, error)
}
