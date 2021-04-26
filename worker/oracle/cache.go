package oracle

import (
	"fmt"
	"time"

	"github.com/fox-one/dirtoracle/core"
)

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
