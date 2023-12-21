package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/mapprotocol/filter/internal/constant"
	"github.com/mapprotocol/filter/internal/dao"
	"strconv"
)

var (
	KeyOfMapMessenger   = "messenger_%d_%d" // messenger_sourceChainId_toChainId
	KeyOfOtherMessenger = "messenger_%d"    // messenger_sourceChainId_toChainId
)

type Redis struct {
	dsn         string
	redisClient *redis.Client
}

func newRds(dsn string) (*Redis, error) {
	m := &Redis{dsn: dsn}
	err := m.init()
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (r *Redis) init() error {
	opt, err := redis.ParseURL(r.dsn)
	if err != nil {
		return err
	}
	rdb := redis.NewClient(opt)
	if err != nil {
		return err
	}
	r.redisClient = rdb
	return nil
}

func (r *Redis) GetType() string {
	return constant.Redis
}

func (r *Redis) Storage(toChainId uint64, event *dao.MosEvent) error {
	var key string
	if event.ChainId == 22776 || event.ChainId == 212 || event.ChainId == 213 {
		if _, ok := constant.OnlineChaId[strconv.FormatUint(toChainId, 10)]; !ok {
			return nil
		}
		key = fmt.Sprintf(KeyOfMapMessenger, event.ChainId, toChainId)
	} else {
		key = fmt.Sprintf(KeyOfOtherMessenger, event.ChainId)
	}
	data, _ := json.Marshal(event)
	_, err := r.redisClient.RPush(context.Background(), key, data).Result()
	if err != nil {
		return err
	}
	return nil
}
