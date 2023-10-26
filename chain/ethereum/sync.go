package ethereum

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/mapprotocol/filter/pkg/mysql"
	"github.com/mapprotocol/filter/pkg/redis"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/mapprotocol/filter/constant"
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
	for {
		select {
		case <-c.stop:
			return errors.New("polling terminated")
		default:
			latestBlock, err := c.conn.Client().BlockNumber(context.Background())
			if err != nil {
				c.log.Error("Unable to get latest block", "block", currentBlock, "err", err)
				time.Sleep(constant.RetryInterval)
				continue
			}
			if !c.cfg.BackUp {
				c.saveValue(constant.FlagOfLatestBlock, latestBlock)
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
			if c.cfg.BackUp {
				c.saveValue(constant.FlagOfBackUpProgress, currentBlock.String())
				res, err := redis.GetClient().Get(context.Background(), fmt.Sprintf(constant.FlagOfBackUpStop, c.cfg.Id)).Result()
				if err == nil {
					redis.GetClient().Del(context.Background(), fmt.Sprintf(constant.FlagOfBackUpProgress, c.cfg.Id))
					c.log.Info("BackUp receive stop signal, will stop process", "res", res)
					return nil
				}
			} else {
				c.saveValue(constant.FlagOfCurrentBlock, currentBlock.String())
				c.getBackupEvent(currentBlock)
			}

			currentBlock.Add(currentBlock, big.NewInt(1))
			if latestBlock-currentBlock.Uint64() <= c.cfg.BlockConfirmations.Uint64() {
				time.Sleep(constant.RetryInterval)
			}
		}
	}
}

func (c *Chain) mosHandler(latestBlock *big.Int) error {
	cid, _ := strconv.ParseInt(c.cfg.Id, 10, 64)
	query := c.BuildQuery(latestBlock, latestBlock)
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
		if !existTopic(l.Topics[0], c.cfg.Events) {
			c.log.Debug("ignore log, because address not match", "blockNumber", l.BlockNumber, "logTopic", l.Topics[0])
			continue
		}

		var topic string
		for idx, t := range l.Topics {
			topic += t.Hex()
			if idx != len(l.Topics)-1 {
				topic += ","
			}
		}
		// save
		_, err = mysql.GetDb().Exec("INSERT INTO mos_event (chain_id, tx_hash, contract_address, topic, block_number, log_index, log_data, tx_timestamp) "+
			"VALUES (?, ?, ?, ?, ?, ?, ?, ?)", cid, l.TxHash.String(), l.Address.String(), topic, l.BlockNumber, l.Index, common.Bytes2Hex(l.Data), header.Time)
		if err != nil {
			if strings.Index(err.Error(), "Duplicate") != -1 {
				c.log.Info("log inserted", "blockNumber", l.BlockNumber, "hash", l.TxHash, "logIndex", l.Index)
				continue
			}
			c.log.Error("insert failed", "hash", l.TxHash, "logIndex", l.Index, "err", err)
			continue
		}
		c.log.Info("insert success", "blockNumber", l.BlockNumber, "hash", l.TxHash, "logIndex", l.Index)
		time.Sleep(time.Millisecond * 50)
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

func existTopic(target common.Hash, dst []constant.EventSig) bool {
	for _, d := range dst {
		if target == d.GetTopic() {
			return true
		}
	}
	return false
}

func (c *Chain) saveValue(format string, value interface{}) {
	key := fmt.Sprintf(format, c.cfg.Id)
	err := redis.GetClient().Set(context.Background(), key, value, 0).Err()
	if err != nil {
		c.log.Error("Save latestBlock to redis failed", "value", value, "err", err)
	}
}

func (c *Chain) getBackupEvent(currentBlock *big.Int) {
	res, err := redis.GetClient().Get(context.Background(), fmt.Sprintf(constant.FlagOfBackUpProgress, c.cfg.Id)).Result()
	if err != nil {
		return
	}
	c.log.Info("Get backup event progress", "backupProgress", res, "currentBlock", currentBlock)
	bac, _ := big.NewInt(0).SetString(res, 10)
	if bac.Uint64() < currentBlock.Uint64() {
		return
	}
	// add event
	backup, err := redis.GetClient().Get(context.Background(), fmt.Sprintf(constant.FlagOfBackUpEvent, c.cfg.Id)).Result()
	if err != nil {
		c.log.Error("Failed to Get backup event", "err", err)
		return
	}
	c.log.Info("Get backup event", "backupEvent", backup)
	bu := constant.BackUpEvent{}
	err = json.Unmarshal([]byte(backup), &bu)
	if err != nil {
		c.log.Error("Failed to Unmarshal", "data", backup, "err", err)
		return
	}
	c.log.Info("Get backup event before", "address", c.cfg.Mcs, "event", c.cfg.Events, "add", bu.Address, "addEvent", c.cfg.Events)
	for _, a := range bu.Address {
		c.cfg.Mcs = append(c.cfg.Mcs, common.HexToAddress(a))
	}
	for _, s := range bu.Event {
		c.cfg.Events = append(c.cfg.Events, constant.EventSig(s))
	}
	c.log.Info("Get backup event after", "address", c.cfg.Mcs, "event", c.cfg.Events)
	// insert backup stop signal
	c.saveValue(constant.FlagOfBackUpStop, 1)
	// save
	c.saveValue(constant.FlagOfAddEvent, backup)
	// del
	redis.GetClient().Del(context.Background(), fmt.Sprintf(constant.FlagOfBackUpEvent, c.cfg.Id))
}
