package httpserver

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"instawall/internal/usecase" // Replace with your module path
)

func NewHTTPServer(e *echo.Echo, authHandler *AuthHandler) *echo.Echo {
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/register", authHandler.Register)
	e.POST("/login", authHandler.Login)
	e.POST("/send-otp", authHandler.SendPasswordResetOTP)
	e.POST("/validate-otp", authHandler.ValidatePasswordResetOTP)
	e.POST("/reset-password", authHandler.ResetPassword)

	// Add a health check endpoint
	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	return e
}
func (h *AuthHandler) Login(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"message": "Login endpoint hit"})
}

func (h *AuthHandler) SendPasswordResetOTP(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"message": "Send OTP endpoint hit"})
}

func (h *AuthHandler) ValidatePasswordResetOTP(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"message": "Validate OTP endpoint hit"})
}

func (h *AuthHandler) ResetPassword(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"message": "Reset Password endpoint hit"})
}

// Placeholder for usecase package - replace with your actual package path
// package usecase
//
// type AuthUseCase interface {
// 	Register(...) error
// 	Login(...) (string, error) // Assuming login returns a token
// 	SendPasswordResetOTP(email string) error
// 	ValidatePasswordResetOTP(email, otp string) error
// 	ResetPassword(email, newPassword string) error
// }