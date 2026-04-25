package models

import "time"

type Alert struct {
	ID              uint      `gorm:"primarykey" json:"id"`
	CreatedAt       time.Time `json:"created_at"`
	UserID          uint      `gorm:"not null" json:"user_id"`
	Neighborhood    string    `gorm:"size:50" json:"neighborhood"`
	PropertyType    string    `gorm:"size:50" json:"property_type"`
	MinRooms        int       `json:"min_rooms"`
	MaxPrice        float64   `json:"max_price"`
	TransactionType string    `gorm:"size:20" json:"transaction_type"` // Location, Vente
	IsActive        bool      `gorm:"default:true" json:"is_active"`
}
