package service

import (
	"context"
	"github.com/mapprotocol/filter/internal/api/store"
	"github.com/mapprotocol/filter/internal/api/store/mysql"
	"github.com/mapprotocol/filter/internal/api/stream"
	"gorm.io/gorm"
)

type BlockSrv interface {
	Get(context.Context, *stream.GetBlockReq) (string, error)
}

type Block struct {
	store store.Blocker
}

func NewBlockSrv(db *gorm.DB) BlockSrv {
	return &Block{store: mysql.NewBlock(db)}
}

func (b *Block) Get(ctx context.Context, req *stream.GetBlockReq) (string, error) {
	block, err := b.store.Get(ctx, &store.BlockCond{ChainId: req.ChainId})
	if err != nil {
		return "", err
	}
	return block.Number, nil
}
