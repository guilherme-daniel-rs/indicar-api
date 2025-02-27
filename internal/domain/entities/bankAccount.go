package entities

import (
	"time"

	"gorm.io/gorm"
)

type BankAccount struct {
	ID         uint      `gorm:"primaryKey;column:id"`
	ReviewerID uint      `gorm:"column:reviewer_id"`
	Reviewer   Reviewer  `gorm:"foreignKey:ReviewerID;references:ID"`
	CreatedAt  time.Time `gorm:"column:created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at"`
}

func (BankAccount) TableName() string {
	return "bank_account"
}

func (bankAccount *BankAccount) BeforeCreate(tx *gorm.DB) error {
	bankAccount.CreatedAt = time.Now()
	return nil
}

func (bankAccount *BankAccount) BeforeUpdate(tx *gorm.DB) (err error) {
	bankAccount.UpdatedAt = time.Now()
	return nil
}
