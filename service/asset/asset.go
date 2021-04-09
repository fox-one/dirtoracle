package asset

import (
	"context"

	"github.com/fox-one/dirtoracle/core"
	"github.com/fox-one/mixin-sdk-go"
)

func New(c *mixin.Client) core.AssetService {
	return &assetService{c: c}
}

type assetService struct {
	c *mixin.Client
}

func (s *assetService) ReadAsset(ctx context.Context, id string) (*core.Asset, error) {
	asset, err := s.c.ReadAsset(ctx, id)
	if err != nil {
		if mixin.IsErrorCodes(err, 10002) {
			err = core.ErrAssetNotExist
		}

		return nil, err
	}

	return convertAsset(asset), nil
}

func (s *assetService) ReadTopAssets(ctx context.Context) ([]*core.Asset, error) {
	assets, err := mixin.ReadTopNetworkAssets(ctx)
	if err != nil {
		return nil, err
	}

	return convertAssets(assets), nil
}

func convertAsset(asset *mixin.Asset) *core.Asset {
	return &core.Asset{
		AssetID: asset.AssetID,
		Symbol:  asset.Symbol,
	}
}

func convertAssets(assets []*mixin.Asset) []*core.Asset {
	out := make([]*core.Asset, len(assets))
	for idx, asset := range assets {
		out[idx] = convertAsset(asset)
	}

	return out
}
