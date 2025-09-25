package entities

import "time"

type ReportStatus string

const (
	ReportStatusDraft     ReportStatus = "draft"
	ReportStatusFinalized ReportStatus = "finalized"
)

type Report struct {
	ID           int          `json:"id" gorm:"primaryKey;autoIncrement"`
	EvaluationID int          `json:"evaluation_id" gorm:"not null;unique"`
	EvaluatorID  int          `json:"evaluator_id" gorm:"not null;index:idx_evaluator_status"`
	Summary      *string      `json:"summary,omitempty" gorm:"type:varchar(255)"`
	Status       ReportStatus `json:"status" gorm:"type:ENUM('draft','finalized');not null;index:idx_evaluator_status"`
	CreatedAt    time.Time    `json:"created_at" gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt    time.Time    `json:"updated_at" gorm:"not null;default:CURRENT_TIMESTAMP;autoUpdateTime"`

	// Relationships
	Evaluation Evaluation `json:"-" gorm:"foreignKey:EvaluationID"`
	Evaluator  User       `json:"-" gorm:"foreignKey:EvaluatorID"`
}

type ReportFile struct {
	ID          int       `json:"id" gorm:"primaryKey;autoIncrement"`
	ReportID    int       `json:"report_id" gorm:"not null;unique"`
	S3Bucket    string    `json:"s3_bucket" gorm:"type:varchar(128);not null;uniqueIndex:idx_s3_location"`
	S3Key       string    `json:"s3_key" gorm:"type:varchar(256);not null;uniqueIndex:idx_s3_location"`
	ContentType string    `json:"content_type" gorm:"type:varchar(80);default:'application/pdf'"`
	SizeBytes   *int      `json:"size_bytes,omitempty"`
	CreatedAt   time.Time `json:"created_at" gorm:"not null;default:CURRENT_TIMESTAMP"`

	// Relationships
	Report Report `json:"-" gorm:"foreignKey:ReportID"`
}
