package api

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

const connection_timeout = 10 * time.Second

func (app *application) ConnectMongoDB() *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), connection_timeout)
	defer cancel()

	client_options := options.Client()

	client_options.
		ApplyURI("mongodb://localhost:27017").
		SetMaxPoolSize(5).
		SetMinPoolSize(1)

	client, err := mongo.Connect(ctx, client_options)

	if err != nil {
		log.Error("error while connecting to mongodb: ", err)
		return nil
	}

	return client
}

func (app *application) DisconnectMongoDB() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Disconnect(ctx); err != nil {
		log.Error("error occured while disconnecting from db: ", err)
		panic(err)
	}

	log.Info("mongodb disconnected successfully")
}
