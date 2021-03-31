package core

import (
	"fmt"

	"github.com/fox-one/dirtoracle/pkg/mtg"
	"github.com/gofrs/uuid"
	"github.com/shopspring/decimal"
)

type (
	PriceData struct {
		Timestamp int64           `json:"t,omitempty"`
		AssetID   string          `json:"a,omitempty"`
		Price     decimal.Decimal `json:"p,omitempty"`
		Signature *CosiSignature  `json:"s,omitempty"`
	}
)

func (p PriceData) Payload() []byte {
	return []byte(fmt.Sprintf("%d:%s:%v", p.Timestamp, p.AssetID, p.Price))
}

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
