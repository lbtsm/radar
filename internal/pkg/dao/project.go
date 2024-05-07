package dao

import (
	"time"
)

type Project struct {
	Id          int64     `gorm:"column:id;primary_key;AUTO_INCREMENT;NOT NULL" gorm:"id" json:"id"`
	Name        string    `gorm:"column:name" gorm:"name" json:"name"`
	Description string    `gorm:"column:description;default:NULL" gorm:"description" json:"description"`
	CreatedAt   time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP;NOT NULL" gorm:"created_at" json:"created_at"`
	DeletedAt   time.Time `gorm:"column:deleted_at;default:CURRENT_TIMESTAMP;NOT NULL" gorm:"deleted_at" json:"deleted_at"`
}

func (p *Project) TableName() string {
	return "project"
}
