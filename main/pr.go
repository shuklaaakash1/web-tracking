package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"

	"github.com/gorilla/mux"
)

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
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Category string  `json:"category"`
	// Add more fields as needed
}

func GetProductRecommendations(userID int, category string) ([]Product, error) {
	// Implement your recommendation algorithm here
	// Fetch user interactions from the database for the given userID
	userInteractions, err := GetUserInteractions(userID)
	if err != nil {
		return nil, err
	}

	// Fetch products from the same category
	relatedProducts, err := GetProductsByCategory(category)
	if err != nil {
		return nil, err
	}

	// Apply the recommendation algorithm to calculate weighted scores
	const inCategory = 1.0
	const other = 0.5

	recommendedProducts := make([]Product, 0)
	for _, product := range relatedProducts {
		var weightedScore float64
		if product.Category == category {
			weightedScore = float64(product.Interactions) * inCategory
		} else {
			weightedScore = float64(product.Interactions) * other
		}
		product.WeightedScore = weightedScore
		recommendedProducts = append(recommendedProducts, product)
	}

	// Sort recommended products by weighted score in descending order
	sort.SliceStable(recommendedProducts, func(i, j int) bool {
		return recommendedProducts[j].WeightedScore < recommendedProducts[i].WeightedScore
	})

	return recommendedProducts, nil
}

type Product struct {
	ID            int     `json:"id"`
	Name          string  `json:"name"`
	Category      string  `json:"category"`
	Price         float64 `json:"price"`
	Interactions  int     `json:"interactions"`
	WeightedScore float64 `json:"weighted_score"` // Add this field
}
