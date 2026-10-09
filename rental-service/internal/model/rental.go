package model

import "time"

type RentalStatus string

const (
	RentalStatusPending   RentalStatus = "PENDING"
	RentalStatusActive    RentalStatus = "ACTIVE"
	RentalStatusOverdue   RentalStatus = "OVERDUE"
	RentalStatusCompleted RentalStatus = "COMPLETED"
	RentalStatusCancelled RentalStatus = "CANCELLED"
)

type Rental struct {
	ID              uint       `gorm:"primaryKey" json:"id"`
	CustomerID      uint       `gorm:"not null" json:"customer_id"`
	RentalDate      time.Time  `gorm:"not null" json:"rental_date"`
	DueDate         time.Time  `gorm:"not null" json:"due_date"`
	ReturnDate      *time.Time `json:"return_date"`
	Status          string     `gorm:"not null" json:"status"`
	TotalRent       float64    `gorm:"not null;default:0" json:"total_rent"`
	LateFee         float64    `gorm:"not null;default:0" json:"late_fee"`
	LastLateFeeDate *time.Time `json:"last_late_fee_date"`
	LastDueDate     *time.Time `json:"last_due_date"`
	DamageFee       float64    `gorm:"not null;default:0" json:"damage_fee"`
}