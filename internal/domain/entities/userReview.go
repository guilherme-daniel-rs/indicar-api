package entities

import (
	"time"

	"gorm.io/gorm"
)

type UserReview struct {
	ID         uint      `gorm:"primaryKey;column:id"`
	UserID     uint      `gorm:"column:user_id"`
	User       User      `gorm:"foreignKey:UserID;references:ID"`
	ReviewerID uint      `gorm:"column:reviewer_id"`
	Reviewer   Reviewer  `gorm:"foreignKey:ReviewerID;references:ID"`
	Grade      float64   `gorm:"column:grade"`
	Content    string    `gorm:"column:content"`
	CreatedAt  time.Time `gorm:"column:created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at"`
}

func (UserReview) TableName() string {
	return "user_review"
}

func (userReview *UserReview) BeforeCreate(tx *gorm.DB) error {
	userReview.CreatedAt = time.Now()
	return nil
}

func (userReview *UserReview) BeforeUpdate(tx *gorm.DB) (err error) {
	userReview.UpdatedAt = time.Now()
	return nil
}
