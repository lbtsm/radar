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
	Format      string `json:"format" validate:"required"`
	BlockNumber string `json:"block_number"`
}

type GetEventReq struct {
	Id        int64  `json:"id" form:"id"`
	ProjectId int64  `json:"project_id"`
	Format    string `json:"format"`
	Topic     string `json:"topic"`
}

type DelEventReq struct {
	Id int64 `json:"id" validate:"required"`
}

type EventListReq struct {
	Id        int64  `json:"id"`
	ProjectId int64  `json:"project_id"`
	Format    string `json:"format"`
	Topic     string `json:"topic"`
	Offset    int64  `json:"offset"`
	Limit     int64  `json:"limit"`
}
