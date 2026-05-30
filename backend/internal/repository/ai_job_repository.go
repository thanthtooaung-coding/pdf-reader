package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/thanthtooaung-coding/pdf-reader/backend/internal/models"
)

type AIJobRepository interface {
	Create(job *models.AIJob) error
	GetByID(id uuid.UUID) (*models.AIJob, error)
	Update(job *models.AIJob) error
	ListByWorkspace(workspaceID uuid.UUID) ([]models.AIJob, error)
}

type aiJobRepository struct{ db *gorm.DB }

func NewAIJobRepository(db *gorm.DB) AIJobRepository { return &aiJobRepository{db: db} }

func (r *aiJobRepository) Create(job *models.AIJob) error { return r.db.Create(job).Error }

func (r *aiJobRepository) GetByID(id uuid.UUID) (*models.AIJob, error) {
	var job models.AIJob
	if err := r.db.First(&job, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &job, nil
}

func (r *aiJobRepository) Update(job *models.AIJob) error { return r.db.Save(job).Error }

func (r *aiJobRepository) ListByWorkspace(workspaceID uuid.UUID) ([]models.AIJob, error) {
	var jobs []models.AIJob
	if err := r.db.Where("workspace_id = ?", workspaceID).Order("created_at DESC").Find(&jobs).Error; err != nil {
		return nil, err
	}
	return jobs, nil
}
