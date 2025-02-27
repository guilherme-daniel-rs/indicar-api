package entities

import (
	"time"

	"gorm.io/gorm"
)

type Review struct {
	ID         uint      `gorm:"primaryKey;column:id"`
	UserID     uint      `gorm:"column:user_id"`
	User       User      `gorm:"foreignKey:UserID;references:ID"`
	Status     int       `gorm:"column:status"`
	ReviewerID uint      `gorm:"column:reviewer_id"`
	Reviewer   Reviewer  `gorm:"foreignKey:ReviewerID;references:ID"`
	Car        string    `gorm:"column:car"`
	Grade      float64   `gorm:"column:grade"`
	Location   string    `gorm:"column:location"`
	CreatedAt  time.Time `gorm:"column:created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at"`
}

func (Review) TableName() string {
	return "review"
}

func (review *Review) BeforeCreate(tx *gorm.DB) error {
	review.CreatedAt = time.Now()
	return nil
}

func (review *Review) BeforeUpdate(tx *gorm.DB) (err error) {
	review.UpdatedAt = time.Now()
	return nil
}
