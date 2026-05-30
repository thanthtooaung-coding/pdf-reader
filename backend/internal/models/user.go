package models

import "github.com/google/uuid"

type User struct {
	Base
	RoleID   uuid.UUID `gorm:"type:uuid;not null;index" json:"role_id"`
	Fullname string    `gorm:"type:varchar(128)" json:"fullname,omitempty"`
	Username string    `gorm:"type:varchar(64);not null;uniqueIndex" json:"username"`
	Email    string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	Password string    `gorm:"type:varchar(255);not null" json:"-"`

	Role *Role `gorm:"foreignKey:RoleID;references:ID" json:"role,omitempty"`
}

func (User) TableName() string { return "users" }
