package controllers

import (
	"indicar-api/internal/application/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type EvaluationController struct {
	evaluationService *services.EvaluationService
}

func NewEvaluationController(evaluationService *services.EvaluationService) *EvaluationController {
	return &EvaluationController{
		evaluationService: evaluationService,
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
