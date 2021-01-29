package exchange

import (
	"context"

	"github.com/fox-one/dirtoracle/core"
)

type (
	Handler interface {
		OnTicker(ctx context.Context, ticker *core.Ticker) error
	}

	Interface interface {
		Name() string
		// Subscribe subscribe exchange market events
		Subscribe(ctx context.Context, a *core.Asset, handler Handler) error
	}
)
