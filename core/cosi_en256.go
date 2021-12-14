package core

import (
	"database/sql/driver"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strconv"

	"github.com/fox-one/dirtoracle/pkg/mtg"
	"github.com/pandodao/blst/en256"
)

type (
	En256CosiSignature struct {
		en256.Signature
		Mask uint64
	}
)

func (s *En256CosiSignature) Bytes() []byte {
	bts, err := mtg.Encode(s.Mask, &s.Signature)
	if err != nil {
		panic(err)
	}
	return bts
}

func (s *En256CosiSignature) FromBytes(bts []byte) error {
	var mask uint64
	left, err := mtg.Scan(bts, &mask)
	if err != nil {
		return err
	}

	if len(left) < 65 || left[0] != 64 {
		return errors.New("empty signature")
	}

	var sig en256.Signature
	if err := sig.FromBytes(left[1:65]); err != nil {
		return err
	}
	s.Mask, s.Signature = mask, sig
	return nil
}

func (s *En256CosiSignature) String() string {
	return base64.StdEncoding.EncodeToString(s.Bytes())
}

func (s *En256CosiSignature) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(s.String())), nil
}

func (s *En256CosiSignature) UnmarshalJSON(b []byte) error {
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

func (s *En256CosiSignature) MarshalBinary() (data []byte, err error) {
	return s.Bytes(), nil
}

func (s *En256CosiSignature) UnmarshalBinary(data []byte) error {
	return s.FromBytes(data)
}

// Scan implements the sql.Scanner interface for database deserialization.
func (s *En256CosiSignature) Scan(value interface{}) error {
	var d []byte
	switch v := value.(type) {
	case string:
		d = []byte(v)
	case []byte:
		d = v
	}
	var sig En256CosiSignature
	if err := json.Unmarshal(d, &sig); err != nil {
		return err
	}
	*s = sig
	return nil
}

// Value implements the driver.Valuer interface for database serialization.
func (s *En256CosiSignature) Value() (driver.Value, error) {
	return s.MarshalJSON()
}
