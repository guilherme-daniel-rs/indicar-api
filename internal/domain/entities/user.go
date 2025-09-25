package entities

import "time"

type UserRole string

const (
	UserRoleUser      UserRole = "user"
	UserRoleEvaluator UserRole = "evaluator"
	UserRoleAdmin     UserRole = "admin"
)

type User struct {
	ID           int       `json:"id" gorm:"primaryKey;autoIncrement"`
	FullName     string    `json:"full_name" gorm:"type:varchar(120);not null"`
	Email        string    `json:"email" gorm:"type:varchar(160);unique;not null"`
	PasswordHash string    `json:"-" gorm:"type:varchar(255);not null;column:password_hash"`
	Phone        *string   `json:"phone,omitempty" gorm:"type:varchar(32)"`
	Role         UserRole  `json:"role" gorm:"type:varchar(20);not null;index:idx_role_active"`
	CreatedAt    time.Time `json:"created_at" gorm:"type:datetime(3);not null;default:current_timestamp(3)"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"type:datetime(3);not null;default:current_timestamp(3) on update current_timestamp(3)"`
	IsActive     bool      `json:"is_active" gorm:"not null;default:true;index:idx_role_active"`
}
