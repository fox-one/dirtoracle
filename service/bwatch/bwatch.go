package bwatch

import (
	"context"

	"github.com/asaskevich/govalidator"
	"github.com/fox-one/dirtoracle/core"
)

const (
	serviceName = "bwatch"
)

type Config struct {
	ApiHost string `valid:"required"`
}

func New(cfg Config) core.PortfolioService {
	if _, err := govalidator.ValidateStruct(cfg); err != nil {
		panic(err)
	}

	if cfg.ApiHost != "" {
		client.SetHostURL(cfg.ApiHost)
	}
	return &bwatchService{}
}

type bwatchService struct{}

func (bwatchService) Name() string {
	return serviceName
}

func (bwatchService) ListPortfolioTokens(ctx context.Context) ([]*core.PortfolioToken, error) {
	r, err := request(ctx).Get("/api/etfs")
	if err != nil {
		return nil, err
	}

	var body struct {
		Etfs []*Etf `json:"etfs"`
	}
	if err = decodeResponse(r, &body); err != nil {
		return nil, err
	}

	var tokens []*core.PortfolioToken
	for _, etf := range body.Etfs {
		items := make([]*core.PortfolioItem, 0, len(etf.Assets))
		for id, amount := range etf.Assets {
			items = append(items, &core.PortfolioItem{
				AssetID: id,
				Amount:  amount,
			})
		}
		tokens = append(tokens, &core.PortfolioToken{
			AssetID: etf.AssetID,
			Items:   items,
		})
	}
	return tokens, nil
}
