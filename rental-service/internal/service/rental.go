package service

import (
	"errors"
	"log"
	"rental-service/internal/client"
	"rental-service/internal/model"
	"rental-service/internal/repository"
	"strings"
	"time"
)

type RentalService struct {
	Repo            *repository.RentalRepository
	RentalDueRepo *repository.RentalDueRepository
	CustomerClient  *client.CustomerClient
	EquipmentClient *client.EquipmentClient
}

func NewRentalService(repo *repository.RentalRepository,rentalDueRepo *repository.RentalDueRepository, customerClient *client.CustomerClient, equipmentClient *client.EquipmentClient) *RentalService {
	return &RentalService{
		Repo:            repo,
		RentalDueRepo: rentalDueRepo,
		CustomerClient:  customerClient,
		EquipmentClient: equipmentClient,
	}
}
func (s *RentalService) CreateRental(rental *model.Rental,equipment []uint,) error {

	if err := s.ValidateCustomer(rental.CustomerID); err != nil {
		log.Println("customer validation failed:", err)
		return err
	}

	if err := s.Repo.CreateRental(rental); err != nil {
		log.Println("failed to create rental:", err)
		return err
	}

	log.Println("RENTAL CREATED ID:", rental.ID)

	if rental.ID == 0 {
		return errors.New("rental id was not generated")
	}

	for _, equipmentID := range equipment {

		item := model.RentalItem{
			RentalID:       rental.ID,
			EquipmentUnitID: equipmentID,
			IssueCondition: "GOOD",
			RentPerDay:     0,
			DamageCharge:   0,
		}

		log.Println(
			"CREATING RENTAL ITEM",
			"rental_id:", item.RentalID,
			"equipment_id:", item.EquipmentUnitID,
		)

		if err := s.Repo.CreateRentalItem(&item); err != nil {
			log.Println("failed to create rental item:", err)
			return err
		}

		log.Println(
			"RENTAL ITEM CREATED",
			"item_id:", item.ID,
			"rental_id:", item.RentalID,
		)
	}

	log.Println("rental created successfully:", rental.ID)

	return nil
}

func (s *RentalService) ValidateCustomer(customerID uint) error {
	customer, err := s.CustomerClient.GetCustomer(customerID)
	if err != nil {
		log.Println("failed to get customer:", err)
		return err
	}

	if customer.KYCStatus != "VERIFIED" {
		return errors.New("customer KYC is not verified")
	}

	if customer.OutstandingDue > 10000 {
		return errors.New("please pay the outstanding due ")
	}

	return nil
}

func (s *RentalService) ValidateEquipment(equipmentIDs []uint) error {
	for _, id := range equipmentIDs {
		equipment, err := s.EquipmentClient.GetEquipmentUnit(id)
		if err != nil {
			log.Printf("WARNING: reached limit pay the due  %d \n", id)
			log.Println("failed to get equipment:", err)
			return err
		}

		if equipment.Status != "AVAILABLE" {
			return errors.New("equipment is not available")
		}

		cond := strings.ToUpper(strings.TrimSpace(equipment.Condition))

		if cond == "DAMAGED" {
			log.Printf("WARNING: equipment unit %d is  damaged\n", id)
			return errors.New("equipment is damaged")
		}

		if cond != "GOOD" && cond != "NEW" {
			return errors.New("equipment is not in rental condition")
		}
	}

	return nil
}

