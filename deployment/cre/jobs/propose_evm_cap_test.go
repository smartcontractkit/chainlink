package jobs_test

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/deployment/cre/jobs/pkg"
	"github.com/smartcontractkit/chainlink/deployment/cre/pkg/offchain"
	tenv "github.com/smartcontractkit/chainlink/deployment/environment/test"

	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	csav1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/csa"
	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"
	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"
	"github.com/smartcontractkit/chainlink/deployment/cre/jobs"
	"github.com/smartcontractkit/chainlink/deployment/cre/test"

	"github.com/smartcontractkit/chainlink/deployment/cre/ocr3"
)

const (
	testForwarderContractType = datastore.ContractType("KeystoneForwarder")
	testOCRQualifier          = "OCR3Capability"
	testForwarderQualifier    = "forwarder-qualifier"
)

type evmCapTestSetup struct {
	env          *cldf.Environment
	nodeIDs      []string
	evmCapInputs []jobs.EVMCapabilityInput
	baseInput    jobs.ProposeEVMCapJobSpecInput
}

func setupEVMCapTest(t *testing.T) evmCapTestSetup {
	t.Helper()
	testEnv := test.SetupEnvV2(t, false)

	selector := testEnv.RegistrySelector
	ds := datastore.NewMemoryDataStore()
	seedAddressesForSelector(t, ds, selector, "0xocr...", "0xfwd...")
	env := testEnv.Env
	env.DataStore = ds.Seal()

	nodes, err := testEnv.TestJD.ListNodes(t.Context(), &node.ListNodesRequest{})
	require.NoError(t, err)

	var nodeIDs []string
	var evmCapInputs []jobs.EVMCapabilityInput
	mockGetter := &tenv.MockJobApproverGetter{JobApprovers: make(map[string]*tenv.MockJobApprover)}
	for _, n := range nodes.GetNodes() {
		if strings.Contains(n.Id, "bootstrap") {
			continue
		}
		nodeIDs = append(nodeIDs, n.Id)
		mockGetter.JobApprovers[n.Id] = &tenv.MockJobApprover{}
		evmCapInputs = append(evmCapInputs, minimalEVMCapInput(n.Id))
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

	baseInput := jobs.ProposeEVMCapJobSpecInput{
		Environment:             "test",
		Zone:                    test.Zone,
		Domain:                  "cre",
		DONName:                 test.DONName,
		ChainSelector:           selector,
		OCRChainSelector:        selector,
		BootstrapperOCR3Urls:    []string{"12D3KooWabc@127.0.0.1:5001"},
		OCRContractQualifier:    testOCRQualifier,
		ForwardersQualifier:     testForwarderQualifier,
		ForwarderLookbackBlocks: 123,
		EVMCapabilityInputs:     evmCapInputs,
	}

	return evmCapTestSetup{
		env:          env,
		nodeIDs:      nodeIDs,
		evmCapInputs: evmCapInputs,
		baseInput:    baseInput,
	}
}

func minimalEVMCapInput(nodeID string) jobs.EVMCapabilityInput {
	return jobs.EVMCapabilityInput{
		NodeID:             nodeID,
		OverrideDefaultCfg: jobs.OverrideDefaultCfg{},
	}
}

func deepCloneInput(in jobs.ProposeEVMCapJobSpecInput) jobs.ProposeEVMCapJobSpecInput {
	clone := in
	if len(in.EVMCapabilityInputs) > 0 {
		clone.EVMCapabilityInputs = append([]jobs.EVMCapabilityInput(nil), in.EVMCapabilityInputs...)
	}
	return clone
}

func freshBase(selector uint64) jobs.ProposeEVMCapJobSpecInput {
	return jobs.ProposeEVMCapJobSpecInput{
		Environment:          "test",
		Zone:                 test.Zone,
		Domain:               "cre",
		DONName:              test.DONName,
		ChainSelector:        selector,
		OCRChainSelector:     selector,
		BootstrapperOCR3Urls: []string{"12D3KooWxyz@127.0.0.1:5001"},
		OCRContractQualifier: testOCRQualifier,
		ForwardersQualifier:  testForwarderQualifier,
		EVMCapabilityInputs:  []jobs.EVMCapabilityInput{minimalEVMCapInput("peer-1")},
	}
}

func seedAddressesForSelector(t *testing.T, ds *datastore.MemoryDataStore, sel uint64, ocrAddr, fwdAddr string) {
	t.Helper()
	require.NoError(t, ds.Addresses().Add(datastore.AddressRef{
		ChainSelector: sel,
		Type:          datastore.ContractType(ocr3.OCR3Capability),
		Version:       semver.MustParse("1.0.0"),
		Address:       ocrAddr,
		Qualifier:     testOCRQualifier,
	}))
	require.NoError(t, ds.Addresses().Add(datastore.AddressRef{
		ChainSelector: sel,
		Type:          testForwarderContractType,
		Version:       semver.MustParse("1.0.0"),
		Address:       fwdAddr,
		Qualifier:     testForwarderQualifier,
	}))
}

func TestProposeEVMCapJobSpec_VerifyPreconditions_success(t *testing.T) {
	var env cldf.Environment

	ds := datastore.NewMemoryDataStore()
	chain := chainsel.ETHEREUM_TESTNET_SEPOLIA
	ocrAddr := "0x1111111111111111111111111111111111111111"
	fwdAddr := "0x2222222222222222222222222222222222222222"
	seedAddressesForSelector(t, ds, chain.Selector, ocrAddr, fwdAddr)
	env.DataStore = ds.Seal()

	in := jobs.ProposeEVMCapJobSpecInput{
		Environment:          "test",
		Zone:                 test.Zone,
		Domain:               "cre",
		DONName:              test.DONName,
		ChainSelector:        chain.Selector,
		OCRChainSelector:     chain.Selector,
		BootstrapperOCR3Urls: []string{"12D3KooWxyz@127.0.0.1:5001"},
		OCRContractQualifier: testOCRQualifier,
		ForwardersQualifier:  testForwarderQualifier,
		EVMCapabilityInputs: []jobs.EVMCapabilityInput{
			minimalEVMCapInput("peer-1"),
			minimalEVMCapInput("peer-2"),
		},
	}

	err := jobs.ProposeEVMCapJobSpec{}.VerifyPreconditions(env, in)
	require.NoError(t, err)
}

func TestProposeEVMCapJobSpec_VerifyPreconditions_requiredFields(t *testing.T) {
	var env cldf.Environment
	ds := datastore.NewMemoryDataStore()
	selector := chainsel.ETHEREUM_TESTNET_SEPOLIA.Selector
	seedAddressesForSelector(t, ds, selector, "0x1111111111111111111111111111111111111111", "0x2222222222222222222222222222222222222222")
	env.DataStore = ds.Seal()

	base := freshBase(selector)

	cases := []struct {
		name    string
		mutate  func(*jobs.ProposeEVMCapJobSpecInput)
		errFrag string
	}{
		{"missing environment", func(in *jobs.ProposeEVMCapJobSpecInput) { in.Environment = "" }, "environment is required"},
		{"missing domain", func(in *jobs.ProposeEVMCapJobSpecInput) { in.Domain = "" }, "domain is required"},
		{"missing don name", func(in *jobs.ProposeEVMCapJobSpecInput) { in.DONName = "" }, "donName is required"},
		{"missing zone", func(in *jobs.ProposeEVMCapJobSpecInput) { in.Zone = "" }, "zone is required"},
		{"missing chain selector", func(in *jobs.ProposeEVMCapJobSpecInput) { in.ChainSelector = 0 }, "chain selector is required"},
		{"missing ocr chain selector", func(in *jobs.ProposeEVMCapJobSpecInput) { in.OCRChainSelector = 0 }, "ocr chain selector is required"},
		{"missing evm inputs", func(in *jobs.ProposeEVMCapJobSpecInput) { in.EVMCapabilityInputs = nil }, "at least one evm capability input is required"},
		{"missing bootstrapper urls", func(in *jobs.ProposeEVMCapJobSpecInput) { in.BootstrapperOCR3Urls = nil }, "at least one bootstrapper OCR3 URL is required"},
		{"empty bootstrapper url element", func(in *jobs.ProposeEVMCapJobSpecInput) { in.BootstrapperOCR3Urls = []string{""} }, "bootstrapper OCR3 URL at index 0 is empty"},
		{"missing OCR qualifier", func(in *jobs.ProposeEVMCapJobSpecInput) { in.OCRContractQualifier = "" }, "ocr contract qualifier is required"},
		{"missing forwarder qualifier", func(in *jobs.ProposeEVMCapJobSpecInput) { in.ForwardersQualifier = "" }, "cre forwarder qualifier is required"},
		{"missing node id", func(in *jobs.ProposeEVMCapJobSpecInput) { in.EVMCapabilityInputs[0].NodeID = "" }, "nodeID is required for evm capability input"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			in := deepCloneInput(base)
			tc.mutate(&in)
			err := jobs.ProposeEVMCapJobSpec{}.VerifyPreconditions(env, in)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tc.errFrag)
		})
	}
}

