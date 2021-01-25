package core

import (
	blst "github.com/supranational/blst/bindings/go"
)

type Member struct {
	ClientID  string
	Name      string
	VerifyKey blst.P1Affine
}
