package jobs_test

import (
	"strings"
	"testing"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/deployment/cre/jobs"
	"github.com/smartcontractkit/chainlink/deployment/cre/ocr3"
	"github.com/smartcontractkit/chainlink/deployment/cre/test"
	tenv "github.com/smartcontractkit/chainlink/deployment/environment/test"

	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	csav1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/csa"
	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"
	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"
)

const (
	testAptosForwarderQualifier = "aptos-forwarder-qualifier"
	testAptosOCRQualifier       = "aptos-ocr-qualifier"
)

func minimalAptosCapInput(nodeID string) jobs.AptosCapabilityInput {
	return jobs.AptosCapabilityInput{
		NodeID:             nodeID,
		OverrideDefaultCfg: jobs.AptosOverrideDefaultCfg{},
	}
}

func seedAptosAddresses(t *testing.T, ds *datastore.MemoryDataStore, ocrSel, aptosSel uint64, ocrAddr, fwdAddr string) {
	t.Helper()
	require.NoError(t, ds.Addresses().Add(datastore.AddressRef{
		ChainSelector: ocrSel,
		Type:          datastore.ContractType(ocr3.OCR3Capability),
		Version:       semver.MustParse("1.0.0"),
		Address:       ocrAddr,
		Qualifier:     testAptosOCRQualifier,
	}))
	require.NoError(t, ds.Addresses().Add(datastore.AddressRef{
		ChainSelector: aptosSel,
		Type:          testForwarderContractType,
		Version:       semver.MustParse("1.0.0"),
		Address:       fwdAddr,
		Qualifier:     testAptosForwarderQualifier,
	}))
}

func freshAptosBase(ocrSel, aptosSel uint64) jobs.ProposeAptosCapJobSpecInput {
	return jobs.ProposeAptosCapJobSpecInput{
		Environment:          "test",
		Zone:                 test.Zone,
		Domain:               "cre",
		DONName:              test.DONName,
		ChainSelector:        aptosSel,
		OCRChainSelector:     ocrSel,
		BootstrapperOCR3Urls: []string{"12D3KooWxyz@127.0.0.1:5001"},
		OCRContractQualifier: testAptosOCRQualifier,
		ForwardersQualifier:  testAptosForwarderQualifier,
		DeltaStage:           10 * time.Second,
		AptosCapabilityInputs: []jobs.AptosCapabilityInput{
			minimalAptosCapInput("peer-1"),
		},
	}
}

func deepCloneAptosInput(in jobs.ProposeAptosCapJobSpecInput) jobs.ProposeAptosCapJobSpecInput {
	clone := in
	if len(in.AptosCapabilityInputs) > 0 {
		clone.AptosCapabilityInputs = append([]jobs.AptosCapabilityInput(nil), in.AptosCapabilityInputs...)
	}
	return clone
}

func TestProposeAptosCapJobSpec_VerifyPreconditions_success(t *testing.T) {
	var env cldf.Environment

	ocrSel := chainsel.ETHEREUM_TESTNET_SEPOLIA.Selector
	aptosSel := chainsel.APTOS_TESTNET.Selector

	ds := datastore.NewMemoryDataStore()
	seedAptosAddresses(t, ds, ocrSel, aptosSel,
		"0x1111111111111111111111111111111111111111",
		"0x2222222222222222222222222222222222222222222222222222222222222222",
	)
	env.DataStore = ds.Seal()

	in := freshAptosBase(ocrSel, aptosSel)
	in.AptosCapabilityInputs = []jobs.AptosCapabilityInput{
		minimalAptosCapInput("peer-1"),
		minimalAptosCapInput("peer-2"),
	}

	err := jobs.ProposeAptosCapJobSpec{}.VerifyPreconditions(env, in)
	require.NoError(t, err)
}

