package core

import (
	"context"
	"fmt"

	"github.com/fox-one/dirtoracle/pkg/blst"
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
	"github.com/shopspring/decimal"
)

type (
	Asset struct {
		AssetID string `json:"asset_id,omitempty"`
		Symbol  string `json:"symbol,omitempty"`
	}

	FeedConfig struct {
		Asset
		Sources []string `json:"sources,omitempty"`
	}

	PriceData struct {
		Timestamp int64           `json:"t,omitempty"`
		AssetID   string          `json:"a,omitempty"`
		Price     decimal.Decimal `json:"p,omitempty"`
		Signature *CosiSignature  `json:"s,omitempty"`
	}

	PriceProposal struct {
		PriceData

		Signatures map[int64]*blst.Signature `json:"sigs,omitempty"`
	}

	Feeder struct {
		gorm.Model
		AssetID   string         `json:"asset_id,omitempty"`
		Threshold uint8          `json:"threshold,omitempty"`
		Opponents pq.StringArray `sql:"type:TEXT" json:"opponents,omitempty"`
	}

	FeederStore interface {
		SaveFeeder(ctx context.Context, f *Feeder) error
		AllFeeders(ctx context.Context) ([]*Feeder, error)
		FindFeeders(ctx context.Context, assetID string) ([]*Feeder, error)
	}
)

func (p PriceData) Payload() []byte {
	return []byte(fmt.Sprintf("%d:%s:%v", p.Timestamp, p.AssetID, p.Price))
}
