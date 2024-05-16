package ethereum

import (
	"context"
	"fmt"
	"github.com/mapprotocol/filter/pkg/utils"
	"time"
)

func (c *Chain) watchdog() {
	tmp := c.currentProgress
	for {
		time.Sleep(time.Minute)
		if tmp != c.currentProgress {
			c.log.Info("watchdog work progress report", "record", tmp, "curr", c.currentProgress)
			tmp = c.currentProgress
			continue
		}
		c.log.Info("watchdog find chain work progress not change in minute, will retry conn", "record", tmp, "curr", c.currentProgress)
		utils.Alarm(context.Background(), fmt.Sprintf("cId(%s) work progress (%d) not change in one minute", c.cfg.Id, tmp))
		c.stop <- struct{}{}
		time.Sleep(time.Minute)
		c.conn.Close()
		newConn := NewConn(c.cfg.Endpoint, c.kp)
		err := newConn.Connect()
		if err != nil {
			c.log.Error("watchdog retry conn", "err", err, "endpoint", c.cfg.Endpoint)
			continue
		}
		c.conn = newConn
		c.log.Info("watchdog retry conn success, will sync", "endpoint", c.cfg.Endpoint)
		go func() {
			err := c.sync()
			if err != nil {
				c.log.Error("Polling blocks failed", "err", err)
			}
		}()
	}
}
