package main

import (
	"log"
	"os"

	"github.com/alecruiz/portfolio-service/internal/handlers"
	"github.com/alecruiz/portfolio-service/internal/middleware"
	"github.com/alecruiz/portfolio-service/internal/repository"
	"github.com/alecruiz/portfolio-service/pkg/database"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Initialize database
	db, err := database.NewPostgresDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize repository
	projectRepo := repository.NewProjectRepository(db)

	// Initialize handlers
	projectHandler := handlers.NewProjectHandler(projectRepo)

	// Setup Gin router
	router := gin.Default()

	// Apply global middleware
	router.Use(middleware.CORSMiddleware())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":   "healthy",
			"service":  "portfolio-service",
			"database": "connected",
		})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		projects := v1.Group("/projects")
		{
			// Public routes (anyone can view projects)
			projects.GET("", projectHandler.GetAllProjects)
			projects.GET("/:id", projectHandler.GetProject)

			// Protected routes (require authentication)
			projects.POST("", middleware.AuthMiddleware(), projectHandler.CreateProject)
			projects.PUT("/:id", middleware.AuthMiddleware(), projectHandler.UpdateProject)
			projects.DELETE("/:id", middleware.AuthMiddleware(), projectHandler.DeleteProject)
		}
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8002"
	}

	log.Printf("Portfolio Service starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
