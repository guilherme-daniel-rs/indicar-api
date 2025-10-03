package controllers

import (
	"indicar-api/internal/application/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ReportController struct {
	reportService *services.ReportService
}

func NewReportController(reportService *services.ReportService) *ReportController {
	return &ReportController{
		reportService: reportService,
	}
}

// @Summary Create or update report
// @Description Create a new report or update an existing one
// @Tags reports
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int false "Report ID (required for updates)"
// @Param input body services.CreateReportInput true "Report data"
// @Success 200,201 {object} entities.Report
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /reports [post]
// @Router /reports/{id} [patch]
func (c *ReportController) CreateOrUpdate(ctx *gin.Context) {
	userID := ctx.GetInt("user_id")
	if userID == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Check if report ID is provided
	if reportID := ctx.Param("id"); reportID != "" {
		// Update existing report
		id, err := strconv.Atoi(reportID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid report ID"})
			return
		}

		var input services.UpdateReportInput
		if err := ctx.ShouldBindJSON(&input); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		report, err := c.reportService.Update(id, userID, input)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, report)
		return
	}

	// Create new report
	var input services.CreateReportInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	report, err := c.reportService.Create(userID, input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, report)
}

// @Summary Get report by ID
// @Description Get a specific report by its ID
// @Tags reports
// @Produce json
// @Security Bearer
// @Param id path int true "Report ID"
// @Success 200 {object} entities.Report
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /reports/{id} [get]
func (c *ReportController) GetByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid report ID"})
		return
	}

	report, err := c.reportService.GetByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, report)
}

// @Summary Get report file URL
// @Description Get a pre-signed URL for downloading the report file
// @Tags reports
// @Produce json
// @Security Bearer
// @Param id path int true "Report ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /reports/{id}/file [get]
func (c *ReportController) GetFileURL(ctx *gin.Context) {
	userID := ctx.GetInt("user_id")
	if userID == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	reportID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid report ID"})
		return
	}

	fileURL, err := c.reportService.GetReportFileURL(reportID, userID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"url": fileURL})
}

// @Summary Upload report file
// @Description Upload a PDF file for a report
// @Tags reports
// @Accept multipart/form-data
// @Produce json
// @Security Bearer
// @Param id path int true "Report ID"
// @Param file formData file true "PDF file (max 50MB)"
// @Success 201 {object} entities.ReportFile
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /reports/{id}/file [post]
func (c *ReportController) UploadFile(ctx *gin.Context) {
	userID := ctx.GetInt("user_id")
	if userID == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	reportID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid report ID"})
		return
	}

	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}

	// Validate file size (max 50MB)
	if file.Size > 50*1024*1024 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "file size exceeds 50MB limit"})
		return
	}

	// Validate file type (only PDF)
	contentType := file.Header.Get("Content-Type")
	if contentType != "application/pdf" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "only PDF files are allowed"})
		return
	}

	f, err := file.Open()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error reading file"})
		return
	}
	defer f.Close()

	fileBytes := make([]byte, file.Size)
	_, err = f.Read(fileBytes)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error reading file"})
		return
	}

	input := services.UploadReportFileInput{
		File:        fileBytes,
		ContentType: contentType,
		SizeBytes:   int(file.Size),
		Filename:    file.Filename,
	}

	reportFile, err := c.reportService.UploadReportFile(reportID, userID, input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, reportFile)
}
