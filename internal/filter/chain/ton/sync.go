package ton

import (
	"context"
	"encoding/base64"
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/mapprotocol/filter/internal/pkg/dao"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/mapprotocol/filter/pkg/blockstore"
	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/liteclient"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/ton"
)

func (c *Chain) sync() error {
	cid, _ := strconv.ParseInt(c.cfg.Id, 10, 64)
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

	//c.log.Info("Fetching and checking proofs since config init block, it may take near a minute...")
	//master, err := api.CurrentMasterchainInfo(context.Background()) // we fetch block just to trigger chain proof check
	//if err != nil {
	//	log.Fatalln("get masterchain info err: ", err.Error())
	//	return err
	//}
	//c.log.Info("Master proof checks are completed successfully, now communication is 100% safe!", "master", master.Shard)

	//// address on which we are accepting payments
	//treasuryAddress := address.MustParseAddr("EQCD39VS5jcptHL8vMjEXrzGaRcCVYto7HUn4bpAOg8xqB2N")
	//
	//acc, err := api.GetAccount(context.Background(), master, treasuryAddress)
	//if err != nil {
	//	log.Fatalln("get masterchain info err: ", err.Error())
	//	return err
	//}
	//c.log.Info("Get account", "acc", acc.LastTxLT)

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

			c.log.Info("------------- ", "LT", lastProcessedLT)
			go api.SubscribeOnTransactions(context.Background(), treasuryAddress, 47839075000001, transactions)

			for {
				select {
				case <-tmp:
					return
				default:
					c.log.Info("Waiting for transfers...", "addr", addr)
					for t := range transactions {
						txHash := base64.StdEncoding.EncodeToString(t.Hash)
						c.log.Info("Get transaction", "addr", addr, "txHash", txHash)
						if t.IO.Out == nil {
							c.log.Info("In transaction", "txHash", txHash)
							continue
						}
						msgs, err := t.IO.Out.ToSlice()
						if err != nil {
							c.log.Error("Tx ToSlice failed", "addr", addr, "txHash", txHash)
							break
						}

						for _, msg := range msgs {
							if msg.MsgType != tlb.MsgTypeExternalOut {
								continue
							}
							externalOut := msg.AsExternalOut()
							log.Println("src: ", externalOut.SrcAddr)
							log.Println("dst: ", externalOut.DstAddr)

							idx := c.match(externalOut.DstAddr.String())
							if idx == -1 {
								c.log.Info("Ignore this tx, because topic not match", "txHash", txHash,
									"topic", externalOut.DstAddr.String())
								continue
							}

							data, err := msg.AsExternalOut().Payload().MarshalJSON()
							if err != nil {
								c.log.Error("TxHash marshal failed", "txHash", txHash, "err", err)
							}
							for _, s := range c.storages {
								err = s.Mos(0, &dao.Mos{
									ChainId:         cid,
									ProjectId:       1,
									TxHash:          txHash,
									ContractAddress: treasuryAddress.String(),
									Topic:           externalOut.DstAddr.String(),
									LogData:         common.Bytes2Hex(data),
								})
								if err != nil {
									c.log.Error("Insert failed", "hash", txHash, "err", err)
									continue
								}
								c.log.Info("Insert success", "hash", txHash)
							}
						}
					}

					c.log.Error("Something went wrong, transaction listening unexpectedly finished")
				}
			}
		}(ele)
	}

	for {
		select {
		case <-c.stop:
			for _, ele := range sig {
				close(ele)
			}
			return errors.New("polling terminated")
		default:
			c.log.Info("is running")
			time.Sleep(time.Minute)
		}
	}
}

func (c *Chain) match(target string) int {
	for idx, v := range c.cfg.Event {
		if strings.HasSuffix(target, v) {
			return idx
		}
	}
	return -1
}
