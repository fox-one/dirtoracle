package bwatch

import (
	"context"

	"github.com/fox-one/dirtoracle/core"
	"github.com/fox-one/dirtoracle/pkg/number"
	"github.com/shopspring/decimal"
)

const (
	etfsKey = "etfs"
)

type (
	Etf struct {
		AssetID               string          `json:"asset_id,omitempty"`
		Name                  string          `json:"name,omitempty"`
		Symbol                string          `json:"symbol,omitempty"`
		Logo                  string          `json:"logo,omitempty"`
		Price                 decimal.Decimal `json:"price"`
		MaxSupply             decimal.Decimal `json:"max_supply"`
		CirculatingSupply     decimal.Decimal `json:"circulating_supply"`
		MaxSubscriptionAmount decimal.Decimal `json:"max_subscription_amount"`
		MinSubscriptionAmount decimal.Decimal `json:"min_subscription_amount"`
		MaxRedemptionAmount   decimal.Decimal `json:"max_redemption_amount"`
		MinRedemptionAmount   decimal.Decimal `json:"min_redemption_amount"`
		SubscriptionFeeRate   decimal.Decimal `json:"subscription_fee_rate"`
		SubscriptionFee       decimal.Decimal `json:"subscription_fee"`
		RedemptionFeeRate     decimal.Decimal `json:"redemption_fee_rate"`
		RedemptionFee         decimal.Decimal `json:"redemption_fee"`
		Assets                number.Values   `json:"assets"`
	}
)

func (b *bwatchService) getETFs(ctx context.Context) ([]*Etf, error) {
	if v, ok := b.cache.Get(etfsKey); ok {
		return v.([]*Etf), nil
	}

	r, err := request(ctx).Get("/api/etfs")
	if err != nil {
		return nil, err
	}

	var body struct {
		Etfs []*Etf `json:"etfs"`
	}
	if err = decodeResponse(r, &body); err != nil {
		return nil, err
	}

	b.cache.SetDefault(etfsKey, body.Etfs)
	return body.Etfs, nil
}

func (b *bwatchService) getETF(ctx context.Context, a *core.Asset) (*Etf, error) {
	etfs, err := b.getETFs(ctx)
	if err != nil {
		return nil, err
	}
	for _, etf := range etfs {
		if etf.AssetID == a.AssetID {
			return etf, nil
		}
	}
	return nil, nil
}
