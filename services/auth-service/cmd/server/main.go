package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Create a Gin router
	router := gin.Default()

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) { // *gin.Context = pointer to a gin.Context struct, c is just a placeholder variable and could be named anything
		c.JSON(200, gin.H{ // gin.H is a shortcut for map[string]interface{} as defined in the gin package. map[string]interface{} is a map with string keys and values of any type.
			"status":  "healthy",
			"service": "auth-service",
		})
	})

	// Start server on port 8001
	log.Println("Auth Service starting on port 8001")
	if err := router.Run(":8001"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
