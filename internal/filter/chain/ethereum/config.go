package ethereum

import (
	"errors"
	"github.com/mapprotocol/filter/internal/filter/config"
	"math/big"
)

type EthConfig struct {
	Name               string // Human-readable chain name
	Id                 string // ChainID
	Endpoint           string // url for rpc endpoint
	StartBlock         *big.Int
	BlockConfirmations *big.Int
	Range              *big.Int
}

func parseConfig(cfg config.RawChainConfig) (*EthConfig, error) {
	ret := &EthConfig{
		Name:               cfg.Name,
		Id:                 cfg.Id,
		Endpoint:           cfg.Endpoint,
		Range:              big.NewInt(0),
		StartBlock:         new(big.Int).SetUint64(0),
		BlockConfirmations: new(big.Int).SetUint64(config.DefaultBlockConfirm),
	}
	if cfg.Opts.StartBlock != "" {
		sb, ok := new(big.Int).SetString(cfg.Opts.StartBlock, 10)
		if !ok {
			return nil, errors.New("startBlock format failed")
		}
		ret.StartBlock = sb
	}

	if cfg.Opts.BlockConfirmations != "" {
		bf, ok := new(big.Int).SetString(cfg.Opts.BlockConfirmations, 10)
		if !ok {
			return nil, errors.New("blockConfirmations format failed")
		}
		ret.BlockConfirmations = bf
	}

	if cfg.Opts.Range != "" {
		bf, ok := new(big.Int).SetString(cfg.Opts.Range, 10)
		if !ok {
			return nil, errors.New("blockConfirmations format failed")
		}
		ret.Range = bf
	}

	return ret, nil
}
