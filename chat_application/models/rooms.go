package models

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type Room struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name"`
	CreatedAt time.Time          `bson:"created_at"`
}

// CreateNewRoom creates a new instance of Room
func CreateNewRoom(name string) *Room {
	return &Room{
		Name:      name,
		CreatedAt: time.Now(),
	}
}

// CreateRoomInDatabase insert Room to the database
func (room *Room) CreateRoomInDatabase(mongoClient *mongo.Client) error {
	roomCollection := mongoClient.Database("chat_app").Collection("rooms")
	_, err := roomCollection.InsertOne(context.TODO(), room)
	if err != nil {
		return err
	}
	return nil

}

// RoomExists give the room with given name
func RoomExists(name string, client *mongo.Client) error {
	collection := client.Database("chat_app").Collection("rooms")

	var room Room
	err := collection.FindOne(context.TODO(), bson.M{"name": name}).Decode(&room)
	if err != nil {
		return err
	}
	return nil
}
