package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"onboarding/common/data/entities"
	"time"
)

// Mongo Backend

type MongoNumber struct {
	dbCollection *mongo.Collection
	resourceName string
}

func NewMongoNumber(dbc DBConnector) *MongoNumber {
	db := dbc.GetDB()
	number := &MongoNumber{
		dbCollection: db.Collection("numbers"),
		resourceName: "Number",
	}
	return number
}

func (mn *MongoNumber) AddNum(n *entities.Number) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(3)*time.Second)
	defer cancel()
	_, err := mn.dbCollection.InsertOne(ctx, bson.D{{"_id", n.ID}, {"number", n.Number}, {"guesses", n.Guesses}})
	if err != nil {
		return n.ID.Hex(), err
	}
	return n.ID.Hex(), nil
}

func (mn *MongoNumber) QueryNumber(num int64) (int64, *[]entities.GuessType, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(3)*time.Second)
	defer cancel()
	var number entities.Number
	err := mn.dbCollection.FindOne(ctx, bson.D{{"number", num}}).Decode(&number)
	guesses := number.Guesses
	if err != nil {
		return num, nil, err
	}
	return num, &guesses, nil
}

func (mn *MongoNumber) RemoveNum(num int64) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(3)*time.Second)
	defer cancel()
	_, err := mn.dbCollection.DeleteOne(ctx, bson.D{{"number", num}})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (mn *MongoNumber) GetNumber(num int64) (*entities.Number, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(3)*time.Second)
	defer cancel()
	var number entities.Number
	err := mn.dbCollection.FindOne(ctx, bson.D{{"number", num}}).Decode(&number)
	if err != nil {
		return nil, err
	}
	return &number, nil
}

func (mn *MongoNumber) UpdateGuessForNumber(num int64, guess *entities.GuessType) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(3)*time.Second)
	defer cancel()
	currNum, err := mn.GetNumber(num)
	if err != nil {
		return "", err
	}
	guesses := currNum.Guesses
	guesses = append(guesses, *guess)
	_, err =  mn.dbCollection.UpdateOne(ctx, bson.D{{"_id", currNum.ID}}, bson.D{{"$set", bson.D{{"guesses", guesses}}}})
	if err != nil {
		return currNum.ID.Hex(), err
	}
	return currNum.ID.Hex(), nil
}
