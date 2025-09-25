package entities

type City struct {
	ID          int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string `json:"name" gorm:"type:varchar(120);not null;uniqueIndex:idx_name_state"`
	StateCode   string `json:"state_code" gorm:"type:varchar(2);not null;uniqueIndex:idx_name_state"`
	CountryCode string `json:"country_code" gorm:"type:varchar(2);not null;default:BR"`
}

type EvaluatorCity struct {
	EvaluatorID int  `json:"evaluator_id" gorm:"primaryKey"`
	CityID      int  `json:"city_id" gorm:"primaryKey"`
	CoverageKm  *int `json:"coverage_km,omitempty" gorm:"default:30"`

	// Relationships
	Evaluator User `json:"-" gorm:"foreignKey:EvaluatorID"`
	City      City `json:"-" gorm:"foreignKey:CityID"`
}
