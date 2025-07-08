package middleware

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// JWTSecret should be loaded from a secure configuration in production.
var JWTSecret = []byte("your-very-secret-key")

// JWTCustomClaims are the claims for a standard logged-in user.
type JWTCustomClaims struct {
	UserID string `json:"user_id"`
	Name   string `json:"name"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// RegistrationClaims are for the temporary token used during registration.
type RegistrationClaims struct {
	MobileNumber string `json:"mobile_number"`
	jwt.RegisteredClaims
}

// GenerateToken creates a standard JWT for an authenticated user.
func GenerateToken(userID, name, role string) (string, error) {
	claims := &JWTCustomClaims{
		UserID: userID,
		Name:   name,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JWTSecret)
}

// GenerateRegistrationToken creates a short-lived token for completing registration.
func GenerateRegistrationToken(mobileNumber string) (string, error) {
	claims := &RegistrationClaims{
		MobileNumber: mobileNumber,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 10)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JWTSecret)
}

// JWTAuthMiddleware validates a standard user token.
func JWTAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": "missing or malformed jwt"})
		}
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
		c.Set("user_id", claims.UserID)
		c.Set("user_name", claims.Name)
		c.Set("user_role", claims.Role)
		return next(c)
	}
}

// AdminOnlyMiddleware must be used *after* JWTAuthMiddleware.
func AdminOnlyMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		role, ok := c.Get("user_role").(string)
		if !ok || role != "admin" {
			return c.JSON(http.StatusForbidden, echo.Map{"error": "Forbidden: Admins only"})
		}
		return next(c)
	}
}

// RegistrationTokenMiddleware validates the temporary registration token.
func RegistrationTokenMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": "missing registration token"})
		}
		tokenString := authHeader[len("Bearer "):]
		token, err := jwt.ParseWithClaims(tokenString, &RegistrationClaims{}, func(token *jwt.Token) (interface{}, error) {
			return JWTSecret, nil
		})
		if err != nil || !token.Valid {
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": "invalid or expired registration token"})
		}
		claims, ok := token.Claims.(*RegistrationClaims)
		if !ok {
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": "invalid registration token claims"})
		}
		c.Set("verified_mobile", claims.MobileNumber)
		return next(c)
	}
}
