package handler

import (
	"github.com/mapprotocol/filter/internal/api/service"
	"gorm.io/gorm"
)

type Mos struct {
	srv service.MosSrv
}

func NewMos(db *gorm.DB) *Mos {
	return &Mos{srv: service.NewMosSrv(db)}
}
