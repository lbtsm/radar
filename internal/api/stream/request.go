package stream

type AddProjectReq struct {
	Name string `json:"name" validate:"required"`
	Desc string `json:"desc"`
}

type GetProjectReq struct {
	Id   int64  `json:"id" form:"id"`
	Name string `json:"name" form:"name" validate:"required"`
}

type AddEventReq struct {
	ProjectId   int64  `json:"project_id"`
	ChainId     int64  `json:"chain_id"`
	Format      string `json:"format" validate:"required"`
	Address     string `json:"address" validate:"required"`
	BlockNumber string `json:"block_number"`
}

type GetEventReq struct {
	Id        int64  `json:"id" form:"id"`
	ProjectId int64  `json:"project_id" form:"project_id"`
	Format    string `json:"format" form:"format"`
	Topic     string `json:"topic" form:"topic"`
}

type DelEventReq struct {
	Id int64 `json:"id" validate:"required"`
}

type EventListReq struct {
	Id        int64  `json:"id" form:"id"`
	ProjectId int64  `json:"project_id" form:"project_id"`
	Format    string `json:"format" form:"format"`
	Topic     string `json:"topic" form:"topic"`
	Offset    int64  `json:"offset" form:"offset"`
	Limit     int64  `json:"limit" form:"limit"`
}

type MosListReq struct {
	Id          int64  `json:"id" form:"id"`
	ProjectId   int64  `json:"project_id" form:"project_id"`
	ChainId     int64  `json:"chain_id" form:"chain_id"`
	Topic       string `json:"topic" form:"topic"`
	Format      string `json:"format" form:"format"`
	TxHash      string `json:"tx_hash" form:"tx_hash"`
	BlockNumber uint64 `json:"block_number" form:"block_number"`
	Limit       int    `json:"limit" form:"limit"`
}

type GetBlockReq struct {
	ChainId int64 `json:"chain_id" form:"chain_id"`
}
