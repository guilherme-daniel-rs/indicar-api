package entities

import (
	"time"

	"gorm.io/gorm"
)

type Payment struct {
	ID            uint        `gorm:"primaryKey;column:id"`
	UserID        uint        `gorm:"column:user_id"`
	User          User        `gorm:"foreignKey:UserID;references:ID"`
	Value         float64     `gorm:"column:value"`
	CreatedAt     time.Time   `gorm:"column:created_at"`
	UpdatedAt     time.Time   `gorm:"column:updated_at"`
	Status        int         `gorm:"column:status"`
	BankAccountID uint        `gorm:"column:bank_account_id;foreignKey:BankAccountID;references:ID"`
	BankAccount   BankAccount `gorm:"foreignKey:BankAccountID;references:ID"`
}

func (Payment) TableName() string {
	return "payment"
}

func (payment *Payment) BeforeCreate(tx *gorm.DB) error {
	payment.CreatedAt = time.Now()
	return nil
}

func (payment *Payment) BeforeUpdate(tx *gorm.DB) (err error) {
	payment.UpdatedAt = time.Now()
	return nil
}
