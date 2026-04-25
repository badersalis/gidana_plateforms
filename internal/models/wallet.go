package models

import "time"

type Wallet struct {
	ID             uint      `gorm:"primarykey" json:"id"`
	CreatedAt      time.Time `json:"created_at"`
	UserID         uint      `gorm:"not null" json:"user_id"`
	Provider       string    `gorm:"size:50;not null" json:"provider"` // Nita, MPesa, Visa, Mastercard, PayPal
	Nature         string    `gorm:"size:50" json:"nature"`
	PhoneNumber    string    `gorm:"size:20;uniqueIndex" json:"phone_number,omitempty"`
	Email          string    `gorm:"size:255;uniqueIndex" json:"email,omitempty"`
	CardNumber     string    `gorm:"size:20;uniqueIndex" json:"card_number,omitempty"`
	CVV            string    `gorm:"size:4" json:"-"`
	ExpirationDate string    `gorm:"size:7" json:"expiration_date,omitempty"`
	Password       string    `json:"-"`
	Balance        float64   `gorm:"default:0.0" json:"balance"`
	Currency       string    `gorm:"size:5;default:'XOF'" json:"currency"`
	Selected       bool      `gorm:"default:false" json:"selected"`

	Transactions []Transaction `gorm:"foreignKey:WalletID" json:"-"`

	// Masked display fields (computed, not stored)
	MaskedPhone  string `gorm:"-" json:"masked_phone,omitempty"`
	MaskedEmail  string `gorm:"-" json:"masked_email,omitempty"`
	MaskedCard   string `gorm:"-" json:"masked_card,omitempty"`
}

func (w *Wallet) ApplyMasks() {
	if w.PhoneNumber != "" && len(w.PhoneNumber) > 4 {
		w.MaskedPhone = w.PhoneNumber[:3] + "****" + w.PhoneNumber[len(w.PhoneNumber)-3:]
	}
	if w.Email != "" {
		atIdx := 0
		for i, c := range w.Email {
			if c == '@' {
				atIdx = i
				break
			}
		}
		if atIdx > 2 {
			w.MaskedEmail = w.Email[:2] + "****" + w.Email[atIdx:]
		} else {
			w.MaskedEmail = "****" + w.Email[atIdx:]
		}
	}
	if w.CardNumber != "" && len(w.CardNumber) >= 4 {
		w.MaskedCard = "**** **** **** " + w.CardNumber[len(w.CardNumber)-4:]
	}
}
