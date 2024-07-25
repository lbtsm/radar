package ton

import (
	"strings"

	"github.com/mapprotocol/filter/internal/filter/config"
)

type Config struct {
	Name     string // Human-readable chain name
	Id       string // ChainID
	Endpoint string // url for rpc endpoint
	Mcs      []string
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

	return ret, nil
}
