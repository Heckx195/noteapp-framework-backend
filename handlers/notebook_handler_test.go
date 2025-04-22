package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"noteapp-framework-backend/config"
	"noteapp-framework-backend/models"
)

func setupNotebookTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	protected := router.Group("/")
	protected.Use(mockNotebookAuthMiddleware())
	{
		protected.POST("/notebooks", CreateNotebook)
		protected.GET("/notebooks", GetNotebooks)
		protected.GET("/notebooks/:id", GetNotebook)
		protected.PUT("/notebooks/:id", UpdateNotebook)
		protected.DELETE("/notebooks/:id", DeleteNotebook)
		protected.GET("/notebookscount/", GetNotebookCount)
		protected.GET("/notescount/:notebookid", GetNoteCount)
		protected.GET("/notebookname/:id", GetNotebookName)
	}

	return router
}

func initNotebookTestDB() {
	// Initialize DB
	config.DBInit()

	// Clean up and reset for test
	// Clean up and reset for test
	config.DB.Exec("DELETE FROM users")
	config.DB.Exec("DELETE FROM notebooks")

	// Insert a test user
	config.DB.Exec("INSERT INTO users (id, username, password) VALUES (1, 'Test', 'password')")
}

func mockNotebookAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Simulate an authenticated user by setting a user ID in the context
		c.Set("user_id", "1")
		c.Next()
	}
}

func TestCreateNotebook(t *testing.T) {
	initNotebookTestDB()
	router := setupNotebookTestRouter()

	notebook := models.Notebook{Name: "Test Notebook"}
	body, _ := json.Marshal(notebook)

	req, _ := http.NewRequest("POST", "/notebooks", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var createdNotebook models.Notebook
	err := json.Unmarshal(w.Body.Bytes(), &createdNotebook)
	assert.NoError(t, err)
	assert.Equal(t, notebook.Name, createdNotebook.Name)
}
