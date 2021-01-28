package market

import (
	"context"
	"errors"
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
	m[ticker.Exchange] = ticker
	s.tickers[ticker.AssetID] = m
	s.lock.Unlock()
	return nil
}

func (s *marketStore) FindTickers(_ context.Context, assetID string) ([]*core.Ticker, error) {
	m, ok := s.tickers[assetID]
	if !ok {
		return nil, errors.New("tickers not avaiable")
	}

	var (
		ts = make([]*core.Ticker, 0, len(m))
		d  = time.Now().Add(-15 * time.Second)
	)
	for _, t := range m {
		if t.UpdatedAt.After(d) && t.Price.IsPositive() && t.VolumeUSD.IsPositive() {
			ts = append(ts, t)
		}
	}

	if len(ts) < 3 {
		return nil, errors.New("tickers outdated")
	}

	return ts, nil
}

func (s *marketStore) AggregateTickerPrices(ctx context.Context, assetID string) (decimal.Decimal, error) {
	tickers, err := s.FindTickers(ctx, assetID)
	if err != nil {
		return decimal.Zero, err
	}

	var (
		volume     = decimal.Zero
		totalValue = decimal.Zero
	)

	for _, t := range tickers {
		volume = volume.Add(t.VolumeUSD)
		totalValue = totalValue.Add(t.Price.Mul(t.VolumeUSD))
	}

	return totalValue.Div(volume).Truncate(8), nil
}
