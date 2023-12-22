package main

import (
	"github.com/ethereum/go-ethereum/log"
	"github.com/mapprotocol/filter/chain"
	"github.com/mapprotocol/filter/config"
	"github.com/mapprotocol/filter/core"
	"github.com/mapprotocol/filter/internal/constant"
	"github.com/mapprotocol/filter/pkg/storage"
	"github.com/mapprotocol/filter/pkg/utils"
	"github.com/urfave/cli/v2"
	"os"
)

var (
	app = cli.NewApp()
)

func main() {
	app.Copyright = "Copyright 2023 MAP Protocol 2023 Authors"
	app.Name = "filter"
	app.Usage = "Filter"
	app.Authors = []*cli.Author{{Name: "MAP Protocol 2023"}}
	app.Version = "1.0.0"
	app.EnableBashCompletion = true
	app.Flags = append(app.Flags, constant.ConfigFileFlag)
	app.Action = func(cli *cli.Context) error {
		log.Root().SetHandler(log.StdoutHandler)
		cfg, err := config.Local(cli.String(constant.ConfigFileFlag.Name))
		if err != nil {
			return err
		}
		//redis.Init(cfg.Other.Redis)
		storages := make([]storage.Saver, 0, len(cfg.Storages))
		for _, s := range cfg.Storages {
			ele, err := storage.NewSaver(s.Type, s.Url)
			if err != nil {
				return err
			}
			storages = append(storages, ele)
		}

		utils.Init(cfg.Other.Env, cfg.Other.MonitorUrl)
		chainers, err := chain.Init(cfg, storages)
		if err != nil {
			return err
		}
		sysErr := make(chan error)
		c := core.New(sysErr)
		for _, ch := range chainers {
			c.AddChain(ch)
		}
		c.Start()
		return nil
	}

	if err := app.Run(os.Args); err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
}
