package main

import (
	"fmt"
	"net/http"
)

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
