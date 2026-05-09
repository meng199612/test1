package models

import (
	"time"
)

type OperationLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    *uint     `gorm:"index" json:"user_id"`
	Username  string    `gorm:"size:50" json:"username"`
	Module    string    `gorm:"size:50" json:"module"`
	Operation string    `gorm:"size:100" json:"operation"`
	Method    string    `gorm:"size:10" json:"method"`
	Path      string    `gorm:"size:255" json:"path"`
	IP        string    `gorm:"size:50" json:"ip"`
	Params    string    `gorm:"type:text" json:"params"`
	Result    string    `gorm:"type:text" json:"result"`
	Status    int8      `gorm:"default:1" json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

func (OperationLog) TableName() string {
	return "operation_logs"
}
