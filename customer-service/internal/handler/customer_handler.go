package handler

import (
	"customer-service/internal/model"
	"customer-service/internal/service"
	
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type CustomerHandler struct {
	Service *service.CustomerService
}

type DepositRequest struct {
	Amount float64 `json:"amount"`
}

type PaymentRequest struct {
	Amount float64 `json:"amount"`
}

type OutstandingRequest struct {
	Amount float64 `json:"amount"`
}


func NewCustomerHandler(handler *service.CustomerService) *CustomerHandler {
	return &CustomerHandler{
		Service: handler,
	}
}

func (r *CustomerHandler) Create(c *gin.Context) {

	var customer model.Customer

	if err := c.ShouldBindJSON(&customer); err != nil {
		log.Println("failed to bind customer request:", err)

		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid request body",
		})
		return
	}

	if customer.Name == "" || customer.Email == "" || customer.Phone == "" {
		log.Println("customer required fields are missing")

		c.JSON(http.StatusCreated, gin.H{
			"success": false,
			"message": "name,email,and phone are required",
		})
		return
	}

	err := r.Service.Create(&customer)

	if err != nil {
		log.Println("failed to create customer:", err)

		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	log.Println("customer created:", customer.ID)

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "customer created successfully",
		"data":    customer,
	})
}

func (r *CustomerHandler) GetAll(c *gin.Context) {

	customer, total, err := r.Service.GetAll(1, 10)
	if err != nil {
		log.Println("failed to get customers:", err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to get customers",
		})
		return
	}

	log.Println("customers fetched successfully, total:", total)

	c.JSON(http.StatusOK, gin.H{
		"data":  customer,
		"total": total,
	})
}

func (h *CustomerHandler) GetById(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		log.Println("invalid customer id:", c.Param("id"))

		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid id ",
		})
		return
	}

	customer, err := h.Service.GetById(uint(id))

	if err != nil {
		log.Println("customer not found:", id)

		c.JSON(http.StatusNotFound, gin.H{
			"message": "customer not found",
		})
		return
	}

	log.Println("customer fetched:", id)

	c.JSON(http.StatusOK, customer)
}


func (h *CustomerHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid id",
		})
		return
	}

	var customer model.Customer
	
		
		
		

	if err := c.ShouldBindJSON(&customer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid input",
		})
		return
	}

	if customer.Name == "" || customer.Email == "" || customer.Phone == "" {
	c.JSON(http.StatusBadRequest, gin.H{
		"message": "name, email and phone are required",
	})
	return
}

	customer.ID = uint(id)

	if err := h.Service.Update(&customer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, customer)
}

func (r *CustomerHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		log.Println("invalid customer id:", c.Param("id"))
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid id",
		})
		return
	}

	err = r.Service.Deactivate(uint(id))
	if err != nil {
		log.Println("failed to deactivate customer:", err)

		if err.Error() == "customer already inactive" {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "empty your wallet or pay outstanding due",
		})
		return
	}

	log.Println("customer status updated to INACTIVE:", id)

	c.JSON(http.StatusOK, gin.H{
		"message": "customer status updated to INACTIVE",
	})
}



func (r *CustomerHandler) VerifyKYC(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid customer id",
		})
		return
	}


	if c.Request.ContentLength == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			
			"message": "enter email",
		})
		return
	}

	var req struct {
		Email string `json:"email"`
		
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			
			"message": "enter email",
		})
		return
	}

	if strings.TrimSpace(req.Email) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			
			"message": "enter email",
		})
		return
	}

	if err := r.Service.VerifyKYC(uint(id), req.Email); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		
		"message": "kyc verified successfully",
	})
}


func (r *CustomerHandler) UpdateDeposit(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		log.Println("invalid customer id:", c.Param("id"))
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid id",
		})
		return
	}

	var input struct {
		Deposit *float64 `json:"deposit"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		log.Println("failed to bind deposit request:", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "please enter deposit amount",
		})
		return
	}

	if input.Deposit == nil {
		log.Println("deposit field missing in request payload")
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "please enter deposit amount",
		})
		return
	}

	err = r.Service.UpdateDeposit(uint(id), *input.Deposit)
	if err != nil {
		log.Println("failed to update deposit:", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	log.Println("deposit updated successfully:", id)
	c.JSON(http.StatusOK, gin.H{
		"message": "deposit updated successfully",
	})
}


func (r *CustomerHandler) UpdateBy(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		log.Println("invalid customer id:", c.Param("id"))

		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid customer id",
		})
		return
	}

	var updates map[string]interface{}

	if err := c.ShouldBindJSON(&updates); err != nil {
		log.Println("failed to bind update fields:", err)

		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid json body",
		})
		return
	}

	if err := r.Service.UpdateBy(uint(id), updates); err != nil {
		log.Println("failed to update customer:", err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	log.Println("customer fields updated:", id)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "customer updated successfully",
	})
}




func (h *CustomerHandler) DeductDeposit(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid customer id",
		})
		return
	}

	var input DepositRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid input",
		})
		return
	}

	if input.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "amount must be greater than zero",
		})
		return
	}

	err = h.Service.DeductDeposit(uint(id), input.Amount)
	if err != nil {
		log.Println("failed to deduct deposit:", err)

		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "deposit updated",
	})
}

func (h *CustomerHandler) PayOutstanding(c *gin.Context) {

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid customer id",
		})
		return
	}

	var req PaymentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("invalid payment request:", err)

		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid payment request",
		})
		return
	}

	customer, err := h.Service.PayOutstanding(uint(id), req.Amount)
	if err != nil {
		log.Println("payment failed:", err)

		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	log.Println("payment successful:", id, req.Amount)

	c.JSON(http.StatusOK, gin.H{
		"message": "payment successful",
		"data":    customer,
	})
}



func (h *CustomerHandler) AddOutstandingDue(c *gin.Context) {

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid customer id",
		})
		return
	}

	var req OutstandingRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("invalid outstanding request:", err)

		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request",
		})
		return
	}

	customer, err := h.Service.AddOutstandingDue(uint(id), req.Amount)
	if err != nil {
		log.Println("failed to update outstanding due:", err)

		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	log.Println("outstanding due updated:", id)

	c.JSON(http.StatusOK, gin.H{
		"message": "outstanding due updated successfully",
		"data": gin.H{
			"ID":             customer.ID,
			"OutstandingDue": customer.OutStandingDue,
		},
	})
}