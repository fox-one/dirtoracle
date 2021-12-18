package ftx

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
)

const (
	Endpoint = "https://ftx.com/api"
)

var httpClient = resty.New().
	SetHeader("Content-Type", "application/json").
	SetBaseURL(Endpoint).
	SetTimeout(10 * time.Second)

type Error struct {
	Success bool   `json:"success"`
	Message string `json:"error,omitempty"`
}

func (err *Error) Error() string {
	return err.Message
}

func Request(ctx context.Context) *resty.Request {
	return httpClient.R().SetContext(ctx)
}

func DecodeResponse(resp *resty.Response) ([]byte, error) {
	if resp.IsError() {
		return nil, &Error{
			Message: fmt.Sprintf("[%d] %s", resp.StatusCode(), resp.Status()),
		}
	}
	var result struct {
		Error
		Body json.RawMessage `json:"result"`
	}
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, err
	} else if !result.Success {
		return nil, &result.Error
	}
	return result.Body, nil
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
