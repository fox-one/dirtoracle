package oracle

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/fox-one/dirtoracle/core"
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
		msg.CreatedAt.Before(time.Now().Add(-m.config.MaxInterval)) {
		return nil
	}

	isMember := false
	for _, m := range m.system.Members {
		if m.ClientID == msg.UserID {
			isMember = true
		}
	}

	if !isMember {
		return nil
	}

	var p = new(core.PriceProposal)
	log := logger.FromContext(ctx)
	data, _ := base64.StdEncoding.DecodeString(msg.Data)
	if err := json.Unmarshal(data, p); err != nil {
		log.WithError(err).Errorln("Unmarshal PriceProposal failed")
		return nil
	}

	if p.Signature == nil && len(p.Signatures) == 0 {
		log.Errorln("empty signature proposal received")
		return nil
	}

	m.proposals <- p
	return nil
}
