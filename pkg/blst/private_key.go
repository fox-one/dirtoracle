package blst

import (
	"encoding/base64"
	"fmt"
	"math/rand"
	"strconv"

	blst "github.com/supranational/blst/bindings/go"
)

func GenerateKey() *PrivateKey {
	var ikm = make([]byte, 32)
	rand.Read(ikm)
	return (*PrivateKey)(blst.KeyGen(ikm))
}

func (k *PrivateKey) Sign(msg []byte) *Signature {
	return (*Signature)(new(blst.P1Affine).Sign((*blst.SecretKey)(k), msg, dst))
}

func (k *PrivateKey) PublicKey() *PublicKey {
	pub := new(blst.P2Affine).From((*blst.SecretKey)(k))
	return (*PublicKey)(pub)
}

func (k *PrivateKey) Bytes() []byte {
	return (*blst.SecretKey)(k).Serialize()
}

func (k *PrivateKey) FromBytes(bts []byte) error {
	secret := new(blst.SecretKey).Deserialize(bts)
	if secret == nil {
		return fmt.Errorf("invalid blst private key")
	}

	*k = (PrivateKey)(*secret)
	return nil
}

func (k *PrivateKey) String() string {
	return base64.StdEncoding.EncodeToString(k.Bytes())
}

func (k *PrivateKey) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(k.String())), nil
}

func (k *PrivateKey) UnmarshalJSON(b []byte) error {
	unquoted, err := strconv.Unquote(string(b))
	if err != nil {
		return err
	}

	bts, err := base64.StdEncoding.DecodeString(unquoted)
	if err != nil {
		return err
	}

	return k.FromBytes(bts)
}
