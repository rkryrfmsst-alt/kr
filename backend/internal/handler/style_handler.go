package handler

import (
	"net/http"
	"web-gallery/internal/repository"

	"github.com/gin-gonic/gin"
)

type StyleHandler struct {
	repo *repository.StyleRepository
}

func NewStyleHandler(repo *repository.StyleRepository) *StyleHandler {
	return &StyleHandler{repo: repo}
}

func (h *StyleHandler) GetAll(c *gin.Context) {
	styles, err := h.repo.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, styles)
}
