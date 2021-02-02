package blst

import (
	"encoding/base64"
	"fmt"
	"strconv"

	blst "github.com/supranational/blst/bindings/go"
)

func (s *Signature) Bytes() []byte {
	return (*blst.P1Affine)(s).Compress()
}

func (s *Signature) FromBytes(bts []byte) error {
	secret := new(blst.P1Affine).Uncompress(bts)
	if secret == nil {
		return fmt.Errorf("invalid blst public key")
	}

	*s = (Signature)(*secret)
	return nil
}

func (s *Signature) String() string {
	return base64.StdEncoding.EncodeToString(s.Bytes())
}

func (s *Signature) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(s.String())), nil
}

func (s *Signature) UnmarshalJSON(b []byte) error {
	unquoted, err := strconv.Unquote(string(b))
	if err != nil {
		return err
	}

	bts, err := base64.StdEncoding.DecodeString(unquoted)
	if err != nil {
		return err
	}

	return s.FromBytes(bts)
}

func (s *Signature) MarshalBinary() (data []byte, err error) {
	return s.Bytes(), nil
}

func (s *Signature) UnmarshalBinary(data []byte) error {
	return s.FromBytes(data)
}

func AggregateSignatures(sigs []*Signature) *Signature {
	agSig := new(blst.P1Aggregate)
	for _, s := range sigs {
		agSig.Add((*blst.P1Affine)(s), false)
	}
	return (*Signature)(agSig.ToAffine())
}
