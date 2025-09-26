package controllers

import (
	"indicar-api/internal/application/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type NotificationController struct {
	notificationService *services.NotificationService
}

func NewNotificationController(notificationService *services.NotificationService) *NotificationController {
	return &NotificationController{
		notificationService: notificationService,
	}
}

// @Summary Register device for push notifications
// @Description Register a device for receiving push notifications
// @Tags devices
// @Accept json
// @Produce json
// @Security Bearer
// @Param input body services.RegisterDeviceInput true "Device registration data"
// @Success 201 {object} entities.PushDevice
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /devices [post]
func (c *NotificationController) RegisterDevice(ctx *gin.Context) {
	userID := ctx.GetInt("user_id")
	if userID == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var input services.RegisterDeviceInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	device, err := c.notificationService.RegisterDevice(userID, input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, device)
}
