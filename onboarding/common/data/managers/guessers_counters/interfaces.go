package guessers_counters

import "onboarding/common/data/entities"

type DBBackend interface {
	CreateGuessersCounter(gc *entities.GuesserCounter) error
	IncreaseGuesserCounter(guesserID int64) error
	GetGuesserCounter(guesserID int64) (int64, error)
}

type GuessersManager interface {
	CreateGuessersCounter(guesserID int64) error
	IncreaseGuesserCounter(guesserID int64) error
	GetGuesserCounter(guesserID int64) (int64, error)
}
