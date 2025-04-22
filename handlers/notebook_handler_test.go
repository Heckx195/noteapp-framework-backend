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
	config.DB.Exec("DELETE FROM notes")

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

func TestUpdateNotebook(t *testing.T) {
	initNotebookTestDB()
	router := setupNotebookTestRouter()

	// Insert a notebook manually for the test
	notebook := models.Notebook{ID: 1, Name: "Old Notebook Name", UserID: 1}
	config.DB.Create(&notebook)

	updatedNotebook := models.Notebook{Name: "Updated Notebook Name"}
	body, _ := json.Marshal(updatedNotebook)

	req, _ := http.NewRequest("PUT", "/notebooks/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Data models.Notebook `json:"data"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, updatedNotebook.Name, response.Data.Name)

	var updatedNotebookInDB models.Notebook
	config.DB.First(&updatedNotebookInDB, 1)
	assert.Equal(t, "Updated Notebook Name", updatedNotebookInDB.Name)
}

func TestDeleteNotebook(t *testing.T) {
	initNotebookTestDB()
	router := setupNotebookTestRouter()

	// Insert a notebook manually for the test
	notebook := models.Notebook{ID: 1, Name: "Test Notebook", UserID: 1}
	config.DB.Create(&notebook)

	req, _ := http.NewRequest("DELETE", "/notebooks/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Message string `json:"message"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Notebook deleted successfully", response.Message)

	// Verify the notebook was deleted from the database
	var deletedNotebook models.Notebook
	result := config.DB.First(&deletedNotebook, 1)
	assert.Error(t, result.Error) // Should return an error as the notebook no longer exists
}

func TestGetNotebookCount(t *testing.T) {
	initNotebookTestDB()
	router := setupNotebookTestRouter()

	// Insert notebooks manually for the test
	config.DB.Create(&models.Notebook{Name: "Notebook 1", UserID: 1})
	config.DB.Create(&models.Notebook{Name: "Notebook 2", UserID: 1})

	req, _ := http.NewRequest("GET", "/notebookscount/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		NotebookCount int `json:"notebook_count"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 2, response.NotebookCount) // Expecting 2 notebooks
}

func TestGetNoteCount(t *testing.T) {
	initNotebookTestDB()
	router := setupNotebookTestRouter()

	// Insert a notebook manually
	notebook := models.Notebook{ID: 1, Name: "Test Notebook", UserID: 1}
	config.DB.Create(&notebook)

	// Insert notes associated with the notebook manually
	config.DB.Create(&models.Note{Title: "Note 1", Content: "Content 1", NotebookID: 1, UserID: 1})
	config.DB.Create(&models.Note{Title: "Note 2", Content: "Content 2", NotebookID: 1, UserID: 1})

	req, _ := http.NewRequest("GET", "/notescount/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		NoteCount int64 `json:"note_count"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), response.NoteCount) // Expecting 2 notes
}

func TestGetNotebookName(t *testing.T) {
	initNotebookTestDB()
	router := setupNotebookTestRouter()

	// Insert a notebook manually for the test
	notebook := models.Notebook{ID: 1, Name: "Test Notebook", UserID: 1}
	config.DB.Create(&notebook)

	req, _ := http.NewRequest("GET", "/notebookname/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		NotebookName string `json:"notebook_name"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Test Notebook", response.NotebookName)
}
