package binance

import (
	"context"
	"time"

	"github.com/fox-one/pkg/logger"
)

const (
	pairsKey = "pairs"

	PairStatusPreTrading   = "PRE_TRADING"
	PairStatusTrading      = "TRADING"
	PairStatusPostTrading  = "POST_TRADING"
	PairStatusEndOfDay     = "END_OF_DAY"
	PairStatusHalt         = "HALT"
	PairStatusAuctionMatch = "AUCTION_MATCH"
	PairStatusBreak        = "BREAK"
)

type (
	PairStatus string

	Pair struct {
		Symbol     string     `json:"symbol,omitempty"`
		Status     PairStatus `json:"status,omitempty"`
		BaseAsset  string     `json:"baseAsset,omitempty"`
		QuoteAsset string     `json:"quoteAsset,omitempty"`
	}
)

func (b *binanceEx) getPairs(ctx context.Context) ([]*Pair, error) {
	if pairs, ok := b.cache.Get(pairsKey); ok {
		return pairs.([]*Pair), nil
	}

	log := logger.FromContext(ctx)
	resp, err := Request(ctx).Get("/exchangeInfo")
	if err != nil {
		log.WithError(err).Errorln("GET /exchangeInfo")
		return nil, err
	}

	var info struct {
		Pairs []*Pair `json:"symbols,omitempty"`
	}
	if err := UnmarshalResponse(resp, &info); err != nil {
		log.WithError(err).Errorln("getPairs.UnmarshalResponse")
		return nil, err
	}

	b.cache.Set(pairsKey, info.Pairs, time.Minute*10)
	return info.Pairs, nil
}

func (b *binanceEx) supported(ctx context.Context, symbol string) (bool, error) {
	pairs, err := b.getPairs(ctx)
	if err != nil {
		return false, err
	}

	for _, pair := range pairs {
		if pair.Symbol == symbol {
			return pair.Status == PairStatusTrading, nil
		}
	}
	return false, nil
}
