package model

import "time"

type EquipmentType struct {
	ID          uint      `gorm:"primaryKey;autoIncrement"`
	Name        string `gorm:"unique;not null"`
	Description string
	DailyRate   int64
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Quantity int `gorm:"not null" json:"quantity"`
}