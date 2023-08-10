package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"time"

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
func getProductDetailsHandler(w http.ResponseWriter, r *http.Request) {
	// Parse product ID from request parameters
	productID := mux.Vars(r)["id"]

	// Fetch product details from the database
	product, err := GetProductByID(productID)
	if err != nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	// Implement recommendation algorithm to get related products
	recommendations, err := GetProductRecommendations(productID)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Combine product details and recommendations
	productWithRecommendations := struct {
		Product         Product   `json:"product"`
		Recommendations []Product `json:"recommendations"`
	}{
		Product:         *product,
		Recommendations: recommendations,
	}

	// Respond with JSON including product details and recommendations
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

type Product struct {
	ID            uint `gorm:"primaryKey"`
	CreatedAt     time.Time
	Name          string  `json:"name"`
	Category      string  `json:"category"`
	Price         float64 `json:"price"`
	Interactions  int     `json:"interactions"`
	WeightedScore float64 `json:"weighted_score"`
}

func GetProductsByCategory(category string, sortBy string, order string) ([]Product, error) {
	var products []Product

	// Construct the ORDER BY clause based on the sortBy and order parameters
	orderClause := fmt.Sprintf("%s %s", sortBy, order)

	// Use the orderClause in your SQL query
	result := DB.Where("category = ?", category).Order(orderClause).Find(&products)
	if result.Error != nil {
		return nil, result.Error
	}
	return products, nil
}

func sortProductsHandler(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category") // Get the category from query parameter
	sortBy := r.URL.Query().Get("sortBy")     // Get the sorting field from query parameter
	order := r.URL.Query().Get("order")       // Get the order (asc or desc) from query parameter

	if category == "" || sortBy == "" || (order != "asc" && order != "desc") {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Fetch products from the database based on category, sortBy, and order
	products, err := GetProductsByCategory(category, sortBy, order)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Sort products based on the specified field and order
	switch sortBy {
	case "price":
		sort.Slice(products, func(i, j int) bool {
			if order == "asc" {
				return products[i].Price < products[j].Price
			}
			return products[i].Price > products[j].Price
		})
	case "created_at":
		sort.Slice(products, func(i, j int) bool {
			if order == "asc" {
				return products[i].CreatedAt.Before(products[j].CreatedAt)
			}
			return products[i].CreatedAt.After(products[j].CreatedAt)
		})
	}

	// Respond with JSON containing sorted products
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}
