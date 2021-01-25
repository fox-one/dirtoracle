package core

import (
	"context"
	"errors"
	"time"

	"github.com/fox-one/dirtoracle/pkg/number"
	"github.com/shopspring/decimal"
)

var (
	ErrAssetNotExist = errors.New("asset not exist")
)

type (
	Asset struct {
		ID            string          `sql:"size:36;PRIMARY_KEY" json:"id,omitempty"`
		UpdatedAt     time.Time       `json:"updated_at,omitempty"`
		Name          string          `sql:"size:64" json:"name,omitempty"`
		Symbol        string          `sql:"size:32" json:"symbol,omitempty"`
		DisplaySymbol string          `sql:"size:32" json:"display_symbol,omitempty"`
		Logo          string          `sql:"size:256" json:"logo,omitempty"`
		ChainID       string          `sql:"size:36" json:"chain_id,omitempty"`
		Price         decimal.Decimal `sql:"type:decimal(24,8)" json:"price_usd,omitempty"`
	}

	// AssetStore defines operations for working with assets on db.
	AssetStore interface {
		Save(ctx context.Context, asset *Asset, columns ...string) error
		Find(ctx context.Context, id string) (*Asset, error)
		ListAll(ctx context.Context) ([]*Asset, error)
		ListPrices(ctx context.Context, ids ...string) (number.Values, error)
	}

	// AssetService provides access to assets information
	// in the remote system like mixin network.
	AssetService interface {
		Find(ctx context.Context, id string) (*Asset, error)
		ListAll(ctx context.Context) ([]*Asset, error)
	}
)
