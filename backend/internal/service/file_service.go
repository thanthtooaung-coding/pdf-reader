package service

import (
	"errors"
	"io"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/thanthtooaung-coding/pdf-reader/backend/internal/mapper"
	"github.com/thanthtooaung-coding/pdf-reader/backend/internal/models"
	"github.com/thanthtooaung-coding/pdf-reader/backend/internal/repository"
	"github.com/thanthtooaung-coding/pdf-reader/backend/internal/response"
	"github.com/thanthtooaung-coding/pdf-reader/backend/pkg/storage"
)

var (
	ErrFileNotFound    = errors.New("file not found")
	ErrInvalidFileType = errors.New("only pdf files are allowed")
)

type FileService interface {
	Upload(userID, workspaceID uuid.UUID, originalName string, r io.Reader) (*response.FileResponse, error)
	Get(userID, fileID uuid.UUID) (*response.FileResponse, error)
	ListByWorkspace(userID, workspaceID uuid.UUID) ([]response.FileResponse, error)
	DownloadPath(userID, fileID uuid.UUID) (string, string, error)
}

type fileServiceImpl struct {
	logger     *logrus.Logger
	files      repository.FileRepository
	workspaces repository.WorkspaceRepository
	storage    *storage.LocalStorage
}

func NewFileService(
	logger *logrus.Logger,
	files repository.FileRepository,
	workspaces repository.WorkspaceRepository,
	store *storage.LocalStorage,
) FileService {
	return &fileServiceImpl{
		logger: logger, files: files, workspaces: workspaces, storage: store,
	}
}

func (s *fileServiceImpl) Upload(userID, workspaceID uuid.UUID, originalName string, r io.Reader) (*response.FileResponse, error) {
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

	ext := strings.ToLower(filepath.Ext(originalName))
	if ext != ".pdf" {
		return nil, ErrInvalidFileType
	}

	storedName, urlPath, err := s.storage.Save(userID, originalName, r)
	if err != nil {
		if strings.Contains(err.Error(), "only pdf") {
			return nil, ErrInvalidFileType
		}
		return nil, err
	}

	f := &models.File{
		UserID:           userID,
		WorkspaceID:      workspaceID,
		OriginalFileName: originalName,
		Name:             storedName,
		URL:              urlPath,
		Type:             "pdf",
		Base:             models.Base{IsEnable: true},
	}
	if err := s.files.Create(f); err != nil {
		return nil, err
	}

	resp := mapper.ToFileResponse(*f)
	return &resp, nil
}

func (s *fileServiceImpl) Get(userID, fileID uuid.UUID) (*response.FileResponse, error) {
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
	resp := mapper.ToFileResponse(*f)
	return &resp, nil
}

func (s *fileServiceImpl) ListByWorkspace(userID, workspaceID uuid.UUID) ([]response.FileResponse, error) {
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

	files, err := s.files.ListByWorkspace(workspaceID)
	if err != nil {
		return nil, err
	}
	return mapper.ToFileResponses(files), nil
}

func (s *fileServiceImpl) DownloadPath(userID, fileID uuid.UUID) (string, string, error) {
	f, err := s.files.GetByID(fileID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", "", ErrFileNotFound
		}
		return "", "", err
	}
	if f.UserID != userID {
		return "", "", ErrForbidden
	}

	path, err := s.storage.Open(userID, f.Name)
	if err != nil {
		return "", "", ErrFileNotFound
	}
	return path, f.OriginalFileName, nil
}
