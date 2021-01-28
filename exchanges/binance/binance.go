package binance

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/fox-one/dirtoracle/core"
	"github.com/fox-one/dirtoracle/core/exchange"
	"github.com/fox-one/pkg/logger"
	"github.com/gorilla/websocket"
	"golang.org/x/sync/errgroup"
)

const (
	exchangeName = "binance"
)

type binanceEx struct {
	dialer websocket.Dialer
}

func New() exchange.Interface {
	return &binanceEx{
		dialer: websocket.Dialer{
			Subprotocols:   []string{"json"},
			ReadBufferSize: 1024,
		},
	}
}

func (b *binanceEx) Name() string {
	return exchangeName
}

func (b *binanceEx) Subscribe(ctx context.Context, asset *core.Asset, handler exchange.MarketHandler) error {
	log := logger.FromContext(ctx)
	log.Info("start")
	defer log.Info("quit")

	var (
		g      errgroup.Group
		ticker *core.Ticker
	)

	g.Go(func() error {
		pairSymbol := b.pairSymbol(b.assetSymbol(asset.Symbol))
		stream := strings.ToLower(pairSymbol) + "@miniTicker"
		url := fmt.Sprintf("%s/stream?streams=%s", WebsocketEndpoint, stream)

		conn, _, err := b.dialer.Dial(url, nil)
		if err != nil {
			log.WithError(err).Errorln("dail failed")
			return err
		}

		var msg struct {
			Ticker Ticker `json:"data"`
			Stream string `json:"stream"`
		}
		for {
			if err := conn.ReadJSON(&msg); err != nil {
				log.WithError(err).Errorln("read json failed")
				return err
			}

			if msg.Stream != stream {
				log.WithField("stream", stream).Debugln("receive unknown message")
				continue
			}
			ticker = convertTicker(&msg.Ticker)
		}
	})

	g.Go(func() error {
		var (
			sleepDur  = time.Second
			updatedAt time.Time
		)
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()

			case <-time.After(sleepDur):
				if ticker == nil || ticker.UpdatedAt.Before(updatedAt) {
					sleepDur = time.Second
					continue
				}
				ticker.AssetID = asset.ID
				if err := handler.OnTicker(ctx, asset, ticker); err != nil {
					log.WithError(err).Errorln("OnTicker failed")
					sleepDur = time.Second
					continue
				}

				updatedAt = ticker.UpdatedAt
				sleepDur = 10 * time.Second
			}
		}
	})

	return g.Wait()
}

func (b *binanceEx) assetSymbol(symbol string) string {
	return symbol
}

func (b *binanceEx) pairSymbol(symbol string) string {
	return symbol + "BUSD"
}
