package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Connect() *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGODB")))
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println("Connected to database")
	return client

}

var Client *mongo.Client = Connect()

func OpenCollection(client *mongo.Client, collname string) *mongo.Collection {
	var collection *mongo.Collection = Client.Database("cluster0").Collection(collname)
	return collection
}
