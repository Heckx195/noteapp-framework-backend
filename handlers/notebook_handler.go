package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"

	"noteapp-framework-backend/config"
	"noteapp-framework-backend/models"
)

// CreateNotebook creates a new notebook
func CreateNotebook(c *gin.Context) {
	var notebook models.Notebook
	if err := c.ShouldBindJSON(&notebook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set user_id in notebook.
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
	notebook.UserID = uint(userIDUint)

	if err := config.DB.Create(&notebook).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create notebook"})
		return
	}

	c.JSON(http.StatusCreated, notebook)
}

// GetNotebooks retrieves all notebooks for user_id
func GetNotebooks(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	var notebooks []models.Notebook
	if err := config.DB.Where("user_id = ?", userID).Find(&notebooks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch notebooks"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": notebooks,
	})
}

// GetNotebook retrieves a single notebook by ID
func GetNotebook(c *gin.Context) {
	// Retrieve user ID from the context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	id := c.Param("id")
	var notebook models.Notebook

	if err := config.DB.Where("id = ? AND user_id = ?", id, userID).First(&notebook).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Notebook not found or access denied"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": notebook})
}

// UpdateNotebook updates an existing notebook
func UpdateNotebook(c *gin.Context) {
	// Retrieve user ID from the context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	id := c.Param("id")
	var notebook models.Notebook

	// Fetch the notebook and ensure it belongs to the authenticated user
	if err := config.DB.Where("id = ? AND user_id = ?", id, userID).First(&notebook).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Notebook not found or access denied"})
		return
	}

	// Bind the updated data (but keep the original user_id)
	var updatedData models.Notebook
	if err := c.ShouldBindJSON(&updatedData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update only the fields that are allowed to change
	notebook.Name = updatedData.Name

	// Save the updated notebook
	if err := config.DB.Save(&notebook).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update notebook"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": notebook})
}

// DeleteNotebook deletes a notebook by ID
func DeleteNotebook(c *gin.Context) {
	// Retrieve user ID from the context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	id := c.Param("id")
	var notebook models.Notebook

	// Fetch the notebook and ensure it belongs to the authenticated user
	if err := config.DB.Where("id = ? AND user_id = ?", id, userID).First(&notebook).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Notebook not found or access denied"})
		return
	}

	// Delete all notes associated with the notebook
	if err := config.DB.Where("notebook_id = ?", id).Delete(&models.Note{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete associated notes"})
		return
	}

	// Delete the notebook
	if err := config.DB.Delete(&notebook).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete notebook"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notebook deleted successfully"})
}

func GetNoteCount(c *gin.Context) {
	// Retrieve user ID from the context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	notebookID := c.Param("notebookid")

	var notebook models.Notebook
	if err := config.DB.Where("id = ? AND user_id = ?", notebookID, userID).First(&notebook).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Notebook not found or access denied"})
		return
	}

	var noteCount int64
	if err := config.DB.Model(&models.Note{}).Where("notebook_id = ?", notebookID).Count(&noteCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count notes"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"note_count": noteCount})
}

func GetNotebookCount(c *gin.Context) {
	// Retrieve user ID from the context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	var notebookCount int64
	if err := config.DB.Model(&models.Notebook{}).Where("user_id = ?", userID).Count(&notebookCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count notebooks"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"notebook_count": notebookCount})
}

func GetNotebookName(c *gin.Context) {
	// Retrieve user ID from the context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	notebookID := c.Param("id")

	var notebook models.Notebook
	if err := config.DB.Where("id = ? AND user_id = ?", notebookID, userID).First(&notebook).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Notebook not found or access denied"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"notebook_name": notebook.Name})
}

// ExportNotebook forwards the export request to the export-service
func ExportNotebook(c *gin.Context) {
	notebookID := c.Param("id")

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	// Fetch the notebook from the database
	var notebook models.Notebook
	if err := config.DB.Where("id = ? AND user_id = ?", notebookID, userID).First(&notebook).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Notebook not found or access denied"})
		return
	}

	// Fetch the notes associated with the notebook
	var notes []models.Note
	if err := config.DB.Where("notebook_id = ?", notebookID).Find(&notes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch notes for the notebook"})
		return
	}

	// Prepare the request body for the export-service
	requestBody := map[string]interface{}{
		"notebook_name": notebook.Name,
		"notes":         notes,
	}

	// Call the export-service
	client := resty.New()
	resp, err := client.R().
		SetBody(requestBody).
		SetHeader("Content-Type", "application/json").
		Post("http://localhost:8081/export/notebook")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to export notebook"})
		return
	}

	// Forward the response from the export-service
	c.Data(resp.StatusCode(), resp.Header().Get("Content-Type"), resp.Body())
}
