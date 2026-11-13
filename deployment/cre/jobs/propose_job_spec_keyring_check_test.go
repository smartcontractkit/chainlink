package jobs

import (
	"encoding/hex"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	chainsel "github.com/smartcontractkit/chain-selectors"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	nodev1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"
	"github.com/smartcontractkit/chainlink/deployment"
	job_types "github.com/smartcontractkit/chainlink/deployment/cre/jobs/types"
	"github.com/smartcontractkit/chainlink/deployment/cre/pkg/offchain"
)

const (
	testSignerA = "02baa006fc5f6e9c138c37fb58c962bb7f40ca7c"
	testSignerB = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	testCSAKey  = "csa-1"
)

func nodeWithEVMBundle(t *testing.T, addrHex, bundleID string) deployment.Node {
	t.Helper()
	b, err := hex.DecodeString(addrHex)
	require.NoError(t, err)

	return deployment.Node{
		Name:   "node-1",
		CSAKey: testCSAKey,
		SelToOCRConfig: map[chainsel.ChainDetails]deployment.OCRConfig{
			{ChainSelector: chainsel.ETHEREUM_MAINNET.Selector}: {OnchainPublicKey: b, KeyBundleID: bundleID},
		},
	}
}

func TestCheckKeyringsAgainstExpected(t *testing.T) {
	t.Parallel()

	padded := testSignerA + "000000000000000000000000"
	tests := []struct {
		name     string
		nodes    deployment.Nodes
		expected map[string]string
		wantErr  string
	}{
		{"padded signer matches", deployment.Nodes{nodeWithEVMBundle(t, testSignerA, "b")}, map[string]string{testCSAKey: padded}, ""},
		{"0x unpadded signer matches", deployment.Nodes{nodeWithEVMBundle(t, testSignerA, "b")}, map[string]string{testCSAKey: "0x" + testSignerA}, ""},
		{"signer mismatch", deployment.Nodes{nodeWithEVMBundle(t, testSignerB, "b")}, map[string]string{testCSAKey: padded}, "expected signer " + testSignerA + " matches no JD EVM bundle"},
		{"empty bundle id", deployment.Nodes{nodeWithEVMBundle(t, testSignerA, "")}, map[string]string{testCSAKey: padded}, "empty OCR key bundle ID"},
		{"node without expected signer", deployment.Nodes{nodeWithEVMBundle(t, testSignerA, "b")}, map[string]string{"other": padded}, "no expected signer"},
		{"no evm bundle", deployment.Nodes{{Name: "node-1", CSAKey: testCSAKey}}, map[string]string{testCSAKey: padded}, "no EVM OCR key bundle in JD"},
		{"malformed signer", deployment.Nodes{nodeWithEVMBundle(t, testSignerA, "b")}, map[string]string{testCSAKey: "zz"}, "not a 20 or 32 byte hex value"},
		{"non-zero padding rejected", deployment.Nodes{nodeWithEVMBundle(t, testSignerA, "b")}, map[string]string{testCSAKey: testSignerA + "000000000000000000000001"}, "not a zero-padded 20 byte address"},
		{"uppercase signer matches", deployment.Nodes{nodeWithEVMBundle(t, testSignerA, "b")}, map[string]string{testCSAKey: "0x" + strings.ToUpper(testSignerA)}, ""},
		{"bootstrap skipped", deployment.Nodes{{Name: "bt", CSAKey: "x", IsBootstrap: true}}, map[string]string{testCSAKey: padded}, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := checkKeyringsAgainstExpected(tt.nodes, tt.expected, specSelectors{})
			if tt.wantErr == "" {
				require.NoError(t, err)

				return
			}
			require.ErrorContains(t, err, tt.wantErr)
		})
	}
}

