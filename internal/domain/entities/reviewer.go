package entities

import (
	"time"

	"gorm.io/gorm"
)

type Reviewer struct {
	ID          uint      `gorm:"primaryKey;column:id"`
	UserID      uint      `gorm:"unique;column:user_id"`
	User        User      `gorm:"foreignKey:UserID;references:ID"`
	Grade       float64   `gorm:"column:grade"`
	Price       float64   `gorm:"column:price"`
	Description string    `gorm:"column:description"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
}

func (Reviewer) TableName() string {
	return "reviewer"
}

func (reviewer *Reviewer) BeforeCreate(tx *gorm.DB) error {
	reviewer.CreatedAt = time.Now()
	return nil
}

func (reviewer *Reviewer) BeforeUpdate(tx *gorm.DB) (err error) {
	reviewer.UpdatedAt = time.Now()
	return nil
}
