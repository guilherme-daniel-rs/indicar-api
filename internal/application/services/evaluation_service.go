package services

import (
	"errors"
	"fmt"
	"indicar-api/internal/domain/entities"
	"indicar-api/internal/infrastructure/aws"
	"time"

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

type EvaluationPhotoService struct {
	db        *gorm.DB
	s3Service *aws.S3Service
}

func NewEvaluationPhotoService(db *gorm.DB) (*EvaluationPhotoService, error) {
	s3Service, err := aws.NewS3Service()
	if err != nil {
		return nil, err
	}

	return &EvaluationPhotoService{
		db:        db,
		s3Service: s3Service,
	}, nil
}

type UploadPhotoInput struct {
	File        []byte
	ContentType string
	SizeBytes   int
}

func (s *EvaluationPhotoService) UploadPhoto(evaluationID int, input UploadPhotoInput) (*entities.EvaluationPhoto, error) {
	var evaluation entities.Evaluation
	if err := s.db.First(&evaluation, evaluationID).Error; err != nil {
		return nil, err
	}

	s3Key := fmt.Sprintf("evaluations/%d/photos/%d.jpg", evaluationID, time.Now().UnixNano())

	if err := s.s3Service.UploadFile(s3Key, input.File, input.ContentType); err != nil {
		return nil, fmt.Errorf("failed to upload file to S3: %w", err)
	}

	photo := &entities.EvaluationPhoto{
		EvaluationID: evaluationID,
		S3Bucket:     s.s3Service.Bucket,
		S3Key:        s3Key,
		ContentType:  &input.ContentType,
		SizeBytes:    &input.SizeBytes,
	}

	if err := s.db.Create(photo).Error; err != nil {
		return nil, err
	}

	return photo, nil
}

func (s *EvaluationPhotoService) ListPhotos(evaluationID int) ([]entities.EvaluationPhoto, error) {
	var photos []entities.EvaluationPhoto

	if err := s.db.Where("evaluation_id = ?", evaluationID).
		Order("created_at DESC").
		Find(&photos).Error; err != nil {
		return nil, err
	}

	return photos, nil
}
