package feeder

import (
	"context"

	"github.com/fox-one/dirtoracle/core"
	"github.com/fox-one/pkg/store/db"
)

func init() {
	db.RegisterMigrate(func(db *db.DB) error {
		tx := db.Update().Model(core.Feeder{})

		if err := tx.AutoMigrate(core.Feeder{}).Error; err != nil {
			return err
		}

		return nil
	})
}

func New(db *db.DB) core.FeederStore {
	return &feederStore{
		db: db,
	}
}

type feederStore struct {
	db *db.DB
}

func (s *feederStore) SaveFeeder(_ context.Context, f *core.Feeder) error {
	return s.db.Update().Model(f).
		Where("asset_id = ? AND threshold = ? AND opponents = ?", f.AssetID, f.Threshold, f.Opponents).
		FirstOrCreate(f).Error
}

func (s *feederStore) AllFeeders(_ context.Context) ([]*core.Feeder, error) {
	var feeders []*core.Feeder
	if err := s.db.View().Find(&feeders).Error; err != nil {
		return nil, err
	}
	return feeders, nil
}

func (s *feederStore) FindFeeders(ctx context.Context, assetID string) ([]*core.Feeder, error) {
	var feeders []*core.Feeder
	if err := s.db.View().Where("asset_id = ?", assetID).Find(&feeders).Error; err != nil {
		return nil, err
	}
	return feeders, nil
}
