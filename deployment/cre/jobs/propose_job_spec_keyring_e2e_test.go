package jobs_test

import (
	"context"
	"encoding/hex"
	"errors"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/cre/jobs"
	job_types "github.com/smartcontractkit/chainlink/deployment/cre/jobs/types"
	"github.com/smartcontractkit/chainlink/deployment/cre/pkg/offchain"
	"github.com/smartcontractkit/chainlink/deployment/cre/test"
	envtest "github.com/smartcontractkit/chainlink/deployment/environment/test"
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
		// corrupt every entry so the check hits a targeted node whichever nodes the proposal selects
		bad := make(map[string]string, len(expected))
		for csa := range expected {
			bad[csa] = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
		}

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

	t.Run("EVM chain config without OCR2 metadata errors instead of panicking", func(t *testing.T) {
		broken := env
		broken.Offchain = rewriteChainConfigsJD{h.TestJD, func(cc *node.ChainConfig) { cc.Ocr2Config = nil }}

		err := jobs.ProposeJobSpec{}.VerifyPreconditions(broken, newInput(expected))
		require.ErrorContains(t, err, "no OCR2 or P2P key bundle")
	})

	t.Run("an unrelated chain config without a key bundle is skipped", func(t *testing.T) {
		broken := env
		broken.Offchain = rewriteChainConfigsJD{h.TestJD, func(cc *node.ChainConfig) {
			if cc.Chain.Type == node.ChainType_CHAIN_TYPE_APTOS {
				cc.Ocr2Config.OcrKeyBundle = nil
			}
		}}

		require.NoError(t, jobs.ProposeJobSpec{}.VerifyPreconditions(broken, newInput(expected)))
	})
}

// rewriteChainConfigsJD serves the harness JD's chain configs after applying rewrite to each copy.
type rewriteChainConfigsJD struct {
	*envtest.JDNodeService
	rewrite func(*node.ChainConfig)
}

func (j rewriteChainConfigsJD) ListNodeChainConfigs(ctx context.Context, req *node.ListNodeChainConfigsRequest, opts ...grpc.CallOption) (*node.ListNodeChainConfigsResponse, error) {
	resp, err := j.JDNodeService.ListNodeChainConfigs(ctx, req, opts...)
	if err != nil {
		return nil, err
	}
	out := make([]*node.ChainConfig, len(resp.ChainConfigs))
	for i, cc := range resp.ChainConfigs {
		c, ok := proto.Clone(cc).(*node.ChainConfig)
		if !ok {
			return nil, errors.New("unexpected chain config type")
		}
		j.rewrite(c)
		out[i] = c
	}

	return &node.ListNodeChainConfigsResponse{ChainConfigs: out}, nil
}
