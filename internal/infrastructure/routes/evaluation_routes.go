package routes

import (
	"fmt"
	"indicar-api/configs"
	"indicar-api/internal/application/controllers"
	"indicar-api/internal/application/services"
	"indicar-api/internal/infrastructure/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupEvaluationRoutes(router *gin.Engine, db *gorm.DB) error {
	evaluationService := services.NewEvaluationService(db)
	evaluationPhotoService, err := services.NewEvaluationPhotoService(db)
	if err != nil {
		return fmt.Errorf("failed to initialize photo service: %w", err)
	}

	evaluationController := controllers.NewEvaluationController(evaluationService, evaluationPhotoService)

	// Initialize report controller for nested route
	reportService, err := services.NewReportService(db)
	if err != nil {
		return fmt.Errorf("failed to initialize report service: %w", err)
	}
	reportController := controllers.NewReportController(reportService)

	authMiddleware := middleware.AuthMiddleware([]byte(configs.Get().JWT.Secret))

	evaluations := router.Group("/evaluations")
	evaluations.Use(authMiddleware)
	{
		evaluations.POST("", evaluationController.Create)
		evaluations.GET("", evaluationController.List)

		// Nested routes must come before /:id route to avoid conflicts
		evaluations.GET("/:id/report", reportController.GetByEvaluationID)
		evaluations.POST("/:id/photos", evaluationController.UploadPhoto)
		evaluations.GET("/:id/photos", evaluationController.ListPhotos)
		evaluations.GET("/:id/photos/:photo_id", evaluationController.GetPhotoURL)

		// Generic routes come last
		evaluations.GET("/:id", evaluationController.GetByID)
		evaluations.PATCH("/:id", evaluationController.Update)
	}

	return nil
}
