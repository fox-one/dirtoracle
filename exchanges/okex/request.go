package okex

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
)

const (
	Endpoint = "https://aws.okex.com/api"
)

var httpClient = resty.New().
	SetHeader("Content-Type", "application/json").
	SetBaseURL(Endpoint).
	SetTimeout(10 * time.Second)

type Error struct {
	Code string `json:"code,omitempty"`
	Msg  string `json:"msg,omitempty"`
}

func (err *Error) Error() string {
	return fmt.Sprintf("[%s] %s", err.Code, err.Msg)
}

func Request(ctx context.Context) *resty.Request {
	return httpClient.R().SetContext(ctx)
}

func DecodeResponse(resp *resty.Response) ([]byte, error) {
	if resp.IsError() {
		return nil, &Error{
			Code: fmt.Sprint(resp.StatusCode()),
			Msg:  resp.Status(),
		}
	}

	var body struct {
		Error
		Body json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(resp.Body(), &body); err != nil {
		return nil, err
	} else if body.Code != "0" {
		return nil, &body.Error
	}
	return body.Body, nil
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
