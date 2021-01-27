package market

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/RussellLuo/timingwheel"
	"github.com/fox-one/dirtoracle/core"
	"github.com/fox-one/dirtoracle/core/exchange"
	"github.com/fox-one/dirtoracle/worker"
	"github.com/fox-one/pkg/logger"
	"github.com/sirupsen/logrus"
)

type Markets struct {
	channels  []*core.Channel
	markets   core.MarketStore
	exchanges map[string]exchange.Interface

	tw *timingwheel.TimingWheel
}

func New(
	markets core.MarketStore,
	channels []*core.Channel,
	exchanges map[string]exchange.Interface,
) worker.Worker {
	m := &Markets{
		channels:  channels,
		markets:   markets,
		exchanges: exchanges,
		tw:        timingwheel.NewTimingWheel(time.Second, 5),
	}

	for _, c := range channels {
		if c.Asset == nil || c.Asset.ID != c.AssetID {
			panic(fmt.Errorf("invalid asset channel: %s", c.AssetID))
		}
		for _, e := range c.Exchanges {
			m.mustExchange(e)
		}
	}

	return m
}

// MarketHandler Implementation
func (m *Markets) OnTicker(ctx context.Context, asset *core.Asset, ticker *core.Ticker) error {
	logger.FromContext(ctx).WithFields(logrus.Fields{
		"price":      ticker.Price,
		"volume":     ticker.VolumeUSD,
		"updated_at": ticker.UpdatedAt,
	}).Debugln("OnTicker")
	return m.markets.SaveTicker(ctx, ticker)
}

func (m *Markets) Run(ctx context.Context) error {
	log := logger.FromContext(ctx).WithField("worker", "markets")
	ctx = logger.WithContext(ctx, log)

	m.tw.Start()
	defer m.tw.Stop()

	for _, c := range m.channels {
		c := c
		for _, e := range c.Exchanges {
			e := e
			log := log.WithFields(logrus.Fields{
				"symbol": c.Asset.Symbol,
				"ex":     e,
			})
			ctx = logger.WithContext(ctx, log)
			go m.subscribe(ctx, c.Asset, e)

		}
	}

	<-ctx.Done()
	return ctx.Err()
}

func (m *Markets) mustExchange(id string) exchange.Interface {
	ex, ok := m.exchanges[id]
	if !ok {
		panic(fmt.Errorf("exchange with id %s not found", id))
	}

	return ex
}

func (m *Markets) subscribe(ctx context.Context, asset *core.Asset, e string) error {
	ex := m.mustExchange(e)
	if err := ex.Subscribe(ctx, asset, m); err != nil {
		logger.FromContext(ctx).WithError(err).Errorln("subscribe markets")

		if errors.Is(err, context.Canceled) {
			return err
		}
	}

	m.tw.AfterFunc(time.Second, func() {
		m.subscribe(ctx, asset, e)
	})
	return nil
}
