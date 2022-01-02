package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"onboarding/common/data/entities"
	"time"
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

func (mn *MongoNumber) AddNum(n *entities.Number) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(3)*time.Second)
	defer cancel()
	_, err := mn.dbCollection.InsertOne(ctx, bson.D{{"_id", n.ID}, {"number", n.Number}, {"active", n.Active}, {"guesses", n.Guesses}})  // UpdateOne(ctx, bson.D{{"_id", n.ID}}, bson.D{{"$set", bson.D{{"number", n.Number}, {"active", n.Active}, {"guesses", n.Guesses}}}})
	if err != nil {
		return n.ID.Hex(), err
	}
	return n.ID.Hex(), nil
}

/*
func (p *MongoPeer) CreateOrUpdate(e *entities.Peer) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(3)*time.Second)
	defer cancel()
	opts := &options.UpdateOptions{}
	opts = opts.SetUpsert(true)
	_, err := p.dbCollection.UpdateOne(ctx, bson.D{{"_id", e.ID}}, bson.D{{"$set", e}}, opts)
	if err != nil {
		return e.ID, p.ResolveError(err, p.resourceName)
	}

	return e.ID, nil
}
*/