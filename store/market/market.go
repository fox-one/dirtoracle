package asset

import (
	"context"
	"errors"
	"sync"

	"github.com/fox-one/dirtoracle/core"
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

func (s *marketStore) SaveTicker(ctx context.Context, ticker *core.Ticker) error {
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

func (s *marketStore) FindTickers(ctx context.Context, assetID string) ([]*core.Ticker, error) {
	m, ok := s.tickers[assetID]
	if !ok || len(m) == 0 {
		return nil, errors.New("tickers not avaiable")
	}

	var tickers = make([]*core.Ticker, 0, len(m))
	for _, t := range m {
		tickers = append(tickers, t)
	}
	return tickers, nil
}
