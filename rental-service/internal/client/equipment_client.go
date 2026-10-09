package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type EquipmentUnit struct {
	ID              uint   `json:"id"`
	EquipmentTypeID uint   `json:"equipment_type_id"`
	Status          string `json:"status"`
	Condition       string `json:"condition"`
}

type EquipmentClient struct {
	BaseURL string
}

type EquipmentType struct {
	ID        uint  `json:"id"`
	DailyRate int64 `json:"dailyrate"`
}

func NewEquipmentClient(baseURL string) *EquipmentClient {
	return &EquipmentClient{
		BaseURL: baseURL,
	}
}
func (c *EquipmentClient) GetEquipmentUnit(id uint) (*EquipmentUnit, error) {

	url := fmt.Sprintf("%s/equipment-unit/%d", c.BaseURL, id)

	log.Println("equipment unit request:", url)

	resp, err := http.Get(url)
	if err != nil {
		log.Println("equipment unit request failed:", err)
		return nil, err
	}

	defer resp.Body.Close()

	log.Println("equipment unit response:", resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		log.Println("equipment unit not found:", id)
		return nil, fmt.Errorf("equipment unit not found")
	}

	var response struct {
		Data EquipmentUnit `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		log.Println("failed to decode equipment unit:", err)
		return nil, err
	}

	log.Println("equipment unit received:",response.Data.ID,"equipment type:",response.Data.EquipmentTypeID,"status:",response.Data.Status,)

	return &response.Data, nil
}
func (c *EquipmentClient) UpdateEquipmentStatus(id uint, status string) error {
url := fmt.Sprintf("%s/equipment-unit/%d/status", c.BaseURL, id)

	body := []byte(fmt.Sprintf(`{"status":"%s"}`, status))

	req, err := http.NewRequest(http.MethodPatch,url,bytes.NewBuffer(body),)
	if err != nil {
		log.Println("failed to create equipment request:", err)
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("failed to update equipment status:", err)
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		log.Println("equipment status update failed:", id)
		return fmt.Errorf("failed to update equipment status")
	}

	log.Println("equipment status updated:", id, status)

	return nil
}

func (c *EquipmentClient) GetEquipmentType(id uint) (*EquipmentType, error) {
	
url := fmt.Sprintf("%s/internal/equipment-type/%d", c.BaseURL, id)

	log.Println("equipment type request:", url)

	resp, err := http.Get(url)
	if err != nil {
		log.Println("equipment type request failed:", err)
		return nil, err
	}

	defer resp.Body.Close()

	log.Println("equipment type response:", resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		log.Println("equipment type not found:", id)
		return nil, fmt.Errorf("equipment type not found")
	}

	var response struct {
		Data EquipmentType `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		log.Println("failed to decode equipment type:", err)
		return nil, err
	}

	log.Println(
		"equipment type received:",
		response.Data.ID,
		"daily rate:",
		response.Data.DailyRate,
	)

	return &response.Data, nil
}