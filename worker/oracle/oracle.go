package oracle

import (
	"context"
	"time"

	"github.com/fox-one/dirtoracle/core"
	"github.com/fox-one/dirtoracle/worker"
	"github.com/fox-one/mixin-sdk-go"
	"github.com/patrickmn/go-cache"
	"github.com/shopspring/decimal"
	"golang.org/x/sync/errgroup"
)

var (
	maxDuration          = time.Minute
	priceChangeThreshold = decimal.New(1, -2)
)

type (
	Oracle struct {
		exchanges   []core.Exchange
		wallets     core.WalletStore
		assetz      core.AssetService
		posrvs      []core.PortfolioService
		subscribers core.SubscriberStore
		client      *mixin.Client
		system      *core.System
		cache       *cache.Cache
	}
)

func New(
	exchanges []core.Exchange,
	client *mixin.Client,
	wallets core.WalletStore,
	assetz core.AssetService,
	posrvs []core.PortfolioService,
	subscribers core.SubscriberStore,
	system *core.System,
) worker.Worker {
	m := &Oracle{
		exchanges:   exchanges,
		client:      client,
		wallets:     wallets,
		assetz:      assetz,
		posrvs:      posrvs,
		subscribers: subscribers,
		system:      system,
		cache:       cache.New(time.Minute*15, time.Minute),
	}

	return m
}

func (m *Oracle) Run(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return m.loopBlaze(ctx)
	})

	g.Go(func() error {
		return m.loopSubscribers(ctx)
	})

	g.Go(func() error {
		return m.loopTopAssetss(ctx)
	})

	return g.Wait()
}
