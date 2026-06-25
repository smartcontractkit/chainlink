package jobs_test

import (
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/deployment/cre/jobs"
	crejobsops "github.com/smartcontractkit/chainlink/deployment/cre/jobs/operations"
	jobspkg "github.com/smartcontractkit/chainlink/deployment/cre/jobs/pkg"
	"github.com/smartcontractkit/chainlink/deployment/cre/test"
	tenv "github.com/smartcontractkit/chainlink/deployment/environment/test"

	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/runtime"
	csav1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/csa"
	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"
	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"
)

const (
	testSolSolanaFwdQualifier  = "test-solana-fwd-qualifier"
	testSolanaForwarderProgram = "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v"
	testSolanaForwarderState   = "Es9vMFrzaCERmJfrF4H2FYD4KCoNkY11McCe8BenwNYB"
	testSolanaTransmitter      = "So11111111111111111111111111111111111111112"
	testSolanaForwarderVersion = "1.0.0"
)

func seedSolanaForwarderAddresses(t *testing.T, ds *datastore.MemoryDataStore, chainSel uint64, qualifier, programAddr, stateAddr string) {
	t.Helper()
	v := semver.MustParse(testSolanaForwarderVersion)
	require.NoError(t, ds.Addresses().Add(datastore.AddressRef{
		ChainSelector: chainSel,
		Type:          jobspkg.SolanaForwarderProgramType,
		Version:       v,
		Qualifier:     qualifier,
		Address:       programAddr,
	}))
	require.NoError(t, ds.Addresses().Add(datastore.AddressRef{
		ChainSelector: chainSel,
		Type:          jobspkg.SolanaForwarderStateType,
		Version:       v,
		Qualifier:     qualifier,
		Address:       stateAddr,
	}))
}

func solanaCapInput(nodeID, transmitter string) jobs.SolanaCapabilityInput {
	return jobs.SolanaCapabilityInput{
		NodeID:             nodeID,
		Transmitter:        transmitter,
		OverrideDefaultCfg: jobs.SolanaOverrideDefaultCfg{},
	}
}

const testSolOCRQualifier = "zone-a-audited"

func seedSolanaCapRegAddress(t *testing.T, ds *datastore.MemoryDataStore, ocrSel uint64, qualifier, addr string) {
	t.Helper()
	require.NoError(t, ds.Addresses().Add(datastore.AddressRef{
		ChainSelector: ocrSel,
		Type:          datastore.ContractType("CapabilitiesRegistry"),
		Version:       semver.MustParse("2.0.0"),
		Address:       addr,
		Qualifier:     qualifier,
	}))
}

func freshSolanaReadsBase(solSel, ocrSel uint64) jobs.ProposeSolanaJobSpecInput {
	in := freshSolanaBase(solSel)
	in.ReadsEnabled = true
	in.OCRChainSelector = ocrSel
	in.BootstrapperOCR3Urls = []string{"12D3KooWxyz@127.0.0.1:5001"}
	in.OCRContractQualifier = testSolOCRQualifier
	return in
}

