package blst

import (
	"encoding/base64"
	"fmt"

	blst "github.com/supranational/blst/bindings/go"
)

func EncodePrivateKey(key *blst.SecretKey) string {
	return base64.StdEncoding.EncodeToString(key.Serialize())
}

func DecodePrivateKey(s string) (*blst.SecretKey, error) {
	secret := &blst.SecretKey{}
	secret = secret.Deserialize(decodeBase64(s))
	if secret == nil {
		return nil, fmt.Errorf("invalid blst private key")
	}

	return secret, nil
}

func EncodePublicKey(key *blst.P1Affine) string {
	return base64.StdEncoding.EncodeToString(key.Serialize())
}

func DecodePublicKey(s string) (*blst.P1Affine, error) {
	key := &blst.P1Affine{}
	key = key.Deserialize(decodeBase64(s))
	if key == nil {
		return nil, fmt.Errorf("invalid blst public key")
	}

	return key, nil
}

func decodeBase64(data string) []byte {
	if b, err := base64.StdEncoding.DecodeString(data); err == nil {
		return b
	}

	if b, err := base64.URLEncoding.DecodeString(data); err == nil {
		return b
	}

	return []byte(data)
}
