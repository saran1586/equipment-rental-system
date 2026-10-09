package service

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"equipment-service/internal/model"
	"equipment-service/internal/repository"
)

type EquipmentTypeService struct {
	Repo        *repository.EquipmentTypeRepository
	UnitService *EquipmentUnitService
}

func NewEquipmentTypeService(repo *repository.EquipmentTypeRepository, unitService *EquipmentUnitService) *EquipmentTypeService {

	return &EquipmentTypeService{
		Repo:        repo,
		UnitService: unitService,
	}
}

func (r *EquipmentTypeService) Create(equipmentType *model.EquipmentType) error {

	if strings.TrimSpace(equipmentType.Name) == "" {
		return errors.New("equipment type name is required")
	}

	if equipmentType.DailyRate <= 0 {
		return errors.New("daily rate must be greater than zero")
	}

	if equipmentType.Quantity <= 0 {
		return errors.New("quantity must be greater than zero")
	}

	if err := r.Repo.Create(equipmentType); err != nil {
		log.Println("failed to create equipment type:", err)
		return err
	}

	for i := 1; i <= equipmentType.Quantity; i++ {

		serialNumber := fmt.Sprintf("%s-%03d", strings.ToUpper(equipmentType.Name), i)

		_, err := r.UnitService.Create(equipmentType.ID, serialNumber)

		if err != nil {
			return err
		}
	}

	log.Println("equipment type created:", equipmentType.ID)

	return nil
}

func (s *EquipmentTypeService) GetAll(limit int) ([]model.EquipmentType, error) {
	return s.Repo.GetAll(limit)
}

func (r *EquipmentTypeService) GetById(id uint) (*model.EquipmentType, error) {

	equipmentType, err := r.Repo.GetById(id)
	if err != nil {
		log.Println("failed to get equipment type:", err)
		return nil, err
	}

	return equipmentType, nil
}

func (r *EquipmentTypeService) Update(id uint,name string,description string,dailyRate int64,	quantity int,
) (*model.EquipmentType, error) {

	equipmentType, err := r.Repo.GetById(id)
	if err != nil {
		log.Println("equipment type not found:", id)
		return nil, err
	}

	log.Println("old quantity:", equipmentType.Quantity)
	log.Println("new quantity:", quantity)

	if strings.TrimSpace(name) == "" {
		return nil, errors.New("equipment type name is required")
	}

	if dailyRate <= 0 {
		return nil, errors.New("daily rate must be greater than zero")
	}

	if quantity <= 0 {
		return nil, errors.New("quantity must be greater than zero")
	}

	oldQuantity := equipmentType.Quantity

	equipmentType.Name = name
	equipmentType.Description = description
	equipmentType.DailyRate = dailyRate
	equipmentType.Quantity = quantity

	if err := r.Repo.Update(equipmentType); err != nil {
		log.Println("failed to update equipment type:", err)
		return nil, err
	}

	if quantity > oldQuantity {

		for i := oldQuantity + 1; i <= quantity; i++ {

			serialNumber := fmt.Sprintf("%s-%03d",strings.ToUpper(name),i,)

			_, err := r.UnitService.Create(
				equipmentType.ID,
				serialNumber,
			)

			if err != nil {
				log.Println("failed to create equipment unit:", err)
				return nil, err
			}

			log.Println("equipment unit created:", serialNumber)
		}
	}

	log.Println("equipment type updated:", id)

	return equipmentType, nil
}