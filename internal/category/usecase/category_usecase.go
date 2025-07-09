package usecase

import (
	"context"
	"fmt"
	"github.com/cavidyrm/instawall/internal/category/domain"
	"github.com/google/uuid"
	"io"
)

// --- Interface Definitions for Dependencies ---
type CategoryRepository interface {
	CreateCategory(ctx context.Context, c *domain.Category) error
	GetCategoryByID(ctx context.Context, categoryID uuid.UUID) (*domain.Category, error)
	GetAllCategories(ctx context.Context) ([]domain.Category, error)
	UpdateCategory(ctx context.Context, c *domain.Category) error
	DeleteCategory(ctx context.Context, categoryID uuid.UUID) error
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

// --- Input DTOs ---
type CreateCategoryInput struct {
	Title       string
	Description string
	ImageFile   io.Reader
	ImageSize   int64
	ImageName   string
}
type UpdateCategoryInput struct {
	CategoryID  uuid.UUID
	Title       string
	Description string
	ImageFile   io.Reader // Optional
	ImageSize   int64
	ImageName   string
}

// --- Usecase Methods ---

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

func (uc *CategoryUsecase) GetCategory(ctx context.Context, categoryID uuid.UUID) (*domain.Category, error) {
	return uc.catRepo.GetCategoryByID(ctx, categoryID)
}

func (uc *CategoryUsecase) GetAllCategories(ctx context.Context) ([]domain.Category, error) {
	return uc.catRepo.GetAllCategories(ctx)
}

func (uc *CategoryUsecase) UpdateCategory(ctx context.Context, input UpdateCategoryInput) (*domain.Category, error) {
	existingCategory, err := uc.catRepo.GetCategoryByID(ctx, input.CategoryID)
	if err != nil {
		return nil, fmt.Errorf("category not found")
	}

	imageURL := existingCategory.ImageURL
	if input.ImageFile != nil {
		newImageURL, err := uc.fileStore.UploadFile(ctx, input.ImageFile, input.ImageSize, input.ImageName)
		if err != nil {
			return nil, err
		}
		imageURL = newImageURL
	}

	categoryToUpdate := &domain.Category{
		ID:          input.CategoryID,
		Title:       input.Title,
		Description: input.Description,
		ImageURL:    imageURL,
	}

	if err := uc.catRepo.UpdateCategory(ctx, categoryToUpdate); err != nil {
		return nil, err
	}
	return categoryToUpdate, nil
}

func (uc *CategoryUsecase) DeleteCategory(ctx context.Context, categoryID uuid.UUID) error {
	return uc.catRepo.DeleteCategory(ctx, categoryID)
}
