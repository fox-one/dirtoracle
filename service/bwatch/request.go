package bwatch

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fox-one/pkg/foxerr"
	"github.com/go-resty/resty/v2"
)

const bwatchHost = "https://etf-api.b.watch"

var client = resty.New().SetHostURL(bwatchHost)

func request(ctx context.Context) *resty.Request {
	return client.R().SetContext(ctx)
}

func decodeResponse(r *resty.Response, v interface{}) error {
	var resp struct {
		foxerr.Error
		Data json.RawMessage `json:"data,omitempty"`
	}

	if err := json.Unmarshal(r.Body(), &resp); err != nil {
		return fmt.Errorf("[%d] %s", r.StatusCode(), r.Status())
	}

	if resp.Code > 0 {
		return &resp.Error
	}

	if v != nil {
		return json.Unmarshal(resp.Data, v)
	}

	return nil
}
