package main

import (
	"log"
	"net/http"
	"os"
	"web-gallery/internal/db"
	"web-gallery/internal/handler"
	"web-gallery/internal/middleware"
	"web-gallery/internal/repository"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	loadDotEnv(".env")

	database := db.Connect()
	defer database.Close()

	paintingRepo := repository.NewPaintingRepository(database)
	authorRepo   := repository.NewAuthorRepository(database)
	materialRepo := repository.NewMaterialRepository(database)
	styleRepo    := repository.NewStyleRepository(database)
	plotRepo     := repository.NewPlotRepository(database)
	userRepo     := repository.NewUserRepository(database)

	paintingHandler := handler.NewPaintingHandler(paintingRepo)
	authorHandler   := handler.NewAuthorHandler(authorRepo)
	materialHandler := handler.NewMaterialHandler(materialRepo)
	styleHandler    := handler.NewStyleHandler(styleRepo)
	plotHandler     := handler.NewPlotHandler(plotRepo)
	authHandler     := handler.NewAuthHandler(userRepo)

	secret := os.Getenv("SESSION_SECRET")
	if secret == "" {
		secret = "change-me-in-production"
	}
	store := cookie.NewStore([]byte(secret))
	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	r := gin.Default()
	r.Use(corsMiddleware())
	r.Use(sessions.Sessions("gallery_session", store))

	r.GET("/static/*filepath", func(c *gin.Context) {
		c.Header("Cache-Control", "public, max-age=31536000, immutable")
		c.File("./static" + c.Param("filepath"))
	})

	api := r.Group("/api")
	{
		api.GET("/paintings",    paintingHandler.GetAll)
		api.GET("/paintings/:id", paintingHandler.GetByID)

		api.GET("/authors",    authorHandler.GetAll)
		api.GET("/authors/:id", authorHandler.GetByID)

		api.GET("/materials", materialHandler.GetAll)
		api.GET("/styles",    styleHandler.GetAll)
		api.GET("/plots",     plotHandler.GetAll)

		auth := api.Group("/auth")
		{
			auth.POST("/login",    authHandler.Login)
			auth.POST("/logout",   authHandler.Logout)
			auth.POST("/register", authHandler.Register)
			auth.GET("/me",        authHandler.Me)
		}

		admin := api.Group("/admin")
		admin.Use(middleware.RequireAdmin())
		{
			admin.POST("/authors",        authorHandler.Create)
			admin.PUT("/authors/:id",     authorHandler.Update)
			admin.DELETE("/authors/:id",  authorHandler.Delete)

			admin.POST("/paintings",       paintingHandler.Create)
			admin.PUT("/paintings/:id",    paintingHandler.Update)
			admin.DELETE("/paintings/:id", paintingHandler.Delete)
		}
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
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
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
