package service

import (
	"context"
	"strings"

	"github.com/mapprotocol/filter/internal/api/store"
	"github.com/mapprotocol/filter/internal/api/store/mysql"
	"github.com/mapprotocol/filter/internal/api/stream"
	"gorm.io/gorm"
)

type MosSrv interface {
	List(context.Context, *stream.MosListReq) (*stream.MosListResp, error)
}

type Mos struct {
	store      store.Moser
	event      store.Evener
	eventCache map[string]int64
}

func NewMosSrv(db *gorm.DB) MosSrv {
	return &Mos{
		store: mysql.NewMos(db), event: mysql.NewEvent(db), eventCache: make(map[string]int64),
	}
}

func (m *Mos) List(ctx context.Context, req *stream.MosListReq) (*stream.MosListResp, error) {
	if req.Limit == 0 || req.Limit > 100 {
		req.Limit = 10
	}

	if req.Id == 0 {
		req.Id = 1
	}

	splits := strings.Split(req.Topic, ",")
	eventIds := make([]int64, 0, len(splits))
	for _, sp := range splits {
		if id, ok := m.eventCache[sp]; ok {
			eventIds = append(eventIds, id)
			continue
		}
		event, err := m.event.Get(ctx, &store.EventCond{Topic: sp})
		if err != nil {
			return nil, err
		}
		m.eventCache[sp] = event.Id
		eventIds = append(eventIds, event.Id)
	}

	list, total, err := m.store.List(ctx, &store.MosCond{
		Id:          req.Id,
		ChainId:     req.ChainId,
		ProjectId:   req.ProjectId,
		EventIds:    eventIds,
		BlockNumber: req.BlockNumber,
		TxHash:      req.TxHash,
		Limit:       req.Limit,
	})
	if err != nil {
		return nil, err
	}
	ret := make([]*stream.GetMosResp, 0, req.Limit)
	for _, ele := range list {
		ret = append(ret, &stream.GetMosResp{
			Id:              ele.Id,
			ProjectId:       ele.ProjectId,
			ChainId:         ele.ChainId,
			EventId:         ele.EventId,
			TxHash:          ele.TxHash,
			ContractAddress: ele.ContractAddress,
			Topic:           ele.Topic,
			BlockNumber:     ele.BlockNumber,
			LogIndex:        ele.LogIndex,
			LogData:         ele.LogData,
			TxIndex:         ele.TxIndex,
			TxTimestamp:     ele.TxTimestamp,
		})
	}
	return &stream.MosListResp{
		Total: total,
		List:  ret,
	}, nil
}
