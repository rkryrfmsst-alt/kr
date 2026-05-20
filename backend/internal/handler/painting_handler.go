package handler

import (
	"database/sql"
	"net/http"
	"strconv"
	"web-gallery/internal/repository"

	"github.com/gin-gonic/gin"
)

type PaintingHandler struct {
	repo *repository.PaintingRepository
}

func NewPaintingHandler(repo *repository.PaintingRepository) *PaintingHandler {
	return &PaintingHandler{repo: repo}
}

func (h *PaintingHandler) GetAll(c *gin.Context) {
	paintings, err := h.repo.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, paintings)
}

func (h *PaintingHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	painting, err := h.repo.GetByID(id)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, painting)
}
