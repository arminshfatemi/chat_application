package apiRouters

import (
	"chatRoom/handlers/apiHandlers"
	"errors"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"os"
)

func AuthAPIRouter(e *echo.Echo, mongoClient *mongo.Client) {
	r := e.Group("/api/")
	// user authentication
	r.POST("user/signup/", apiHandlers.ClientSignUpHandler(mongoClient))
	r.POST("user/login/", apiHandlers.ClientLogInHandler(mongoClient))

	// protected routes
	r.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		ContextKey: "userToken",
		SigningKey: []byte(os.Getenv("SECRET_KEY")),
	}))

	r.POST("user/logout/", func(c echo.Context) error {
		token, ok := c.Get("userToken").(*jwt.Token)
		if !ok {
			return errors.New("JWT token missing or invalid")
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return errors.New("failed to cast claims as jwt.MapClaims")
		}

		// TODO: if need put the JWT in the black list

		return c.JSON(http.StatusOK, claims)
	})

}
