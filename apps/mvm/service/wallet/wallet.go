package wallet

import (
	"context"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/fox-one/dirtoracle/apps/mvm/core"
	"github.com/fox-one/mixin-sdk-go"
)

type Config struct {
	Pin string `valid:"required"`
}

func New(client *mixin.Client, cfg Config) core.WalletService {
	if _, err := govalidator.ValidateStruct(cfg); err != nil {
		panic(err)
	}

	return &walletService{
		client: client,
		pin:    cfg.Pin,
	}
}

type walletService struct {
	client *mixin.Client
	pin    string
}

func (s *walletService) Poll(ctx context.Context, offset time.Time, limit int) ([]*core.Snapshot, error) {
	snapshots, err := s.client.ReadNetworkSnapshots(ctx, "", offset, "ASC", limit)
	if err != nil {
		return nil, err
	}

	return convertSnapshots(snapshots), nil
}

func (m *walletService) Transfer(ctx context.Context, req *core.Transfer) error {
	input := &mixin.TransferInput{
		AssetID: req.AssetID,
		Amount:  req.Amount,
		TraceID: req.TraceID,
		Memo:    req.Memo,
	}

	input.OpponentMultisig.Threshold = req.Threshold
	input.OpponentMultisig.Receivers = req.Opponents
	_, err := m.client.Transaction(ctx, input, m.pin)
	if e, ok := err.(*mixin.Error); ok && e.Code == mixin.InvalidTraceID {
		return core.ErrInvalidTrace
	}

	return err
}

func convertSnapshots(items []*mixin.Snapshot) []*core.Snapshot {
	var snapshots = make([]*core.Snapshot, len(items))
	for i, s := range items {
		snapshots[i] = &core.Snapshot{
			CreatedAt:  s.CreatedAt,
			SnapshotID: s.SnapshotID,
			UserID:     s.UserID,
			OpponentID: s.OpponentID,
			TraceID:    s.TraceID,
			AssetID:    s.AssetID,
			Amount:     s.Amount,
			Memo:       s.Memo,
		}
	}
	return snapshots
}
