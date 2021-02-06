package oracle

import (
	"context"
	"time"

	"github.com/fox-one/dirtoracle/core"
	"github.com/fox-one/dirtoracle/worker"
	"github.com/fox-one/mixin-sdk-go"
	"github.com/fox-one/pkg/logger"
	"github.com/patrickmn/go-cache"
	"github.com/shopspring/decimal"
	"golang.org/x/sync/errgroup"
)

type (
	Config struct {
		MaxInterval          time.Duration   `json:"max_interval"`
		PriceChangeThreshold decimal.Decimal `json:"price_change_threshold"`
	}

	Oracle struct {
		config    *Config
		feeds     []*core.FeedConfig
		markets   core.MarketStore
		feeders   core.FeederStore
		client    *mixin.Client
		system    *core.System
		me        *core.Member
		cache     *cache.Cache
		proposals chan *core.PriceProposal
	}
)

func New(
	client *mixin.Client,
	markets core.MarketStore,
	feeders core.FeederStore,
	feeds []*core.FeedConfig,
	system *core.System,
	config *Config,
) worker.Worker {
	m := &Oracle{
		config:    config,
		client:    client,
		feeds:     feeds,
		markets:   markets,
		feeders:   feeders,
		system:    system,
		me:        system.Me(),
		cache:     cache.New(time.Minute*15, time.Minute),
		proposals: make(chan *core.PriceProposal),
	}

	return m
}

func (m *Oracle) Run(ctx context.Context) error {
	var g errgroup.Group

	g.Go(func() error {
		return m.loopBlaze(ctx)
	})

	g.Go(func() error {
		return m.loopPrice(ctx)
	})

	g.Go(func() error {
		return m.loopProposals(ctx)
	})

	return g.Wait()
}

func (m *Oracle) loopPrice(ctx context.Context) error {
	log := logger.FromContext(ctx).WithField("worker", "oracle")
	ctx = logger.WithContext(ctx, log)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case <-time.After(time.Second):
			for _, feed := range m.feeds {
				log := log.WithField("symbol", feed.Symbol)
				ticker, err := m.markets.AggregateTickers(ctx, feed.ID)
				if err != nil {
					log.WithError(err).Errorln("AggregateTickers failed")
					continue
				}

				m.proposals <- ticker.ExportProposal()
			}
		}
	}
}
