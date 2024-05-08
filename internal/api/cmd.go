package api

import (
	"github.com/ethereum/go-ethereum/log"
	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"github.com/mapprotocol/filter/internal/api/config"
	"github.com/mapprotocol/filter/internal/pkg/constant"
	"github.com/urfave/cli/v2"
)

var Command = &cli.Command{
	Name:  "api",
	Flags: []cli.Flag{constant.ConfigFileFlag},
	Action: func(cli *cli.Context) error {
		log.Root().SetHandler(log.StdoutHandler)
		cfg, err := config.Local(cli.String(constant.ConfigFileFlag.Name))
		if err != nil {
			return err
		}

		g := gin.Default()
		initMiddleware(g)
		err = initController(g, cfg.Dsn)
		if err != nil {
			log.Error("init failed", "err", err)
			return err
		}

		endless.ListenAndServe(cfg.Listen, g)
		return nil
	},
}
