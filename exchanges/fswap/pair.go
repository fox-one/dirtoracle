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

func (f *fswapEx) getPairs(ctx context.Context) ([]*fswapsdk.Pair, error) {
	if pairs, ok := f.cache.Get(pairsKey); ok {
		return pairs.([]*fswapsdk.Pair), nil
	}

	log := logger.FromContext(ctx)
	pairs, err := fswapsdk.ListPairs(ctx)
	if err != nil {
		log.WithError(err).Errorln("GET /pairs")
		return nil, err
	}

	f.cache.Set(pairsKey, pairs, time.Second*10)
	return pairs, nil
}
