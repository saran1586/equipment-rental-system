package repository

import (
	"log"

	"rental-service/internal/model"

	"gorm.io/gorm"
)

type RentalItemRepository struct {
	Db *gorm.DB
}

func NewRentalItemRepository(db *gorm.DB) *RentalItemRepository {
	return &RentalItemRepository{Db: db}
}

func (r *RentalItemRepository) CreateRentalItem(item *model.RentalItem) error {
	err := r.Db.Create(item).Error
	if err != nil {
		log.Println("rental item create error:", err)
		return err
	}

	log.Println("rental item created:", item.ID)
	return nil
}

func (r *RentalItemRepository) GetRentalItems(rentalID uint) ([]model.RentalItem, error) {
	var items []model.RentalItem

	err := r.Db.Where("rental_id = ?", rentalID).Find(&items).Error
	if err != nil {
		log.Println("rental items get error:", err)
		return nil, err
	}

	return items, nil
}

func (r *RentalItemRepository) UpdateRentalItem(item *model.RentalItem) error {
	err := r.Db.Save(item).Error
	if err != nil {
		log.Println("rental item update error:", err)
		return err
	}

	log.Println("rental item updated:", item.ID)
	return nil
}

func (r *RentalItemRepository) DeleteRentalItem(id uint) error {
	err := r.Db.Delete(&model.RentalItem{}, id).Error
	if err != nil {
		log.Println("failed to delete rental item:", id)
		return err
	}

	log.Println("rental item deleted:", id)

	return nil
}