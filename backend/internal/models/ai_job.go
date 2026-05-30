package models

import "github.com/google/uuid"

const (
	AIJobTypeTranslate = "TRANSLATE"
	AIJobTypeSummarize = "SUMMARIZE"
	AIJobTypeComment   = "COMMENT"

	AIJobStatusPending    = "pending"
	AIJobStatusProcessing = "processing"
	AIJobStatusCompleted  = "completed"
	AIJobStatusFailed     = "failed"
)

type AIJob struct {
	Base
	UserID      uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	WorkspaceID uuid.UUID `gorm:"type:uuid;not null;index" json:"workspace_id"`
	FileID      uuid.UUID `gorm:"type:uuid;not null;index" json:"file_id"`
	Type        string    `gorm:"type:varchar(32);not null" json:"type"`
	Status      string    `gorm:"type:varchar(32);not null;default:pending" json:"status"`
	Input       string    `gorm:"type:text" json:"input,omitempty"`
	Output      string    `gorm:"type:text" json:"output,omitempty"`
	Duration    int64     `gorm:"not null;default:0" json:"duration_ms"`
	RetryCount  int       `gorm:"not null;default:0" json:"retry_count"`

	User      *User      `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
	Workspace *Workspace `gorm:"foreignKey:WorkspaceID;references:ID" json:"workspace,omitempty"`
	File      *File      `gorm:"foreignKey:FileID;references:ID" json:"file,omitempty"`
}

func (AIJob) TableName() string { return "ai_jobs" }
