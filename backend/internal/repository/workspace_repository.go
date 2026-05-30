package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/thanthtooaung-coding/pdf-reader/backend/internal/models"
)

type WorkspaceRepository interface {
	Create(ws *models.Workspace) error
	GetByID(id uuid.UUID) (*models.Workspace, error)
	ListByUser(userID uuid.UUID, page, limit int) ([]models.Workspace, int64, error)
}

type workspaceRepository struct{ db *gorm.DB }

func NewWorkspaceRepository(db *gorm.DB) WorkspaceRepository {
	return &workspaceRepository{db: db}
}

func (r *workspaceRepository) Create(ws *models.Workspace) error { return r.db.Create(ws).Error }

func (r *workspaceRepository) GetByID(id uuid.UUID) (*models.Workspace, error) {
	var ws models.Workspace
	if err := r.db.First(&ws, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &ws, nil
}

func (r *workspaceRepository) ListByUser(userID uuid.UUID, page, limit int) ([]models.Workspace, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}

	q := r.db.Model(&models.Workspace{}).Where("user_id = ?", userID)
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var items []models.Workspace
	if err := q.Order("updated_at DESC").Offset((page - 1) * limit).Limit(limit).Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}
