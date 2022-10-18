package api

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"

	"github.com/prybintsev/validation_cloud/internal/db"
)

type AuthHandler struct {
	usersRepo Users
	auth      Authorizer
}

type Users interface {
	CreateUser(ctx context.Context, userName, passwordHash string) error
	GetPasswordHashByUsername(ctx context.Context, userName string) (string, error)
}

type Authorizer interface {
	GenerateToken(username string) (string, error)
}

func NewAuthHandler(usersRepo Users, auth Authorizer) AuthHandler {
	return AuthHandler{usersRepo: usersRepo, auth: auth}
}

type UserCredentials struct {
	UserName *string `json:"user-name"`
	Password *string `json:"password"`
}

type CreateUserResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// validateUserCredentials checks whether current request contains user credentials.
// it only validates that the request contains all the required fields, but does not check the validity
// of the credentials themselves
func (h AuthHandler) validateUserCredentials(c *gin.Context) (*UserCredentials, error) {
	var req UserCredentials
	err := c.BindJSON(&req)
	var msg string
	if err != nil {
		msg = "malformed request"
		writeErrorResponse(c, http.StatusBadRequest, msg)
		return nil, err
	}
	if req.UserName == nil {
		msg = "missing user-name"
		writeErrorResponse(c, http.StatusBadRequest, msg)
		return nil, errors.New(msg)
	}
	if req.Password == nil {
		msg = "missing password"
		writeErrorResponse(c, http.StatusBadRequest, msg)
		return nil, errors.New(msg)
	}

	return &req, nil
}

func (h AuthHandler) SignUp(c *gin.Context) {
	req, err := h.validateUserCredentials(c)
	if err != nil {
		log.WithError(err).Error("invalid signup request")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.WithError(err).Error("failed to generate hash for password")
		writeErrorResponse(c, http.StatusInternalServerError, "signup has failed due to an internal error")
		return
	}

	err = h.usersRepo.CreateUser(c, *req.UserName, string(hashedPassword))
	if err != nil {
		log.WithError(err).Error("failed to create a user")

		if errors.Is(err, db.ErrorUserAlreadyExists) {
			writeErrorResponse(c, http.StatusBadRequest, "user already exists")
			return
		}
		writeErrorResponse(c, http.StatusInternalServerError, "failed to signup a user")
		return
	}
	writeCreateUserResponse(c)
}

func writeCreateUserResponse(c *gin.Context) {
	c.JSON(http.StatusOK, CreateUserResponse{Code: http.StatusOK, Message: "ok"})
}

func writeErrorResponse(c *gin.Context, code int, message string) {
	c.JSON(code, ErrorResponse{Code: code, Message: message})
}

type GenerateTokenResponse struct {
	Code  int    `json:"code"`
	Token string `json:"token"`
}

func (h AuthHandler) GenerateToken(c *gin.Context) {
	req, err := h.validateUserCredentials(c)
	if err != nil {
		log.WithError(err).Error("invalid signup request")
		return
	}

	hash, err := h.usersRepo.GetPasswordHashByUsername(c, *req.UserName)
	if err != nil {
		log.WithError(err).Error("failed to retrieve a user")

		if errors.Is(err, db.ErrorUserNotFound) {
			writeErrorResponse(c, http.StatusBadRequest, "user not found")
			return
		}
		writeErrorResponse(c, http.StatusInternalServerError, "failed retrieve the user")
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(*req.Password))
	if err != nil {
		writeErrorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	token, err := h.auth.GenerateToken(*req.UserName)
	if err != nil {
		writeErrorResponse(c, http.StatusInternalServerError, "failed to generate token")
		return
	}
	writeGenerateTokenResponse(c, token)
}

func writeGenerateTokenResponse(c *gin.Context, token string) {
	c.JSON(http.StatusOK, GenerateTokenResponse{Code: http.StatusOK, Token: token})
}
