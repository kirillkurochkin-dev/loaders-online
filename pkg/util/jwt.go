package util

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
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

func ValidateJWT(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	if claims.ExpiresAt < time.Now().Unix() {
		return nil, errors.New("token is expired")
	}

	logrus.WithFields(logrus.Fields{
		"user_id": claims.UserID,
		"role":    claims.Role,
	})

	return claims, nil
}
