package blst

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrivateKeyMarshal(t *testing.T) {
	assert := assert.New(t)
	p := GenerateKey()
	p1 := new(PrivateKey)
	{
		err := p1.FromBytes(p.Bytes())
		assert.Nil(err, "PrivateKey.FromBytes")
		assert.Equal(p.String(), p1.String(), "p != p1")
	}

	P := p.PublicKey()
	P1 := new(PublicKey)
	{
		err := P1.FromBytes(P.Bytes())
		assert.Nil(err, "PublicKey.FromBytes")
		assert.Equal(P.String(), P1.String(), "P != P1")
	}

	msg := []byte("just a test")
	s := p.Sign(msg)
	s1 := new(Signature)
	{
		err := s1.FromBytes(s.Bytes())
		assert.Nil(err, "Signature.FromBytes")
		assert.Equal(s.String(), s1.String(), "s != s1")
	}

	assert.True(P.Verify(msg, s1), "P verify signature s1 failed")
	assert.True(P1.Verify(msg, s), "P1 verify signature s failed")
}

func TestAggregatePublicKeys(t *testing.T) {
	assert := assert.New(t)

	var (
		privs []*PrivateKey
		pubs  []*PublicKey
		sigs  []*Signature
		msg   = []byte("just a test")
	)

	for i := 0; i < 20; i++ {
		var (
			p = GenerateKey()
			P = p.PublicKey()
			s = p.Sign(msg)
		)
		assert.True(P.Verify(msg, s), "verify signature failed")

		privs = append(privs, p)
		pubs = append(pubs, P)
		sigs = append(sigs, s)
	}
	pub := AggregatePublicKeys(pubs)
	sig := AggregateSignatures(sigs)
	assert.True(pub.Verify(msg, sig), "AggregatePublicKeys verify AggregateSignatures failed")
}
