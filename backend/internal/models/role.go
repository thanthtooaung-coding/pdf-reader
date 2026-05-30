package models

type Role struct {
	Base
	Name string `gorm:"type:varchar(64);not null;uniqueIndex" json:"name"`
	Code string `gorm:"type:varchar(64);not null;uniqueIndex" json:"code"`
}

func (Role) TableName() string { return "roles" }
