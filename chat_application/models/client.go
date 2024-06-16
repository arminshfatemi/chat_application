package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Client struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	Name         string             `bson:"name"`
	Email        string             `bson:"email"`
	PasswordHash string             `bson:"password_hash"`
	CreatedAt    time.Time          `bson:"created_at"`
}
