package exinswap

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	fswapsdk "github.com/fox-one/4swap-sdk-go"
	"github.com/fox-one/dirtoracle/core"
	"github.com/fox-one/dirtoracle/core/exchange"
	"github.com/fox-one/dirtoracle/pkg/number"
	"github.com/fox-one/pkg/logger"
	"github.com/shopspring/decimal"
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

func (e *eswapEx) Subscribe(ctx context.Context, a *core.Asset, h exchange.Handler) error {
	log := logger.FromContext(ctx)
	log.Info("start")
	defer log.Info("quit")

	e.once.Do(func() {
		go e.syncPairs(ctx)
	})

	var (
		sleepDur      = time.Second
		lastTimestamp int64
	)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case <-time.After(sleepDur):
			if e.pairs == nil || e.pairs.Timestamp <= lastTimestamp {
				sleepDur = time.Second
				continue
			}

			ticker, err := e.readTicker(ctx, a)
			if err != nil {
				log.WithError(err).Errorln("readTicker failed")
				sleepDur = time.Second
				continue
			}
			if err := h.OnTicker(ctx, ticker); err != nil {
				log.WithError(err).Errorln("OnTicker failed")
				sleepDur = time.Second
				continue
			}

			sleepDur = 3 * time.Second
			lastTimestamp = e.pairs.Timestamp
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
			e.updatePairs(ctx)
			sleepDur = time.Second
		}
	}
}

func (e *eswapEx) updatePairs(ctx context.Context) error {
	const uri = "/pairs"

	log := logger.FromContext(ctx)
	resp, err := Request(ctx).Get(uri)
	if err != nil {
		log.WithError(err).Errorln("fetch pairs failed")
		return err
	}

	var body struct {
		Pairs []struct {
			Asset0 struct {
				ID    string          `json:"uuid"`
				Price decimal.Decimal `json:"priceUsdt"`
			} `json:"asset0"`
			Asset1 struct {
				ID    string          `json:"uuid"`
				Price decimal.Decimal `json:"priceUsdt"`
			} `json:"asset1"`
			Balance0  decimal.Decimal `json:"asset0Balance"`
			Balance1  decimal.Decimal `json:"asset1Balance"`
			Volume    decimal.Decimal `json:"usdtTradeVolume24hours"`
			TradeType string          `json:"tradeType"`
		} `json:"data"`
		Timestamp int64 `json:"timestampMs"`
	}

	if err := UnmarshalResponse(resp, &body); err != nil {
		log.WithError(err).Errorln("fetch pairs failed")
		return err
	}

	var pairs = PairResp{
		Timestamp: body.Timestamp,
		Pairs:     make([]*fswapsdk.Pair, len(body.Pairs)),
	}

	fee := decimal.New(3, -3)
	curveFee := decimal.New(4, -4)
	for i, p := range body.Pairs {
		pair := &fswapsdk.Pair{
			RouteID:        int64(i),
			BaseAssetID:    p.Asset0.ID,
			BaseAmount:     p.Balance0,
			QuoteAssetID:   p.Asset1.ID,
			QuoteAmount:    p.Balance1,
			FeePercent:     fee,
			Volume24h:      p.Volume,
			BaseVolume24h:  p.Volume.Div(p.Asset0.Price),
			QuoteVolume24h: p.Volume.Div(p.Asset1.Price),
		}

		if p.TradeType == "curve" {
			pair.SwapMethod = "curve"
			pair.FeePercent = curveFee
		}

		pairs.Pairs[i] = pair
	}
	e.pairs = &pairs
	return nil
}

func (e *eswapEx) readTicker(ctx context.Context, asset *core.Asset) (*core.Ticker, error) {
	if e.pairs == nil {
		return nil, errors.New("pairs not avaialbe")
	}

	var (
		funds = decimal.New(1, 3)
		pairs = e.pairs.Pairs
		t     = core.Ticker{
			Source:    exchangeName,
			AssetID:   asset.AssetID,
			Timestamp: e.pairs.Timestamp,
		}
	)

	order, err := fswapsdk.Route(pairs, pusdAsset, asset.AssetID, funds)
	if err != nil {
		return nil, err
	}

	t.Price = order.PayAmount.Div(order.FillAmount).Truncate(8)

	volumes := number.Values{}
	for _, p := range pairs {
		volumes.Set(fmt.Sprint(p.RouteID), p.Volume24h)
	}

	for _, id := range fswapsdk.DecodeRoutes(order.Routes) {
		if v := volumes.Get(fmt.Sprint(id)); t.VolumeUSD.IsZero() || v.LessThan(t.VolumeUSD) {
			t.VolumeUSD = v
		}
	}

	return &t, nil
}
