package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
)

type GethHandler struct {
	samplesRepo SamplesRepo
}

func NewGethHandler(repo SamplesRepo) GethHandler {
	return GethHandler{
		samplesRepo: repo,
	}
}

type SamplesRepo interface {
	GetAverageGrowth(ctx context.Context) (float64, error)
}

type AvgGrowthResponse struct {
	AverageGrowth float64
}

func (h GethHandler) AverageGrowth(c *gin.Context) {
	avgGowth, err := h.samplesRepo.GetAverageGrowth(c)
	if err != nil {
		WriteErrorResponse(c, http.StatusInternalServerError, "failed to retrieve average growth")
		return
	}
	c.JSON(http.StatusOK, &AvgGrowthResponse{avgGowth})
}
