package core

import (
	"context"
	"time"
)

type (
	Subscriber struct {
		ID         uint       `sql:"PRIMARY_KEY" json:"id"`
		CreatedAt  time.Time  `json:"created_at"`
		UpdatedAt  time.Time  `json:"updated_at"`
		DeletedAt  *time.Time `sql:"index" json:"deleted_at"`
		Name       string     `sql:"SIZE:255;" json:"name,omitempty"`
		RequestURL string     `sql:"SIZE:255;" json:"request_url,omitempty"`
	}

	SubscriberStore interface {
		Save(ctx context.Context, f *Subscriber) error
		All(ctx context.Context) ([]*Subscriber, error)
	}
)
