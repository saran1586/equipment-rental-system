package model

import "time"

type RentalDue struct {
	ID         uint       `gorm:"primaryKey" json:"id"`
	RentalID   uint       `gorm:"not null;index" json:"rental_id"`
	CustomerID uint       `gorm:"not null;index" json:"customer_id"`
	DueDate    time.Time  `gorm:"not null" json:"due_date"`
	DailyRent  float64    `gorm:"not null" json:"daily_rent"`
	Penalty    float64    `gorm:"not null;default:0" json:"penalty"`
	TotalDue   float64    `gorm:"not null" json:"total_due"`
	Status     string     `gorm:"not null;default:PENDING" json:"status"`
	PaidAt     *time.Time `json:"paid_at"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}