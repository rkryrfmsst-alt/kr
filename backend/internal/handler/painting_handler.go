package handler

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
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

func (h *PaintingHandler) Create(c *gin.Context) {
	title := c.PostForm("title")
	if title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Название обязательно"})
		return
	}

	var description *string
	if d := c.PostForm("description"); d != "" {
		description = &d
	}

	var year *int
	if y := c.PostForm("year"); y != "" {
		yi, err := strconv.Atoi(y)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный год"})
			return
		}
		year = &yi
	}

	var materialID *int
	if m := c.PostForm("material_id"); m != "" {
		mi, err := strconv.Atoi(m)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректная техника"})
			return
		}
		materialID = &mi
	}

	parseIDs := func(key string) []int {
		strs := c.PostFormArray(key)
		ids := make([]int, 0, len(strs))
		for _, s := range strs {
			if id, err := strconv.Atoi(s); err == nil {
				ids = append(ids, id)
			}
		}
		return ids
	}
	authorIDs := parseIDs("author_ids[]")
	styleIDs  := parseIDs("style_ids[]")
	plotIDs   := parseIDs("plot_ids[]")

	var imagePath *string
	file, err := c.FormFile("image")
	if err == nil {
		ext := filepath.Ext(file.Filename)
		filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
		dst := "./static/uploads/" + filename
		if err := os.MkdirAll("./static/uploads", 0755); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка создания директории"})
			return
		}
		if err := c.SaveUploadedFile(file, dst); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сохранения файла"})
			return
		}
		p := "/static/uploads/" + filename
		imagePath = &p
	}

	painting, err := h.repo.Create(title, description, year, materialID, imagePath, authorIDs, styleIDs, plotIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, painting)
}

func (h *PaintingHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	title := c.PostForm("title")
	if title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Название обязательно"})
		return
	}

	var description *string
	if d := c.PostForm("description"); d != "" {
		description = &d
	}

	var year *int
	if y := c.PostForm("year"); y != "" {
		yi, err := strconv.Atoi(y)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный год"})
			return
		}
		year = &yi
	}

	var materialID *int
	if m := c.PostForm("material_id"); m != "" {
		mi, err := strconv.Atoi(m)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректная техника"})
			return
		}
		materialID = &mi
	}

	parseIDs := func(key string) []int {
		strs := c.PostFormArray(key)
		ids := make([]int, 0, len(strs))
		for _, s := range strs {
			if id, err := strconv.Atoi(s); err == nil {
				ids = append(ids, id)
			}
		}
		return ids
	}
	authorIDs := parseIDs("author_ids[]")
	styleIDs  := parseIDs("style_ids[]")
	plotIDs   := parseIDs("plot_ids[]")

	var imagePath *string
	file, err := c.FormFile("image")
	if err == nil {
		ext := filepath.Ext(file.Filename)
		filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
		dst := "./static/uploads/" + filename
		if err := os.MkdirAll("./static/uploads", 0755); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка создания директории"})
			return
		}
		if err := c.SaveUploadedFile(file, dst); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сохранения файла"})
			return
		}
		p := "/static/uploads/" + filename
		imagePath = &p
	}

	painting, err := h.repo.Update(id, title, description, year, materialID, imagePath, authorIDs, styleIDs, plotIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, painting)
}

func (h *PaintingHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.repo.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
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
