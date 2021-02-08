package wallet

import (
	"context"

	"github.com/fox-one/dirtoracle/core"
	"github.com/fox-one/pkg/store/db"
)

func init() {
	db.RegisterMigrate(func(db *db.DB) error {
		tx := db.Update().Model(core.Transfer{})
		if err := tx.AutoMigrate(core.Transfer{}).Error; err != nil {
			return err
		}

		if err := tx.AddUniqueIndex("idx_transfers_trace", "trace_id").Error; err != nil {
			return err
		}

		return nil
	})
}

func New(db *db.DB) core.WalletStore {
	return &walletStore{db: db}
}

type walletStore struct {
	db *db.DB
}

func (s *walletStore) ListTransfers(ctx context.Context) ([]*core.Transfer, error) {
	var ts []*core.Transfer
	err := s.db.View().Limit(100).Find(&ts).Error
	return ts, err
}

func (s *walletStore) CreateTransfers(ctx context.Context, transfers []*core.Transfer) error {
	return s.db.Tx(func(tx *db.DB) error {
		for _, transfer := range transfers {
			if err := tx.Update().Where("trace_id = ?", transfer.TraceID).FirstOrCreate(transfer).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *walletStore) ExpireTransfers(ctx context.Context, transfers []*core.Transfer) error {
	return s.db.Tx(func(tx *db.DB) error {
		for _, transfer := range transfers {
			if err := tx.Update().Where("trace_id = ?", transfer.TraceID).Delete(transfer).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
