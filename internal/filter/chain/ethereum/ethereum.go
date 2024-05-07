package ethereum

import (
	"fmt"
	"github.com/ethereum/go-ethereum/log"
	"github.com/mapprotocol/filter/internal/filter/config"
	"github.com/mapprotocol/filter/internal/pkg/dao"
	"github.com/mapprotocol/filter/internal/pkg/storage"
	"github.com/mapprotocol/filter/pkg/blockstore"
	"github.com/pkg/errors"
)

type Chain struct {
	eventId  int64
	conn     Conner
	log      log.Logger
	cfg      *EthConfig
	stop     chan struct{}
	bs       blockstore.BlockStorer
	storages []storage.Saver
	events   []*dao.Event
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

	ret := &Chain{
		conn:     conn,
		log:      log.New("chain", eCfg.Name),
		cfg:      eCfg,
		stop:     make(chan struct{}),
		bs:       bs,
		storages: storages,
		events:   make([]*dao.Event, 0),
	}
	ret.log.SetHandler(log.LvlFilterHandler(log.LvlInfo, log.StdoutHandler))

	return ret, nil
}

func (c *Chain) Start() error {
	for _, s := range c.storages {
		events, err := s.GetEvent(0)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("%s get events failed", s.Type()))
		}
		if len(events) == 0 {
			continue
		}
		c.events = append(c.events, events...)
		c.eventId = events[len(events)-1].Id
	}
	go func() {
		err := c.sync()
		if err != nil {
			c.log.Error("Polling blocks failed", "err", err)
		}
	}()
	go func() {
		err := c.renewEvent()
		if err != nil {
			c.log.Error("Renew event failed", "err", err)
		}
	}()
	c.log.Info("Starting filter ...")
	return nil
}

func (c *Chain) Stop() {
	close(c.stop)
}
