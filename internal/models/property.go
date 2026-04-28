package models

import (
	"time"

	"gorm.io/gorm"
)

type Property struct {
	ID              uint           `gorm:"primarykey" json:"id"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
	Title           string         `gorm:"size:100;not null" json:"title"`
	Description     string         `gorm:"type:text" json:"description"`
	Neighborhood    string         `gorm:"size:50;not null" json:"neighborhood"`
	Country         string         `gorm:"size:50;not null" json:"country"`
	PropertyType    string         `gorm:"size:50;not null" json:"property_type"` // Studio, Appartement, Maison
	TransactionType string         `gorm:"size:20;not null" json:"transaction_type"` // À louer, À vendre
	Rooms           int            `gorm:"not null" json:"rooms"`
	Bathrooms       int            `gorm:"not null" json:"bathrooms"`
	ShowerType      string         `gorm:"size:20" json:"shower_type"` // interne, externe
	Surface         float64        `json:"surface"`
	HasCourtyard    bool           `gorm:"default:false" json:"has_courtyard"`
	HasWater        bool           `gorm:"default:false" json:"has_water"`
	HasElectricity  bool           `gorm:"default:false" json:"has_electricity"`
	ExactAddress    string         `gorm:"size:200" json:"exact_address"`
	WhatsappContact string         `gorm:"size:20" json:"-"`
	PhoneContact    string         `gorm:"size:20" json:"-"`
	Price           float64        `gorm:"not null" json:"price"`
	Currency        string         `gorm:"size:5;default:'XOF'" json:"currency"`
	IsAvailable     bool           `gorm:"default:true" json:"is_available"`
	OwnerID         uint           `gorm:"not null" json:"owner_id"`

	Owner   User            `gorm:"foreignKey:OwnerID" json:"owner,omitempty"`
	Images  []PropertyImage `gorm:"foreignKey:PropertyID;constraint:OnDelete:CASCADE" json:"images,omitempty"`
	Rentals []Rental        `gorm:"foreignKey:PropertyID;constraint:OnDelete:CASCADE" json:"-"`
	Reviews []Review        `gorm:"foreignKey:PropertyID;constraint:OnDelete:CASCADE" json:"reviews,omitempty"`
	Favorites []Favorite    `gorm:"foreignKey:PropertyID;constraint:OnDelete:CASCADE" json:"-"`

	AverageRating float64 `gorm:"-" json:"average_rating"`
	ReviewCount   int     `gorm:"-" json:"review_count"`
	IsFavorited   bool    `gorm:"-" json:"is_favorited"`
}

func (p *Property) ComputeRating() {
	if len(p.Reviews) == 0 {
		p.AverageRating = 0
		p.ReviewCount = 0
		return
	}
	total := 0
	for _, r := range p.Reviews {
		total += r.Rating
	}
	p.ReviewCount = len(p.Reviews)
	p.AverageRating = float64(total) / float64(p.ReviewCount)
}
