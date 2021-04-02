package exchanges

import (
	"context"

	"github.com/fox-one/dirtoracle/core"
)

type Exchange struct {
	// some assets were only listed on 4swap,
	//	to avoid reading wrong asset prices with same symbol,
	//	should check if the asset was in blacklist before sendding price requests
	blacklist map[string]bool
}

func New() *Exchange {
	return &Exchange{
		blacklist: map[string]bool{
			"XIN": true,
			"BOX": true,
		},
	}
}

func (e *Exchange) IsAssetBlocked(ctx context.Context, a *core.Asset) bool {
	_, ok := e.blacklist[a.Symbol]
	return ok
}
