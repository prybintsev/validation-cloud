package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Auth struct {
	SecretKey string
}

func NewAuth(secretKey string) *Auth {
	return &Auth{SecretKey: secretKey}
}

func (a *Auth) GenerateToken(username string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(10 * time.Minute)
	claims["user"] = username

	tokenString, err := token.SignedString([]byte(a.SecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil

}
