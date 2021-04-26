package asset

import (
	"context"

	"github.com/fox-one/dirtoracle/core"
	"github.com/patrickmn/go-cache"
	"golang.org/x/sync/singleflight"
)

func Cache(assetz core.AssetService) core.AssetService {
	return &cacheAssetz{
		AssetService: assetz,
		caches:       cache.New(cache.NoExpiration, cache.NoExpiration),
		sf:           &singleflight.Group{},
	}
}

type cacheAssetz struct {
	core.AssetService
	caches *cache.Cache
	sf     *singleflight.Group
}

func (s *cacheAssetz) ReadAsset(ctx context.Context, id string) (*core.Asset, error) {
	v, err, _ := s.sf.Do(id, func() (interface{}, error) {
		if v, ok := s.caches.Get(id); ok {
			return v, nil
		}

		asset, err := s.AssetService.ReadAsset(ctx, id)
		if err != nil {
			return nil, err
		}

		s.caches.SetDefault(id, asset)
		return asset, nil
	})

	if err != nil {
		return nil, err
	}

	return v.(*core.Asset), nil
}
