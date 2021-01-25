package core

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
)

type (
	Ticker struct {
		Exchange string `json:"exchange,omitempty"`
		AssetID  string `json:"asset_id,omitempty"`

		UpdatedAt time.Time `json:"updated_at,omitempty"`
		// AskPrice 卖一参考价格，根据交易所深度算出，仅用于估算
		AskPrice decimal.Decimal `json:"ask_price,omitempty"`
		// BidPrice 买一参考价格，根据交易所深度算出，仅用于估算
		BidPrice decimal.Decimal `json:"bid_price,omitempty"`
		// LastPrice 交易所最新的成交价格
		LastPrice decimal.Decimal `json:"last_price,omitempty"`
		// Change24h 24h 价格变化比例
		Change24h decimal.Decimal `json:"change_24h,omitempty"`
	}

	MarketStore interface {
		// ticker
		SaveTicker(ctx context.Context, ticker *Ticker) error
		FindTickers(ctx context.Context, assetID string) ([]*Ticker, error)
	}
)
