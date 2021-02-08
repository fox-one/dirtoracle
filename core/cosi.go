package core

import (
	"encoding/base64"
	"strconv"

	"github.com/fox-one/dirtoracle/pkg/blst"
	"github.com/fox-one/dirtoracle/pkg/mtg"
)

type (
	CosiSignature struct {
		blst.Signature
		Mask uint64
	}
)

func (s *CosiSignature) Bytes() []byte {
	bts, err := mtg.Encode(s.Mask, &s.Signature)
	if err != nil {
		panic(err)
	}
	return bts
}

func (s *CosiSignature) FromBytes(bts []byte) error {
	var sig CosiSignature
	_, err := mtg.Scan(bts, &sig.Mask, &sig.Signature)
	if err != nil {
		return err
	}
	s.Mask, s.Signature = sig.Mask, sig.Signature
	return nil
}

func (s *CosiSignature) String() string {
	return base64.StdEncoding.EncodeToString(s.Bytes())
}

func (s *CosiSignature) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(s.String())), nil
}

func (s *CosiSignature) UnmarshalJSON(b []byte) error {
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

func (s *CosiSignature) MarshalBinary() (data []byte, err error) {
	return s.Bytes(), nil
}

func (s *CosiSignature) UnmarshalBinary(data []byte) error {
	return s.FromBytes(data)
}