func (s *RentalService) CreateRentalWithItems(rental *model.Rental,equipmentIDs []uint,) error {

	if err := s.ValidateCustomer(rental.CustomerID); err != nil {
		log.Println("customer validation failed:", err)
		return err
	}

	if len(equipmentIDs) == 0 {
		return errors.New("at least one equipment is required")
	}

	if err := s.ValidateEquipment(equipmentIDs); err != nil {
		log.Println("equipment validation failed:", err)
		return err
	}

	if rental.DueDate.Before(rental.RentalDate) {
		return errors.New("due date must be after rental date")
	}

	rentalDays := int(
		rental.DueDate.Sub(rental.RentalDate).Hours() / 24,
	)

	if rentalDays <= 0 {
		rentalDays = 1
	}

	rental.TotalRent = 0

	if err := s.Repo.CreateRental(rental); err != nil {
		log.Println("failed to create rental:", err)
		return err
	}

	log.Println("rental created:", rental.ID)

	for _, equipmentID := range equipmentIDs {

		equipment, err := s.EquipmentClient.GetEquipmentUnit(equipmentID)
		if err != nil {
			log.Println("failed to get equipment unit:", err)
			return err
		}

		equipmentType, err := s.EquipmentClient.GetEquipmentType(
			equipment.EquipmentTypeID,
		)
		if err != nil {
			log.Println("failed to get equipment type:", err)
			return err
		}

		rentPerDay := float64(equipmentType.DailyRate)
		itemTotal := rentPerDay * float64(rentalDays)

		item := model.RentalItem{
			RentalID:        rental.ID,
			EquipmentUnitID: equipmentID,
			IssueCondition:  "GOOD",
			RentPerDay:      rentPerDay,
			DamageCharge:    0,
		}

		if err := s.Repo.CreateRentalItem(&item); err != nil {
			log.Println("failed to create rental item:", err)
			return err
		}

		if err := s.EquipmentClient.UpdateEquipmentStatus(
			equipmentID,
			"RENTED",
		); err != nil {
			log.Println("failed to update equipment status:", err)
			return err
		}

		rental.TotalRent += itemTotal

		log.Println(
			"rental item calculated:",
			"equipment:", equipmentID,
			"daily rate:", rentPerDay,
			"days:", rentalDays,
			"total:", itemTotal,
		)
	}

	if err := s.Repo.UpdateRental(rental); err != nil {
		log.Println("failed to update rental total:", err)
		return err
	}

	dailyRent := rental.TotalRent / float64(rentalDays)

	log.Println(
		"rental total calculated:",
		"rental:", rental.ID,
		"total:", rental.TotalRent,
		"days:", rentalDays,
		"daily rent:", dailyRent,
	)

	log.Println(
		"rental created successfully:",
		rental.ID,
	)

	return nil
}

func (r *RentalService) GetRentalById(id uint) (*model.Rental, error) {
	rental, err := r.Repo.GetRentalById(id)
	if err != nil {
		log.Println("failed to get rental:", err)
		return nil, err
	}

	return rental, nil
}

func (r *RentalService) GetAllRental() ([]model.Rental, error) {
	rentals, err := r.Repo.GetAllRental()
	if err != nil {
		log.Println("failed to get rentals:", err)
		return nil, err
	}

	return rentals, nil
}

func (r *RentalService) UpdateRental(rental *model.Rental) error {
	if err := r.Repo.UpdateRental(rental); err != nil {
		log.Println("failed to update rental:", err)
		return err
	}

	log.Println("rental updated:", rental.ID)

	return nil
}

func (r *RentalService) GetRentalItem(rentalID uint) ([]model.RentalItem, error) {
	items, err := r.Repo.GetRentalItem(rentalID)
	if err != nil {
		log.Println("failed to get rental items:", err)
		return nil, err
	}

	return items, nil
}

