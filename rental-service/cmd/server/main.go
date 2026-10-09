package main

import (
	"log"
	"os"
	"rental-service/internal/client"
	"rental-service/internal/config"
	"rental-service/internal/handler"
	"rental-service/internal/model"
	"rental-service/internal/repository"
	"rental-service/internal/service"
     "github.com/robfig/cron/v3"
	"github.com/gin-gonic/gin"
)

func main() {

	logFile, err := os.OpenFile(
		"../application/application.log",
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0666,
	)
	if err != nil {
		log.Fatal("failed to create log file:", err)
	}

	defer logFile.Close()

	log.SetOutput(logFile)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	db, err := config.ConnectDB()
	if err != nil {
		log.Fatal("database connection failed :", err)
	}

	err = db.AutoMigrate(&model.Rental{},&model.RentalItem{},	&model.RentalDue{},)
	if err != nil {
		log.Fatal("migration failed :", err)
	}

	log.Println("database connected")
	log.Println("migration ")

	router := gin.Default()

	CustomerClient := client.NewCustomerClient("http://localhost:8081")
	EquipmentClient := client.NewEquipmentClient("http://localhost:8082")
	Repo := repository.NewRentalRepository(db)
	
	RentalDueRepo := repository.NewRentalDueRepository(db) 
	Service := service.NewRentalService( Repo, RentalDueRepo, CustomerClient, EquipmentClient, )
	
	
c := cron.New()
/*
_, err = c.AddFunc("5 0 * * *", func() {
	log.Println("rental scheduler started")

	if err := Service.ProcessDailyRentalDue(); err != nil {
		log.Println("daily rental due failed:", err)
		return
	}

	if err := Service.MarkOverdueRentals(); err != nil {
		log.Println("overdue rental check failed:", err)
		return
	}

	log.Println("rental scheduler completed")
})

if err != nil {
	log.Fatal("failed to schedule rental job:", err)
}

c.Start()

log.Println("rental scheduler started")
defer c.Stop()*/

_, err = c.AddFunc("@every 1m", func() {
	log.Println("rental scheduler started")

	if err := Service.ProcessDailyRentalDue(); err != nil {
		log.Println("daily rental due failed:", err)
		return
	}

	if err := Service.MarkOverdueRentals(); err != nil {
		log.Println("overdue rental check failed:", err)
		return
	}

	log.Println("rental scheduler completed")
})

if err != nil {
	log.Fatal("failed to schedule rental job:", err)
}

c.Start()

log.Println("rental scheduler started")
defer c.Stop()



	Handler := handler.NewRentalHandler(Service)

	RentalRepo := repository.NewRentalItemRepository(db)
	RentalService := service.NewRentalItemService(RentalRepo)
	RentalHandler := handler.NewRentalItemHandler(RentalService)

	router.POST("/rental", Handler.Create)
	router.GET("/rental", Handler.GetAll)
	router.GET("/rental/:id", Handler.GetById)
	router.DELETE("/rental/:id", Handler.Delete)

	router.POST("/rental-item", RentalHandler.Create)
	router.GET("/rental-item/:rental_id", RentalHandler.GetByRentalID)
	router.PUT("/rental-item/:id", RentalHandler.Update)
	router.PATCH("/rental/:id/return", Handler.Return)
	router.DELETE("/rental-item/:id", RentalHandler.Delete)

	

	router.Run(":8080")

}
