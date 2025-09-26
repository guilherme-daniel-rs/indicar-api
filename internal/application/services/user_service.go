package services

import (
	"errors"
	"indicar-api/internal/domain/entities"

	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{
		db: db,
	}
}

func (s *UserService) GetCurrentUser(userID int) (*entities.User, error) {
	var user entities.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) UpdateUser(userID int, input UpdateUserInput) (*entities.User, error) {
	user, err := s.GetCurrentUser(userID)
	if err != nil {
		return nil, err
	}

	// Update allowed fields
	if input.FullName != "" {
		user.FullName = input.FullName
	}
	if input.Phone != nil {
		user.Phone = input.Phone
	}

	if err := s.db.Save(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) GetEvaluator(evaluatorID int) (*entities.User, *entities.Evaluator, error) {
	var user entities.User
	var evaluator entities.Evaluator

	err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&user, evaluatorID).Error; err != nil {
			return err
		}

		if user.Role != entities.UserRoleEvaluator {
			return errors.New("user is not an evaluator")
		}

		if err := tx.Where("user_id = ?", evaluatorID).First(&evaluator).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, nil, err
	}

	return &user, &evaluator, nil
}

type UpdateUserInput struct {
	FullName string  `json:"full_name"`
	Phone    *string `json:"phone"`
}
