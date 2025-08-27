package stream

import "fmt"

type BridgeDataResponse struct {
	Errno      int      `json:"errno"`
	StatusCode int      `json:"statusCode"`
	Message    string   `json:"message"`
	Data       []Bridge `json:"data"`
}

type Bridge struct {
	Data     string `json:"data"`
	Relay    bool   `json:"relay"`
	Receiver string `json:"receiver"`
}

type BridgeDataRequest struct {
	Entrance        string   `json:"entrance,omitempty"`
	Affiliate       []string `json:"affiliate,omitempty"`
	FromChainID     string   `json:"fromChainID"`
	ToChainID       string   `json:"toChainID"`
	Amount          string   `json:"amount"`
	TokenInAddress  string   `json:"tokenInAddress"`
	TokenOutAddress string   `json:"tokenOutAddress"`
	MinAmountOut    string   `json:"minAmountOut"`
	Receiver        string   `json:"receiver"`
	Caller          string   `json:"caller"`
	EntranceId      string   `json:"entranceId"` // only sol
}

type ExternalRequestError struct {
	URL    string
	Msg    string
	Code   string
	Detail error
}

func (e *ExternalRequestError) Error() string {
	if e.Detail != nil {
		return fmt.Sprintf("ExternalRequestError: URL=%s, Message=%s, Code=%s, Detail=%v", e.URL, e.Msg, e.Code, e.Detail)
	}
	return fmt.Sprintf("ExternalRequestError: URL=%s, Message=%s, Code=%s", e.URL, e.Msg, e.Code)
}