func TestProposeEVMCapJobSpec_VerifyPreconditions_missingAddresses(t *testing.T) {
	var env cldf.Environment
	ds := datastore.NewMemoryDataStore()
	selector := chainsel.ETHEREUM_TESTNET_SEPOLIA.Selector
	// Only seed forwarder so OCR lookup fails
	require.NoError(t, ds.Addresses().Add(datastore.AddressRef{
		ChainSelector: selector,
		Type:          testForwarderContractType,
		Version:       semver.MustParse("1.0.0"),
		Address:       "0x2222222222222222222222222222222222222222",
		Qualifier:     testForwarderQualifier,
	}))
	env.DataStore = ds.Seal()

	in := freshBase(selector)
	in.OCRChainSelector = selector

	err := jobs.ProposeEVMCapJobSpec{}.VerifyPreconditions(env, in)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get OCR contract address")

	// Now seed OCR only and remove forwarder by using a fresh DS
	ds2 := datastore.NewMemoryDataStore()
	require.NoError(t, ds2.Addresses().Add(datastore.AddressRef{
		ChainSelector: selector,
		Type:          datastore.ContractType(ocr3.OCR3Capability),
		Version:       semver.MustParse("1.0.0"),
		Address:       "0x1111111111111111111111111111111111111111",
		Qualifier:     testOCRQualifier,
	}))
	env.DataStore = ds2.Seal()

	err = jobs.ProposeEVMCapJobSpec{}.VerifyPreconditions(env, in)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get CRE forwarder address")
}

