package subscriber

import (
	"context"

	"github.com/fox-one/dirtoracle/core"
	"github.com/fox-one/pkg/store/db"
)

func init() {
	db.RegisterMigrate(func(db *db.DB) error {
		tx := db.Update().Model(core.Subscriber{})

		if err := tx.AutoMigrate(core.Subscriber{}).Error; err != nil {
			return err
		}

		return nil
	})
}

func New(db *db.DB) core.SubscriberStore {
	return &subscriberStore{
		db: db,
	}
}

type subscriberStore struct {
	db *db.DB
}

func (s *subscriberStore) SaveSubscriber(_ context.Context, f *core.Subscriber) error {
	return s.db.Update().Model(f).
		Where("asset_id = ? AND threshold = ? AND opponents = ?", f.AssetID, f.Threshold, f.Opponents).
		FirstOrCreate(f).Error
}

func (s *subscriberStore) AllSubscribers(_ context.Context) ([]*core.Subscriber, error) {
	var subscribers []*core.Subscriber
	if err := s.db.View().Find(&subscribers).Error; err != nil {
		return nil, err
	}
	return subscribers, nil
}

func (s *subscriberStore) FindSubscribers(ctx context.Context, assetID string) ([]*core.Subscriber, error) {
	var subscribers []*core.Subscriber
	if err := s.db.View().Where("asset_id = ?", assetID).Find(&subscribers).Error; err != nil {
		return nil, err
	}
	return subscribers, nil
}
