package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/prybintsev/validation_cloud/internal/api"
	"net/http"
	"strings"
)

func Auth(secretKey string) gin.HandlerFunc {
	return func(context *gin.Context) {
		bearerToken := context.Request.Header.Get("Authorization")

		if len(strings.Split(bearerToken, " ")) != 2 {
			api.WriteErrorResponse(context, http.StatusUnauthorized, "invalid bearer token")
			context.Abort()
			return
		}
		tokenString := strings.Split(bearerToken, " ")[1]

		_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}))
		if err != nil {
			api.WriteErrorResponse(context, http.StatusUnauthorized, err.Error())
			context.Abort()
			return
		}

		context.Next()
	}
}