func TestProposeEVMCapJobSpec_VerifyPreconditions_mismatchAndMinimums(t *testing.T) {
	var env cldf.Environment
	ds := datastore.NewMemoryDataStore()
	chain := chainsel.ETHEREUM_TESTNET_SEPOLIA
	seedAddressesForSelector(t, ds, chain.Selector, "0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", "0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb")
	env.DataStore = ds.Seal()

	chainIDStr, _ := chainsel.GetChainIDFromSelector(chain.Selector) // e.g., "11155111"

	base := freshBase(chain.Selector)

	t.Run("chainID mismatch when provided", func(t *testing.T) {
		in := deepCloneInput(base)
		// Provide a mismatching chainID (increment the parsed one)
		var wrong uint64 = 1
		if idNum, err := strconv.ParseUint(chainIDStr, 10, 64); err == nil {
			wrong = idNum + 1
		}
		in.EVMCapabilityInputs[0].OverrideDefaultCfg.ChainID = wrong

		err := jobs.ProposeEVMCapJobSpec{}.VerifyPreconditions(env, in)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "chainID in override config")
	})

	t.Run("network must be evm if provided", func(t *testing.T) {
		in := deepCloneInput(base)
		in.EVMCapabilityInputs[0].OverrideDefaultCfg.Network = "solana"
		err := jobs.ProposeEVMCapJobSpec{}.VerifyPreconditions(env, in)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "network in override config must be")
	})

	t.Run("forwarder address mismatch when provided", func(t *testing.T) {
		in := deepCloneInput(base)
		in.EVMCapabilityInputs[0].OverrideDefaultCfg.CREForwarderAddress = "0xdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef"
		err := jobs.ProposeEVMCapJobSpec{}.VerifyPreconditions(env, in)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "CRE forwarder address in override config")
	})

	t.Run("below-minimum values are rejected, zeros allowed because of defaults", func(t *testing.T) {
		// Zeros OK (treated as "use defaults"), so verify passes:
		in := deepCloneInput(base)
		require.NoError(t, jobs.ProposeEVMCapJobSpec{}.VerifyPreconditions(env, in))

		in = deepCloneInput(base)
		in.EVMCapabilityInputs[0].OverrideDefaultCfg.LogTriggerPollInterval = 1500000000 // ns
		in.EVMCapabilityInputs[0].OverrideDefaultCfg.ReceiverGasMinimum = 500
		in.EVMCapabilityInputs[0].OverrideDefaultCfg.LogTriggerSendChannelBufferSize = 3000
		require.NoError(t, jobs.ProposeEVMCapJobSpec{}.VerifyPreconditions(env, in))

		// Now set below-minimums independently each time
		in = deepCloneInput(base)
		//nolint:nolintlint,gosec // disable G115
		in.EVMCapabilityInputs[0].OverrideDefaultCfg.LogTriggerPollInterval = uint64(1499 * time.Millisecond)
		err := jobs.ProposeEVMCapJobSpec{}.VerifyPreconditions(env, in)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "logTriggerPollInterval")

		in = deepCloneInput(base)
		in.EVMCapabilityInputs[0].OverrideDefaultCfg.ReceiverGasMinimum = 499
		err = jobs.ProposeEVMCapJobSpec{}.VerifyPreconditions(env, in)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "receiverGasMinimum")

		in = deepCloneInput(base)
		in.EVMCapabilityInputs[0].OverrideDefaultCfg.LogTriggerSendChannelBufferSize = 2999
		err = jobs.ProposeEVMCapJobSpec{}.VerifyPreconditions(env, in)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "logTriggerSendChannelBufferSize")
	})
}

