package okex

import (
	"context"
	"time"

	"github.com/fox-one/dirtoracle/pkg/route"
	"github.com/fox-one/pkg/logger"
)

const (
	pairsKey = "pairs"

	PairStateLive    = "live"
	PairStateSuspend = "suspend"
	PairStatePreopen = "preopen"
)

type (
	PairState string

	Pair struct {
		Symbol        string    `json:"instId,omitempty"`
		BaseCurrency  string    `json:"baseCcy,omitempty"`
		QuoteCurrency string    `json:"quoteCcy,omitempty"`
		Type          string    `json:"instType,omitempty"`
		Status        PairState `json:"state,omitempty"`
	}

	Pairs []*Pair
)

func (pair Pair) IsOnline() bool {
	return pair.Status == PairStateLive
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

func (exch *okexEx) getPairs(ctx context.Context) ([]*route.Pair, error) {
	if pairs, ok := exch.cache.Get(pairsKey); ok {
		return pairs.([]*route.Pair), nil
	}

	log := logger.FromContext(ctx)
	uri := "/v5/public/instruments"
	resp, err := Request(ctx).SetQueryParam("instType", "SPOT").Get(uri)
	if err != nil {
		log.WithError(err).Errorln("GET", uri)
		return nil, err
	}

	var pairs Pairs
	if err := UnmarshalResponse(resp, &pairs); err != nil {
		log.WithError(err).Errorln("getPairs.UnmarshalResponse")
		return nil, err
	}

	exported := pairs.export()
	exch.cache.Set(pairsKey, exported, time.Minute*10)
	return exported, nil
}
