package delivery

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"instawall/internal/usecase"
)

type ResetPasswordHandler struct {
	uc usecase.ResetPasswordUseCase
}

func NewResetPasswordHandler(e *echo.Echo, uc usecase.ResetPasswordUseCase) {
	h := &ResetPasswordHandler{uc: uc}

	e.POST("/request-reset", h.RequestReset)
	e.POST("/reset-password", h.ResetPassword)
}

type requestResetInput struct {
	Email string `json:"email" validate:"required,email"`
}

func (h *ResetPasswordHandler) RequestReset(c echo.Context) error {
	var input requestResetInput
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid input")
	}

	token, err := h.uc.RequestReset(input.Email)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{"reset_token": token})
}

type resetPasswordInput struct {
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=6"`
}

func (h *ResetPasswordHandler) ResetPassword(c echo.Context) error {
	var input resetPasswordInput
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid input")
	}

	if err := h.uc.ResetPassword(input.Token, input.NewPassword); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Password has been reset successfully."})
}
