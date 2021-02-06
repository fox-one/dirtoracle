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
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

func (m *Oracle) loopProposals(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case p := <-m.proposals:
			m.execWithTimeout(ctx, time.Second*20, func() error {
				return m.handlePriceProposal(ctx, p)
			})
		}
	}
}

func (m *Oracle) handlePriceProposal(ctx context.Context, p *core.PriceProposal) error {
	log := logger.FromContext(ctx).WithFields(logrus.Fields{
		"timestamp": p.Timestamp,
		"asset_id":  p.AssetID,
		"price":     p.Price,
	})
	ctx = logger.WithContext(ctx, log)

	// 全新 proposal
	if len(p.Signatures) == 0 {
		return m.handleNewPriceProposal(ctx, p)
	}

	// 非法 proposal，直接遗弃
	if !m.validatePriceProposal(ctx, p) {
		return nil
	}

	// 与历史 proposal 合并
	if p1 := m.cachedProposal(p); p1 != nil {
		// 之前已发送
		if p1.Signature != nil {
			return nil
		}

		if p.Signature != nil {
			// 若收到的 proposal 未集齐签名，则与本地记录合并签名
			p = m.system.MergeProposals(p1, p)
		}
	}

	// 已收集签名不足
	if p.Signature == nil {
		// 若已签过名，可以直接退出
		if _, ok := p.Signatures[m.me.ID]; ok {
			return nil
		}
		p = m.system.SignProposal(p)
		// 第一次签名，将签名后的 proposal 同步给其他节点
		if err := m.sendPriceProposal(ctx, p); err != nil {
			return err
		}
	}

	if !m.system.VerifyData(&p.PriceData) {
		log.WithField("signature", p.Signature).Errorln("Verify PriceData failed")
		return nil
	}

	return m.sendPriceData(ctx, p)
}

func (m *Oracle) handleNewPriceProposal(ctx context.Context, p *core.PriceProposal) error {
	log := logger.FromContext(ctx)

	if lastProposal := m.latestProposal(p.AssetID); lastProposal != nil {
		change := p.Price.Sub(lastProposal.Price).Div(p.Price)
		timeDelta := p.Timestamp - lastProposal.Timestamp
		log = log.WithFields(logrus.Fields{
			"time_delta": timeDelta,
			"change":     change,
			"price":      p.Price,
		})
		if timeDelta < 0 {
			log.Debugln("newer proposal has been sent")
			return nil
		}

		if change.Abs().LessThan(m.config.PriceChangeThreshold) &&
			timeDelta < m.config.MaxInterval.Milliseconds() {
			log.Debugln("price diff too small")
			return nil
		}
	}

	p = m.system.SignProposal(p)
	return m.sendPriceProposal(ctx, p)
}

func (m *Oracle) validatePriceProposal(ctx context.Context, p *core.PriceProposal) bool {
	log := logger.FromContext(ctx)

	// validate if price is newest
	{
		if p.Timestamp > time.Now().Unix()*1000 {
			log.Debugln("ignore proposal from the future")
			return false
		}

		ticker, err := m.markets.AggregateTickers(ctx, p.AssetID)
		if err != nil || !ticker.Price.IsPositive() {
			log.WithError(err).Errorln("AggregateTickers failed")
			return false
		}

		if p.Timestamp < ticker.Timestamp-m.config.MaxInterval.Milliseconds() {
			log.Debugln("ignore old proposal")
			return false
		}

		if p.Price.Div(ticker.Price).Sub(decimal.New(1, 0)).Abs().GreaterThan(m.config.PriceChangeThreshold) {
			log.Debugln("price diff too large")
			return false
		}
	}

	if !m.system.VerifyProposal(p) {
		log.WithField("signatures", p.Signatures).Errorln("Verify PriceProposal failed")
		return false
	}

	return true
}

func (m *Oracle) validatePriceData(ctx context.Context, p *core.PriceProposal) bool {
	log := logger.FromContext(ctx)
	if v := m.latestPriceData(p.AssetID); v != nil {
		change := p.Price.Sub(v.Price).Div(p.Price)
		timeDelta := p.Timestamp - v.Timestamp
		log = log.WithFields(logrus.Fields{
			"time_delta": timeDelta,
			"change":     change,
			"price":      p.Price,
		})

		if timeDelta < 0 {
			log.Debugln("newer price data has been sent")
			return false
		}

		if change.Abs().LessThan(m.config.PriceChangeThreshold) &&
			timeDelta < m.config.MaxInterval.Milliseconds() {

			log.Debugln("price diff too small")
			return false
		}
	}
	log.Infoln("price data")
	return true
}

func (m *Oracle) sendPriceProposal(ctx context.Context, p *core.PriceProposal) error {
	log := logger.FromContext(ctx).WithFields(logrus.Fields{
		"method":    "sendPriceProposal",
		"timestamp": p.Timestamp,
	})
	ctx = logger.WithContext(ctx, log)

	// if !m.validatePriceData(ctx, p) {
	// 	return nil
	// }

	bts, _ := json.MarshalIndent(p, "", "    ")
	msg := &mixin.MessageRequest{
		ConversationID: m.system.ConversationID,
		MessageID:      uuid.New(),
		Category:       mixin.MessageCategoryPlainPost,
		Data:           base64.StdEncoding.EncodeToString(bts),
	}
	if err := m.client.SendMessage(ctx, msg); err != nil {
		logger.FromContext(ctx).WithError(err).Errorln("SendMessage failed")
		return err
	}

	log.WithField("mid", msg.MessageID).Debugln("SendMessage")
	m.cacheProposal(p)
	return nil
}

func (m *Oracle) sendPriceData(ctx context.Context, p *core.PriceProposal) error {
	log := logger.FromContext(ctx)

	if !m.validatePriceData(ctx, p) {
		return nil
	}

	// asset, _ := uuid.FromString(p.AssetID)
	// bts, _ := mtg.Encode(p.Timestamp, asset, p.Price, p.Mask, p.Signature)

	{
		bts, _ := json.MarshalIndent(p, "", "    ")
		u := "170e40f0-627f-4af2-acf5-0f25c009e523"
		c := mixin.UniqueConversationID(m.system.ClientID, u)
		if err := m.client.SendMessage(ctx, &mixin.MessageRequest{
			ConversationID: c,
			RecipientID:    u,
			MessageID:      uuid.MD5(string(p.Payload())),
			Category:       mixin.MessageCategoryPlainPost,
			Data:           base64.StdEncoding.EncodeToString(bts),
		}); err != nil {
			logger.FromContext(ctx).WithError(err).Errorln("SendMessage failed")
			return err
		}
	}

	log.Infoln("cachePriceData")
	m.cachePriceData(p)
	return m.sendPriceProposal(ctx, p)
}
