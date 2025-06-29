package delivery

import (
	"net/http"

	appMiddleware "github.com/cavidyrm/instawall/internal/middleware" // Alias the import
	"github.com/cavidyrm/instawall/internal/user/usecase"
	"github.com/labstack/echo/v4"
)

// UserHandler handles HTTP requests for users.
type UserHandler struct {
	userUsecase *usecase.UserUsecase
}

// NewUserHandler creates a new UserHandler and registers the routes.
func NewUserHandler(e *echo.Echo, userUsecase *usecase.UserUsecase) {
	handler := &UserHandler{userUsecase: userUsecase}

	// Public routes
	e.POST("/register", handler.Register)
	e.POST("/login", handler.Login)

	// Protected routes
	profileGroup := e.Group("/profile")
	profileGroup.Use(appMiddleware.JWTAuthMiddleware) // Apply JWT middleware
	profileGroup.GET("", handler.GetProfile)
}

// RegisterRequest and LoginRequest structs remain the same.
type RegisterRequest struct {
	MobileNumber string `json:"mobile_number" validate:"required"`
	Password     string `json:"password" validate:"required,min=8"`
	Name         string `json:"name" validate:"required"`
	Email        string `json:"email" validate:"required,email"`
}

type LoginRequest struct {
	MobileNumber string `json:"mobile_number" validate:"required"`
	Password     string `json:"password" validate:"required"`
}

// Register handler remains the same.
func (h *UserHandler) Register(c echo.Context) error {
	var req RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request payload"})
	}
	// Add validation for the request here.

	ctx := c.Request().Context()
	newUser, err := h.userUsecase.Register(ctx, req.MobileNumber, req.Password, req.Name, req.Email)
	if err != nil {
		// In a real app, you should check for specific errors, like duplicate mobile/email
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to register user"})
	}

	return c.JSON(http.StatusCreated, newUser)
}

// Login handler is updated to return the token.
func (h *UserHandler) Login(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request payload"})
	}
	// Add validation for the request here.

	ctx := c.Request().Context()
	token, err := h.userUsecase.Login(ctx, req.MobileNumber, req.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid credentials"})
	}

	return c.JSON(http.StatusOK, echo.Map{"token": token})
}

// ProfileResponse defines the structure for the user profile response.
type ProfileResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// GetProfile is a new handler for the protected profile route.
func (h *UserHandler) GetProfile(c echo.Context) error {
	// Retrieve user ID from the context, set by the JWT middleware
	userID := c.Get("user_id").(string)

	ctx := c.Request().Context()
	userProfile, err := h.userUsecase.GetProfile(ctx, userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "User not found"})
	}

	// We don't want to expose all user data, so we map it to a response struct.
	resp := &ProfileResponse{
		ID:    userProfile.ID.String(),
		Name:  userProfile.Name,
		Email: userProfile.Email,
	}

	return c.JSON(http.StatusOK, resp)
}