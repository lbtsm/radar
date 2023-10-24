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
	BackUpFlag = &cli.BoolFlag{
		Name:  "backup",
		Usage: "BackUp flag",
		Value: false,
	}
)

const (
	FlagOfLatestBlock    = "latest_blockNumber_%s" // chainId
	FlagOfCurrentBlock   = "currentHandler_number_%s"
	FlagOfBackUpProgress = "backup_event_progress_%s"
	FlagOfBackUpEvent    = "backup_event_%s"
	FlagOfBackUpStop     = "backup_stop_%s"
	FlagOfAddEvent       = "extra_event_%s"
)

type BackUpEvent struct {
	Event   []string `json:"event"`
	Address []string `json:"address"`
	ChainId string   `json:"chain_id"`
}
