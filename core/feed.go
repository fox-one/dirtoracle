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
	FeedConfig struct {
		Asset
		Sources []string `json:"sources,omitempty"`
	}

	PriceData struct {
		Timestamp int64           `json:"timestamp,omitempty"`
		AssetID   string          `json:"asset_id,omitempty"`
		Price     decimal.Decimal `json:"price,omitempty"`
		Mask      int64           `json:"mask,omitempty"`
		Signature *blst.Signature `json:"signature,omitempty"`
	}

	PriceProposal struct {
		PriceData

		Signatures map[int64]*blst.Signature `json:"signatures,omitempty"`
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
	}
)

func (p PriceData) Payload() []byte {
	return []byte(fmt.Sprintf("%d:%s:%v", p.Timestamp, p.AssetID, p.Price))
}
