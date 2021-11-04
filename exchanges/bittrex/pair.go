package bittrex

import (
	"context"
	"time"

	"github.com/fox-one/dirtoracle/pkg/route"
	"github.com/fox-one/pkg/logger"
)

const (
	pairsKey = "pairs"

	TradingStatusOnline = "ONLINE"
)

type (
	TradingStatus string

	Pair struct {
		Symbol        string `json:"symbol,omitempty"`
		BaseCurrency  string `json:"baseCurrencySymbol,omitempty"`
		QuoteCurrency string `json:"quoteCurrencySymbol,omitempty"`
		Status        string `json:"status,omitempty"`
	}

	Pairs []*Pair
)

func (pair Pair) IsOnline() bool {
	return pair.Status == TradingStatusOnline
}

func (pairs Pairs) export() []*route.Pair {
	items := make([]*route.Pair, 0, len(pairs))
	for _, pair := range pairs {
		if !pair.IsOnline() {
			continue
		}
		items = append(items, &route.Pair{
			Symbol: pair.Symbol,
			Base:   pair.BaseCurrency,
			Quote:  pair.QuoteCurrency,
		})
	}
	return items
}

func (exch *bittrexEx) getPairs(ctx context.Context) ([]*route.Pair, error) {
	if pairs, ok := exch.cache.Get(pairsKey); ok {
		return pairs.([]*route.Pair), nil
	}

	log := logger.FromContext(ctx)
	resp, err := Request(ctx).Get("/markets")
	if err != nil {
		log.WithError(err).Errorln("GET /markets")
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
