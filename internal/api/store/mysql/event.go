package mysql

import (
	"context"
	"github.com/mapprotocol/filter/internal/api/store"
	"github.com/mapprotocol/filter/internal/pkg/dao"
	"gorm.io/gorm"
)

type Event struct {
	db *gorm.DB
}

func NewEvent(db *gorm.DB) *Event {
	return &Event{db: db}
}

func (e *Event) Create(ctx context.Context, ele *dao.Event) error {
	return e.db.WithContext(ctx).Create(ele).Error
}

func (e *Event) Delete(ctx context.Context, id int64) error {
	return e.db.WithContext(ctx).Where("id = ?", id).Delete(&dao.Event{}).Error
}

func (e *Event) Get(ctx context.Context, c *store.EventCond) (*dao.Event, error) {
	db := e.db.WithContext(ctx)
	if c.Id != 0 {
		db = db.Where("id = ?", c.Id)
	}
	if c.ProjectId != 0 {
		db = db.Where("project_id = ?", c.ProjectId)
	}
	if c.Format != "" {
		db = db.Where("format = ?", c.Format)
	}
	if c.Topic != "" {
		db = db.Where("topic = ?", c.Topic)
	}
	ret := dao.Event{}
	err := db.First(&ret).Error
	return &ret, err
}

func (e *Event) List(ctx context.Context, c *store.EventCond) ([]*dao.Event, int64, error) {
	db := e.db.WithContext(ctx)
	if c.Id != 0 {
		db = db.Where("id > ?", c.Id)
	}
	if c.ProjectId != 0 {
		db = db.Where("project_id = ?", c.ProjectId)
	}
	if c.Format != "" {
		db = db.Where("format = ?", c.Format)
	}
	if c.Topic != "" {
		db = db.Where("topic = ?", c.Topic)
	}
	ret := make([]*dao.Event, 0)
	total := int64(0)
	err := db.Model(&dao.Event{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	err = db.Find(&ret).Error
	if err != nil {
		return nil, 0, err
	}
	return ret, total, err
}
