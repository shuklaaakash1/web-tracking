package main

import "sort"

func GetProductRecommendations(productID string) ([]Product, error) {
	// Fetch the product from the database by its ID
	product, err := GetProductByID(productID)
	if err != nil {
		return nil, err
	}

	// Fetch products from the same category
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
