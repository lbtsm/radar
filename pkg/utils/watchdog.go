package utils

import (
	"context"
	"fmt"
	"sync"
	"time"
)

var (
	dogLock = sync.RWMutex{}
	set     = make(map[string]int64)
)

func AddProgress(cId string) {
	dogLock.Lock()
	set[cId] = time.Now().Unix()
	dogLock.Unlock()
}

func init() {
	go func() {
		for {
			time.Sleep(time.Minute)
			for cId, past := range set {
				if time.Now().Unix()-past >= 300 { // five minute no change，alarm
					Alarm(context.Background(), fmt.Sprintf("cId(%s) work progress not change in five minute", cId))
				}
			}
		}
	}()
}
