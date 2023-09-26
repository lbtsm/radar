package ethereum

import (
	"errors"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/mapprotocol/filter/config"
	"github.com/mapprotocol/filter/constant"
)

type EthConfig struct {
	Name               string // Human-readable chain name
	Id                 string // ChainID
	Endpoint           string // url for rpc endpoint
	StartBlock         *big.Int
	BlockConfirmations *big.Int
	Mcs                []common.Address
	Events             []constant.EventSig
}

func parseConfig(cfg config.RawChainConfig) (*EthConfig, error) {
	ret := &EthConfig{
		Name:               cfg.Name,
		Id:                 cfg.Id,
		Endpoint:           cfg.Endpoint,
		StartBlock:         new(big.Int).SetUint64(0),
		BlockConfirmations: new(big.Int).SetUint64(20),
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

	for _, addr := range strings.Split(cfg.Opts.Mcs, ",") {
		ret.Mcs = append(ret.Mcs, common.HexToAddress(addr))
	}

	vs := strings.Split(cfg.Opts.Event, "|")
	for _, s := range vs {
		ret.Events = append(ret.Events, constant.EventSig(s))
	}

	return ret, nil
}
