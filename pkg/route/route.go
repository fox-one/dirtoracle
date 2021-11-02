package route

import (
	"fmt"
	"sort"
	"strings"
)

const (
	MaxRouteLevel = 5
)

type (
	Route struct {
		Symbol  string `json:"symbol"`
		Asset   string `json:"asset"`
		Reverse bool   `json:"reverse,omitempty"`
	}

	Routes []*Route
)

func (routes Routes) String() string {
	var items = make([]string, len(routes))
	for i, route := range routes {
		if route.Reverse {
			items[i] = fmt.Sprintf("-%s", route.Symbol)
		} else {
			items[i] = route.Symbol
		}
	}
	return strings.Join(items, ";")
}

// FindRoutes: BFS search routes, returning the shortest routes
func FindRoutes(pairs []*Pair, from, to string) (Routes, bool) {
	var allRoutes = []Routes{{{Asset: from}}}
	for i := 0; i < MaxRouteLevel; i++ {
		newRoutes, state := findRoutes(pairs, allRoutes, to)
		switch state {
		case -1:
			return nil, false

		case 0:
			allRoutes = newRoutes

		case 1:
			return newRoutes[0][1:], true
		}
	}
	return nil, false
}

// findRoutes: return (routes, state)
//	state:
//		0 for pending
//		1 for success
//		-1 for failed
func findRoutes(pairs []*Pair, allRoutes []Routes, expectSymbol string) ([]Routes, int) {
	newRoutes := make([]Routes, 0, len(allRoutes))
	for _, routes := range allRoutes {
		if len(routes) == 0 {
			continue
		}

		from := routes[len(routes)-1].Asset
		routeAssets := make(map[string]bool, len(routes))
		for _, route := range routes {
			routeAssets[route.Asset] = true
		}

		for _, pair := range pairs {
			if pair.Base != from && pair.Quote != from {
				continue
			}

			route := Route{
				Symbol: pair.Symbol,
				Asset:  pair.Quote,
			}

			if pair.Quote == from {
				route.Asset = pair.Base
				route.Reverse = true
			}

			if route.Asset == expectSymbol {
				return []Routes{append(routes, &route)}, 1
			}

			// pair was not added to the previous routes
			if _, ok := routeAssets[route.Asset]; !ok {
				newRoutes = append(
					newRoutes,
					append(routes, &route),
				)
			}
		}
	}

	if len(newRoutes) == 0 {
		return nil, -1
	}

	sort.Slice(newRoutes, func(i, j int) bool {
		// sort routes with last route's quote asset
		a1, a2 := newRoutes[i][len(newRoutes[i])-1].Asset, newRoutes[j][len(newRoutes[j])-1].Asset
		i1, i2 := quoteIndex(a1), quoteIndex(a2)
		if i1 < i2 {
			return true
		} else if i1 == i2 {
			return a1 < a2
		}
		return false
	})

	return newRoutes, 0
}

var (
	quoteAssets = []string{
		"USDT",
		"USDC",
		"BUSD",
		"HUSD",
		"USDK",
		"BTC",
		"ETH",
		"BNB",
		"DAI",
	}
)

func quoteIndex(asset string) int {
	for i, s := range quoteAssets {
		if s == asset {
			return i
		}
	}

	return len(quoteAssets)
}