func TestProposeAptosCapJobSpec_VerifyPreconditions_requiredFields(t *testing.T) {
	var env cldf.Environment

	ocrSel := chainsel.ETHEREUM_TESTNET_SEPOLIA.Selector
	aptosSel := chainsel.APTOS_TESTNET.Selector

	ds := datastore.NewMemoryDataStore()
	seedAptosAddresses(t, ds, ocrSel, aptosSel,
		"0x1111111111111111111111111111111111111111",
		"0x2222222222222222222222222222222222222222222222222222222222222222",
	)
	env.DataStore = ds.Seal()

	base := freshAptosBase(ocrSel, aptosSel)

	cases := []struct {
		name    string
		mutate  func(*jobs.ProposeAptosCapJobSpecInput)
		errFrag string
	}{
		{"missing environment", func(in *jobs.ProposeAptosCapJobSpecInput) { in.Environment = "" }, "environment is required"},
		{"missing domain", func(in *jobs.ProposeAptosCapJobSpecInput) { in.Domain = "" }, "domain is required"},
		{"missing zone", func(in *jobs.ProposeAptosCapJobSpecInput) { in.Zone = "" }, "zone is required"},
		{"missing don name", func(in *jobs.ProposeAptosCapJobSpecInput) { in.DONName = "" }, "donName is required"},
		{"missing chain selector", func(in *jobs.ProposeAptosCapJobSpecInput) { in.ChainSelector = 0 }, "chain selector is required"},
		{"missing ocr chain selector", func(in *jobs.ProposeAptosCapJobSpecInput) { in.OCRChainSelector = 0 }, "ocr chain selector is required"},
		{"missing aptos inputs", func(in *jobs.ProposeAptosCapJobSpecInput) { in.AptosCapabilityInputs = nil }, "at least one aptos capability input is required"},
		{"missing bootstrapper urls", func(in *jobs.ProposeAptosCapJobSpecInput) { in.BootstrapperOCR3Urls = nil }, "at least one bootstrapper OCR3 URL is required"},
		{"empty bootstrapper url element", func(in *jobs.ProposeAptosCapJobSpecInput) { in.BootstrapperOCR3Urls = []string{""} }, "bootstrapper OCR3 URL at index 0 is empty"},
		{"missing OCR qualifier", func(in *jobs.ProposeAptosCapJobSpecInput) { in.OCRContractQualifier = "" }, "ocr contract qualifier is required"},
		{"missing forwarder qualifier", func(in *jobs.ProposeAptosCapJobSpecInput) { in.ForwardersQualifier = "" }, "cre forwarder qualifier is required"},
		{"missing node id", func(in *jobs.ProposeAptosCapJobSpecInput) { in.AptosCapabilityInputs[0].NodeID = "" }, "nodeID is required for aptos capability input"},
		{"missing delta stage", func(in *jobs.ProposeAptosCapJobSpecInput) { in.DeltaStage = 0 }, "deltaStage"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			in := deepCloneAptosInput(base)
			tc.mutate(&in)
			err := jobs.ProposeAptosCapJobSpec{}.VerifyPreconditions(env, in)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tc.errFrag)
		})
	}
}

func TestProposeAptosCapJobSpec_VerifyPreconditions_missingAddresses(t *testing.T) {
	var env cldf.Environment

	ocrSel := chainsel.ETHEREUM_TESTNET_SEPOLIA.Selector
	aptosSel := chainsel.APTOS_TESTNET.Selector

	t.Run("missing OCR address", func(t *testing.T) {
		ds := datastore.NewMemoryDataStore()
		require.NoError(t, ds.Addresses().Add(datastore.AddressRef{
			ChainSelector: aptosSel,
			Type:          testForwarderContractType,
			Version:       semver.MustParse("1.0.0"),
			Address:       "0x2222222222222222222222222222222222222222222222222222222222222222",
			Qualifier:     testAptosForwarderQualifier,
		}))
		env.DataStore = ds.Seal()

		in := freshAptosBase(ocrSel, aptosSel)
		err := jobs.ProposeAptosCapJobSpec{}.VerifyPreconditions(env, in)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get OCR contract address")
	})

	t.Run("missing forwarder address", func(t *testing.T) {
		ds := datastore.NewMemoryDataStore()
		require.NoError(t, ds.Addresses().Add(datastore.AddressRef{
			ChainSelector: ocrSel,
			Type:          datastore.ContractType(ocr3.OCR3Capability),
			Version:       semver.MustParse("1.0.0"),
			Address:       "0x1111111111111111111111111111111111111111",
			Qualifier:     testAptosOCRQualifier,
		}))
		env.DataStore = ds.Seal()

		in := freshAptosBase(ocrSel, aptosSel)
		err := jobs.ProposeAptosCapJobSpec{}.VerifyPreconditions(env, in)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get CRE forwarder address")
	})
}

