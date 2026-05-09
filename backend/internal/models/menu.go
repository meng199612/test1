package models

import (
	"time"

	"gorm.io/gorm"
)

type Menu struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"size:50;not null" json:"name"`
	Path      string         `gorm:"size:255;not null" json:"path"`
	Icon      string         `gorm:"size:50" json:"icon"`
	ParentID  uint           `gorm:"default:0" json:"parent_id"`
	Component string         `gorm:"size:255" json:"component"`
	Sort      int            `gorm:"default:0" json:"sort"`
	Visible   int8           `gorm:"default:1" json:"visible"`
	Children  []*Menu        `gorm:"-" json:"children,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
