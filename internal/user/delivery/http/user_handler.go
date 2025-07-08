package http

import (
	"net/http"

	appMiddleware "github.com/cavidyrm/instawall/internal/middleware"
	"github.com/cavidyrm/instawall/internal/user/usecase"
	"github.com/labstack/echo/v4"
)

// RegisterHandlers registers all handlers for the application.
func RegisterHandlers(e *echo.Echo, userUsecase *usecase.UserUsecase) {
	// Group all handlers under a single struct
	h := &handler{
		userUsecase: userUsecase,
	}

	// --- Authentication Flow ---
	authGroup := e.Group("/auth")
	authGroup.POST("/send-otp", h.SendOTP)
	authGroup.POST("/verify-otp", h.VerifyOTP)
	authGroup.POST("/login", h.Login)

	// This endpoint requires the special registration token
	regGroup := authGroup.Group("/complete-registration")
	regGroup.Use(appMiddleware.RegistrationTokenMiddleware)
	regGroup.POST("", h.CompleteRegistration)

	// --- User-Specific Routes (require standard login) ---
	userGroup := e.Group("/users")
	userGroup.Use(appMiddleware.JWTAuthMiddleware)
	userGroup.GET("/profile", h.GetProfile)

	// --- Admin-Only Routes (require login AND admin role) ---
	adminGroup := e.Group("/admin")
	adminGroup.Use(appMiddleware.JWTAuthMiddleware, appMiddleware.AdminOnlyMiddleware)
	adminGroup.GET("/dashboard", h.AdminDashboard)
}

// handler holds all dependencies for the HTTP handlers.
type handler struct {
	userUsecase *usecase.UserUsecase
}

// Request/Response Structs
type SendOTPRequest struct {
	MobileNumber string `json:"mobile_number" validate:"required"`
}
type VerifyOTPRequest struct {
	MobileNumber string `json:"mobile_number" validate:"required"`
	OTP          string `json:"otp" validate:"required"`
}
type CompleteRegistrationRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}
type LoginRequest struct {
	MobileNumber string `json:"mobile_number" validate:"required"`
	Password     string `json:"password" validate:"required"`
}
type ProfileResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

// --- Handler Methods ---

func (h *handler) SendOTP(c echo.Context) error {
	var req SendOTPRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid request")
	}
	err := h.userUsecase.SendOTP(c.Request().Context(), req.MobileNumber)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Failed to send OTP")
	}
	return c.JSON(http.StatusOK, echo.Map{"message": "OTP sent successfully"})
}

func (h *handler) VerifyOTP(c echo.Context) error {
	var req VerifyOTPRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid request")
	}
	token, err := h.userUsecase.VerifyOTP(c.Request().Context(), req.MobileNumber, req.OTP)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, echo.Map{"registration_token": token})
}

func (h *handler) CompleteRegistration(c echo.Context) error {
	mobileNumber := c.Get("verified_mobile").(string)
	var req CompleteRegistrationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid request")
	}
	_, err := h.userUsecase.CompleteRegistration(c.Request().Context(), mobileNumber, req.Password, req.Name, req.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Registration failed"})
	}
	return c.JSON(http.StatusCreated, echo.Map{"message": "User registered successfully"})
}

func (h *handler) Login(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid request")
	}
	token, err := h.userUsecase.Login(c.Request().Context(), req.MobileNumber, req.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid credentials"})
	}
	return c.JSON(http.StatusOK, echo.Map{"token": token})
}

func (h *handler) GetProfile(c echo.Context) error {
	userID := c.Get("user_id").(string)
	userProfile, err := h.userUsecase.GetProfile(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "User not found"})
	}
	resp := &ProfileResponse{
		ID:    userProfile.ID.String(),
		Name:  userProfile.Name,
		Email: userProfile.Email,
		Role:  userProfile.Role,
	}
	return c.JSON(http.StatusOK, resp)
}

func (h *handler) AdminDashboard(c echo.Context) error {
	userName := c.Get("user_name").(string)
	return c.JSON(http.StatusOK, echo.Map{
		"message": "Welcome to the admin dashboard, " + userName,
	})
}
