package entities

import (
	"time"

	"gorm.io/gorm"
)

type ReviewComment struct {
	ID          uint      `gorm:"primaryKey;column:id"`
	ReviewID    uint      `gorm:"column:review_id"`
	Review      Review    `gorm:"foreignKey:ReviewID;references:ID"`
	Description string    `gorm:"column:description"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
}

func (ReviewComment) TableName() string {
	return "review_comment"
}

func (reviewComment *ReviewComment) BeforeCreate(tx *gorm.DB) error {
	reviewComment.CreatedAt = time.Now()
	return nil
}

func (reviewComment *ReviewComment) BeforeUpdate(tx *gorm.DB) (err error) {
	reviewComment.UpdatedAt = time.Now()
	return nil
}
