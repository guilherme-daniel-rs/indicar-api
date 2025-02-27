package entities

import (
	"time"

	"gorm.io/gorm"
)

type Message struct {
	ID        uint      `gorm:"primaryKey;column:id"`
	ChatID    uint      `gorm:"column:chat_id"`
	Chat      Chat      `gorm:"foreignKey:ChatID;references:ID"`
	UserID    uint      `gorm:"column:user_id"`
	User      User      `gorm:"foreignKey:UserID;references:ID"`
	Content   string    `gorm:"column:content"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
	Status    int       `gorm:"column:status"`
}

func (Message) TableName() string {
	return "message"
}

func (message *Message) BeforeCreate(tx *gorm.DB) error {
	message.CreatedAt = time.Now()
	return nil
}

func (message *Message) BeforeUpdate(tx *gorm.DB) (err error) {
	message.UpdatedAt = time.Now()
	return nil
}
