package dao

import "time"

type MosEvent struct {
	Id              int64     `gorm:"column:id;primary_key;AUTO_INCREMENT;NOT NULL" gorm:"id" json:"id"`
	ChainId         int64     `gorm:"column:chain_id;default:NULL" gorm:"chain_id" json:"chain_id"`
	TxHash          string    `gorm:"column:tx_hash;default:NULL" gorm:"tx_hash" json:"tx_hash"`
	ContractAddress string    `gorm:"column:contract_address;default:NULL" gorm:"contract_address" json:"contract_address"`
	Topic           string    `gorm:"column:topic;default:NULL" gorm:"topic" json:"topic"`
	BlockNumber     uint64    `gorm:"column:block_number;default:NULL" gorm:"block_number" json:"block_number"`
	LogIndex        uint      `gorm:"column:log_index;default:NULL" gorm:"log_index" json:"log_index"`
	LogData         string    `gorm:"column:log_data" gorm:"log_data" json:"log_data"`
	TxTimestamp     uint64    `gorm:"column:tx_timestamp;default:NULL" gorm:"tx_timestamp" json:"tx_timestamp"`
	CreateAt        time.Time `gorm:"column:create_at;default:CURRENT_TIMESTAMP;NOT NULL" gorm:"create_at" json:"create_at"`
}

func (m *MosEvent) TableName() string {
	return "mos_event"
}
