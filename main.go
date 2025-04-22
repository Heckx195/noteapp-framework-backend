package main

import (
	"log"

	"noteapp-framework-backend/config"
	"noteapp-framework-backend/handlers"
	"noteapp-framework-backend/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables.
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Initialize DB
	config.DBInit()

	r := gin.Default()

	// Enable CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Public routes
	r.POST("/register", handlers.Register)
	r.POST("/login", handlers.Login)
	r.POST("/refresh-token", handlers.RefreshToken)
	r.POST("/logout", handlers.Logout)

	// Protected routes
	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		// Notebook Routes
		protected.POST("/notebooks", handlers.CreateNotebook)
		protected.GET("/notebooks", handlers.GetNotebooks)
		protected.GET("/notebooks/:id", handlers.GetNotebook)
		protected.PUT("/notebooks/:id", handlers.UpdateNotebook)
		protected.DELETE("/notebooks/:id", handlers.DeleteNotebook)
		protected.GET("/notebookscount/", handlers.GetNotebookCount)
		protected.GET("/notescount/:notebookid", handlers.GetNoteCount)
		protected.GET("/notebookname/:id", handlers.GetNotebookName)

		// Note Routes
		protected.POST("/notes", handlers.CreateNote)
		protected.GET("/notes/:notebookid", handlers.GetNotes)
		protected.GET("/notes/:notebookid/pagination", handlers.GetNotesWithPagination)
		protected.GET("/notebyid/:notebookid/:noteid", handlers.GetNote)
		protected.PUT("/notes/:id", handlers.UpdateNote)
		protected.DELETE("/notes/:id", handlers.DeleteNote)

		// User Info Route
		protected.GET("/me", handlers.GetUserInfo)
		protected.POST("/changeusername", handlers.ChangeUsername)

	}

	r.Run(":8080")
}
