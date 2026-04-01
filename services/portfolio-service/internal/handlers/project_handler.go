package handlers

import (
	"net/http"
	"strconv"

	"github.com/alecruiz/portfolio-service/internal/models"
	"github.com/alecruiz/portfolio-service/internal/repository"
	"github.com/gin-gonic/gin"
)

// ProjectHandler handles project-related requests
type ProjectHandler struct {
	projectRepo *repository.ProjectRepository
}

// NewProjectHandler creates a new project handler
func NewProjectHandler(projectRepo *repository.ProjectRepository) *ProjectHandler {
	return &ProjectHandler{
		projectRepo: projectRepo,
	}
}

// CreateProject creates a new project
func (h *ProjectHandler) CreateProject(c *gin.Context) {
	var req models.CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user_id from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Create project
	project := &models.Project{
		UserID:       userID.(int),
		Title:        req.Title,
		Description:  req.Description,
		Technologies: req.Technologies,
		GithubURL:    req.GithubURL,
		LiveURL:      req.LiveURL,
		ImageURL:     req.ImageURL,
		IsFeatured:   req.IsFeatured,
	}

	if err := h.projectRepo.Create(project); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create project"})
		return
	}

	c.JSON(http.StatusCreated, project)
}

// GetAllProjects retrieves all projects
func (h *ProjectHandler) GetAllProjects(c *gin.Context) {
	// Optional: filter by user_id from query param
	userIDStr := c.Query("user_id")
	var userID *int
	if userIDStr != "" {
		id, err := strconv.Atoi(userIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id"})
			return
		}
		userID = &id
	}

	projects, err := h.projectRepo.GetAll(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get projects"})
		return
	}

	c.JSON(http.StatusOK, projects)
}

// GetProject retrieves a single project by ID
func (h *ProjectHandler) GetProject(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	project, err := h.projectRepo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	c.JSON(http.StatusOK, project)
}

// UpdateProject updates a project
func (h *ProjectHandler) UpdateProject(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	// Get user_id from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Check if project exists and belongs to user
	existingProject, err := h.projectRepo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	if existingProject.UserID != userID.(int) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to update this project"})
		return
	}

	var req models.UpdateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update fields (only if provided)
	if req.Title != "" {
		existingProject.Title = req.Title
	}
	if req.Description != "" {
		existingProject.Description = req.Description
	}
	if req.Technologies != nil {
		existingProject.Technologies = req.Technologies
	}
	if req.GithubURL != "" {
		existingProject.GithubURL = req.GithubURL
	}
	if req.LiveURL != "" {
		existingProject.LiveURL = req.LiveURL
	}
	if req.ImageURL != "" {
		existingProject.ImageURL = req.ImageURL
	}
	existingProject.DisplayOrder = req.DisplayOrder
	existingProject.IsFeatured = req.IsFeatured

	if err := h.projectRepo.Update(id, existingProject); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update project"})
		return
	}

	c.JSON(http.StatusOK, existingProject)
}

// DeleteProject deletes a project
func (h *ProjectHandler) DeleteProject(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	// Get user_id from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Check if project exists and belongs to user
	existingProject, err := h.projectRepo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	if existingProject.UserID != userID.(int) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to delete this project"})
		return
	}

	if err := h.projectRepo.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete project"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Project deleted successfully"})
}
