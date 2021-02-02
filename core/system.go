package core

import (
	"github.com/fox-one/dirtoracle/pkg/blst"
)

type (
	Member struct {
		ID        int64           `json:"id"`
		ClientID  string          `json:"client_id"`
		Name      string          `json:"name"`
		VerifyKey *blst.PublicKey `json:"verify_key"`
	}

	// System stores system information.
	System struct {
		Admins         []string
		ConversationID string
		ClientID       string
		Members        []*Member
		Threshold      uint8
		SignKey        *blst.PrivateKey
		Version        string
	}
)

func (m *Member) Mask() int64 {
	return 0x1 << m.ID
}

func (m *Member) VerifyProposal(p *PriceData) bool {
	return m.VerifyKey.Verify(p.Payload(), p.Signature)
}

func (s *System) Me() *Member {
	for _, m := range s.Members {
		if m.ClientID == s.ClientID {
			return m
		}
	}
	return nil
}

func (s *System) Member(id int64) *Member {
	for _, m := range s.Members {
		if m.ID == id {
			return m
		}
	}
	return nil
}

func (s *System) MemberIDs() []string {
	ids := make([]string, len(s.Members))
	for idx, m := range s.Members {
		ids[idx] = m.ClientID
	}

	return ids
}

func (s *System) MergeProposals(p0, p1 *PriceProposal) *PriceProposal {
	p := &PriceProposal{
		PriceData: p0.PriceData,
	}
	p.Mask = (p0.Mask | p1.Mask)

	sigMap := map[int64]*blst.Signature{}
	for id, s := range p0.Signatures {
		sigMap[id] = s
	}
	for id, s := range p1.Signatures {
		sigMap[id] = s
	}

	if len(sigMap) >= int(s.Threshold) {
		sigs := make([]*blst.Signature, len(p.Signatures))
		for _, s := range sigMap {
			sigs = append(sigs, s)
		}
		p.Signature = blst.AggregateSignatures(sigs)
	} else {
		p.Signatures = sigMap
	}
	return p
}

func (s *System) SignProposal(p *PriceProposal) *PriceProposal {
	me := s.Me()
	p1 := &PriceProposal{
		Signatures: map[int64]*blst.Signature{
			me.ID: s.SignKey.Sign(p.Payload()),
		},
	}
	p1.Mask = me.Mask()
	return s.MergeProposals(p, p1)
}

func (s *System) VerifyData(p *PriceData) bool {
	var pubs []*blst.PublicKey
	for _, m := range s.Members {
		if p.Mask&(0x1<<m.ID) != 0 {
			pubs = append(pubs, m.VerifyKey)
		}
	}

	return len(pubs) >= int(s.Threshold) &&
		blst.AggregatePublicKeys(pubs).Verify(p.Payload(), p.Signature)
}

func (s *System) VerifyProposal(p *PriceProposal) bool {
	payload := p.Payload()
	for id, sig := range p.Signatures {
		if m := s.Member(id); m == nil || !m.VerifyKey.Verify(payload, sig) {
			return false
		}
	}

	return true
}
