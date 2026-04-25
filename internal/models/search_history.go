package models

import "time"

type SearchHistory struct {
	ID         uint      `gorm:"primarykey" json:"id"`
	CreatedAt  time.Time `json:"created_at"`
	UserID     uint      `gorm:"not null" json:"user_id"`
	SearchTerm string    `gorm:"size:255" json:"search_term"`
}
