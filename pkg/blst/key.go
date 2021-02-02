package blst

import (
	blst "github.com/supranational/blst/bindings/go"
)

var (
	dst = []byte("BLS_SIG_BLS12381G1_XMD:SHA-256_SSWU_RO_NUL_")
)

type (
	// For minimal-signature-size operations
	PrivateKey blst.SecretKey
	PublicKey  blst.P2Affine
	Signature  blst.P1Affine
)
