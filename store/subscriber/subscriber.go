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

func (s *subscriberStore) Save(_ context.Context, f *core.Subscriber) error {
	return s.db.Update().Model(f).Assign(core.Subscriber{Name: f.Name}).
		Where("request_url = ?", f.RequestURL).
		FirstOrCreate(f).Error
}

func (s *subscriberStore) All(_ context.Context) ([]*core.Subscriber, error) {
	var subscribers []*core.Subscriber
	if err := s.db.View().Find(&subscribers).Error; err != nil {
		return nil, err
	}
	return subscribers, nil
}
