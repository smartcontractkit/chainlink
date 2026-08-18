package report

import (
	"math/big"
	"testing"

	"github.com/stellar/go-stellar-sdk/xdr"
	"github.com/stretchr/testify/require"
)

func TestEncodeEntriesRoundTrips(t *testing.T) {
	var dataID [32]byte
	for i := range dataID {
		dataID[i] = 0xAB
	}
	b, err := EncodeEntries(dataID, 123456, 100)
	require.NoError(t, err)

	var v xdr.ScVal
	require.NoError(t, xdr.SafeUnmarshal(b, &v), "payload is not a single ScVal")
	vec, ok := v.GetVec()
	require.True(t, ok)
	require.NotNil(t, vec)
	require.Len(t, *vec, 1)
	m, ok := (*vec)[0].GetMap()
	require.True(t, ok)
	require.NotNil(t, m)
	require.Len(t, *m, 3)
	keys := []string{"answer", "data_id", "timestamp"}
	for i, e := range *m {
		sym, _ := e.Key.GetSym()
		require.Equal(t, keys[i], string(sym))
	}
	i256, ok := (*m)[0].Val.GetI256()
	require.True(t, ok, "answer is not i256")
	require.Zero(t, int256ToBig(i256).Cmp(big.NewInt(123456)))
	db, ok := (*m)[1].Val.GetBytes()
	require.True(t, ok)
	require.Equal(t, dataID[:], []byte(db))
	ts, ok := (*m)[2].Val.GetU64()
	require.True(t, ok)
	require.Equal(t, uint64(100), uint64(ts))
}

func TestEncodeEntriesNegativeAnswer(t *testing.T) {
	var dataID [32]byte
	dataID[0] = 1
	b, err := EncodeEntries(dataID, -42, 1)
	require.NoError(t, err)
	var v xdr.ScVal
	require.NoError(t, xdr.SafeUnmarshal(b, &v))
	vec, _ := v.GetVec()
	m, _ := (*vec)[0].GetMap()
	i256, _ := (*m)[0].Val.GetI256()
	require.Zero(t, int256ToBig(i256).Cmp(big.NewInt(-42)))
}

func int256ToBig(p xdr.Int256Parts) *big.Int {
	n := new(big.Int).SetInt64(int64(p.HiHi))
	n.Lsh(n, 64).Or(n, new(big.Int).SetUint64(uint64(p.HiLo)))
	n.Lsh(n, 64).Or(n, new(big.Int).SetUint64(uint64(p.LoHi)))
	n.Lsh(n, 64).Or(n, new(big.Int).SetUint64(uint64(p.LoLo)))
	return n
}
