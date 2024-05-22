package util

import (
	"github.com/dgrijalva/jwt-go"
	"strconv"
	"time"
)

var jwtKey = []byte("QvVg6AvdQZ2wh3fs9jQMBvZcTB2y7i6cU9/W9Yb74+s=")

type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.StandardClaims
}

func GenerateJWT(userID int, role string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID: strconv.Itoa(userID),
		Role:   role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tok, err := token.SignedString(jwtKey)

	return tok, err
}
