package model

type RentalItem struct {
 	ID uint `gorm:"primaryKey" json:"id"`
 	RentalID uint `gorm:"not null" json:"rental_id"`
 	EquipmentUnitID uint `gorm:"not null" json:"equipment_unit_id"`
 	IssueCondition string `gorm:"not null" json:"issue_condition"`
 	ReturnCondition string `json:"return_condition"`
 	RentPerDay float64 `gorm:"not null" json:"rent_per_day"`
 	DamageCharge float64 `gorm:"not null;default:0" json:"damage_charge"`
}