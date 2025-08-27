package butter

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/mapprotocol/filter/internal/pkg/client"
	"github.com/mapprotocol/filter/internal/pkg/stream"
)

const (
	UrlOfMessageInBridgeData = "/messageInBridgeData"
)

const (
	SuccessCode = 0
)

var defaultButter = New()

type Butter struct {
}

func New() *Butter {
	return &Butter{}
}

func (b *Butter) MessageInBridgeData(domain, query string) ([]byte, error) {
	fmt.Println("ExecMessageInBridgeData uri ", fmt.Sprintf("%s%s?%s", domain, UrlOfMessageInBridgeData, query))
	return client.JsonGet(fmt.Sprintf("%s%s?%s", domain, UrlOfMessageInBridgeData, query))
}

func MessageInBridgeData(domain, query string) ([]byte, error) {
	return defaultButter.MessageInBridgeData(domain, query)
}

func RequestBridgeData(domain, txHash string, request *stream.BridgeDataRequest) (*stream.Bridge, error) {
	affiliate := ""
	if len(request.Affiliate) > 0 {
		affiliateParams := make([]string, len(request.Affiliate))
		for i, a := range request.Affiliate {
			if i == 0 {
				affiliateParams[i] = "affiliate=" + url.QueryEscape(a)
			} else {
				affiliateParams[i] = "&affiliate=" + url.QueryEscape(a)
			}
		}
		affiliate = strings.Join(affiliateParams, "")
	}
	params := fmt.Sprintf(
		"fromChainId=%s&caller=%s&toChainId=%s&amount=%s&tokenInAddress=%s&tokenOutAddress=%s&minAmountOut=%s&receiver=%s&entrance=%s&%s",
		url.QueryEscape(request.FromChainID),
		url.QueryEscape(request.Caller),
		url.QueryEscape(request.ToChainID),
		url.QueryEscape(request.Amount),
		url.QueryEscape(request.TokenInAddress),
		url.QueryEscape(request.TokenOutAddress),
		url.QueryEscape(request.MinAmountOut),
		url.QueryEscape(request.Receiver),
		url.QueryEscape(request.Entrance),
		affiliate,
	)

	fullURL := domain + UrlOfMessageInBridgeData + "?" + params
	fmt.Printf("request butter bridge data url: %s\n", fullURL)

	resp, err := http.Get(fullURL)
	if err != nil {
		return nil, &stream.ExternalRequestError{
			URL:  fullURL,
			Msg:  err.Error(),
			Code: "",
		}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &stream.ExternalRequestError{
			URL:  fullURL,
			Msg:  err.Error(),
			Code: "",
		}
	}

	var ret stream.BridgeDataResponse
	if err := json.Unmarshal(body, &ret); err != nil {
		return nil, &stream.ExternalRequestError{
			URL:  fullURL,
			Msg:  err.Error(),
			Code: "",
		}
	}

	if ret.StatusCode != SuccessCode {
		return nil, &stream.ExternalRequestError{
			URL:  fullURL,
			Msg:  ret.Message,
			Code: fmt.Sprintf("%d", ret.StatusCode),
		}
	}

	if ret.Errno != SuccessCode {
		return nil, &stream.ExternalRequestError{
			URL:  fullURL,
			Msg:  ret.Message,
			Code: fmt.Sprintf("%d", ret.Errno),
		}
	}

	if ret.Data == nil || len(ret.Data) == 0 {
		return nil, &stream.ExternalRequestError{
			URL:    fullURL,
			Msg:    ret.Message,
			Code:   fmt.Sprintf("%d", ret.Errno),
			Detail: fmt.Errorf("transaction data not found"),
		}
	}

	fmt.Printf("request butter bridge back data: %+v\n", ret)
	return &ret.Data[0], nil
}
