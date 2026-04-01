package main

import (
	"log"
	"os"

	"github.com/alecruiz/portfolio-auth-service/internal/handlers"
	"github.com/alecruiz/portfolio-auth-service/internal/middleware"
	"github.com/alecruiz/portfolio-auth-service/internal/repository"
	"github.com/alecruiz/portfolio-auth-service/pkg/database"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Connect to database
	db, err := database.NewPostgresDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Connect to Redis
	redisClient := database.NewRedisClient()
	if redisClient != nil {
		defer redisClient.Close()
	}

	// Initialize repositories and handlers
	userRepo := repository.NewUserRepository(db)
	authHandler := handlers.NewAuthHandler(userRepo, redisClient)

	// Create a Gin router
	router := gin.Default()

	// Apply global middleware
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.RateLimitMiddleware(redisClient))

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) { // *gin.Context = pointer to a gin.Context struct, c is just a placeholder variable and could be named anything
		c.JSON(200, gin.H{ // gin.H is a shortcut for map[string]interface{} as defined in the gin package. map[string]interface{} is a map with string keys and values of any type.
			"status":   "healthy",
			"service":  "auth-service",
			"database": "connected",
		})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)

			// Protected routes (require authentication)
			auth.POST("/logout", middleware.AuthMiddleware(), authHandler.Logout)
			auth.GET("/me", middleware.AuthMiddleware(), authHandler.GetCurrentUser)
		}
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8001"
	}
	log.Printf("Auth Service starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
