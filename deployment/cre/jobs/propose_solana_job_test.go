package jobs_test

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/deployment/cre/jobs"
	"github.com/smartcontractkit/chainlink/deployment/cre/test"
	tenv "github.com/smartcontractkit/chainlink/deployment/environment/test"

	chainsel "github.com/smartcontractkit/chain-selectors"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	csav1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/csa"
	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"
	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"
)

// Valid base58 Solana pubkeys (32 bytes) for tests — capability config unmarshaling expects base58.
const (
	testSolanaForwarderProgram = "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v"
	testSolanaForwarderState   = "Es9vMFrzaCERmJfrF4H2FYD4KCoNkY11McCe8BenwNYB"
	testSolanaTransmitter      = "So11111111111111111111111111111111111111112"
)

func minimalSolanaCapInput(nodeID string) jobs.SolanaCapabilityInput {
	return jobs.SolanaCapabilityInput{
		NodeID:             nodeID,
		Transmitter:        testSolanaTransmitter,
		OverrideDefaultCfg: jobs.SolanaOverrideDefaultCfg{},
	}
}

func freshSolanaBase(solSel uint64) jobs.ProposeSolanaJobSpecInput {
	return jobs.ProposeSolanaJobSpecInput{
		Environment:         "test",
		Zone:                test.Zone,
		Domain:              "cre",
		DONName:             test.DONName,
		ChainSelector:       solSel,
		DeltaStage:          10 * time.Second,
		CREForwarderAddress: testSolanaForwarderProgram,
		CREForwarderState:   testSolanaForwarderState,
		SolanaCapabilityInputs: []jobs.SolanaCapabilityInput{
			minimalSolanaCapInput("peer-1"),
		},
	}
}

func deepCloneSolanaInput(in jobs.ProposeSolanaJobSpecInput) jobs.ProposeSolanaJobSpecInput {
	clone := in
	if len(in.SolanaCapabilityInputs) > 0 {
		clone.SolanaCapabilityInputs = append([]jobs.SolanaCapabilityInput(nil), in.SolanaCapabilityInputs...)
	}
	return clone
}

func TestProposeSolanaJobSpec_VerifyPreconditions_success(t *testing.T) {
	var env cldf.Environment
	solSel := chainsel.SOLANA_DEVNET.Selector

	in := freshSolanaBase(solSel)
	in.SolanaCapabilityInputs = []jobs.SolanaCapabilityInput{
		minimalSolanaCapInput("peer-1"),
		minimalSolanaCapInput("peer-2"),
	}

	err := jobs.ProposeSolanaJobSpec{}.VerifyPreconditions(env, in)
	require.NoError(t, err)
}

func TestProposeSolanaJobSpec_VerifyPreconditions_requiredFields(t *testing.T) {
	var env cldf.Environment
	solSel := chainsel.SOLANA_DEVNET.Selector
	base := freshSolanaBase(solSel)

	cases := []struct {
		name    string
		mutate  func(*jobs.ProposeSolanaJobSpecInput)
		errFrag string
	}{
		{"missing environment", func(in *jobs.ProposeSolanaJobSpecInput) { in.Environment = "" }, "environment is required"},
		{"missing domain", func(in *jobs.ProposeSolanaJobSpecInput) { in.Domain = "" }, "domain is required"},
		{"missing zone", func(in *jobs.ProposeSolanaJobSpecInput) { in.Zone = "" }, "zone is required"},
		{"missing don name", func(in *jobs.ProposeSolanaJobSpecInput) { in.DONName = "" }, "donName is required"},
		{"missing chain selector", func(in *jobs.ProposeSolanaJobSpecInput) { in.ChainSelector = 0 }, "chain selector is required"},
		{"missing solana inputs", func(in *jobs.ProposeSolanaJobSpecInput) { in.SolanaCapabilityInputs = nil }, "at least one solana capability input is required"},
		{"missing node id", func(in *jobs.ProposeSolanaJobSpecInput) { in.SolanaCapabilityInputs[0].NodeID = "" }, "nodeID is required for solana capability input"},
		{"missing transmitter", func(in *jobs.ProposeSolanaJobSpecInput) {
			in.SolanaCapabilityInputs[0].Transmitter = ""
			in.SolanaCapabilityInputs[0].OverrideDefaultCfg.Transmitter = ""
		}, "transmitter is required"},
		{"missing delta stage", func(in *jobs.ProposeSolanaJobSpecInput) { in.DeltaStage = 0 }, "deltaStage"},
		{"missing cre forwarder address", func(in *jobs.ProposeSolanaJobSpecInput) { in.CREForwarderAddress = "" }, "creForwarderAddress is required"},
		{"missing cre forwarder state", func(in *jobs.ProposeSolanaJobSpecInput) { in.CREForwarderState = "" }, "creForwarderState is required"},
		{"wrong chain family", func(in *jobs.ProposeSolanaJobSpecInput) {
			in.ChainSelector = chainsel.ETHEREUM_TESTNET_SEPOLIA.Selector
		}, "expected \"solana\""},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			in := deepCloneSolanaInput(base)
			tc.mutate(&in)
			err := jobs.ProposeSolanaJobSpec{}.VerifyPreconditions(env, in)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tc.errFrag)
		})
	}
}

