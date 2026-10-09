package handler

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"equipment-service/internal/model"
	"equipment-service/internal/service"
"strings"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type EquipmentTypeHandler struct {
	Service *service.EquipmentTypeService
}

func NewEquipmentTypeHandler(handler *service.EquipmentTypeService) *EquipmentTypeHandler {
	return &EquipmentTypeHandler{
		Service: handler,
	}
}

type CreateEquipmentTypeRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	DailyRate   int64  `json:"dailyrate"`
	Quantity    int    `json:"quantity"`
}
func (r *EquipmentTypeHandler) Create(c *gin.Context) {

	var req CreateEquipmentTypeRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request body",
		})
		return
	}
	equipmentType := &model.EquipmentType{
	Name:        req.Name,
	Description: req.Description,
	DailyRate:   req.DailyRate,
	Quantity:    req.Quantity,
}
	if err := r.Service.Create(equipmentType); err != nil {
		log.Println("failed to create equipment type:", err)
       
if strings.Contains(err.Error(), "23505") || strings.Contains(err.Error(), "duplicate key") {
			c.JSON(http.StatusConflict, gin.H{
				"message": "Equipment type name already exists",
				"data" : equipmentType.Name,
			})
			return
		}
		
	c.JSON(http.StatusBadRequest, gin.H{
	"message": err.Error(),
		})
	return
}

    log.Println("equipment type created:", equipmentType.ID)

	c.JSON(http.StatusCreated, gin.H{
	"data":    equipmentType,
 })
}

func (h *EquipmentTypeHandler) GetAll(c *gin.Context) {
	limitStr := c.Query("limit")
	log.Println("GetAll requested with limit:", limitStr)

	if limitStr == "" {
		log.Println("Error: limit param missing")
		c.JSON(http.StatusBadRequest, gin.H{"error": "limit query parameter is required"})
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		log.Println("Error: invalid limit value:", limitStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "limit must be a valid positive integer"})
		return
	}

	log.Println("Fetching records with limit:", limit)
	equipmentTypes, err := h.Service.GetAll(limit)
	if err != nil {
		log.Println("Database error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch data"})
		return
	}

	log.Println("Success: fetched rows count:", len(equipmentTypes))
	c.JSON(http.StatusOK, gin.H{
		"data": equipmentTypes,
	})
}



func (h *EquipmentTypeHandler) GetById(c *gin.Context) {

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
if err != nil {
	c.JSON(http.StatusBadRequest, gin.H{
	"message": "invalid equipment type id",
	})
	return
}
equipmentType, err := h.Service.GetById(uint(id))

if err != nil {
	if errors.Is(err, gorm.ErrRecordNotFound) {
c.JSON(http.StatusNotFound, gin.H{
	"message": "equipment type not found",
	})
	return
}

log.Println("failed to get equipment type:", err)

	c.JSON(http.StatusInternalServerError, gin.H{
	"message": "failed to get equipment type",
	})
	return
}

	c.JSON(http.StatusOK, gin.H{
		"data": equipmentType,
	})
}

func (h *EquipmentTypeHandler) Update(c *gin.Context) {

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid equipment type id",
		})
		return
	}

	var req CreateEquipmentTypeRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid message body",
		})
		return
	}

	equipmentType, err := h.Service.Update(
		uint(id),
		req.Name,
		req.Description,
		req.DailyRate,
		req.Quantity,
	)

	if err != nil {
		log.Println("failed to update equipment type:", err)

		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "equipment type not found",
			})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	log.Println("equipment type updated:", id)

	c.JSON(http.StatusOK, gin.H{
		"message": "equipment type updated successfully",
		"data":    equipmentType,
	})
}