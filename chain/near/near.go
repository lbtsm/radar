package near

import (
	"github.com/ethereum/go-ethereum/log"
	"github.com/mapprotocol/filter/config"
	"github.com/mapprotocol/filter/pkg/blockstore"
)

type Chain struct {
	log  log.Logger
	cfg  *Config
	stop chan struct{}
	bs   blockstore.BlockStorer
}

func New(cfg config.RawChainConfig) (*Chain, error) {
	eCfg, err := parseConfig(cfg)
	if err != nil {
		return nil, err
	}

	ret := &Chain{
		log:  log.New("chain", eCfg.Name),
		cfg:  eCfg,
		stop: make(chan struct{}),
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
