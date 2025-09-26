package routes

import (
	"indicar-api/configs"
	"indicar-api/internal/application/controllers"
	"indicar-api/internal/application/services"
	"indicar-api/internal/infrastructure/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupUserRoutes(router *gin.Engine, db *gorm.DB) {
	userService := services.NewUserService(db)
	userController := controllers.NewUserController(userService)

	// Middleware de autenticação
	authMiddleware := middleware.AuthMiddleware([]byte(configs.Get().JWT.Secret))

	// Current user endpoints (protegidos)
	me := router.Group("/me")
	me.Use(authMiddleware)
	{
		me.GET("", userController.GetMe)
		me.PUT("", userController.UpdateMe)
	}

	// Evaluator endpoints (protegidos)
	evaluators := router.Group("/evaluators")
	evaluators.Use(authMiddleware)
	{
		evaluators.GET("/:id", userController.GetEvaluator)
	}
}
