package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/thanthtooaung-coding/pdf-reader/backend/internal/models"
)

type UserRepository interface {
	Create(user *models.User) error
	GetByID(id uuid.UUID) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	GetByUsername(username string) (*models.User, error)
	Update(user *models.User) error
}

type userRepository struct{ db *gorm.DB }

func NewUserRepository(db *gorm.DB) UserRepository { return &userRepository{db: db} }

func (r *userRepository) Create(u *models.User) error { return r.db.Create(u).Error }

func (r *userRepository) GetByID(id uuid.UUID) (*models.User, error) {
	var u models.User
	if err := r.db.Preload("Role").First(&u, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *userRepository) GetByEmail(email string) (*models.User, error) {
	var u models.User
	if err := r.db.Preload("Role").First(&u, "email = ?", email).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *userRepository) GetByUsername(username string) (*models.User, error) {
	var u models.User
	if err := r.db.Preload("Role").First(&u, "username = ?", username).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *userRepository) Update(u *models.User) error { return r.db.Save(u).Error }
