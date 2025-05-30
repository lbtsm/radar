package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/mapprotocol/filter/internal/api/service"
	"github.com/mapprotocol/filter/internal/api/stream"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type Block struct {
	srv service.BlockSrv
}

func NewBlock(db *gorm.DB) *Block {
	return &Block{srv: service.NewBlockSrv(db)}
}

func (p *Block) Get(c *gin.Context) {
	var req stream.GetBlockReq
	if err := c.ShouldBind(&req); err != nil {
		WriteResponse(c, err, nil)
		return
	}

	if req.ChainId == 0 {
		WriteResponse(c, errors.New("param chain id is empty"), nil)
		return
	}

	ret, err := p.srv.Get(c, &req)
	if err != nil {
		WriteResponse(c, errors.Wrap(err, "get block failed"), nil)
		return
	}
	WriteResponse(c, nil, ret)
}

func (p *Block) GetCurrentScan(c *gin.Context) {
	var req stream.GetBlockReq
	if err := c.ShouldBind(&req); err != nil {
		WriteResponse(c, err, nil)
		return
	}

	if req.ChainId == 0 {
		WriteResponse(c, errors.New("param chain id is empty"), nil)
		return
	}

	ret, err := p.srv.GetCurrentScan(c, &req)
	if err != nil {
		WriteResponse(c, errors.Wrap(err, "get block failed"), nil)
		return
	}
	WriteResponse(c, nil, ret)
}
