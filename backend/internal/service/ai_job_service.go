package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/thanthtooaung-coding/pdf-reader/backend/internal/mapper"
	"github.com/thanthtooaung-coding/pdf-reader/backend/internal/models"
	"github.com/thanthtooaung-coding/pdf-reader/backend/internal/repository"
	"github.com/thanthtooaung-coding/pdf-reader/backend/internal/request"
	"github.com/thanthtooaung-coding/pdf-reader/backend/internal/response"
)

var ErrAIJobNotFound = errors.New("ai job not found")

type AIJobService interface {
	Create(userID uuid.UUID, req request.CreateAIJobRequest) (*response.AIJobResponse, error)
	Get(userID, jobID uuid.UUID) (*response.AIJobResponse, error)
	ListByWorkspace(userID, workspaceID uuid.UUID) ([]response.AIJobResponse, error)
}

type aiJobServiceImpl struct {
	logger     *logrus.Logger
	jobs       repository.AIJobRepository
	files      repository.FileRepository
	workspaces repository.WorkspaceRepository
}

func NewAIJobService(
	logger *logrus.Logger,
	jobs repository.AIJobRepository,
	files repository.FileRepository,
	workspaces repository.WorkspaceRepository,
) AIJobService {
	return &aiJobServiceImpl{
		logger: logger, jobs: jobs, files: files, workspaces: workspaces,
	}
}

func (s *aiJobServiceImpl) Create(userID uuid.UUID, req request.CreateAIJobRequest) (*response.AIJobResponse, error) {
	f, err := s.files.GetByID(req.FileID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrFileNotFound
		}
		return nil, err
	}
	if f.UserID != userID {
		return nil, ErrForbidden
	}

	ws, err := s.workspaces.GetByID(f.WorkspaceID)
	if err != nil || ws.UserID != userID {
		return nil, ErrForbidden
	}

	job := &models.AIJob{
		UserID:      userID,
		WorkspaceID: f.WorkspaceID,
		FileID:      req.FileID,
		Type:        req.Type,
		Status:      models.AIJobStatusPending,
		Input:       req.Input,
		Base:        models.Base{IsEnable: true},
	}
	if err := s.jobs.Create(job); err != nil {
		return nil, err
	}

	go s.processJob(job.ID)

	resp := mapper.ToAIJobResponse(*job)
	return &resp, nil
}

func (s *aiJobServiceImpl) Get(userID, jobID uuid.UUID) (*response.AIJobResponse, error) {
	job, err := s.jobs.GetByID(jobID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrAIJobNotFound
		}
		return nil, err
	}
	if job.UserID != userID {
		return nil, ErrForbidden
	}
	resp := mapper.ToAIJobResponse(*job)
	return &resp, nil
}

func (s *aiJobServiceImpl) ListByWorkspace(userID, workspaceID uuid.UUID) ([]response.AIJobResponse, error) {
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

	jobs, err := s.jobs.ListByWorkspace(workspaceID)
	if err != nil {
		return nil, err
	}
	return mapper.ToAIJobResponses(jobs), nil
}

func (s *aiJobServiceImpl) processJob(jobID uuid.UUID) {
	start := time.Now()

	job, err := s.jobs.GetByID(jobID)
	if err != nil {
		return
	}

	job.Status = models.AIJobStatusProcessing
	_ = s.jobs.Update(job)

	time.Sleep(500 * time.Millisecond)

	var output string
	switch job.Type {
	case models.AIJobTypeTranslate:
		lang := job.Input
		if lang == "" {
			lang = "en"
		}
		output = fmt.Sprintf("[stub] Translated document to %s. Connect an AI provider to replace this output.", lang)
	case models.AIJobTypeSummarize:
		output = "[stub] Summary of the PDF content. Connect an AI provider to replace this output."
	case models.AIJobTypeComment:
		output = "[stub] AI-generated commentary on the document. Connect an AI provider to replace this output."
	default:
		job.Status = models.AIJobStatusFailed
		job.Output = "unsupported job type"
		_ = s.jobs.Update(job)
		return
	}

	job.Status = models.AIJobStatusCompleted
	job.Output = output
	job.Duration = time.Since(start).Milliseconds()
	_ = s.jobs.Update(job)
}
