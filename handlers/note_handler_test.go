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

func setupNoteTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	protected := router.Group("/")
	protected.Use(mockNoteAuthMiddleware())
	{
		protected.POST("/notes", CreateNote)                                   // Check
		protected.GET("/notes/:notebookid", GetNotes)                          // Check
		protected.GET("/notes/:notebookid/pagination", GetNotesWithPagination) //
		protected.GET("/notebyid/:notebookid/:noteid", GetNote)                //
		protected.PUT("/notes/:id", UpdateNote)                                //
		protected.DELETE("/notes/:id", DeleteNote)                             //
	}

	return router
}

func initNoteTestDB() {
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

func mockNoteAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Simulate an authenticated user by setting a user ID in the context
		c.Set("user_id", "1")
		c.Next()
	}
}

func TestCreateNote(t *testing.T) {
	initNoteTestDB()
	router := setupNoteTestRouter()

	// Insert a notebook manually for the test
	notebook := models.Notebook{ID: 1, Name: "Test Notebook", UserID: 1}
	config.DB.Create(&notebook)

	note := models.Note{Title: "Test Note", Content: "Test Content", NotebookID: 1}
	body, _ := json.Marshal(note)

	req, _ := http.NewRequest("POST", "/notes", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	// Parse the response body
	var response struct {
		Data models.Note `json:"data"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	createdNote := response.Data
	assert.Equal(t, note.Title, createdNote.Title)
	assert.Equal(t, note.Content, createdNote.Content)
	assert.Equal(t, note.NotebookID, createdNote.NotebookID)
}

func TestGetNotes(t *testing.T) {
	initNoteTestDB()
	router := setupNoteTestRouter()

	// Insert a notebook manually for the test
	notebook := models.Notebook{ID: 1, Name: "Test Notebook", UserID: 1}
	config.DB.Create(&notebook)

	// Insert notes associated with the notebook
	note1 := models.Note{ID: 1, Title: "Note 1", Content: "Content 1", NotebookID: 1, UserID: 1}
	note2 := models.Note{ID: 2, Title: "Note 2", Content: "Content 2", NotebookID: 1, UserID: 1}
	config.DB.Create(&note1)
	config.DB.Create(&note2)

	req, _ := http.NewRequest("GET", "/notes/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Data []models.Note `json:"data"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, 2, len(response.Data))
	assert.Equal(t, "Note 1", response.Data[0].Title)
	assert.Equal(t, "Note 2", response.Data[1].Title)
}
