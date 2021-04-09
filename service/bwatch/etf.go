package bwatch

import (
	"github.com/fox-one/dirtoracle/pkg/number"
	"github.com/shopspring/decimal"
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
