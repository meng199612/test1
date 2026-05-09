package models

import (
	"time"

	"gorm.io/gorm"
)

type Permission struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"size:50;not null" json:"name"`
	Code      string         `gorm:"size:100;not null" json:"code"`
	Type      int8           `gorm:"default:1" json:"type"`
	ParentID  uint           `gorm:"default:0" json:"parent_id"`
	Path      string         `gorm:"size:255" json:"path"`
	Method    string         `gorm:"size:10" json:"method"`
	Sort      int            `gorm:"default:0" json:"sort"`
	Children  []*Permission  `gorm:"-" json:"children,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
