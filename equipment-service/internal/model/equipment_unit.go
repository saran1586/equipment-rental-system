
package model

import "time"

type EquipmentUnit struct {
	ID             uint   `gorm:"primaryKey;autoIncrement"`
	EquipmentTypeId uint `json:"equipment_type_id"`
	SerialNumber   string `gorm:"not null;uniqueIndex"`
	Condition      EquipmentCondition `gorm:"not null"`
	Status         EquipmentStatus    `gorm:"not null"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

