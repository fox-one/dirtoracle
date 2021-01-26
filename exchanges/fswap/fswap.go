package fswap

import (
	"context"
	"errors"
	"sync"
	"time"

	fswapsdk "github.com/fox-one/4swap-sdk-go"
	"github.com/fox-one/dirtoracle/core"
	"github.com/fox-one/dirtoracle/core/exchange"
	"github.com/fox-one/pkg/logger"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

const (
	pusdAsset    = "31d2ea9c-95eb-3355-b65b-ba096853bc18"
	exchangeName = "4swap"
)

type (
	PairResp struct {
		Pairs     []*fswapsdk.Pair `json:"pairs"`
		Timestamp int64            `json:"ts"`
	}

	fswapEx struct {
		once  sync.Once
		pairs *PairResp
	}
)

func init() {
	fswapsdk.UseEndpoint("https://f1-mtgswap-api.firesbox.com")
}

func New() exchange.Interface {
	return &fswapEx{}
}

func (b *fswapEx) Name() string {
	return exchangeName
}

func (f *fswapEx) Subscribe(ctx context.Context, asset *core.Asset, handler exchange.MarketHandler) error {
	log := logger.FromContext(ctx).WithFields(logrus.Fields{
		"ex":     "4swap",
		"worker": "subscribe",
	})
	ctx = logger.WithContext(ctx, log)

	f.once.Do(func() {
		f.syncPairs(ctx)
	})

	log = log.WithField("asset", asset.Symbol)
	log.Info("start")
	defer log.Info("quit")

	var (
		sleepDur = time.Millisecond
	)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case <-time.After(sleepDur):
			sleepDur = time.Second
			ticker, err := f.readTicker(ctx, asset)
			if err != nil {
				log.WithError(err).Errorln("readTicker failed")
				continue
			}
			if err := handler.OnTicker(ctx, asset, ticker); err != nil {
				log.WithError(err).Errorln("onTicker failed")
				continue
			}
		}
	}
}

func (f *fswapEx) syncPairs(ctx context.Context) {
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
			f.updatePairs(ctx)
			sleepDur = time.Second
		}
	}
}

func (f *fswapEx) updatePairs(ctx context.Context) error {
	const uri = "/api/pairs"
	resp, err := fswapsdk.Request(ctx).Get(uri)
	if err != nil {
		return err
	}

	var body PairResp
	if err := fswapsdk.UnmarshalResponse(resp, &body); err != nil {
		return err
	}

	f.pairs = &body
	return nil
}

func (f *fswapEx) readTicker(ctx context.Context, asset *core.Asset) (*core.Ticker, error) {
	if f.pairs == nil {
		return nil, errors.New("pairs not avaialbe")
	}

	var (
		funds = decimal.New(1, 3)
	)
	bidOrder, err := fswapsdk.Route(f.pairs.Pairs, pusdAsset, asset.ID, funds)
	if err != nil {
		return nil, err
	}
	askOrder, err := fswapsdk.ReverseRoute(f.pairs.Pairs, asset.ID, pusdAsset, funds)
	if err != nil {
		return nil, err
	}

	x := bidOrder.PayAmount.Div(bidOrder.FillAmount).Truncate(8)
	y := askOrder.FillAmount.Div(askOrder.PayAmount).Truncate(8)

	t := &core.Ticker{
		Exchange: exchangeName,
		AssetID:  asset.ID,

		UpdatedAt: time.Unix(0, f.pairs.Timestamp*1000000),
		AskPrice:  decimal.Max(x, y),
		BidPrice:  decimal.Min(x, y),
		LastPrice: decimal.Avg(x, y),
	}

	return t, nil
}
