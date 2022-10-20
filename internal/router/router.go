package router

import (
	"context"
	"database/sql"
	"github.com/prybintsev/validation_cloud/internal/api/middleware"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/prybintsev/validation_cloud/internal/api"
	"github.com/prybintsev/validation_cloud/internal/auth"
	"github.com/prybintsev/validation_cloud/internal/db/samples"
	"github.com/prybintsev/validation_cloud/internal/db/users"
	samplesJob "github.com/prybintsev/validation_cloud/internal/jobs/samples"
)

func StartHttpServer(ctx context.Context, dbCon *sql.DB, secretKey string) error {
	router := gin.Default()
	authGroup := router.Group("auth")

	usersRepo := users.NewUsersRepo(dbCon)
	authorizer := auth.NewAuth(secretKey)
	authHandler := api.NewAuthHandler(usersRepo, authorizer)
	authGroup.POST("signup", authHandler.SignUp)
	authGroup.POST("generate-token", authHandler.GenerateToken)

	samplesRepo := samples.NewSamplesRepo(dbCon)
	gethHandler := api.NewGethHandler(samplesRepo)
	gethGroup := router.Group("geth")
	gethGroup.Use(middleware.Auth(secretKey))
	gethGroup.GET("avg-growth", gethHandler.AverageGrowth)

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

	// Run a job that collects samples of frequency
	s := samplesJob.NewHeightSamplesCollector(time.Minute, samplesRepo)
	go s.Run(ctx)

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}
