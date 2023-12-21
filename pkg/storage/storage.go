package storage

import (
	"errors"
	"github.com/mapprotocol/filter/internal/constant"
	"github.com/mapprotocol/filter/internal/dao"
)

var (
	ErrorOfStorageType = errors.New("storage type is miss")
)

type Saver interface {
	GetType() string
	Storage(uint64, *dao.MosEvent) error
}

func NewSaver(tp, url string) (Saver, error) {
	switch tp {
	case constant.Redis:
		return newMysql(url)
	case constant.Mysql:
		return newRds(url)
	default:
		return nil, ErrorOfStorageType
	}
}
