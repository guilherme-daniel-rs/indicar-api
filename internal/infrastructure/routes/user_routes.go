package routes

import (
	"indicar-api/configs"
	"indicar-api/internal/application/controllers"
	"indicar-api/internal/application/services"
	"indicar-api/internal/infrastructure/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupUserRoutes(router *gin.Engine, db *gorm.DB) error {
	userService := services.NewUserService(db)
	userController := controllers.NewUserController(userService)

	authMiddleware := middleware.AuthMiddleware([]byte(configs.Get().JWT.Secret))

	me := router.Group("/me")
	me.Use(authMiddleware)
	{
		me.GET("", userController.GetMe)
		me.PUT("", userController.UpdateMe)
	}

	evaluators := router.Group("/evaluators")
	evaluators.Use(authMiddleware)
	{
		evaluators.GET("/:id", userController.GetEvaluator)
	}

	return nil
}
