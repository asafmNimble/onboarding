package entities

type GuesserCounter struct {
	GuesserID int64              `bson:"guesser_id"`
	Counter   int64              `bson:"counter"`
}

func NewGuesserCounter(guesserID int64) *GuesserCounter {
	return &GuesserCounter{
		GuesserID: guesserID,
		Counter:   1,
	}
}
