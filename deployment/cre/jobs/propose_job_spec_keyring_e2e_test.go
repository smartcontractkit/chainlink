package jobs_test

import (
	"encoding/hex"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/cre/jobs"
	job_types "github.com/smartcontractkit/chainlink/deployment/cre/jobs/types"
	"github.com/smartcontractkit/chainlink/deployment/cre/pkg/offchain"
	"github.com/smartcontractkit/chainlink/deployment/cre/test"
)

func TestProposeJobSpec_VerifyPreconditions_KeyringCheck(t *testing.T) {
	h := test.NewTestHarness(t)
	env := h.Runtime.Environment()

	resp, err := h.TestJD.ListNodes(t.Context(), &node.ListNodesRequest{})
	require.NoError(t, err)
	ids := make([]string, 0, len(resp.Nodes))
	for _, n := range resp.Nodes {
		ids = append(ids, n.Id)
	}
	infos, err := deployment.NodeInfo(ids, h.TestJD)
	require.NoError(t, err)

	expected := make(map[string]string)
	for _, n := range infos {
		if n.IsBootstrap {
			continue
		}
		cfg, ok := n.OCRConfigForChainSelector(h.RegistrySelector)
		require.True(t, ok)
		expected[n.CSAKey] = hex.EncodeToString(cfg.OnchainPublicKey)
	}
	require.NotEmpty(t, expected)

	newInput := func(signers map[string]string) jobs.ProposeJobSpecInput {
		return jobs.ProposeJobSpecInput{
			Environment: test.EnvironmentName,
			Domain:      "cre",
			JobName:     "consensus-cap-job",
			DONName:     test.DONName,
			Template:    job_types.Consensus,
			DONFilters: []offchain.TargetDONFilter{
				{Key: offchain.FilterKeyDONName, Value: test.DONName},
				{Key: "environment", Value: test.EnvironmentName},
				{Key: "product", Value: offchain.ProductLabel},
			},
			Inputs: job_types.JobSpecInput{
				"command":          "consensus",
				"chainSelectorEVM": strconv.FormatUint(h.RegistrySelector, 10),
			},
			ExpectedSigners: signers,
		}
	}

	t.Run("matching signers pass", func(t *testing.T) {
		require.NoError(t, jobs.ProposeJobSpec{}.VerifyPreconditions(env, newInput(expected)))
	})

	t.Run("no expected signers skips the check", func(t *testing.T) {
		require.NoError(t, jobs.ProposeJobSpec{}.VerifyPreconditions(env, newInput(nil)))
	})

	t.Run("a wrong signer is rejected", func(t *testing.T) {
		bad := make(map[string]string, len(expected))
		var victim string
		for csa, addr := range expected {
			bad[csa] = addr
			victim = csa
		}
		bad[victim] = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"

		err := jobs.ProposeJobSpec{}.VerifyPreconditions(env, newInput(bad))
		require.ErrorContains(t, err, "matches no JD EVM bundle")

		in := newInput(bad)
		in.SkipKeyringCheck = true
		require.NoError(t, jobs.ProposeJobSpec{}.VerifyPreconditions(env, in))
	})

	t.Run("non-EVM bundles shared between nodes are rejected", func(t *testing.T) {
		in := newInput(expected)
		in.Inputs["chainSelectorAptos"] = strconv.FormatUint(h.AptosSelector, 10)

		// the harness gives every node the same fake aptos bundle ID
		err := jobs.ProposeJobSpec{}.VerifyPreconditions(env, in)
		require.ErrorContains(t, err, "is shared by nodes")
	})

	t.Run("a node missing from expected signers is rejected", func(t *testing.T) {
		// one entry leaves the other targeted nodes without an expected signer
		partial := make(map[string]string)
		for csa, addr := range expected {
			partial[csa] = addr

			break
		}

		err := jobs.ProposeJobSpec{}.VerifyPreconditions(env, newInput(partial))
		require.ErrorContains(t, err, "no expected signer for csa")
	})
}
