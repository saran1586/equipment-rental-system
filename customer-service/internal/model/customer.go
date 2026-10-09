package model

import (
	"time"
)

type Customer struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	Name           string    `gorm:"not null" json:"name"`
	Email          string    `gorm:"unique;not null" json:"email"`
	Phone          string    `gorm:"unique" json:"phone"`
	KYCStatus      string    `gorm:"not null;default:PENDING" json:"kyc_status"`
	Deposit        float64   `gorm:"not null;default:0" json:"deposit"`
	OutStandingDue float64   `gorm:"not null;default:0" json:"outstanding_due"`
	Status         string    `gorm:"not null;default:ACTIVE" json:"status"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
