package services

import (
	"indicar-api/internal/domain/entities"
	"time"

	"gorm.io/gorm"
)

type NotificationService struct {
	db *gorm.DB
}

func NewNotificationService(db *gorm.DB) *NotificationService {
	return &NotificationService{
		db: db,
	}
}

type RegisterDeviceInput struct {
	Platform    string `json:"platform" binding:"required,oneof=ios android"`
	DeviceToken string `json:"device_token" binding:"required"`
}

func (s *NotificationService) RegisterDevice(userID int, input RegisterDeviceInput) (*entities.PushDevice, error) {
	// Try to find an existing device
	var device entities.PushDevice
	err := s.db.Where("user_id = ? AND platform = ?", userID, input.Platform).First(&device).Error

	if err == nil {
		// Update existing device
		device.DeviceToken = input.DeviceToken
		device.LastSeenAt = &time.Time{}
		*device.LastSeenAt = time.Now()

		if err := s.db.Save(&device).Error; err != nil {
			return nil, err
		}
		return &device, nil
	}

	// Create new device
	device = entities.PushDevice{
		UserID:      userID,
		Platform:    input.Platform,
		DeviceToken: input.DeviceToken,
		LastSeenAt:  &time.Time{},
	}
	*device.LastSeenAt = time.Now()

	if err := s.db.Create(&device).Error; err != nil {
		return nil, err
	}

	return &device, nil
}

func (s *NotificationService) CreateNotification(userID int, title string, message string) (*entities.Notification, error) {
	notification := &entities.Notification{
		UserID:  userID,
		Channel: "push",
		Title:   title,
		Message: message,
		Status:  "queued",
	}

	if err := s.db.Create(notification).Error; err != nil {
		return nil, err
	}

	return notification, nil
}
