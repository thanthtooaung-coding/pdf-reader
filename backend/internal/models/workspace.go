package models

import "github.com/google/uuid"

type Workspace struct {
	Base
	UserID uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	Name   string    `gorm:"type:varchar(255);not null" json:"name"`

	User  *User  `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
	Files []File `gorm:"foreignKey:WorkspaceID" json:"files,omitempty"`
}

func (Workspace) TableName() string { return "workspaces" }