func TestProposeAptosCapJobSpec_VerifyPreconditions_overrideMismatches(t *testing.T) {
	var env cldf.Environment

	ocrSel := chainsel.ETHEREUM_TESTNET_SEPOLIA.Selector
	aptosSel := chainsel.APTOS_TESTNET.Selector
	fwdAddr := "0x2222222222222222222222222222222222222222222222222222222222222222"

	ds := datastore.NewMemoryDataStore()
	seedAptosAddresses(t, ds, ocrSel, aptosSel,
		"0x1111111111111111111111111111111111111111",
		fwdAddr,
	)
	env.DataStore = ds.Seal()

	base := freshAptosBase(ocrSel, aptosSel)

	t.Run("chainID mismatch when provided", func(t *testing.T) {
		in := deepCloneAptosInput(base)
		in.AptosCapabilityInputs[0].OverrideDefaultCfg.ChainID = "999999"
		err := jobs.ProposeAptosCapJobSpec{}.VerifyPreconditions(env, in)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "chainID in override config")
	})

	t.Run("matching chainID is accepted", func(t *testing.T) {
		in := deepCloneAptosInput(base)
		chainIDStr, err := chainsel.GetChainIDFromSelector(aptosSel)
		require.NoError(t, err)
		in.AptosCapabilityInputs[0].OverrideDefaultCfg.ChainID = chainIDStr
		require.NoError(t, jobs.ProposeAptosCapJobSpec{}.VerifyPreconditions(env, in))
	})

	t.Run("network must be aptos if provided", func(t *testing.T) {
		in := deepCloneAptosInput(base)
		in.AptosCapabilityInputs[0].OverrideDefaultCfg.Network = "evm"
		err := jobs.ProposeAptosCapJobSpec{}.VerifyPreconditions(env, in)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "network in override config must be")
	})

	t.Run("matching network is accepted", func(t *testing.T) {
		in := deepCloneAptosInput(base)
		in.AptosCapabilityInputs[0].OverrideDefaultCfg.Network = "aptos"
		require.NoError(t, jobs.ProposeAptosCapJobSpec{}.VerifyPreconditions(env, in))
	})

	t.Run("forwarder address mismatch when provided", func(t *testing.T) {
		in := deepCloneAptosInput(base)
		in.AptosCapabilityInputs[0].OverrideDefaultCfg.CREForwarderAddress = "0xdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef"
		err := jobs.ProposeAptosCapJobSpec{}.VerifyPreconditions(env, in)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "CRE forwarder address in override config")
	})

	t.Run("matching forwarder address is accepted", func(t *testing.T) {
		in := deepCloneAptosInput(base)
		in.AptosCapabilityInputs[0].OverrideDefaultCfg.CREForwarderAddress = fwdAddr
		require.NoError(t, jobs.ProposeAptosCapJobSpec{}.VerifyPreconditions(env, in))
	})
}

type aptosCapTestSetup struct {
	env            *cldf.Environment
	nodeIDs        []string
	aptosCapInputs []jobs.AptosCapabilityInput
	baseInput      jobs.ProposeAptosCapJobSpecInput
}

func setupAptosCapTest(t *testing.T) aptosCapTestSetup {
	t.Helper()
	testEnv := test.SetupEnvV2(t, false)

	ocrSel := testEnv.RegistrySelector
	aptosSel := testEnv.AptosSelector

	ds := datastore.NewMemoryDataStore()
	seedAptosAddresses(t, ds, ocrSel, aptosSel,
		"0x1111111111111111111111111111111111111111",
		"0x2222222222222222222222222222222222222222222222222222222222222222",
	)
	env := testEnv.Env
	env.DataStore = ds.Seal()

	nodes, err := testEnv.TestJD.ListNodes(t.Context(), &node.ListNodesRequest{})
	require.NoError(t, err)

	var nodeIDs []string
	var aptosCapInputs []jobs.AptosCapabilityInput
	mockGetter := &tenv.MockJobApproverGetter{JobApprovers: make(map[string]*tenv.MockJobApprover)}
	for _, n := range nodes.GetNodes() {
		if strings.Contains(n.Id, "bootstrap") {
			continue
		}
		nodeIDs = append(nodeIDs, n.Id)
		mockGetter.JobApprovers[n.Id] = &tenv.MockJobApprover{}
		aptosCapInputs = append(aptosCapInputs, minimalAptosCapInput(n.Id))
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

	baseInput := jobs.ProposeAptosCapJobSpecInput{
		Environment:            "test",
		Zone:                   test.Zone,
		Domain:                 "cre",
		DONName:                test.DONName,
		ChainSelector:          aptosSel,
		OCRChainSelector:       ocrSel,
		BootstrapperOCR3Urls:   []string{"12D3KooWabc@127.0.0.1:5001"},
		OCRContractQualifier:   testAptosOCRQualifier,
		ForwardersQualifier:    testAptosForwarderQualifier,
		DeltaStage:             time.Second,
		TxSearchStartingBuffer: 30 * time.Second,
		AptosCapabilityInputs:  aptosCapInputs,
	}

	return aptosCapTestSetup{
		env:            env,
		nodeIDs:        nodeIDs,
		aptosCapInputs: aptosCapInputs,
		baseInput:      baseInput,
	}
}

func TestProposeAptosCapJobSpec_Apply_success(t *testing.T) {
	setup := setupAptosCapTest(t)
	env := setup.env

	input := setup.baseInput

	require.NoError(t, jobs.ProposeAptosCapJobSpec{}.VerifyPreconditions(*env, input))

	out, err := jobs.ProposeAptosCapJobSpec{}.Apply(*env, input)
	require.NoError(t, err)
	assert.Len(t, out.Reports, 1)
}

func TestProposeAptosCapJobSpec_Apply_duplicateNodeIDs(t *testing.T) {
	setup := setupAptosCapTest(t)
	env := setup.env

	input := setup.baseInput
	require.GreaterOrEqual(t, len(setup.aptosCapInputs), 2, "need at least 2 nodes")
	input.AptosCapabilityInputs = []jobs.AptosCapabilityInput{
		setup.aptosCapInputs[0],
		setup.aptosCapInputs[0],
	}

	_, err := jobs.ProposeAptosCapJobSpec{}.Apply(*env, input)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "duplicate nodeID")
}
