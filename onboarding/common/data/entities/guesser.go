package entities

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Guesser struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	GuesserID   int64              `bson:"guesser_id,omitempty"`
	BeginAt     int64              `bson:"begin_at"`
	IncrementBy int64              `bson:"increment_by"`
	Sleep       int64              `bson:"sleep"`
	Active      bool               `bson:"active"`
	GuessesMade []Guess            `bson:"guesses_made"`
}

type Guess struct {
	GuessNum  int64
	GuessedAt time.Time
}

func NewGuesser(guesserID int64, beginAt int64, incrementBy int64, sleep int64) (*Guesser, error) {
	var guesses []Guess
	guesser := &Guesser{
		ID:          primitive.NewObjectID(),
		GuesserID:   guesserID,
		BeginAt:     beginAt,
		IncrementBy: incrementBy,
		Sleep:       sleep,
		Active:      true,
		GuessesMade: guesses,
	}
	return guesser, nil
}
