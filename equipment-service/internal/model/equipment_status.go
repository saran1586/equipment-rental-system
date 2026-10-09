package model

type EquipmentStatus string

const (
	StatusAvailable EquipmentStatus = "AVAILABLE"
	StatusRented   EquipmentStatus = "RENTED"
	StatusRepair   EquipmentStatus = "REPAIR"
	StatusInactive EquipmentStatus = "INACTIVE"
)

type EquipmentCondition string

const (
	ConditionNew     EquipmentCondition = "new"
	ConditionGood    EquipmentCondition = "good"
	ConditionDamaged EquipmentCondition = "damaged"
)