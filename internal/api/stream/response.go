package stream

type CommonResp struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type GetProjectResp struct {
	Id      int64  `json:"id"`
	Name    string `json:"name"`
	Desc    string `json:"desc"`
	Created int64  `json:"created"`
}

type AddEventResp struct {
	ProjectId   int64  `json:"project_id"`
	Format      string `json:"format"`
	BlockNumber string `json:"block_number"`
}

type GetEventResp struct {
	Id        int64  `json:"id"`
	ProjectId int64  `json:"project_id"`
	Format    string `json:"format"`
	Topic     string `json:"topic"`
	Created   int64  `json:"created"`
}

type EventListResp struct {
	Total int64           `json:"total"`
	Page  int64           `json:"page"`
	Limit int64           `json:"limit"`
	List  []*GetEventResp `json:"list"`
}

type GetMosResp struct {
	Id              int64  `json:"id"`
	ProjectId       int64  `json:"project_id"`
	ChainId         int64  `json:"chain_id"`
	EventId         int64  `json:"event_id"`
	TxHash          string `json:"tx_hash"`
	ContractAddress string `json:"contract_address"`
	Topic           string `json:"topic"`
	BlockNumber     uint64 `json:"block_number"`
	LogIndex        uint   `json:"log_index"`
	LogData         string `json:"log_data"`
	TxTimestamp     uint64 `json:"tx_timestamp"`
}

type MosListResp struct {
	Total int64         `json:"total"`
	List  []*GetMosResp `json:"list"`
}
