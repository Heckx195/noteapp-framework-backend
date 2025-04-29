package main

import (
	"export-service/handlers"
	"export-service/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Apply Logging Middleware
	r.Use(middleware.LoggingMiddleware())

	// Routes
	r.POST("/export/notebook", handlers.ExportNotebook)
	r.POST("/export/note", handlers.ExportNote)

	r.Run(":8081")
}
