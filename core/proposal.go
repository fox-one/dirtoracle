package core

import (
	"github.com/pandodao/blst"
	"github.com/pandodao/blst/en256"
	"github.com/shopspring/decimal"
)

type (
	Proposal struct {
		PriceRequest
		ProposalRequest `json:"-"`

		Signatures      map[uint64]*blst.Signature  `json:"sigs,omitempty"`
		En256Signatures map[uint64]*en256.Signature `json:"en256_sigs,omitempty"`
	}

	ProposalRequest struct {
		Asset

		TraceID   string          `json:"trace_id,omitempty"`
		Timestamp int64           `json:"timestamp,omitempty"`
		Price     decimal.Decimal `json:"price,omitempty"`
		Signers   []*Signer       `json:"signers,omitempty"`
		Signature *ProposalResp   `json:"signature,omitempty"`
	}

	ProposalResp struct {
		TraceID        string           `json:"trace_id,omitempty"`
		Index          uint64           `json:"index"`
		Signature      *blst.Signature  `json:"signature,omitempty"`
		En256Signature *en256.Signature `json:"en256_signature,omitempty"`
	}
)

func (p Proposal) Export() *PriceData {
	var (
		cosi      *CosiSignature
		cosiEn256 *En256CosiSignature
	)

	if len(p.Signatures) >= int(p.Threshold) {
		cosi = new(CosiSignature)
		var sigs = make([]*blst.Signature, 0, len(p.Signatures))
		for id, sig := range p.Signatures {
			cosi.Mask = cosi.Mask | (1 << id)
			sigs = append(sigs, sig)
		}
		cosi.Signature = *blst.AggregateSignatures(sigs)
	}

	if len(p.En256Signatures) >= int(p.Threshold) {
		cosiEn256 = new(En256CosiSignature)
		var sigs = make([]*en256.Signature, 0, len(p.En256Signatures))
		for id, sig := range p.En256Signatures {
			cosiEn256.Mask = cosiEn256.Mask | (1 << id)
			sigs = append(sigs, sig)
		}
		cosiEn256.Signature = *en256.AggregateSignatures(sigs)
	}

	return &PriceData{
		AssetID:        p.PriceRequest.AssetID,
		Timestamp:      p.Timestamp,
		Price:          p.Price,
		Signature:      cosi,
		En256Signature: cosiEn256,
	}
}

func (p ProposalRequest) Payload() []byte {
	return PriceData{
		Timestamp: p.Timestamp,
		AssetID:   p.AssetID,
		Price:     p.Price,
	}.Payload()
}

func (p ProposalRequest) PayloadV1() ([]byte, error) {
	return PriceData{
		Timestamp: p.Timestamp,
		AssetID:   p.AssetID,
		Price:     p.Price,
	}.PayloadV1()
}

func (p ProposalRequest) Verify(resp *ProposalResp) bool {
	for _, signer := range p.Signers {
		if signer.Index == resp.Index {
			if resp.Signature != nil {
				return signer.VerifyKey.Verify(p.Payload(), resp.Signature)
			}
			if signer.En256VerifyKey != nil && resp.En256Signature != nil {
				if payload, err := p.PayloadV1(); err == nil {
					return signer.En256VerifyKey.Verify(payload, resp.En256Signature)
				}
			}
			return false
		}
	}
	return false
}
