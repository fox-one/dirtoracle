package oracle

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
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

	// 如果发送过更新数据 或者 短时间内价格变化不大，则可以直接跳过
	if lastProposal := m.latestProposal(p.AssetID); lastProposal != nil {
		if m.proposalKey(p) != m.proposalKey(lastProposal) {
			var (
				change    = p.Price.Sub(lastProposal.Price).Div(p.Price)
				timeDelta = p.Timestamp - lastProposal.Timestamp
				log       = log.WithFields(logrus.Fields{
					"time_delta": timeDelta,
					"change":     change,
					"price":      p.Price,
				})
			)
			if timeDelta < 0 {
				log.Debugln("newer proposal has been sent")
				return nil
			}

			if change.Abs().LessThan(m.config.PriceChangeThreshold) &&
				timeDelta < m.config.MaxInterval.Milliseconds() {
				log.Debugln("price diff too small")
				return nil
			}
		} else if lastProposal.Signature != nil {
			return nil
		}
	}

	// 全新 proposal
	if len(p.Signatures) == 0 {
		p = m.system.SignProposal(p)
		return m.sendPriceProposal(ctx, p)
	}

	// 非法 proposal，直接遗弃
	if !m.validatePriceProposal(ctx, p) {
		return nil
	}

	// 与历史 proposal 合并
	if p1 := m.cachedProposal(p); p1 != nil {
		// 之前已收集到足够签名，可直接退出
		if p1.Signature != nil {
			return nil
		}

		if p.Signature == nil {
			// 若收到的 proposal 未集齐签名，则与本地记录合并签名
			p = m.system.MergeProposals(p1, p)
		}
	}

	// 已收集签名不足
	if p.Signature == nil {
		// 若已签过名，可以直接退出
		if _, ok := p.Signatures[m.me.ID]; ok {
			m.cacheProposal(p)
			return nil
		}
		p = m.system.SignProposal(p)
		// 第一次签名，将签名后的 proposal 同步给其他节点
		if err := m.sendPriceProposal(ctx, p); err != nil {
			return err
		}

		m.cacheProposal(p)
		if p.Signature == nil {
			return nil
		}
	}
	if !m.system.VerifyData(&p.PriceData) {
		log.WithField("signature", p.Signature).Errorln("Verify PriceData failed")
		return nil
	}

	return m.sendPriceData(ctx, p)
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

func (m *Oracle) sendPriceProposal(ctx context.Context, p *core.PriceProposal) error {
	log := logger.FromContext(ctx).WithFields(logrus.Fields{
		"method":    "sendPriceProposal",
		"timestamp": p.Timestamp,
	})

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

	return nil
}

func (m *Oracle) sendPriceData(ctx context.Context, p *core.PriceProposal) error {
	log := logger.FromContext(ctx)

	feeders, err := m.feeders.FindFeeders(ctx, p.AssetID)
	if err != nil {
		log.WithError(err).Errorln("FindFeeders failed")
		return err
	}

	if len(feeders) == 0 {
		return nil
	}

	p.Signatures = nil
	memo, _ := json.Marshal(p)
	var ts = make([]*core.Transfer, len(feeders))
	trace := uuid.MD5(fmt.Sprintf("price_data:%s;%d;%v;", p.AssetID, p.Timestamp, p.Price))
	for i, f := range feeders {
		ts[i] = &core.Transfer{
			TraceID:   uuid.MD5(fmt.Sprintf("trace:%s;%d;%s;", trace, f.Threshold, strings.Join(f.Opponents, ";"))),
			AssetID:   m.system.GasAsset,
			Amount:    m.system.GasAmount,
			Memo:      string(memo),
			Threshold: f.Threshold,
			Opponents: f.Opponents,
		}
	}
	return m.wallets.CreateTransfers(ctx, ts)
}
