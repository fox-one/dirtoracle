package core

import (
	"github.com/fox-one/dirtoracle/pkg/blst"
	"github.com/shopspring/decimal"
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
		GasAsset       string
		GasAmount      decimal.Decimal
	}
)

func (m *Member) Mask() uint64 {
	return 0x1 << m.ID
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

	var (
		sigMap = map[int64]*blst.Signature{}
		mask   uint64
	)

	for id, s := range p0.Signatures {
		sigMap[id] = s
		mask = mask | (0x1 << id)
	}
	for id, s := range p1.Signatures {
		sigMap[id] = s
		mask = mask | (0x1 << id)
	}

	if len(sigMap) >= int(s.Threshold) {
		p.Signature = &CosiSignature{
			Mask: mask,
		}
		sigs := make([]*blst.Signature, 0, len(sigMap))
		for _, s := range sigMap {
			sigs = append(sigs, s)
		}
		p.Signature.Signature = *blst.AggregateSignatures(sigs)
		p.Signatures = nil
	} else {
		p.Signatures = sigMap
	}
	return p
}

func (s *System) SignProposal(p *PriceProposal) *PriceProposal {
	me := s.Me()
	p1 := &PriceProposal{
		PriceData: p.PriceData,
		Signatures: map[int64]*blst.Signature{
			me.ID: s.SignKey.Sign(p.Payload()),
		},
	}
	if len(p.Signatures) == 0 {
		return p1
	}
	return s.MergeProposals(p, p1)
}

func (s *System) VerifyData(p *PriceData) bool {
	var pubs []*blst.PublicKey
	for _, m := range s.Members {
		if p.Signature.Mask&(0x1<<m.ID) != 0 {
			pubs = append(pubs, m.VerifyKey)
		}
	}

	return len(pubs) >= int(s.Threshold) &&
		blst.AggregatePublicKeys(pubs).Verify(p.Payload(), &p.Signature.Signature)
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
