package constant

import (
	"github.com/urfave/cli/v2"
	"time"
)

const (
	Near     = "near"
	Ethereum = "ethereum"
)

const (
	Redis = "redis"
	Mysql = "mysql"
)

var (
	RetryInterval = time.Second * 5
)

var (
	ConfigFileFlag = &cli.StringFlag{
		Name:  "config",
		Usage: "JSON configuration file",
	}
)

var (
	OnlineChaId = map[string]struct{}{}
)

const (
	KeyOfLatestBlock = "chain_%s_latest_block"
)
