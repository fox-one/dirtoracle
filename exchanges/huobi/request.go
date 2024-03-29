package huobi

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
)

const (
	Endpoint = "https://api.huobi.pro/"
)

var httpClient = resty.New().
	SetHeader("Content-Type", "application/json").
	SetBaseURL(Endpoint).
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
	if resp.IsError() {
		return nil, &Error{
			Code: resp.StatusCode(),
			Msg:  resp.Status(),
		}
	}

	var e Error
	if err := json.Unmarshal(resp.Body(), &e); err == nil && e.Code > 0 {
		return nil, &e
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