func TestVerifyOCRKeyrings_Gating(t *testing.T) {
	t.Parallel()

	base := ProposeJobSpecInput{
		Environment:     "staging",
		Domain:          "cre",
		DONName:         "don",
		DONFilters:      []offchain.TargetDONFilter{{Key: "k", Value: "v"}},
		JobName:         "job",
		Template:        job_types.Consensus,
		ExpectedSigners: map[string]string{testCSAKey: testSignerA},
		Inputs:          job_types.JobSpecInput{"command": "consensus"},
	}

	t.Run("no expected signers is a no-op", func(t *testing.T) {
		t.Parallel()
		in := base
		in.ExpectedSigners = nil
		require.NoError(t, verifyOCRKeyrings(cldf.Environment{}, in))
	})
	t.Run("skip flag bypasses JD", func(t *testing.T) {
		t.Parallel()
		in := base
		in.SkipKeyringCheck = true
		require.NoError(t, verifyOCRKeyrings(cldf.Environment{}, in))
	})
	t.Run("non-OCR template bypasses JD", func(t *testing.T) {
		t.Parallel()
		in := base
		in.Template = job_types.Cron
		require.NoError(t, verifyOCRKeyrings(cldf.Environment{}, in))
	})
	t.Run("EVM template bypasses JD", func(t *testing.T) {
		t.Parallel()
		in := base
		in.Template = job_types.EVM
		require.NoError(t, verifyOCRKeyrings(cldf.Environment{}, in))
	})
	t.Run("minimal oracle factory bypasses JD", func(t *testing.T) {
		t.Parallel()
		in := base
		in.Inputs = job_types.JobSpecInput{"command": "consensus", "minimalOracleFactory": true}
		require.NoError(t, verifyOCRKeyrings(cldf.Environment{}, in))
	})
	t.Run("explicit oracle factory bypasses JD", func(t *testing.T) {
		t.Parallel()
		in := base
		in.Inputs = job_types.JobSpecInput{"command": "consensus", "oracleFactory": map[string]any{"enabled": true}}
		require.NoError(t, verifyOCRKeyrings(cldf.Environment{}, in))
	})
	t.Run("invalid inputs fail before any JD call", func(t *testing.T) {
		t.Parallel()
		in := base
		in.Inputs = job_types.JobSpecInput{"chainSelectorEVM": "not-a-number"}
		require.ErrorContains(t, verifyOCRKeyrings(cldf.Environment{}, in), "failed to parse")
	})
	t.Run("OCR template without JD client errors", func(t *testing.T) {
		t.Parallel()
		require.ErrorContains(t, verifyOCRKeyrings(cldf.Environment{}, base), "needs a job distributor client")
	})
}

func selectorKeys(f *nodev1.ListNodesRequest_Filter) []string {
	keys := make([]string, 0, len(f.GetSelectors()))
	for _, s := range f.GetSelectors() {
		keys = append(keys, s.GetKey())
	}

	return keys
}

func TestTargetNodeFilter_MatchesProposingOperation(t *testing.T) {
	t.Parallel()

	base := ProposeJobSpecInput{
		Environment: "staging",
		Domain:      "cre",
		DONName:     "don",
		DONFilters:  []offchain.TargetDONFilter{{Key: "zone", Value: "a"}},
	}

	for _, tmpl := range []job_types.JobSpecTemplate{job_types.OCR3, job_types.Ring} {
		in := base
		in.Template = tmpl
		keys := selectorKeys(targetNodeFilter(in))
		require.ElementsMatch(t, []string{"zone", "type"}, keys, tmpl.String())
	}

	for _, tmpl := range []job_types.JobSpecTemplate{job_types.Consensus, job_types.EVM, job_types.Aptos, job_types.Solana, job_types.Stellar} {
		in := base
		in.Template = tmpl
		keys := selectorKeys(targetNodeFilter(in))
		require.Contains(t, keys, "don-don", tmpl.String())
		require.Contains(t, keys, "environment", tmpl.String())
		require.Contains(t, keys, "product", tmpl.String())
	}

	in := base
	in.Template = job_types.Consensus
	in.JDDONName = "jd-don"
	require.Contains(t, selectorKeys(targetNodeFilter(in)), "don-jd-don")
}

func nodeWithBundles(name string, bundles map[uint64]string) deployment.Node {
	n := deployment.Node{Name: name, CSAKey: name, SelToOCRConfig: map[chainsel.ChainDetails]deployment.OCRConfig{}}
	for sel, id := range bundles {
		n.SelToOCRConfig[chainsel.ChainDetails{ChainSelector: sel}] = deployment.OCRConfig{KeyBundleID: id}
	}

	return n
}

