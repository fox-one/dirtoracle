package oracle

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/fox-one/dirtoracle/core"
	"github.com/fox-one/mixin-sdk-go"
	"github.com/fox-one/pkg/logger"
	"github.com/fox-one/pkg/uuid"
	"github.com/sirupsen/logrus"
)

func (m *Oracle) handleProposalMessage(ctx context.Context, msg *mixin.MessageView) error {
	log := logger.FromContext(ctx)
	var p = new(core.ProposalRequest)
	data, _ := base64.StdEncoding.DecodeString(msg.Data)
	if err := json.Unmarshal(data, p); err != nil {
		log.WithError(err).Errorln("Unmarshal ProposalRequest failed")
		return nil
	}

	log = log.WithFields(logrus.Fields{
		"timestamp": p.Timestamp,
		"asset_id":  p.AssetID,
		"price":     p.Price,
	})
	ctx = logger.WithContext(ctx, log)

	// validate proposal request
	{ // proposal too old, just ignore
		if time.Unix(p.Timestamp, 0).Before(time.Now().Add(-maxDuration)) {
			log.Infoln("ignore:", "proposal old proposal")
			return nil
		}

		if !p.Verify(p.Signature) {
			log.Infoln("ignore:", "invalid proposal signature")
			return nil
		}
	}

	var signer *core.Signer

	// skip if local node is not expected as a signer
	{
		for _, s := range p.Signers {
			if s.VerifyKey != nil && s.VerifyKey.String() == m.system.VerifyKey.String() {
				s.En256VerifyKey = nil
				signer = s
			} else if s.En256VerifyKey != nil && s.En256VerifyKey.String() == m.system.En256VerifyKey.String() {
				signer = s
			}
		}

		if signer == nil {
			log.Infoln("ignore:", "node not approved by the requester")
			return nil
		}
	}

	// compare with local price
	{
		price, err := m.getPrice(ctx, &p.Asset)
		if err != nil {
			return err
		}

		change := p.Price.Sub(price).Div(p.Price)
		log = log.WithFields(logrus.Fields{
			"change":         change,
			"price":          price,
			"price.proposal": p.Price,
		})

		if change.Abs().GreaterThan(priceChangeThreshold) {
			log.Infoln("ignore:", "price diff too large")
			return nil
		}
	}

	// send proposal response back to the mixin conversation
	{
		resp, err := m.system.SignProposal(p, signer)
		if err != nil {
			log.WithError(err).Errorln("SignProposal failed")
			return err
		}

		bts, _ := json.MarshalIndent(resp, "", "    ")
		reply := &mixin.MessageRequest{
			ConversationID: msg.ConversationID,
			QuoteMessageID: msg.MessageID,
			RecipientID:    msg.UserID,
			MessageID:      uuid.Modify(msg.MessageID, fmt.Sprintf("reply:%d", signer.Index)),
			Category:       mixin.MessageCategoryPlainPost,
			Data:           base64.StdEncoding.EncodeToString(bts),
		}
		if err := m.client.SendMessage(ctx, reply); err != nil {
			log.WithError(err).Errorln("SendMessage failed")
			return err
		}
	}

	log.Infoln("ProposalResp sent")
	return nil
}
