package confidentialcompute

import (
	"testing"

	"github.com/stretchr/testify/require"

	cctypes "github.com/smartcontractkit/chainlink-confidential-compute/types"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
)

func nodeSetWithValues(name string, values map[string]any) *cre.NodeSet {
	return &cre.NodeSet{
		CapabilityConfigs: map[cre.CapabilityFlag]cre.CapabilityConfig{
			name: {Values: values},
		},
	}
}

func TestEnclavesFromConfig(t *testing.T) {
	t.Parallel()

	const name = "confidential-workflows"

	t.Run("parses declared enclaves", func(t *testing.T) {
		t.Parallel()

		ns := nodeSetWithValues(name, map[string]any{
			EnclavesConfigKey: `[{"enclaveURL":"http://10.0.0.1:8080","enclaveAuthHeader":"key-a"},` +
				`{"enclaveURL":"http://10.0.0.1:8081"}]`,
		})

		got, err := EnclavesFromConfig(ns, name)
		require.NoError(t, err)
		require.Len(t, got, 2)
		require.Equal(t, "http://10.0.0.1:8080", got[0].EnclaveURL)
		require.Equal(t, "key-a", got[0].EnclaveAuthHeader)
		require.Equal(t, "http://10.0.0.1:8081", got[1].EnclaveURL)
	})

	t.Run("round-trips a marshalled enclave list", func(t *testing.T) {
		t.Parallel()

		encoded, err := MarshalEnclaves([]cctypes.Enclave{{
			EnclaveURL:    "http://10.0.0.2:8080",
			TrustedValues: [][]byte{[]byte("fake-measurements")},
			Region:        "us-west-2",
		}})
		require.NoError(t, err)

		got, err := EnclavesFromConfig(nodeSetWithValues(name, map[string]any{EnclavesConfigKey: encoded}), name)
		require.NoError(t, err)
		require.Len(t, got, 1)
		require.Equal(t, "http://10.0.0.2:8080", got[0].EnclaveURL)
		require.Equal(t, "us-west-2", got[0].Region)
		require.Equal(t, [][]byte{[]byte("fake-measurements")}, got[0].TrustedValues)
	})

	t.Run("returns nil when unconfigured", func(t *testing.T) {
		t.Parallel()

		for _, tc := range []struct {
			desc    string
			nodeSet *cre.NodeSet
		}{
			{"nil node set", nil},
			{"no values", nodeSetWithValues(name, nil)},
			{"no enclaves key", nodeSetWithValues(name, map[string]any{"other": "x"})},
			{"different capability", nodeSetWithValues("confidential-http", map[string]any{EnclavesConfigKey: "[]"})},
		} {
			t.Run(tc.desc, func(t *testing.T) {
				t.Parallel()

				got, err := EnclavesFromConfig(tc.nodeSet, name)
				require.NoError(t, err)
				require.Nil(t, got)
			})
		}
	})

	t.Run("errors on a non-string value", func(t *testing.T) {
		t.Parallel()

		_, err := EnclavesFromConfig(nodeSetWithValues(name, map[string]any{EnclavesConfigKey: 42}), name)
		require.ErrorContains(t, err, "must be a JSON string")
	})

	t.Run("errors on malformed JSON", func(t *testing.T) {
		t.Parallel()

		_, err := EnclavesFromConfig(nodeSetWithValues(name, map[string]any{EnclavesConfigKey: "{not json"}), name)
		require.ErrorContains(t, err, "failed to parse")
	})
}
