package bitstamp

import (
	"context"
	"strings"
	"time"

	"github.com/fox-one/dirtoracle/pkg/route"
	"github.com/fox-one/pkg/logger"
)

const (
	pairsKey = "pairs"

	TradingStatusEnabled  = "Enabled"
	TradingStatusDisabled = "Disabled"
)

type (
	TradingStatus string

	Pair struct {
		Name          string        `json:"name,omitempty"`
		UrlSymbol     string        `json:"url_symbol,omitempty"`
		Trading       TradingStatus `json:"trading,omitempty"`
		BaseCurrency  string        `json:"base_currency,omitempty"`
		QuoteCurrency string        `json:"quote_currency,omitempty"`
	}

	Pairs []*Pair
)

func (pair Pair) IsOnline() bool {
	return pair.Trading == TradingStatusEnabled
}

func (pairs Pairs) export() []*route.Pair {
	items := make([]*route.Pair, 0, len(pairs))
	for _, pair := range pairs {
		if !pair.IsOnline() {
			continue
		}
		items = append(items, &route.Pair{
			Symbol: pair.UrlSymbol,
			Base:   pair.BaseCurrency,
			Quote:  pair.QuoteCurrency,
		})
	}
	return items
}

func (exch *bitstampEx) getPairs(ctx context.Context) ([]*route.Pair, error) {
	if pairs, ok := exch.cache.Get(pairsKey); ok {
		return pairs.([]*route.Pair), nil
	}

	log := logger.FromContext(ctx)
	resp, err := Request(ctx).Get("/trading-pairs-info/")
	if err != nil {
		log.WithError(err).Errorln("GET /trading-pairs-info/")
		return nil, err
	}

	var pairs Pairs
	if err := UnmarshalResponse(resp, &pairs); err != nil {
		log.WithError(err).Errorln("getPairs.UnmarshalResponse")
		return nil, err
	}

	for _, pair := range pairs {
		if arr := strings.Split(pair.Name, "/"); len(arr) == 2 {
			pair.BaseCurrency = arr[0]
			pair.QuoteCurrency = arr[1]
		}
	}

	exported := pairs.export()
	exch.cache.Set(pairsKey, exported, time.Hour)
	return exported, nil
}
