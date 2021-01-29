package market

import (
	"context"
	"errors"
	"sort"
	"sync"
	"time"

	"github.com/fox-one/dirtoracle/core"
	"github.com/shopspring/decimal"
)

type (
	tickerMap map[string]*core.Ticker
)

func New() core.MarketStore {
	return &marketStore{
		tickers: map[string]tickerMap{},
	}
}

type marketStore struct {
	tickers map[string]tickerMap
	lock    sync.Mutex
}

func (s *marketStore) SaveTicker(_ context.Context, ticker *core.Ticker) error {
	s.lock.Lock()
	m, ok := s.tickers[ticker.AssetID]
	if !ok {
		m = tickerMap{}
	}
	m[ticker.Source] = ticker
	s.tickers[ticker.AssetID] = m
	s.lock.Unlock()
	return nil
}

func (s *marketStore) FindTickers(_ context.Context, assetID string) ([]*core.Ticker, error) {
	m, ok := s.tickers[assetID]
	if !ok {
		return nil, errors.New("tickers not avaiable")
	}

	ts := make([]*core.Ticker, 0, len(m))
	// remove outdated prices
	for _, t := range m {
		ts = append(ts, t)
	}
	return ts, nil
}

func (s *marketStore) AggregateTickerPrices(ctx context.Context, assetID string) (decimal.Decimal, error) {
	ts, err := s.FindTickers(ctx, assetID)
	if err != nil {
		return decimal.Zero, err
	}

	sort.Slice(ts, func(i, j int) bool {
		return ts[i].Price.LessThan(ts[j].Price)
	})

	{
		var (
			index = 0
			d     = time.Now().Add(-15*time.Second).Unix() * 1000
		)

		for _, t := range ts {
			// validate ticker:
			// 	price must be positive
			// 	volume must be positive
			// 	updated within 15s
			if t.Price.IsPositive() &&
				t.VolumeUSD.IsPositive() &&
				t.Timestamp > d {
				ts[index] = t
				index++
			}
		}

		if index < 2 {
			return decimal.Zero, errors.New("no enough valid tickers")
		}
		ts = ts[:index]
	}

	{
		var (
			index      = 0
			one        = decimal.New(1, 0)
			threshold  = decimal.New(5, -2)
			mid        = ts[len(ts)/2].Price
			volume     = decimal.Zero
			totalValue = decimal.Zero
		)

		for _, t := range ts {
			// 	price diff less than threshold
			if t.Price.Div(mid).Sub(one).Abs().LessThan(threshold) {
				index++
				volume = volume.Add(t.VolumeUSD)
				totalValue = totalValue.Add(t.Price.Mul(t.VolumeUSD))
			}
		}
		if index >= 2 {
			return totalValue.Div(volume).Truncate(8), nil
		}
	}
	return decimal.Zero, errors.New("no enough valid tickers")
}
