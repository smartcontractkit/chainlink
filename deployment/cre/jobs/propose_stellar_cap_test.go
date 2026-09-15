package jobs_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/deployment/cre/forwarder/stellar"
	"github.com/smartcontractkit/chainlink/deployment/cre/jobs"
	"github.com/smartcontractkit/chainlink/deployment/cre/ocr3"
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
	testStellarOCRQualifier     = "stellar-ocr-qualifier"
	testStellarFwdQualifier     = "test-stellar-fwd-qualifier"
	testStellarForwarderAddress = "CDLZFC3SYJYDZT7K67VZ75HPJVIEUVNIXF47ZG2FB2RMQQVU2HHGCYSC"
)

func minimalStellarCapInput(nodeID string) jobs.StellarCapabilityInput {
	return jobs.StellarCapabilityInput{
		NodeID:             nodeID,
		OverrideDefaultCfg: jobs.StellarOverrideDefaultCfg{},
	}
}

// seedStellarOCR3 seeds only the OCR3 contract. Used by the Apply tests, which run
// through NewTestHarness — the harness deploys the CapabilitiesRegistry itself (and
// asserts exactly one exists), so those tests must NOT seed one.
func seedStellarOCR3(t *testing.T, ds *datastore.MemoryDataStore, ocrSel uint64) {
	t.Helper()
	require.NoError(t, ds.Addresses().Add(datastore.AddressRef{
		ChainSelector: ocrSel,
		Type:          datastore.ContractType(ocr3.OCR3Capability),
		Version:       semver.MustParse("1.0.0"),
		Address:       "0x1111111111111111111111111111111111111111",
		Qualifier:     testStellarOCRQualifier,
	}))
}

// seedStellarAddresses seeds the OCR3 contract + CapabilitiesRegistry + forwarders for the
// VerifyPreconditions tests, which run against a bare MemoryDataStore (no harness), so
// resolveCapRegAddress needs the registry present. No forwarder address is seeded: the
// Stellar cap proposal no longer resolves a deployed forwarder from the datastore.
func seedStellarAddresses(t *testing.T, ds *datastore.MemoryDataStore, ocrSel, stellarSel uint64) {
	t.Helper()
	seedStellarOCR3(t, ds, ocrSel)
	require.NoError(t, ds.Addresses().Add(datastore.AddressRef{
		ChainSelector: ocrSel,
		Type:          datastore.ContractType("CapabilitiesRegistry"),
		Version:       semver.MustParse("2.0.0"),
		Address:       "0x2222222222222222222222222222222222222222",
		Qualifier:     testStellarOCRQualifier,
	}))
	require.NoError(t, ds.Addresses().Add(datastore.AddressRef{
		ChainSelector: stellarSel,
		Type:          stellar.ForwarderContract,
		Version:       semver.MustParse("1.0.0"),
		Address:       testStellarForwarderAddress,
		Qualifier:     testStellarFwdQualifier,
	}))
}

func freshStellarBase(ocrSel, stellarSel uint64) jobs.ProposeStellarCapJobSpecInput {
	return jobs.ProposeStellarCapJobSpecInput{
		Environment:          test.EnvironmentName,
		Zone:                 test.Zone,
		Domain:               "cre",
		DONName:              test.DONName,
		ChainSelector:        stellarSel,
		OCRChainSelector:     ocrSel,
		BootstrapperOCR3Urls: []string{"12D3KooWxyz@127.0.0.1:5001"},
		OCRContractQualifier: testStellarOCRQualifier,
		ForwardersQualifier:  testStellarFwdQualifier,
		DeltaStage:           10 * time.Second,
		StellarCapabilityInputs: []jobs.StellarCapabilityInput{
			minimalStellarCapInput("peer-1"),
		},
	}
}

func deepCloneStellarInput(in jobs.ProposeStellarCapJobSpecInput) jobs.ProposeStellarCapJobSpecInput {
	clone := in
	if len(in.StellarCapabilityInputs) > 0 {
		clone.StellarCapabilityInputs = append([]jobs.StellarCapabilityInput(nil), in.StellarCapabilityInputs...)
	}
	return clone
}

func TestProposeStellarCapJobSpec_VerifyPreconditions_success(t *testing.T) {
	var env cldf.Environment

	ocrSel := chainsel.ETHEREUM_TESTNET_SEPOLIA.Selector
	stellarSel := chainsel.STELLAR_LOCALNET.Selector

	ds := datastore.NewMemoryDataStore()
	seedStellarAddresses(t, ds, ocrSel, stellarSel)
	env.DataStore = ds.Seal()

	in := freshStellarBase(ocrSel, stellarSel)
	err := jobs.ProposeStellarCapJobSpec{}.VerifyPreconditions(env, in)
	require.NoError(t, err)
}

