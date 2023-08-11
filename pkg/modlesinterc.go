package pkg

type Interaction struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint `json:"user_id"`
	ProductID uint `json:"product_id"`
	// ... other fields ...
}
