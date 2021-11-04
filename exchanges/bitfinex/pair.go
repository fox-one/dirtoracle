package bitfinex

import (
	"context"
	"strings"
	"time"

	"github.com/fox-one/dirtoracle/pkg/route"
	"github.com/fox-one/pkg/logger"
)

const (
	pairsKey = "pairs"
)

func (exch *bitfinexEx) getPairs(ctx context.Context) ([]*route.Pair, error) {
	if pairs, ok := exch.cache.Get(pairsKey); ok {
		return pairs.([]*route.Pair), nil
	}

	log := logger.FromContext(ctx)
	resp, err := Request(ctx).Get("/conf/pub:map:currency:sym,pub:list:pair:exchange")
	if err != nil {
		log.WithError(err).Errorln("GET /conf/pub:map:currency:sym,pub:list:pair:exchange")
		return nil, err
	}

	var body [][]interface{}
	if err := UnmarshalResponse(resp, &body); err != nil {
		log.WithError(err).Errorln("getPairs.UnmarshalResponse")
		return nil, err
	}

	assetM := make(map[string]string)
	for _, m := range body[0] {
		arr := m.([]interface{})
		assetM[arr[0].(string)] = arr[1].(string)
	}

	var pairs []*route.Pair
	for _, m := range body[1] {
		symbol := m.(string)
		base := strings.TrimRight(symbol[:len(symbol)-3], ":")
		if s, ok := assetM[base]; ok {
			base = s
		}
		quote := symbol[len(symbol)-3:]
		if s, ok := assetM[quote]; ok {
			quote = s
		}
		pairs = append(pairs, &route.Pair{
			Symbol: "t" + symbol,
			Base:   strings.ToUpper(base),
			Quote:  strings.ToUpper(quote),
		})
	}

	exch.cache.Set(pairsKey, pairs, time.Hour)
	return pairs, nil
}
