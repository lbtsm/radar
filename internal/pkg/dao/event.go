package dao

import "time"

type Event struct {
	Id          int64     `gorm:"column:id;primary_key;AUTO_INCREMENT;NOT NULL" gorm:"id" json:"id"`
	ProjectId   int64     `gorm:"column:project_id;default:NULL" gorm:"project_id" json:"project_id"`
	Format      string    `gorm:"column:format" gorm:"format" json:"format"`
	Topic       string    `gorm:"column:topic;default:NULL" gorm:"topic" json:"topic"`
	BlockNumber string    `gorm:"column:block_number;default:NULL" gorm:"block_number" json:"block_number"`
	CreateAt    time.Time `gorm:"column:create_at;default:CURRENT_TIMESTAMP;NOT NULL" gorm:"create_at" json:"create_at"`
	DeletedAt   time.Time `gorm:"column:deleted_at;default:CURRENT_TIMESTAMP;NOT NULL" gorm:"deleted_at" json:"deleted_at"`
}

func (e *Event) TableName() string {
	return "event"
}
