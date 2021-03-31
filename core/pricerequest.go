package core

import (
	"github.com/pandodao/blst"
)

type (
	Signer struct {
		Index     uint64          `json:"index,omitempty"`
		VerifyKey *blst.PublicKey `json:"verify_key,omitempty"`
	}

	Receiver struct {
		Members   []string `json:"members,omitempty"`
		Threshold uint8    `json:"threshold"`
	}

	PriceRequest struct {
		Asset

		TraceID   string    `json:"trace_id,omitempty"`
		Receiver  *Receiver `json:"receiver,omitempty"`
		Signers   []*Signer `json:"signers,omitempty"`
		Threshold uint8     `json:"threshold"`
	}
)
