package entities

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID         uint      `gorm:"primaryKey;column:id"`
	Name       string    `gorm:"column:name"`
	Email      string    `gorm:"unique;column:email"`
	Password   string    `gorm:"column:password"`
	City       string    `gorm:"column:city"`
	PictureURL string    `gorm:"column:picture_url"`
	CreatedAt  time.Time `gorm:"column:created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at"`
}

func (User) TableName() string {
	return "user"
}

func (user *User) BeforeCreate(tx *gorm.DB) error {
	user.CreatedAt = time.Now()
	return nil
}

func (user *User) BeforeUpdate(tx *gorm.DB) (err error) {
	user.UpdatedAt = time.Now()
	return nil
}
