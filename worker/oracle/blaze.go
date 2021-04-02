package oracle

import (
	"context"
	"time"

	"github.com/fox-one/mixin-sdk-go"
	"github.com/fox-one/pkg/logger"
)

func (m *Oracle) loopBlaze(ctx context.Context) error {
	log := logger.FromContext(ctx).WithField("worker", "blaze")
	ctx = logger.WithContext(ctx, log)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case <-time.After(time.Second):
			if err := m.client.LoopBlaze(ctx, mixin.BlazeListenFunc(m.OnMessage)); err != nil {
				log.WithError(err).Errorln("LoopBlaze failed")
			}
		}
	}
}

func (m *Oracle) OnMessage(ctx context.Context, msg *mixin.MessageView, userID string) error {
	if msg.Category != mixin.MessageCategoryPlainPost ||
		msg.ConversationID != m.system.ConversationID ||
		msg.CreatedAt.Before(time.Now().Add(-maxDuration)) {
		return nil
	}

	if msg.QuoteMessageID != "" {
		return m.handleProposalRespMessage(ctx, msg)
	}
	return m.handleProposalMessage(ctx, msg)
}
