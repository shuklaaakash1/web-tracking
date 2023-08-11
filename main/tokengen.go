package main

import (
	"crypto/rand"
	"encoding/base64"
	"errors"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey []byte

type JWTClaim struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.StandardClaims
}

func ValidateToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return err
	}

	if !token.Valid {
		return errors.New("Invalid token")
	}

	return nil
}

func generateRandomSecretKey() (string, error) {
	key := make([]byte, 32) // 32 bytes for a strong secret key
	_, err := rand.Read(key)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(key), nil
}
