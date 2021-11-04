package binance

import (
	"context"
	"time"

	"github.com/fox-one/dirtoracle/pkg/route"
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
		Symbol               string     `json:"symbol,omitempty"`
		Status               PairStatus `json:"status,omitempty"`
		BaseAsset            string     `json:"baseAsset,omitempty"`
		QuoteAsset           string     `json:"quoteAsset,omitempty"`
		IsSpotTradingAllowed bool       `json:"isSpotTradingAllowed,omitempty"`
	}

	Pairs []*Pair
)

func (pair Pair) IsOnline() bool {
	return pair.Status == PairStatusTrading && pair.IsSpotTradingAllowed
}

func (pairs Pairs) export() []*route.Pair {
	items := make([]*route.Pair, 0, len(pairs))
	for _, pair := range pairs {
		if !pair.IsOnline() {
			continue
		}
		items = append(items, &route.Pair{
			Symbol: pair.Symbol,
			Base:   pair.BaseAsset,
			Quote:  pair.QuoteAsset,
		})
	}
	return items
}

func (exch *binanceEx) getPairs(ctx context.Context) ([]*route.Pair, error) {
	if pairs, ok := exch.cache.Get(pairsKey); ok {
		return pairs.([]*route.Pair), nil
	}

	log := logger.FromContext(ctx)
	resp, err := Request(ctx).Get("/exchangeInfo")
	if err != nil {
		log.WithError(err).Errorln("GET /exchangeInfo")
		return nil, err
	}

	var info struct {
		Pairs Pairs `json:"symbols,omitempty"`
	}
	if err := UnmarshalResponse(resp, &info); err != nil {
		log.WithError(err).Errorln("getPairs.UnmarshalResponse")
		return nil, err
	}

	pairs := info.Pairs.export()
	exch.cache.Set(pairsKey, pairs, time.Minute*10)
	return pairs, nil
}
