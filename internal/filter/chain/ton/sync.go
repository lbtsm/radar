package ton

import (
	"context"
	"encoding/base64"
	"errors"
	"time"

	"github.com/mapprotocol/filter/pkg/blockstore"
	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/liteclient"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/ton"
)

func (c *Chain) sync() error {
	client := liteclient.NewConnectionPool()
	cfg, err := liteclient.GetConfigFromUrl(context.Background(), "https://ton.org/global.config.json")
	if err != nil {
		c.log.Error("Get config failed", "err", err.Error())
		return err
	}

	err = client.AddConnectionsFromConfig(context.Background(), cfg)
	if err != nil {
		c.log.Error("Connection failed: ", "err", err.Error())
		return err
	}
	api := ton.NewAPIClient(client, ton.ProofCheckPolicySecure).WithRetry()
	api.SetTrustedBlockFromConfig(cfg)

	sig := make([]chan struct{}, 0)
	for _, v := range c.cfg.Mcs {
		ele := v
		tmp := make(chan struct{})
		sig = append(sig, tmp)
		go func(addr string) {
			bs, err := blockstore.New(blockstore.PathPostfix, addr+"-"+c.cfg.Id)
			if err != nil {
				c.log.Error("New BlockStore failed", "addr", addr, "err", err)
				close(c.stop)
				return
			}

			treasuryAddress := address.MustParseAddr(addr)
			transactions := make(chan *tlb.Transaction)
			lastProcessedLT, err := bs.TryLoadLatestBlock()
			if err != nil {
				c.log.Error("TryLoadLatestBlock failed", "addr", addr, "err", err)
				close(c.stop)
				return
			}

			go api.SubscribeOnTransactions(context.Background(), treasuryAddress, lastProcessedLT.Uint64(), transactions)

			for {
				select {
				case <-tmp:
					return
				default:
					c.log.Info("Waiting for transfers...", "addr", addr)
					for t := range transactions {
						if t.IO.Out == nil {
							c.log.Info("In transaction", "addr", addr, "txHash", base64.StdEncoding.EncodeToString(t.Hash))
							continue
						}
						msgs, err := t.IO.Out.ToSlice()
						if err != nil {
							c.log.Error("Tx ToSlice failed", "addr", addr, "txHash", base64.StdEncoding.EncodeToString(t.Hash))
							break
						}

						for _, v := range msgs {
							if v.MsgType != tlb.MsgTypeExternalOut {
								continue
							}
							data := v.AsExternalOut().Payload()
							c.log.Error("Tx ToSlice failed", "addr", addr, "txHash", base64.StdEncoding.EncodeToString(t.Hash), "data", data)
						}
					}

					c.log.Error("something went wrong, transaction listening unexpectedly finished")
				}
			}
		}(ele)
	}

	for {
		select {
		case <-c.stop:
			for _, c := range sig {
				close(c)
			}
			return errors.New("polling terminated")
		default:
			c.log.Info("is running")
			time.Sleep(time.Minute)
		}
	}
}
