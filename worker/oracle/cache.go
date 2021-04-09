package oracle

import (
	"fmt"
	"time"

	"github.com/fox-one/dirtoracle/core"
)

func (m *Oracle) assetKey(id string) string {
	return fmt.Sprintf("asset:%s", id)
}

func (m *Oracle) proposalKey(trace string) string {
	return fmt.Sprintf("price_proposal:%s", trace)
}

func (m *Oracle) cacheProposal(p *core.Proposal) error {
	m.cache.Set(m.proposalKey(p.PriceRequest.TraceID), p, time.Minute*2)
	return nil
}

func (m *Oracle) cachedProposal(trace string) *core.Proposal {
	if v, ok := m.cache.Get(m.proposalKey(trace)); ok {
		return v.(*core.Proposal)
	}
	return nil
}

func (m *Oracle) cacheAssets(assets ...*core.Asset) error {
	for _, a := range assets {
		m.cache.Set(m.assetKey(a.AssetID), a, time.Hour)
	}
	return nil
}

func (m *Oracle) cachedAsset(id string) *core.Asset {
	if v, ok := m.cache.Get(m.assetKey(id)); ok {
		return v.(*core.Asset)
	}
	return nil
}
