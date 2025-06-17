package changeset_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/internal"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/test"
)

func TestDeployOCR3(t *testing.T) {
	t.Parallel()
	lggr := logger.Test(t)
	cfg := memory.MemoryEnvironmentConfig{
		Nodes:  1, // nodes unused but required in config
		Chains: 2,
	}
	env := memory.NewMemoryEnvironment(t, lggr, zapcore.DebugLevel, cfg)

	registrySel := env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))[0]

	resp, err := changeset.DeployOCR3(env, registrySel)
	require.NoError(t, err)
	require.NotNil(t, resp)
	// OCR3 should be deployed on chain 0
	addrs, err := resp.AddressBook.AddressesForChain(registrySel)
	require.NoError(t, err)
	require.Len(t, addrs, 1)

	// nothing on chain 1
	require.NotEqual(t, registrySel, env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))[1])
	oaddrs, _ := resp.AddressBook.AddressesForChain(env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))[1])
	assert.Empty(t, oaddrs)
}

func TestConfigureOCR3(t *testing.T) {
	t.Parallel()

	nWfNodes := 4
	c := internal.OracleConfig{
		MaxFaultyOracles:     1,
		DeltaProgressMillis:  12345,
		TransmissionSchedule: []int{nWfNodes},
	}

	t.Run("no mcms", func(t *testing.T) {
		te := test.SetupContractTestEnv(t, test.EnvWrapperConfig{
			WFDonConfig:     test.DonConfig{Name: "wfDon", N: nWfNodes},
			AssetDonConfig:  test.DonConfig{Name: "assetDon", N: 4},
			WriterDonConfig: test.DonConfig{Name: "writerDon", N: 4},
			NumChains:       1,
		})

		wfNodes := te.GetP2PIDs("wfDon").Strings()

		w := &bytes.Buffer{}
		cfg := changeset.ConfigureOCR3Config{
			ChainSel:             te.RegistrySelector,
			NodeIDs:              wfNodes,
			OCR3Config:           &c,
			WriteGeneratedConfig: w,
		}

		csOut, err := changeset.ConfigureOCR3Contract(te.Env, cfg)
		require.NoError(t, err)
		var got internal.OCR2OracleConfig
		err = json.Unmarshal(w.Bytes(), &got)
		require.NoError(t, err)
		assert.Len(t, got.Signers, 4)
		assert.Len(t, got.Transmitters, 4)
		assert.Nil(t, csOut.MCMSTimelockProposals)
	})

	t.Run("success multiple OCR3 contracts", func(t *testing.T) {
		te := test.SetupContractTestEnv(t, test.EnvWrapperConfig{
			WFDonConfig:     test.DonConfig{Name: "wfDon", N: nWfNodes},
			AssetDonConfig:  test.DonConfig{Name: "assetDon", N: 4},
			WriterDonConfig: test.DonConfig{Name: "writerDon", N: 4},
			NumChains:       1,
		})

		registrySel := te.Env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))[0]

		existingContracts, err := te.Env.ExistingAddresses.AddressesForChain(registrySel)
		require.NoError(t, err)
		require.Len(t, existingContracts, 4)

		// Find existing OCR3 contract
		var existingOCR3Addr string
		for addr, tv := range existingContracts {
			if tv.Type == internal.OCR3Capability {
				existingOCR3Addr = addr
				break
			}
		}

		// Deploy a new OCR3 contract
		resp, err := changeset.DeployOCR3(te.Env, registrySel)
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.NoError(t, te.Env.ExistingAddresses.Merge(resp.AddressBook))

		// Verify after merge there are three original contracts plus one new one
		addrs, err := te.Env.ExistingAddresses.AddressesForChain(registrySel)
		require.NoError(t, err)
		require.Len(t, addrs, 5)

		// Find new OCR3 contract
		var newOCR3Addr string
		for addr, tv := range addrs {
			if tv.Type == internal.OCR3Capability && addr != existingOCR3Addr {
				newOCR3Addr = addr
				break
			}
		}

		wfNodes := te.GetP2PIDs("wfDon").Strings()

		na := common.HexToAddress(newOCR3Addr)
		w := &bytes.Buffer{}
		cfg := changeset.ConfigureOCR3Config{
			ChainSel:             te.RegistrySelector,
			NodeIDs:              wfNodes,
			Address:              &na, // Use the new OCR3 contract to configure
			OCR3Config:           &c,
			WriteGeneratedConfig: w,
		}

		csOut, err := changeset.ConfigureOCR3Contract(te.Env, cfg)
		require.NoError(t, err)
		var got internal.OCR2OracleConfig
		err = json.Unmarshal(w.Bytes(), &got)
		require.NoError(t, err)
		assert.Len(t, got.Signers, 4)
		assert.Len(t, got.Transmitters, 4)
		assert.Nil(t, csOut.MCMSTimelockProposals)
	})

	t.Run("fails multiple OCR3 contracts but unspecified address", func(t *testing.T) {
		te := test.SetupContractTestEnv(t, test.EnvWrapperConfig{
			WFDonConfig:     test.DonConfig{Name: "wfDon", N: nWfNodes},
			AssetDonConfig:  test.DonConfig{Name: "assetDon", N: 4},
			WriterDonConfig: test.DonConfig{Name: "writerDon", N: 4},
			NumChains:       1,
		})

		registrySel := te.Env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))[0]

		existingContracts, err := te.Env.ExistingAddresses.AddressesForChain(registrySel)
		require.NoError(t, err)
		require.Len(t, existingContracts, 4)

		// Deploy a new OCR3 contract
		resp, err := changeset.DeployOCR3(te.Env, registrySel)
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.NoError(t, te.Env.ExistingAddresses.Merge(resp.AddressBook))

		// Verify after merge there are original contracts plus one new one
		addrs, err := te.Env.ExistingAddresses.AddressesForChain(registrySel)
		require.NoError(t, err)
		require.Len(t, addrs, 5)

		wfNodes := te.GetP2PIDs("wfDon").Strings()

		w := &bytes.Buffer{}
		cfg := changeset.ConfigureOCR3Config{
			ChainSel:             te.RegistrySelector,
			NodeIDs:              wfNodes,
			OCR3Config:           &c,
			WriteGeneratedConfig: w,
		}

		_, err = changeset.ConfigureOCR3Contract(te.Env, cfg)
		require.Error(t, err)
		require.ErrorContains(t, err, "OCR contract address is unspecified")
	})

	t.Run("fails multiple OCR3 contracts but address not found", func(t *testing.T) {
		te := test.SetupContractTestEnv(t, test.EnvWrapperConfig{
			WFDonConfig:     test.DonConfig{Name: "wfDon", N: nWfNodes},
			AssetDonConfig:  test.DonConfig{Name: "assetDon", N: 4},
			WriterDonConfig: test.DonConfig{Name: "writerDon", N: 4},
			NumChains:       1,
		})

		registrySel := te.Env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))[0]

		existingContracts, err := te.Env.ExistingAddresses.AddressesForChain(registrySel)
		require.NoError(t, err)
		require.Len(t, existingContracts, 4)

		// Deploy a new OCR3 contract
		resp, err := changeset.DeployOCR3V2(te.Env, &changeset.DeployRequestV2{
			ChainSel:  registrySel,
			Qualifier: "test-ocr-contract"})
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.NoError(t, te.Env.ExistingAddresses.Merge(resp.AddressBook))
		refs := resp.DataStore.Addresses().Filter(datastore.AddressRefByQualifier("test-ocr-contract"))
		require.Len(t, refs, 1)

		// Verify after merge there are original contracts plus one new one
		addrs, err := te.Env.ExistingAddresses.AddressesForChain(registrySel)
		require.NoError(t, err)
		require.Len(t, addrs, 5)

		wfNodes := te.GetP2PIDs("wfDon").Strings()

		nfa := common.HexToAddress("0x1234567890123456789012345678901234567890")
		w := &bytes.Buffer{}
		cfg := changeset.ConfigureOCR3Config{
			ChainSel:             te.RegistrySelector,
			NodeIDs:              wfNodes,
			OCR3Config:           &c,
			Address:              &nfa,
			WriteGeneratedConfig: w,
		}

		_, err = changeset.ConfigureOCR3Contract(te.Env, cfg)
		require.Error(t, err)
		require.ErrorContains(t, err, "not found in contract set")
	})

	t.Run("mcms", func(t *testing.T) {
		te := test.SetupContractTestEnv(t, test.EnvWrapperConfig{
			WFDonConfig:     test.DonConfig{Name: "wfDon", N: nWfNodes},
			AssetDonConfig:  test.DonConfig{Name: "assetDon", N: 4},
			WriterDonConfig: test.DonConfig{Name: "writerDon", N: 4},
			NumChains:       1,
			UseMCMS:         true,
		})

		wfNodes := te.GetP2PIDs("wfDon").Strings()

		registrySel := te.Env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))[0]
		// Verify after merge there are three original contracts plus one new one
		addrs, err := te.Env.ExistingAddresses.AddressesForChain(registrySel)
		require.NoError(t, err)

		// Find new OCR3 contract
		var existingOCR3Addr string
		for addr, tv := range addrs {
			if tv.Type == internal.OCR3Capability {
				existingOCR3Addr = addr
				break
			}
		}

		w := &bytes.Buffer{}
		addr := common.HexToAddress(existingOCR3Addr)
		cfg := changeset.ConfigureOCR3Config{
			ChainSel:             te.RegistrySelector,
			NodeIDs:              wfNodes,
			OCR3Config:           &c,
			WriteGeneratedConfig: w,
			Address:              &addr,
			MCMSConfig:           &changeset.MCMSConfig{MinDuration: 0},
		}

		csOut, err := changeset.ConfigureOCR3Contract(te.Env, cfg)
		require.NoError(t, err)
		var got internal.OCR2OracleConfig
		err = json.Unmarshal(w.Bytes(), &got)
		require.NoError(t, err)
		assert.Len(t, got.Signers, 4)
		assert.Len(t, got.Transmitters, 4)
		assert.NotNil(t, csOut.MCMSTimelockProposals)
		t.Logf("got: %v", csOut.MCMSTimelockProposals[0])

		// now apply the changeset such that the proposal is signed and execed
		w2 := &bytes.Buffer{}
		cfg.WriteGeneratedConfig = w2

		err = applyProposal(t, te, commonchangeset.Configure(cldf.CreateLegacyChangeSet(changeset.ConfigureOCR3Contract), cfg))
		require.NoError(t, err)
	})
}
