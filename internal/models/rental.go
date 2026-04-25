package models

import "time"

type Rental struct {
	ID           uint       `gorm:"primarykey" json:"id"`
	CreatedAt    time.Time  `json:"created_at"`
	PropertyID   uint       `gorm:"not null" json:"property_id"`
	TenantID     uint       `gorm:"not null" json:"tenant_id"`
	StartDate    time.Time  `gorm:"not null" json:"start_date"`
	EndDate      *time.Time `json:"end_date,omitempty"`
	MonthlyPrice float64    `gorm:"not null" json:"monthly_price"`
	Status       string     `gorm:"size:20;default:'pending'" json:"status"` // pending, occupied, available, completed

	Property Property `gorm:"foreignKey:PropertyID" json:"property,omitempty"`
	Tenant   User     `gorm:"foreignKey:TenantID" json:"tenant,omitempty"`
}
