package pricedata

import (
	"context"

	"github.com/fox-one/dirtoracle/core"
	"github.com/fox-one/pkg/store/db"
)

func init() {
	db.RegisterMigrate(func(db *db.DB) error {
		tx := db.Update().Model(core.PriceData{})
		if err := tx.AutoMigrate(&core.PriceData{}).Error; err != nil {
			return err
		}

		if err := tx.AddUniqueIndex("idx_price_data_asset_timestamp", "asset_id", "timestamp").Error; err != nil {
			return err
		}

		return nil
	})
}

func New(db *db.DB) core.PriceDataStore {
	return &pdataStore{db: db}
}

type pdataStore struct {
	db *db.DB
}

func (s *pdataStore) SavePriceData(ctx context.Context, p *core.PriceData) error {
	return s.db.Update().Where("asset_id = ? AND timestamp = ?", p.AssetID, p.Timestamp).FirstOrCreate(p).Error
}

func (s *pdataStore) LatestPriceData(ctx context.Context, asset string) (*core.PriceData, error) {
	var p core.PriceData
	if err := s.db.View().Where("asset_id = ?", asset).Order("timestamp DESC").First(&p).Error; err != nil {
		return nil, err
	}
	return &p, nil
}
