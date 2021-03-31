package core

import (
	"github.com/pandodao/blst"
	"github.com/shopspring/decimal"
)

type (
	// System stores system information.
	System struct {
		ConversationID string
		GasAsset       string
		GasAmount      decimal.Decimal
		SignKey        *blst.PrivateKey
		VerifyKey      *blst.PublicKey
	}
)

func (s *System) SignProposal(p *ProposalRequest, index uint64) *ProposalResp {
	return &ProposalResp{
		TraceID:   p.TraceID,
		Index:     index,
		Signature: s.SignKey.Sign(p.Payload()),
	}
}

func (s *System) VerifyData(req *PriceRequest, p *PriceData) bool {
	var pubs []*blst.PublicKey
	for _, signer := range req.Signers {
		if p.Signature.Mask&(0x1<<signer.Index) != 0 {
			pubs = append(pubs, signer.VerifyKey)
		}
	}

	return len(pubs) >= int(req.Threshold) &&
		blst.AggregatePublicKeys(pubs).Verify(p.Payload(), &p.Signature.Signature)
}
