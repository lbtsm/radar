package ethereum

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	vr "github.com/go-redis/redis/v8"
	"github.com/mapprotocol/filter/config"
	"github.com/mapprotocol/filter/constant"
	"github.com/mapprotocol/filter/pkg/redis"
	"math/big"
	"strings"
)

type EthConfig struct {
	Name               string // Human-readable chain name
	Id                 string // ChainID
	Endpoint           string // url for rpc endpoint
	StartBlock         *big.Int
	BlockConfirmations *big.Int
	Mcs                []common.Address
	Events             []constant.EventSig
	BackUp             bool
}

func parseConfig(cfg config.RawChainConfig, backup bool) (*EthConfig, error) {
	ret := &EthConfig{
		Name:               cfg.Name,
		Id:                 cfg.Id,
		Endpoint:           cfg.Endpoint,
		StartBlock:         new(big.Int).SetUint64(0),
		BlockConfirmations: new(big.Int).SetUint64(20),
		BackUp:             backup,
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
	if backup {
		ele := constant.BackUpEvent{
			Event:   vs,
			Address: strings.Split(cfg.Opts.Mcs, ","),
			ChainId: ret.Id,
		}
		data, _ := json.Marshal(&ele)
		err := redis.GetClient().Set(context.Background(), fmt.Sprintf(constant.FlagOfBackUpEvent, ret.Id), string(data), 0).Err()
		if err != nil {
			return nil, err
		}
	} else {
		extra, err := redis.GetClient().Get(context.Background(), fmt.Sprintf(constant.FlagOfAddEvent, ret.Id)).Result()
		if err != nil && !errors.Is(err, vr.Nil) {
			return nil, err
		}
		if extra == "" {
			return ret, nil
		}
		log.Info("Init Config Get backup event", "extraEvent", extra)
		bu := constant.BackUpEvent{}
		err = json.Unmarshal([]byte(extra), &bu)
		if err != nil {
			log.Error("Failed to Unmarshal", "data", extra, "err", err)
			return nil, err
		}
		for _, a := range bu.Address {
			ret.Mcs = append(ret.Mcs, common.HexToAddress(a))
		}
		for _, s := range bu.Event {
			ret.Events = append(ret.Events, constant.EventSig(s))
		}
	}

	return ret, nil
}
