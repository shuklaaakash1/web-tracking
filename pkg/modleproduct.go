package pkg

import "time"

type Product struct {
	ID            uint `gorm:"primaryKey"`
	CreatedAt     time.Time
	Name          string  `json:"name"`
	Category      string  `json:"category"`
	Price         float64 `json:"price"`
	Interactions  int     `json:"interactions"`
	WeightedScore float64 `json:"weighted_score"`
}
