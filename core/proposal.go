package core

import (
	"github.com/pandodao/blst"
	"github.com/shopspring/decimal"
)

type (
	Proposal struct {
		PriceRequest
		ProposalRequest `json:"-"`

		Signatures map[uint64]*blst.Signature `json:"sigs,omitempty"`
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
		TraceID   string          `json:"trace_id,omitempty"`
		Index     uint64          `json:"index"`
		Signature *blst.Signature `json:"signature,omitempty"`
	}
)

func (p Proposal) Export() *PriceData {
	var (
		cosi CosiSignature
		sigs = make([]*blst.Signature, 0, len(p.Signatures))
	)
	for id, sig := range p.Signatures {
		cosi.Mask = cosi.Mask | (1 << id)
		sigs = append(sigs, sig)
	}
	cosi.Signature = *blst.AggregateSignatures(sigs)

	return &PriceData{
		AssetID:   p.PriceRequest.AssetID,
		Timestamp: p.Timestamp,
		Price:     p.Price,
		Signature: &cosi,
	}
}

func (p ProposalRequest) Payload() []byte {
	return PriceData{
		Timestamp: p.Timestamp,
		AssetID:   p.AssetID,
		Price:     p.Price,
	}.Payload()
}

func (p ProposalRequest) Verify(resp *ProposalResp) bool {
	for _, signer := range p.Signers {
		if signer.Index == resp.Index {
			return signer.VerifyKey.Verify(p.Payload(), resp.Signature)
		}
	}
	return false
}
