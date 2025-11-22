package controllers

import (
	"indicar-api/internal/application/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CityController struct {
	cityService *services.CityService
}

func NewCityController(cityService *services.CityService) *CityController {
	return &CityController{
		cityService: cityService,
	}
}

// @Summary List all cities
// @Description Get a list of all cities available in the system
// @Tags cities
// @Produce json
// @Success 200 {array} services.CityResponse
// @Failure 500 {object} map[string]interface{}
// @Router /cities [get]
func (c *CityController) List(ctx *gin.Context) {
	cities, err := c.cityService.ListAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, cities)
}
