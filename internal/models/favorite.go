package models

import "time"

type Favorite struct {
	ID         uint      `gorm:"primarykey" json:"id"`
	CreatedAt  time.Time `json:"created_at"`
	UserID     uint      `gorm:"not null;index:idx_user_property,unique" json:"user_id"`
	PropertyID uint      `gorm:"not null;index:idx_user_property,unique" json:"property_id"`

	Property Property `gorm:"foreignKey:PropertyID" json:"property,omitempty"`
}
