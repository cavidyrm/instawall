package usecase

import (
	"context"
	"io"

	"github.com/cavidyrm/instawall/internal/page/domain"
	"github.com/google/uuid"
)

// --- Interface Definitions for Dependencies ---
type PageRepository interface {
	CreatePage(ctx context.Context, p *domain.Page) error
	LinkPageToCategories(ctx context.Context, pageID uuid.UUID, categoryIDs []uuid.UUID) error
}
type FileStore interface {
	UploadFile(ctx context.Context, file io.Reader, fileSize int64, originalFilename string) (string, error)
}

// --- Usecase Implementation ---
type PageUsecase struct {
	pageRepo  PageRepository
	fileStore FileStore
}

func NewPageUsecase(pr PageRepository, fs FileStore) *PageUsecase {
	return &PageUsecase{pageRepo: pr, fileStore: fs}
}

// --- Input DTO ---
type CreatePageInput struct {
	UserID      uuid.UUID
	Title       string
	Description string
	Link        string
	HasIssue    bool
	CategoryIDs []uuid.UUID
	ImageFile   io.Reader
	ImageSize   int64
	ImageName   string
}

// --- Usecase Method ---
func (uc *PageUsecase) CreatePage(ctx context.Context, input CreatePageInput) (*domain.Page, error) {
	imageURL, err := uc.fileStore.UploadFile(ctx, input.ImageFile, input.ImageSize, input.ImageName)
	if err != nil {
		return nil, err
	}

	newPage := &domain.Page{
		UserID:      input.UserID,
		Title:       input.Title,
		Description: input.Description,
		Link:        input.Link,
		HasIssue:    input.HasIssue,
		ImageURL:    imageURL,
	}

	if err := uc.pageRepo.CreatePage(ctx, newPage); err != nil {
		return nil, err
	}

	if len(input.CategoryIDs) > 0 {
		if err := uc.pageRepo.LinkPageToCategories(ctx, newPage.ID, input.CategoryIDs); err != nil {
			return nil, err
		}
	}

	return newPage, nil
}