func freshSolanaBase(solSel uint64) jobs.ProposeSolanaJobSpecInput {
	return jobs.ProposeSolanaJobSpecInput{
		Environment:         test.EnvironmentName,
		Zone:                test.Zone,
		Domain:              "cre",
		DONName:             test.DONName,
		ChainSelector:       solSel,
		DeltaStage:          10 * time.Second,
		ForwardersQualifier: testSolSolanaFwdQualifier,
		SolanaCapabilityInputs: []jobs.SolanaCapabilityInput{
			solanaCapInput("peer-1", testSolanaTransmitter),
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
	solSel := chainsel.SOLANA_DEVNET.Selector
	ds := datastore.NewMemoryDataStore()
	seedSolanaForwarderAddresses(t, ds, solSel, testSolSolanaFwdQualifier, testSolanaForwarderProgram, testSolanaForwarderState)
	env := cldf.Environment{DataStore: ds.Seal()}

	in := freshSolanaBase(solSel)
	in.SolanaCapabilityInputs = []jobs.SolanaCapabilityInput{
		solanaCapInput("peer-1", testSolanaTransmitter),
		solanaCapInput("peer-2", testSolanaTransmitter),
	}

	err := jobs.ProposeSolanaJobSpec{}.VerifyPreconditions(env, in)
	require.NoError(t, err)
}

func TestProposeSolanaJobSpec_VerifyPreconditions_requiredFields(t *testing.T) {
	solSel := chainsel.SOLANA_DEVNET.Selector
	ds := datastore.NewMemoryDataStore()
	seedSolanaForwarderAddresses(t, ds, solSel, testSolSolanaFwdQualifier, testSolanaForwarderProgram, testSolanaForwarderState)
	env := cldf.Environment{DataStore: ds.Seal()}
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
		{"missing transmitter without JD", func(in *jobs.ProposeSolanaJobSpecInput) {
			in.SolanaCapabilityInputs[0].Transmitter = ""
			in.SolanaCapabilityInputs[0].OverrideDefaultCfg.Transmitter = ""
		}, "offchain client is required"},
		{"missing delta stage", func(in *jobs.ProposeSolanaJobSpecInput) { in.DeltaStage = 0 }, "deltaStage"},
		{"missing forwarder qualifier", func(in *jobs.ProposeSolanaJobSpecInput) { in.ForwardersQualifier = "" }, "cre forwarder qualifier is required"},
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

func TestProposeSolanaJobSpec_VerifyPreconditions_missingDatastore(t *testing.T) {
	solSel := chainsel.SOLANA_DEVNET.Selector
	ds := datastore.NewMemoryDataStore()
	env := cldf.Environment{DataStore: ds.Seal()}
	in := freshSolanaBase(solSel)

	err := jobs.ProposeSolanaJobSpec{}.VerifyPreconditions(env, in)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get Solana forwarder program")
}

func TestProposeSolanaJobSpec_VerifyPreconditions_overrideMismatches(t *testing.T) {
	solSel := chainsel.SOLANA_DEVNET.Selector
	ds := datastore.NewMemoryDataStore()
	seedSolanaForwarderAddresses(t, ds, solSel, testSolSolanaFwdQualifier, testSolanaForwarderProgram, testSolanaForwarderState)
	env := cldf.Environment{DataStore: ds.Seal()}
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
	rt              *runtime.Runtime
	h               *test.Harness
	solanaCapInputs []jobs.SolanaCapabilityInput
	baseInput       jobs.ProposeSolanaJobSpecInput
}

func setupSolanaJobTest(t *testing.T) solanaJobTestSetup {
	t.Helper()

	var (
		solSel = chainsel.SOLANA_DEVNET.Selector
		ds     = datastore.NewMemoryDataStore()
	)

	seedSolanaForwarderAddresses(t, ds, solSel, testSolSolanaFwdQualifier, testSolanaForwarderProgram, testSolanaForwarderState)

	// Harness deploys CapabilitiesRegistry on RegistrySelector (TEST_90000001); nodes have OCR configs for that chain.
	h := test.NewTestHarness(t, test.WithDatastore(ds))
	env := h.Runtime.Environment()

	nodes, err := h.TestJD.ListNodes(t.Context(), &node.ListNodesRequest{})
	require.NoError(t, err)

	var solanaCapInputs []jobs.SolanaCapabilityInput
	mockGetter := &tenv.MockJobApproverGetter{JobApprovers: make(map[string]*tenv.MockJobApprover)}
	for _, n := range nodes.GetNodes() {
		if strings.Contains(n.Id, "bootstrap") {
			continue
		}
		mockGetter.JobApprovers[n.Id] = &tenv.MockJobApprover{}
		solanaCapInputs = append(solanaCapInputs, solanaCapInput(n.Id, ""))
	}

	client := tenv.NewJobServiceClient(mockGetter)
	h.TestJD.JobServiceClient = client

	env.Offchain = struct {
		jobv1.JobServiceClient
		node.NodeServiceClient
		csav1.CSAServiceClient
	}{
		JobServiceClient:  client,
		NodeServiceClient: env.Offchain,
		CSAServiceClient:  env.Offchain,
	}

	// We need to create a new runtime from the updated environment
	h.Runtime = runtime.NewFromEnvironment(env)

	baseInput := jobs.ProposeSolanaJobSpecInput{
		Environment:            test.EnvironmentName,
		Zone:                   test.Zone,
		Domain:                 "cre",
		DONName:                test.DONName,
		ChainSelector:          solSel,
		DeltaStage:             time.Second,
		ForwardersQualifier:    testSolSolanaFwdQualifier,
		SolanaCapabilityInputs: solanaCapInputs,
	}

	return solanaJobTestSetup{
		rt:              h.Runtime,
		h:               h,
		solanaCapInputs: solanaCapInputs,
		baseInput:       baseInput,
	}
}

func TestProposeSolanaJobSpec_VerifyPreconditions_readsEnabled_success(t *testing.T) {
	solSel := chainsel.SOLANA_DEVNET.Selector
	ocrSel := chainsel.ETHEREUM_TESTNET_SEPOLIA.Selector

	ds := datastore.NewMemoryDataStore()
	seedSolanaForwarderAddresses(t, ds, solSel, testSolSolanaFwdQualifier, testSolanaForwarderProgram, testSolanaForwarderState)
	seedSolanaCapRegAddress(t, ds, ocrSel, testSolOCRQualifier, "0xAb5801a7D398351b8bE11C439e05C5B3259aeC9B")
	env := cldf.Environment{DataStore: ds.Seal()}

	in := freshSolanaReadsBase(solSel, ocrSel)
	err := jobs.ProposeSolanaJobSpec{}.VerifyPreconditions(env, in)
	require.NoError(t, err)
}

func TestProposeSolanaJobSpec_VerifyPreconditions_readsEnabled_requiresOCRFields(t *testing.T) {
	solSel := chainsel.SOLANA_DEVNET.Selector
	ocrSel := chainsel.ETHEREUM_TESTNET_SEPOLIA.Selector

	ds := datastore.NewMemoryDataStore()
	seedSolanaForwarderAddresses(t, ds, solSel, testSolSolanaFwdQualifier, testSolanaForwarderProgram, testSolanaForwarderState)
	seedSolanaCapRegAddress(t, ds, ocrSel, testSolOCRQualifier, "0xAb5801a7D398351b8bE11C439e05C5B3259aeC9B")
	env := cldf.Environment{DataStore: ds.Seal()}

	base := freshSolanaReadsBase(solSel, ocrSel)

	cases := []struct {
		name    string
		mutate  func(*jobs.ProposeSolanaJobSpecInput)
		errFrag string
	}{
		{"missing ocr chain selector", func(in *jobs.ProposeSolanaJobSpecInput) { in.OCRChainSelector = 0 }, "ocr chain selector is required"},
		{"missing bootstrapper urls", func(in *jobs.ProposeSolanaJobSpecInput) { in.BootstrapperOCR3Urls = nil }, "at least one bootstrapper OCR3 URL is required"},
		{"missing ocr contract qualifier", func(in *jobs.ProposeSolanaJobSpecInput) { in.OCRContractQualifier = "" }, "ocr contract qualifier is required"},
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

func TestProposeSolanaJobSpec_VerifyPreconditions_readsEnabled_missingCapReg(t *testing.T) {
	solSel := chainsel.SOLANA_DEVNET.Selector
	ocrSel := chainsel.ETHEREUM_TESTNET_SEPOLIA.Selector

	ds := datastore.NewMemoryDataStore()
	seedSolanaForwarderAddresses(t, ds, solSel, testSolSolanaFwdQualifier, testSolanaForwarderProgram, testSolanaForwarderState)
	env := cldf.Environment{DataStore: ds.Seal()}

	in := freshSolanaReadsBase(solSel, ocrSel)
	err := jobs.ProposeSolanaJobSpec{}.VerifyPreconditions(env, in)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get CapabilitiesRegistry address")
}

func TestProposeSolanaJobSpec_Apply_readsEnabled_includesOracleFactory(t *testing.T) {
	setup := setupSolanaJobTest(t)
	input := setup.baseInput
	input.ReadsEnabled = true
	input.OCRChainSelector = setup.h.RegistrySelector
	input.BootstrapperOCR3Urls = []string{"12D3KooWxyz@127.0.0.1:5001"}
	input.OCRContractQualifier = test.RegistryQualifier

	task := runtime.ChangesetTask(jobs.ProposeSolanaJobSpec{}, input)
	err := setup.rt.Exec(task)
	require.NoError(t, err)

	out := setup.rt.State().Outputs[task.ID()]
	require.Len(t, out.Reports, 1)

	output, ok := out.Reports[0].Output.(crejobsops.ProposeStandardCapabilityJobOutput)
	require.True(t, ok, "unexpected output type: %T", out.Reports[0].Output)
	require.NotEmpty(t, output.Specs)

	checked := 0
	for _, specs := range output.Specs {
		for _, spec := range specs {
			if !strings.Contains(spec, "solana-cap-v2") {
				continue
			}
			checked++
			assert.Contains(t, spec, `[oracle_factory]`)
			assert.Contains(t, spec, `enabled = true`)
			assert.Contains(t, spec, `strategyName = "multi-chain"`)
			assert.Contains(t, spec, `"readsEnabled":true`)
		}
	}
	require.Positive(t, checked, "expected at least one solana-cap-v2 job spec")
}

func TestProposeSolanaJobSpec_ReusesUUIDOnRepropose(t *testing.T) {
	setup := setupSolanaJobTest(t)

	input := setup.baseInput
	task1 := runtime.ChangesetTask(jobs.ProposeSolanaJobSpec{}, input)
	require.NoError(t, setup.rt.Exec(task1))
	firstExternalJobIDs := solanaCapExternalJobIDsByNodeFromTask(t, setup.rt, task1.ID())
	require.Len(t, firstExternalJobIDs, len(setup.solanaCapInputs))

	// Enabling reads changes per-node config (and adds oracle factory), which would hash to new UUIDs without lookup.
	input.ReadsEnabled = true
	input.OCRChainSelector = setup.h.RegistrySelector
	input.OCRContractQualifier = test.RegistryQualifier
	input.BootstrapperOCR3Urls = []string{"12D3KooWxyz@127.0.0.1:5001"}

	task2 := runtime.ChangesetTask(jobs.ProposeSolanaJobSpec{}, input)
	require.NoError(t, setup.rt.Exec(task2))
	secondExternalJobIDs := solanaCapExternalJobIDsByNodeFromTask(t, setup.rt, task2.ID())
	assert.Equal(t, firstExternalJobIDs, secondExternalJobIDs)
}

var solanaCapExternalJobIDRE = regexp.MustCompile(`externalJobID\s*=\s*"([^"]+)"`)

func solanaCapExternalJobIDsByNodeFromTask(t *testing.T, rt *runtime.Runtime, taskID string) map[string]string {
	t.Helper()

	out := rt.State().Outputs[taskID]
	require.Len(t, out.Reports, 1)

	output, ok := out.Reports[0].Output.(crejobsops.ProposeStandardCapabilityJobOutput)
	require.True(t, ok, "unexpected output type: %T", out.Reports[0].Output)

	ids := make(map[string]string, len(output.Specs))
	for nodeID, specs := range output.Specs {
		var found bool
		for _, spec := range specs {
			if !strings.Contains(spec, "solana-cap-v2") {
				continue
			}
			m := solanaCapExternalJobIDRE.FindStringSubmatch(spec)
			require.Len(t, m, 2, "solana-cap-v2 spec for node %s must include externalJobID", nodeID)
			ids[nodeID] = m[1]
			found = true
			break
		}
		require.True(t, found, "node %s missing solana-cap-v2 spec", nodeID)
	}

	return ids
}

func TestProposeSolanaJobSpec_Apply_success(t *testing.T) {
	setup := setupSolanaJobTest(t)
	input := setup.baseInput

	task := runtime.ChangesetTask(jobs.ProposeSolanaJobSpec{}, input)
	err := setup.rt.Exec(task)
	require.NoError(t, err)

	out := setup.rt.State().Outputs[task.ID()]
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

	_, err := jobs.ProposeSolanaJobSpec{}.Apply(setup.rt.Environment(), input)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "duplicate nodeID")
}
