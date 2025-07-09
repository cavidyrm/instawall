package http

import (
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"strings"

	"github.com/cavidyrm/instawall/internal/category/usecase"
	appMiddleware "github.com/cavidyrm/instawall/internal/middleware"
	"github.com/labstack/echo/v4"
)

type CategoryHandler struct {
	categoryUsecase *usecase.CategoryUsecase
}

func RegisterCategoryHandlers(e *echo.Echo, uc *usecase.CategoryUsecase) {
	h := &CategoryHandler{categoryUsecase: uc}
	categoryGroup := e.Group("/categories")

	// Public routes to view categories
	categoryGroup.GET("", h.GetAllCategories)
	categoryGroup.GET("/:id", h.GetCategory)

	// Admin-only routes to manage categories
	adminCategoryGroup := categoryGroup.Group("")
	adminCategoryGroup.Use(appMiddleware.JWTAuthMiddleware, appMiddleware.AdminOnlyMiddleware)
	adminCategoryGroup.POST("", h.CreateCategory)
	adminCategoryGroup.PUT("/:id", h.UpdateCategory)
	adminCategoryGroup.DELETE("/:id", h.DeleteCategory)
}

// --- Handler Methods ---

func (h *CategoryHandler) CreateCategory(c echo.Context) error {
	title := c.FormValue("title")
	description := c.FormValue("description")

	fileHeader, err := c.FormFile("image")
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Image file is required")
	}
	src, err := fileHeader.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Failed to open image file")
	}
	defer src.Close()

	input := usecase.CreateCategoryInput{
		Title:       title,
		Description: description,
		ImageFile:   src,
		ImageSize:   fileHeader.Size,
		ImageName:   fileHeader.Filename,
	}

	newCategory, err := h.categoryUsecase.CreateCategory(c.Request().Context(), input)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Failed to create category: %v", err))
	}

	return c.JSON(http.StatusCreated, newCategory)
}

func (h *CategoryHandler) GetCategory(c echo.Context) error {
	categoryID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid category ID")
	}

	cat, err := h.categoryUsecase.GetCategory(c.Request().Context(), categoryID)
	if err != nil {
		return c.JSON(http.StatusNotFound, "Category not found")
	}

	return c.JSON(http.StatusOK, cat)
}

func (h *CategoryHandler) GetAllCategories(c echo.Context) error {
	categories, err := h.categoryUsecase.GetAllCategories(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Failed to retrieve categories")
	}
	return c.JSON(http.StatusOK, categories)
}

func (h *CategoryHandler) UpdateCategory(c echo.Context) error {
	categoryID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid category ID")
	}

	title := c.FormValue("title")
	description := c.FormValue("description")

	input := usecase.UpdateCategoryInput{
		CategoryID:  categoryID,
		Title:       title,
		Description: description,
	}

	fileHeader, err := c.FormFile("image")
	if err == nil {
		src, err := fileHeader.Open()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, "Failed to open image file")
		}
		defer src.Close()
		input.ImageFile = src
		input.ImageSize = fileHeader.Size
		input.ImageName = fileHeader.Filename
	}

	updatedCategory, err := h.categoryUsecase.UpdateCategory(c.Request().Context(), input)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.JSON(http.StatusNotFound, err.Error())
		}
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Failed to update category: %v", err))
	}

	return c.JSON(http.StatusOK, updatedCategory)
}

func (h *CategoryHandler) DeleteCategory(c echo.Context) error {
	categoryID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid category ID")
	}

	if err := h.categoryUsecase.DeleteCategory(c.Request().Context(), categoryID); err != nil {
		return c.JSON(http.StatusInternalServerError, "Failed to delete category")
	}

	return c.NoContent(http.StatusNoContent)
}
