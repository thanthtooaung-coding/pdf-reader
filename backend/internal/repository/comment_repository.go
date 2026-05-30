package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/thanthtooaung-coding/pdf-reader/backend/internal/models"
)

type CommentRepository interface {
	Create(c *models.Comment) error
	ListByFile(fileID uuid.UUID) ([]models.Comment, error)
}

type commentRepository struct{ db *gorm.DB }

func NewCommentRepository(db *gorm.DB) CommentRepository { return &commentRepository{db: db} }

func (r *commentRepository) Create(c *models.Comment) error { return r.db.Create(c).Error }

func (r *commentRepository) ListByFile(fileID uuid.UUID) ([]models.Comment, error) {
	var comments []models.Comment
	if err := r.db.Preload("User").Where("file_id = ?", fileID).Order("created_at ASC").Find(&comments).Error; err != nil {
		return nil, err
	}
	return comments, nil
}
