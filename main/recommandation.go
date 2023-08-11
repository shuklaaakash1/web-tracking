package main

import (
	"sort"
	"strconv"
)

func GetProductRecommendations(productID string) ([]Product, error) {
	// Fetch the product from the database by its ID
	product, err := GetProductByID(productID)
	if err != nil {
		return nil, err
	}

	relatedProducts, err := GetProductsByCategory(product.Category, "", "")
	if err != nil {
		return nil, err
	}

	// Apply the recommendation algorithm to calculate weighted scores
	const inCategoryWeight = 1.0
	const otherCategoryWeight = 0.5

	recommendedProducts := make([]Product, 0)
	for _, relatedProduct := range relatedProducts {
		var weightedScore float64
		if relatedProduct.Category == product.Category {
			weightedScore = float64(relatedProduct.Interactions) * inCategoryWeight
		} else {
			weightedScore = float64(relatedProduct.Interactions) * otherCategoryWeight
		}
		relatedProduct.WeightedScore = weightedScore
		recommendedProducts = append(recommendedProducts, relatedProduct)
	}

	// Sort recommended products by weighted score in descending order
	sort.SliceStable(recommendedProducts, func(i, j int) bool {
		return recommendedProducts[j].WeightedScore < recommendedProducts[i].WeightedScore
	})

	return recommendedProducts, nil
}

// type Interaction struct {
// 	ID        uint `gorm:"primaryKey"`
// 	UserID    uint `json:"user_id"`
// 	ProductID uint `json:"product_id"`
// 	// ... other fields ...
// }

func UpdateProductInteractions(productID uint) error {
	// Convert uint to string
	productIDStr := strconv.FormatUint(uint64(productID), 10)

	product, err := GetProductByID(productIDStr)
	if err != nil {
		return err
	}
	product.Interactions++
	return DB.Save(product).Error
}

func GetRelatedUsersForProduct(productID uint) ([]User, error) {
	// Fetch interactions for the given product
	var interactions []Interaction
	result := DB.Where("product_id = ?", productID).Find(&interactions)
	if result.Error != nil {
		return nil, result.Error
	}

	// Collect user IDs from interactions as uint values
	userIDs := make([]uint, len(interactions))
	for i, interaction := range interactions {
		userIDs[i] = interaction.UserID
	}

	// Fetch users with the collected user IDs
	var relatedUsers []User
	result = DB.Where("id IN ?", userIDs).Find(&relatedUsers)
	if result.Error != nil {
		return nil, result.Error
	}

	return relatedUsers, nil
}
