package mysql

import (
	"context"
	"github.com/mapprotocol/filter/internal/api/store"
	"github.com/mapprotocol/filter/internal/pkg/dao"
	"gorm.io/gorm"
)

type Block struct {
	db *gorm.DB
}

func NewBlock(db *gorm.DB) *Block {
	return &Block{db: db}
}

func (e *Block) Get(ctx context.Context, c *store.BlockCond) (*dao.Block, error) {
	ret := dao.Block{}
	err := e.db.WithContext(ctx).Where("chain_id = ?", c.ChainId).First(&ret).Error
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

func (e *Block) GetCurrentScan(ctx context.Context, c *store.BlockCond) (*dao.ScanBlock, error) {
	ret := dao.ScanBlock{}
	err := e.db.WithContext(ctx).Where("chain_id = ?", c.ChainId).First(&ret).Error
	if err != nil {
		return nil, err
	}
	return &ret, nil
}
