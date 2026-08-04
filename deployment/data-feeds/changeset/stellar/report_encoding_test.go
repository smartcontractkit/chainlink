package stellar

import (
	"math/big"
	"testing"

	"github.com/stellar/go-stellar-sdk/xdr"
	"github.com/stretchr/testify/require"
)

// TestBuildReportMetadataLayout pins buildReportMetadata's 64-byte on_report
// metadata header: [0:32) workflow_cid zero | [32:42) workflow_name |
// [42:62) workflow_owner | [62:64) report_id zero.
func TestBuildReportMetadataLayout(t *testing.T) {
	var name [10]byte
	for i := range name {
		name[i] = byte(i + 1) // 0x01..0x0a, distinguishable from zero padding
	}
	var owner [20]byte
	for i := range owner {
		owner[i] = byte(i + 100) // 0x64..0x77
	}

	got := buildReportMetadata(name, owner)
	require.Len(t, got, 64)

	want := make([]byte, 64)
	copy(want[32:42], name[:])
	copy(want[42:62], owner[:])
	require.Equal(t, want, got)

	require.Equal(t, make([]byte, 32), got[0:32], "workflow_cid bytes must be zero")
	require.Equal(t, name[:], got[32:42], "workflow_name occupies [32:42)")
	require.Equal(t, owner[:], got[42:62], "workflow_owner occupies [42:62)")
	require.Equal(t, []byte{0, 0}, got[62:64], "report_id bytes must be zero")
}

// TestBuildReportScMapKeyOrder pins that buildReport's ScMap struct fields
// decode back in alphabetical key order (answer < data_id < timestamp), per
// scval.BuildStructScVal's sort.Strings over field names.
func TestBuildReportScMapKeyOrder(t *testing.T) {
	var wireHi16 [16]byte
	for i := range wireHi16 {
		wireHi16[i] = byte(i + 1)
	}
	answer := big.NewInt(123456789)
	const ts = uint64(1700000000)

	b := buildReport(t, wireHi16, answer, ts)

	var decoded xdr.ScVal
	require.NoError(t, decoded.UnmarshalBinary(b))
	require.Equal(t, xdr.ScValTypeScvVec, decoded.Type)
	require.NotNil(t, decoded.Vec)
	vec := *decoded.Vec
	require.Len(t, *vec, 1)

	entry := (*vec)[0]
	require.Equal(t, xdr.ScValTypeScvMap, entry.Type)
	require.NotNil(t, entry.Map)
	scMap := *entry.Map
	require.Len(t, *scMap, 3)

	keys := make([]string, len(*scMap))
	for i, e := range *scMap {
		require.Equal(t, xdr.ScValTypeScvSymbol, e.Key.Type)
		require.NotNil(t, e.Key.Sym)
		keys[i] = string(*e.Key.Sym)
	}
	require.Equal(t, []string{"answer", "data_id", "timestamp"}, keys)
}
