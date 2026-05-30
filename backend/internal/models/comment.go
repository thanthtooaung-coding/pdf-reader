package models

import "github.com/google/uuid"

type Comment struct {
	Base
	FileID  uuid.UUID `gorm:"type:uuid;not null;index" json:"file_id"`
	UserID  uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	Message string    `gorm:"type:text;not null" json:"message"`

	File *File `gorm:"foreignKey:FileID;references:ID" json:"file,omitempty"`
	User *User `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
}

func (Comment) TableName() string { return "comments" }
