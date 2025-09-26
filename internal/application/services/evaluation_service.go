package services

import (
	"errors"
	"indicar-api/internal/domain/entities"

	"gorm.io/gorm"
)

type EvaluationService struct {
	db *gorm.DB
}

func NewEvaluationService(db *gorm.DB) *EvaluationService {
	return &EvaluationService{
		db: db,
	}
}

type CreateEvaluationInput struct {
	CityID       int     `json:"city_id" binding:"required"`
	VehicleMake  string  `json:"vehicle_make" binding:"required"`
	VehicleModel string  `json:"vehicle_model" binding:"required"`
	VehicleYear  *int    `json:"vehicle_year"`
	VehiclePlate *string `json:"vehicle_plate"`
	Notes        *string `json:"notes"`
}

type UpdateEvaluationInput struct {
	EvaluatorID *int    `json:"evaluator_id"`
	Status      *string `json:"status"`
	Notes       *string `json:"notes"`
}

func (s *EvaluationService) Create(userID int, input CreateEvaluationInput) (*entities.Evaluation, error) {
	evaluation := &entities.Evaluation{
		RequesterID:  userID,
		CityID:       input.CityID,
		VehicleMake:  input.VehicleMake,
		VehicleModel: input.VehicleModel,
		VehicleYear:  input.VehicleYear,
		VehiclePlate: input.VehiclePlate,
		Notes:        input.Notes,
		Status:       entities.EvaluationStatusCreated,
	}

	if err := s.db.Create(evaluation).Error; err != nil {
		return nil, err
	}

	return evaluation, nil
}

func (s *EvaluationService) GetByID(id int) (*entities.Evaluation, error) {
	var evaluation entities.Evaluation
	if err := s.db.First(&evaluation, id).Error; err != nil {
		return nil, err
	}
	return &evaluation, nil
}

func (s *EvaluationService) List(status string) ([]entities.Evaluation, error) {
	var evaluations []entities.Evaluation
	query := s.db.Order("created_at DESC")

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Find(&evaluations).Error; err != nil {
		return nil, err
	}

	return evaluations, nil
}

func (s *EvaluationService) Update(id int, input UpdateEvaluationInput) (*entities.Evaluation, error) {
	evaluation, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	if input.EvaluatorID != nil {
		if evaluation.Status != entities.EvaluationStatusCreated {
			return nil, errors.New("evaluator can only be assigned to evaluations in 'created' status")
		}
		evaluation.EvaluatorID = input.EvaluatorID
		evaluation.Status = entities.EvaluationStatusAccepted
	}

	if input.Status != nil {
		newStatus := entities.EvaluationStatus(*input.Status)
		if !isValidStatusTransition(evaluation.Status, newStatus) {
			return nil, errors.New("invalid status transition")
		}
		evaluation.Status = newStatus
	}

	if input.Notes != nil {
		evaluation.Notes = input.Notes
	}

	if err := s.db.Save(evaluation).Error; err != nil {
		return nil, err
	}

	return evaluation, nil
}

func isValidStatusTransition(current, new entities.EvaluationStatus) bool {
	switch current {
	case entities.EvaluationStatusCreated:
		return new == entities.EvaluationStatusAccepted || new == entities.EvaluationStatusCanceled
	case entities.EvaluationStatusAccepted:
		return new == entities.EvaluationStatusInProgress || new == entities.EvaluationStatusCanceled
	case entities.EvaluationStatusInProgress:
		return new == entities.EvaluationStatusCompleted || new == entities.EvaluationStatusCanceled
	case entities.EvaluationStatusCompleted, entities.EvaluationStatusCanceled:
		return false
	default:
		return false
	}
}
