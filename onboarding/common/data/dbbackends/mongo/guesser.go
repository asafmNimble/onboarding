package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"onboarding/common/data/entities"
	"time"
)

// Mongo Backend

type MongoGuesser struct {
	dbCollection *mongo.Collection
	resourceName string
}

func NewMongoGuesser(dbc DBConnector) *MongoGuesser {
	db := dbc.GetDB()
	guesser := &MongoGuesser{
		dbCollection: db.Collection("guessers"),
		resourceName: "Guesser",
	}
	return guesser
}

func (mg *MongoGuesser) AddGuesser(g *entities.Guesser) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(3)*time.Second)
	defer cancel()
	_, err := mg.dbCollection.InsertOne(ctx, bson.D{{"_id", g.ID}, {"guesser_id", g.GuesserID}, {"begin_at", g.BeginAt},
		{"increment_by", g.IncrementBy}, {"sleep", g.Sleep}, {"active", g.Active},
		{"guesses_made", g.GuessesMade}})
	if err != nil {
		return g.ID.Hex(), err
	}
	return g.ID.Hex(), nil
}

func (mn *MongoGuesser) RemoveGuesser(guesserID int64) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(3)*time.Second)
	defer cancel()
	_, err := mn.dbCollection.DeleteOne(ctx, bson.D{{"guesser_id", guesserID}})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (mg *MongoGuesser) QueryGuesser(guesserID int64) (string, *[]entities.Guess, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(3)*time.Second)
	defer cancel()
	var guesser entities.Guesser
	//err := mg.dbCollection.FindOne(ctx, bson.D{{"ID", g.GuesserID}, {"guesses", g.GuessesMade}, {"active", g.Active}}).Decode(&guesser)
	err := mg.dbCollection.FindOne(ctx, bson.D{{"guesser_id", guesserID}}).Decode(&guesser)
	guessesMade := &guesser.GuessesMade
	if err != nil {
		return guesser.ID.Hex(), nil, false, err
	}
	return guesser.ID.Hex(), guessesMade, guesser.Active, nil
}

func (mg *MongoGuesser) GetGuesser(guesserID int64) (*entities.Guesser, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(3)*time.Second)
	defer cancel()
	var guesser entities.Guesser
	//err := mg.dbCollection.FindOne(ctx, bson.D{{"ID", g.GuesserID}, {"guesses", g.GuessesMade}, {"active", g.Active}}).Decode(&guesser)
	err := mg.dbCollection.FindOne(ctx, bson.D{{"guesser_id", guesserID}}).Decode(&guesser)
	if err != nil {
		return nil, err
	}
	return &guesser, nil

}
