package mysql

import (
	"context"
	"github.com/mapprotocol/filter/internal/api/store"
	"github.com/mapprotocol/filter/internal/pkg/dao"
	"gorm.io/gorm"
)

type Mos struct {
	db *gorm.DB
}

func NewMos(db *gorm.DB) *Mos {
	return &Mos{db: db}
}

func (m *Mos) Create(ctx context.Context, ele *dao.Mos) error {
	return m.db.WithContext(ctx).Create(ele).Error
}

func (m *Mos) Delete(ctx context.Context, id int64) error {
	return m.db.WithContext(ctx).Where("id = ?", id).Delete(&dao.Mos{}).Error
}

func (m *Mos) Get(ctx context.Context, c *store.MosCond) (*dao.Mos, error) {
	db := m.db.WithContext(ctx)
	if c.Id != 0 {
		db = db.Where("id = ?", c.Id)
	}
	if c.ChainId != 0 {
		db = db.Where("chain_id = ?", c.ChainId)
	}
	if c.ProjectId != 0 {
		db = db.Where("project_id = ?", c.ProjectId)
	}
	if c.BlockNumber != 0 {
		db = db.Where("block_number = ?", c.BlockNumber)
	}
	if c.EventId != 0 {
		db = db.Where("event_id = ?", c.EventId)
	}
	if c.TxHash != "" {
		db = db.Where("tx_hash = ?", c.TxHash)
	}
	ret := dao.Mos{}
	err := db.First(&ret).Error
	return &ret, err
}

func (m *Mos) List(ctx context.Context, c *store.MosCond) ([]*dao.Mos, int64, error) {
	db := m.db.WithContext(ctx)
	if c.Id != 0 {
		db = db.Where("id > ?", c.Id)
	}
	if c.BlockNumber != 0 {
		db = db.Where("block_number >= ?", c.BlockNumber)
	}
	if c.ChainId != 0 {
		db = db.Where("chain_id = ?", c.ChainId)
	}
	if c.ProjectId != 0 {
		db = db.Where("project_id = ?", c.ProjectId)
	}
	if c.EventId != 0 {
		db = db.Where("event_id = ?", c.EventId)
	}
	if c.TxHash != "" {
		db = db.Where("tx_hash = ?", c.TxHash)
	}
	total := int64(0)
	err := db.Model(&dao.Event{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	ret := make([]*dao.Mos, 0)
	err = db.Find(&ret).Limit(c.Limit).Error
	return ret, total, err
}
