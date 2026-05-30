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

var (
	ErrWorkspaceNotFound = errors.New("workspace not found")
	ErrForbidden         = errors.New("forbidden")
)

type WorkspaceService interface {
	Create(userID uuid.UUID, req request.CreateWorkspaceRequest) (*response.WorkspaceResponse, error)
	Get(userID, workspaceID uuid.UUID) (*response.WorkspaceResponse, error)
	List(userID uuid.UUID, page, limit int) (*response.PagedResponse[response.WorkspaceResponse], error)
}

type workspaceServiceImpl struct {
	logger    *logrus.Logger
	workspaces repository.WorkspaceRepository
	files     repository.FileRepository
}

func NewWorkspaceService(
	logger *logrus.Logger,
	workspaces repository.WorkspaceRepository,
	files repository.FileRepository,
) WorkspaceService {
	return &workspaceServiceImpl{logger: logger, workspaces: workspaces, files: files}
}

func (s *workspaceServiceImpl) Create(userID uuid.UUID, req request.CreateWorkspaceRequest) (*response.WorkspaceResponse, error) {
	ws := &models.Workspace{
		UserID: userID,
		Name:   req.Name,
		Base:   models.Base{IsEnable: true},
	}
	if err := s.workspaces.Create(ws); err != nil {
		return nil, err
	}
	resp := mapper.ToWorkspaceResponse(*ws, 0)
	return &resp, nil
}

func (s *workspaceServiceImpl) Get(userID, workspaceID uuid.UUID) (*response.WorkspaceResponse, error) {
	ws, err := s.workspaces.GetByID(workspaceID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrWorkspaceNotFound
		}
		return nil, err
	}
	if ws.UserID != userID {
		return nil, ErrForbidden
	}
	count, _ := s.files.CountByWorkspace(workspaceID)
	resp := mapper.ToWorkspaceResponse(*ws, count)
	return &resp, nil
}

func (s *workspaceServiceImpl) List(userID uuid.UUID, page, limit int) (*response.PagedResponse[response.WorkspaceResponse], error) {
	items, total, err := s.workspaces.ListByUser(userID, page, limit)
	if err != nil {
		return nil, err
	}

	out := make([]response.WorkspaceResponse, 0, len(items))
	for _, ws := range items {
		count, _ := s.files.CountByWorkspace(ws.ID)
		out = append(out, mapper.ToWorkspaceResponse(ws, count))
	}

	totalPages := int(total) / limit
	if int(total)%limit != 0 {
		totalPages++
	}

	return &response.PagedResponse[response.WorkspaceResponse]{
		Items: out,
		Meta: response.PageMeta{
			Page: page, Limit: limit, Total: total, TotalPages: totalPages,
		},
	}, nil
}
