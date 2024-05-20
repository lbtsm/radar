package ethereum

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/mapprotocol/filter/internal/pkg/constant"
	"github.com/mapprotocol/filter/internal/pkg/dao"
	"github.com/pkg/errors"

	"github.com/mapprotocol/filter/pkg/utils"
)

func (c *Chain) sync() error {
	var currentBlock = c.cfg.StartBlock
	local, err := c.bs.TryLoadLatestBlock()
	if err != nil {
		return err
	}
	c.log.Info("sync start", "config", currentBlock, "local", local)
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
						c.log.Error("Save latest height failed", "storage", s.Type(), "err", err)
					}
				}
			}

			if currentBlock.Uint64() == 0 {
				currentBlock = big.NewInt(0).SetUint64(latestBlock)
				time.Sleep(constant.RetryInterval)
				continue
			}

			if latestBlock-currentBlock.Uint64() < c.cfg.BlockConfirmations.Uint64() {
				c.log.Debug("Block not ready, will retry", "currentBlock", currentBlock, "latest", latestBlock)
				time.Sleep(constant.RetryInterval)
				continue
			}
			err = c.mosHandler(currentBlock)
			if err != nil && !errors.Is(err, types.ErrInvalidSig) {
				c.log.Error("Failed to get events for block", "block", currentBlock, "err", err)
				utils.Alarm(context.Background(), fmt.Sprintf("filter failed, chain=%s, err is %s", c.cfg.Name, err.Error()))
				time.Sleep(constant.RetryInterval)
				continue
			}

			err = c.bs.StoreBlock(currentBlock)
			if err != nil {
				c.log.Error("Failed to write latest block to blockStore", "block", currentBlock, "err", err)
			}

			c.currentProgress = currentBlock.Int64()
			currentBlock.Add(currentBlock, big.NewInt(1))
			if latestBlock-currentBlock.Uint64() <= c.cfg.BlockConfirmations.Uint64() {
				time.Sleep(constant.RetryInterval)
			}
		}
	}
}

type oldStruct struct {
	*dao.Event
	End int64 `json:"end"`
}

func (c *Chain) rangeScan(event *dao.Event, end int64) {
	start, ok := big.NewInt(0).SetString(event.BlockNumber, 10)
	if !ok {
		return
	}
	if start.Int64() >= end {
		c.log.Info("Find a event of appoint blockNumber, but block gather than current block",
			"appoint", event.BlockNumber, "current", end)
		return
	}
	c.log.Info("Find a event of appoint blockNumber, begin start scan", "appoint", event.BlockNumber, "current", end)
	// todo store redis
	sold := &oldStruct{
		Event: event,
		End:   end,
	}
	data, _ := json.Marshal(sold)
	filename := fmt.Sprintf("%s-%s-old.json", event.ChainId, event.Topic)
	err := c.bs.CustomStore(filename, data)
	if err != nil {
		c.log.Error("Find a event of appoint blockNumber, but store local filed", "format", event.Format, "topic", event.Topic)
		return
	}
	topics := make([]common.Hash, 0, 1)
	topics = append(topics, common.HexToHash(event.Topic))
	for i := start.Int64(); i < end; i += 20 {
		// querying for logs
		logs, err := c.conn.Client().FilterLogs(context.Background(), ethereum.FilterQuery{
			FromBlock: big.NewInt(i),
			ToBlock:   big.NewInt(i + 20),
			Addresses: []common.Address{common.HexToAddress(event.Address)},
			Topics:    [][]common.Hash{topics},
		})
		if err != nil {
			continue
		}
		if len(logs) == 0 {
			continue
		}
		for _, l := range logs {
			ele := l
			err = c.insert(&ele, event)
			if err != nil {
				c.log.Error("RangeScan insert failed", "hash", l.TxHash, "logIndex", l.Index, "err", err)
				continue
			}
		}
		time.Sleep(time.Millisecond * 3)
	}
	c.log.Info("Range scan finish", "appoint", event.BlockNumber, "current", end)
	_ = c.bs.DelFile(filename)
}

func (c *Chain) mosHandler(latestBlock *big.Int) error {
	query := c.BuildQuery(latestBlock, latestBlock)
	logs, err := c.conn.Client().FilterLogs(context.Background(), query)
	if err != nil {
		return fmt.Errorf("unable to Filter Logs: %w", err)
	}
	if len(logs) == 0 {
		return nil
	}

	for _, l := range logs {
		ele := l
		idx := c.match(&ele)
		if idx == -1 {
			c.log.Debug("ignore log, because topic or address not match", "blockNumber", l.BlockNumber, "logTopic", l.Topics, "address", l.Address)
			continue
		}
		event := c.events[idx]
		err = c.insert(&ele, event)
		if err != nil {
			c.log.Error("insert failed", "hash", l.TxHash, "logIndex", l.Index, "err", err)
			continue
		}
	}

	return nil
}

func (c *Chain) insert(l *types.Log, event *dao.Event) error {
	var (
		topic     string
		toChainId uint64
		cid, _    = strconv.ParseInt(c.cfg.Id, 10, 64)
	)
	header, err := c.conn.Client().HeaderByNumber(context.Background(), big.NewInt(0).SetUint64(l.BlockNumber))
	if err != nil && strings.Index(err.Error(), "server returned non-empty transaction list but block header indicates no transactions") == -1 {
		return err
	}
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
			ProjectId:       event.ProjectId,
			EventId:         event.Id,
			TxHash:          l.TxHash.String(),
			ContractAddress: l.Address.String(),
			Topic:           topic,
			BlockNumber:     l.BlockNumber,
			LogIndex:        l.Index,
			TxIndex:         l.TxIndex,
			BlockHash:       l.BlockHash.Hex(),
			LogData:         common.Bytes2Hex(l.Data),
			TxTimestamp:     header.Time,
		})
		if err != nil {
			c.log.Error("insert failed", "hash", l.TxHash, "logIndex", l.Index, "err", err)
			continue
		}
		c.log.Info("insert success", "blockNumber", l.BlockNumber, "hash", l.TxHash, "logIndex", l.Index, "txIndex", l.TxIndex)
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

func (c *Chain) match(l *types.Log) int {
	for idx, d := range c.events {
		if l.Address.String() != d.Address {
			continue
		}
		if l.Topics[0].Hex() != d.Topic {
			continue
		}
		return idx
	}
	return -1
}
