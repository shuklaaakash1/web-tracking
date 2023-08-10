package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"errors"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
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

func main() {
	ConnectDatabase()
	DB.AutoMigrate(&User{})
	secretKey, err := generateRandomSecretKey()
	if err != nil {
		fmt.Println("Error generating secret key:", err)
		return
	}
	jwtKey = []byte(secretKey)

	r := mux.NewRouter()

	r.HandleFunc("/login", loginHandler).Methods("POST")
	r.HandleFunc("/protected", protectedHandler).Methods("GET")
	r.HandleFunc("/create-user", createUserHandler).Methods("POST")
	// r.HandleFunc("/create-product", createProductHandler).Methods("POST")
	// r.HandleFunc("/product/{id}", getProductDetailsHandler).Methods("GET")

	http.Handle("/", r)

	port := "3000"
	fmt.Printf("Server is running on port %s\n", port)
	http.ListenAndServe(":"+port, nil)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	// Handle user login and generate JWT token
	// In a real application, validate user credentials and generate a token
	// user := User{
	// 	ID:       1,
	// 	Username: "user1",
	// 	Password: "password1",
	// 	Role:     "customer",
	// }
	// b, err := io.ReadAll(r.Body)
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// fmt.Println(string(b))
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

// func loginHandler(w http.ResponseWriter, r *http.Request) {
// 	var requestBody struct {
// 		Username string `json:"username"`
// 		Password string `json:"password"`
// 	}

// 	err := json.NewDecoder(r.Body).Decode(&requestBody)
// 	if err != nil {
// 		http.Error(w, "Bad Request", http.StatusBadRequest)
// 		return
// 	}
// 	username := requestBody.Username
// 	// password := requestBody.Password

// 	expirationTime := time.Now().Add(1 * time.Hour)
// 	claims := &JWTClaim{

// 		Username: username,
// 		StandardClaims: jwt.StandardClaims{
// 			ExpiresAt: expirationTime.Unix(),
// 		},
// 	}

// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 	tokenString, err := token.SignedString(jwtKey)
// 	if err != nil {
// 		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	fmt.Fprintf(w, `{"token": "%s"}`, tokenString)
// }

// Rest of your code for GetUserByUsername, CheckPasswordHash, User model, and DB setup

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

// Rest of the code for your User struct and database handling
