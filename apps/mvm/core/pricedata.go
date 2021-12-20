package core

import (
	"errors"
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
		En256Signature *En256CosiSignature `json:"es,omitempty"`
	}
)

func (p PriceData) Payload() []byte {
	return []byte(fmt.Sprintf("%d:%s:%v", p.Timestamp, p.AssetID, p.Price))
}

func (p PriceData) PayloadV1() ([]byte, error) {
	asset, err := uuid.FromString(p.AssetID)
	if err != nil {
		return nil, err
	}

	return mtg.Encode(p.Timestamp, asset, p.Price)
}

func (p *PriceData) MarshalBinary() (data []byte, err error) {
	asset, err := uuid.FromString(p.AssetID)
	if err != nil {
		return nil, err
	}
	if p.En256Signature != nil {
		return mtg.Encode(p.Timestamp, asset, p.Price, p.En256Signature)
	} else if p.Signature != nil {
		return mtg.Encode(p.Timestamp, asset, p.Price, p.Signature)
	}
	return nil, errors.New("empty signature")
}

func (p *PriceData) UnmarshalBinary(data []byte) error {
	var (
		timestamp int64
		price     decimal.Decimal
		asset     uuid.UUID
	)
	left, err := mtg.Scan(data, &timestamp, &asset, &price)
	if err != nil {
		return err
	}

	if len(left) == 52 {
		var signature CosiSignature
		if _, err := mtg.Scan(left, &signature); err != nil {
			return err
		}
		p.Signature = &signature
	} else if len(left) == 37 || len(left) == 68 {
		var signature En256CosiSignature
		if _, err := mtg.Scan(left, &signature); err != nil {
			return err
		}
		p.En256Signature = &signature
	} else {
		return errors.New("unknown signature")
	}
	p.Timestamp = timestamp
	p.AssetID = asset.String()
	p.Price = price

	return nil
}
