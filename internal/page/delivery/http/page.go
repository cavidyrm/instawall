package http

import (
	"fmt"
	"net/http"
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
	pageGroup.Use(appMiddleware.JWTAuthMiddleware)
	pageGroup.POST("", h.CreatePage)
}

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
