package routes

import (
	"indicar-api/internal/application/controllers"
	"indicar-api/internal/application/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupAuthRoutes(router *gin.Engine, db *gorm.DB) {
	authService := services.NewAuthService(db)
	authController := controllers.NewAuthController(authService)

	auth := router.Group("/auth")
	{
		auth.POST("/signup", authController.Signup)
		auth.POST("/login", authController.Login)
		auth.POST("/refresh", authController.RefreshToken)
	}
}
