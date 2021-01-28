package coinbase

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
)

const (
	Endpoint          = "https://api.pro.coinbase.com"
	WebsocketEndpoint = "wss://ws-feed.pro.coinbase.com"
)

var httpClient = resty.New().
	SetHeader("Content-Type", "application/json").
	SetHostURL(Endpoint).
	SetTimeout(10 * time.Second)

type Error struct {
	Code int    `json:"code,omitempty"`
	Msg  string `json:"message,omitempty"`
}

func (err *Error) Error() string {
	return fmt.Sprintf("[%d] %s", err.Code, err.Msg)
}

func Request(ctx context.Context) *resty.Request {
	return httpClient.R().SetContext(ctx)
}

func DecodeResponse(resp *resty.Response) ([]byte, error) {
	var respErr = Error{
		Code: resp.StatusCode(),
	}
	if err := json.Unmarshal(resp.Body(), &respErr); err != nil {
		if resp.IsError() {
			return nil, &Error{
				Code: resp.StatusCode(),
				Msg:  resp.Status(),
			}
		}

		return nil, err
	}

	if len(respErr.Msg) > 0 {
		return nil, &respErr
	}

	return json.RawMessage(resp.Body()), nil
}

func UnmarshalResponse(resp *resty.Response, v interface{}) error {
	data, err := DecodeResponse(resp)
	if err != nil {
		return err
	}

	if v != nil {
		return json.Unmarshal(data, v)
	}

	return nil
}
