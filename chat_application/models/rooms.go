package models

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type Room struct {
	ID        primitive.ObjectID   `bson:"_id,omitempty"`
	Name      string               `bson:"name"`
	Members   []primitive.ObjectID `bson:"members"`
	CreatedAt time.Time            `bson:"created_at"`
}

// CreateNewRoom creates a new instance of Room with given name and creator in the members
func CreateNewRoom(name string, id primitive.ObjectID) *Room {
	return &Room{
		Name:      name,
		Members:   []primitive.ObjectID{id},
		CreatedAt: time.Now(),
	}
}

// CreateRoomInDatabase insert Room to the database
func (room *Room) CreateRoomInDatabase(mongoClient *mongo.Client) (primitive.ObjectID, error) {
	roomCollection := mongoClient.Database("chat_app").Collection("rooms")
	result, err := roomCollection.InsertOne(context.TODO(), room)
	if err != nil {
		return primitive.ObjectID{}, err
	}
	return result.InsertedID.(primitive.ObjectID), nil
}

// RoomExists give the room with given name
func RoomExists(name string, client *mongo.Client) (*Room, error) {
	collection := client.Database("chat_app").Collection("rooms")

	var room Room
	err := collection.FindOne(context.TODO(), bson.M{"name": name}).Decode(&room)
	if err != nil {
		return &Room{}, err
	}
	return &room, nil
}
