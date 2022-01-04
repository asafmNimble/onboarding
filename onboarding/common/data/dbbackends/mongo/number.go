package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	_, err := mn.dbCollection.InsertOne(ctx, bson.D{{"_id", n.ID}, {"number", n.Number}, {"guesses", n.Guesses}})  // UpdateOne(ctx, bson.D{{"_id", n.ID}}, bson.D{{"$set", bson.D{{"number", n.Number}, {"active", n.Active}, {"guesses", n.Guesses}}}})
	if err != nil {
		return n.ID.Hex(), err
	}
	return n.ID.Hex(), nil
}

func (mn *MongoNumber) QueryNumber(num int64) (primitive.ObjectID, *entities.Number, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(3)*time.Second)
	defer cancel()
	var res entities.Number
	err := mn.dbCollection.FindOne(ctx, bson.D{{"Number", num}}).Decode(&res)
	if err != nil {
		return res.ID, nil, err
	}
	return res.ID, &res, nil
}


func (mn *MongoNumber) RemoveNum(num int64) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(300)*time.Second)  //TODO: change back to 3 secs
	defer cancel()
	_, err := mn.dbCollection.DeleteOne(ctx, bson.D{{"Number", num}})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (mn *MongoNumber) Get(num int64) (*entities.Number, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(3)*time.Second)
	defer cancel()
	var res entities.Number
	err := mn.dbCollection.FindOne(ctx, bson.D{{"number", num}}).Decode(&res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
