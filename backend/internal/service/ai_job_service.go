package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/thanthtooaung-coding/pdf-reader/backend/internal/mapper"
	"github.com/thanthtooaung-coding/pdf-reader/backend/internal/models"
	"github.com/thanthtooaung-coding/pdf-reader/backend/internal/repository"
	"github.com/thanthtooaung-coding/pdf-reader/backend/internal/request"
	"github.com/thanthtooaung-coding/pdf-reader/backend/internal/response"
	"github.com/thanthtooaung-coding/pdf-reader/backend/pkg/openai"
	pdfextract "github.com/thanthtooaung-coding/pdf-reader/backend/pkg/pdf"
	"github.com/thanthtooaung-coding/pdf-reader/backend/pkg/storage"
)

const maxDocumentChars = 12000

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
	storage    *storage.LocalStorage
	openai     openai.Client
}

func NewAIJobService(
	logger *logrus.Logger,
	jobs repository.AIJobRepository,
	files repository.FileRepository,
	workspaces repository.WorkspaceRepository,
	store *storage.LocalStorage,
	openaiClient openai.Client,
) AIJobService {
	return &aiJobServiceImpl{
		logger:     logger,
		jobs:       jobs,
		files:      files,
		workspaces: workspaces,
		storage:    store,
		openai:     openaiClient,
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
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	job, err := s.jobs.GetByID(jobID)
	if err != nil {
		return
	}

	job.Status = models.AIJobStatusProcessing
	_ = s.jobs.Update(job)

	output, err := s.runJob(ctx, job)
	job.Duration = time.Since(start).Milliseconds()

	if err != nil {
		s.logger.WithError(err).WithField("job_id", jobID).Error("ai job failed")
		job.Status = models.AIJobStatusFailed
		job.Output = err.Error()
		_ = s.jobs.Update(job)
		return
	}

	job.Status = models.AIJobStatusCompleted
	job.Output = output
	_ = s.jobs.Update(job)
}

func (s *aiJobServiceImpl) runJob(ctx context.Context, job *models.AIJob) (string, error) {
	switch job.Type {
	case models.AIJobTypeTranslate, models.AIJobTypeSummarize:
		return s.runOpenAIJob(ctx, job)
	case models.AIJobTypeComment:
		return "[stub] AI-generated commentary on the document.", nil
	default:
		return "", fmt.Errorf("unsupported job type")
	}
}

func (s *aiJobServiceImpl) runOpenAIJob(ctx context.Context, job *models.AIJob) (string, error) {
	documentText, err := s.loadDocumentText(job)
	if err != nil {
		if !s.openai.Enabled() {
			return s.stubOutput(job)
		}
		return "", err
	}

	documentText = truncateText(documentText, maxDocumentChars)

	switch job.Type {
	case models.AIJobTypeSummarize:
		return s.openai.Summarize(ctx, documentText)
	case models.AIJobTypeTranslate:
		return s.openai.Translate(ctx, documentText, job.Input)
	default:
		return "", fmt.Errorf("unsupported job type")
	}
}

func (s *aiJobServiceImpl) loadDocumentText(job *models.AIJob) (string, error) {
	f, err := s.files.GetByID(job.FileID)
	if err != nil {
		return "", fmt.Errorf("load file: %w", err)
	}

	path, err := s.storage.Open(f.UserID, f.Name)
	if err != nil {
		return "", fmt.Errorf("open stored pdf: %w", err)
	}

	return pdfextract.ExtractText(path)
}

func (s *aiJobServiceImpl) stubOutput(job *models.AIJob) (string, error) {
	switch job.Type {
	case models.AIJobTypeSummarize:
		return s.openai.Summarize(context.Background(), "")
	case models.AIJobTypeTranslate:
		return s.openai.Translate(context.Background(), "", job.Input)
	default:
		return "", fmt.Errorf("unsupported job type")
	}
}

func truncateText(text string, maxChars int) string {
	if maxChars <= 0 || utf8.RuneCountInString(text) <= maxChars {
		return text
	}

	runes := []rune(text)
	return strings.TrimSpace(string(runes[:maxChars])) + "\n\n[truncated]"
}
