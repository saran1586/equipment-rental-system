
package service

import (
	"errors"
	"log"
	"strings"

	"equipment-service/internal/model"
	"equipment-service/internal/repository"
)

type EquipmentUnitService struct {
	EquipmentTypeRepo *repository.EquipmentTypeRepository
	UnitRepo          *repository.EquipmentUnitRepository
}

func NewEquipmentUnitService(
	repo *repository.EquipmentUnitRepository,
	equipmentTypeRepo *repository.EquipmentTypeRepository,
) *EquipmentUnitService {
	return &EquipmentUnitService{
		UnitRepo:          repo,
		EquipmentTypeRepo: equipmentTypeRepo,
	}
}

func (s *EquipmentUnitService) Create(equipmentTypeID uint,serialNumber string,) (*model.EquipmentUnit, error) {

	if equipmentTypeID == 0 {
		return nil, errors.New("equipment type id is required")
	}

	if strings.TrimSpace(serialNumber) == "" {
		return nil, errors.New("serial number is required")
	}
	

	exists, err := s.EquipmentTypeRepo.Exists(equipmentTypeID)
	if err != nil {
		log.Println("failed to check equipment type:", err)
		return nil, err
	}

	if !exists {
		return nil, errors.New("equipment type not found")
	}
exists, err = s.UnitRepo.Exists(serialNumber)
if err != nil {
	log.Println("failed to check serial number:", err)
	return nil, err
}

if exists {
	return nil, errors.New("serial number already exists")
}

	
	unit := &model.EquipmentUnit{
		EquipmentTypeId: equipmentTypeID,
		SerialNumber:    serialNumber,
		Condition:       model.ConditionNew,
		Status:          model.StatusAvailable,
	}

	if err := s.UnitRepo.Create(unit); err != nil {
		log.Println("failed to create equipment unit:", err)
		return nil, err
	}

	log.Println("equipment unit created:", unit.ID)

	return unit, nil
}

func (s *EquipmentUnitService) GetById(id uint) (*model.EquipmentUnit, error) {

	if id == 0 {
		return nil, errors.New("equipment unit id is required")
	}

	unit, err := s.UnitRepo.GetById(id)
	if err != nil {
		log.Println("failed to get equipment unit:", err)
		return nil, err
	}

	return unit, nil
}

func (s *EquipmentUnitService) Update(id uint,serialNumber string,condition model.EquipmentCondition,status model.EquipmentStatus,) (*model.EquipmentUnit, error) {

	if id == 0 {
		return nil, errors.New("equipment unit id is required")
	}

	if strings.TrimSpace(serialNumber) == "" {
		return nil, errors.New("serial number is required")
	}

	unit, err := s.UnitRepo.GetById(id)
	if err != nil {
		return nil, errors.New("equipment unit not found")
	}

	if unit.SerialNumber != serialNumber {
		exists, err := s.UnitRepo.Exists(serialNumber)
		if err != nil {
			log.Println("failed to check serial number:", err)
			return nil, err
		}

		if exists {
			return nil, errors.New("serial number already exists")
		}
	}

	unit.SerialNumber = serialNumber
	unit.Condition = condition
	unit.Status = status

	if err := s.UnitRepo.Update(unit); err != nil {
		log.Println("failed to update equipment unit:", err)
		return nil, err
	}

	log.Println("equipment unit updated:", id)

	return unit, nil
}

func (s *EquipmentUnitService) Delete(id uint) error {

	if id == 0 {
		return errors.New("equipment unit id is required")
	}

	unit, err := s.UnitRepo.GetById(id)
	if err != nil {
		return errors.New("equipment unit not found")
	}

	if unit.Status == model.StatusRented {
		return errors.New("rented equipment unit cannot be deactivated")
	}

	unit.Status = model.StatusInactive

	if err := s.UnitRepo.Deactivate(unit); err != nil {
		log.Println("failed to deactivate equipment unit:", err)
		return err
	}

	log.Println("equipment unit deactivated:", id)

	return nil
}

func (s *EquipmentUnitService) UpdateStatus(id uint,status model.EquipmentStatus,) (*model.EquipmentUnit, error) {

	unit, err := s.UnitRepo.GetById(id)
	if err != nil {
		return nil, err
	}

	unit.Status = status

	if err := s.UnitRepo.Update(unit); err != nil {
		return nil, err
	}

	log.Println("equipment status updated:", id, status)

	return unit, nil
}


func (s *EquipmentUnitService) GetAll(limit int) ([]model.EquipmentUnit, error) {
	return s.UnitRepo.GetAll(limit)
}


func (s *EquipmentUnitService) UpdateCondition(id uint,condition model.EquipmentCondition,) (*model.EquipmentUnit, error) {

	if id == 0 {
		return nil, errors.New("equipment unit id is required")
	}

	unit, err := s.UnitRepo.GetById(id)
	if err != nil {
		log.Println("equipment unit not found:", id)
		return nil, errors.New("equipment unit not found")
	}

	if condition != model.ConditionGood &&
		condition != model.ConditionDamaged &&
		condition != model.ConditionNew {
		return nil, errors.New("invalid equipment condition")
	}

	unit.Condition = condition

	if err := s.UnitRepo.Update(unit); err != nil {
		log.Println("failed to update equipment condition:", err)
		return nil, err
	}

	log.Println("equipment condition updated:", id, condition)

	return unit, nil
}