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
