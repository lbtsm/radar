package constant

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/urfave/cli/v2"
	"time"
)

const (
	Near     = "near"
	Ton      = "ton"
	Ethereum = "ethereum"
)

const (
	Redis = "redis"
	Mysql = "mysql"
)

const (
	LatestBlock = "latest"
)

var (
	RetryInterval = time.Second * 5
)

var (
	ConfigFileFlag = &cli.StringFlag{
		Name:  "config",
		Usage: "JSON configuration file",
	}
	KeyPathFlag = &cli.StringFlag{
		Name:  "keystorePath",
		Usage: "Path to keystore",
	}
	LatestFlag = &cli.BoolFlag{
		Name:  "latest",
		Usage: "use latest block height",
	}
)

var (
	OnlineChaId = map[string]struct{}{}
)

const (
	KeyOfLatestBlock = "chain_%s_latest_block"
)

const (
	ReqInterval = int64(2)
)

var (
	ZeroAddress = common.HexToAddress("0x0000000000000000000000000000000000000000")
)
