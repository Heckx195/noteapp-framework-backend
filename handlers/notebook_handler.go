package handlers

import (
	"net/http"
	"noteapp-framework-backend/config"
	"noteapp-framework-backend/models"

	"github.com/gin-gonic/gin"
)

// CreateNotebook creates a new notebook
func CreateNotebook(c *gin.Context) {
	var notebook models.Notebook
	if err := c.ShouldBindJSON(&notebook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Create(&notebook).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create notebook"})
		return
	}

	c.JSON(http.StatusCreated, notebook)
}

// GetNotebooks retrieves all notebooks
func GetNotebooks(c *gin.Context) {
	var notebooks []models.Notebook
	if err := config.DB.Preload("Notes").Find(&notebooks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch notebooks"})
		return
	}

	c.JSON(http.StatusOK, notebooks)
}

// GetNotebook retrieves a single notebook by ID
func GetNotebook(c *gin.Context) {
	id := c.Param("id")
	var notebook models.Notebook

	if err := config.DB.Preload("Notes").First(&notebook, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Notebook not found"})
		return
	}

	c.JSON(http.StatusOK, notebook)
}

// UpdateNotebook updates an existing notebook
func UpdateNotebook(c *gin.Context) {
	id := c.Param("id")
	var notebook models.Notebook

	if err := config.DB.First(&notebook, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Notebook not found"})
		return
	}

	if err := c.ShouldBindJSON(&notebook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Save(&notebook).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update notebook"})
		return
	}

	c.JSON(http.StatusOK, notebook)
}

// DeleteNotebook deletes a notebook by ID
func DeleteNotebook(c *gin.Context) {
	id := c.Param("id")
	if err := config.DB.Delete(&models.Notebook{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete notebook"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notebook deleted successfully"})
}
