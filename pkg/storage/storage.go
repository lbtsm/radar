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
	Type() string
	Event(uint64, *dao.MosEvent) error
	LatestBlockNumber(chainId string, latest uint64) error
}

func NewSaver(tp, url string) (Saver, error) {
	switch tp {
	case constant.Redis:
		return newRds(url)
	case constant.Mysql:
		return newMysql(url)
	default:
		return nil, ErrorOfStorageType
	}
}
