package repository

import (
	"rental-service/internal/model"

	"gorm.io/gorm"
)

type RentalDueRepository struct {
	DB *gorm.DB
}

func NewRentalDueRepository(db *gorm.DB) *RentalDueRepository {
	return &RentalDueRepository{
		DB: db,
	}
}

func (r *RentalDueRepository) Create(due *model.RentalDue) error {
	return r.DB.Create(due).Error
}

func (r *RentalDueRepository) GetByID(id uint) (*model.RentalDue, error) {
	var due model.RentalDue

	if err := r.DB.First(&due, id).Error; err != nil {
		return nil, err
	}

	return &due, nil
}

func (r *RentalDueRepository) GetPendingByRental(rentalID uint) ([]model.RentalDue, error) {
	var dues []model.RentalDue

	err := r.DB.Where(
		"rental_id = ? AND status = ?",
		rentalID,
		"PENDING",
	).Find(&dues).Error

	return dues, err
}

func (r *RentalDueRepository) Update(due *model.RentalDue) error {
	return r.DB.Save(due).Error
}

