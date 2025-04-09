package main

import (
	"noteapp-framework-backend/config"
	"noteapp-framework-backend/handlers"
	"noteapp-framework-backend/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize DB
	config.DBInit()

	r := gin.Default()

	// Public routes
	r.POST("/register", handlers.Register)
	r.POST("/login", handlers.Login)

	// Protected routes
	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		// Note Routes
		r.POST("/notes", handlers.CreateNote)
		r.GET("/notes", handlers.GetNotes)
		r.GET("/notes/:id", handlers.GetNote)
		r.PUT("/notes/:id", handlers.UpdateNote)
		r.DELETE("/notes/:id", handlers.DeleteNote)

		// Notebook Routes
		r.POST("/notebooks", handlers.CreateNotebook)
		r.GET("/notebooks", handlers.GetNotebooks)
		r.GET("/notebooks/:id", handlers.GetNotebook)
		r.PUT("/notebooks/:id", handlers.UpdateNotebook)
		r.DELETE("/notebooks/:id", handlers.DeleteNotebook)
	}

	r.Run(":8080")
}
