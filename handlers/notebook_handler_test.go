package handlers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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
		protected.POST("/notebooks", CreateNotebook)           // Check
		protected.GET("/notebooks", GetNotebooks)              // Check
		protected.GET("/notebooks/:id", GetNotebook)           //
		protected.PUT("/notebooks/:id", UpdateNotebook)        //
		protected.DELETE("/notebooks/:id", DeleteNotebook)     //
		protected.GET("/notebookscount/", GetNotebookCount)    //
		protected.GET("/notescount/:notebookid", GetNoteCount) //
		protected.GET("/notebookname/:id", GetNotebookName)    //
	}

	return router
}

func initNotebookTestDB() {
	// Load environment variables.
	err := godotenv.Load("./../.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Initialize DB
	config.DBInit()

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

func TestGetNotebooks(t *testing.T) {
	initNotebookTestDB()
	router := setupNotebookTestRouter()

	// Insert notebook manually for test
	config.DB.Create(&models.Notebook{Name: "Test Notebook", UserID: 1})

	req, _ := http.NewRequest("GET", "/notebooks", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Data []models.Notebook `json:"data"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(response.Data), 1)
	assert.Equal(t, "Test Notebook", response.Data[0].Name)
}

func TestGetNotebook(t *testing.T) {
	initNotebookTestDB()
	router := setupNotebookTestRouter()

	// Insert notebook manually for test
	notebook := models.Notebook{ID: 1, Name: "Test Notebook", UserID: 1}
	config.DB.Create(&notebook)

	req, _ := http.NewRequest("GET", "/notebooks/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Data models.Notebook `json:"data"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, notebook.Name, response.Data.Name)
}
