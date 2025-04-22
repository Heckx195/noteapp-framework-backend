package handlers

import (
	"net/http"

	"noteapp-framework-backend/config"
	"noteapp-framework-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetUserInfo handles the /me route
func GetUserInfo(c *gin.Context) {
	// Retrieve user ID from the Gin context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	// Fetch user from the database
	user, err := FindUserByID(userID.(string))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// If the user doesn't exist, clear the invalid user_id from the context
			c.Set("user_id", nil)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			return
		}
		// Handle other database errors
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
		return
	}

	// Remove sensitive fields (e.g., password) before sending the response
	user.Password = ""

	// Respond with the user data
	c.JSON(http.StatusOK, gin.H{
		"data": user,
	})
}

// FindUserByID fetches a user by their ID using GORM
func FindUserByID(userID string) (*models.User, error) {
	var user models.User
	result := config.DB.First(&user, "id = ?", userID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
