package constant

import (
	"github.com/urfave/cli/v2"
	"time"
)

const (
	Near     = "near"
	Ethereum = "ethereum"
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

const (
	FlagOfLatestBlock  = "latest_blockNumber_%s" // chainId
	FlagOfCurrentBlock = "currentHandler_number_%s"
)
