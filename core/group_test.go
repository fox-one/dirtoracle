package core

import (
	"encoding/base64"
	"fmt"
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

func TestPriceDataMarshal(t *testing.T) {
	memo := "BarexZsMELdkcgWgSzrTqqGo9ejmaJQI/////gIBXGxDATZAETWNrfLfexVSO/yOhblgkKWzzsim3Sz2eahnvEZLJJ8coEfkknqFyEconF2EgPy9GbzMnu8xujmURo00yVseGQ=="
	data, _ := base64.StdEncoding.DecodeString(memo)

	var p PriceData
	if err := p.UnmarshalBinary(data); err != nil {
		panic(err)
	}

	fmt.Println(p)
}
