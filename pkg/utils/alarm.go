package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/log"
)

var (
	prefix, hooksUrl = "", ""
	lock             = sync.RWMutex{}
	monitor          = make(map[string]int64)
)

func Init(env, hooks string) {
	prefix = env
	hooksUrl = hooks
}

func Alarm(ctx context.Context, msg string) {
	if hooksUrl == "" {
		log.Info("hooks is empty")
		return
	}
	lock.RLock()
	if v, ok := monitor[msg]; ok {
		if time.Now().Unix()-v < 300 { // ignore same alarm in five minute
			return
		}
	}
	lock.RUnlock()
	fmt.Println("send alarm ulock")
	lock.Lock()
	monitor[msg] = time.Now().Unix()
	lock.Unlock()
	fmt.Println("send alarm ulock")
	body, err := json.Marshal(map[string]interface{}{
		"text": fmt.Sprintf("%s %s", prefix, msg),
	})
	if err != nil {
		return
	}

	req, err := http.NewRequestWithContext(ctx, "POST", hooksUrl, io.NopCloser(bytes.NewReader(body)))
	if err != nil {
		return
	}
	req.Header.Set("Content-type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Warn("read resp failed", "err", err)
		return
	}
	fmt.Println("send alarm message", "resp", string(data))
}
