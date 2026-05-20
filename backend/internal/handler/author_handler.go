package handler

import (
	"database/sql"
	"net/http"
	"strconv"
	"web-gallery/internal/repository"

	"github.com/gin-gonic/gin"
)

type AuthorHandler struct {
	repo *repository.AuthorRepository
}

func NewAuthorHandler(repo *repository.AuthorRepository) *AuthorHandler {
	return &AuthorHandler{repo: repo}
}

func (h *AuthorHandler) GetAll(c *gin.Context) {
	authors, err := h.repo.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, authors)
}

func (h *AuthorHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	author, err := h.repo.GetByID(id)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, author)
}
