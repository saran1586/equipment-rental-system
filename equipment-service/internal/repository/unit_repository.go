
package repository

import (
	"equipment-service/internal/model"
     "log"
	"gorm.io/gorm"
)

type EquipmentUnitRepository struct {
	Db *gorm.DB
}

func NewEquipmentUnitRepository(db *gorm.DB) *EquipmentUnitRepository {
	return &EquipmentUnitRepository{
		Db: db,
	}
}

func (r *EquipmentUnitRepository) Create(unit *model.EquipmentUnit) error {
	return r.Db.Create(unit).Error
}

func (r *EquipmentUnitRepository) GetById(id uint) (*model.EquipmentUnit, error) {
	var equipmentUnit model.EquipmentUnit

	err := r.Db.First(&equipmentUnit, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	return &equipmentUnit, nil
}

func (r *EquipmentUnitRepository) Update(unit *model.EquipmentUnit) error {
	return r.Db.Save(unit).Error
}
func (r *EquipmentUnitRepository) Exists(serialNumber string) (bool, error) {
	var count int64

	err := r.Db.Model(&model.EquipmentUnit{}).
		Where("serial_number = ?", serialNumber).
		Count(&count).Error

	return count > 0, err
}

func (r *EquipmentUnitRepository) Deactivate(unit *model.EquipmentUnit) error {
	return r.Db.Save(unit).Error
}

func (r *EquipmentUnitRepository) Update1(unit *model.EquipmentUnit) error {
	return r.Db.Save(unit).Error
}

func (r *EquipmentUnitRepository) GetAll(limit int) ([]model.EquipmentUnit, error) {
	var equipmentUnits []model.EquipmentUnit

	err := r.Db.Limit(limit).Find(&equipmentUnits).Error
	if err != nil {
		log.Println("failed to get equipment units:", err)
		return nil, err
	}

	return equipmentUnits, nil
}