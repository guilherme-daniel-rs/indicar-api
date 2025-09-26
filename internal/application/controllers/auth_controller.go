package controllers

import (
	"indicar-api/internal/application/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthController handles authentication related requests
type AuthController struct {
	authService *services.AuthService
}

// NewAuthController creates a new instance of AuthController
func NewAuthController(authService *services.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

// @Summary User signup
// @Description Register a new user
// @Tags auth
// @Accept json
// @Produce json
// @Param input body services.SignupInput true "Signup credentials"
// @Success 201 {object} services.AuthResponse
// @Failure 400 {object} map[string]interface{}
// @Router /auth/signup [post]
func (c *AuthController) Signup(ctx *gin.Context) {
	var input services.SignupInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := c.authService.Signup(input)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, response)
}

// @Summary User login
// @Description Authenticate an existing user
// @Tags auth
// @Accept json
// @Produce json
// @Param input body services.LoginInput true "Login credentials"
// @Success 200 {object} services.AuthResponse
// @Failure 401 {object} map[string]interface{}
// @Router /auth/login [post]
func (c *AuthController) Login(ctx *gin.Context) {
	var input services.LoginInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := c.authService.Login(input)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// @Summary Refresh token
// @Description Get new access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param X-Refresh-Token header string true "Refresh token"
// @Success 200 {object} services.AuthResponse
// @Failure 401 {object} map[string]interface{}
// @Router /auth/refresh [post]
func (c *AuthController) RefreshToken(ctx *gin.Context) {
	refreshToken := ctx.GetHeader("X-Refresh-Token")
	if refreshToken == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "refresh token is required"})
		return
	}

	response, err := c.authService.RefreshToken(refreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}
