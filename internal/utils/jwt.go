package utils

import (
	"crypto/rand"
	"errors"

	"github.com/golang-jwt/jwt/v5"
	"time"
	"fmt"
)

func GenerateRandomToken(length int) ([]byte, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return nil, errors.New("failed to generate random token")
	}
	return b, nil
}

func GenerateJWT(userID, ip, jti, secret string, expiry time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"sub":  userID,
		"ip":   ip,
		"jti":  jti,
		"exp":  time.Now().Add(expiry).Unix(),
		"iat":  time.Now().Unix(),
		"type": "access",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return token.SignedString([]byte(secret))
}

type CustomClaims struct {
	IP       string `json:"ip"`
	TokenType string `json:"type"`
	JTI      string `json:"jti"`
	jwt.RegisteredClaims
}

func ValidateJWT(tokenString, secret string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(secret), nil
	})
	/*
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	*/
	if claims, ok := token.Claims.(*CustomClaims); ok {
		return claims, nil
	}
	return nil, err
}