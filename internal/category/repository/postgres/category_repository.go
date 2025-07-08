package postgres

import (
	"context"
	"github.com/cavidyrm/instawall/internal/category/domain"
	"github.com/jmoiron/sqlx"
)

// CategoryRepository provides a database implementation for category operations.
type CategoryRepository struct {
	db *sqlx.DB
}

// NewCategoryRepository creates a new CategoryRepository.
func NewCategoryRepository(db *sqlx.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

// CreateCategory saves a new category to the database.
func (r *CategoryRepository) CreateCategory(ctx context.Context, c *domain.Category) error {
	query := `INSERT INTO categories (title, description, image_url)
			  VALUES ($1, $2, $3) RETURNING id, created_at`
	return r.db.QueryRowxContext(ctx, query, c.Title, c.Description, c.ImageURL).Scan(&c.ID, &c.CreatedAt)
}
