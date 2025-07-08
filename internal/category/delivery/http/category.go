package http

import (
	"fmt"
	"net/http"

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
	categoryGroup.Use(appMiddleware.JWTAuthMiddleware)
	categoryGroup.POST("", h.CreateCategory)
}

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
