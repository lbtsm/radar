package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/mapprotocol/filter/internal/api/service"
	"github.com/mapprotocol/filter/internal/api/stream"
	"github.com/mapprotocol/filter/internal/pkg/constant"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type Event struct {
	srv service.EventSrv
}

func NewEvent(db *gorm.DB) *Event {
	return &Event{srv: service.NewEventSrv(db)}
}

func (p *Event) Add(c *gin.Context) {
	var req stream.AddEventReq
	if err := c.ShouldBindJSON(&req); err != nil {
		WriteResponse(c, err, nil)
		return
	}

	if req.Format == "" {
		WriteResponse(c, errors.New("param format is empty"), nil)
		return
	}
	if req.Address == "" || req.Address == constant.ZeroAddress.String() {
		WriteResponse(c, errors.New("param address is empty"), nil)
		return
	}
	if req.ProjectId == 0 {
		WriteResponse(c, errors.New("param project is zero"), nil)
		return
	}
	if req.BlockNumber != "" && req.ChainId == 0 {
		WriteResponse(c, errors.New("appoint blockNumber must appoint chain_id"), nil)
		return
	}

	err := p.srv.Add(c, &req)
	if err != nil {
		WriteResponse(c, errors.Wrap(err, "add Event failed"), nil)
		return
	}
	WriteResponse(c, nil, nil)
}

func (p *Event) Get(c *gin.Context) {
	var req stream.GetEventReq
	if err := c.ShouldBindJSON(&req); err != nil {
		WriteResponse(c, err, nil)
		return
	}

	ret, err := p.srv.Get(c, &req)
	if err != nil {
		WriteResponse(c, errors.Wrap(err, "get Event failed"), nil)
		return
	}
	WriteResponse(c, nil, ret)
}

func (p *Event) Delete(c *gin.Context) {
	var req stream.DelEventReq
	if err := c.ShouldBindJSON(&req); err != nil {
		WriteResponse(c, err, nil)
		return
	}
	if req.Id == 0 {
		WriteResponse(c, errors.New("param event id is zero"), nil)
		return
	}

	err := p.srv.Del(c, &req)
	if err != nil {
		WriteResponse(c, errors.Wrap(err, "del Event failed"), nil)
		return
	}
	WriteResponse(c, nil, nil)
}

func (p *Event) List(c *gin.Context) {
	var req stream.EventListReq
	if err := c.ShouldBind(&req); err != nil {
		WriteResponse(c, err, nil)
		return
	}
	if req.Id == 0 {
		WriteResponse(c, errors.New("Param event id is zero"), nil)
		return
	}

	ret, err := p.srv.List(c, &req)
	if err != nil {
		WriteResponse(c, errors.Wrap(err, "Get Event list failed"), nil)
		return
	}
	WriteResponse(c, nil, ret)
}