func TestNonEVMBundleProblems(t *testing.T) {
	t.Parallel()

	aptos := chainsel.APTOS_MAINNET.Selector
	tests := []struct {
		name  string
		nodes deployment.Nodes
		want  []string
	}{
		{
			name: "distinct bundles pass",
			nodes: deployment.Nodes{
				nodeWithBundles("n1", map[uint64]string{aptos: "b1"}),
				nodeWithBundles("n2", map[uint64]string{aptos: "b2"}),
			},
		},
		{
			name: "missing bundle",
			nodes: deployment.Nodes{
				nodeWithBundles("n1", map[uint64]string{aptos: "b1"}),
				nodeWithBundles("n2", nil),
			},
			want: []string{"n2: no OCR key bundle in JD for chain selector"},
		},
		{
			name:  "empty bundle id",
			nodes: deployment.Nodes{nodeWithBundles("n1", map[uint64]string{aptos: ""})},
			want:  []string{"n1: JD has an empty OCR key bundle ID"},
		},
		{
			name: "shared bundle id",
			nodes: deployment.Nodes{
				nodeWithBundles("n1", map[uint64]string{aptos: "same"}),
				nodeWithBundles("n2", map[uint64]string{aptos: "same"}),
			},
			want: []string{"bundle ID same is shared by nodes n1, n2"},
		},
		{
			name:  "bootstrap skipped",
			nodes: deployment.Nodes{{Name: "bt", IsBootstrap: true}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := nonEVMBundleProblems(tt.nodes, []uint64{aptos})
			require.Len(t, got, len(tt.want))
			for i, w := range tt.want {
				require.Contains(t, got[i], w)
			}
		})
	}
}

func TestSpecChainSelectors(t *testing.T) {
	t.Parallel()

	in := ProposeJobSpecInput{
		Template: job_types.Consensus,
		JobName:  "j",
		Inputs: job_types.JobSpecInput{
			"command": "consensus", "chainSelectorEVM": 5, "chainSelectorAptos": 7, "chainSelectorSolana": 8,
		},
	}
	got, err := specChainSelectors(in)
	require.NoError(t, err)
	require.Equal(t, uint64(5), got.EVM)
	require.ElementsMatch(t, []uint64{7, 8}, got.NonEVM)

	in.Inputs["ocrChainSelector"] = 6
	got, err = specChainSelectors(in)
	require.NoError(t, err)
	require.Equal(t, uint64(6), got.EVM, "ocrChainSelector wins")

	in.Template = job_types.OCR3
	in.Inputs["templateName"] = "t"
	in.Inputs["contractQualifier"] = "q"
	in.Inputs["bootstrapperOCR3Urls"] = []any{"u"}
	got, err = specChainSelectors(in)
	require.NoError(t, err)
	require.Equal(t, uint64(5), got.EVM)
	require.ElementsMatch(t, []uint64{7, 8}, got.NonEVM)

	in.Template = job_types.Ring
	got, err = specChainSelectors(in)
	require.NoError(t, err)
	require.Equal(t, uint64(5), got.EVM)
	require.Empty(t, got.NonEVM)
}

func TestEVMBundleProblems_PinnedSelector(t *testing.T) {
	t.Parallel()

	selA := chainsel.ETHEREUM_MAINNET.Selector
	selB := chainsel.ETHEREUM_TESTNET_SEPOLIA.Selector
	a, err := hex.DecodeString(testSignerA)
	require.NoError(t, err)
	b, err := hex.DecodeString(testSignerB)
	require.NoError(t, err)
	n := deployment.Node{Name: "n1", SelToOCRConfig: map[chainsel.ChainDetails]deployment.OCRConfig{
		{ChainSelector: selA}: {OnchainPublicKey: a, KeyBundleID: "ba"},
		{ChainSelector: selB}: {OnchainPublicKey: b, KeyBundleID: "bb"},
	}}

	require.Empty(t, evmBundleProblems(n, testSignerA, selA))
	require.Empty(t, evmBundleProblems(n, testSignerA, 0), "unpinned accepts any EVM chain")
	require.NotEmpty(t, evmBundleProblems(n, testSignerA, selB), "pinned to the chain whose bundle differs")
	require.NotEmpty(t, evmBundleProblems(n, testSignerA, chainsel.ETHEREUM_TESTNET_HOLESKY.Selector), "pinned chain missing")
}
