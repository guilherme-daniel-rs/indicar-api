package entities

import (
	"time"

	"gorm.io/gorm"
)

type Chat struct {
	ID         uint      `gorm:"primaryKey;column:id"`
	UserID     uint      `gorm:"column:user_id"`
	User       User      `gorm:"foreignKey:UserID;references:ID"`
	ReviewerID uint      `gorm:"column:reviewer_id"`
	Reviewer   Reviewer  `gorm:"foreignKey:ReviewerID;references:ID"`
	CreatedAt  time.Time `gorm:"column:created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at"`
}

func (Chat) TableName() string {
	return "chat"
}

func (chat *Chat) BeforeCreate(tx *gorm.DB) error {
	chat.CreatedAt = time.Now()
	return nil
}

func (chat *Chat) BeforeUpdate(tx *gorm.DB) (err error) {
	chat.UpdatedAt = time.Now()
	return nil
}
