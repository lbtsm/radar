package xrp

import (
	"context"
	"fmt"
	"time"

	"github.com/mapprotocol/filter/pkg/utils"
)

func (c *Chain) watchdog() {
	tmp := c.currentProgress
	for {
		select {
		case <-c.dog:
			c.log.Info("watchdog receive stop signal")
			return
		default:
			time.Sleep(time.Minute)
			if tmp != c.currentProgress {
				c.log.Info("watchdog progress report", "record", tmp, "curr", c.currentProgress, "latest", c.latest)
				tmp = c.currentProgress
				continue
			}
			if tmp == c.latest {
				c.log.Info("watchdog progress report, curr = latest", "record", tmp)
				continue
			}
			c.log.Info("watchdog work progress not change in minute, will retry conn", "record", tmp, "curr", c.currentProgress)
			utils.Alarm(context.Background(), fmt.Sprintf("chain(%s) work progress (%d) not change in one minute", c.cfg.Name, tmp))
			c.log.Info("watchdog work progress not change in minute, send alarm ok")
			c.stop <- struct{}{}
			time.Sleep(time.Second)
			c.conn.Close()
			time.Sleep(time.Minute)
			for {
				c.log.Info("watchdog will retry conn ", "endpoint", c.cfg.Endpoint)
				newConn := NewConn(c.cfg.Endpoint)
				err := newConn.Connect()
				if err != nil {
					c.log.Error("watchdog retry conn", "err", err, "endpoint", c.cfg.Endpoint)
					time.Sleep(time.Second * 30)
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
				break
			}
		}
	}
}
