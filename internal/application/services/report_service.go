package services

import (
	"errors"
	"fmt"
	"indicar-api/internal/domain/entities"
	"indicar-api/internal/infrastructure/aws"
	"time"

	"gorm.io/gorm"
)

type ReportService struct {
	db        *gorm.DB
	s3Service *aws.S3Service
}

func NewReportService(db *gorm.DB) (*ReportService, error) {
	s3Service, err := aws.NewS3Service()
	if err != nil {
		return nil, err
	}

	return &ReportService{
		db:        db,
		s3Service: s3Service,
	}, nil
}

type CreateReportInput struct {
	EvaluationID int     `json:"evaluation_id" binding:"required"`
	Summary      *string `json:"summary"`
}

type UpdateReportInput struct {
	Summary *string                `json:"summary"`
	Status  *entities.ReportStatus `json:"status"`
}

type UploadReportFileInput struct {
	File        []byte
	ContentType string
	SizeBytes   int
	Filename    string
}

func (s *ReportService) Create(evaluatorID int, input CreateReportInput) (*entities.Report, error) {
	var evaluation entities.Evaluation
	if err := s.db.First(&evaluation, input.EvaluationID).Error; err != nil {
		return nil, errors.New("evaluation not found")
	}

	if evaluation.EvaluatorID == nil || *evaluation.EvaluatorID != evaluatorID {
		return nil, errors.New("unauthorized: only the assigned evaluator can create a report")
	}

	report := &entities.Report{
		EvaluationID: input.EvaluationID,
		EvaluatorID:  evaluatorID,
		Summary:      input.Summary,
		Status:       entities.ReportStatusDraft,
	}

	if err := s.db.Create(report).Error; err != nil {
		return nil, err
	}

	return report, nil
}

func (s *ReportService) GetByID(id int) (*entities.Report, error) {
	var report entities.Report
	if err := s.db.First(&report, id).Error; err != nil {
		return nil, err
	}
	return &report, nil
}

func (s *ReportService) Update(id int, evaluatorID int, input UpdateReportInput) (*entities.Report, error) {
	report, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	if report.EvaluatorID != evaluatorID {
		return nil, errors.New("unauthorized: only the report's evaluator can update it")
	}

	if input.Summary != nil {
		report.Summary = input.Summary
	}

	if input.Status != nil {
		if !isValidReportStatusTransition(report.Status, *input.Status) {
			return nil, errors.New("invalid status transition")
		}
		report.Status = *input.Status
	}

	if err := s.db.Save(report).Error; err != nil {
		return nil, err
	}

	return report, nil
}

func (s *ReportService) UploadReportFile(reportID int, evaluatorID int, input UploadReportFileInput) (*entities.ReportFile, error) {
	allowedTypes := []string{"pdf", "application/pdf"}
	if err := s.s3Service.ValidateFileType(input.Filename, allowedTypes); err != nil {
		return nil, fmt.Errorf("invalid file type: %w", err)
	}

	maxSize := int64(50 * 1024 * 1024) // 50MB
	if err := s.s3Service.ValidateFileSize(int64(input.SizeBytes), maxSize); err != nil {
		return nil, fmt.Errorf("file too large: %w", err)
	}

	var report entities.Report
	if err := s.db.First(&report, reportID).Error; err != nil {
		return nil, errors.New("report not found")
	}

	if report.EvaluatorID != evaluatorID {
		return nil, errors.New("unauthorized: only the report's evaluator can upload files")
	}

	var existingFile entities.ReportFile
	if err := s.db.Where("report_id = ?", reportID).First(&existingFile).Error; err == nil {
		if err := s.s3Service.DeleteFile(existingFile.S3Key); err != nil {
			return nil, fmt.Errorf("failed to delete existing file: %w", err)
		}
		if err := s.db.Delete(&existingFile).Error; err != nil {
			return nil, fmt.Errorf("failed to delete existing file record: %w", err)
		}
	}

	s3Key := fmt.Sprintf("reports/%d/report_%d.pdf", reportID, time.Now().UnixNano())

	if err := s.s3Service.UploadFile(s3Key, input.File, input.ContentType); err != nil {
		return nil, fmt.Errorf("failed to upload file to S3: %w", err)
	}

	reportFile := &entities.ReportFile{
		ReportID:    reportID,
		S3Bucket:    s.s3Service.Bucket,
		S3Key:       s3Key,
		ContentType: input.ContentType,
		SizeBytes:   &input.SizeBytes,
	}

	if err := s.db.Create(reportFile).Error; err != nil {
		s.s3Service.DeleteFile(s3Key)
		return nil, fmt.Errorf("failed to create report file record: %w", err)
	}

	return reportFile, nil
}

func (s *ReportService) GetReportFileURL(reportID int, evaluatorID int) (string, error) {
	var reportFile entities.ReportFile
	if err := s.db.Where("report_id = ?", reportID).First(&reportFile).Error; err != nil {
		return "", errors.New("report file not found")
	}

	var report entities.Report
	if err := s.db.First(&report, reportID).Error; err != nil {
		return "", errors.New("report not found")
	}

	if report.EvaluatorID != evaluatorID {
		return "", errors.New("unauthorized: only the report's evaluator can access the file")
	}

	url, err := s.s3Service.GetPresignedURL(reportFile.S3Key, time.Hour)
	if err != nil {
		return "", fmt.Errorf("failed to generate pre-signed URL: %w", err)
	}

	return url, nil
}

func isValidReportStatusTransition(current, new entities.ReportStatus) bool {
	switch current {
	case entities.ReportStatusDraft:
		return new == entities.ReportStatusFinalized
	case entities.ReportStatusFinalized:
		return false
	default:
		return false
	}
}
