package repository

import (
	"log"

	"equipment-service/internal/model"

	"gorm.io/gorm"
)

type EquipmentTypeRepository struct {
	Db *gorm.DB
}

func NewEquipmentTypeRepository(db *gorm.DB) *EquipmentTypeRepository {
	return &EquipmentTypeRepository{
		Db: db,
	}
}

func (r *EquipmentTypeRepository) Create(equipmentType *model.EquipmentType) error {
	err := r.Db.Create(equipmentType).Error
	if err != nil {
		log.Println("failed to create equipment type:", err)
		return err
	}

	log.Println("equipment type created:", equipmentType.ID)

	return nil
}

func (r *EquipmentTypeRepository) GetAll(limit int) ([]model.EquipmentType, error) {
	var equipmentTypes []model.EquipmentType

	err := r.Db.Limit(limit).Find(&equipmentTypes).Error
	if err != nil {
		log.Println("failed to get equipment types:", err)
		return nil, err
	}

	return equipmentTypes, nil
}

func (r *EquipmentTypeRepository) GetById(id uint) (*model.EquipmentType, error) {
	var equipmentType model.EquipmentType

	err := r.Db.First(&equipmentType, "id = ?", id).Error
	if err != nil {
		log.Println("failed to get equipment type:", id)
		return nil, err
	}

	return &equipmentType, nil
}

func (r *EquipmentTypeRepository) Update(equipmentType *model.EquipmentType) error {
	if err := r.Db.Save(equipmentType).Error; err != nil {
		log.Println("failed to update equipment type:", err)
		return err
	}

	log.Println("equipment type updated:", equipmentType.ID)
	return nil
}

func (r EquipmentTypeRepository) Exists(id uint) (bool, error) {
	var count int64

	err := r.Db.Model(&model.EquipmentType{}).
		Where("id = ?", id).
		Count(&count).Error

	if err != nil {
		log.Println("failed to check equipment type:", id)
		return false, err
	}

	return count > 0, nil
}