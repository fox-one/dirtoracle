package core

import (
	"context"
	"errors"
	"time"

	"github.com/lib/pq"
	"github.com/shopspring/decimal"
)

var (
	ErrInvalidTrace = errors.New("invalid trace")
)

type (
	Transfer struct {
		ID        int64           `sql:"PRIMARY_KEY" json:"id,omitempty"`
		CreatedAt time.Time       `json:"created_at,omitempty"`
		UpdatedAt time.Time       `json:"updated_at,omitempty"`
		TraceID   string          `sql:"type:char(36)" json:"trace_id,omitempty"`
		AssetID   string          `sql:"type:char(36)" json:"asset_id,omitempty"`
		Amount    decimal.Decimal `sql:"type:decimal(64,8)" json:"amount,omitempty"`
		Memo      string          `sql:"size:200" json:"memo,omitempty"`
		Threshold uint8           `json:"threshold,omitempty"`
		Opponents pq.StringArray  `sql:"type:varchar(1024)" json:"opponents,omitempty"`
	}

	WalletStore interface {
		ListTransfers(ctx context.Context, limit int) ([]*Transfer, error)
		CreateTransfers(ctx context.Context, transfers []*Transfer) error
		ExpireTransfers(ctx context.Context, transfers []*Transfer) error
	}

	WalletService interface {
		Transfer(ctx context.Context, transfer *Transfer) error
	}
)
