package handler

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"equipment-service/internal/model"
	"equipment-service/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type EquipmentUnitHandler struct {
	Service *service.EquipmentUnitService
}

func NewEquipmentUnitHandler(handler *service.EquipmentUnitService) *EquipmentUnitHandler {
	return &EquipmentUnitHandler{
		Service: handler,
	}
}

type CreateEquipmentUnitRequest struct {
	EquipmentTypeId uint   `json:"equipment_type_id" binding:"required"`
	SerialNumber    string `json:"serial_number" binding:"required"`
}

type UpdateEquipmentUnitRequest struct {
	SerialNumber string                   `json:"serial_number"`
	Condition    model.EquipmentCondition `json:"condition"`
	Status       model.EquipmentStatus    `json:"status"`
}

type UpdateEquipmentConditionRequest struct {
	Condition model.EquipmentCondition `json:"condition" binding:"required"`
}



func (h *EquipmentUnitHandler) GetById(c *gin.Context) {

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid equipment unit id",
		})
		return
	}

	equipmentUnit, err := h.Service.GetById(uint(id))

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "equipment unit not found",
			})
			return
		}

		log.Println("failed to get equipment unit:", err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to get equipment unit",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": equipmentUnit,
	})
}

func (h *EquipmentUnitHandler) Update(c *gin.Context) {

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid equipment unit id",
		})
		return
	}

	var req UpdateEquipmentUnitRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid message body",
		})
		return
	}

	unit, err := h.Service.Update(
		uint(id),
		req.SerialNumber,
		req.Condition,
		req.Status,
	)

	if err != nil {
		log.Println("failed to update equipment unit:", err)

		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	log.Println("equipment unit updated:", id)

	c.JSON(http.StatusOK, gin.H{
		"message": "equipment unit updated",
		"data":    unit,
	})
}

func (r *EquipmentUnitHandler) Delete(c *gin.Context) {

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid equipment unit id",
		})
		return
	}

	if err := r.Service.Delete(uint(id)); err != nil {
		log.Println("failed to delete equipment unit:", err)

		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	log.Println("equipment unit deleted:", id)

	c.JSON(http.StatusOK, gin.H{
		"message": "equipment unit deleted",
	})
}

type UpdateEquipmentStatusRequest struct {
	Status model.EquipmentStatus `json:"status" binding:"required"`
}

func (h *EquipmentUnitHandler) UpdateStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid equipment unit id",
		})
		return
	}

	var req UpdateEquipmentStatusRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request body",
		})
		return
	}

	unit, err := h.Service.UpdateStatus(
		uint(id),
		req.Status,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "equipment status updated",
		"data":    unit,
	})
}


func (h *EquipmentUnitHandler) GetAll(c *gin.Context) {
	limitStr := c.Query("limit")
	log.Println("EquipmentUnit GetAll requested with limit:", limitStr)

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

	log.Println("Fetching equipment units with limit:", limit)
	equipmentUnits, err := h.Service.GetAll(limit)
	if err != nil {
		log.Println("Database error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch data"})
		return
	}

	log.Println("Success: fetched units count:", len(equipmentUnits))
	c.JSON(http.StatusOK, gin.H{
		"data": equipmentUnits,
	})
}




func (h *EquipmentUnitHandler) UpdateCondition(c *gin.Context) {

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid equipment unit id",
		})
		return
	}

	var req UpdateEquipmentConditionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("invalid condition request:", err)

		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request body",
		})
		return
	}

	unit, err := h.Service.UpdateCondition(
		uint(id),
		req.Condition,
	)

	if err != nil {
		log.Println("failed to update equipment condition:", err)

		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	log.Println("equipment condition updated:", id)

	c.JSON(http.StatusOK, gin.H{
		"message": "equipment condition updated successfully",
		"data":    unit,
	})
}