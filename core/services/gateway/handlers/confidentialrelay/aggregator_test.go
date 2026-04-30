package confidentialrelay

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

func TestAggregator_tiedMajoritiesPickDigestDeterministically(t *testing.T) {
	t.Parallel()

	lggr := logger.Test(t)
	agg := &aggregator{}

	const id = "req-tie"
	makeResp := func(payload string) jsonrpc.Response[json.RawMessage] {
		buf, err := json.Marshal(map[string]string{"payload": payload})
		require.NoError(t, err)
		rd := make(json.RawMessage, len(buf))
		copy(rd, buf)
		return jsonrpc.Response[json.RawMessage]{
			Version: jsonrpc.JsonRpcVersion,
			ID:      id,
			Method:  MethodCapabilityExec,
			Result:  &rd,
		}
	}

	a := makeResp("aaa")
	b := makeResp("zzz")
	digestA, err := a.Digest()
	require.NoError(t, err)
	digestB, err := b.Digest()
	require.NoError(t, err)
	require.NotEqual(t, digestA, digestB, "fixtures must produce distinct digests")

	wantWinnerDigest := digestA
	if digestB < digestA {
		wantWinnerDigest = digestB
	}

	// Two nodes report A, two report B: each side has F+1 when F=1. Map iteration order
	// must not change which digest wins.
	for range 300 {
		m := map[string]jsonrpc.Response[json.RawMessage]{
			"n0": a,
			"n1": a,
			"n2": b,
			"n3": b,
		}
		got, err := agg.Aggregate(m, 1, 4, lggr)
		require.NoError(t, err)
		require.NotNil(t, got)
		gotDigest, derr := got.Digest()
		require.NoError(t, derr)
		require.Equal(t, wantWinnerDigest, gotDigest,
			"with tied majorities the chosen digest must be order-independent (lexicographically smallest qualifying digest)")
	}
}
