package coinbase

import (
	"context"
	"time"

	"github.com/fox-one/dirtoracle/pkg/route"
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

	Pairs []*Pair
)

func (pair Pair) IsOnline() bool {
	return pair.Status == PairStatusOnline && !pair.CancelOnly && !pair.PostOnly && !pair.TradingDisabled
}

func (pairs Pairs) export() []*route.Pair {
	items := make([]*route.Pair, 0, len(pairs))
	for _, pair := range pairs {
		if !pair.IsOnline() {
			continue
		}
		items = append(items, &route.Pair{
			Symbol: pair.ID,
			Base:   pair.BaseCurrency,
			Quote:  pair.QuoteCurrency,
		})
	}
	return items
}

func (exch *coinbaseEx) getPairs(ctx context.Context) ([]*route.Pair, error) {
	if pairs, ok := exch.cache.Get(pairsKey); ok {
		return pairs.([]*route.Pair), nil
	}

	log := logger.FromContext(ctx)
	resp, err := Request(ctx).Get("/products")
	if err != nil {
		log.WithError(err).Errorln("GET /products")
		return nil, err
	}

	var pairs Pairs
	if err := UnmarshalResponse(resp, &pairs); err != nil {
		log.WithError(err).Errorln("getPairs.UnmarshalResponse")
		return nil, err
	}

	exported := pairs.export()
	exch.cache.Set(pairsKey, exported, time.Hour)
	return exported, nil
}
