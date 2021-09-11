package db

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"os"
)

var client *mongo.Client
var db *mongo.Collection
var mongoCtx context.Context

func GetContext() context.Context {
	return mongoCtx
}

func Disconnect() {
	log.Println("disconnecting database ...")
	client.Disconnect(mongoCtx)
}

func Connect() error {
	log.Println("starting database connection ...")
	client, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("DBCN")))
	if err != nil {
		return err
	}

	mongoCtx = context.Background()
	err = client.Connect(mongoCtx)
	if err != nil {
		return err
	}

	err = client.Ping(mongoCtx, readpref.Primary())
	if err != nil {
		return err
	}

	db = client.Database("film-voting").Collection("rank")
	log.Println("database successfully connected !")
	return nil
}
