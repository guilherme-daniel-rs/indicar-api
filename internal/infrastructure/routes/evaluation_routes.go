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

	authMiddleware := middleware.AuthMiddleware([]byte(configs.Get().JWT.Secret))

	evaluations := router.Group("/evaluations")
	evaluations.Use(authMiddleware)
	{
		evaluations.POST("", evaluationController.Create)
		evaluations.GET("/:id", evaluationController.GetByID)
		evaluations.GET("", evaluationController.List)
		evaluations.PATCH("/:id", evaluationController.Update)

		evaluations.POST("/:id/photos", evaluationController.UploadPhoto)
		evaluations.GET("/:id/photos", evaluationController.ListPhotos)
	}

	return nil
}
