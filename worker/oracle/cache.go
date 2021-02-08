package oracle

import (
	"fmt"
	"time"

	"github.com/fox-one/dirtoracle/core"
)

func (m *Oracle) assetProposalKey(assetID string) string {
	return fmt.Sprintf("asset:proposals:%s", assetID)
}

func (m *Oracle) proposalKey(p *core.PriceProposal) string {
	return fmt.Sprintf("price_proposal:%s;%d;%s;", p.AssetID, p.Timestamp, p.Price)
}

func (m *Oracle) cacheProposal(p *core.PriceProposal) error {
	if v, ok := m.cache.Get(m.assetProposalKey(p.AssetID)); !ok || v.(*core.PriceProposal).Timestamp < p.Timestamp {
		m.cache.Set(m.assetProposalKey(p.AssetID), p, m.config.MaxInterval)
	}

	m.cache.Set(m.proposalKey(p), p, time.Minute*2)
	return nil
}

func (m *Oracle) cachedProposal(p *core.PriceProposal) *core.PriceProposal {
	if v, ok := m.cache.Get(m.proposalKey(p)); ok {
		return v.(*core.PriceProposal)
	}
	return nil
}

func (m *Oracle) latestProposal(assetID string) *core.PriceProposal {
	if v, ok := m.cache.Get(m.assetProposalKey(assetID)); ok {
		return v.(*core.PriceProposal)
	}
	return nil
}
