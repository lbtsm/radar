package filter

import (
	"github.com/ethereum/go-ethereum/log"
	"github.com/mapprotocol/filter/internal/filter/chain"
	"github.com/mapprotocol/filter/internal/filter/config"
	"github.com/mapprotocol/filter/internal/filter/core"
	"github.com/mapprotocol/filter/internal/pkg/constant"
	"github.com/mapprotocol/filter/internal/pkg/storage"
	"github.com/mapprotocol/filter/pkg/utils"
	"github.com/urfave/cli/v2"
)

var Command = &cli.Command{
	Name: "cli",
	Action: func(cli *cli.Context) error {
		log.Root().SetHandler(log.StdoutHandler)
		cfg, err := config.Local(cli.String(constant.ConfigFileFlag.Name))
		if err != nil {
			return err
		}

		storages := make([]storage.Saver, 0, len(cfg.Storages))
		for _, s := range cfg.Storages {
			ele, err := storage.NewSaver(s.Type, s.Url)
			if err != nil {
				return err
			}
			storages = append(storages, ele)
		}

		utils.Init(cfg.Other.Env, cfg.Other.MonitorUrl)
		chains, err := chain.Init(cfg, storages)
		if err != nil {
			return err
		}
		sysErr := make(chan error)
		c := core.New(sysErr)
		for _, ch := range chains {
			c.AddChain(ch)
		}
		c.Start()
		return nil
	},
}
