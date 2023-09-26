package core

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/ethereum/go-ethereum/log"
	"github.com/mapprotocol/filter/chain"
)

type Core struct {
	registry []chain.Chainer
	sysErr   <-chan error
}

func New(sysErr <-chan error) *Core {
	return &Core{
		registry: make([]chain.Chainer, 0),
		sysErr:   sysErr,
	}
}

// AddChain registers the chain in the registry and calls Chain.SetRouter()
func (c *Core) AddChain(chain chain.Chainer) {
	c.registry = append(c.registry, chain)
}

// Start will call all registered chains' Start methods and block forever (or until signal is received)
func (c *Core) Start() {
	for _, r := range c.registry {
		err := r.Start()
		if err != nil {
			return
		}
	}

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigc)

	select {
	case err := <-c.sysErr:
		log.Error("FATAL ERROR. Shutting down.", "err", err)
	case <-sigc:
		log.Warn("Interrupt received, shutting down now.")
	}

	for _, r := range c.registry {
		r.Stop()
	}
}

func (c *Core) Errors() <-chan error {
	return c.sysErr
}
