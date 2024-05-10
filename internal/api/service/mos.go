package service

import (
	"context"
	"github.com/mapprotocol/filter/internal/api/store"
	"github.com/mapprotocol/filter/internal/api/store/mysql"
	"github.com/mapprotocol/filter/internal/api/stream"
	"gorm.io/gorm"
)

type MosSrv interface {
	List(context.Context, *stream.MosListReq) (*stream.MosListResp, error)
}

type Mos struct {
	store store.Moser
}

func NewMosSrv(db *gorm.DB) MosSrv {
	return &Mos{store: mysql.NewMos(db)}
}

func (m *Mos) List(ctx context.Context, req *stream.MosListReq) (*stream.MosListResp, error) {
	if req.Limit == 0 || req.Limit > 100 {
		req.Limit = 10
	}
	list, total, err := m.store.List(ctx, &store.MosCond{
		Id:          req.Id,
		ChainId:     req.ChainId,
		ProjectId:   req.ProjectId,
		EventId:     0,
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
			TxTimestamp:     ele.TxTimestamp,
		})
	}
	return &stream.MosListResp{
		Total: total,
		List:  ret,
	}, nil
}
