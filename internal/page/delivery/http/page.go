package http

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	appMiddleware "github.com/cavidyrm/instawall/internal/middleware"
	"github.com/cavidyrm/instawall/internal/page/usecase"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type PageHandler struct {
	pageUsecase *usecase.PageUsecase
}

func RegisterPageHandlers(e *echo.Echo, uc *usecase.PageUsecase) {
	h := &PageHandler{pageUsecase: uc}
	pageGroup := e.Group("/pages")

	// Public routes to view pages
	pageGroup.GET("", h.GetAllPages)
	pageGroup.GET("/:id", h.GetPage)

	// Authenticated routes to manage pages
	pageGroup.POST("", h.CreatePage, appMiddleware.JWTAuthMiddleware)
	pageGroup.PUT("/:id", h.UpdatePage, appMiddleware.JWTAuthMiddleware)
	pageGroup.DELETE("/:id", h.DeletePage, appMiddleware.JWTAuthMiddleware)
}

// --- Handler Methods ---

func (h *PageHandler) CreatePage(c echo.Context) error {
	userIDStr := c.Get("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, "Invalid user ID in token")
	}

	title := c.FormValue("title")
	description := c.FormValue("description")
	link := c.FormValue("link")
	hasIssue := c.FormValue("has_issue") == "true"
	categoryIDsStr := c.FormValue("category_ids")

	categoryIDs, err := parseUUIDs(categoryIDsStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("Invalid category_ids format: %v", err))
	}

	fileHeader, err := c.FormFile("image")
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Image file is required")
	}
	src, err := fileHeader.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Failed to open image file")
	}
	defer src.Close()

	input := usecase.CreatePageInput{
		UserID:      userID,
		Title:       title,
		Description: description,
		Link:        link,
		HasIssue:    hasIssue,
		CategoryIDs: categoryIDs,
		ImageFile:   src,
		ImageSize:   fileHeader.Size,
		ImageName:   fileHeader.Filename,
	}

	newPage, err := h.pageUsecase.CreatePage(c.Request().Context(), input)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Failed to create page: %v", err))
	}

	return c.JSON(http.StatusCreated, newPage)
}

func (h *PageHandler) GetPage(c echo.Context) error {
	pageID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid page ID")
	}

	p, err := h.pageUsecase.GetPage(c.Request().Context(), pageID)
	if err != nil {
		return c.JSON(http.StatusNotFound, "Page not found")
	}

	return c.JSON(http.StatusOK, p)
}

func (h *PageHandler) GetAllPages(c echo.Context) error {
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	offset, _ := strconv.Atoi(c.QueryParam("offset"))
	if limit == 0 {
		limit = 10 // Default limit
	}

	pages, err := h.pageUsecase.GetAllPages(c.Request().Context(), limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Failed to retrieve pages")
	}

	return c.JSON(http.StatusOK, pages)
}

func (h *PageHandler) UpdatePage(c echo.Context) error {
	userIDStr := c.Get("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, "Invalid user ID in token")
	}

	pageID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid page ID")
	}

	title := c.FormValue("title")
	description := c.FormValue("description")
	link := c.FormValue("link")
	hasIssue := c.FormValue("has_issue") == "true"
	categoryIDsStr := c.FormValue("category_ids")
	categoryIDs, err := parseUUIDs(categoryIDsStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("Invalid category_ids format: %v", err))
	}

	input := usecase.UpdatePageInput{
		PageID:      pageID,
		UserID:      userID,
		Title:       title,
		Description: description,
		Link:        link,
		HasIssue:    hasIssue,
		CategoryIDs: categoryIDs,
	}

	// Handle optional image update
	fileHeader, err := c.FormFile("image")
	if err == nil { // An image was provided
		src, err := fileHeader.Open()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, "Failed to open image file")
		}
		defer src.Close()
		input.ImageFile = src
		input.ImageSize = fileHeader.Size
		input.ImageName = fileHeader.Filename
	}

	updatedPage, err := h.pageUsecase.UpdatePage(c.Request().Context(), input)
	if err != nil {
		// Differentiate between not found/forbidden and other errors
		if strings.Contains(err.Error(), "forbidden") {
			return c.JSON(http.StatusForbidden, err.Error())
		}
		if strings.Contains(err.Error(), "not found") {
			return c.JSON(http.StatusNotFound, err.Error())
		}
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Failed to update page: %v", err))
	}

	return c.JSON(http.StatusOK, updatedPage)
}

func (h *PageHandler) DeletePage(c echo.Context) error {
	userIDStr := c.Get("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, "Invalid user ID in token")
	}

	pageID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid page ID")
	}

	if err := h.pageUsecase.DeletePage(c.Request().Context(), pageID, userID); err != nil {
		return c.JSON(http.StatusInternalServerError, "Failed to delete page")
	}

	return c.NoContent(http.StatusNoContent)
}

func parseUUIDs(s string) ([]uuid.UUID, error) {
	if s == "" {
		return nil, nil
	}
	parts := strings.Split(s, ",")
	uuids := make([]uuid.UUID, len(parts))
	for i, part := range parts {
		id, err := uuid.Parse(strings.TrimSpace(part))
		if err != nil {
			return nil, err
		}
		uuids[i] = id
	}
	return uuids, nil
}
