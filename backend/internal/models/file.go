package models

import "github.com/google/uuid"

type File struct {
	Base
	UserID           uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	WorkspaceID      uuid.UUID `gorm:"type:uuid;not null;index" json:"workspace_id"`
	OriginalFileName string    `gorm:"type:varchar(512);not null" json:"original_file_name"`
	Name             string    `gorm:"type:varchar(512);not null" json:"name"`
	URL              string    `gorm:"type:varchar(1024);not null" json:"url"`
	Type             string    `gorm:"type:varchar(32);not null;default:pdf" json:"type"`

	User      *User      `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
	Workspace *Workspace `gorm:"foreignKey:WorkspaceID;references:ID" json:"workspace,omitempty"`
	Comments  []Comment  `gorm:"foreignKey:FileID" json:"comments,omitempty"`
}

func (File) TableName() string { return "files" }
