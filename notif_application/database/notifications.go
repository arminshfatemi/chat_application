package database

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

// Message struct to unmarshal sent data from producer
type Message struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Content   string             `bson:"content"`
	SenderID  primitive.ObjectID `bson:"sender_id"`
	RoomID    primitive.ObjectID `bson:"room_id"`
	Timestamp time.Time          `bson:"timestamp"`
}

// NotificationsMessage struct for saving notifications in database
type NotificationsMessage struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Content    string             `bson:"content"`
	SenderID   primitive.ObjectID `bson:"sender_id"`
	RoomID     primitive.ObjectID `bson:"room_id"`
	ReceiverID primitive.ObjectID `bson:"receiver_id"` // determine who is notification for
	Seen       bool               `bson:"seen"`        // does the receiver seen the notification or no
	Timestamp  time.Time          `bson:"timestamp"`
}

// CreateNewNotificationsMessage responsible to create a new NotificationsMessage
func CreateNewNotificationsMessage(message *Message, receiverID primitive.ObjectID) *NotificationsMessage {
	return &NotificationsMessage{
		Content:    message.Content,
		SenderID:   message.SenderID,
		RoomID:     message.RoomID,
		ReceiverID: receiverID,
		Seen:       false,
		Timestamp:  time.Now(),
	}
}

// InsertNotificationInDatabase is for Inserting the NotificationsMessage for the given client in the database
func (nm *NotificationsMessage) InsertNotificationInDatabase(mongodb *mongo.Client) error {
	_, err := mongodb.Database("chat_app").Collection("notifications").InsertOne(context.Background(), nm)
	if err != nil {
		return err
	}
	return nil
}
