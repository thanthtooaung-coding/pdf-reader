package mapper

import (
	"github.com/thanthtooaung-coding/pdf-reader/backend/internal/models"
	"github.com/thanthtooaung-coding/pdf-reader/backend/internal/response"
)

func ToUserResponse(u models.User) response.UserResponse {
	roleCode := ""
	if u.Role != nil {
		roleCode = u.Role.Code
	}
	return response.UserResponse{
		ID:        u.ID,
		RoleID:    u.RoleID,
		RoleCode:  roleCode,
		Fullname:  u.Fullname,
		Username:  u.Username,
		Email:     u.Email,
		IsEnable:  u.IsEnable,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func ToWorkspaceResponse(ws models.Workspace, fileCount int64) response.WorkspaceResponse {
	return response.WorkspaceResponse{
		ID:        ws.ID,
		UserID:    ws.UserID,
		Name:      ws.Name,
		FileCount: fileCount,
		CreatedAt: ws.CreatedAt,
		UpdatedAt: ws.UpdatedAt,
	}
}

func ToFileResponse(f models.File) response.FileResponse {
	return response.FileResponse{
		ID:               f.ID,
		UserID:           f.UserID,
		WorkspaceID:      f.WorkspaceID,
		OriginalFileName: f.OriginalFileName,
		Name:             f.Name,
		URL:              f.URL,
		Type:             f.Type,
		CreatedAt:        f.CreatedAt,
		UpdatedAt:        f.UpdatedAt,
	}
}

func ToFileResponses(files []models.File) []response.FileResponse {
	out := make([]response.FileResponse, 0, len(files))
	for _, f := range files {
		out = append(out, ToFileResponse(f))
	}
	return out
}

func ToCommentResponse(c models.Comment) response.CommentResponse {
	username := ""
	if c.User != nil {
		username = c.User.Username
	}
	return response.CommentResponse{
		ID:        c.ID,
		FileID:    c.FileID,
		UserID:    c.UserID,
		Username:  username,
		Message:   c.Message,
		CreatedAt: c.CreatedAt,
	}
}

func ToCommentResponses(comments []models.Comment) []response.CommentResponse {
	out := make([]response.CommentResponse, 0, len(comments))
	for _, c := range comments {
		out = append(out, ToCommentResponse(c))
	}
	return out
}

func ToAIJobResponse(job models.AIJob) response.AIJobResponse {
	return response.AIJobResponse{
		ID:          job.ID,
		UserID:      job.UserID,
		WorkspaceID: job.WorkspaceID,
		FileID:      job.FileID,
		Type:        job.Type,
		Status:      job.Status,
		Input:       job.Input,
		Output:      job.Output,
		DurationMS:  job.Duration,
		RetryCount:  job.RetryCount,
		CreatedAt:   job.CreatedAt,
		UpdatedAt:   job.UpdatedAt,
	}
}

func ToAIJobResponses(jobs []models.AIJob) []response.AIJobResponse {
	out := make([]response.AIJobResponse, 0, len(jobs))
	for _, j := range jobs {
		out = append(out, ToAIJobResponse(j))
	}
	return out
}
