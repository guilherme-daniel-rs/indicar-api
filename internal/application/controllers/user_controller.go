package controllers

import (
	"indicar-api/internal/application/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService *services.UserService
}

func NewUserController(userService *services.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

// @Summary Get current user
// @Description Get the currently authenticated user's profile
// @Tags users
// @Produce json
// @Security Bearer
// @Success 200 {object} entities.User
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /me [get]
func (c *UserController) GetMe(ctx *gin.Context) {
	userID := ctx.GetInt("user_id")
	if userID == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	user, err := c.userService.GetCurrentUser(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

// @Summary Update current user
// @Description Update the currently authenticated user's profile
// @Tags users
// @Accept json
// @Produce json
// @Security Bearer
// @Param input body services.UpdateUserInput true "User update data"
// @Success 200 {object} entities.User
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /me [put]
func (c *UserController) UpdateMe(ctx *gin.Context) {
	userID := ctx.GetInt("user_id")
	if userID == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var input services.UpdateUserInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := c.userService.UpdateUser(userID, input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

// @Summary Get evaluator profile
// @Description Get an evaluator's public profile
// @Tags evaluators
// @Produce json
// @Security Bearer
// @Param id path int true "Evaluator ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /evaluators/{id} [get]
func (c *UserController) GetEvaluator(ctx *gin.Context) {
	evaluatorID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid evaluator ID"})
		return
	}

	user, evaluator, err := c.userService.GetEvaluator(evaluatorID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	response := gin.H{
		"id":            user.ID,
		"full_name":     user.FullName,
		"phone":         user.Phone,
		"role":          user.Role,
		"rating":        evaluator.Rating,
		"total_reviews": evaluator.TotalReviews,
		"bio":           evaluator.Bio,
	}

	ctx.JSON(http.StatusOK, response)
}
