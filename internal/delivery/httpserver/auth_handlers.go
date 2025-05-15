package httpserver

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"

	"instawall/internal/domain"
	"instawall/internal/usecase"
)

type AuthHandler struct {
	AuthUseCase usecase.AuthUseCase
}

func NewAuthHandler(e *echo.Echo, authUseCase usecase.AuthUseCase) *AuthHandler {
	handler := &AuthHandler{AuthUseCase: authUseCase}
	e.POST("/register", handler.Register)
	e.POST("/login", handler.Login)
	e.POST("/password-reset/send-otp", handler.SendPasswordResetOTP)
	e.POST("/password-reset/validate-otp", handler.ValidatePasswordResetOTP)
	return handler
}

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SendOTPRequest struct {
	Email string `json:"email"`
}

type ValidateOTPRequest struct {
	Email string `json:"email"`
	OTP   string `json:"otp"`
}

func (h *AuthHandler) Register(c echo.Context) error {
	req := new(RegisterRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to hash password")
	}

	user := &domain.User{
		Username:  req.Username,
		Email:     req.Email,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
	}

	err = h.AuthUseCase.Register(c.Request().Context(), user)
	if err != nil {
		// TODO: Handle specific errors (e.g., email already exists)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "User registered successfully"})
}

func (h *AuthHandler) Login(c echo.Context) error {
	req := new(LoginRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// In a real application, you would generate and return a JWT token here
	token, err := h.AuthUseCase.Login(c.Request().Context(), req.Email, req.Password)
	if err != nil {
		// TODO: Handle specific errors (e.g., invalid credentials)
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid credentials")
	}

	return c.JSON(http.StatusOK, map[string]string{"token": token})
}

func (h *AuthHandler) SendPasswordResetOTP(c echo.Context) error {
	req := new(SendOTPRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err := h.AuthUseCase.SendPasswordResetOTP(c.Request().Context(), req.Email)
	if err != nil {
		// TODO: Handle specific errors (e.g., user not found, error sending email/SMS)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Password reset OTP sent"})
}

func (h *AuthHandler) ValidatePasswordResetOTP(c echo.Context) error {
	req := new(ValidateOTPRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// You would typically also include the new password in this request
	// and update the user's password in the use case after validation.
	err := h.AuthUseCase.ValidatePasswordResetOTP(c.Request().Context(), req.Email, req.OTP)
	if err != nil {
		// TODO: Handle specific errors (e.g., invalid or expired OTP)
		return echo.NewHTTPError(http.StatusBadRequest, "invalid or expired OTP")
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "OTP validated successfully"})
}