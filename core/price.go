package core

import (
	"context"

	"github.com/fox-one/dirtoracle/pkg/blst"
	"github.com/shopspring/decimal"
)

type (
	PriceData struct {
		Timestamp int64           `json:"t,omitempty"`
		AssetID   string          `json:"a,omitempty"`
		Price     decimal.Decimal `gorm:"TYPE:DECIMAL(16,8);" json:"p,omitempty"`
		Signature *CosiSignature  `gorm:"TYPE:TEXT;" json:"s,omitempty"`
	}

	PriceProposal struct {
		PriceData

		Signatures map[int64]*blst.Signature `json:"sigs,omitempty"`
	}

	PriceDataStore interface {
		SavePriceData(ctx context.Context, p *PriceData) error
		LatestPriceData(ctx context.Context, asset string) (*PriceData, error)
	}
)
