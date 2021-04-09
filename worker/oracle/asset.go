package oracle

import (
	"context"
	"time"

	"github.com/fox-one/dirtoracle/core"
	"github.com/fox-one/pkg/logger"
)

func (m *Oracle) loopTopAssets(ctx context.Context) error {
	log := logger.FromContext(ctx).WithField("loop", "top-assets")
	ctx = logger.WithContext(ctx, log)

	var sleepDur = time.Millisecond
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case <-time.After(sleepDur):
			assets, err := m.assetz.ReadTopAssets(ctx)
			if err != nil {
				log.WithError(err).Errorln("ReadTopAssets")
				sleepDur = time.Second
				break
			}

			m.cacheAssets(assets...)
			sleepDur = time.Minute * 10
		}
	}
}

func (m *Oracle) getAsset(ctx context.Context, id string) (*core.Asset, error) {
	if a := m.cachedAsset(id); a != nil {
		return a, nil
	}
	a, err := m.assetz.ReadAsset(ctx, id)
	if err != nil {
		if err == core.ErrAssetNotExist {
			return nil, nil
		}
		return nil, err
	}
	m.cacheAssets(a)
	return a, nil
}
