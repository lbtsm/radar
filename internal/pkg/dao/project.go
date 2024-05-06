package dao

import "time"

type Project struct {
	Id          int64     `gorm:"column:id;primary_key;AUTO_INCREMENT;NOT NULL" gorm:"id" json:"id"`
	Name        string    `gorm:"column:name" gorm:"name" json:"name"`
	Description string    `gorm:"column:description;default:NULL" gorm:"description" json:"description"`
	CreateAt    time.Time `gorm:"column:create_at;default:CURRENT_TIMESTAMP;NOT NULL" gorm:"create_at" json:"create_at"`
	DeletedAt   time.Time `gorm:"column:deleted_at;default:CURRENT_TIMESTAMP;NOT NULL" gorm:"deleted_at" json:"deleted_at"`
}

func (p *Project) TableName() string {
	return "project"
}
