package repository

import (
	"database/sql"
	"fmt"

	"github.com/alecruiz/portfolio-service/internal/models"
	"github.com/lib/pq"
)

// ProjectRepository handles database operations for projects
type ProjectRepository struct {
	db *sql.DB
}

// NewProjectRepository creates a new project repository
func NewProjectRepository(db *sql.DB) *ProjectRepository {
	return &ProjectRepository{db: db}
}

// Create inserts a new project into the database
func (r *ProjectRepository) Create(project *models.Project) error {
	query := `
		INSERT INTO projects (user_id, title, description, technologies, github_url, live_url, image_url, is_featured)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, display_order, created_at, updated_at
	`

	err := r.db.QueryRow(
		query,
		project.UserID,
		project.Title,
		project.Description,
		pq.Array(project.Technologies), // Convert Go slice to Postgres array
		project.GithubURL,
		project.LiveURL,
		project.ImageURL,
		project.IsFeatured,
	).Scan(&project.ID, &project.DisplayOrder, &project.CreatedAt, &project.UpdatedAt)

	if err != nil {
		return fmt.Errorf("error creating project: %w", err)
	}

	return nil
}

// GetAll retrieves all projects, optionally filtered by user_id
func (r *ProjectRepository) GetAll(userID *int) ([]*models.Project, error) {
	var query string
	var args []interface{}

	if userID != nil {
		query = `
			SELECT id, user_id, title, description, technologies, github_url, live_url, 
			       image_url, display_order, is_featured, created_at, updated_at
			FROM projects
			WHERE user_id = $1
			ORDER BY display_order ASC, created_at DESC
		`
		args = append(args, *userID)
	} else {
		query = `
			SELECT id, user_id, title, description, technologies, github_url, live_url, 
			       image_url, display_order, is_featured, created_at, updated_at
			FROM projects
			ORDER BY display_order ASC, created_at DESC
		`
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("error getting projects: %w", err)
	}
	defer rows.Close()

	var projects []*models.Project
	for rows.Next() {
		project := &models.Project{}
		err := rows.Scan(
			&project.ID,
			&project.UserID,
			&project.Title,
			&project.Description,
			pq.Array(&project.Technologies), // Convert Postgres array to Go slice
			&project.GithubURL,
			&project.LiveURL,
			&project.ImageURL,
			&project.DisplayOrder,
			&project.IsFeatured,
			&project.CreatedAt,
			&project.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning project: %w", err)
		}
		projects = append(projects, project)
	}

	return projects, nil
}

// GetByID retrieves a project by ID
func (r *ProjectRepository) GetByID(id int) (*models.Project, error) {
	project := &models.Project{}
	query := `
		SELECT id, user_id, title, description, technologies, github_url, live_url, 
		       image_url, display_order, is_featured, created_at, updated_at
		FROM projects
		WHERE id = $1
	`

	err := r.db.QueryRow(query, id).Scan(
		&project.ID,
		&project.UserID,
		&project.Title,
		&project.Description,
		pq.Array(&project.Technologies),
		&project.GithubURL,
		&project.LiveURL,
		&project.ImageURL,
		&project.DisplayOrder,
		&project.IsFeatured,
		&project.CreatedAt,
		&project.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("project not found")
	}

	if err != nil {
		return nil, fmt.Errorf("error getting project: %w", err)
	}

	return project, nil
}

// Update updates a project
func (r *ProjectRepository) Update(id int, project *models.Project) error {
	query := `
		UPDATE projects
		SET title = $1, description = $2, technologies = $3, github_url = $4, 
		    live_url = $5, image_url = $6, display_order = $7, is_featured = $8,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $9
		RETURNING updated_at
	`

	err := r.db.QueryRow(
		query,
		project.Title,
		project.Description,
		pq.Array(project.Technologies),
		project.GithubURL,
		project.LiveURL,
		project.ImageURL,
		project.DisplayOrder,
		project.IsFeatured,
		id,
	).Scan(&project.UpdatedAt)

	if err == sql.ErrNoRows {
		return fmt.Errorf("project not found")
	}

	if err != nil {
		return fmt.Errorf("error updating project: %w", err)
	}

	return nil
}

// Delete deletes a project
func (r *ProjectRepository) Delete(id int) error {
	query := `DELETE FROM projects WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error deleting project: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("project not found")
	}

	return nil
}
