package blst

import (
	"encoding/base64"
	"fmt"
	"strconv"

	blst "github.com/supranational/blst/bindings/go"
)

func (k *PublicKey) Verify(msg []byte, s *Signature) bool {
	return (*blst.P1Affine)(s).Verify(false, (*blst.P2Affine)(k), false, msg, dst)
}

func (k *PublicKey) Bytes() []byte {
	return (*blst.P2Affine)(k).Compress()
}

func (k *PublicKey) FromBytes(bts []byte) error {
	secret := new(blst.P2Affine).Uncompress(bts)
	if secret == nil {
		return fmt.Errorf("invalid blst public key")
	}

	*k = (PublicKey)(*secret)
	return nil
}

func (k *PublicKey) String() string {
	return base64.StdEncoding.EncodeToString(k.Bytes())
}

func (k *PublicKey) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(k.String())), nil
}

func (k *PublicKey) UnmarshalJSON(b []byte) error {
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

func AggregatePublicKeys(pubs []*PublicKey) *PublicKey {
	agPk := new(blst.P2Aggregate)
	for _, p := range pubs {
		agPk.Add((*blst.P2Affine)(p), false)
	}
	return (*PublicKey)(agPk.ToAffine())
}
