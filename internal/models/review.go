package models

import "time"

type Review struct {
	ID         uint      `gorm:"primarykey" json:"id"`
	CreatedAt  time.Time `json:"created_at"`
	PropertyID uint      `gorm:"not null" json:"property_id"`
	UserID     uint      `gorm:"not null" json:"user_id"`
	Rating     int       `gorm:"not null" json:"rating"` // 1-5
	Comment    string    `gorm:"type:text" json:"comment"`
	IsVerified bool      `gorm:"default:false" json:"is_verified"`

	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
