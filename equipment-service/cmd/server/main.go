package main

import (
	"equipment-service/admin"
	"equipment-service/internal/config"
	"equipment-service/internal/handler"
	"equipment-service/internal/model"
	"equipment-service/internal/repository"
	"equipment-service/internal/service"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {

	logFile, err := os.OpenFile(
		"../application/application.log",
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
		log.Fatal("database connection failed :", err)
	}

	err = db.AutoMigrate(&model.EquipmentType{}, &model.EquipmentUnit{})
	if err != nil {
		log.Fatal("migration failed :", err)
	}

	log.Println("database connected")
	log.Println("migration ")

	router := gin.Default()
	typeRepo := repository.NewEquipmentTypeRepository(db)
	unitRepo := repository.NewEquipmentUnitRepository(db)
	unitService := service.NewEquipmentUnitService(unitRepo, typeRepo)
	unitHandler := handler.NewEquipmentUnitHandler(unitService)
	typeService := service.NewEquipmentTypeService(typeRepo,unitService)
	typeHandler := handler.NewEquipmentTypeHandler(typeService)



	admin.SetupRoutes(router, typeHandler)

	router.GET("/equipment-unit", unitHandler.GetAll)
	router.GET("/equipment-unit/:id", unitHandler.GetById)
	router.PUT("/equipment-unit/:id", unitHandler.Update)
	router.DELETE("/equipment-unit/:id", unitHandler.Delete)

	router.PATCH("/equipment-unit/:id/status", unitHandler.UpdateStatus)
     router.GET("/internal/equipment-type/:id", typeHandler.GetById)
    router.PATCH("/equipment-unit/:id/condition", unitHandler.UpdateCondition)

	router.Run(":8082")

}
