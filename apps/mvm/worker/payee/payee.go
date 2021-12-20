package payee

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/fox-one/dirtoracle/apps/mvm/core"
	"github.com/fox-one/dirtoracle/apps/mvm/encoding"
	"github.com/fox-one/pkg/logger"
	"github.com/fox-one/pkg/uuid"
	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
)

const (
	snapshotKey = "snapshot:%s"
)

type (
	Worker struct {
		system  core.System
		assets  core.AssetStore
		walletz core.WalletService
		cache   *cache.Cache

		checkpoint time.Time
	}
)

func New(
	system core.System,
	assets core.AssetStore,
	walletz core.WalletService,
) *Worker {
	return &Worker{
		system:  system,
		assets:  assets,
		walletz: walletz,
		cache:   cache.New(time.Hour, time.Hour),

		checkpoint: time.Now().Add(-time.Minute * 10),
	}
}

func (w *Worker) Run(ctx context.Context) error {
	log := logger.FromContext(ctx).WithField("worker", "payee")
	ctx = logger.WithContext(ctx, log)

	sleepDur := time.Millisecond
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(sleepDur):
			if err := w.run(ctx); err != nil {
				sleepDur = time.Second
			} else {
				sleepDur = time.Millisecond * 100
			}
		}
	}
}

func (w *Worker) run(ctx context.Context) error {
	const LIMIT = 500
	snapshots, err := w.walletz.Poll(ctx, w.checkpoint, LIMIT)
	if err != nil {
		logger.FromContext(ctx).WithError(err).Errorln("list snapshots")
		return err
	}

	for _, snapshot := range snapshots {
		if snapshot.UserID != w.system.ClientID || snapshot.Amount.IsNegative() {
			w.checkpoint = snapshot.CreatedAt
			continue
		}

		if _, ok := w.cache.Get(fmt.Sprintf(snapshotKey, snapshot.SnapshotID)); ok {
			continue
		}
		if err := w.handleSnapshot(ctx, snapshot); err != nil {
			return nil
		}
		w.cache.Set(fmt.Sprintf(snapshotKey, snapshot.SnapshotID), true, time.Hour)
		w.checkpoint = snapshot.CreatedAt
	}

	if len(snapshots) < LIMIT {
		return fmt.Errorf("no more snapshots")
	}

	return nil
}

func (w *Worker) handleSnapshot(ctx context.Context, snapshot *core.Snapshot) error {
	log := logger.FromContext(ctx).WithFields(logrus.Fields{
		"snapshot": snapshot.SnapshotID,
		"ss_time":  snapshot.CreatedAt,
	})
	data, err := base64.StdEncoding.DecodeString(snapshot.Memo)
	if err != nil {
		log.WithError(err).Errorln("base64 decode")
		return nil
	}

	var p core.PriceData
	if err := p.UnmarshalBinary(data); err != nil {
		log.WithError(err).Errorln("PriceData.UnmarshalBinary")
		return nil
	}

	log = log.WithFields(logrus.Fields{
		"price":      p.Price,
		"price_time": time.Unix(p.Timestamp, 0),
		"asset":      p.AssetID,
	})

	if !w.system.VerifyData(&p) {
		log.WithError(err).Errorln("VerifyData failed")
		return nil
	}

	data, _ = p.MarshalBinary()
	op := &encoding.Operation{
		Purpose: encoding.OperationPurposeGroupEvent,
		Process: w.system.MvmProcess,
		Extra:   data,
	}

	if err := w.walletz.Transfer(ctx, &core.Transfer{
		AssetID:   w.system.GasAsset,
		Amount:    w.system.GasAmount,
		TraceID:   uuid.Modify(snapshot.SnapshotID, "forward-to-mvm"),
		Memo:      base64.RawURLEncoding.EncodeToString(op.Encode()),
		Opponents: w.system.MvmGroups,
		Threshold: w.system.MvmThreshold,
	}); err != nil {
		log.WithError(err).Errorln("walletz.Transfer")
		return err
	}

	if asset, err := w.assets.Find(ctx, p.AssetID); err != nil {
		log.WithError(err).Errorln("assets.Find", p.AssetID)
		return err
	} else if asset.ID == 0 {
		log.Infoln("asset not found", p.AssetID)
		return err
	} else {
		updatedAt := time.Unix(p.Timestamp, 0)
		asset.Price = p.Price
		asset.PriceUpdatedAt = &updatedAt
		if err := w.assets.Update(ctx, asset); err != nil {
			log.WithError(err).Errorln("assets.Update", p.AssetID, p.Price, asset.PriceUpdatedAt)
			return err
		}
	}

	return nil
}
