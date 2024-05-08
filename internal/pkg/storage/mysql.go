package storage

import (
	"github.com/ethereum/go-ethereum/log"
	_ "github.com/go-sql-driver/mysql"
	"github.com/mapprotocol/filter/internal/pkg/constant"
	"github.com/mapprotocol/filter/internal/pkg/dao"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	glog "log"
	"os"
	"strings"
)

type Mysql struct {
	dsn string
	db  *gorm.DB
}

func NewMysql(dsn string) (*Mysql, error) {
	m := &Mysql{dsn: dsn}
	err := m.init()
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (m *Mysql) init() error {
	mdb, err := gorm.Open(mysql.Open(m.dsn), &gorm.Config{
		Logger: logger.New(
			glog.New(os.Stdout, "\r\n", glog.LstdFlags),
			logger.Config{
				LogLevel: logger.Info,
				Colorful: false,
			},
		),
	})

	if err != nil {
		return err
	}
	m.db = mdb
	return nil
}

func (m *Mysql) Type() string {
	return constant.Mysql
}

func (m *Mysql) Mos(toChainId uint64, event *dao.Mos) error {
	err := m.db.Create(event).Error
	if err != nil {
		if strings.Index(err.Error(), "Duplicate") != -1 {
			log.Info("log is inserted", "blockNumber", event.BlockNumber, "hash", event.TxHash,
				"logIndex", event.LogIndex)
		}
		return err
	}
	return nil
}

func (m *Mysql) LatestBlockNumber(chainId string, latest uint64) error {
	return nil
}

func (m *Mysql) GetEvent(id int64) ([]*dao.Event, error) {
	ret := make([]*dao.Event, 0)
	err := m.db.Model(&dao.Event{}).Where("id > ?", id).Find(&ret).Error
	return ret, err
}
