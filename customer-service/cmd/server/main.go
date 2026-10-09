package main

import (
	"customer-service/internal/config"
	"customer-service/internal/handler"
	"customer-service/internal/model"
	"customer-service/internal/repository"
	"customer-service/internal/service"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {

	logFile, err := os.OpenFile("../application/application.log",
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0666,
	)
	if err != nil {
		log.Fatal("failed to open log file:", err)
	}

	defer logFile.Close()

	log.SetOutput(logFile)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	db, err := config.ConnectDB()
	if err != nil {
		log.Fatal("database connection failed:", err)
	}

	if err := db.AutoMigrate(&model.Customer{}); err != nil {
		log.Fatal("migration failed:", err)
	}

	log.Println("database connected")
	log.Println("migration completed")

	router := gin.Default()

	repo := repository.NewCustomerRepository(db)
	service := service.NewCustomerService(repo)
	handler := handler.NewCustomerHandler(service)

	

	router.POST("/customer", handler.Create)
	router.GET("/customer", handler.GetAll)
	router.GET("/customer/:id", handler.GetById)

	router.PUT("/customer/:id", handler.Update)
	router.DELETE("/customer/:id", handler.Delete)
	router.PATCH("/customer/:id/VerifyKYC", handler.VerifyKYC)
	router.PATCH("/customer/:id/UpdateDeposit", handler.UpdateDeposit)
	router.PATCH("/customer/:id/UpdateBy", handler.UpdateBy)
	router.PATCH("/customer/:id/deposit", handler.DeductDeposit)
	router.PATCH("/customer/:id/payment", handler.PayOutstanding)
	router.PATCH("/customer/:id/outstanding", handler.AddOutstandingDue)

	log.Println("customer service running on port 8081")

	if err := router.Run(":8081"); err != nil {
		log.Fatal("server failed:", err)
	}
}
