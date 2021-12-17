package asset

import (
	"context"

	"github.com/fox-one/dirtoracle/apps/mvm/core"
	"github.com/fox-one/pkg/store/db"
)

func init() {
	db.RegisterMigrate(func(db *db.DB) error {
		tx := db.Update().Model(core.Asset{})

		if err := tx.AutoMigrate(core.Asset{}).Error; err != nil {
			return err
		}

		return nil
	})
}

func New(db *db.DB) core.AssetStore {
	return &assetStore{
		db: db,
	}
}

type assetStore struct {
	db *db.DB
}

func (s *assetStore) List(ctx context.Context) ([]*core.Asset, error) {
	var assets []*core.Asset
	if err := s.db.View().Find(&assets).Error; err != nil {
		return nil, err
	}

	return assets, nil
}

func (s *assetStore) Find(ctx context.Context, assetID string) (*core.Asset, error) {
	var asset core.Asset
	if err := s.db.View().Where("asset_id = ?", assetID).First(&asset).Error; err != nil {
		if db.IsErrorNotFound(err) {
			return &core.Asset{}, nil
		}
		return nil, err
	}

	return &asset, nil
}

func (s *assetStore) Update(ctx context.Context, asset *core.Asset) error {
	updates := toUpdateParams(asset)
	if tx := s.db.Update().Model(asset).Updates(updates); tx.Error != nil {
		return tx.Error
	}
	return nil
}

func toUpdateParams(asset *core.Asset) map[string]interface{} {
	return map[string]interface{}{
		"price":            asset.Price,
		"price_updated_at": asset.PriceUpdatedAt,
	}
}
