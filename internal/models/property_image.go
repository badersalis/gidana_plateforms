package models

import "time"

type PropertyImage struct {
	ID         uint      `gorm:"primarykey" json:"id"`
	CreatedAt  time.Time `json:"created_at"`
	Filename   string    `gorm:"size:500;not null" json:"filename"`
	PropertyID uint      `gorm:"not null" json:"property_id"`
	IsMain     bool      `gorm:"default:false" json:"is_main"`
}
