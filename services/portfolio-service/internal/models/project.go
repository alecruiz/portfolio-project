package models

import "time"

// Project represents a portfolio project
type Project struct {
	ID           int       `json:"id"`
	UserID       int       `json:"user_id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	Technologies []string  `json:"technologies"` // We'll store as TEXT[] in Postgres
	GithubURL    string    `json:"github_url"`
	LiveURL      string    `json:"live_url"`
	ImageURL     string    `json:"image_url"`
	DisplayOrder int       `json:"display_order"`
	IsFeatured   bool      `json:"is_featured"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// CreateProjectRequest represents a request to create a project
type CreateProjectRequest struct {
	Title        string   `json:"title" binding:"required"`
	Description  string   `json:"description" binding:"required"`
	Technologies []string `json:"technologies" binding:"required"`
	GithubURL    string   `json:"github_url"`
	LiveURL      string   `json:"live_url"`
	ImageURL     string   `json:"image_url"`
	IsFeatured   bool     `json:"is_featured"`
}

// UpdateProjectRequest represents a request to update a project
type UpdateProjectRequest struct {
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	Technologies []string `json:"technologies"`
	GithubURL    string   `json:"github_url"`
	LiveURL      string   `json:"live_url"`
	ImageURL     string   `json:"image_url"`
	DisplayOrder int      `json:"display_order"`
	IsFeatured   bool     `json:"is_featured"`
}
