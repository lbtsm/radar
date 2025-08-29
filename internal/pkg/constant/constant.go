package constant

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/urfave/cli/v2"
)

const (
	Near     = "near"
	Ton      = "ton"
	Xrp      = "xrp"
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
	BackUpFlag = &cli.BoolFlag{
		Name:  "back",
		Usage: "is back up program",
	}
)

var (
	OnlineChaId = map[string]struct{}{}
)

const (
	KeyOfLatestBlock = "chain_%s_latest_block"
	KeyOfScanBlock   = "chain_%s_scan_block"
)

const (
	ReqInterval = int64(2)
)

var (
	ZeroAddress = common.HexToAddress("0x0000000000000000000000000000000000000000")
)

const (
	TopicMessageOut = "0x2aaebc938a5ab70e98644b0c6a8472fe02b7edece7e6e6dc71959dc34e109dfc"
	TopicMessageIn  = "0xf01fbdd2fdbc5c2f201d087d588789d600e38fe56427e813d9dced2cdb25bcac"
)
