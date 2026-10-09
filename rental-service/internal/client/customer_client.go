package client

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"bytes"
)

type Customer struct {
	ID             uint    `json:"id"`
	KYCStatus      string  `json:"kyc_status"`
	OutstandingDue float64 `json:"outstanding_due"`
}

type CustomerClient struct {
	BaseURL string
}

func NewCustomerClient(baseURL string) *CustomerClient {
	return &CustomerClient{
		BaseURL: baseURL,
	}
}

func (c *CustomerClient) GetCustomer(id uint) (*Customer, error) {
	url := fmt.Sprintf("%s/customer/%d", c.BaseURL, id)

	resp, err := http.Get(url)
	if err != nil {
		log.Println("failed to get customer:", err)
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("customer not found:", id)
		return nil, fmt.Errorf("customer not found")
	}

	var customer Customer

	if err := json.NewDecoder(resp.Body).Decode(&customer); err != nil {
		log.Println("failed to decode customer response:", err)
		return nil, err
	}

	log.Println("customer fetched:", id)

	return &customer, nil
}

func (c *CustomerClient) UpdateDeposit(id uint, damage float64) error {
	url := fmt.Sprintf("%s/customer/%d/UpdateDeposit", c.BaseURL, id)

	body := []byte(fmt.Sprintf(`{"deposit":%f}`, damage))

	req, err := http.NewRequest(
		http.MethodPatch,
		url,
		bytes.NewBuffer(body),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("failed to update customer deposit")
	}

	return nil
}

func (c *CustomerClient) AddOutstandingDue(id uint, amount float64) error {

	url := fmt.Sprintf("%s/customer/%d/outstanding", c.BaseURL, id)

	body := []byte(fmt.Sprintf(`{"amount":%f}`, amount))

	req, err := http.NewRequest(
		http.MethodPatch,
		url,
		bytes.NewBuffer(body),
	)
	if err != nil {
		log.Println("failed to create outstanding request:", err)
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("failed to update outstanding due:", err)
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		log.Println("failed to update outstanding due for customer:", id)
		return fmt.Errorf("failed to update outstanding due")
	}

	log.Println(
		"outstanding due updated:",
		"customer:", id,
		"amount:", amount,
	)

	return nil
}