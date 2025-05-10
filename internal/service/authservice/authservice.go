package authservice

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID uint64
}

func (c Claims) Valid() error {
	return nil

}
func createToken(userID uint64, signKey string) (string, error) {
	claims := Claims{
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
		},
		userID,
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := accessToken.SignedString(signKey)
	if err != nil {
		return "", err
	}
	// Creat token string
	return token, nil
}
