package postgres

import (
	"context"
	"github.com/cavidyrm/instawall/internal/page/domain"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// PageRepository provides a database implementation for page operations.
type PageRepository struct {
	db *sqlx.DB
}

// NewPageRepository creates a new PageRepository.
func NewPageRepository(db *sqlx.DB) *PageRepository {
	return &PageRepository{db: db}
}

// CreatePage saves a new page to the database.
func (r *PageRepository) CreatePage(ctx context.Context, p *domain.Page) error {
	query := `INSERT INTO pages (user_id, title, description, image_url, link, has_issue)
			  VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at, updated_at`
	return r.db.QueryRowxContext(ctx, query, p.UserID, p.Title, p.Description, p.ImageURL, p.Link, p.HasIssue).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
}

// LinkPageToCategories associates a page with multiple categories in the join table.
func (r *PageRepository) LinkPageToCategories(ctx context.Context, pageID uuid.UUID, categoryIDs []uuid.UUID) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	stmt, err := tx.PreparexContext(ctx, `INSERT INTO page_categories (page_id, category_id) VALUES ($1, $2)`)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	for _, catID := range categoryIDs {
		if _, err := stmt.ExecContext(ctx, pageID, catID); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

// GetPageByID retrieves a single page by its ID.
func (r *PageRepository) GetPageByID(ctx context.Context, pageID uuid.UUID) (*domain.Page, error) {
	var p domain.Page
	query := `SELECT * FROM pages WHERE id = $1`
	err := r.db.GetContext(ctx, &p, query, pageID)
	return &p, err
}

// GetAllPages retrieves a paginated list of all pages.
func (r *PageRepository) GetAllPages(ctx context.Context, limit, offset int) ([]domain.Page, error) {
	var pages []domain.Page
	query := `SELECT * FROM pages ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	err := r.db.SelectContext(ctx, &pages, query, limit, offset)
	return pages, err
}

// UpdatePage updates an existing page's details in the database.
func (r *PageRepository) UpdatePage(ctx context.Context, p *domain.Page) error {
	query := `UPDATE pages SET title = $1, description = $2, image_url = $3, link = $4, has_issue = $5, updated_at = NOW()
			  WHERE id = $6 AND user_id = $7`
	_, err := r.db.ExecContext(ctx, query, p.Title, p.Description, p.ImageURL, p.Link, p.HasIssue, p.ID, p.UserID)
	return err
}

// DeletePage removes a page from the database.
func (r *PageRepository) DeletePage(ctx context.Context, pageID, userID uuid.UUID) error {
	query := `DELETE FROM pages WHERE id = $1 AND user_id = $2`
	_, err := r.db.ExecContext(ctx, query, pageID, userID)
	return err
}
