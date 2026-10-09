package repository

import (
	"log"

	"rental-service/internal/model"
  "time"
	"gorm.io/gorm"
)

type RentalRepository struct {
	Db *gorm.DB
}

func NewRentalRepository(db *gorm.DB) *RentalRepository {
	return &RentalRepository{
		Db: db,
	}
}

func (r *RentalRepository) CreateRental(rental *model.Rental) error {
	err := r.Db.Create(rental).Error
	if err != nil {
		log.Println("failed to create rental:", err)
		return err
	}

	log.Println("rental created:", rental.ID)

	return nil
}



func (r *RentalRepository) CreateRentalItem(item *model.RentalItem) error {
	err := r.Db.Create(item).Error
	if err != nil {
		log.Println("failed to create rental item:", err)
		return err
	}

	log.Println("rental item created:", item.ID)

	return nil
}

func (r *RentalRepository) GetAllRental() ([]model.Rental, error) {
	var rental []model.Rental

	err := r.Db.Find(&rental).Error
	if err != nil {
		log.Println("failed to get rentals:", err)
		return nil, err
	}

	return rental, nil
}

func (r *RentalRepository) GetRentalById(id uint) (*model.Rental, error) {
	var rental model.Rental

	err := r.Db.First(&rental, id).Error
	if err != nil {
		log.Println("failed to get rental:", id)
		return nil, err
	}

	return &rental, nil
}

func (r *RentalRepository) UpdateRental(rental *model.Rental) error {
	err := r.Db.Save(rental).Error
	if err != nil {
		log.Println("failed to update rental:", rental.ID)
		return err
	}

	log.Println("rental updated:", rental.ID)

	return nil
}

func (r *RentalRepository) GetRentalItem(rentalID uint) ([]model.RentalItem, error) {
	var item []model.RentalItem

	err := r.Db.Where("rental_id = ?", rentalID).Find(&item).Error
	if err != nil {
		log.Println("failed to get rental items:", rentalID)
		return nil, err
	}

	return item, nil
}

func (r *RentalRepository) Delete(id uint) error {
	err := r.Db.Delete(&model.Rental{}, id).Error
	if err != nil {
		log.Println("failed to delete rental:", id)
		return err
	}

	log.Println("rental deleted:", id)

	return nil
}

func (r *RentalRepository) UpdateRentalItem(item *model.RentalItem) error {
	return r.Db.Save(item).Error
}



func (r *RentalRepository) GetOverdueRentals() ([]model.Rental, error) {
	var rentals []model.Rental

	log.Println("checking for overdue rentals")

	err := r.Db.
		Where("status IN ?", []string{"ACTIVE", "OVERDUE"}).
		Where("due_date < ?", time.Now()).
		Find(&rentals).Error

	if err != nil {
		log.Println("failed to get overdue rentals:", err)
		return nil, err
	}

	log.Println("overdue rentals found:", len(rentals))

	return rentals, nil
}

func (r *RentalRepository) GetActiveRentals() ([]model.Rental, error) {
	var rentals []model.Rental

	err := r.Db.Where(
		"status IN ?",
		[]string{"ACTIVE", "OVERDUE"},
	).Find(&rentals).Error

	return rentals, err
}