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
	priceChangeThreshold = decimal.New(2, -2)
)

type (
	Oracle struct {
		exchange    core.Exchange
		wallets     core.WalletStore
		subscribers core.SubscriberStore
		client      *mixin.Client
		system      *core.System
		cache       *cache.Cache
	}
)

func New(
	exchange core.Exchange,
	client *mixin.Client,
	wallets core.WalletStore,
	subscribers core.SubscriberStore,
	system *core.System,
) worker.Worker {
	m := &Oracle{
		exchange:    exchange,
		client:      client,
		wallets:     wallets,
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

	return g.Wait()
}
