package bittrex

import (
	"context"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/fox-one/dirtoracle/core"
	"github.com/fox-one/dirtoracle/core/exchange"
	"github.com/fox-one/pkg/logger"
)

const (
	exchangeName = "bittrex"
)

type bittrexEx struct {
	once    sync.Once
	tickers map[string]*core.Ticker
}

func New() exchange.Interface {
	return &bittrexEx{}
}

func (c *bittrexEx) Name() string {
	return exchangeName
}

func (c *bittrexEx) Subscribe(ctx context.Context, a *core.Asset, h exchange.Handler) error {
	log := logger.FromContext(ctx)
	log.Info("start")
	defer log.Info("quit")

	c.once.Do(func() {
		go c.syncPairs(ctx)
	})

	var (
		sleepDur   = time.Duration(rand.Int63n(int64(time.Second * 5)))
		pairSymbol = c.pairSymbol(c.assetSymbol(a.Symbol))
	)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case <-time.After(sleepDur):
			t, ok := c.tickers[pairSymbol]
			if !ok {
				log.Errorln("ticker not found")
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

func (c *bittrexEx) syncPairs(ctx context.Context) {
	log := logger.FromContext(ctx).WithField("thread", "sync_pairs")

	log.Info("start")
	defer log.Info("quit")

	var (
		sleepDur = time.Millisecond
	)

	for {
		select {
		case <-ctx.Done():
			return

		case <-time.After(sleepDur):
			if tickers, err := readTickers(ctx); err != nil {
				log.WithError(err).Errorln("readTickers failed")
				sleepDur = time.Second * 2
			} else {
				c.tickers = tickers
				sleepDur = time.Second * 5
			}
		}
	}
}

func (*bittrexEx) assetSymbol(symbol string) string {
	return symbol
}

func (*bittrexEx) pairSymbol(symbol string) string {
	return strings.ToLower(symbol) + "-USD"
}
