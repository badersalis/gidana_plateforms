package models

import "time"

type TransactionStatus string

const (
	StatusDone    TransactionStatus = "done"
	StatusFailed  TransactionStatus = "failed"
	StatusOngoing TransactionStatus = "ongoing"
)

type Transaction struct {
	ID                  uint              `gorm:"primarykey" json:"id"`
	CreatedAt           time.Time         `json:"created_at"`
	UserID              uint              `json:"user_id"`
	WalletID            uint              `json:"wallet_id"`
	Amount              float64           `json:"amount"`
	Nature              string            `gorm:"size:20" json:"nature"` // expense, income
	Service             string            `gorm:"size:100;not null" json:"service"`
	ServiceProvider     string            `gorm:"size:100;not null" json:"service_provider"`
	RelatedTransactionID *uint            `json:"related_transaction_id,omitempty"`
	Currency            string            `gorm:"size:5;default:'XOF'" json:"currency"`
	Status              TransactionStatus `gorm:"size:20;default:'ongoing'" json:"status"`

	Wallet Wallet `gorm:"foreignKey:WalletID" json:"wallet,omitempty"`
}
