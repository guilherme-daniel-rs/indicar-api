package entities

type Evaluator struct {
	UserID       int     `json:"user_id" gorm:"primaryKey"`
	DocumentID   *string `json:"document_id,omitempty" gorm:"type:varchar(32)"`
	Rating       float64 `json:"rating" gorm:"type:decimal(3,2);default:0.00;index:idx_rating"`
	TotalReviews int     `json:"total_reviews" gorm:"default:0;index:idx_rating"`
	Bio          *string `json:"bio,omitempty" gorm:"type:varchar(255)"`

	// Relationships
	User User `json:"-" gorm:"foreignKey:UserID"`
}
