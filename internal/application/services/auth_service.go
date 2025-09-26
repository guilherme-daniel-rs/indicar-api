package services

import (
	"errors"
	"indicar-api/configs"
	"indicar-api/internal/domain/entities"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	db          *gorm.DB
	jwtSecret   []byte
	tokenExpiry time.Duration
}

func NewAuthService(db *gorm.DB) *AuthService {
	return &AuthService{
		db:          db,
		jwtSecret:   []byte(configs.Get().JWT.Secret),
		tokenExpiry: 24 * time.Hour,
	}
}

type SignupInput struct {
	FullName   string  `json:"full_name" binding:"required"`
	Email      string  `json:"email" binding:"required,email"`
	Password   string  `json:"password" binding:"required,min=6"`
	Phone      string  `json:"phone"`
	Role       string  `json:"role" binding:"required,oneof=user evaluator"`
	DocumentID *string `json:"document_id,omitempty"`
	Bio        *string `json:"bio,omitempty"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	AccessToken  string         `json:"access_token"`
	RefreshToken string         `json:"refresh_token"`
	User         *entities.User `json:"user"`
}

func (s *AuthService) Signup(input SignupInput) (*AuthResponse, error) {
	var existingUser entities.User
	if result := s.db.Where("email = ?", input.Email).First(&existingUser); result.Error == nil {
		return nil, errors.New("email already registered")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Start a transaction since we might need to create multiple records
	err = s.db.Transaction(func(tx *gorm.DB) error {
		user := &entities.User{
			FullName:     input.FullName,
			Email:        input.Email,
			PasswordHash: string(hashedPassword),
			Phone:        &input.Phone,
			Role:         entities.UserRole(input.Role),
			IsActive:     true,
		}

		if err := tx.Create(user).Error; err != nil {
			return err
		}

		if user.Role == entities.UserRoleEvaluator {
			evaluator := &entities.Evaluator{
				UserID:     user.ID,
				DocumentID: input.DocumentID,
				Bio:        input.Bio,
			}

			if err := tx.Create(evaluator).Error; err != nil {
				return err
			}
		}

		existingUser = *user
		return nil
	})

	if err != nil {
		return nil, err
	}

	return s.generateTokens(&existingUser)
}

func (s *AuthService) Login(input LoginInput) (*AuthResponse, error) {
	var user entities.User
	if err := s.db.Where("email = ?", input.Email).First(&user).Error; err != nil {
		return nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	return s.generateTokens(&user)
}

func (s *AuthService) RefreshToken(refreshToken string) (*AuthResponse, error) {
	var token entities.AuthRefreshToken
	if err := s.db.Where("token = ? AND revoked = ? AND expires_at > ?", refreshToken, false, time.Now()).First(&token).Error; err != nil {
		return nil, errors.New("invalid refresh token")
	}

	var user entities.User
	if err := s.db.First(&user, token.UserID).Error; err != nil {
		return nil, errors.New("user not found")
	}

	// Revoke the used refresh token
	s.db.Model(&token).Update("revoked", true)

	return s.generateTokens(&user)
}

func (s *AuthService) generateTokens(user *entities.User) (*AuthResponse, error) {
	// Generate access token
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
		"exp":     time.Now().Add(s.tokenExpiry).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return nil, err
	}

	// Generate refresh token
	refreshToken := &entities.AuthRefreshToken{
		UserID:    user.ID,
		Token:     generateRandomToken(),
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		Revoked:   false,
	}

	if err := s.db.Create(refreshToken).Error; err != nil {
		return nil, err
	}

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken.Token,
		User:         user,
	}, nil
}

func generateRandomToken() string {
	// In a real application, implement a secure random token generation
	return "random-token-" + time.Now().Format("20060102150405")
}
