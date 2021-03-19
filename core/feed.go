package core

import (
	"context"
	"fmt"
	"time"

	"github.com/lib/pq"
)

type (
	Asset struct {
		AssetID string `json:"asset_id,omitempty"`
		Symbol  string `json:"symbol,omitempty"`
	}

	FeedConfig struct {
		Asset
		Sources []string `json:"sources,omitempty"`
	}

	Subscriber struct {
		ID        uint           `sql:"PRIMARY_KEY" json:"id"`
		CreatedAt time.Time      `json:"created_at"`
		UpdatedAt time.Time      `json:"updated_at"`
		DeletedAt *time.Time     `sql:"index" json:"deleted_at"`
		AssetID   string         `json:"asset_id,omitempty"`
		Threshold uint8          `json:"threshold,omitempty"`
		Opponents pq.StringArray `sql:"type:TEXT" json:"opponents,omitempty"`
	}

	SubscriberStore interface {
		SaveSubscriber(ctx context.Context, f *Subscriber) error
		AllSubscribers(ctx context.Context) ([]*Subscriber, error)
		FindSubscribers(ctx context.Context, assetID string) ([]*Subscriber, error)
	}
)

func (p PriceData) Payload() []byte {
	return []byte(fmt.Sprintf("%d:%s:%v", p.Timestamp, p.AssetID, p.Price))
}
