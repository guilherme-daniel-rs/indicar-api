package routes

import (
	"indicar-api/configs"
	"indicar-api/internal/application/controllers"
	"indicar-api/internal/application/services"
	"indicar-api/internal/infrastructure/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupNotificationRoutes(router *gin.Engine, db *gorm.DB) error {
	notificationService := services.NewNotificationService(db)
	notificationController := controllers.NewNotificationController(notificationService)

	authMiddleware := middleware.AuthMiddleware([]byte(configs.Get().JWT.Secret))

	devices := router.Group("/devices")
	devices.Use(authMiddleware)
	{
		devices.POST("", notificationController.RegisterDevice)
	}

	return nil
}
