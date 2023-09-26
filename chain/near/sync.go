package near

import (
	"context"
	"encoding/json"
	"errors"
	rds "github.com/go-redis/redis/v8"
	"github.com/mapprotocol/filter/constant"
	"github.com/mapprotocol/filter/pkg/mysql"
	"github.com/mapprotocol/filter/pkg/redis"
	"strconv"
	"strings"
	"time"
)

func (c *Chain) sync() error {
	cid, _ := strconv.ParseInt(c.cfg.Id, 10, 64)
	for {
		select {
		case <-c.stop:
			return errors.New("polling terminated")
		default:
			ctx := context.Background()
			cmd := redis.GetClient().RPop(ctx, redis.ListKey)
			result, err := cmd.Result()
			if err != nil && !errors.Is(err, rds.Nil) {
				c.log.Error("Unable to get latest block", "err", err)
				time.Sleep(constant.RetryInterval)
				continue
			}
			if result == "" {
				time.Sleep(constant.RetryInterval)
				continue
			}

			data := StreamerMessage{}
			err = json.Unmarshal([]byte(result), &data)
			if err != nil {
				c.log.Error("json marshal failed", "err", err, "data", result)
				time.Sleep(constant.RetryInterval)
				continue
			}
			idx := 0
			for _, shard := range data.Shards {
				for _, outcome := range shard.ReceiptExecutionOutcomes {
					idx++
					if c.Idx(outcome.ExecutionOutcome.Outcome.ExecutorID) == -1 {
						continue
					}
					if len(outcome.ExecutionOutcome.Outcome.Logs) == 0 {
						continue
					}
					match := false
					for _, ls := range outcome.ExecutionOutcome.Outcome.Logs {
						match = c.match(ls)
						if match {
							break
						}
					}
					if !match {
						c.log.Info("Event Not Match", "log", outcome.ExecutionOutcome.Outcome.Logs)
						continue
					}

					sData, _ := json.Marshal(outcome)
					txHash := redis.GetClient().Get(context.Background(), outcome.ExecutionOutcome.ID.String())
					c.log.Info("Event found", "log", outcome.ExecutionOutcome.Outcome.Logs, "contract", outcome.ExecutionOutcome.Outcome.ExecutorID)
					_, err = mysql.GetDb().Exec("INSERT INTO mos_event (chain_id, tx_hash, contract_address, topic, block_number, log_index, log_data, tx_timestamp) "+
						"VALUES (?, ?, ?, ?, ?, ?, ?, ?)", cid, txHash, outcome.ExecutionOutcome.Outcome.ExecutorID, "",
						data.Block.Header.Height, idx, string(sData), data.Block.Header.Timestamp)
					if err != nil {
						if strings.Index(err.Error(), "Duplicate") != -1 {
							c.log.Info("log inserted", "blockNumber", data.Block.Header.Height, "hash", txHash, "logIndex", idx)
							continue
						}
						c.log.Error("insert failed", "hash", txHash, "logIndex", idx, "err", err)
						continue
					}
					c.log.Info("insert success", "blockNumber", data.Block.Header.Height, "hash", txHash, "logIndex", idx)
					time.Sleep(time.Millisecond * 50)
				}
			}
		}
	}
}

func (c *Chain) Idx(contract string) int {
	ret := -1
	for idx, addr := range c.cfg.Mcs {
		if addr == contract {
			ret = idx
			break
		}
	}

	return ret
}

func (c *Chain) match(log string) bool {
	for _, e := range c.cfg.Events {
		if strings.HasPrefix(log, e) {
			return true
		}
	}

	return false
}
