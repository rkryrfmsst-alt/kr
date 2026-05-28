package handler

import (
	"net/http"
	"web-gallery/internal/repository"

	"github.com/gin-gonic/gin"
)

type PlotHandler struct {
	repo *repository.PlotRepository
}

func NewPlotHandler(repo *repository.PlotRepository) *PlotHandler {
	return &PlotHandler{repo: repo}
}

func (h *PlotHandler) GetAll(c *gin.Context) {
	plots, err := h.repo.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, plots)
}
