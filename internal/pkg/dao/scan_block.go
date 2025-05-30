package dao

type ScanBlock struct {
	Id      int64  `gorm:"column:id;primary_key;AUTO_INCREMENT;NOT NULL" gorm:"id" json:"id"`
	ChainId string `gorm:"column:chain_id;default:NULL" gorm:"chain_id" json:"chain_id"`
	Number  string `gorm:"column:number" gorm:"number" json:"number"`
}

func (b *ScanBlock) TableName() string {
	return "scan_block"
}
