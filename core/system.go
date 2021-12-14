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

func (s *System) SignProposal(p *ProposalRequest, index uint64) (*ProposalResp, error) {
	payload := p.Payload()
	sig := s.SignKey.Sign(payload)
	en256Sig, err := s.En256SignKey.Sign(payload)
	if err != nil {
		return nil, err
	}

	return &ProposalResp{
		TraceID:        p.TraceID,
		Index:          index,
		Signature:      sig,
		En256Signature: en256Sig,
	}, nil
}

func (s *System) VerifyData(req *PriceRequest, p *PriceData) bool {
	{
		var pubs []*blst.PublicKey
		for _, signer := range req.Signers {
			if p.Signature.Mask&(0x1<<signer.Index) != 0 {
				pubs = append(pubs, signer.VerifyKey)
			}
		}

		if len(pubs) < int(req.Threshold) ||
			!blst.AggregatePublicKeys(pubs).Verify(p.Payload(), &p.Signature.Signature) {
			return false
		}
	}

	if p.En256Signature != nil {
		var pubs []*en256.PublicKey
		for _, signer := range req.Signers {
			if p.En256Signature.Mask&(0x1<<signer.Index) != 0 {
				pubs = append(pubs, signer.En256VerifyKey)
			}
		}

		if len(pubs) < int(req.Threshold) ||
			!en256.AggregatePublicKeys(pubs).Verify(p.Payload(), &p.En256Signature.Signature) {
			return false
		}
	}

	return true
}