func TestProposeEVMCapJobSpec_Apply_success(t *testing.T) {
	setup := setupEVMCapTest(t)
	env := setup.env

	const (
		inputLookback  int64 = 123 // non-zero input-level default
		overrideCustom int64 = 999 // per-node explicit override
	)

	input := setup.baseInput
	input.ForwarderLookbackBlocks = inputLookback

	// Use the nodes from setup but ensure we have at least 4 nodes for this test
	require.GreaterOrEqual(t, len(setup.evmCapInputs), 4, "need at least 4 nodes for this test")
	input.EVMCapabilityInputs = setup.evmCapInputs[:4]

	// Explicit per-node override on first node (should be preserved).
	input.EVMCapabilityInputs[0].OverrideDefaultCfg.ForwarderLookbackBlocks = overrideCustom

	// Verify should pass
	require.NoError(t, jobs.ProposeEVMCapJobSpec{}.VerifyPreconditions(*env, input))

	out, err := jobs.ProposeEVMCapJobSpec{}.Apply(*env, input)
	require.NoError(t, err)
	assert.Len(t, out.Reports, 1)

	// Validate exactly one override and three defaults
	outputStr := fmt.Sprintf("%v", out.Reports[0].Output)
	count999 := strings.Count(outputStr, `"forwarderLookbackBlocks":999`)
	count123 := strings.Count(outputStr, `"forwarderLookbackBlocks":123`)
	assert.Equal(t, 1, count999, "expected exactly one override lookbackBlocks=999")
	assert.Equal(t, 3, count123, "expected exactly three defaulted lookbackBlocks=123")
}

func TestProposeEVMCapJobSpec_Apply_duplicateNodeIDs(t *testing.T) {
	setup := setupEVMCapTest(t)
	env := setup.env

	input := setup.baseInput
	// duplicate
	require.GreaterOrEqual(t, len(setup.evmCapInputs), 1, "need at least 1 node for this test")
	input.EVMCapabilityInputs = []jobs.EVMCapabilityInput{
		setup.evmCapInputs[0],
		setup.evmCapInputs[0],
	}

	_, err := jobs.ProposeEVMCapJobSpec{}.Apply(*env, input)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "duplicate nodeID")
}

func TestProposeStandardCapabilityJob_ReusesUUIDForEvmCapabilitiesV2(t *testing.T) {
	setup := setupEVMCapTest(t)
	env := setup.env
	nodeIDs := setup.nodeIDs
	baseInput := setup.baseInput

	verifyJobProposal := func(revision int64) {
		jobProposals, err := env.Offchain.ListProposals(t.Context(), &jobv1.ListProposalsRequest{})
		require.NoError(t, err)

		jobsList, err := env.Offchain.ListJobs(t.Context(), &jobv1.ListJobsRequest{Filter: &jobv1.ListJobsRequest_Filter{NodeIds: nodeIDs}})
		require.NoError(t, err)

		proposalsAtRevision := 0
		for _, jp := range jobProposals.GetProposals() {
			require.Equal(t, jobv1.ProposalStatus_PROPOSAL_STATUS_APPROVED, jp.Status)
			if revision == jp.Revision {
				proposalsAtRevision++
			}
		}

		require.Len(t, jobsList.GetJobs(), len(nodeIDs))
		require.Len(t, jobProposals.GetProposals(), len(nodeIDs)*int(revision))
		require.Equal(t, len(nodeIDs), proposalsAtRevision)
	}

	// verify that the jobs have been distributed and accepted
	input := baseInput
	_, err := jobs.ProposeEVMCapJobSpec{}.Apply(*env, input)
	require.NoError(t, err)

	verifyJobProposal(1)

	// different config generates different uuid, but evm cap jobs should lookup the old id and reuse it
	input.ForwarderLookbackBlocks = 999
	_, err = jobs.ProposeEVMCapJobSpec{}.Apply(*env, input)
	require.NoError(t, err)
	verifyJobProposal(2)

	// Verify that jobs use the new name format (evm-cap-v2)
	jobProposals, err := env.Offchain.ListProposals(t.Context(), &jobv1.ListProposalsRequest{})
	require.NoError(t, err)

	hasNewFormat := false
	for _, jp := range jobProposals.GetProposals() {
		if jp.Status == jobv1.ProposalStatus_PROPOSAL_STATUS_APPROVED {
			if strings.Contains(jp.Spec, `name = "evm-cap-v2`) {
				hasNewFormat = true
				break
			}
		}
	}
	assert.True(t, hasNewFormat, "should find job with new name format (evm-cap-v2)")
}

