package secp256k1

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
)

var numPointSamples = 10

var randomStreamPoint = cltest.NewStream(&testing.T{}, 0)

func TestPoint_CloneAndEqual(t *testing.T) {
	f := newPoint()
	for i := 0; i < numPointSamples; i++ {
		g := f.Clone()
		f.Pick(randomStreamPoint)
		assert.NotEqual(t, f, g,
			"modifying original shouldn't change clone")
		g, h := f.Clone(), f.Clone()
		assert.Equal(t, f, g, "clones should be equal")
		g.Add(g, f)
		assert.Equal(t, h, f,
			"modifying a clone shouldn't change originial")
	}
}

func TestPoint_NullAndAdd(t *testing.T) {
	f, g := newPoint(), newPoint()
	for i := 0; i < numPointSamples; i++ {
		g.Null()
		f.Pick(randomStreamPoint)
		g.Add(f, g)
		assert.Equal(t, f, g, "adding zero should have no effect")
	}
}

func TestPoint_Set(t *testing.T) {
	p := newPoint()
	base := newPoint().Base()
	assert.NotEqual(t, p, base, "generator should not be zero")
	p.Set(base)
	assert.Equal(t, p, base, "setting to generator should yield generator")
}

func TestPoint_Embed(t *testing.T) {
	p := newPoint()
	for i := 0; i < numPointSamples; i++ {
		data := make([]byte, p.EmbedLen())
		_, err := rand.Read(data)
		require.Nil(t, err)
		p.Embed(data, randomStreamPoint)
		require.True(t, s256.IsOnCurve(p.X.int(), p.Y.int()),
			"should embed to a secp256k1 point")
		output, err := p.Data()
		require.True(t, bytes.Equal(data, output),
			"should get same value back after round-trip "+
				"embedding, got %v, then %v", data, output)
	}
}

func TestPoint_AddSubAndNeg(t *testing.T) {
	zero := newPoint().Null()
	p := newPoint()
	for i := 0; i < numPointSamples; i++ {
		p.Pick(randomStreamPoint)
		q := p.Clone()
		p.Sub(p, q)
		require.True(t, p.Equal(zero),
			"subtracting a point from itself should give zero, "+
				"got %v - %v = %v ≠ %v", q, q, p, zero)
		p.Neg(q)
		r := newPoint().Add(p, q)
		require.True(t, r.Equal(zero),
			"adding a point to its negative should give zero"+
				" got %v+%v=%v≠%v", q, p, r, zero)
		r.Neg(q)
		p.Sub(q, r)
		s := newPoint().Add(q, q)
		require.True(t, p.Equal(s), "q-(-q)=q+q?"+
			" got %v-%v=%v≠%v", q, r, p, s)
	}
}

func TestPoint_Mul(t *testing.T) {
	zero := newPoint().Null()
	multiplier := newScalar(bigZero)
	one := newScalar(big.NewInt(int64(1)))
	var p *secp256k1Point
	for i := 0; i < numPointSamples/5; i++ {
		if i%20 == 0 {
			p = nil // Test default to generator point
		} else {
			p := newPoint()
			p.Pick(randomStreamPoint)
		}
		multiplier.Pick(randomStreamPoint)
		q := newPoint().Mul(one, p)
		comparee := newPoint()
		if p == (*secp256k1Point)(nil) {
			comparee.Base()
		} else {
			comparee = p.Clone().(*secp256k1Point)
		}
		require.True(t, comparee.Equal(q), "1*p=p? %v * %v ≠ %v", one,
			comparee, q)
		q.Mul(multiplier, p)
		negMultiplier := newScalar(bigZero).Neg(multiplier)
		r := newPoint().Mul(negMultiplier, p)
		s := newPoint().Add(q, r)
		require.True(t, s.Equal(zero), "s*p+(-s)*p=0? got "+
			"%v*%v + %v*%v = %v + %v = %v ≠ %v", multiplier, p,
		)
	}
}

func TestPoint_Marshal(t *testing.T) {
	p := newPoint()
	for i := 0; i < numPointSamples; i++ {
		p.Pick(randomStreamPoint)
		serialized, err := p.MarshalBinary()
		require.Nil(t, err)
		q := newPoint()
		err = q.UnmarshalBinary(serialized)
		require.Nil(t, err)
		require.True(t, p.Equal(q), "%v marshalled to %x, which "+
			"unmarshalled to %v", p, serialized, q)
	}
}

func TestPoint_BaseTakesCopy(t *testing.T) {
	p := newPoint().Base()
	p.Add(p, p)
	q := newPoint().Base()
	assert.False(t, p.Equal(q),
		"modifying output from Base changes S256.G{x,y}")
}

func TestPoint_EthereumAddress(t *testing.T) {
	// Example taken from
	// https://theethereum.wiki/w/index.php/Accounts,_Addresses,_Public_And_Private_Keys,_And_Tokens
	pString := "3a1076bf45ab87712ad64ccb3b10217737f7faacbf2872e88fdd9a537d8fe266"
	pInt, ok := big.NewInt(0).SetString(pString, 16)
	require.True(t, ok, "failed to parse private key")
	private := newScalar(pInt)
	public := newPoint().Mul(private, nil)
	address, err := EthereumAddress(public)
	require.Nil(t, err)
	assert.Equal(t, fmt.Sprintf("%x", address),
		"c2d7cf95645d33006175b78989035c7c9061d3f9")
}
