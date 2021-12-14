package core

import (
	"fmt"

	"github.com/fox-one/dirtoracle/pkg/mtg"
	"github.com/gofrs/uuid"
	"github.com/shopspring/decimal"
)

type (
	PriceData struct {
		Timestamp      int64               `json:"t,omitempty"`
		AssetID        string              `json:"a,omitempty"`
		Price          decimal.Decimal     `json:"p,omitempty"`
		Signature      *CosiSignature      `json:"s,omitempty"`
		En256Signature *CosiEn256Signature `json:"es,omitempty"`
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
	if p.En256Signature != nil {
		return mtg.Encode(p.Timestamp, asset, p.Price, p.Signature, p.En256Signature)
	}
	return mtg.Encode(p.Timestamp, asset, p.Price, p.Signature)
}

func (p *PriceData) UnmarshalBinary(data []byte) error {
	var (
		timestamp      int64
		price          decimal.Decimal
		signature      CosiSignature
		en256Signature CosiEn256Signature
		asset          uuid.UUID
	)
	left, err := mtg.Scan(data, &timestamp, &asset, &price, &signature)
	if err != nil {
		return err
	}
	p.Timestamp = timestamp
	p.AssetID = asset.String()
	p.Price = price
	p.Signature = &signature

	if _, err := mtg.Scan(left, &en256Signature); err == nil {
		p.En256Signature = &en256Signature
	}
	return nil
}
