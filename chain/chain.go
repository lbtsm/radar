package chain

import (
	"errors"
	"github.com/mapprotocol/filter/chain/ethereum"
	"github.com/mapprotocol/filter/chain/near"
	"github.com/mapprotocol/filter/config"
	"github.com/mapprotocol/filter/constant"
)

type Chainer interface {
	Start() error
	Stop()
}

func Init(cfgs *config.Config, backup bool) ([]Chainer, error) {
	ret := make([]Chainer, 0)
	for _, cfg := range cfgs.Chains {
		var (
			err error
			c   Chainer
		)
		switch cfg.Type {
		case constant.Ethereum:
			c, err = ethereum.New(cfg, backup)
		case constant.Near:
			c, err = near.New(cfg, backup)
		default:
			return nil, errors.New("unrecognized Chain Type")
		}
		if err != nil {
			return nil, err
		}
		ret = append(ret, c)
	}
	return ret, nil
}
