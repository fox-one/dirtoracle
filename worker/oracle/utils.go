package oracle

import (
	"context"
	"time"
)

func (m *Oracle) execWithTimeout(ctx context.Context, timeout time.Duration, f func() error) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	sleepDur := time.Millisecond
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case <-time.After(sleepDur):
			if err := f(); err == nil {
				return nil
			}
			sleepDur = time.Second
		}
	}
}
