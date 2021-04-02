package coinbase

import (
	"context"
	"time"

	"github.com/fox-one/pkg/logger"
)

const (
	pairsKey = "pairs"

	PairStatusOnline = "online"
)

type (
	PairStatus string

	Pair struct {
		ID              string     `json:"id,omitempty"`
		DisplayName     string     `json:"display_name,omitempty"`
		BaseCurrency    string     `json:"base_currency,omitempty"`
		QuoteCurrency   string     `json:"quote_currency,omitempty"`
		Status          PairStatus `json:"status,omitempty"`
		StatusMessage   string     `json:"status_message,omitempty"`
		CancelOnly      bool       `json:"cancel_only"`
		LimitOnly       bool       `json:"limit_only"`
		PostOnly        bool       `json:"post_only"`
		TradingDisabled bool       `json:"trading_disabled"`
	}
)

func (c *coinbaseEx) getPairs(ctx context.Context) ([]*Pair, error) {
	if pairs, ok := c.cache.Get(pairsKey); ok {
		return pairs.([]*Pair), nil
	}

	log := logger.FromContext(ctx)
	resp, err := Request(ctx).Get("/products")
	if err != nil {
		log.WithError(err).Errorln("GET /products")
		return nil, err
	}

	var pairs []*Pair
	if err := UnmarshalResponse(resp, &pairs); err != nil {
		log.WithError(err).Errorln("getPairs.UnmarshalResponse")
		return nil, err
	}

	c.cache.Set(pairsKey, pairs, time.Hour)
	return pairs, nil
}

func (c *coinbaseEx) supported(ctx context.Context, symbol string) (bool, error) {
	pairs, err := c.getPairs(ctx)
	if err != nil {
		return false, err
	}

	for _, pair := range pairs {
		if pair.ID == symbol {
			return pair.Status == PairStatusOnline && !pair.CancelOnly &&
				!pair.PostOnly && !pair.TradingDisabled, nil
		}
	}
	return false, nil
}
