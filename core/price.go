package core

import (
	"context"

	"github.com/fox-one/dirtoracle/pkg/blst"
	"github.com/fox-one/dirtoracle/pkg/mtg"
	"github.com/gofrs/uuid"
	"github.com/shopspring/decimal"
)

type (
	PriceData struct {
		Timestamp int64           `json:"t,omitempty"`
		AssetID   string          `sql:"SIZE:36;" json:"a,omitempty"`
		Price     decimal.Decimal `sql:"TYPE:DECIMAL(16,8);" json:"p,omitempty"`
		Signature *CosiSignature  `sql:"TYPE:TEXT;" json:"s,omitempty"`
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

func (p *PriceData) MarshalBinary() (data []byte, err error) {
	asset, err := uuid.FromString(p.AssetID)
	if err != nil {
		return nil, err
	}
	return mtg.Encode(p.Timestamp, asset, p.Price, p.Signature)
}

func (p *PriceData) UnmarshalBinary(data []byte) error {
	var (
		d     PriceData
		asset uuid.UUID
	)
	_, err := mtg.Scan(data, &d.Timestamp, &asset, &d.Price, d.Signature)
	if err != nil {
		return err
	}
	p.Timestamp = d.Timestamp
	p.AssetID = asset.String()
	p.Price = d.Price
	p.Signature = d.Signature
	return nil
}
