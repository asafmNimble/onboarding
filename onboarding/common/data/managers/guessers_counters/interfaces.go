package guessers_counters

type DBBackend interface {
	CreateGuessersCounter(guesserID int64) error
	IncreaseGuesserCounter(guesserID int64) error
	GetGuesserCounter(guesserID int64) (int64, error)
}

type GuessersManager interface {
	CreateGuessersCounter(guesserID int64) error
	IncreaseGuesserCounter(guesserID int64) error
	GetGuesserCounter(guesserID int64) (int64, error)
}