func TestProposeStandardCapabilityJob_ReusesUUIDWithLegacyNameFormat(t *testing.T) {
	setup := setupEVMCapTest(t)
	env := setup.env
	nodeIDs := setup.nodeIDs
	baseInput := setup.baseInput

	// First create a job with legacy name format using ProposeStandardCapabilityJob
	// This simulates an existing job that was created with the old name format
	legacyJobInput := jobs.ProposeStandardCapabilityJobInput{
		JobName: "evm-capabilities-v2--test-zone-1",
		Command: "/usr/local/bin/evm",
		DONName: test.DONName,
		Domain:  "cre",
		DONFilters: []offchain.TargetDONFilter{
			{Key: "zone", Value: test.Zone},
		},
		GenerateOracleFactory: true,
		ContractQualifier:     testOCRQualifier,
		ChainSelectorEVM:      pkg.ChainSelector(baseInput.ChainSelector),
		BootstrapPeers:        baseInput.BootstrapperOCR3Urls,
	}

	_, err := jobs.ProposeStandardCapabilityJob{}.Apply(*env, legacyJobInput)
	require.NoError(t, err)

	// Get initial proposal count after creating legacy job
	initialProposals, err := env.Offchain.ListProposals(t.Context(), &jobv1.ListProposalsRequest{})
	require.NoError(t, err)
	initialCount := len(initialProposals.GetProposals())

	verifyJobProposal := func(expectedNewProposals int) {
		jobProposals, err := env.Offchain.ListProposals(t.Context(), &jobv1.ListProposalsRequest{})
		require.NoError(t, err)

		jobsList, err := env.Offchain.ListJobs(t.Context(), &jobv1.ListJobsRequest{Filter: &jobv1.ListJobsRequest_Filter{NodeIds: nodeIDs}})
		require.NoError(t, err)

		require.Len(t, jobsList.GetJobs(), len(nodeIDs))
		require.Len(t, jobProposals.GetProposals(), initialCount+expectedNewProposals)
	}

	// Now propose again with ProposeEVMCapJobSpec (which uses new format evm-cap-v2)
	// The system should detect the existing legacy job and handle it correctly
	input := baseInput
	_, err = jobs.ProposeEVMCapJobSpec{}.Apply(*env, input)
	require.NoError(t, err)

	verifyJobProposal(len(nodeIDs))

	// different config generates different uuid, but evm cap jobs should lookup the old id and reuse it
	input.ForwarderLookbackBlocks = 999
	_, err = jobs.ProposeEVMCapJobSpec{}.Apply(*env, input)
	require.NoError(t, err)
	verifyJobProposal(len(nodeIDs) * 2)

	// Verify that jobs use the legacy name format (evm-capabilities-v2) when UUIDs are reused
	jobProposals, err := env.Offchain.ListProposals(t.Context(), &jobv1.ListProposalsRequest{})
	require.NoError(t, err)

	hasLegacyFormat := false
	for _, jp := range jobProposals.GetProposals() {
		if jp.Status == jobv1.ProposalStatus_PROPOSAL_STATUS_APPROVED {
			if strings.Contains(jp.Spec, `name = "evm-capabilities-v2`) {
				hasLegacyFormat = true
				break
			}
		}
	}
	assert.True(t, hasLegacyFormat, "should find job with legacy name format (evm-capabilities-v2) when UUIDs are reused")
}
