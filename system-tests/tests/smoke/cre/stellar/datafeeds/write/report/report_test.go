package report

import (
	"bytes"
	"math/big"
	"testing"

	"github.com/stellar/go-stellar-sdk/xdr"
)

func TestEncodeEntriesRoundTrips(t *testing.T) {
	var dataID [32]byte
	for i := range dataID {
		dataID[i] = 0xAB
	}
	b, err := EncodeEntries(dataID, 123456, 100)
	if err != nil {
		t.Fatal(err)
	}

	var v xdr.ScVal
	if err := xdr.SafeUnmarshal(b, &v); err != nil {
		t.Fatalf("payload is not a single ScVal: %v", err)
	}
	vec, ok := v.GetVec()
	if !ok || vec == nil || len(*vec) != 1 {
		t.Fatalf("expected ScVec with one entry, got %+v", v)
	}
	m, ok := (*vec)[0].GetMap()
	if !ok || m == nil || len(*m) != 3 {
		t.Fatalf("expected 3-entry ScMap, got %+v", (*vec)[0])
	}
	keys := []string{"answer", "data_id", "timestamp"}
	for i, e := range *m {
		sym, _ := e.Key.GetSym()
		if string(sym) != keys[i] {
			t.Fatalf("key %d = %q, want %q", i, sym, keys[i])
		}
	}
	i256, ok := (*m)[0].Val.GetI256()
	if !ok {
		t.Fatal("answer is not i256")
	}
	if got := int256ToBig(i256); got.Cmp(big.NewInt(123456)) != 0 {
		t.Fatalf("answer = %s, want 123456", got)
	}
	db, ok := (*m)[1].Val.GetBytes()
	if !ok || !bytes.Equal(db, dataID[:]) {
		t.Fatal("data_id mismatch")
	}
	ts, ok := (*m)[2].Val.GetU64()
	if !ok || uint64(ts) != 100 {
		t.Fatal("timestamp mismatch")
	}
}

func TestEncodeEntriesNegativeAnswer(t *testing.T) {
	var dataID [32]byte
	dataID[0] = 1
	b, err := EncodeEntries(dataID, -42, 1)
	if err != nil {
		t.Fatal(err)
	}
	var v xdr.ScVal
	if err := xdr.SafeUnmarshal(b, &v); err != nil {
		t.Fatal(err)
	}
	vec, _ := v.GetVec()
	m, _ := (*vec)[0].GetMap()
	i256, _ := (*m)[0].Val.GetI256()
	if got := int256ToBig(i256); got.Cmp(big.NewInt(-42)) != 0 {
		t.Fatalf("answer = %s, want -42", got)
	}
}

func int256ToBig(p xdr.Int256Parts) *big.Int {
	n := new(big.Int).SetInt64(int64(p.HiHi))
	n.Lsh(n, 64).Or(n, new(big.Int).SetUint64(uint64(p.HiLo)))
	n.Lsh(n, 64).Or(n, new(big.Int).SetUint64(uint64(p.LoHi)))
	n.Lsh(n, 64).Or(n, new(big.Int).SetUint64(uint64(p.LoLo)))
	return n
}
