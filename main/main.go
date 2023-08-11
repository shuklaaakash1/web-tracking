package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// var jwtKey []byte

// type JWTClaim struct {
// 	Username string `json:"username"`
// 	Email    string `json:"email"`
// 	jwt.StandardClaims
// }

// func ValidateToken(tokenString string) error {
// 	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
// 		return jwtKey, nil
// 	})
// 	if err != nil {
// 		return err
// 	}

// 	if !token.Valid {
// 		return errors.New("Invalid token")
// 	}

// 	return nil
// }

// func generateRandomSecretKey() (string, error) {
// 	key := make([]byte, 32) // 32 bytes for a strong secret key
// 	_, err := rand.Read(key)
// 	if err != nil {
// 		return "", err
// 	}
// 	return base64.StdEncoding.EncodeToString(key), nil
// }

func main() {
	ConnectDatabase()
	DB.AutoMigrate(&User{})
	DB.AutoMigrate(&Product{})
	DB.AutoMigrate(&Interaction{})
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
	r.HandleFunc("/create-product", createProductHandler).Methods("POST")
	r.HandleFunc("/product/{id}", getProductDetailsHandler).Methods("GET")
	r.HandleFunc("/products/sort", sortProductsHandler).Methods("GET")

	http.Handle("/", r)

	port := "3000"

	fmt.Printf("Server is running on port %s\n", port)
	http.ListenAndServe(":"+port, nil)
}

// Rest of the code for your User struct and database handling
