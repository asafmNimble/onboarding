package mongo

import (
	"onboarding/common/data/entities"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

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

func (mn *MongoNumber) AddNum(n *entities.Number) (primitive.ObjectID, error) {
	return primitive.NewObjectID(), nil
}