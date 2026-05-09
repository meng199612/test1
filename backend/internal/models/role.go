package models

import (
	"time"

	"gorm.io/gorm"
)

type Role struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:50;not null" json:"name"`
	Code        string         `gorm:"uniqueIndex;size:50;not null" json:"code"`
	Description string         `gorm:"size:255" json:"description"`
	Status      int8           `gorm:"default:1" json:"status"`
	Permissions []Permission   `gorm:"many2many:role_permissions;" json:"permissions,omitempty"`
	Menus       []Menu         `gorm:"many2many:role_menus;" json:"menus,omitempty"`
	Users       []User         `gorm:"many2many:user_roles;" json:"-"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
