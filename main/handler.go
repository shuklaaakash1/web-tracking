package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func loginHandler(w http.ResponseWriter, r *http.Request) {

	var requestBody struct {
		Username string `json:"username"`
		Email    string `json:"email"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	username := requestBody.Username
	email := requestBody.Email
	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &JWTClaim{
		Username: username, // You need to set the email here
		Email:    email,    // You need to set the username here
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"token": "%s"}`, tokenString)
}

func protectedHandler(w http.ResponseWriter, r *http.Request) {
	tokenString := r.Header.Get("Authorization")
	strippedToken := extractTokenFromHeader(tokenString)

	if strippedToken == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err := ValidateToken(strippedToken)
	fmt.Println("err+", err)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Proceed with handling protected content
	fmt.Fprintf(w, "user authorized")
}

func extractTokenFromHeader(header string) string {
	const bearer = "Bearer "
	if len(header) > len(bearer) && header[:len(bearer)] == bearer {
		return header[len(bearer):]
	}
	return header
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	var newUser User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	err = CreateUser(&newUser)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, `{"message": "User created successfully"}`)
}
