package handlers

import (
	"fmt"
	"net/http"
	"noteapp-framework-backend/config"
	"noteapp-framework-backend/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateNote creates a new note
func CreateNote(c *gin.Context) {
	var input struct {
		Title      string `json:"title" binding:"required"`
		Content    string `json:"content" binding:"required"`
		NotebookID uint   `json:"notebook_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Log the input for debugging
	fmt.Printf("Received input: %+v\n", input)

	// Create a new note
	note := models.Note{
		Title:      input.Title,
		Content:    input.Content,
		NotebookID: input.NotebookID,
	}

	// Set user_id in note.
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}
	userIDStr, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID is not a valid string"})
		return
	}
	userIDUint, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}
	note.UserID = uint(userIDUint)

	if err := config.DB.Create(&note).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create note"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": note})
}

// DELETE: me not necessary to load their contents for the list
// TODO: Create new GetNotes where only a list is sent
// GetNotes retrieves all notes
func GetNotes(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	var notes []models.Note
	notebookID := c.Param("id")

	if err := config.DB.Where("user_id = ? AND notebook_id = ?", userID, notebookID).Find(&notes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch notes"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": notes,
	})
}

// GetNote retrieves a single note by ID
func GetNote(c *gin.Context) {
	// Retrieve user ID from the context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	id := c.Param("id")
	var note models.Note

	// Fetch the note and ensure it belongs to the authenticated user
	if err := config.DB.Where("id = ? AND user_id = ?", id, userID).First(&note).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"}) // TODO: Add "or access denied" to msg
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": note})
}

// UpdateNote updates an existing note
func UpdateNote(c *gin.Context) {
	// Retrieve user ID from the context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	id := c.Param("id")
	var note models.Note

	// Fetch the note and ensure it belongs to the authenticated user
	if err := config.DB.Where("id = ? AND user_id = ?", id, userID).First(&note).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"}) // TODO: Add "or access denied" to msg
		return
	}

	// Bind the updated data
	if err := c.ShouldBindJSON(&note); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Save the updated note
	if err := config.DB.Save(&note).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update note"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": note})
}

// DeleteNote deletes a note by ID
func DeleteNote(c *gin.Context) {
	// Retrieve user ID from the context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	id := c.Param("id")
	var note models.Note

	// Fetch the note and ensure it belongs to the authenticated user
	if err := config.DB.Where("id = ? AND user_id = ?", id, userID).First(&note).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"}) // TODO: Add "or access denied" to msg
		return
	}

	// Delete the note
	if err := config.DB.Delete(&note).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete note"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Note deleted successfully"})
}