func (s *RentalService) Delete(id uint) error {
	if err := s.Repo.Delete(id); err != nil {
		log.Println("failed to delete rental:", err)
		return err
	}

	log.Println("rental deleted:", id)

	return nil
}
func (s *RentalService) ReturnRental(rentalID uint,condition string,damageCharge float64,) error {

	rental, err := s.Repo.GetRentalById(rentalID)
	if err != nil {
		log.Println("failed to get rental:", err)
		return errors.New("rental not found")
	}

	if rental.Status == "COMPLETED" {
		return errors.New("rental already returned")
	}

	items, err := s.Repo.GetRentalItem(rentalID)
	if err != nil {
		log.Println("failed to get rental items:", err)
		return err
	}

	condition = strings.ToUpper(strings.TrimSpace(condition))

	if condition == "" {
		return errors.New("return condition is required")
	}

	if condition == "GOOD" || condition == "NEW" {
		damageCharge = 0
	}

	if damageCharge < 0 {
		return errors.New("damage charge cannot be negative")
	}

	for _, item := range items {
		item.ReturnCondition = condition
		item.DamageCharge = damageCharge

		if err := s.Repo.UpdateRentalItem(&item); err != nil {
			log.Println("failed to update rental item:", err)
			return err
		}

		if err := s.EquipmentClient.UpdateEquipmentStatus(
			item.EquipmentUnitID,
			"AVAILABLE",
		); err != nil {
			log.Println("failed to update equipment status:", err)
			return err
		}

		log.Println("equipment returned:", item.EquipmentUnitID)
	}

	lateFee := rental.LateFee
	extraAmount := lateFee + damageCharge

	if extraAmount > 0 {
		if err := s.CustomerClient.AddOutstandingDue(
			rental.CustomerID,
			extraAmount,
		); err != nil {
			log.Println("failed to add outstanding due:", err)
			return err
		}

		log.Println(
			"extra amount added to outstanding:",
			"customer:", rental.CustomerID,
			"late fee:", lateFee,
			"damage charge:", damageCharge,
			"total:", extraAmount,
		)
	}

	rental.Status = "COMPLETED"
	rental.DamageFee = damageCharge

	returnDate := time.Now()
	rental.ReturnDate = &returnDate

	if err := s.Repo.UpdateRental(rental); err != nil {
		log.Println("failed to update rental:", err)
		return err
	}

	log.Println("rental returned successfully:", rentalID)

	return nil
}
func (s *RentalService) ProcessDailyRentalDue() error {
	log.Println("starting daily rental due check")

	rentals, err := s.Repo.GetActiveRentals()
	if err != nil {
		log.Println("failed to get active rentals:", err)
		return err
	}

	for _, rental := range rentals {

		if rental.Status != "ACTIVE" && rental.Status != "OVERDUE" {
			continue
		}

		now := time.Now()
/*
		if rental.LastDueDate != nil {
			lastDue := rental.LastDueDate.Format("2006-01-02")
			today := now.Format("2006-01-02")

			if lastDue == today {
				log.Println(
					"daily due already added:",
					"rental:", rental.ID,
				)
				continue
			}
		}*/

		if rental.LastDueDate != nil {
	lastDue := rental.LastDueDate.Format("2006-01-02 15:04")
	currentMinute := now.Format("2006-01-02 15:04")

	if lastDue == currentMinute {
		log.Println(
			"daily due already added this minute:",
			"rental:", rental.ID,
		)
		continue
	}
}

		totalDays := int(
			rental.DueDate.Sub(rental.RentalDate).Hours() / 24,
		)

		if totalDays <= 0 {
			totalDays = 1
		}

		dailyRent := rental.TotalRent / float64(totalDays)

		customer, err := s.CustomerClient.GetCustomer(rental.CustomerID)
		if err != nil {
			log.Println(
				"failed to get customer:",
				rental.CustomerID,
				err,
			)
			continue
		}

		penalty := 0.0

		if customer.OutstandingDue > 0 {
			penalty = customer.OutstandingDue * 0.02

			log.Println(
				"WARNING: previous due not paid",
				"rental:", rental.ID,
				"customer_id:", rental.CustomerID,
				"previous_due:", customer.OutstandingDue,
				"penalty:", penalty,
			)
		}

		dueAmount := dailyRent + penalty

		due := &model.RentalDue{
			RentalID:   rental.ID,
			CustomerID: rental.CustomerID,
			DueDate:    now,
			DailyRent:  dailyRent,
			Penalty:    penalty,
			TotalDue:   dueAmount,
			Status:     "PENDING",
		}

		if err := s.RentalDueRepo.Create(due); err != nil {
			log.Println(
				"failed to save rental due:",
				"rental:", rental.ID,
				"error:", err,
			)
			continue
		}

		if err := s.CustomerClient.AddOutstandingDue(
			rental.CustomerID,
			dueAmount,
		); err != nil {
			log.Println(
				"failed to add outstanding due:",
				"rental:", rental.ID,
				"customer:", rental.CustomerID,
				"error:", err,
			)
			continue
		}

		rental.LastDueDate = &now

		if penalty > 0 {
			rental.Status = "OVERDUE"
		}

		if err := s.Repo.UpdateRental(&rental); err != nil {
			log.Println(
				"failed to update rental:",
				rental.ID,
				err,
			)
			continue
		}

		log.Println(
			"ATTENTION: RENTAL DUE ADDED",
			"rental:", rental.ID,
			"customer_id:", rental.CustomerID,
			"daily_rent:", dailyRent,
			"penalty:", penalty,
			"due_added:", dueAmount,
		)

		log.Println(
			"WARNING: PAY THE DUE",
			"rental:", rental.ID,
			"customer_id:", rental.CustomerID,
			"total_outstanding:",
			customer.OutstandingDue+dueAmount,
		)
	}

	log.Println("daily rental due check completed")

	return nil
}


func (s *RentalService) MarkOverdueRentals() error {

	log.Println("starting overdue rental check")

	rentals, err := s.Repo.GetOverdueRentals()
	if err != nil {
		log.Println(
			"failed to get overdue rentals:",
			err,
		)
		return err
	}

	for _, rental := range rentals {

		if rental.Status == "COMPLETED" {
			continue
		}

		if time.Now().Before(rental.DueDate) {
			continue
		}

		if rental.Status == "OVERDUE" {
			log.Println(
				"rental already marked overdue:",
				rental.ID,
			)
			continue
		}

		rental.Status = "OVERDUE"

		if err := s.Repo.UpdateRental(&rental); err != nil {
			log.Println(
				"failed to update overdue rental:",
				rental.ID,
				err,
			)
			return err
		}

		log.Println(
			"rental marked overdue:",
			"rental:", rental.ID,
			"customer_id:", rental.CustomerID,
		)
	}

	log.Println("overdue rental check completed")

	return nil
}