package exinswap

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-resty/resty/v2"
)

const (
	Endpoint = "https://app.exinswap.com/api/v1"
)

var httpClient = resty.New().
	SetHeader("Content-Type", "application/json").
	SetBaseURL(Endpoint).
	SetTimeout(10 * time.Second)

func Request(ctx context.Context) *resty.Request {
	return httpClient.R().SetContext(ctx)
}

func DecodeResponse(resp *resty.Response) ([]byte, error) {
	var respErr Error
	if err := json.Unmarshal(resp.Body(), &respErr); err != nil {
		if resp.IsError() {
			return nil, &Error{
				Code: resp.StatusCode(),
				Msg:  resp.Status(),
			}
		}

		return nil, err
	}

	if respErr.Code > 0 {
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
