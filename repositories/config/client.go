package config

import (
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongoClient(uri string, timeout time.Duration) (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(uri)
	timeoutContext, cancel := TimeoutContext(timeout)
	defer cancel()

	client, err := mongo.Connect(timeoutContext, clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Ping(timeoutContext, nil)
	if err != nil {
		return nil, err
	}

	return client, nil
}
