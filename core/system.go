package core

import (
	"github.com/pandodao/blst"
	"github.com/pandodao/blst/en256"
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

		En256SignKey   *en256.PrivateKey
		En256VerifyKey *en256.PublicKey
	}
)

func (s *System) SignProposal(p *ProposalRequest, signer *Signer) (*ProposalResp, error) {
	resp := ProposalResp{
		TraceID: p.TraceID,
		Index:   signer.Index,
	}

	payload := p.Payload()

	if signer.VerifyKey != nil {
		resp.Signature = s.SignKey.Sign(payload)
	}

	if signer.En256VerifyKey != nil {
		sig, err := s.En256SignKey.Sign(payload)
		if err != nil {
			return nil, err
		}
		resp.En256Signature = sig
	}

	return &resp, nil
}

func (s *System) VerifyData(req *PriceRequest, p *PriceData) bool {
	if p.Signature != nil {
		var pubs []*blst.PublicKey
		for _, signer := range req.Signers {
			if p.Signature.Mask&(0x1<<signer.Index) != 0 {
				pubs = append(pubs, signer.VerifyKey)
			}
		}

		return len(pubs) >= int(req.Threshold) &&
			blst.AggregatePublicKeys(pubs).Verify(p.Payload(), &p.Signature.Signature)
	} else if p.En256Signature != nil {
		var pubs []*en256.PublicKey
		for _, signer := range req.Signers {
			if p.En256Signature.Mask&(0x1<<signer.Index) != 0 {
				pubs = append(pubs, signer.En256VerifyKey)
			}
		}

		return len(pubs) >= int(req.Threshold) &&
			en256.AggregatePublicKeys(pubs).Verify(p.Payload(), &p.En256Signature.Signature)
	}

	return false
}
