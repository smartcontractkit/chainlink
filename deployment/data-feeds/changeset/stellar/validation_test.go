package stellar

import (
	"encoding/hex"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// mustHex16 decodes s (<=32 hex chars) and left-justifies it into a [16]byte,
// matching the layout dataIDsToBytes is expected to produce.
func mustHex16(t *testing.T, s string) [16]byte {
	t.Helper()
	b, err := hex.DecodeString(s)
	require.NoError(t, err)
	require.LessOrEqual(t, len(b), 16)
	var out [16]byte
	copy(out[:], b)
	return out
}

func TestValidateAddress(t *testing.T) {
	require.NoError(t, validateAddress("GAAZI4TCR3TY5OJHCTJC2A4QSY6CJWJH5IAJTGKIN2ER7LBNVKOCCWN7")) // account
	require.NoError(t, validateAddress("CA7QYNF7SOWQ3GLR2BGMZEHXAVIRZA4KVWLTJJFC7MGXUA74P7UJUWDA")) // contract
	require.Error(t, validateAddress("not-a-key"))
	require.Error(t, validateAddress(""))
}

func TestDataIDsToBytes(t *testing.T) {
	// Full-width (exactly 16-byte) input: exact round-trip check. NOTE: for a
	// value this wide, left-justified copy(b[:], v.Bytes()) and
	// right-justified v.FillBytes(b[:]) produce byte-identical output, so
	// this case alone cannot catch a regression to FillBytes — see the short
	// case below for that.
	got, err := dataIDsToBytes([]string{"0x018e16c39e0003200000000000000000"})
	require.NoError(t, err)
	require.Len(t, got, 1)
	require.Equal(t, mustHex16(t, "018e16c39e0003200000000000000000"), got[0])

	// Short (sub-16-byte) shorthand input: this is the discriminating case.
	// Left-justified (copy), the nonzero prefix "18e16c39e0003200000" lands in
	// the LEADING bytes b[0:10] with zero padding trailing at b[10:16]:
	//   18 e1 6c 39 e0 00 32 00 00 00 00 00 00 00 00 00
	// Right-justified (FillBytes) would instead push the same nonzero bytes
	// down to END at b[15], i.e.:
	//   00 00 00 00 00 00 18 e1 6c 39 e0 00 32 00 00 00
	got, err = dataIDsToBytes([]string{"0x018e16c39e00032000000"})
	require.NoError(t, err)
	require.Len(t, got, 1)
	require.Equal(t, byte(0x18), got[0][0])
	require.Equal(t, mustHex16(t, "18e16c39e00032000000000000000000"), got[0])

	// >128-bit input (17 bytes) must be rejected.
	_, err = dataIDsToBytes([]string{"0x01" + strings.Repeat("00", 16)})
	require.Error(t, err)

	// Invalid hex string must be rejected.
	_, err = dataIDsToBytes([]string{"zzz"})
	require.Error(t, err)
}

func TestWorkflowFieldConversions(t *testing.T) {
	n, err := workflowNameToBytes("abc")
	require.NoError(t, err)
	require.Equal(t, byte('a'), n[0])
	_, err = workflowNameToBytes("elevenchars!")
	require.Error(t, err)

	o, err := workflowOwnerToBytes("0x0102030405060708090a0b0c0d0e0f1011121314")
	require.NoError(t, err)
	require.Equal(t, byte(0x14), o[19])
	_, err = workflowOwnerToBytes("0x01")
	require.Error(t, err)
}
