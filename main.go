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
		protected.GET("/notebooks/:id", handlers.GetNotebook) // postman // checked
		protected.PUT("/notebooks/:id", handlers.UpdateNotebook)
		protected.DELETE("/notebooks/:id", handlers.DeleteNotebook)
		protected.GET("/notebookscount/:id", handlers.GetNotebookCount)
		protected.GET("/notebookname/:id", handlers.GetNotebookName)

		// Note Routes
		protected.POST("/notes", handlers.CreateNote)
		protected.GET("/notes/:notebookid", handlers.GetNotes)
		protected.GET("/notes/:notebookid/pagination", handlers.GetNotesWithPagination)
		protected.GET("/notebyid/:notebookid/:noteid", handlers.GetNote) // postman // checked
		protected.PUT("/notes/:id", handlers.UpdateNote)
		protected.DELETE("/notes/:id", handlers.DeleteNote) // postman // checked

		// User Info Route
		protected.GET("/me", handlers.GetUserInfo)

		// TODO: Further 10 endpoints (in total 20)
		// TODO: 10 unit tests for 10 endpoints
		// TODO: golang migration -> remove AutoMigration

		// TODO: Fix random order of notes because of id 		-- checked
		// TODO: Implement pagenation							-- checked

		// TODO: Delete function Note frontend					-- checked
		// TODO: Delete function Note backend 					-- checked

		// TODO: Delete function Notebook frontend				-- checked
		// TODO: Delete function Notebook backend 				-- checked

		// TODO: Update name function Notebook frontend
		// TODO: Update name function Notebook backend			-- checked
	}

	r.Run(":8080")
}
