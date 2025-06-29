package middleware

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// JWTSecret is the secret key for signing the JWT.
// In a real-world application, this should be loaded from a secure configuration,
// like an environment variable, and should not be hardcoded.
var JWTSecret = []byte("your-very-secret-key")

// JWTCustomClaims are custom claims extending default ones.
type JWTCustomClaims struct {
	UserID string `json:"user_id"`
	Name   string `json:"name"`
	jwt.RegisteredClaims
}

// GenerateToken generates a new JWT for a given user.
func GenerateToken(userID, name string) (string, error) {
	// Set custom claims
	claims := &JWTCustomClaims{
		UserID: userID,
		Name:   name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)), // Token expires in 72 hours
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and return it
	t, err := token.SignedString(JWTSecret)
	if err != nil {
		return "", err
	}

	return t, nil
}

// JWTAuthMiddleware is the middleware for validating JWTs.
func JWTAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": "missing or malformed jwt"})
		}

		// The token should be in the format "Bearer <token>"
		tokenString := authHeader[len("Bearer "):]

		token, err := jwt.ParseWithClaims(tokenString, &JWTCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			return JWTSecret, nil
		})

		if err != nil || !token.Valid {
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": "invalid or expired jwt"})
		}

		claims, ok := token.Claims.(*JWTCustomClaims)
		if !ok {
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": "invalid jwt claims"})
		}

		// Attach user info to the context
		c.Set("user_id", claims.UserID)
		c.Set("user_name", claims.Name)

		return next(c)
	}
}