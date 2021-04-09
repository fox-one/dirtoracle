package oracle

import (
	"context"
	"time"

	"github.com/fox-one/dirtoracle/core"
	"github.com/fox-one/pkg/logger"
	"golang.org/x/sync/errgroup"
)

func (m *Oracle) loopPortfolioTokens(ctx context.Context) error {
	log := logger.FromContext(ctx).WithField("loop", "portfolio-tokens")
	ctx = logger.WithContext(ctx, log)

	var sleepDur = time.Millisecond
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case <-time.After(sleepDur):
			var g errgroup.Group
			for _, p := range m.posrvs {
				p := p
				g.Go(func() error {
					arr, err := p.ListPortfolioTokens(ctx)
					if err != nil {
						log.WithError(err).Errorln("ListPortfolioTokens")
						return err
					}
					tokens := make(map[string]*core.PortfolioToken, len(arr))
					for _, t := range arr {
						tokens[t.AssetID] = t
					}
					return m.cacheServicePortfolioTokens(p.Name(), tokens)
				})
			}
			sleepDur = time.Second * 10
			if err := g.Wait(); err != nil {
				sleepDur = time.Second
			}
		}
	}
}

func (m *Oracle) unpackAsset(ctx context.Context, id string) ([]*core.PortfolioItem, error) {
	for _, p := range m.posrvs {
		tokens := m.cachedServicePortfolioTokens(p.Name())
		if t, ok := tokens[id]; ok {
			return t.Items, nil
		}
	}
	return nil, nil
}
