package exinswap

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
	exchangeName = "exinswap"
)

type (
	PairResp struct {
		Pairs     []*fswapsdk.Pair `json:"pairs"`
		Timestamp int64            `json:"ts"`
	}

	eswapEx struct {
		once  sync.Once
		pairs *PairResp
	}
)

func New() exchange.Interface {
	return &eswapEx{}
}

func (*eswapEx) Name() string {
	return exchangeName
}

func (e *eswapEx) Subscribe(ctx context.Context, asset *core.Asset, handler exchange.MarketHandler) error {
	log := logger.FromContext(ctx).WithFields(logrus.Fields{
		"ex":     "exinswap",
		"worker": "subscribe",
	})
	ctx = logger.WithContext(ctx, log)

	e.once.Do(func() {
		e.syncPairs(ctx)
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
			ticker, err := e.readTicker(ctx, asset)
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

func (e *eswapEx) syncPairs(ctx context.Context) {
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
			sleepDur = time.Second

			resp, err := e.listPairs(ctx)
			if err != nil {
				log.WithError(err).Errorln("list pairs")
				continue
			}
			e.pairs = resp
		}
	}
}

func (e *eswapEx) listPairs(ctx context.Context) (*PairResp, error) {
	const uri = "/pairs"

	log := logger.FromContext(ctx)
	resp, err := Request(ctx).Get(uri)
	if err != nil {
		log.WithError(err).Errorln("fetch pairs failed")
		return nil, err
	}

	var body struct {
		Pairs []struct {
			Asset0 struct {
				ID string `json:"uuid"`
			} `json:"asset0"`
			Asset1 struct {
				ID string `json:"uuid"`
			} `json:"asset1"`
			Balance0 decimal.Decimal `json:"asset0Balance"`
			Balance1 decimal.Decimal `json:"asset1Balance"`
		} `json:"body"`
		Timestamp int64 `json:"timestampMs"`
	}

	if err := UnmarshalResponse(resp, &body); err != nil {
		log.WithError(err).Errorln("fetch pairs failed")
		return nil, err
	}

	var pairs = PairResp{
		Timestamp: body.Timestamp,
		Pairs:     make([]*fswapsdk.Pair, len(body.Pairs)),
	}

	fee := decimal.New(3, -3)
	for i, p := range body.Pairs {
		pairs.Pairs[i] = &fswapsdk.Pair{
			RouteID:      int64(i),
			BaseAssetID:  p.Asset0.ID,
			BaseAmount:   p.Balance0,
			QuoteAssetID: p.Asset1.ID,
			QuoteAmount:  p.Balance1,
			FeePercent:   fee,
		}
	}
	return &pairs, nil
}

func (f *eswapEx) readTicker(ctx context.Context, asset *core.Asset) (*core.Ticker, error) {
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
