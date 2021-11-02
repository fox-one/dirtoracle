package huobi

import (
	"context"
	"strings"
	"time"

	"github.com/fox-one/dirtoracle/pkg/route"
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
		Status     PairState `json:"state,omitempty"`
		BaseAsset  string    `json:"base-currency,omitempty"`
		QuoteAsset string    `json:"quote-currency,omitempty"`
	}

	Pairs []*Pair
)

func (pair Pair) IsOnline() bool {
	return pair.Status == PairStateOnline
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

func (exch *huobiEx) getPairs(ctx context.Context) ([]*route.Pair, error) {
	if pairs, ok := exch.cache.Get(pairsKey); ok {
		return pairs.([]*route.Pair), nil
	}

	log := logger.FromContext(ctx)
	uri := "/v1/common/symbols"
	resp, err := Request(ctx).Get(uri)
	if err != nil {
		log.WithError(err).Errorln("GET", uri)
		return nil, err
	}

	var body struct {
		Pairs Pairs `json:"data"`
	}
	if err := UnmarshalResponse(resp, &body); err != nil {
		log.WithError(err).Errorln("getPairs.UnmarshalResponse")
		return nil, err
	}

	for _, pair := range body.Pairs {
		pair.BaseAsset = strings.ToUpper(pair.BaseAsset)
		pair.QuoteAsset = strings.ToUpper(pair.QuoteAsset)
	}

	exported := body.Pairs.export()
	exch.cache.Set(pairsKey, exported, time.Minute*10)
	return exported, nil
}
