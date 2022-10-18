package api

import (
	"context"
	"errors"
	"github.com/prybintsev/validation_cloud/internal/db"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	usersRepo Users
}

type Users interface {
	CreateUser(ctx context.Context, userName, passwordHash string) error
}

func NewAuthHandler(usersRepo Users) AuthHandler {
	return AuthHandler{usersRepo: usersRepo}
}

type SignupRequest struct {
	UserName *string `json:"user-name"`
	Password *string `json:"password"`
}

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (h AuthHandler) SignUp(c *gin.Context) {
	var req SignupRequest
	err := c.BindJSON(&req)
	if err != nil {
		writeResponse(c, http.StatusBadRequest, "malformed signup request")
		return
	}

	if req.UserName == nil {
		writeResponse(c, http.StatusBadRequest, "missing user-name")
		return
	}
	if req.Password == nil {
		writeResponse(c, http.StatusBadRequest, "missing password")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.WithError(err).Error("failed to generate hash for password")
		writeResponse(c, http.StatusInternalServerError, "signup has failed due to an internal error")
		return
	}

	err = h.usersRepo.CreateUser(c, *req.UserName, string(hashedPassword))
	if err != nil {
		log.WithError(err).Error("failed to create a user")

		if errors.Is(err, db.ErrorUserAlreadyExists) {
			writeResponse(c, http.StatusBadRequest, "user already exists")
			return
		}
		writeResponse(c, http.StatusInternalServerError, "failed to signup a user")
		return
	}
	writeResponse(c, http.StatusOK, "ok")
}

func writeResponse(c *gin.Context, code int, message string) {
	c.JSON(code, Response{Code: code, Message: message})
}
