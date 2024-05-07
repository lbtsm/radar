package mysql

import (
	"context"
	"github.com/mapprotocol/filter/internal/api/store"
	"github.com/mapprotocol/filter/internal/pkg/dao"
	"gorm.io/gorm"
)

type Project struct {
	db *gorm.DB
}

func NewProject(db *gorm.DB) *Project {
	return &Project{db: db}
}

func (e *Project) Create(ctx context.Context, ele *dao.Project) error {
	return e.db.WithContext(ctx).Create(ele).Error
}

func (e *Project) Delete(ctx context.Context, id int64) error {
	return e.db.WithContext(ctx).Where("id = ?", id).Delete(&dao.Project{}).Error
}

func (e *Project) Get(ctx context.Context, c *store.ProjectCond) (*dao.Project, error) {
	db := e.db.WithContext(ctx)
	if c.Id != 0 {
		db = db.Where("id = ?", c.Id)
	}
	if c.Name != "" {
		db = db.Where("name = ?", c.Name)
	}
	ret := dao.Project{}
	err := db.First(&ret).Error
	return &ret, err
}

func (e *Project) List(ctx context.Context, c *store.ProjectCond) ([]*dao.Project, error) {
	db := e.db.WithContext(ctx)
	if c.Id != 0 {
		db = db.Where("id > ?", c.Id)
	}
	if c.Name != "" {
		db = db.Where("name = ?", c.Name)
	}
	ret := make([]*dao.Project, 0)
	err := db.Find(&ret).Error
	return ret, err
}
