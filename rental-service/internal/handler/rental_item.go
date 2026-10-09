package handler

import (
	"log"
	"net/http"
	"strconv"

	"rental-service/internal/model"
	"rental-service/internal/service"

	"github.com/gin-gonic/gin"
)

type RentalItemHandler struct {
	Service *service.RentalItemService
}

func NewRentalItemHandler(service *service.RentalItemService) *RentalItemHandler {
	return &RentalItemHandler{Service: service}
}

func (h *RentalItemHandler) Create(c *gin.Context) {
	var item model.RentalItem

	if err := c.ShouldBindJSON(&item); err != nil {
		log.Println("invalid rental item request:", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := h.Service.CreateRentalItem(&item); err != nil {
		log.Println("create rental item failed:", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	log.Println("rental item created:", item.ID)

	c.JSON(http.StatusCreated, gin.H{
		"data":    item,
		"message": "Rental item created",
	})
}

func (h *RentalItemHandler) GetByRentalID(c *gin.Context) {
	rentalID, err := strconv.ParseUint(c.Param("rental_id"), 10, 64)
	if err != nil {
		log.Println("invalid rental id:", c.Param("rental_id"))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid rental_id",
		})
		return
	}

	items, err := h.Service.GetRentalItems(uint(rentalID))
	if err != nil {
		log.Println("get rental items failed:", err)
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	log.Println("rental items fetched for:", rentalID)

	c.JSON(http.StatusOK, gin.H{
		"data": items,
	})
}

func (h *RentalItemHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		log.Println("invalid rental item id:", c.Param("id"))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid id",
		})
		return
	}

	var item model.RentalItem

	if err := c.ShouldBindJSON(&item); err != nil {
		log.Println("invalid update request:", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	item.ID = uint(id)

	if err := h.Service.UpdateRentalItem(&item); err != nil {
		log.Println("update rental item failed:", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	log.Println("rental item updated:", id)

	c.JSON(http.StatusOK, gin.H{
		"data":    item,
		"message": "Rental item updated",
	})
}

func (h *RentalItemHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		log.Println("invalid rental item id:", c.Param("id"))

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid id",
		})
		return
	}

	var input struct {
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		log.Println("invalid delete request body")

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	if err := h.Service.DeleteRentalItem(uint(id), input.Password); err != nil {
		log.Println("failed to delete rental item:", err)

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	log.Println("rental item deleted successfully:", id)

	c.JSON(http.StatusOK, gin.H{
		"message": "rental item deleted",
	})
}