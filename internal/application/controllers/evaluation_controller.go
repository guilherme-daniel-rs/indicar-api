package controllers

import (
	"indicar-api/internal/application/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type EvaluationController struct {
	evaluationService      *services.EvaluationService
	evaluationPhotoService *services.EvaluationPhotoService
}

func NewEvaluationController(evaluationService *services.EvaluationService, evaluationPhotoService *services.EvaluationPhotoService) *EvaluationController {
	return &EvaluationController{
		evaluationService:      evaluationService,
		evaluationPhotoService: evaluationPhotoService,
	}
}

func (c *EvaluationController) Create(ctx *gin.Context) {
	userID := ctx.GetInt("user_id")
	if userID == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var input services.CreateEvaluationInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	evaluation, err := c.evaluationService.Create(userID, input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, evaluation)
}

func (c *EvaluationController) GetByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid evaluation ID"})
		return
	}

	evaluation, err := c.evaluationService.GetByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, evaluation)
}

func (c *EvaluationController) List(ctx *gin.Context) {
	status := ctx.Query("status")
	evaluations, err := c.evaluationService.List(status)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, evaluations)
}

func (c *EvaluationController) Update(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid evaluation ID"})
		return
	}

	var input services.UpdateEvaluationInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	evaluation, err := c.evaluationService.Update(id, input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, evaluation)
}

func (c *EvaluationController) UploadPhoto(ctx *gin.Context) {
	evaluationID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid evaluation ID"})
		return
	}

	file, err := ctx.FormFile("photo")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "photo file is required"})
		return
	}

	if file.Size > 10*1024*1024 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "file size exceeds 10MB limit"})
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

	input := services.UploadPhotoInput{
		File:        fileBytes,
		ContentType: file.Header.Get("Content-Type"),
		SizeBytes:   int(file.Size),
	}

	photo, err := c.evaluationPhotoService.UploadPhoto(evaluationID, input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, photo)
}

func (c *EvaluationController) ListPhotos(ctx *gin.Context) {
	evaluationID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid evaluation ID"})
		return
	}

	photos, err := c.evaluationPhotoService.ListPhotos(evaluationID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, photos)
}
