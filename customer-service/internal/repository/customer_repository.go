package repository

import (
	"customer-service/internal/model"
	"log"
    
	"gorm.io/gorm"
)

type CustomerRepository struct {
	Db *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) *CustomerRepository {
	return &CustomerRepository{
		Db: db,
	}
}

func (r *CustomerRepository) Create(customer *model.Customer) error {
	err := r.Db.Create(customer).Error
	if err != nil {
		log.Println("failed to create customer:", err)
		return err
	}

	log.Println("customer created:", customer.ID)

	return nil
}

func (r *CustomerRepository) GetAll(offset int, limit int) ([]model.Customer, int64, error) {
	var customer []model.Customer
	var total int64

	if err := r.Db.Model(&model.Customer{}).Count(&total).Error; err != nil {
		log.Println("failed to count customers:", err)
		return nil, 0, err
	}

	err := r.Db.Offset(offset).Limit(limit).Find(&customer).Error
	if err != nil {
		log.Println("failed to get customers:", err)
		return nil, 0, err
	}

	log.Println("customers fetched:", len(customer))

	return customer, total, nil
}

func (r *CustomerRepository) GetById(id uint) (*model.Customer, error) {

	var customer model.Customer

	err := r.Db.First(&customer, id).Error
	if err != nil {
		log.Println("failed to get customer:", id)
		return nil, err
	}

	log.Println("customer found:", id)

	return &customer, nil
}

func (r *CustomerRepository) GetByEmail(email string) (*model.Customer, error) {

	var customer model.Customer

	err := r.Db.Where("email = ?", email).First(&customer).Error
	if err != nil {
		return nil, err
	}

	log.Println("customer found by email:", email)

	return &customer, nil
}

func (r *CustomerRepository) Update(customer *model.Customer) error {

	err := r.Db.Save(customer).Error
	if err != nil {
		log.Println("failed to update customer:", customer.ID)
		return err
	}

	log.Println("customer updated:", customer.ID)

	return nil
}

func (r *CustomerRepository) UpdateFields(id uint, data map[string]interface{}) error {
	return r.Db.Model(&model.Customer{}).
		Where("id = ?", id).
		Updates(data).Error
}

func (r CustomerRepository) Exists(id uint) (bool, error) {

	var count int64

	err := r.Db.Model(&model.Customer{}).
		Where("id = ?", id).
		Count(&count).Error

	if err != nil {
		log.Println("failed to check customer:", id)
		return false, err
	}

	return count > 0, nil
}

func (r *CustomerRepository) Delete(id uint) error {
	err := r.Db.Model(&model.Customer{}).
		Where("id = ?", id).
		Update("status", "INACTIVE").Error

	if err != nil {
		log.Println("failed to deactivate customer:", id)
		return err
	}

	log.Println("customer status changed to INACTIVE:", id)
	return nil
}


func (r *CustomerRepository) VerifyKYC(id uint) error {
	err := r.Db.Model(&model.Customer{}).
		Where("id = ?", id).
		Update("kyc_status", "VERIFIED").Error

	if err != nil {
		log.Printf("repository error: failed to update kyc status for customer %d: %v", id, err)
		return err
	}

	log.Printf("repository success: kyc status updated to VERIFIED for customer %d", id)
	return nil
}

func (r *CustomerRepository) UpdateDeposit(id uint, deposit float64) error {
	err := r.Db.Model(&model.Customer{}).
		Where("id = ?", id).
		Update("deposit", deposit).Error

	if err != nil {
		log.Println("failed to update deposit in database:", id)
		return err
	}

	log.Println("database deposit updated:", id)
	return nil
}


func (r CustomerRepository) UpdateBy(id uint, updates map[string]interface{}) error {

	err := r.Db.Model(&model.Customer{}).
		Where("id = ?", id).
		Updates(updates).Error

	if err != nil {
		log.Println("failed to update customer fields:", id)
		return err
	}

	log.Println("customer fields updated:", id)

	return nil
}

func (r *CustomerRepository) DeductDeposit(id uint, amount float64) error {
	var customer model.Customer

	log.Println("deducting deposit for customer:", id)

	err := r.Db.First(&customer, id).Error
	if err != nil {
		log.Println("customer not found:", err)
		return err
	}

	if amount <= customer.Deposit {
		customer.Deposit -= amount
	} else {
		remaining := amount - customer.Deposit

		customer.Deposit = 0
		customer.OutStandingDue += remaining
	}

	err = r.Db.Save(&customer).Error
	if err != nil {
		log.Println("failed to update customer:", err)
		return err
	}

	log.Println("deposit updated for customer:", id)

	return nil
}

