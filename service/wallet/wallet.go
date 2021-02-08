package wallet

import (
	"context"

	"github.com/fox-one/dirtoracle/core"
	"github.com/fox-one/mixin-sdk-go"
)

type mixinBot struct {
	client *mixin.Client
	pin    string
}

func New(client *mixin.Client, pin string) core.WalletService {
	return &mixinBot{
		client: client,
		pin:    pin,
	}
}

func (m *mixinBot) Transfer(ctx context.Context, req *core.Transfer) error {
	input := &mixin.TransferInput{
		AssetID: req.AssetID,
		Amount:  req.Amount,
		TraceID: req.TraceID,
		Memo:    req.Memo,
	}

	var err error
	if len(req.Opponents) == 1 {
		input.OpponentID = req.Opponents[0]
		_, err = m.client.Transfer(ctx, input, m.pin)
	} else {
		input.OpponentMultisig.Threshold = req.Threshold
		input.OpponentMultisig.Receivers = req.Opponents
		_, err = m.client.Transaction(ctx, input, m.pin)
	}

	if err != nil {
		if e, ok := err.(*mixin.Error); ok && e.Code == mixin.InvalidTraceID {
			return core.ErrInvalidTrace
		}
	}
	return err
}
