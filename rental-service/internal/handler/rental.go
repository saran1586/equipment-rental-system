package handler

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"rental-service/internal/model"
	"rental-service/internal/service"

	"github.com/gin-gonic/gin"
)

type RentalHandler struct {
	Service *service.RentalService
}

func NewRentalHandler(handler *service.RentalService) *RentalHandler {
	return &RentalHandler{
		Service: handler,
	}
}

type CreateRentalRequest struct {
	CustomerID      uint      `json:"customer_id"`
	DueDate         time.Time `json:"due_date"`
	EquipmentUnitID []uint    `json:"equipment_unit_ids"`
}

type ReturnRentalRequest struct {
	Condition    string  `json:"condition"`
	DamageCharge float64 `json:"damage_charge"`
}

type PayRentalDueRequest struct {
	Amount float64 `json:"amount"`
}



func (h *RentalHandler) Create(c *gin.Context) {
	var input CreateRentalRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		log.Println("invalid rental input:", err)

		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid input",
		})
		return
	}

	rental := model.Rental{
		CustomerID: input.CustomerID,
		RentalDate: time.Now(),
		DueDate:    input.DueDate,
		Status:     "ACTIVE",
	}

	err := h.Service.CreateRentalWithItems(
		&rental,
		input.EquipmentUnitID,
	)

	if err != nil {
		log.Println("failed to create rental:", err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	log.Println("rental created:", rental.ID)

	c.JSON(http.StatusCreated, gin.H{
		"data":    rental,
		"message": "rental created",
	})
}

func (r *RentalHandler) GetAll(c *gin.Context) {
	rental, err := r.Service.GetAllRental()
	if err != nil {
		log.Println("failed to get rentals:", err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to get rentals",
		})
		return
	}

	log.Println("rentals fetched successfully")

	c.JSON(http.StatusOK, gin.H{
		"data": rental,
	})
}

func (h *RentalHandler) GetById(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		log.Println("invalid rental id:", c.Param("id"))

		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid id",
		})
		return
	}

	rental, err := h.Service.GetRentalById(uint(id))
	if err != nil {
		log.Println("failed to get rental:", err)

		c.JSON(http.StatusNotFound, gin.H{
			"message": "rental not found",
		})
		return
	}

	log.Println("rental fetched:", id)

	c.JSON(http.StatusOK, rental)
}

func (h *RentalHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		log.Println("invalid rental id:", c.Param("id"))

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid id",
		})
		return
	}

	if err := h.Service.Delete(uint(id)); err != nil {
		log.Println("failed to delete rental:", err)

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	log.Println("rental deleted:", id)

	c.JSON(http.StatusOK, gin.H{
		"message": "rental deleted",
	})
}


func (h *RentalHandler) Return(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid rental id",
		})
		return
	}

	var input ReturnRentalRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid input",
		})
		return
	}

	err = h.Service.ReturnRental(
		uint(id),
		input.Condition,
		input.DamageCharge,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "rental returned successfully",
	})
}

