package ethereum

import (
	"github.com/ethereum/go-ethereum/log"
	"github.com/go-redis/redis/v8"
	"github.com/mapprotocol/filter/config"
	"github.com/mapprotocol/filter/pkg/blockstore"
	"github.com/mapprotocol/filter/pkg/storage"
)

type Chain struct {
	conn     Conner
	log      log.Logger
	cfg      *EthConfig
	stop     chan struct{}
	rdb      *redis.Client
	bs       blockstore.BlockStorer
	storages []storage.Saver
}

func New(cfg config.RawChainConfig, storages []storage.Saver) (*Chain, error) {
	eCfg, err := parseConfig(cfg)
	if err != nil {
		return nil, err
	}

	conn := NewConn(eCfg.Endpoint)
	err = conn.Connect()
	if err != nil {
		return nil, err
	}
	bs, err := blockstore.New(blockstore.PathPostfix, eCfg.Id)
	if err != nil {
		return nil, err
	}

	opt, err := redis.ParseURL(cfg.Opts.Redis)
	if err != nil {
		panic(err)
	}
	ret := &Chain{
		conn:     conn,
		log:      log.New("chain", eCfg.Name),
		cfg:      eCfg,
		stop:     make(chan struct{}),
		bs:       bs,
		rdb:      redis.NewClient(opt),
		storages: storages,
	}
	ret.log.SetHandler(log.LvlFilterHandler(log.LvlInfo, log.StdoutHandler))

	return ret, nil
}

func (c *Chain) Start() error {
	go func() {
		err := c.sync()
		if err != nil {
			c.log.Error("Polling blocks failed", "err", err)
		}
	}()
	c.log.Info("Starting filter ...")
	return nil
}

func (c *Chain) Stop() {
	close(c.stop)
}
