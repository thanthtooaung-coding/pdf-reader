package service

import (
	"errors"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/thanthtooaung-coding/pdf-reader/backend/internal/mapper"
	"github.com/thanthtooaung-coding/pdf-reader/backend/internal/repository"
	"github.com/thanthtooaung-coding/pdf-reader/backend/internal/response"
)

var ErrUserNotFound = errors.New("user not found")

type UserService interface {
	GetMe(userID uuid.UUID) (*response.UserResponse, error)
}

type userServiceImpl struct {
	logger *logrus.Logger
	users  repository.UserRepository
}

func NewUserService(logger *logrus.Logger, users repository.UserRepository) UserService {
	return &userServiceImpl{logger: logger, users: users}
}

func (s *userServiceImpl) GetMe(userID uuid.UUID) (*response.UserResponse, error) {
	u, err := s.users.GetByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	resp := mapper.ToUserResponse(*u)
	return &resp, nil
}
