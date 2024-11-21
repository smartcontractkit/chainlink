package clientdss

import (
	"crypto/rand"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/services/signatures/cryptotest"
	"github.com/smartcontractkit/chainlink/v2/core/services/signatures/ethschnorr"
	"github.com/smartcontractkit/chainlink/v2/core/services/signatures/secp256k1"

	"go.dedis.ch/kyber/v3"
	dkg "go.dedis.ch/kyber/v3/share/dkg/rabin"
)

var suite = secp256k1.NewBlakeKeccackSecp256k1()

var nbParticipants = 7
var t = nbParticipants/2 + 1

var partPubs []kyber.Point
var partSec []kyber.Scalar

var longterms []*dkg.DistKeyShare
var randoms []*dkg.DistKeyShare

var msg *big.Int

var randomStream = cryptotest.NewStream(&testing.T{}, 0)

func init() {
	partPubs = make([]kyber.Point, nbParticipants)
	partSec = make([]kyber.Scalar, nbParticipants)
	for i := 0; i < nbParticipants; i++ {
		kp := secp256k1.Generate(randomStream)
		partPubs[i] = kp.Public
		partSec[i] = kp.Private
	}
	// Corresponds to section 4.2, step 1 of Stinson, 2001 paper
	longterms = genDistSecret(true) // Keep trying until valid public key
	randoms = genDistSecret(false)

	var err error
	msg, err = rand.Int(rand.Reader, big.NewInt(0).Lsh(big.NewInt(1), 256))
	if err != nil {
		panic(err)
	}
}

func TestDSSNew(t *testing.T) {
	dssArgs := DSSArgs{secret: partSec[0], participants: partPubs,
		long: longterms[0], random: randoms[0], msg: msg, T: 4}
	dss, err := NewDSS(dssArgs)
	assert.NotNil(t, dss)
	assert.Nil(t, err)
	dssArgs.secret = suite.Scalar().Zero()
	dss, err = NewDSS(dssArgs)
	assert.Nil(t, dss)
	assert.Error(t, err)
}

func TestDSSPartialSigs(t *testing.T) {
	dss0 := getDSS(0)
	dss1 := getDSS(1)
	ps0, err := dss0.PartialSig()
	assert.Nil(t, err)
	assert.NotNil(t, ps0)
	assert.Len(t, dss0.partials, 1)
	// second time should not affect list
	ps0, err = dss0.PartialSig()
	assert.Nil(t, err)
	assert.NotNil(t, ps0)
	assert.Len(t, dss0.partials, 1)

	// wrong index
	goodI := ps0.Partial.I
	ps0.Partial.I = 100
	err = dss1.ProcessPartialSig(ps0)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid index")
	ps0.Partial.I = goodI

	// wrong sessionID
	goodSessionID := ps0.SessionID
	ps0.SessionID = []byte("ahhhhhhhhhhhhhhhhhhhhhhhhhhhhhhh")
	err = dss1.ProcessPartialSig(ps0)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "dss: session id")
	ps0.SessionID = goodSessionID

	// wrong Signature
	goodSig := ps0.Signature
	ps0.Signature = ethschnorr.NewSignature()
	copy(ps0.Signature.CommitmentPublicAddress[:], randomBytes(20))
	badSig := secp256k1.ToInt(suite.Scalar().Pick(randomStream))
	ps0.Signature.Signature.Set(badSig)
	assert.Error(t, dss1.ProcessPartialSig(ps0))
	ps0.Signature = goodSig

	// invalid partial sig
	goodV := ps0.Partial.V
	ps0.Partial.V = suite.Scalar().Zero()
	ps0.Signature, err = ethschnorr.Sign(dss0.secret, ps0.Hash())
	require.Nil(t, err)
	err = dss1.ProcessPartialSig(ps0)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not valid")
	ps0.Partial.V = goodV
	ps0.Signature = goodSig

	// fine
	err = dss1.ProcessPartialSig(ps0)
	assert.Nil(t, err)

	// already received
	assert.Error(t, dss1.ProcessPartialSig(ps0))

	// if not enough partial signatures, can't generate signature
	sig, err := dss1.Signature()
	assert.Nil(t, sig) // XXX: Should also check err is nil?
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not enough")

	// enough partial sigs ?
	for i := 2; i < nbParticipants; i++ {
		dss := getDSS(i)
		ps, e := dss.PartialSig()
		require.Nil(t, e)
		require.Nil(t, dss1.ProcessPartialSig(ps))
	}
	assert.True(t, dss1.EnoughPartialSig())
	sig, err = dss1.Signature()
	assert.NoError(t, err)
	assert.NoError(t, Verify(dss1.long.Commitments()[0], msg, sig))
}

var printTests = false

