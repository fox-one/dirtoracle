package huobi

import (
	"context"
	"time"

	"github.com/fox-one/pkg/logger"
)

const (
	pairsKey = "pairs"

	PairStateOnline  = "online"
	PairStateOffline = "offline"
	PairStateSuspend = "suspend"
)

type (
	PairState string

	Pair struct {
		Symbol     string    `json:"symbol,omitempty"`
		State      PairState `json:"state,omitempty"`
		BaseAsset  string    `json:"base-currency,omitempty"`
		QuoteAsset string    `json:"quote-currency,omitempty"`
	}
)

func (b *huobiEx) getPairs(ctx context.Context) ([]*Pair, error) {
	if pairs, ok := b.cache.Get(pairsKey); ok {
		return pairs.([]*Pair), nil
	}

	log := logger.FromContext(ctx)
	uri := "/v1/common/symbols"
	resp, err := Request(ctx).Get(uri)
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

func (b *huobiEx) supported(ctx context.Context, symbol string) (bool, error) {
	pairs, err := b.getPairs(ctx)
	if err != nil {
		return false, err
	}

	for _, pair := range pairs {
		if pair.Symbol == symbol {
			return pair.State == PairStateOnline, nil
		}
	}
	return false, nil
}
