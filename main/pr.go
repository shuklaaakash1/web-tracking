package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func GetProductByID(productID string) (*Product, error) {
	var product Product
	result := DB.Where("id = ?", productID).First(&product)
	if result.Error != nil {
		return nil, result.Error
	}
	return &product, nil
}

// func getProductDetailsHandler(w http.ResponseWriter, r *http.Request) {
// 	// Parse product ID from request parameters
// 	productID := mux.Vars(r)["id"]

// 	// Fetch product details from the database
// 	product, err := GetProductByID(productID)
// 	if err != nil {
// 		http.Error(w, "Product not found", http.StatusNotFound)
// 		return
// 	}

// 	// Implement recommendation algorithm to get related products
// 	recommendations, err := GetProductRecommendations(productID)
// 	if err != nil {
// 		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
// 		return
// 	}

// 	// Combine product details and recommendations
// 	productWithRecommendations := struct {
// 		Product         Product   `json:"product"`
// 		Recommendations []Product `json:"recommendations"`
// 	}{
// 		Product:         *product,
// 		Recommendations: recommendations,
// 	}

//		// Respond with JSON including product details and recommendations
//		w.Header().Set("Content-Type", "application/json")
//		json.NewEncoder(w).Encode(productWithRecommendations)
//	}
func getProductDetailsHandler(w http.ResponseWriter, r *http.Request) {
	productID := mux.Vars(r)["id"]

	product, err := GetProductByID(productID)
	if err != nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	recommendations, err := GetProductRecommendations(productID)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	relatedUsers, err := GetRelatedUsersForProduct(product.ID)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	relatedUsernames := make([]string, len(relatedUsers))
	for i, user := range relatedUsers {
		relatedUsernames[i] = user.Username
	}

	productWithRecommendations := struct {
		Product          Product   `json:"product"`
		Recommendations  []Product `json:"recommendations"`
		RelatedUsernames []string  `json:"related_usernames"`
	}{
		Product:          *product,
		Recommendations:  recommendations,
		RelatedUsernames: relatedUsernames,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(productWithRecommendations)
}

func createProductHandler(w http.ResponseWriter, r *http.Request) {
	var newProduct Product
	err := json.NewDecoder(r.Body).Decode(&newProduct)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	err = CreateProduct(&newProduct)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, `{"message": "Product created successfully"}`)
}

// Create a new product in the database
func CreateProduct(product *Product) error {
	result := DB.Create(product)
	return result.Error
}

func GetProductsByCategory(category string, sortBy string, order string) ([]Product, error) {
	var products []Product

	// Construct the ORDER BY clause based on the sortBy and order parameters
	orderClause := ""
	if sortBy != "" && order != "" {
		orderClause = fmt.Sprintf("ORDER BY %s %s", sortBy, order)
	}

	// Construct the SQL query
	query := fmt.Sprintf("SELECT * FROM products WHERE category = ? %s", orderClause)

	// Execute the query and retrieve the products
	result := DB.Raw(query, category).Scan(&products)
	if result.Error != nil {
		return nil, result.Error
	}

	return products, nil
}
