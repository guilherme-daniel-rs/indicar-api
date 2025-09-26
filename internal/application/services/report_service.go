package services

import (
	"errors"
	"indicar-api/internal/domain/entities"
	"indicar-api/internal/infrastructure/aws"

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

	// Generate pre-signed URL for the report file
	url := s.s3Service.GetFileURL(reportFile.S3Key)
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
