package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
)

func sortProductsHandler(w http.ResponseWriter, r *http.Request) {
	sortBy := r.URL.Query().Get("sortBy") // Get the sorting field from query parameter
	order := r.URL.Query().Get("order")   // Get the order (asc or desc) from query parameter
	fmt.Println("sortby", sortBy)
	fmt.Println("order", order)
	// Validate sortBy and order values
	if sortBy == "" || order == "" {
		http.Error(w, "Bad Request: Invalid or missing parameters", http.StatusBadRequest)
		return
	}

	// Fetch products from the database based on sorting criteria
	products, err := GetProductsSortedBy(sortBy, order)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Sort products based on the specified field and order
	var sortFunc func(i, j int) bool
	switch sortBy {
	case "price":
		sortFunc = func(i, j int) bool {
			if order == "asc" {
				return products[i].Price < products[j].Price
			}
			return products[i].Price > products[j].Price
		}
	case "created_at":
		sortFunc = func(i, j int) bool {
			if order == "asc" {
				return products[i].CreatedAt.Before(products[j].CreatedAt)
			}
			return products[i].CreatedAt.After(products[j].CreatedAt)
		}
	default:
		http.Error(w, "Bad Request: Invalid sortBy value", http.StatusBadRequest)
		return
	}

	// Sort the products using the specified sorting function
	sort.Slice(products, sortFunc)

	// Respond with JSON containing sorted products
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}
func GetProductsSortedBy(sortBy string, order string) ([]Product, error) {
	var products []Product

	// Construct the ORDER BY clause based on the sortBy and order parameters
	orderClause := fmt.Sprintf("%s %s", sortBy, order)

	// Use the orderClause in your SQL query
	result := DB.Order(orderClause).Find(&products)
	if result.Error != nil {
		return nil, result.Error
	}
	return products, nil
}
