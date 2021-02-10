package bitstamp

import (
	"context"
	"math/rand"
	"strings"
	"time"

	"github.com/fox-one/dirtoracle/core"
	"github.com/fox-one/dirtoracle/core/exchange"
	"github.com/fox-one/pkg/logger"
)

const (
	exchangeName = "bitstamp"
)

type bitstampEx struct{}

func New() exchange.Interface {
	return &bitstampEx{}
}

func (c *bitstampEx) Name() string {
	return exchangeName
}

func (c *bitstampEx) Subscribe(ctx context.Context, a *core.Asset, h exchange.Handler) error {
	log := logger.FromContext(ctx)
	log.Info("start")
	defer log.Info("quit")

	var (
		sleepDur   = time.Duration(rand.Int63n(int64(time.Second * 5)))
		pairSymbol = c.pairSymbol(c.assetSymbol(a.Symbol))
	)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case <-time.After(sleepDur):
			t, err := readTicker(ctx, pairSymbol)
			if err != nil {
				log.WithError(err).Errorln("readTicker failed")
				sleepDur = 5 * time.Second
				continue
			}
			t.AssetID = a.AssetID
			if err := h.OnTicker(ctx, t); err != nil {
				log.WithError(err).Errorln("OnTicker failed")
				sleepDur = time.Second
				continue
			}
			sleepDur = 10 * time.Second
		}
	}
}

func (*bitstampEx) assetSymbol(symbol string) string {
	return symbol
}

func (*bitstampEx) pairSymbol(symbol string) string {
	return strings.ToLower(symbol) + "usd"
}
