package main

import (
	"log"
	"net/http"
	"os"
	"web-gallery/internal/db"
	"web-gallery/internal/handler"
	"web-gallery/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	loadDotEnv(".env")

	database := db.Connect()
	defer database.Close()

	paintingRepo := repository.NewPaintingRepository(database)
	authorRepo := repository.NewAuthorRepository(database)
	materialRepo := repository.NewMaterialRepository(database)

	paintingHandler := handler.NewPaintingHandler(paintingRepo)
	authorHandler := handler.NewAuthorHandler(authorRepo)
	materialHandler := handler.NewMaterialHandler(materialRepo)

	r := gin.Default()
	r.Use(corsMiddleware())
	r.GET("/static/*filepath", func(c *gin.Context) {
		c.Header("Cache-Control", "public, max-age=31536000, immutable")
		c.File("./static" + c.Param("filepath"))
	})

	api := r.Group("/api")
	{
		api.GET("/paintings", paintingHandler.GetAll)
		api.GET("/paintings/:id", paintingHandler.GetByID)

		api.GET("/authors", authorHandler.GetAll)
		api.GET("/authors/:id", authorHandler.GetByID)

		api.GET("/materials", materialHandler.GetAll)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("сервер запущен на :%s\n", port)
	r.Run(":" + port)
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

func loadDotEnv(path string) {
	godotenv.Load(path)
}
