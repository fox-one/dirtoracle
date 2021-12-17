package encoding

import (
	"bytes"
	"encoding/binary"

	"github.com/gofrs/uuid"
)

const (
	OperationPurposeUnknown       = 0
	OperationPurposeGroupEvent    = 1
	OperationPurposeAddProcess    = 11
	OperationPurposeCreditProcess = 12
)

type Operation struct {
	Purpose  int
	Process  string
	Platform string
	Address  string
	Extra    []byte
}

func (o *Operation) Encode() []byte {
	enc := &bytes.Buffer{}
	writeInt(enc, o.Purpose)
	writeUUID(enc, o.Process)
	writeBytes(enc, []byte(o.Platform))
	writeBytes(enc, []byte(o.Address))
	writeBytes(enc, o.Extra)
	return enc.Bytes()
}

func writeInt(enc *bytes.Buffer, d int) {
	b := uint16ToByte(uint16(d))
	enc.Write(b)
}

func writeUUID(enc *bytes.Buffer, id string) {
	uid, err := uuid.FromString(id)
	if err != nil {
		panic(err)
	}
	enc.Write(uid.Bytes())
}

func writeBytes(enc *bytes.Buffer, b []byte) {
	if len(b) > 128 {
		panic(b)
	}
	writeInt(enc, len(b))
	enc.Write(b)
}

func uint16ToByte(d uint16) []byte {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, d)
	return b
}
