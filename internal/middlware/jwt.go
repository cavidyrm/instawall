package middlware

import (
	"github.com/labstack/echo/v4"
	"instawall/pkg/jwt"
	"net/http"
	"strings"
)

func JWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Missing token"})
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		userID, err := jwt.ParseToken(tokenStr)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid token"})
		}

		c.Set("user_id", userID)
		return next(c)
	}
}
