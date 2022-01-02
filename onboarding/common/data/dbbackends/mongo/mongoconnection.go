package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"log"
	"sync"
	"time"
)

type DBConnector interface {
	GetDB() *mongo.Database
}

type MongoConnector struct{}

var once sync.Once
var singletonClient *mongo.Client

func NewMongoConnector() *MongoConnector {
	connector := &MongoConnector{}
	return connector
}

func (*MongoConnector) connect() (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10)*time.Second)
	defer cancel()

	// explicit default client level options
	opts := options.Client().ApplyURI("10s")
	opts = opts.SetReadConcern(readconcern.Majority()).SetWriteConcern(writeconcern.New(writeconcern.WMajority()))

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (m *MongoConnector) GetDB() *mongo.Database {
	// TODO: support multi-db?
	once.Do(func() {
		client, err := m.connect()
		if err != nil {
			log.Panicf("error connecting to DB, panicking. error: %v", err)
		}
		singletonClient = client
	})
	return singletonClient.Database("AwesomeProject1")
}
