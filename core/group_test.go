package core

import (
	"encoding/base64"
	"testing"

	"github.com/fox-one/dirtoracle/pkg/mtg"
	"github.com/stretchr/testify/require"
)

func TestDecodeMemberAction(t *testing.T) {
	memo := "QD+ML36DKnOAXqrJdcuoNdI3s8eZok8c+RJkgRFi3SjyWv/EQkRaw52tpeAaQdTN5NaLcnxNLFwBkVd4QmE5ABC4Vt6z6S9MGZcz7ENSb5XOENFctuRC30MopxxJxJ5kbf4BAg"
	data, _ := base64.StdEncoding.DecodeString(memo)

	body, sig, err := mtg.Unpack(data)
	require.Nil(t, err)

	t.Log(len(sig), len(body))
}
