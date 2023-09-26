package near

import (
	"github.com/mapprotocol/filter/config"
	"strings"
)

type Config struct {
	Name, Id, Endpoint string // Human-readable chain name
	Mcs, Events        []string
}

func parseConfig(cfg config.RawChainConfig) (*Config, error) {
	ret := &Config{
		Name:     cfg.Name,
		Id:       cfg.Id,
		Endpoint: cfg.Endpoint,
	}

	mcs := strings.Split(cfg.Opts.Mcs, ",")
	for _, s := range mcs {
		ret.Mcs = append(ret.Mcs, s)
	}

	vs := strings.Split(cfg.Opts.Event, "|")
	for _, s := range vs {
		ret.Events = append(ret.Events, s)
	}

	return ret, nil
}
