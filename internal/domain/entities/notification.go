package entities

import "time"

type Notification struct {
	ID        int        `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID    int        `json:"user_id" gorm:"not null;index:idx_user_status"`
	Channel   string     `json:"channel" gorm:"type:varchar(16);not null"`
	Title     string     `json:"title" gorm:"type:varchar(120);not null"`
	Message   string     `json:"message" gorm:"type:varchar(255);not null"`
	Status    string     `json:"status" gorm:"type:varchar(16);not null;default:queued;index:idx_user_status"`
	CreatedAt time.Time  `json:"created_at" gorm:"not null;default:CURRENT_TIMESTAMP;index"`
	SentAt    *time.Time `json:"sent_at,omitempty"`

	// Relationships
	User User `json:"-" gorm:"foreignKey:UserID"`
}

type PushDevice struct {
	ID          int        `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID      int        `json:"user_id" gorm:"not null;index:idx_user_platform"`
	Platform    string     `json:"platform" gorm:"type:varchar(8);not null;index:idx_user_platform"`
	DeviceToken string     `json:"device_token" gorm:"type:varchar(255);not null;unique"`
	CreatedAt   time.Time  `json:"created_at" gorm:"not null;default:CURRENT_TIMESTAMP"`
	LastSeenAt  *time.Time `json:"last_seen_at,omitempty"`

	// Relationships
	User User `json:"-" gorm:"foreignKey:UserID"`
}
