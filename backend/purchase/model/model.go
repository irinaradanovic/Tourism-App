package model

import "time"

// ShoppingCart - tourists cart for purchasing tours
type ShoppingCart struct {
	ID         uint        `gorm:"primaryKey" json:"id"`
	TouristID  uint        `gorm:"not null;uniqueIndex" json:"tourist_id"` // one active cart per tourist
	TotalPrice float64     `gorm:"type:numeric(10,2);default:0" json:"total_price"`
	Items      []OrderItem `gorm:"foreignKey:ShoppingCartID;constraint:OnDelete:CASCADE" json:"items"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
}

// an item (tour) in the shopping cart
type OrderItem struct {
	ID             uint    `gorm:"primaryKey" json:"id"`
	ShoppingCartID uint    `gorm:"not null" json:"shopping_cart_id"`
	TourID         uint    `gorm:"not null" json:"tourid"`
	TourName       string  `gorm:"type:varchar(255);not null" json:"tour_name"`
	Price          float64 `gorm:"type:numeric(10,2);not null" json:"price"`
}
