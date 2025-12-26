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
	testEnv := test.SetupEnvV2(t, false)
	env := testEnv.Env

	_, err := testEnv.TestJD.ListNodes(t.Context(), &node.ListNodesRequest{})
	require.NoError(t, err)

	selector := testEnv.RegistrySelector // use the test environment's selector
	ds := datastore.NewMemoryDataStore()

	seedAddressesForSelector(t, ds, selector, "0xAb5801a7D398351b8bE11C439e05C5B3259aeC9B", "0x9fE46736679d2D9a65F0992F2272dE9f3c7fa6e0")
	env.DataStore = ds.Seal()

	const (
		inputLookback  int64 = 123 // non-zero input-level default
		overrideCustom int64 = 999 // per-node explicit override
	)

	input := jobs.ProposeEVMCapJobSpecInput{
		Environment:             "test",
		Zone:                    test.Zone,
		Domain:                  "cre",
		DONName:                 test.DONName,
		ChainSelector:           selector,
		OCRChainSelector:        selector,
		BootstrapperOCR3Urls:    []string{"12D3KooWabc@127.0.0.1:5001"},
		OCRContractQualifier:    testOCRQualifier,
		ForwardersQualifier:     testForwarderQualifier,
		ForwarderLookbackBlocks: inputLookback,
		EVMCapabilityInputs: []jobs.EVMCapabilityInput{
			minimalEVMCapInput("node_test-don-0"),
			minimalEVMCapInput("node_test-don-1"),
			minimalEVMCapInput("node_test-don-2"),
			minimalEVMCapInput("node_test-don-3"),
		},
	}

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
	testEnv := test.SetupEnvV2(t, false)
	env := testEnv.Env

	selector := testEnv.RegistrySelector
	ds := datastore.NewMemoryDataStore()
	seedAddressesForSelector(t, ds, selector, "0xocr...", "0xfwd...")
	env.DataStore = ds.Seal()

	input := freshBase(selector)
	input.OCRChainSelector = selector
	// duplicate
	input.EVMCapabilityInputs = []jobs.EVMCapabilityInput{
		minimalEVMCapInput("node_test-don-0"),
		minimalEVMCapInput("node_test-don-0"),
	}

	_, err := jobs.ProposeEVMCapJobSpec{}.Apply(*env, input)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "duplicate nodeID")
}

func TestProposeStandardCapabilityJob_ReusesUUIDForEvmCapabilitiesV2(t *testing.T) {
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

	input := jobs.ProposeEVMCapJobSpecInput{
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
	_, err = jobs.ProposeEVMCapJobSpec{}.Apply(*env, input)
	require.NoError(t, err)

	// verify that the jobs have been distributed and accepted
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
	verifyJobProposal(1)

	// different config generates different uuid, but evm cap jobs should lookup the old id and reuse it
	input.ForwarderLookbackBlocks = 999
	_, err = jobs.ProposeEVMCapJobSpec{}.Apply(*env, input)
	require.NoError(t, err)
	verifyJobProposal(2)
}
