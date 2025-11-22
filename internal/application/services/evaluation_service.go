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

type EvaluationResponse struct {
	ID           int                       `json:"id"`
	RequesterID  int                       `json:"requester_id"`
	EvaluatorID  *int                      `json:"evaluator_id,omitempty"`
	CityID       int                       `json:"city_id"`
	VehicleMake  string                    `json:"vehicle_make"`
	VehicleModel string                    `json:"vehicle_model"`
	VehicleYear  *int                      `json:"vehicle_year,omitempty"`
	VehiclePlate *string                   `json:"vehicle_plate,omitempty"`
	Status       entities.EvaluationStatus `json:"status"`
	Notes        *string                   `json:"notes,omitempty"`
	PhotosCount  int                       `json:"photos_count"`
	CreatedAt    time.Time                 `json:"created_at"`
	UpdatedAt    time.Time                 `json:"updated_at"`
}

func (s *EvaluationService) List(status string) ([]EvaluationResponse, error) {
	var evaluations []entities.Evaluation
	query := s.db.Order("created_at DESC")

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Find(&evaluations).Error; err != nil {
		return nil, err
	}

	// Build response with photos count
	responses := make([]EvaluationResponse, len(evaluations))
	for i, eval := range evaluations {
		var photosCount int64
		s.db.Model(&entities.EvaluationPhoto{}).Where("evaluation_id = ?", eval.ID).Count(&photosCount)

		responses[i] = EvaluationResponse{
			ID:           eval.ID,
			RequesterID:  eval.RequesterID,
			EvaluatorID:  eval.EvaluatorID,
			CityID:       eval.CityID,
			VehicleMake:  eval.VehicleMake,
			VehicleModel: eval.VehicleModel,
			VehicleYear:  eval.VehicleYear,
			VehiclePlate: eval.VehiclePlate,
			Status:       eval.Status,
			Notes:        eval.Notes,
			PhotosCount:  int(photosCount),
			CreatedAt:    eval.CreatedAt,
			UpdatedAt:    eval.UpdatedAt,
		}
	}

	return responses, nil
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

	// Validate file type (only images allowed)
	allowedTypes := []string{"image/jpeg", "image/jpg", "image/png", "image/gif", "image/webp"}
	if err := s.s3Service.ValidateFileType("photo.jpg", allowedTypes); err != nil {
		return nil, fmt.Errorf("invalid file type: %w", err)
	}

	// Validate file size (max 10MB for photos)
	maxSize := int64(10 * 1024 * 1024) // 10MB
	if err := s.s3Service.ValidateFileSize(int64(input.SizeBytes), maxSize); err != nil {
		return nil, fmt.Errorf("file too large: %w", err)
	}

	// Generate appropriate file extension based on content type
	var ext string
	switch input.ContentType {
	case "image/jpeg", "image/jpg":
		ext = "jpg"
	case "image/png":
		ext = "png"
	case "image/gif":
		ext = "gif"
	case "image/webp":
		ext = "webp"
	default:
		ext = "jpg" // default fallback
	}

	s3Key := fmt.Sprintf("evaluations/%d/photos/%d.%s", evaluationID, time.Now().UnixNano(), ext)

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
		// If database creation fails, clean up S3 file
		s.s3Service.DeleteFile(s3Key)
		return nil, err
	}

	return photo, nil
}

type PhotoResponse struct {
	ID           int       `json:"id"`
	EvaluationID int       `json:"evaluation_id"`
	S3Bucket     string    `json:"s3_bucket"`
	S3Key        string    `json:"s3_key"`
	ContentType  *string   `json:"content_type,omitempty"`
	SizeBytes    *int      `json:"size_bytes,omitempty"`
	URL          string    `json:"url"` // Presigned URL
	CreatedAt    time.Time `json:"created_at"`
}

func (s *EvaluationPhotoService) ListPhotos(evaluationID int) ([]PhotoResponse, error) {
	var photos []entities.EvaluationPhoto

	if err := s.db.Where("evaluation_id = ?", evaluationID).
		Order("created_at DESC").
		Find(&photos).Error; err != nil {
		return nil, err
	}

	// Generate presigned URLs for each photo
	responses := make([]PhotoResponse, len(photos))
	for i, photo := range photos {
		url, err := s.s3Service.GetPresignedURL(photo.S3Key, 24*time.Hour) // 24 hours expiration
		if err != nil {
			// If presigned URL generation fails, use empty string (log error but don't fail entire request)
			url = ""
		}

		responses[i] = PhotoResponse{
			ID:           photo.ID,
			EvaluationID: photo.EvaluationID,
			S3Bucket:     photo.S3Bucket,
			S3Key:        photo.S3Key,
			ContentType:  photo.ContentType,
			SizeBytes:    photo.SizeBytes,
			URL:          url,
			CreatedAt:    photo.CreatedAt,
		}
	}

	return responses, nil
}

// GetPhotoURL returns a presigned URL for a specific photo
func (s *EvaluationPhotoService) GetPhotoURL(photoID int) (string, error) {
	var photo entities.EvaluationPhoto
	if err := s.db.First(&photo, photoID).Error; err != nil {
		return "", errors.New("photo not found")
	}

	url, err := s.s3Service.GetPresignedURL(photo.S3Key, 24*time.Hour)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return url, nil
}
