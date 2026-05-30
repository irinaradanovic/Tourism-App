package model

import "time"

type ShoppingCart struct {
	ID         uint        `gorm:"primaryKey" json:"id"`
	TouristID  int64       `gorm:"not null;uniqueIndex" json:"tourist_id"`
	TotalPrice float64     `gorm:"type:numeric(10,2);default:0" json:"total_price"`
	Items      []OrderItem `gorm:"foreignKey:ShoppingCartID;constraint:OnDelete:CASCADE" json:"items"`
	CreatedAt  time.Time   `gorm:"type:timestamp" json:"created_at"`
	UpdatedAt  time.Time   `gorm:"type:timestamp" json:"updated_at"`
}

type OrderItem struct {
	ID             uint    `gorm:"primaryKey" json:"id"`
	ShoppingCartID uint    `gorm:"not null" json:"shopping_cart_id"`
	TourID         string  `gorm:"not null" json:"tourid"`
	TourName       string  `gorm:"type:varchar(255);not null" json:"tour_name"`
	Price          float64 `gorm:"type:numeric(10,2);not null" json:"price"`
}

type TourPurchaseToken struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	TouristID int64     `gorm:"not null" json:"tourist_id"`
	TourID    string    `gorm:"not null" json:"tour_id"`
	TourName  string    `gorm:"type:varchar(255)" json:"tour_name"`
	CreatedAt time.Time `gorm:"type:timestamp" json:"created_at"`
}