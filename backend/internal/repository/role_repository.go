package repository

import (
	"gorm.io/gorm"

	"github.com/thanthtooaung-coding/pdf-reader/backend/internal/models"
)

type RoleRepository interface {
	GetByCode(code string) (*models.Role, error)
}

type roleRepository struct{ db *gorm.DB }

func NewRoleRepository(db *gorm.DB) RoleRepository { return &roleRepository{db: db} }

func (r *roleRepository) GetByCode(code string) (*models.Role, error) {
	var role models.Role
	if err := r.db.First(&role, "code = ?", code).Error; err != nil {
		return nil, err
	}
	return &role, nil
}
