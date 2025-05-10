package httpserver

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func healthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, echo.Map{
		"message": "Server is ok!",
	})
}
