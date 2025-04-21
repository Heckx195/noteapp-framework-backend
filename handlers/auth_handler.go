package handlers

import (
	"fmt"
	"net/http"
	"noteapp-framework-backend/config"
	"noteapp-framework-backend/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create new user
	user := models.User{
		Username: input.Username,
		Password: input.Password,
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	user.Password = string(hashedPassword)

	// Create user
	if err := config.DB.Create(&user).Error; err != nil {
		// Check for unique constraint violation
		if isUniqueConstraintError(err) {
			c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"message":  "User registered successfully",
	})
}

func Login(c *gin.Context) {
	var loginData struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username and password are required"})
		return
	}

	var user models.User
	if err := config.DB.Where("username = ?", loginData.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Compare the stored hashed password with the input password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate Access Token (short-lived)
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Minute * 15).Unix(), // Access token expires in 15 minutes
	})

	// Signing access token
	accessTokenString, err := accessToken.SignedString([]byte(config.GetJWTSecret()))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	// Generate refresh token (long-lived)
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(), // Refresh token expires in 7 days
	})

	// Signing refresh token
	refreshTokenString, err := refreshToken.SignedString([]byte(config.GetJWTSecret()))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
		return
	}

	// Set cookies.
	c.SetCookie("access_token", accessTokenString, 15*60, "/", "localhost", false, true) // Expires in 15 minutes
	c.SetCookie("refresh_token", refreshTokenString, 7*24*60*60, "/", "localhost", false, false)

	// Respond with both tokens
	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessTokenString,
		"refresh_token": refreshTokenString,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
		},
	})
}

func RefreshToken(c *gin.Context) {
	// Get the refresh token from the cookie
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token is required"})
		return
	}

	// Validate the refresh token
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.GetJWTSecret()), nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		return
	}

	// Extract user ID from the refresh token
	userID, ok := claims["user_id"].(string)
	if !ok {
		// Try to read user_id as float64 and convert it to string
		if userIDFloat, ok := claims["user_id"].(float64); ok {
			userID = fmt.Sprintf("%.0f", userIDFloat) // Convert float64 to string
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in token"})
			c.Abort()
			return
		}
	}

	// Generate a new access token
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Minute * 15).Unix(), // Access token expires in 15 minutes
	})

	// Signing access token
	accessTokenString, err := accessToken.SignedString([]byte(config.GetJWTSecret()))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	// Respond with the new access token
	c.JSON(http.StatusOK, gin.H{
		"access_token": accessTokenString,
	})
}

func Logout(c *gin.Context) {
	c.SetCookie("access_token", "", -1, "/", "localhost", false, true)
	c.SetCookie("refresh_token", "", -1, "/", "localhost", false, false)

	c.JSON(http.StatusOK, gin.H{
		"message": "Logged out successfully",
	})
}

func ChangeUsername(c *gin.Context) {
	var request struct {
		NewUsername string `json:"new_username"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	user, err := FindUserByID(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
		return
	}

	user.Username = request.NewUsername
	if err := config.DB.Save(user).Error; err != nil {
		// Check for unique constraint violation
		if isUniqueConstraintError(err) {
			c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update username"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Username updated successfully"})
}

// Private helper functions.
func isUniqueConstraintError(err error) bool {
	return err != nil && err.Error() == `ERROR: duplicate key value violates unique constraint "uni_users_username" (SQLSTATE 23505)`
}
