package ftx

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
		Name          string  `json:"name,omitempty"`
		Enabled       bool    `json:"enabled"`
		Restricted    bool    `json:"restricted"`
		PostOnly      bool    `json:"postOnly"`
		Last          float64 `json:"last"`
		Bid           float64 `json:"bid"`
		Ask           float64 `json:"ask"`
		Price         float64 `json:"price"`
		Type          string  `json:"type,omitempty"`
		BaseCurrency  string  `json:"baseCurrency,omitempty"`
		QuoteCurrency string  `json:"quoteCurrency,omitempty"`
	}

	Pairs []*Pair
)

func (pair Pair) IsOnline() bool {
	return pair.Enabled && !pair.PostOnly
}

func (pairs Pairs) export() []*route.Pair {
	items := make([]*route.Pair, 0, len(pairs))
	for _, pair := range pairs {
		if !pair.IsOnline() {
			continue
		}
		items = append(items, &route.Pair{
			Symbol: pair.Name,
			Base:   pair.BaseCurrency,
			Quote:  pair.QuoteCurrency,
		})
	}
	return items
}

func (exch *ftxEx) getPairs(ctx context.Context) ([]*route.Pair, error) {
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
