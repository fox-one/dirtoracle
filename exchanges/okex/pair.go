package okex

import (
	"context"
	"time"

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
		Symbol string    `json:"instId,omitempty"`
		Basse  string    `json:"baseCcy,omitempty"`
		Quote  string    `json:"quoteCcy,omitempty"`
		Type   string    `json:"instType,omitempty"`
		State  PairState `json:"state,omitempty"`
	}
)

func (b *okexEx) getPairs(ctx context.Context) ([]*Pair, error) {
	if pairs, ok := b.cache.Get(pairsKey); ok {
		return pairs.([]*Pair), nil
	}

	log := logger.FromContext(ctx)
	uri := "/v5/public/instruments"
	resp, err := Request(ctx).SetQueryParam("instType", "SPOT").Get(uri)
	if err != nil {
		log.WithError(err).Errorln("GET", uri)
		return nil, err
	}

	var pairs []*Pair
	if err := UnmarshalResponse(resp, &pairs); err != nil {
		log.WithError(err).Errorln("getPairs.UnmarshalResponse")
		return nil, err
	}

	b.cache.Set(pairsKey, pairs, time.Minute*10)
	return pairs, nil
}

func (b *okexEx) supported(ctx context.Context, symbol string) (bool, error) {
	pairs, err := b.getPairs(ctx)
	if err != nil {
		return false, err
	}

	for _, pair := range pairs {
		if pair.Symbol == symbol {
			return pair.State == PairStateLive, nil
		}
	}
	return false, nil
}
