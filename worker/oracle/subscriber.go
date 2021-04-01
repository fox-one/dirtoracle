package oracle

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/fox-one/dirtoracle/core"
	"github.com/fox-one/mixin-sdk-go"
	"github.com/fox-one/pkg/logger"
	"github.com/fox-one/pkg/uuid"
	"github.com/pandodao/blst"
)

func (m *Oracle) loopSubscribers(ctx context.Context) error {
	log := logger.FromContext(ctx).WithField("loop", "subscribers")
	ctx = logger.WithContext(ctx, log)

	var sleepDur = time.Millisecond
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case <-time.After(sleepDur):
			subscribers, err := m.subscribers.All(ctx)
			if err != nil {
				log.WithError(err).Errorln("subscribers.All")
				sleepDur = time.Second
				break
			}

			for _, s := range subscribers {
				s := s
				go m.execWithTimeout(ctx, time.Second*20, func() error {
					return m.handleSubscriber(ctx, s)
				})
				time.Sleep(time.Second)
			}
			sleepDur = time.Second * 10
		}
	}
}

func (m *Oracle) handleSubscriber(ctx context.Context, subscriber *core.Subscriber) error {
	log := logger.FromContext(ctx)

	var req *core.PriceRequest

	// TODO fetch price request from the subscriber's request url
	{
	}

	if req == nil {
		return nil
	}

	if p := m.cachedProposal(req.TraceID); p != nil {
		return nil
	}

	proposal := core.Proposal{
		PriceRequest: *req,
		Signatures:   map[uint64]*blst.Signature{},
	}

	// make Proposal
	{
		proposal.ProposalRequest = core.ProposalRequest{
			Asset:     req.Asset,
			TraceID:   req.TraceID,
			Signers:   req.Signers,
			Timestamp: time.Now().Unix(),
		}

		price, err := m.getPrice(ctx, &req.Asset)
		if err != nil || price.IsZero() {
			return err
		}
		proposal.ProposalRequest.Price = price

		var signer *core.Signer
		for _, s := range req.Signers {
			if s.VerifyKey.String() == m.system.VerifyKey.String() {
				signer = s
				break
			}
		}

		if signer == nil {
			log.Infoln("ignore:", "node not approved by the requester")
			return nil
		}

		resp := m.system.SignProposal(&proposal.ProposalRequest, signer.Index)
		proposal.ProposalRequest.Signature = resp
		proposal.Signatures[signer.Index] = resp.Signature
	}

	// send and cache proposal
	{
		if err := m.sendProposalRequest(ctx, &proposal.ProposalRequest); err != nil {
			return err
		}
		m.cacheProposal(&proposal)
	}
	return nil
}

func (m *Oracle) sendProposalRequest(ctx context.Context, p *core.ProposalRequest) error {
	log := logger.FromContext(ctx)

	bts, _ := json.MarshalIndent(p, "", "    ")
	msg := &mixin.MessageRequest{
		ConversationID: m.system.ConversationID,
		MessageID:      uuid.New(),
		Category:       mixin.MessageCategoryPlainPost,
		Data:           base64.StdEncoding.EncodeToString(bts),
	}
	if err := m.client.SendMessage(ctx, msg); err != nil {
		log.WithError(err).Errorln("SendMessage failed")
		return err
	}
	log.Infoln("Proposal sent")
	return nil
}
