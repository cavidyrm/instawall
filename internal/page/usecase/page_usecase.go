package usecase

import (
	"context"
	"fmt"
	"io"

	"github.com/cavidyrm/instawall/internal/page/domain"
	"github.com/google/uuid"
)

// --- Interface Definitions for Dependencies ---
type PageRepository interface {
	CreatePage(ctx context.Context, p *domain.Page) error
	GetPageByID(ctx context.Context, pageID uuid.UUID) (*domain.Page, error)
	GetAllPages(ctx context.Context, limit, offset int) ([]domain.Page, error)
	UpdatePage(ctx context.Context, p *domain.Page) error
	DeletePage(ctx context.Context, pageID, userID uuid.UUID) error
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

// --- Input DTOs ---
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
type UpdatePageInput struct {
	PageID      uuid.UUID
	UserID      uuid.UUID
	Title       string
	Description string
	Link        string
	HasIssue    bool
	CategoryIDs []uuid.UUID
	ImageFile   io.Reader // Optional: nil if not updating image
	ImageSize   int64
	ImageName   string
}

// --- Usecase Methods ---

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

func (uc *PageUsecase) GetPage(ctx context.Context, pageID uuid.UUID) (*domain.Page, error) {
	return uc.pageRepo.GetPageByID(ctx, pageID)
}

func (uc *PageUsecase) GetAllPages(ctx context.Context, limit, offset int) ([]domain.Page, error) {
	return uc.pageRepo.GetAllPages(ctx, limit, offset)
}

func (uc *PageUsecase) UpdatePage(ctx context.Context, input UpdatePageInput) (*domain.Page, error) {
	// First, get the existing page to ensure it exists and to have its current data.
	existingPage, err := uc.pageRepo.GetPageByID(ctx, input.PageID)
	if err != nil {
		return nil, fmt.Errorf("page not found")
	}

	// Authorization check: Ensure the user owns the page.
	if existingPage.UserID != input.UserID {
		return nil, fmt.Errorf("forbidden: user does not own this page")
	}

	// If a new image file is provided, upload it and update the URL.
	// Otherwise, keep the existing image URL.
	imageURL := existingPage.ImageURL
	if input.ImageFile != nil {
		newImageURL, err := uc.fileStore.UploadFile(ctx, input.ImageFile, input.ImageSize, input.ImageName)
		if err != nil {
			return nil, err
		}
		imageURL = newImageURL
	}

	// Update the page object with new data.
	pageToUpdate := &domain.Page{
		ID:          input.PageID,
		UserID:      input.UserID,
		Title:       input.Title,
		Description: input.Description,
		Link:        input.Link,
		HasIssue:    input.HasIssue,
		ImageURL:    imageURL,
	}

	// Save the updated page to the database.
	if err := uc.pageRepo.UpdatePage(ctx, pageToUpdate); err != nil {
		return nil, err
	}

	// Update category links.
	if err := uc.pageRepo.LinkPageToCategories(ctx, pageToUpdate.ID, input.CategoryIDs); err != nil {
		return nil, err
	}

	return pageToUpdate, nil
}

func (uc *PageUsecase) DeletePage(ctx context.Context, pageID, userID uuid.UUID) error {
	// In a real-world app, you might get the page first to ensure ownership
	// before deleting, but the SQL query also enforces this.
	return uc.pageRepo.DeletePage(ctx, pageID, userID)
}
