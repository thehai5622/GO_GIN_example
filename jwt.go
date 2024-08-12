package main

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func encodeJWT(username string) (string, error) {
	tokenInstance := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(10 * time.Minute),
	})
	secretKey := os.Getenv("SECRET_KEY")
	token, err := tokenInstance.SignedString([]byte(secretKey))

	if err != nil {
		return "", err
	}

	return token, nil
}

func decodeJWT(token string) (interface{}, error) {
	secretKey := os.Getenv("SECRET_KEY")
	tokenInstance, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != jwt.SigningMethodES256.Name {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}
	if !tokenInstance.Valid {
		return nil, jwt.ErrTokenSignatureInvalid
	}
	return tokenInstance.Claims, nil
}
