package utils

import (
	"chatRoom/models"
	"github.com/golang-jwt/jwt"
	"os"
	"time"
)

// JWTCreator is the JWT token creator with given username
func JWTCreator(client models.Client) (string, error) {
	claims := jwt.MapClaims{
		"username": client.Username,
		"id":       client.ID.Hex(),
		"exp":      time.Now().Add(time.Hour * 12).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
