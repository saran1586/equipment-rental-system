package service

import (
	"log"
    "os"
	"errors"
	"rental-service/internal/model"
	"rental-service/internal/repository"
)

type RentalItemService struct {
	Repo *repository.RentalItemRepository
}

func NewRentalItemService(repo *repository.RentalItemRepository) *RentalItemService {
	return &RentalItemService{
		Repo: repo,
	}
}

func (s *RentalItemService) CreateRentalItem(item *model.RentalItem) error {
	log.Println("creating rental item")

	if err := s.Repo.CreateRentalItem(item); err != nil {
		log.Println("failed to create rental item:", err)
		return err
	}

	log.Println("rental item created:", item.ID)

	return nil
}

func (s *RentalItemService) GetRentalItems(rentalID uint) ([]model.RentalItem, error) {
	log.Println("getting rental items:", rentalID)

	items, err := s.Repo.GetRentalItems(rentalID)
	if err != nil {
		log.Println("failed to get rental items:", err)
		return nil, err
	}

	return items, nil
}

func (s *RentalItemService) UpdateRentalItem(item *model.RentalItem) error {
	log.Println("updating rental item:", item.ID)

	if err := s.Repo.UpdateRentalItem(item); err != nil {
		log.Println("failed to update rental item:", err)
		return err
	}

	log.Println("rental item updated:", item.ID)

	return nil
}

func (s *RentalItemService) DeleteRentalItem(id uint, password string) error {
	log.Println("delete request for rental item:", id)

	adminPassword := os.Getenv("ADMIN_PASSWORD")

	if password != adminPassword {
		log.Println("invalid admin password for rental item:", id)
		return errors.New("invalid admin password")
	}

	err := s.Repo.DeleteRentalItem(id)
	if err != nil {
		log.Println("failed to delete rental item:", err)
		return err
	}

	log.Println("rental item deleted:", id)

	return nil
}