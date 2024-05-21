package starkkey

import (
	cryptorand "crypto/rand"
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/NethermindEth/juno/core/felt"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// msg to hash
// [
//   '0x4acf99cb25a4803916f086440c661295b105a485efdc649ac4de9536da25b', // digest
//   1, // epoch_and_round
//   1, // extra_hash
//   1, // timestamp
//   '0x00010203000000000000000000000000000000000000000000000000000000', // observers
//   4, // len
//   99, // reports
//   99,
//   99,
//   99,
//   1 juels_per_fee_coin
// ]
// hash 0x1332a8dabaabef63b03438ca50760cb9f5c0292cbf015b2395e50e6157df4e3
// --> privKey 2137244795266879235401249500471353867704187908407744160927664772020405449078 r 2898571078985034687500959842265381508927681132188252715370774777831313601543 s 1930849708769648077928186998643944706551011476358007177069185543644456022504 pubKey 1118148281956858477519852250235501663092798578871088714409528077622994994907
//     privKey 3571531812827697194985986636869245829152430835021673171507607525908246940354 r 3242770073040892094735101607173275538752888766491356946211654602282309624331 s 2150742645846855766116236144967953798077492822890095121354692808525999221887 pubKey 2445157821578193538289426656074203099996547227497157254541771705133209838679

func TestStarknetKeyring_TestVector(t *testing.T) {
	var kr1 OCR2Key
	bigKey, _ := new(big.Int).SetString("2137244795266879235401249500471353867704187908407744160927664772020405449078", 10)
	feltKey, err := new(felt.Felt).SetString(bigKey.String())
	require.NoError(t, err)
	bytesKey := feltKey.Bytes()
	err = kr1.Unmarshal(bytesKey[:])
	require.NoError(t, err)
	// kr2, err := NewOCR2Key(cryptorand.Reader)
	// require.NoError(t, err)

	bytes, err := hex.DecodeString("0004acf99cb25a4803916f086440c661295b105a485efdc649ac4de9536da25b")
	require.NoError(t, err)
	configDigest, err := ocrtypes.BytesToConfigDigest(bytes)
	require.NoError(t, err)

	ctx := ocrtypes.ReportContext{
		ReportTimestamp: ocrtypes.ReportTimestamp{
			ConfigDigest: configDigest,
			Epoch:        0,
			Round:        1,
		},
		ExtraHash: [32]byte{
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1,
		},
	}

	var report []byte
	b1 := new(felt.Felt).SetUint64(1).Bytes()
	report = append(report, b1[:]...)
	b2Bytes, err := hex.DecodeString("00010203000000000000000000000000000000000000000000000000000000")
	require.NoError(t, err)
	b2 := new(felt.Felt).SetBytes(b2Bytes).Bytes()
	report = append(report, b2[:]...)
	b3 := new(felt.Felt).SetUint64(4).Bytes()
	report = append(report, b3[:]...)
	b4 := new(felt.Felt).SetUint64(99).Bytes()
	report = append(report, b4[:]...)
	report = append(report, b4[:]...)
	report = append(report, b4[:]...)
	report = append(report, b4[:]...)
	report = append(report, b1[:]...)

	// check that report hash matches expected
	msg, err := ReportToSigData(ctx, report)
	require.NoError(t, err)

	expected, err := new(felt.Felt).SetString("0x1332a8dabaabef63b03438ca50760cb9f5c0292cbf015b2395e50e6157df4e3")
	expectedBytes := expected.Bytes()
	require.NoError(t, err)
	assert.Equal(t, expectedBytes[:], msg.Bytes())

	// check that signature matches expected
	sig, err := kr1.Sign(ctx, report)
	require.NoError(t, err)

	pub := new(felt.Felt).SetBytes(sig[0:32])
	r := new(felt.Felt).SetBytes(sig[32:64])
	s := new(felt.Felt).SetBytes(sig[64:])

	bigPubExpected, _ := new(big.Int).SetString("1118148281956858477519852250235501663092798578871088714409528077622994994907", 10)
	feltPubExpected := new(felt.Felt).SetBytes(bigPubExpected.Bytes())
	assert.Equal(t, feltPubExpected, pub)

	bigRExpected, _ := new(big.Int).SetString("2898571078985034687500959842265381508927681132188252715370774777831313601543", 10)
	feltRExpected := new(felt.Felt).SetBytes(bigRExpected.Bytes())
	assert.Equal(t, feltRExpected, r)

	// test for malleability
	otherS, _ := new(big.Int).SetString("1930849708769648077928186998643944706551011476358007177069185543644456022504", 10)
	bigSExpected, _ := new(big.Int).SetString("1687653079896483135769135784451125398975732275358080312084893914240056843079", 10)

	feltSExpected := new(felt.Felt).SetBytes(bigSExpected.Bytes())
	assert.NotEqual(t, otherS, s, "signature not in canonical form")
	assert.Equal(t, feltSExpected, s)
}

func TestStarknetKeyring_Sign_Verify(t *testing.T) {
	kr1, err := NewOCR2Key(cryptorand.Reader)
	require.NoError(t, err)
	kr2, err := NewOCR2Key(cryptorand.Reader)
	require.NoError(t, err)

	digest := "00044e5d4f35325e464c87374b13c512f60e09d1236dd902f4bef4c9aedd7300"
	bytes, err := hex.DecodeString(digest)
	require.NoError(t, err)
	configDigest, err := ocrtypes.BytesToConfigDigest(bytes)
	require.NoError(t, err)

	ctx := ocrtypes.ReportContext{
		ReportTimestamp: ocrtypes.ReportTimestamp{
			ConfigDigest: configDigest,
			Epoch:        1,
			Round:        1,
		},
		ExtraHash: [32]byte{
			255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
			255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
		},
	}
	report := ocrtypes.Report{
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 97, 91, 43, 83, // observations_timestamp
		0, 1, 2, 3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, // observers
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, // len
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 73, 150, 2, 210, // observation 1
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 73, 150, 2, 211, // observation 2
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 13, 224, 182, 179, 167, 100, 0, 0, // juels per fee coin (1 with 18 decimal places)
	}

	t.Run("can verify", func(t *testing.T) {
		sig, err := kr1.Sign(ctx, report)
		require.NoError(t, err)
		result := kr2.Verify(kr1.PublicKey(), ctx, report, sig)
		require.True(t, result)
	})

	t.Run("invalid sig", func(t *testing.T) {
		result := kr2.Verify(kr1.PublicKey(), ctx, report, []byte{0x01})
		require.False(t, result)

		longSig := [100]byte{}
		result = kr2.Verify(kr1.PublicKey(), ctx, report, longSig[:])
		require.False(t, result)
	})

	t.Run("invalid pubkey", func(t *testing.T) {
		sig, err := kr1.Sign(ctx, report)
		require.NoError(t, err)

		pk := []byte{0x01}
		result := kr2.Verify(pk, ctx, report, sig)
		require.False(t, result)

		pk = big.NewInt(int64(31337)).Bytes()
		result = kr2.Verify(pk, ctx, report, sig)
		require.False(t, result)
	})
}

func TestStarknetKeyring_Marshal(t *testing.T) {
	kr1, err := NewOCR2Key(cryptorand.Reader)
	require.NoError(t, err)
	m, err := kr1.Marshal()
	require.NoError(t, err)
	kr2 := OCR2Key{}
	err = kr2.Unmarshal(m)
	require.NoError(t, err)
	assert.True(t, kr1.priv.Cmp(kr2.priv) == 0)

	// Invalid seed size should error
	require.Error(t, kr2.Unmarshal([]byte{0x01}))
}
