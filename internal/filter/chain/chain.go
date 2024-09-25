package chain

import (
	"github.com/mapprotocol/filter/internal/filter/chain/ethereum"
	"github.com/mapprotocol/filter/internal/filter/chain/near"
	"github.com/mapprotocol/filter/internal/filter/chain/ton"
	"github.com/mapprotocol/filter/internal/filter/config"
	"github.com/mapprotocol/filter/internal/pkg/constant"
	"github.com/mapprotocol/filter/internal/pkg/storage"
)

type Chainer interface {
	Start() error
	Stop()
}

func Init(cfg *config.Config, storages []storage.Saver, latest bool) ([]Chainer, error) {
	ret := make([]Chainer, 0)
	for _, ccfg := range cfg.Chains {
		var (
			err error
			c   Chainer
		)

		if ccfg.KeystorePath == "" {
			ccfg.KeystorePath = cfg.KeystorePath
		}
		switch ccfg.Type {
		case constant.Near:
			c, err = near.New(ccfg, storages)
		case constant.Ton:
			c, err = ton.New(ccfg, storages)
		default:
			c, err = ethereum.New(ccfg, storages, latest)
		}
		if err != nil {
			return nil, err
		}
		ret = append(ret, c)
		constant.OnlineChaId[ccfg.Id] = struct{}{}
	}
	return ret, nil
}
