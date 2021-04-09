package core

import (
	"context"
	"errors"
)

var (
	ErrAssetNotExist = errors.New("asset not exist")
)

type (
	Asset struct {
		AssetID string `json:"asset_id,omitempty"`
		Symbol  string `json:"symbol,omitempty"`
	}

	AssetService interface {
		ReadAsset(ctx context.Context, id string) (*Asset, error)
		ReadTopAssets(ctx context.Context) ([]*Asset, error)
	}
)
