package ethereum

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/mapprotocol/filter/internal/pkg/constant"
	"github.com/mapprotocol/filter/internal/pkg/dao"
	"github.com/pkg/errors"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/mapprotocol/filter/pkg/utils"
)

func (c *Chain) sync() error {
	var currentBlock = c.cfg.StartBlock
	local, err := c.bs.TryLoadLatestBlock()
	if err != nil {
		return err
	}
	if local.Cmp(currentBlock) == 1 {
		currentBlock = local
	}
	savedBN := uint64(0)
	for {
		select {
		case <-c.stop:
			return errors.New("polling terminated")
		default:
			latestBlock, err := c.conn.LatestBlock()
			if err != nil {
				c.log.Error("Unable to get latest block", "block", currentBlock, "err", err)
				time.Sleep(constant.RetryInterval)
				continue
			}

			if latestBlock != savedBN {
				savedBN = latestBlock
				for _, s := range c.storages {
					err = s.LatestBlockNumber(c.cfg.Id, latestBlock)
					if err != nil {
						c.log.Error("Mos latest block height failed", "storage", s.Type(), "err", err)
					}
				}
			}

			if latestBlock-currentBlock.Uint64() < c.cfg.BlockConfirmations.Uint64() {
				c.log.Debug("Block not ready, will retry", "currentBlock", currentBlock, "latest", latestBlock)
				time.Sleep(constant.RetryInterval)
				continue
			}
			err = c.mosHandler(currentBlock)
			if err != nil && !errors.Is(err, types.ErrInvalidSig) {
				c.log.Error("Failed to get events for block", "block", currentBlock, "err", err)
				utils.Alarm(context.Background(), fmt.Sprintf("mos failed, chain=%s, err is %s", c.cfg.Name, err.Error()))
				time.Sleep(constant.RetryInterval)
				continue
			}

			err = c.bs.StoreBlock(currentBlock)
			if err != nil {
				c.log.Error("Failed to write latest block to blockStore", "block", currentBlock, "err", err)
			}

			currentBlock.Add(currentBlock, big.NewInt(1))
			if latestBlock-currentBlock.Uint64() <= c.cfg.BlockConfirmations.Uint64() {
				time.Sleep(constant.RetryInterval)
			}
		}
	}
}

func (c *Chain) renewEvent() error {
	for {
		select {
		case <-c.stop:
			return errors.New("renewEvent polling terminated")
		default:
			for _, s := range c.storages {
				if s.Type() != constant.Mysql {
					continue
				}
				events, err := s.GetEvent(c.eventId)
				if err != nil {
					return errors.Wrap(err, fmt.Sprintf("%s get events failed", s.Type()))
				}
				for _, e := range events {
					tmp := e
					c.events = append(c.events, tmp)
					c.eventId = tmp.Id
					c.log.Info("Add new event", "event", e.Format, "topic", e.Topic, "project", e.ProjectId)
					// todo old block
				}
			}
		}
	}
}

func (c *Chain) mosHandler(latestBlock *big.Int) error {
	cid, _ := strconv.ParseInt(c.cfg.Id, 10, 64)
	query := c.BuildQuery(latestBlock, latestBlock)
	// querying for logs
	logs, err := c.conn.Client().FilterLogs(context.Background(), query)
	if err != nil {
		return fmt.Errorf("unable to Filter Logs: %w", err)
	}
	if len(logs) == 0 {
		return nil
	}
	header, err := c.conn.Client().HeaderByNumber(context.Background(), latestBlock)
	if err != nil && strings.Index(err.Error(), "server returned non-empty transaction list but block header indicates no transactions") == -1 {
		return err
	}
	for _, l := range logs {
		if !exist(l.Address, c.cfg.Mcs) {
			c.log.Debug("ignore log, because address not match", "blockNumber", l.BlockNumber, "logAddress", l.Address)
			continue
		}
		idx := existTopic(l.Topics[0], c.events)
		if idx == -1 {
			c.log.Debug("ignore log, because topic not match", "blockNumber", l.BlockNumber, "logTopic", l.Topics[0])
			continue
		}

		var (
			topic     string
			toChainId uint64
		)
		for idx, t := range l.Topics {
			topic += t.Hex()
			if idx != len(l.Topics)-1 {
				topic += ","
			}
			if idx == len(l.Topics)-1 {
				tmp, ok := big.NewInt(0).SetString(strings.TrimPrefix(t.Hex(), "0x"), 16)
				if ok {
					toChainId = tmp.Uint64()
				}
			}
		}
		for _, s := range c.storages {
			err = s.Mos(toChainId, &dao.Mos{
				ChainId:         cid,
				ProjectId:       c.events[idx].ProjectId,
				EventId:         c.events[idx].Id,
				TxHash:          l.TxHash.String(),
				ContractAddress: l.Address.String(),
				Topic:           topic,
				BlockNumber:     l.BlockNumber,
				LogIndex:        l.Index,
				LogData:         common.Bytes2Hex(l.Data),
				TxTimestamp:     header.Time,
			})
			if err != nil {
				c.log.Error("insert failed", "hash", l.TxHash, "logIndex", l.Index, "err", err)
				continue
			}
			c.log.Info("insert success", "blockNumber", l.BlockNumber, "hash", l.TxHash, "logIndex", l.Index)
		}
	}

	return nil
}

func (c *Chain) BuildQuery(startBlock *big.Int, endBlock *big.Int) ethereum.FilterQuery {
	query := ethereum.FilterQuery{
		FromBlock: startBlock,
		ToBlock:   endBlock,
	}
	return query
}

func exist(target common.Address, dst []common.Address) bool {
	for _, d := range dst {
		if target == d {
			return true
		}
	}
	return false
}

func existTopic(target common.Hash, dst []*dao.Event) int {
	for idx, d := range dst {
		if target.Hex() == d.Topic {
			return idx
		}
	}
	return -1
}
