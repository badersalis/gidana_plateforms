package models

import "time"

// DeletedAccount holds a snapshot of user data when a deletion is requested.
// Legal reviews the record before the deletion is finalized.
// Status values: "pending" | "approved" | "rejected"
type DeletedAccount struct {
	ID             uint       `gorm:"primarykey" json:"id"`
	UserID         uint       `gorm:"uniqueIndex;not null" json:"user_id"`
	FirstName      string     `gorm:"size:100" json:"first_name"`
	LastName       string     `gorm:"size:100" json:"last_name"`
	Email          string     `gorm:"size:255" json:"email"`
	PhoneNumber    string     `gorm:"size:20" json:"phone_number"`
	ProfilePicture string     `json:"profile_picture"`
	Gender         string     `gorm:"size:20" json:"gender"`
	DateOfBirth    *time.Time `json:"date_of_birth"`
	MemberSince    time.Time  `json:"member_since"`
	Locale         string     `gorm:"size:10" json:"locale"`
	Timezone       string     `gorm:"size:50" json:"timezone"`
	RequestedAt    time.Time  `gorm:"not null" json:"requested_at"`
	Status         string     `gorm:"size:20;default:'pending'" json:"status"`
	Notes          string     `json:"notes"`
}
