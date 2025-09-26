package routes

import (
	"indicar-api/configs"
	"indicar-api/internal/application/controllers"
	"indicar-api/internal/application/services"
	"indicar-api/internal/infrastructure/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupEvaluationRoutes(router *gin.Engine, db *gorm.DB) {
	evaluationService := services.NewEvaluationService(db)
	evaluationController := controllers.NewEvaluationController(evaluationService)

	// Middleware de autenticação
	authMiddleware := middleware.AuthMiddleware([]byte(configs.Get().JWT.Secret))

	// Evaluation endpoints (protected)
	evaluations := router.Group("/evaluations")
	evaluations.Use(authMiddleware)
	{
		evaluations.POST("", evaluationController.Create)
		evaluations.GET("/:id", evaluationController.GetByID)
		evaluations.GET("", evaluationController.List)
		evaluations.PATCH("/:id", evaluationController.Update)
	}
}
