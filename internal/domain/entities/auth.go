package entities

import "time"

type AuthRefreshToken struct {
	ID        int       `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID    int       `json:"user_id" gorm:"not null;index:idx_user_revoked"`
	Token     string    `json:"token" gorm:"type:varchar(255);not null;unique"`
	ExpiresAt time.Time `json:"expires_at" gorm:"not null;index"`
	CreatedAt time.Time `json:"created_at" gorm:"type:datetime(3);not null;default:current_timestamp(3)"`
	Revoked   bool      `json:"revoked" gorm:"not null;default:false;index:idx_user_revoked"`

	// Relationships
	User User `json:"-" gorm:"foreignKey:UserID"`
}
