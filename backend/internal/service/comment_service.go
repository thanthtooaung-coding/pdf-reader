package service

import (
	"errors"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/thanthtooaung-coding/pdf-reader/backend/internal/mapper"
	"github.com/thanthtooaung-coding/pdf-reader/backend/internal/models"
	"github.com/thanthtooaung-coding/pdf-reader/backend/internal/repository"
	"github.com/thanthtooaung-coding/pdf-reader/backend/internal/request"
	"github.com/thanthtooaung-coding/pdf-reader/backend/internal/response"
)

type CommentService interface {
	Create(userID, fileID uuid.UUID, req request.CreateCommentRequest) (*response.CommentResponse, error)
	List(userID, fileID uuid.UUID) ([]response.CommentResponse, error)
}

type commentServiceImpl struct {
	logger  *logrus.Logger
	comments repository.CommentRepository
	files   repository.FileRepository
}

func NewCommentService(
	logger *logrus.Logger,
	comments repository.CommentRepository,
	files repository.FileRepository,
) CommentService {
	return &commentServiceImpl{logger: logger, comments: comments, files: files}
}

func (s *commentServiceImpl) Create(userID, fileID uuid.UUID, req request.CreateCommentRequest) (*response.CommentResponse, error) {
	f, err := s.files.GetByID(fileID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrFileNotFound
		}
		return nil, err
	}
	if f.UserID != userID {
		return nil, ErrForbidden
	}

	c := &models.Comment{
		FileID:  fileID,
		UserID:  userID,
		Message: req.Message,
		Base:    models.Base{IsEnable: true},
	}
	if err := s.comments.Create(c); err != nil {
		return nil, err
	}

	list, err := s.comments.ListByFile(fileID)
	if err != nil {
		resp := mapper.ToCommentResponse(*c)
		return &resp, nil
	}
	for _, item := range list {
		if item.ID == c.ID {
			resp := mapper.ToCommentResponse(item)
			return &resp, nil
		}
	}
	resp := mapper.ToCommentResponse(*c)
	return &resp, nil
}

func (s *commentServiceImpl) List(userID, fileID uuid.UUID) ([]response.CommentResponse, error) {
	f, err := s.files.GetByID(fileID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrFileNotFound
		}
		return nil, err
	}
	if f.UserID != userID {
		return nil, ErrForbidden
	}

	comments, err := s.comments.ListByFile(fileID)
	if err != nil {
		return nil, err
	}
	return mapper.ToCommentResponses(comments), nil
}
