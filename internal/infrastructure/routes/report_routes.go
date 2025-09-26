package routes

import (
	"indicar-api/configs"
	"indicar-api/internal/application/controllers"
	"indicar-api/internal/application/services"
	"indicar-api/internal/infrastructure/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupReportRoutes(router *gin.Engine, db *gorm.DB) error {
	reportService, err := services.NewReportService(db)
	if err != nil {
		return err
	}
	reportController := controllers.NewReportController(reportService)

	authMiddleware := middleware.AuthMiddleware([]byte(configs.Get().JWT.Secret))

	reports := router.Group("/reports")
	reports.Use(authMiddleware)
	{
		reports.POST("", reportController.CreateOrUpdate)
		reports.GET("/:id", reportController.GetByID)
		reports.GET("/:id/file", reportController.GetFileURL)
		reports.PATCH("/:id", reportController.CreateOrUpdate)
	}

	return nil
}