func TestProposeStellarCapJobSpec_VerifyPreconditions_requiredFields(t *testing.T) {
	var env cldf.Environment

	ocrSel := chainsel.ETHEREUM_TESTNET_SEPOLIA.Selector
	stellarSel := chainsel.STELLAR_LOCALNET.Selector

	ds := datastore.NewMemoryDataStore()
	seedStellarAddresses(t, ds, ocrSel, stellarSel)
	env.DataStore = ds.Seal()

	base := freshStellarBase(ocrSel, stellarSel)

	cases := []struct {
		name    string
		mutate  func(*jobs.ProposeStellarCapJobSpecInput)
		errFrag string
	}{
		{"missing environment", func(in *jobs.ProposeStellarCapJobSpecInput) { in.Environment = "" }, "environment is required"},
		{"missing domain", func(in *jobs.ProposeStellarCapJobSpecInput) { in.Domain = "" }, "domain is required"},
		{"missing zone", func(in *jobs.ProposeStellarCapJobSpecInput) { in.Zone = "" }, "zone is required"},
		{"missing don name", func(in *jobs.ProposeStellarCapJobSpecInput) { in.DONName = "" }, "donName is required"},
		{"missing chain selector", func(in *jobs.ProposeStellarCapJobSpecInput) { in.ChainSelector = 0 }, "chain selector is required"},
		{"missing stellar inputs", func(in *jobs.ProposeStellarCapJobSpecInput) { in.StellarCapabilityInputs = nil }, "at least one stellar capability input is required"},
		{"missing node id", func(in *jobs.ProposeStellarCapJobSpecInput) { in.StellarCapabilityInputs[0].NodeID = "" }, "nodeID is required for stellar capability input"},
		{"missing delta stage", func(in *jobs.ProposeStellarCapJobSpecInput) { in.DeltaStage = 0 }, "deltaStage"},
		{"missing forwarder qualifier", func(in *jobs.ProposeStellarCapJobSpecInput) { in.ForwardersQualifier = "" }, "cre forwarder qualifier is required"},
		{"wrong chain family", func(in *jobs.ProposeStellarCapJobSpecInput) {
			in.ChainSelector = chainsel.ETHEREUM_TESTNET_SEPOLIA.Selector
		}, "expected \"stellar\""},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			in := deepCloneStellarInput(base)
			tc.mutate(&in)
			err := jobs.ProposeStellarCapJobSpec{}.VerifyPreconditions(env, in)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tc.errFrag)
		})
	}
}

func TestProposeStellarCapJobSpec_VerifyPreconditions_missingDatastore(t *testing.T) {
	var env cldf.Environment

	ocrSel := chainsel.ETHEREUM_TESTNET_SEPOLIA.Selector
	stellarSel := chainsel.STELLAR_LOCALNET.Selector

	env.DataStore = datastore.NewMemoryDataStore().Seal()
	in := freshStellarBase(ocrSel, stellarSel)

	err := jobs.ProposeStellarCapJobSpec{}.VerifyPreconditions(env, in)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "CapabilitiesRegistry")
}

type stellarCapTestSetup struct {
	rt        *runtime.Runtime
	baseInput jobs.ProposeStellarCapJobSpecInput
}

