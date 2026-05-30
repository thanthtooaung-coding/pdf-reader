package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/thanthtooaung-coding/pdf-reader/backend/internal/models"
)

type FileRepository interface {
	Create(f *models.File) error
	GetByID(id uuid.UUID) (*models.File, error)
	ListByWorkspace(workspaceID uuid.UUID) ([]models.File, error)
	CountByWorkspace(workspaceID uuid.UUID) (int64, error)
}

type fileRepository struct{ db *gorm.DB }

func NewFileRepository(db *gorm.DB) FileRepository { return &fileRepository{db: db} }

func (r *fileRepository) Create(f *models.File) error { return r.db.Create(f).Error }

func (r *fileRepository) GetByID(id uuid.UUID) (*models.File, error) {
	var f models.File
	if err := r.db.Preload("Workspace").First(&f, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &f, nil
}

func (r *fileRepository) ListByWorkspace(workspaceID uuid.UUID) ([]models.File, error) {
	var files []models.File
	if err := r.db.Where("workspace_id = ?", workspaceID).Order("created_at DESC").Find(&files).Error; err != nil {
		return nil, err
	}
	return files, nil
}

func (r *fileRepository) CountByWorkspace(workspaceID uuid.UUID) (int64, error) {
	var count int64
	if err := r.db.Model(&models.File{}).Where("workspace_id = ?", workspaceID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
