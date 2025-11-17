package jobs

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"

	"github.com/smartcontractkit/chainlink/deployment/cre/test"

	"github.com/smartcontractkit/chainlink/deployment/cre/ocr3"
)

const (
	// ⚠️ Adjust this if your pkg.GetKeystoneForwarderCapabilityAddressRefKey
	// encodes a different ContractType in its key.
	testForwarderContractType = datastore.ContractType("KeystoneForwarder")

	testOCRQualifier       = "OCR3Capability"
	testForwarderQualifier = "forwarder-qualifier"
)

// helper to build a minimal valid input (no user overrides; Apply will fill defaults)
func minimalEVMCapInput(nodeID string) EVMCapabilityInput {
	return EVMCapabilityInput{
		NodeID:             nodeID,
		OverrideDefaultCfg: OverrideDefaultCfg{},
	}
}

func deepCloneInput(in ProposeEVMCapJobSpecInput) ProposeEVMCapJobSpecInput {
	clone := in
	if len(in.EVMCapabilityInputs) > 0 {
		clone.EVMCapabilityInputs = append([]EVMCapabilityInput(nil), in.EVMCapabilityInputs...)
	}
	return clone
}

func freshBase(selector uint64) ProposeEVMCapJobSpecInput {
	return ProposeEVMCapJobSpecInput{
		Environment:          "test",
		Zone:                 test.Zone,
		Domain:               "cre",
		DONName:              test.DONName,
		ChainSelector:        selector,
		BootstrapperOCR3Urls: []string{"12D3KooWxyz@127.0.0.1:5001"},
		OCRContractQualifier: testOCRQualifier,
		ForwardersQualifier:  testForwarderQualifier,
		EVMCapabilityInputs:  []EVMCapabilityInput{minimalEVMCapInput("peer-1")},
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

	in := ProposeEVMCapJobSpecInput{
		Environment:          "test",
		Zone:                 test.Zone,
		Domain:               "cre",
		DONName:              test.DONName,
		ChainSelector:        chain.Selector,
		BootstrapperOCR3Urls: []string{"12D3KooWxyz@127.0.0.1:5001"},
		OCRContractQualifier: testOCRQualifier,
		ForwardersQualifier:  testForwarderQualifier,
		EVMCapabilityInputs: []EVMCapabilityInput{
			minimalEVMCapInput("peer-1"),
			minimalEVMCapInput("peer-2"),
		},
	}

	err := ProposeEVMCapJobSpec{}.VerifyPreconditions(env, in)
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
		mutate  func(*ProposeEVMCapJobSpecInput)
		errFrag string
	}{
		{"missing environment", func(in *ProposeEVMCapJobSpecInput) { in.Environment = "" }, "environment is required"},
		{"missing domain", func(in *ProposeEVMCapJobSpecInput) { in.Domain = "" }, "domain is required"},
		{"missing don name", func(in *ProposeEVMCapJobSpecInput) { in.DONName = "" }, "donName is required"},
		{"missing zone", func(in *ProposeEVMCapJobSpecInput) { in.Zone = "" }, "zone is required"},
		{"missing chain selector", func(in *ProposeEVMCapJobSpecInput) { in.ChainSelector = 0 }, "chain selector is required"},
		{"missing evm inputs", func(in *ProposeEVMCapJobSpecInput) { in.EVMCapabilityInputs = nil }, "at least one evm capability input is required"},
		{"missing bootstrapper urls", func(in *ProposeEVMCapJobSpecInput) { in.BootstrapperOCR3Urls = nil }, "at least one bootstrapper OCR3 URL is required"},
		{"missing OCR qualifier", func(in *ProposeEVMCapJobSpecInput) { in.OCRContractQualifier = "" }, "ocr contract qualifier is required"},
		{"missing forwarder qualifier", func(in *ProposeEVMCapJobSpecInput) { in.ForwardersQualifier = "" }, "cre forwarder qualifier is required"},
		{"missing node id", func(in *ProposeEVMCapJobSpecInput) { in.EVMCapabilityInputs[0].NodeID = "" }, "nodeID in evm capability input is required"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			in := deepCloneInput(base)
			tc.mutate(&in)
			err := ProposeEVMCapJobSpec{}.VerifyPreconditions(env, in)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tc.errFrag)
		})
	}
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

		err := ProposeEVMCapJobSpec{}.VerifyPreconditions(env, in)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "chainID in override config")
	})

	t.Run("network must be evm if provided", func(t *testing.T) {
		in := deepCloneInput(base)
		in.EVMCapabilityInputs[0].OverrideDefaultCfg.Network = "solana"
		err := ProposeEVMCapJobSpec{}.VerifyPreconditions(env, in)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "network in override config must be")
	})

	t.Run("forwarder address mismatch when provided", func(t *testing.T) {
		in := deepCloneInput(base)
		in.EVMCapabilityInputs[0].OverrideDefaultCfg.CREForwarderAddress = "0xdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef"
		err := ProposeEVMCapJobSpec{}.VerifyPreconditions(env, in)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "CRE forwarder address in override config")
	})

	t.Run("below-minimum values are rejected, zeros allowed", func(t *testing.T) {
		// Zeros OK (treated as "use defaults"), so verify passes:
		in := deepCloneInput(base)
		require.NoError(t, ProposeEVMCapJobSpec{}.VerifyPreconditions(env, in))

		// Now set below-minimums independently each time
		in = deepCloneInput(base)
		in.EVMCapabilityInputs[0].OverrideDefaultCfg.LogTriggerPollInterval = 1499 * time.Millisecond // 1ms below min
		err := ProposeEVMCapJobSpec{}.VerifyPreconditions(env, in)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "logTriggerPollInterval")

		in = deepCloneInput(base)
		in.EVMCapabilityInputs[0].OverrideDefaultCfg.ReceiverGasMinimum = 499
		err = ProposeEVMCapJobSpec{}.VerifyPreconditions(env, in)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "receiverGasMinimum")

		in = deepCloneInput(base)
		in.EVMCapabilityInputs[0].OverrideDefaultCfg.LogTriggerSendChannelBufferSize = 2999
		err = ProposeEVMCapJobSpec{}.VerifyPreconditions(env, in)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "logTriggerSendChannelBufferSize")
	})
}

