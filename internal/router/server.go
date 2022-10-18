package router

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/prybintsev/validation_cloud/internal/api"
	"github.com/prybintsev/validation_cloud/internal/db/users"
)

func StartHttpServer(ctx context.Context, dbCon *sql.DB) error {
	router := gin.Default()
	authGroup := router.Group("auth")

	usersRepo := users.NewUsersRepo(dbCon)
	authHandler := api.NewAuthHandler(usersRepo)
	authGroup.POST("/signup", authHandler.SignUp)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		err := srv.Shutdown(shutdownCtx)
		if err != nil {
			log.WithError(err).Error("Could not gracefully shut down http server")
		}
	}()

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}