func setupStellarCapTest(t *testing.T) stellarCapTestSetup {
	t.Helper()

	ds := datastore.NewMemoryDataStore()
	seedStellarAddresses(t, ds, chainsel.ETHEREUM_TESTNET_SEPOLIA.Selector, chainsel.STELLAR_LOCALNET.Selector)

	h := test.NewTestHarness(t, test.WithDatastore(ds))
	env := h.Runtime.Environment()

	nodes, err := h.TestJD.ListNodes(t.Context(), &node.ListNodesRequest{})
	require.NoError(t, err)

	var stellarCapInputs []jobs.StellarCapabilityInput
	mockGetter := &tenv.MockJobApproverGetter{JobApprovers: make(map[string]*tenv.MockJobApprover)}
	for _, n := range nodes.GetNodes() {
		if strings.Contains(n.Id, "bootstrap") {
			continue
		}
		mockGetter.JobApprovers[n.Id] = &tenv.MockJobApprover{}
		stellarCapInputs = append(stellarCapInputs, minimalStellarCapInput(n.Id))
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

	rt := runtime.NewFromEnvironment(env)

	baseInput := jobs.ProposeStellarCapJobSpecInput{
		Environment:             test.EnvironmentName,
		Zone:                    test.Zone,
		Domain:                  "cre",
		DONName:                 test.DONName,
		ChainSelector:           chainsel.STELLAR_LOCALNET.Selector,
		OCRChainSelector:        h.RegistrySelector,
		BootstrapperOCR3Urls:    []string{"12D3KooWabc@127.0.0.1:5001"},
		OCRContractQualifier:    test.RegistryQualifier,
		ForwardersQualifier:     testStellarFwdQualifier,
		DeltaStage:              time.Second,
		StellarCapabilityInputs: stellarCapInputs,
	}

	return stellarCapTestSetup{rt: rt, baseInput: baseInput}
}

func TestProposeStellarCapJobSpec_Apply_success(t *testing.T) {
	setup := setupStellarCapTest(t)
	task := runtime.ChangesetTask(jobs.ProposeStellarCapJobSpec{}, setup.baseInput)
	err := setup.rt.Exec(task)
	require.NoError(t, err)

	out := setup.rt.State().Outputs[task.ID()]
	assert.Len(t, out.Reports, 1)
}

func TestProposeStellarCapJobSpec_Apply_forwarderLookbackLedgers(t *testing.T) {
	setup := setupStellarCapTest(t)

	const (
		inputLookback  int64 = 123 // DON-wide value
		overrideCustom int64 = 999 // explicit per-node override
	)

	input := setup.baseInput
	input.ForwarderLookbackLedgers = inputLookback

	nodeCount := len(input.StellarCapabilityInputs)
	require.GreaterOrEqual(t, nodeCount, 2, "need at least 2 nodes to distinguish override from default")

	// Explicit per-node override on the first node only; the rest inherit the DON-wide value.
	input.StellarCapabilityInputs[0].OverrideDefaultCfg.ForwarderLookbackLedgers = overrideCustom

	task := runtime.ChangesetTask(jobs.ProposeStellarCapJobSpec{}, input)
	require.NoError(t, setup.rt.Exec(task))

	out := setup.rt.State().Outputs[task.ID()]
	require.NotNil(t, out)
	require.Len(t, out.Reports, 1)

	outputStr := fmt.Sprintf("%v", out.Reports[0].Output)
	assert.Equal(t, 1, strings.Count(outputStr, `"forwarderLookbackLedgers":999`),
		"expected exactly one per-node override to be preserved")
	assert.Equal(t, nodeCount-1, strings.Count(outputStr, `"forwarderLookbackLedgers":123`),
		"expected every other node to inherit the DON-wide value")
}

// With neither the DON-wide input nor any per-node override set, the field is
// omitted entirely and the worker applies actions.DefaultForwarderLookbackLedgers.
// The changeset deliberately does not duplicate that default — same division of
// responsibility as EVM's ForwarderLookbackBlocks.
func TestProposeStellarCapJobSpec_Apply_forwarderLookbackLedgersOmittedWhenUnset(t *testing.T) {
	setup := setupStellarCapTest(t)

	input := setup.baseInput
	input.ForwarderLookbackLedgers = 0
	for i := range input.StellarCapabilityInputs {
		input.StellarCapabilityInputs[i].OverrideDefaultCfg.ForwarderLookbackLedgers = 0
	}

	task := runtime.ChangesetTask(jobs.ProposeStellarCapJobSpec{}, input)
	require.NoError(t, setup.rt.Exec(task))

	out := setup.rt.State().Outputs[task.ID()]
	require.NotNil(t, out)
	require.Len(t, out.Reports, 1)

	outputStr := fmt.Sprintf("%v", out.Reports[0].Output)
	assert.NotContains(t, outputStr, `"forwarderLookbackLedgers"`,
		"omitempty should drop the field entirely so the worker's own default applies")
}

func TestProposeStellarCapJobSpec_Apply_duplicateNodeIDs(t *testing.T) {
	setup := setupStellarCapTest(t)
	env := setup.rt.Environment()

	input := setup.baseInput
	require.GreaterOrEqual(t, len(input.StellarCapabilityInputs), 1, "need at least 1 node")
	input.StellarCapabilityInputs = []jobs.StellarCapabilityInput{
		input.StellarCapabilityInputs[0],
		input.StellarCapabilityInputs[0],
	}

	_, err := jobs.ProposeStellarCapJobSpec{}.Apply(env, input)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "duplicate nodeID")
}
