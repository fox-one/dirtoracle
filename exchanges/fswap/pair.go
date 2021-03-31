package fswap

import (
	"context"
	"time"

	fswapsdk "github.com/fox-one/4swap-sdk-go"
	"github.com/fox-one/pkg/logger"
)

const (
	pairsKey = "pairs"
)

func (c *fswapEx) getPairs(ctx context.Context) ([]*fswapsdk.Pair, error) {
	if pairs, ok := c.cache.Get(pairsKey); ok {
		return pairs.([]*fswapsdk.Pair), nil
	}

	log := logger.FromContext(ctx)
	const uri = "/api/pairs"
	resp, err := fswapsdk.Request(ctx).Get(uri)
	if err != nil {
		log.WithError(err).Errorln("GET", uri)
		return nil, err
	}

	var body struct {
		Pairs []*fswapsdk.Pair `json:"pairs"`
	}
	if err := fswapsdk.UnmarshalResponse(resp, &body); err != nil {
		log.WithError(err).Errorln("getPairs.UnmarshalResponse")
		return nil, err
	}

	c.cache.Set(pairsKey, body.Pairs, time.Second*10)
	return body.Pairs, nil
}
