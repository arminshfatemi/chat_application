package apiHandlers

import (
	"chatRoom/models"
	"chatRoom/utils"
	"context"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"time"
)

type RegisterUserRequest struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func ClientSignUpHandler(mongoClient *mongo.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		// validate the request body sent by user
		var req RegisterUserRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error(), "message": "invalid request1"})
		}

		if err := c.Validate(req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error(), "message": "invalid request2"})
		}

		// check if there is any user with given username and email
		// if err is not nil means we find a user with given email and username
		clientCollection := mongoClient.Database("chat_app").Collection("clients")

		var existingUser models.Client
		err := clientCollection.FindOne(context.TODO(), bson.M{
			"$or": []bson.M{
				{"username": req.Username},
				{"email": req.Email},
			}}).Decode(&existingUser)
		if err == nil {
			return c.JSON(http.StatusConflict, map[string]string{"error": "Username with give username and email already exists"})
		}

		// hash the password of the user
		hashedPassword, err := utils.HashPassword(req.Password)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

		// insert the new user into database
		newUser := models.Client{
			Username:     req.Username,
			Email:        req.Email,
			PasswordHash: hashedPassword,
			CreatedAt:    time.Now(),
		}
		_, err = clientCollection.InsertOne(context.TODO(), newUser)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

		return c.JSON(http.StatusCreated, map[string]string{"message": "User registered successfully"})
	}

}
