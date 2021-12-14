package oracle

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/fox-one/dirtoracle/core"
	"github.com/fox-one/mixin-sdk-go"
	"github.com/fox-one/pkg/logger"
	"github.com/fox-one/pkg/uuid"
	"github.com/sirupsen/logrus"
)

func (m *Oracle) handleProposalRespMessage(ctx context.Context, msg *mixin.MessageView) error {
	var resp = new(core.ProposalResp)
	log := logger.FromContext(ctx)
	data, _ := base64.StdEncoding.DecodeString(msg.Data)
	if err := json.Unmarshal(data, resp); err != nil {
		log.WithError(err).Errorln("Unmarshal ProposalResp failed")
		return nil
	}

	p := m.cachedProposal(resp.TraceID)
	if p == nil {
		log.Infoln("ignored:", "Proposal not found")
		return nil
	}

	if len(p.Signatures) == int(p.Threshold) {
		log.Infoln("ignored:", "Proposal already passed")
		return nil
	}

	if !p.Verify(resp) {
		log.Infoln("ignored:", "ProposalResp verify failed")
		return nil
	}

	p.Signatures[resp.Index] = resp.Signature
	if resp.En256Signature != nil {
		p.En256Signatures[resp.Index] = resp.En256Signature
	}
	if len(p.Signatures) == int(p.Threshold) {
		// create a final transaction to the receiver
		if err := m.sendPriceData(ctx, msg, p); err != nil {
			return err
		}
	}
	m.cacheProposal(p)
	return nil
}

func (m *Oracle) sendPriceData(ctx context.Context, msg *mixin.MessageView, p *core.Proposal) error {
	bts, _ := p.Export().MarshalBinary()
	memo := base64.StdEncoding.EncodeToString(bts)
	if err := m.wallets.CreateTransfers(ctx, []*core.Transfer{
		{
			TraceID:   uuid.MD5(fmt.Sprintf("price_data:trace:%s;proposal_message:%s", p.PriceRequest.TraceID, msg.QuoteMessageID)),
			AssetID:   m.system.GasAsset,
			Amount:    m.system.GasAmount,
			Memo:      memo,
			Threshold: p.Receiver.Threshold,
			Opponents: p.Receiver.Members,
		},
	}); err != nil {
		logger.FromContext(ctx).WithError(err).Errorln("CreateTransfers failed")
		return err
	}
	logger.FromContext(ctx).WithFields(logrus.Fields{
		"asset": p.ProposalRequest.Asset.Symbol,
		"price": p.Price,
	}).Infoln("PriceData sent")
	return nil
}
