package models

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
)

// DatabaseInit is responsible to connect to the mongoDB and returning client variable
func DatabaseInit() (*mongo.Client, error) {
	// Get MongoDB URI and Database name from environment variables
	mongoURI := os.Getenv("MONGO_URI")

	// MongoDB client options
	clientOptions := options.Client().ApplyURI(mongoURI)

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}
	//defer func() {
	//	err := client.Disconnect(context.TODO())
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//}()

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, err
	}

	log.Println("Connected to MongoDB!")
	return client, nil
}
