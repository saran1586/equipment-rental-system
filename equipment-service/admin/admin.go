package admin

import (
	
	"equipment-service/internal/handler"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, typeHandler *handler.EquipmentTypeHandler) {
	adminProtected := router.Group("/equipment-type")

	adminProtected.Use(func(c *gin.Context) {
			adminPassword := os.Getenv("ADMIN_PASSWORD")
	inputPassword := c.GetHeader("password")
		


		if strings.TrimSpace(inputPassword) == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "this is admin portal : enter password",
			})
			return
		}
		if strings.TrimSpace(adminPassword) == "" {
    c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
        "success": false,
        "message": "server configuration error: admin password not set in environment",
    })
    return
}

		
		if inputPassword != adminPassword {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "invalid password",
			})
			return
		}

		c.Next()
	})

	adminProtected.POST("", typeHandler.Create)
	adminProtected.GET("", typeHandler.GetAll)
	adminProtected.GET("/:id", typeHandler.GetById)
	adminProtected.PUT("/:id", typeHandler.Update)

}