func TestProposeEVMCapJobSpec_Apply_success(t *testing.T) {
	testEnv := test.SetupEnvV2(t, false)
	env := testEnv.Env

	allNodes, err := testEnv.TestJD.ListNodes(t.Context(), &node.ListNodesRequest{})
	require.NoError(t, err)

	fmt.Println("allNodes:", allNodes)
	selector := testEnv.RegistrySelector // use the test environment's selector
	ds := datastore.NewMemoryDataStore()

	fmt.Println("selector:", selector)
	// Seed required addresses (OCR for VerifyPreconditions; forwarder for both Verify & Apply)
	require.NoError(t, ds.Addresses().Add(datastore.AddressRef{
		ChainSelector: selector,
		Type:          datastore.ContractType(ocr3.OCR3Capability),
		Version:       semver.MustParse("1.0.0"),
		Address:       "0xAb5801a7D398351b8bE11C439e05C5B3259aeC9B",
		Qualifier:     testOCRQualifier,
	}))
	require.NoError(t, ds.Addresses().Add(datastore.AddressRef{
		ChainSelector: selector,
		Type:          testForwarderContractType,
		Version:       semver.MustParse("1.0.0"),
		Address:       "0x9fE46736679d2D9a65F0992F2272dE9f3c7fa6e0",
		Qualifier:     testForwarderQualifier,
	}))

	env.DataStore = ds.Seal()

	input := ProposeEVMCapJobSpecInput{
		Environment:          "test",
		Zone:                 test.Zone,
		Domain:               "cre",
		DONName:              test.DONName,
		ChainSelector:        selector,
		BootstrapperOCR3Urls: []string{"12D3KooWabc@127.0.0.1:5001"},
		OCRContractQualifier: testOCRQualifier,
		ForwardersQualifier:  testForwarderQualifier,
		EVMCapabilityInputs: []EVMCapabilityInput{
			// leave values zero to exercise defaults; Apply should fill them
			minimalEVMCapInput("node_test-don-0"),
			minimalEVMCapInput("node_test-don-1"),
			minimalEVMCapInput("node_test-don-2"),
			minimalEVMCapInput("node_test-don-3"),
		},
	}

	// Verify should pass
	require.NoError(t, ProposeEVMCapJobSpec{}.VerifyPreconditions(*env, input))

	// Apply should succeed and return one report
	out, err := ProposeEVMCapJobSpec{}.Apply(*env, input)
	require.NoError(t, err)
	assert.Len(t, out.Reports, 1)

	fmt.Println("out.Reports:", out.Reports)
}
