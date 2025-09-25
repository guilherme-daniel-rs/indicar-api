package entities

import "time"

type Payment struct {
	ID               int       `json:"id" gorm:"primaryKey;autoIncrement"`
	EvaluationID     int       `json:"evaluation_id" gorm:"not null;unique"`
	Provider         string    `json:"provider" gorm:"type:varchar(24);not null"`
	ProviderChargeID string    `json:"provider_charge_id" gorm:"type:varchar(64);not null;uniqueIndex:idx_provider_charge"`
	AmountCents      int       `json:"amount_cents" gorm:"not null"`
	Currency         string    `json:"currency" gorm:"type:char(3);not null;default:BRL"`
	Status           string    `json:"status" gorm:"type:varchar(24);not null;index:idx_status_created"`
	CreatedAt        time.Time `json:"created_at" gorm:"type:datetime(3);not null;default:current_timestamp(3);index:idx_status_created"`
	UpdatedAt        time.Time `json:"updated_at" gorm:"type:datetime(3);not null;default:current_timestamp(3) on update current_timestamp(3)"`

	// Relationships
	Evaluation Evaluation `json:"-" gorm:"foreignKey:EvaluationID"`
}
