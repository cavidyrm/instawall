package usecase

import (
	"context"
	"github.com/cavidyrm/instawall/internal/category/domain"
	"io"
)

// --- Interface Definitions for Dependencies ---
type CategoryRepository interface {
	CreateCategory(ctx context.Context, c *domain.Category) error
}
type FileStore interface {
	UploadFile(ctx context.Context, file io.Reader, fileSize int64, originalFilename string) (string, error)
}

// --- Usecase Implementation ---
type CategoryUsecase struct {
	catRepo   CategoryRepository
	fileStore FileStore
}

func NewCategoryUsecase(cr CategoryRepository, fs FileStore) *CategoryUsecase {
	return &CategoryUsecase{catRepo: cr, fileStore: fs}
}

// --- Input DTO ---
type CreateCategoryInput struct {
	Title       string
	Description string
	ImageFile   io.Reader
	ImageSize   int64
	ImageName   string
}

// --- Usecase Method ---
func (uc *CategoryUsecase) CreateCategory(ctx context.Context, input CreateCategoryInput) (*domain.Category, error) {
	imageURL, err := uc.fileStore.UploadFile(ctx, input.ImageFile, input.ImageSize, input.ImageName)
	if err != nil {
		return nil, err
	}

	newCategory := &domain.Category{
		Title:       input.Title,
		Description: input.Description,
		ImageURL:    imageURL,
	}

	if err := uc.catRepo.CreateCategory(ctx, newCategory); err != nil {
		return nil, err
	}

	return newCategory, nil
}
