package ethereum

import (
	"github.com/ethereum/go-ethereum/log"
	"github.com/mapprotocol/filter/internal/filter/config"
	"github.com/mapprotocol/filter/internal/pkg/dao"
	"github.com/mapprotocol/filter/internal/pkg/storage"
	"github.com/mapprotocol/filter/pkg/blockstore"
	"github.com/mapprotocol/filter/pkg/keystore"
	"github.com/pkg/errors"
)

type Chain struct {
	conn                     Conner
	log                      log.Logger
	cfg                      *EthConfig
	stop                     chan struct{}
	bs                       blockstore.BlockStorer
	storages                 []storage.Saver
	events                   []*dao.Event
	eventId, currentProgress int64
}

func New(cfg config.RawChainConfig, storages []storage.Saver) (*Chain, error) {
	eCfg, err := parseConfig(cfg)
	if err != nil {
		return nil, err
	}

	kpI, err := keystore.KeypairFromEth(cfg.KeystorePath)
	if err != nil {
		return nil, err
	}

	conn := NewConn(eCfg.Endpoint, kpI)
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
	err := c.getMatch(true)
	if err != nil {
		return errors.Wrap(err, "init getMatch failed")
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
