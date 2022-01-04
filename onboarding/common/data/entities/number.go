package entities

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Number struct {
	ID      primitive.ObjectID `bson:"_id"`
	Number  int64              `bson:"number,omitempty"`
	Guesses []GuessType `bson:"guesses,omitempty"`
}

type GuessType struct {
	FoundBy string    `bson:"found_by,omitempty"` //guesserID
	FoundAt time.Time `bson:"found_at,omitempty"`
}

func NewNumber(num int64) (*Number, error) {
	number := &Number{
		ID:      primitive.NewObjectID(),
		Number:  num,
		Guesses: nil,
	}
	return number, nil
}
