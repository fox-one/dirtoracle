package core

import (
	"github.com/pandodao/blst"
	"github.com/pandodao/blst/en256"
	"github.com/shopspring/decimal"
)

type (
	// System stores system information.
	System struct {
		Signers         []*Signer
		SignerThreshold uint8

		ClientID string

		MvmProcess   string
		MvmGroups    []string
		MvmThreshold uint8

		GasAsset  string
		GasAmount decimal.Decimal
	}
)

func (s *System) VerifyData(p *PriceData) bool {
	if p.Signature != nil {
		var pubs []*blst.PublicKey
		for _, signer := range s.Signers {
			if p.Signature.Mask&(0x1<<signer.Index) != 0 {
				pubs = append(pubs, signer.VerifyKey)
			}
		}

		return len(pubs) >= int(s.SignerThreshold) &&
			blst.AggregatePublicKeys(pubs).Verify(p.Payload(), &p.Signature.Signature)
	} else if p.En256Signature != nil {
		var pubs []*en256.PublicKey
		for _, signer := range s.Signers {
			if p.En256Signature.Mask&(0x1<<signer.Index) != 0 {
				pubs = append(pubs, signer.En256VerifyKey)
			}
		}

		payload, err := p.PayloadV1()
		return len(pubs) >= int(s.SignerThreshold) &&
			err == nil &&
			en256.AggregatePublicKeys(pubs).Verify(payload, &p.En256Signature.Signature)
	}

	return false
}
