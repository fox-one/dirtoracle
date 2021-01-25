package core

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx/types"
	"github.com/lib/pq"
	"github.com/shopspring/decimal"
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
		Handled   types.BitBool   `sql:"type:bit(1)" json:"handled,omitempty"`
		Passed    types.BitBool   `sql:"type:bit(1)" json:"passed,omitempty"`
		Threshold uint8           `json:"threshold,omitempty"`
		Opponents pq.StringArray  `sql:"type:varchar(1024)" json:"opponents,omitempty"`
	}

	RawTransaction struct {
		ID        int64     `sql:"PRIMARY_KEY" json:"id,omitempty"`
		CreatedAt time.Time `json:"created_at,omitempty"`
		TraceID   string    `sql:"type:char(36);" json:"trace_id,omitempty"`
		Data      string    `sql:"type:TEXT" json:"data,omitempty"`
	}

	WalletStore interface {
		// Transfers
		CreateTransfers(ctx context.Context, transfers []*Transfer) error
	}
)
