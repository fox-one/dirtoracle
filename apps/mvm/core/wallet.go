package core

import (
	"context"
	"errors"
	"time"

	"github.com/shopspring/decimal"
)

var (
	ErrInvalidTrace = errors.New("invalid trace")
)

type (
	Snapshot struct {
		SnapshotID string
		TraceID    string
		CreatedAt  time.Time
		UserID     string
		OpponentID string
		AssetID    string
		Amount     decimal.Decimal
		Memo       string
	}

	Transfer struct {
		AssetID   string
		TraceID   string
		Amount    decimal.Decimal
		Memo      string
		Opponents []string
		Threshold uint8
	}

	WalletService interface {
		Poll(ctx context.Context, offset time.Time, limit int) ([]*Snapshot, error)
		Transfer(ctx context.Context, transfer *Transfer) error
	}
)
