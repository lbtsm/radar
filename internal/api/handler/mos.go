package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/mapprotocol/filter/internal/api/service"
	"github.com/mapprotocol/filter/internal/api/stream"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type Mos struct {
	srv service.MosSrv
}

func NewMos(db *gorm.DB) *Mos {
	return &Mos{srv: service.NewMosSrv(db)}
}

func (m *Mos) List(c *gin.Context) {
	var req stream.MosListReq
	if err := c.ShouldBindJSON(&req); err != nil {
		WriteResponse(c, err, nil)
		return
	}
	if req.Id == 0 {
		WriteResponse(c, errors.New("param event id is zero"), nil)
		return
	}

	ret, err := m.srv.List(c, &req)
	if err != nil {
		WriteResponse(c, errors.Wrap(err, "get Mos list failed"), nil)
		return
	}
	WriteResponse(c, nil, ret)
}
