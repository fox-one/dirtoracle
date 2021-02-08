package cashier

import (
	"context"
	"errors"
	"time"

	"github.com/fox-one/dirtoracle/core"
	"github.com/fox-one/pkg/logger"
	"golang.org/x/sync/errgroup"
)

func New(
	wallets core.WalletStore,
	walletz core.WalletService,
) *Cashier {
	return &Cashier{
		walletz: walletz,
		wallets: wallets,
	}
}

type Cashier struct {
	wallets core.WalletStore
	walletz core.WalletService
}

func (c *Cashier) Run(ctx context.Context) error {
	log := logger.FromContext(ctx).WithField("worker", "cashier")
	ctx = logger.WithContext(ctx, log)

	dur := time.Millisecond
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(dur):
			if err := c.run(ctx); err == nil {
				dur = 100 * time.Millisecond
			} else {
				dur = time.Second
			}
		}
	}
}

func (c *Cashier) run(ctx context.Context) error {
	log := logger.FromContext(ctx)

	payments, err := c.wallets.ListTransfers(ctx, 10)
	if err != nil {
		log.WithError(err).Errorln("list transfers")
		return err
	}

	if len(payments) == 0 {
		return errors.New("end of list")
	}

	var w errgroup.Group
	for _, payment := range payments {
		payment := payment
		w.Go(func() error {
			if err := c.handlePayment(ctx, payment); err != nil {
				return err
			}
			if err := c.wallets.ExpireTransfers(ctx, []*core.Transfer{payment}); err != nil {
				log.WithError(err).Errorln("delete finish transfers")
				return err
			}
			return nil
		})
	}

	return w.Wait()

}

func (c *Cashier) handlePayment(ctx context.Context, transfer *core.Transfer) error {
	log := logger.FromContext(ctx).WithField("trace", transfer.TraceID)

	if err := c.walletz.Transfer(ctx, transfer); err != nil && err != core.ErrInvalidTrace {
		log.WithError(err).Errorln("Transfer failed")
		return err
	}

	return nil
}
