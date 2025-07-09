package postgres

import (
	"context"

	"github.com/cavidyrm/instawall/internal/category/domain"
	"github.com/google/uuid"
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

func (r *CategoryRepository) GetCategoryByID(ctx context.Context, categoryID uuid.UUID) (*domain.Category, error) {
	var c domain.Category
	query := `SELECT * FROM categories WHERE id = $1`
	err := r.db.GetContext(ctx, &c, query, categoryID)
	return &c, err
}

// GetAllCategories retrieves a list of all categories.
func (r *CategoryRepository) GetAllCategories(ctx context.Context) ([]domain.Category, error) {
	var categories []domain.Category
	query := `SELECT * FROM categories ORDER BY title ASC`
	err := r.db.SelectContext(ctx, &categories, query)
	return categories, err
}

// UpdateCategory updates an existing category's details.
func (r *CategoryRepository) UpdateCategory(ctx context.Context, c *domain.Category) error {
	query := `UPDATE categories SET title = $1, description = $2, image_url = $3 WHERE id = $4`
	_, err := r.db.ExecContext(ctx, query, c.Title, c.Description, c.ImageURL, c.ID)
	return err
}

// DeleteCategory removes a category from the database.
func (r *CategoryRepository) DeleteCategory(ctx context.Context, categoryID uuid.UUID) error {
	query := `DELETE FROM categories WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, categoryID)
	return err
}
