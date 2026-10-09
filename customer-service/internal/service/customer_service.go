package service

import (
	"customer-service/internal/model"
	"customer-service/internal/repository"
	"errors"
	"log"
	"net/mail"
	"strings"
)

type CustomerService struct {
	Repo *repository.CustomerRepository
}


func NewCustomerService(repo *repository.CustomerRepository) *CustomerService {
	return &CustomerService{Repo: repo}
}

func (r *CustomerService) Create(customer *model.Customer) error {

	existingCustomer, err := r.Repo.GetByEmail(customer.Email)

	if err == nil && existingCustomer != nil {
		log.Println("customer email already exists:", customer.Email)
		return errors.New("customer with this email already exists")
	}

	customer.Email = strings.TrimSpace(strings.ToLower(customer.Email))

	if _, err := mail.ParseAddress(customer.Email); err != nil {
		log.Println("invalid email:", customer.Email)
		return errors.New("invalid email entered")
	}

	customer.KYCStatus = "PENDING"
	customer.Deposit = 0
	customer.OutStandingDue = 0
	customer.Status = "ACTIVE"

	if err := r.Repo.Create(customer); err != nil {
		log.Println("failed to create customer:", err)
		return err
	}

	log.Println("customer created:", customer.ID)

	return nil
}

func (r *CustomerService) GetAll(page int, limit int) ([]model.Customer, int64, error) {

	if page < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	customers, total, err := r.Repo.GetAll(offset, limit)
	if err != nil {
		log.Println("failed to get customers:", err)
		return nil, 0, err
	}

	log.Println("customers fetched, total:", total)

	return customers, total, nil
}

func (r *CustomerService) GetById(id uint) (*model.Customer, error) {

	customer, err := r.Repo.GetById(id)
	if err != nil {
		log.Println("failed to get customer:", id)
		return nil, err
	}

	log.Println("customer fetched:", id)

	return customer, nil
}

func (s *CustomerService) Update(customer *model.Customer) error {
	oldCustomer, err := s.Repo.GetById(customer.ID)
	if err != nil {
		return errors.New("customer not found")
	}

	data := map[string]interface{}{}

	if customer.Name != "" {
		data["name"] = customer.Name
	}

	if customer.Email != "" {
		_, err := mail.ParseAddress(customer.Email)
		if err != nil {
			return errors.New("invalid email format")
		}

		if strings.EqualFold(customer.Email, oldCustomer.Email) {
			return errors.New("email is already the same")
		}

		data["email"] = customer.Email
	}

	if len(data) == 0 {
		return errors.New("user should give some data")
	}

	err = s.Repo.UpdateFields(customer.ID, data)
	if err != nil {
		return err
	}

	log.Println("customer updated:", customer.ID)

	return nil
}

func (s *CustomerService) Deactivate(id uint) error {

	customer, err := s.Repo.GetById(id)
	if err != nil {
		return errors.New("customer not found")
	}

	if customer.Status == "INACTIVE" {
		return errors.New("customer already inactive")
	}

	if customer.OutStandingDue > 0 {
		return errors.New("customer has outstanding due")
	}

	if customer.Deposit > 0 {
		return errors.New("money is still in your wallet")
	}

	err = s.Repo.Delete(id)
	if err != nil {
		return err
	}

	log.Println("customer deactivated:", id)

	return nil
}

func (s *CustomerService) VerifyKYC(id uint, inputEmail string) error {

	customer, err := s.Repo.GetById(id)
	if err != nil {
		return errors.New("customer not found")
	}

	inputEmail = strings.TrimSpace(strings.ToLower(inputEmail))
	email := strings.TrimSpace(strings.ToLower(customer.Email))

	if inputEmail != email {
		return errors.New("enter the correct email")
	}

	if customer.Deposit < 2000 {
		return errors.New("please deposit minimum 2000")
	}

	customer.KYCStatus = "VERIFIED"

	err = s.Repo.Update(customer)
	if err != nil {
		return err
	}

	log.Println("KYC verified:", id)

	return nil
}

func (s *CustomerService) UpdateDeposit(id uint, deposit float64) error {
	if deposit < 0 {
		log.Println("negative deposit rejected:", id)
		return errors.New("deposit cannot be negative")
	}

	if err := s.Repo.UpdateDeposit(id, deposit); err != nil {
		log.Println("failed to update deposit:", id)
		return err
	}

	log.Println("deposit updated:", id)
	return nil
}

func (s *CustomerService) UpdateBy(id uint, updates map[string]interface{}) error {

	if err := s.Repo.UpdateBy(id, updates); err != nil {
		log.Println("failed to update customer fields:", id)
		return err
	}

	log.Println("customer fields updated:", id)

	return nil
}

func (s *CustomerService) DeductDeposit(id uint, amount float64) error {
	log.Println("updating customer deposit:", id)

	err := s.Repo.DeductDeposit(id, amount)
	if err != nil {
		log.Println("failed to update customer deposit:", err)
		return err
	}

	log.Println("customer deposit updated:", id)

	return nil
}

func (s *CustomerService) PayOutstanding(id uint, amount float64) (*model.Customer, error) {

	if id == 0 {
		return nil, errors.New("customer id is required")
	}

	if amount <= 0 {
		return nil, errors.New("payment amount must be greater than zero")
	}

	customer, err := s.Repo.GetById(id)
	if err != nil {
		log.Println("customer not found:", id)
		return nil, errors.New("customer not found")
	}

	if customer.OutStandingDue <= 0 {
		return nil, errors.New("no outstanding due")
	}

	if amount > customer.OutStandingDue {
		return nil, errors.New("payment amount exceeds outstanding due")
	}

	customer.OutStandingDue -= amount

	if err := s.Repo.Update(customer); err != nil {
		log.Println("failed to update outstanding due:", err)
		return nil, err
	}

	log.Println("outstanding due updated:",id,"paid:",amount,"remaining:",customer.OutStandingDue,)

	return customer, nil
}

func (s *CustomerService) AddOutstandingDue(id uint, amount float64) (*model.Customer, error) {

	if id == 0 {
		return nil, errors.New("customer id is required")
	}

	if amount <= 0 {
		return nil, errors.New("amount must be greater than zero")
	}

	customer, err := s.Repo.GetById(id)
	if err != nil {
		log.Println("customer not found:", id)
		return nil, errors.New("customer not found")
	}

	customer.OutStandingDue += amount

	if err := s.Repo.Update(customer); err != nil {
		log.Println("failed to update outstanding due:", err)
		return nil, err
	}

	log.Println(
		"outstanding due added:",
		id,
		"amount:",
		amount,
		"total:",
		customer.OutStandingDue,
	)

	return customer, nil
}