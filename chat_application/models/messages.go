package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Message struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Content   string             `bson:"content"`
	SenderID  primitive.ObjectID `bson:"sender_id"` // ID of the
	RoomID    primitive.ObjectID `bson:"room_id"`
	Timestamp time.Time          `bson:"timestamp"`
}

type MessageJson struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}
