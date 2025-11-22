package routes

import (
	"indicar-api/internal/application/controllers"
	"indicar-api/internal/application/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupCityRoutes(router *gin.Engine, db *gorm.DB) error {
	cityService := services.NewCityService(db)
	cityController := controllers.NewCityController(cityService)

	cities := router.Group("/cities")
	{
		cities.GET("", cityController.List)
	}

	return nil
}
