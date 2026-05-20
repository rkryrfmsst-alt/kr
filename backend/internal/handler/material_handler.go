package handler

import (
	"net/http"
	"web-gallery/internal/repository"

	"github.com/gin-gonic/gin"
)

type MaterialHandler struct {
	repo *repository.MaterialRepository
}

func NewMaterialHandler(repo *repository.MaterialRepository) *MaterialHandler {
	return &MaterialHandler{repo: repo}
}

func (h *MaterialHandler) GetAll(c *gin.Context) {
	materials, err := h.repo.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, materials)
}