func printTest(t *testing.T, msg *big.Int, public kyber.Point,
	signature ethschnorr.Signature) {
	pX, pY := secp256k1.Coordinates(public)
	t.Logf("  ['%064x',\n   '%064x',\n   '%064x',\n   '%064x',\n   '%040x'],\n",
		msg, pX, pY, signature.Signature,
		signature.CommitmentPublicAddress)
}

func TestDSSSignature(t *testing.T) {
	dsss := make([]*DSS, nbParticipants)
	pss := make([]*PartialSig, nbParticipants)
	for i := 0; i < nbParticipants; i++ {
		dsss[i] = getDSS(i)
		ps, err := dsss[i].PartialSig()
		require.Nil(t, err)
		require.NotNil(t, ps)
		pss[i] = ps
	}
	for i, dss := range dsss {
		for j, ps := range pss {
			if i == j {
				continue
			}
			require.Nil(t, dss.ProcessPartialSig(ps))
		}
	}
	// issue and verify signature
	dss0 := dsss[0]
	sig, err := dss0.Signature()
	assert.NotNil(t, sig)
	assert.Nil(t, err)
	assert.NoError(t, ethschnorr.Verify(longterms[0].Public(), dss0.msg, sig))
	// Original contains this second check. Unclear why.
	assert.NoError(t, ethschnorr.Verify(longterms[0].Public(), dss0.msg, sig))
	if printTests {
		printTest(t, dss0.msg, dss0.long.Commitments()[0], sig)
	}
}

func TestPartialSig_Hash(t *testing.T) {
	observedHashes := make(map[*big.Int]bool)
	for i := 0; i < nbParticipants; i++ {
		psig, err := getDSS(i).PartialSig()
		require.NoError(t, err)
		hash := psig.Hash()
		require.False(t, observedHashes[hash])
		observedHashes[hash] = true
	}
}

func getDSS(i int) *DSS {
	dss, err := NewDSS(DSSArgs{secret: partSec[i], participants: partPubs,
		long: longterms[i], random: randoms[i], msg: msg, T: t})
	if dss == nil || err != nil {
		panic("nil dss")
	}
	return dss
}

func _genDistSecret() []*dkg.DistKeyShare {
	dkgs := make([]*dkg.DistKeyGenerator, nbParticipants)
	for i := 0; i < nbParticipants; i++ {
		dkg, err := dkg.NewDistKeyGenerator(suite, partSec[i], partPubs, nbParticipants/2+1)
		if err != nil {
			panic(err)
		}
		dkgs[i] = dkg
	}
	// full secret sharing exchange
	// 1. broadcast deals
	resps := make([]*dkg.Response, 0, nbParticipants*nbParticipants)
	for _, dkg := range dkgs {
		deals, err := dkg.Deals()
		if err != nil {
			panic(err)
		}
		for i, d := range deals {
			resp, err := dkgs[i].ProcessDeal(d)
			if err != nil {
				panic(err)
			}
			if !resp.Response.Approved {
				panic("wrong approval")
			}
			resps = append(resps, resp)
		}
	}
	// 2. Broadcast responses
	for _, resp := range resps {
		for h, dkg := range dkgs {
			// ignore all messages from ourself
			if resp.Response.Index == uint32(h) {
				continue
			}
			j, err := dkg.ProcessResponse(resp)
			if err != nil || j != nil {
				panic("wrongProcessResponse")
			}
		}
	}
	// 4. Broadcast secret commitment
	for i, dkg := range dkgs {
		scs, err := dkg.SecretCommits()
		if err != nil {
			panic("wrong SecretCommits")
		}
		for j, dkg2 := range dkgs {
			if i == j {
				continue
			}
			cc, err := dkg2.ProcessSecretCommits(scs)
			if err != nil || cc != nil {
				panic("wrong ProcessSecretCommits")
			}
		}
	}

	// 5. reveal shares
	dkss := make([]*dkg.DistKeyShare, len(dkgs))
	for i, dkg := range dkgs {
		dks, err := dkg.DistKeyShare()
		if err != nil {
			panic(err)
		}
		dkss[i] = dks
	}
	return dkss
}

func genDistSecret(checkValidPublicKey bool) []*dkg.DistKeyShare {
	rv := _genDistSecret()
	if checkValidPublicKey {
		// Because of the trick we're using to verify the signatures on-chain, we
		// need to make sure that the  ordinate of this distributed public key is
		// in the lower half of {0,...,}
		for !secp256k1.ValidPublicKey(rv[0].Public()) {
			rv = _genDistSecret() // Keep trying until valid distributed public key.
		}
	}
	return rv
}

func randomBytes(n int) []byte {
	var buff = make([]byte, n)
	_, _ = rand.Read(buff)
	return buff
}
