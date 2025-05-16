package delivery

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"instawall/internal/domain"
	"instawall/internal/usecase"
	"net/http"
)

type AuthHandler struct {
	uc usecase.AuthUseCase
}

func NewAuthHandler(e *echo.Echo, uc usecase.AuthUseCase) {
	handler := &AuthHandler{uc: uc}
	e.POST("/login", handler.Login)
	e.POST("/register", handler.Register)
	e.POST("/forgot-password", handler.ForgotPassword)
	e.POST("/reset-password", handler.ResetPassword)
}

func (h *AuthHandler) Login(c echo.Context) error {
	var req domain.LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request"})
	}

	token, err := h.uc.Login(c.Request().Context(), &req)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid email or password"})
	}

	return c.JSON(http.StatusOK, echo.Map{"token": token})
}

func (h *AuthHandler) Register(c echo.Context) error {
	var req domain.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request"})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	token, err := h.uc.Register(c.Request().Context(), &req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"token": token})
}

func (h *AuthHandler) ForgotPassword(c echo.Context) error {
	var req domain.ForgotPasswordRequest
	if err := c.Bind(&req); err != nil || c.Validate(&req) != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid email"})
	}

	if err := h.uc.ForgotPassword(c.Request().Context(), req.Email); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "If your email is registered, a reset link was sent"})
}

func (h *AuthHandler) ResetPassword(c echo.Context) error {
	var req domain.ResetPasswordRequest
	if err := c.Bind(&req); err != nil || c.Validate(&req) != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid input"})
	}

	if err := h.uc.ResetPassword(c.Request().Context(), req.Token, req.NewPassword); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "Password updated successfully"})
}
