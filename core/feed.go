package core

import (
	"context"

	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
)

type (
	FeedConfig struct {
		Asset
		Sources []string `json:"sources"`
	}

	Feeder struct {
		gorm.Model
		AssetID   string         `json:"asset_id"`
		Threshold uint8          `json:"threshold,omitempty"`
		Opponents pq.StringArray `sql:"type:TEXT" json:"opponents,omitempty"`
	}

	FeederStore interface {
		SaveFeeder(ctx context.Context, f *Feeder) error
		FindFeeders(ctx context.Context, assetID string) ([]*Feeder, error)
		AllFeeders(ctx context.Context) ([]*Feeder, error)
	}
)
