package models

import (
	"time"

	"gorm.io/gorm"
)

// BaseModel 基础模型
type BaseModel struct {
	// 主键ID
	ID uint `gorm:"primarykey" json:"id"`
	// 创建时间
	CreatedAt time.Time `json:"created_at"`
	// 更新时间
	UpdatedAt time.Time `json:"updated_at"`
	// 软删除时间（GORM 的 DeletedAt 类型）
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
