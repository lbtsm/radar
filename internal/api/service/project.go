package service

import (
	"context"
	"github.com/mapprotocol/filter/internal/api/store"
	"github.com/mapprotocol/filter/internal/api/store/mysql"
	"github.com/mapprotocol/filter/internal/api/stream"
	"github.com/mapprotocol/filter/internal/pkg/dao"
	"gorm.io/gorm"
	"time"
)

type ProjectSrv interface {
	Add(context.Context, *stream.AddProjectReq) error
	Get(context.Context, *stream.GetProjectReq) (*stream.GetProjectResp, error)
}

type Project struct {
	store store.Projector
}

func NewProSrv(db *gorm.DB) ProjectSrv {
	return &Project{store: mysql.NewProject(db)}
}

func (p *Project) Add(ctx context.Context, req *stream.AddProjectReq) error {
	err := p.store.Create(ctx, &dao.Project{
		Name:        req.Name,
		Description: req.Desc,
		CreatedAt:   time.Now(),
	})
	return err
}

func (p *Project) Get(ctx context.Context, req *stream.GetProjectReq) (*stream.GetProjectResp, error) {
	pro, err := p.store.Get(ctx, &store.ProjectCond{
		Id:   req.Id,
		Name: req.Name,
	})
	if err != nil {
		return nil, err
	}
	ret := &stream.GetProjectResp{
		Id:      pro.Id,
		Name:    pro.Name,
		Desc:    pro.Description,
		Created: pro.CreatedAt.Unix(),
	}
	return ret, nil
}
