package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/mapprotocol/filter/internal/api/service"
	"github.com/mapprotocol/filter/internal/api/stream"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type Project struct {
	srv service.ProjectSrv
}

func NewProject(db *gorm.DB) *Project {
	return &Project{srv: service.NewProSrv(db)}
}

func (p *Project) Add(c *gin.Context) {
	var req stream.AddProjectReq

	if err := c.ShouldBindJSON(&req); err != nil {
		WriteResponse(c, err, nil)
		return
	}

	if req.Name == "" {
		WriteResponse(c, errors.New("param  name is empty"), nil)
		return
	}

	err := p.srv.Add(c, &req)
	if err != nil {
		WriteResponse(c, errors.Wrap(err, "add project failed"), nil)
		return
	}
	WriteResponse(c, nil, nil)
}

func (p *Project) Get(c *gin.Context) {
	var req stream.GetProjectReq

	if err := c.ShouldBindJSON(&req); err != nil {
		WriteResponse(c, err, nil)
		return
	}

	if req.Name == "" {
		WriteResponse(c, errors.New("param  name is empty"), nil)
		return
	}

	ret, err := p.srv.Get(c, &req)
	if err != nil {
		WriteResponse(c, errors.Wrap(err, "add project failed"), nil)
		return
	}
	WriteResponse(c, nil, ret)
}
