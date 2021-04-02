package bitstamp

import (
	"context"
	"time"

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
		Name      string        `json:"name,omitempty"`
		UrlSymbol string        `json:"url_symbol,omitempty"`
		Trading   TradingStatus `json:"trading,omitempty"`
	}
)

func (b *bitstampEx) getPairs(ctx context.Context) ([]*Pair, error) {
	if pairs, ok := b.cache.Get(pairsKey); ok {
		return pairs.([]*Pair), nil
	}

	log := logger.FromContext(ctx)
	resp, err := Request(ctx).Get("/trading-pairs-info/")
	if err != nil {
		log.WithError(err).Errorln("GET /trading-pairs-info/")
		return nil, err
	}

	var pairs []*Pair
	if err := UnmarshalResponse(resp, &pairs); err != nil {
		log.WithError(err).Errorln("getPairs.UnmarshalResponse")
		return nil, err
	}

	b.cache.Set(pairsKey, pairs, time.Hour)
	return pairs, nil
}

func (b *bitstampEx) supported(ctx context.Context, symbol string) (bool, error) {
	pairs, err := b.getPairs(ctx)
	if err != nil {
		return false, err
	}

	for _, pair := range pairs {
		if pair.UrlSymbol == symbol {
			return pair.Trading == TradingStatusEnabled, nil
		}
	}
	return false, nil
}
