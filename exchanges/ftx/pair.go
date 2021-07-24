package ftx

import (
	"context"
	"time"

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
)

func (c *ftxEx) getPairs(ctx context.Context) ([]*Pair, error) {
	if pairs, ok := c.cache.Get(pairsKey); ok {
		return pairs.([]*Pair), nil
	}

	log := logger.FromContext(ctx)
	resp, err := Request(ctx).Get("/markets")
	if err != nil {
		log.WithError(err).Errorln("GET /markets")
		return nil, err
	}

	var pairs []*Pair
	if err := UnmarshalResponse(resp, &pairs); err != nil {
		log.WithError(err).Errorln("getPairs.UnmarshalResponse")
		return nil, err
	}

	c.cache.Set(pairsKey, pairs, time.Hour)
	return pairs, nil
}

func (c *ftxEx) supported(ctx context.Context, symbol string) (bool, error) {
	pairs, err := c.getPairs(ctx)
	if err != nil {
		return false, err
	}

	for _, pair := range pairs {
		if pair.Name == symbol {
			return pair.Enabled && !pair.PostOnly, nil
		}
	}
	return false, nil
}
