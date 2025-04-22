package handlers

import (
	"noteapp-framework-backend/config"

	"github.com/gin-gonic/gin"
)

func setupNoteTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	protected := router.Group("/")
	protected.Use(mockNotebookAuthMiddleware())
	{
		protected.POST("/notes", CreateNote)
		protected.GET("/notes/:notebookid", GetNotes)
		protected.GET("/notes/:notebookid/pagination", GetNotesWithPagination)
		protected.GET("/notebyid/:notebookid/:noteid", GetNote)
		protected.PUT("/notes/:id", UpdateNote)
		protected.DELETE("/notes/:id", DeleteNote)
	}

	return router
}

func initNoteTestDB() {
	// Initialize DB
	config.DBInit()

	// Clean up and reset for test
	config.DB.Exec("DELETE FROM users")
}

func mockNoteAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Simulate an authenticated user by setting a user ID in the context
		c.Set("userID", "1")
		c.Next()
	}
}
