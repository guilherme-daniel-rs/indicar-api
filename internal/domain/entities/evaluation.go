package entities

import "time"

type EvaluationStatus string

const (
	EvaluationStatusCreated    EvaluationStatus = "created"
	EvaluationStatusAccepted   EvaluationStatus = "accepted"
	EvaluationStatusInProgress EvaluationStatus = "in_progress"
	EvaluationStatusCompleted  EvaluationStatus = "completed"
	EvaluationStatusCanceled   EvaluationStatus = "canceled"
)

type Evaluation struct {
	ID           int              `json:"id" gorm:"primaryKey;autoIncrement"`
	RequesterID  int              `json:"requester_id" gorm:"not null;index:idx_requester_status"`
	EvaluatorID  *int             `json:"evaluator_id,omitempty" gorm:"index:idx_evaluator_status"`
	CityID       int              `json:"city_id" gorm:"not null;index:idx_city_status"`
	VehicleMake  string           `json:"vehicle_make" gorm:"type:varchar(80);not null;index:idx_vehicle"`
	VehicleModel string           `json:"vehicle_model" gorm:"type:varchar(120);not null;index:idx_vehicle"`
	VehicleYear  *int             `json:"vehicle_year,omitempty"`
	VehiclePlate *string          `json:"vehicle_plate,omitempty" gorm:"type:varchar(16)"`
	Status       EvaluationStatus `json:"status" gorm:"type:ENUM('created', 'accepted', 'in_progress', 'completed', 'canceled');not null;index:idx_requester_status,idx_evaluator_status,idx_city_status"`
	Notes        *string          `json:"notes,omitempty" gorm:"type:text"`
	CreatedAt    time.Time        `json:"created_at" gorm:"type:datetime(3);not null;default:current_timestamp(3)"`
	UpdatedAt    time.Time        `json:"updated_at" gorm:"type:datetime(3);not null;default:current_timestamp(3) on update current_timestamp(3)"`

	// Relationships
	Requester User  `json:"-" gorm:"foreignKey:RequesterID"`
	Evaluator *User `json:"-" gorm:"foreignKey:EvaluatorID"`
	City      City  `json:"-" gorm:"foreignKey:CityID"`
}

type EvaluationPhoto struct {
	ID           int       `json:"id" gorm:"primaryKey;autoIncrement"`
	EvaluationID int       `json:"evaluation_id" gorm:"not null;index:idx_evaluation_created"`
	S3Bucket     string    `json:"s3_bucket" gorm:"type:varchar(128);not null;uniqueIndex:idx_s3_location"`
	S3Key        string    `json:"s3_key" gorm:"type:varchar(256);not null;uniqueIndex:idx_s3_location"`
	ContentType  *string   `json:"content_type,omitempty" gorm:"type:varchar(80)"`
	SizeBytes    *int      `json:"size_bytes,omitempty"`
	CreatedAt    time.Time `json:"created_at" gorm:"type:datetime(3);not null;default:current_timestamp(3);index:idx_evaluation_created"`

	// Relationships
	Evaluation Evaluation `json:"-" gorm:"foreignKey:EvaluationID"`
}
