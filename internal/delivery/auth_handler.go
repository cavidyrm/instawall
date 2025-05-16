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
