package xrp

import (
	"github.com/pkg/errors"
	"math/big"
	"strings"

	"github.com/mapprotocol/filter/internal/filter/config"
)

type Config struct {
	Name               string // Human-readable chain name
	Id                 string // ChainID
	Endpoint           string // url for rpc endpoint
	Mcs                []string
	Event              []string
	StartBlock         *big.Int
	BlockConfirmations *big.Int
	Range              *big.Int
}

func parseConfig(cfg config.RawChainConfig) (*Config, error) {
	ret := &Config{
		Name:               cfg.Name,
		Id:                 cfg.Id,
		Endpoint:           cfg.Endpoint,
		StartBlock:         new(big.Int).SetUint64(0),
		BlockConfirmations: new(big.Int).SetUint64(config.DefaultBlockConfirm),
	}

	mcs := strings.Split(cfg.Opts.Mcs, ",")
	for _, s := range mcs {
		ret.Mcs = append(ret.Mcs, s)
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

	event := strings.Split(cfg.Opts.Event, ",")
	for _, s := range event {
		ret.Event = append(ret.Event, s)
	}

	return ret, nil
}
