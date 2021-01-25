package exchange

import (
	"context"

	"github.com/fox-one/dirtoracle/core"
)

type (
	MarketHandler interface {
		OnTicker(ctx context.Context, asset *core.Asset, ticker *core.Ticker) error
	}

	Interface interface {
		Name() string
		// Subscribe subscribe exchange market events
		Subscribe(ctx context.Context, asset *core.Asset, handler MarketHandler) error
	}
)