func TestProposeSolanaJobSpec_VerifyPreconditions_overrideMismatches(t *testing.T) {
	var env cldf.Environment
	solSel := chainsel.SOLANA_DEVNET.Selector
	base := freshSolanaBase(solSel)

	t.Run("chainID mismatch when provided", func(t *testing.T) {
		in := deepCloneSolanaInput(base)
		in.SolanaCapabilityInputs[0].OverrideDefaultCfg.ChainID = "wrong-chain"
		err := jobs.ProposeSolanaJobSpec{}.VerifyPreconditions(env, in)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "chainID in override config")
	})

	t.Run("network must be solana if provided", func(t *testing.T) {
		in := deepCloneSolanaInput(base)
		in.SolanaCapabilityInputs[0].OverrideDefaultCfg.Network = "evm"
		err := jobs.ProposeSolanaJobSpec{}.VerifyPreconditions(env, in)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "network in override config must be")
	})

	t.Run("forwarder program override mismatch", func(t *testing.T) {
		in := deepCloneSolanaInput(base)
		in.SolanaCapabilityInputs[0].OverrideDefaultCfg.CREForwarderAddress = "11111111111111111111111111111111"
		err := jobs.ProposeSolanaJobSpec{}.VerifyPreconditions(env, in)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "CRE forwarder address")
	})
}

type solanaJobTestSetup struct {
	env             *cldf.Environment
	solanaCapInputs []jobs.SolanaCapabilityInput
	baseInput       jobs.ProposeSolanaJobSpecInput
}

func setupSolanaJobTest(t *testing.T) solanaJobTestSetup {
	t.Helper()
	testEnv := test.SetupEnvV2(t, false)
	solSel := chainsel.SOLANA_DEVNET.Selector

	env := testEnv.Env

	nodes, err := testEnv.TestJD.ListNodes(t.Context(), &node.ListNodesRequest{})
	require.NoError(t, err)

	var solanaCapInputs []jobs.SolanaCapabilityInput
	mockGetter := &tenv.MockJobApproverGetter{JobApprovers: make(map[string]*tenv.MockJobApprover)}
	for _, n := range nodes.GetNodes() {
		if strings.Contains(n.Id, "bootstrap") {
			continue
		}
		mockGetter.JobApprovers[n.Id] = &tenv.MockJobApprover{}
		solanaCapInputs = append(solanaCapInputs, minimalSolanaCapInput(n.Id))
	}

	client := tenv.NewJobServiceClient(mockGetter)
	testEnv.TestJD.JobServiceClient = client

	env.Offchain = struct {
		jobv1.JobServiceClient
		node.NodeServiceClient
		csav1.CSAServiceClient
	}{
		JobServiceClient:  client,
		NodeServiceClient: env.Offchain,
		CSAServiceClient:  env.Offchain,
	}

	baseInput := jobs.ProposeSolanaJobSpecInput{
		Environment:            "test",
		Zone:                   test.Zone,
		Domain:                 "cre",
		DONName:                test.DONName,
		ChainSelector:          solSel,
		DeltaStage:             time.Second,
		CREForwarderAddress:    testSolanaForwarderProgram,
		CREForwarderState:      testSolanaForwarderState,
		SolanaCapabilityInputs: solanaCapInputs,
	}

	return solanaJobTestSetup{
		env:             env,
		solanaCapInputs: solanaCapInputs,
		baseInput:       baseInput,
	}
}

func TestProposeSolanaJobSpec_Apply_success(t *testing.T) {
	setup := setupSolanaJobTest(t)
	input := setup.baseInput

	require.NoError(t, jobs.ProposeSolanaJobSpec{}.VerifyPreconditions(*setup.env, input))

	out, err := jobs.ProposeSolanaJobSpec{}.Apply(*setup.env, input)
	require.NoError(t, err)
	assert.Len(t, out.Reports, 1)
}

func TestProposeSolanaJobSpec_Apply_duplicateNodeIDs(t *testing.T) {
	setup := setupSolanaJobTest(t)
	input := setup.baseInput
	require.GreaterOrEqual(t, len(setup.solanaCapInputs), 2, "need at least 2 nodes")
	input.SolanaCapabilityInputs = []jobs.SolanaCapabilityInput{
		setup.solanaCapInputs[0],
		setup.solanaCapInputs[0],
	}

	_, err := jobs.ProposeSolanaJobSpec{}.Apply(*setup.env, input)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "duplicate nodeID")
}
