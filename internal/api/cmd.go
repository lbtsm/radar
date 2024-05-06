package api

import (
	"github.com/ethereum/go-ethereum/log"
	"github.com/mapprotocol/filter/internal/filter/config"
	"github.com/mapprotocol/filter/internal/pkg/constant"
	"github.com/mapprotocol/filter/pkg/utils"
	"github.com/urfave/cli/v2"
)

var Command = &cli.Command{
	Name: "api",
	Action: func(cli *cli.Context) error {
		log.Root().SetHandler(log.StdoutHandler)
		cfg, err := config.Local(cli.String(constant.ConfigFileFlag.Name))
		if err != nil {
			return err
		}

		utils.Init(cfg.Other.Env, cfg.Other.MonitorUrl)

		return nil
	},
}
