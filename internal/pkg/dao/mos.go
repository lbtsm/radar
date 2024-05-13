package dao

import "time"

type Mos struct {
	Id              int64     `gorm:"column:id;primary_key;AUTO_INCREMENT;NOT NULL" gorm:"id" json:"id"`
	ChainId         int64     `gorm:"column:chain_id;default:NULL" gorm:"chain_id" json:"chain_id"`
	ProjectId       int64     `gorm:"column:project_id;default:NULL" gorm:"project_id" json:"project_id"`
	EventId         int64     `gorm:"column:event_id;default:NULL" gorm:"event_id" json:"event_id"`
	TxHash          string    `gorm:"column:tx_hash;default:NULL" gorm:"tx_hash" json:"tx_hash"`
	ContractAddress string    `gorm:"column:contract_address;default:NULL" gorm:"contract_address" json:"contract_address"`
	BlockHash       string    `gorm:"column:block_hash;default:NULL" gorm:"block_hash" json:"block_hash"`
	Topic           string    `gorm:"column:topic;default:NULL" gorm:"topic" json:"topic"`
	BlockNumber     uint64    `gorm:"column:block_number;default:NULL" gorm:"block_number" json:"block_number"`
	LogIndex        uint      `gorm:"column:log_index;default:NULL" gorm:"log_index" json:"log_index"`
	TxIndex         uint      `gorm:"column:tx_index;default:NULL" gorm:"tx_index" json:"tx_index"`
	LogData         string    `gorm:"column:log_data" gorm:"log_data" json:"log_data"`
	TxTimestamp     uint64    `gorm:"column:tx_timestamp;default:NULL" gorm:"tx_timestamp" json:"tx_timestamp"`
	CreatedAt       time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP;NOT NULL" gorm:"created_at" json:"created_at"`
}

func (m *Mos) TableName() string {
	return "mos"
}
