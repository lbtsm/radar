package chain

import (
	"errors"
	"github.com/mapprotocol/filter/chain/ethereum"
	"github.com/mapprotocol/filter/chain/near"
	"github.com/mapprotocol/filter/config"
	"github.com/mapprotocol/filter/internal/constant"
	"github.com/mapprotocol/filter/pkg/storage"
)

type Chainer interface {
	Start() error
	Stop()
}

func Init(cfg *config.Config, storages []storage.Saver) ([]Chainer, error) {
	ret := make([]Chainer, 0)
	for _, ccfg := range cfg.Chains {
		var (
			err error
			c   Chainer
		)

		switch ccfg.Type {
		case constant.Ethereum:
			c, err = ethereum.New(ccfg, storages)
		case constant.Near:
			c, err = near.New(ccfg, storages)
		default:
			return nil, errors.New("unrecognized Chain Type")
		}
		if err != nil {
			return nil, err
		}
		ret = append(ret, c)
		constant.OnlineChaId[ccfg.Id] = struct{}{}
	}
	return ret, nil
}
