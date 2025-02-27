package entities

import (
	"time"

	"gorm.io/gorm"
)

type ReviewImage struct {
	ID        uint      `gorm:"primaryKey;column:id"`
	ReviewID  uint      `gorm:"column:review_id"`
	Review    Review    `gorm:"foreignKey:ReviewID;references:ID"`
	URL       string    `gorm:"column:url"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (ReviewImage) TableName() string {
	return "review_image"
}

func (reviewImage *ReviewImage) BeforeCreate(tx *gorm.DB) error {
	reviewImage.CreatedAt = time.Now()
	return nil
}

func (reviewImage *ReviewImage) BeforeUpdate(tx *gorm.DB) (err error) {
	reviewImage.UpdatedAt = time.Now()
	return nil
}
