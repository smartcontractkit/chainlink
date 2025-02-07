package changeset_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
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

	registrySel := env.AllChainSelectors()[0]

	resp, err := changeset.DeployOCR3(env, registrySel)
	require.NoError(t, err)
	require.NotNil(t, resp)
	// OCR3 should be deployed on chain 0
	addrs, err := resp.AddressBook.AddressesForChain(registrySel)
	require.NoError(t, err)
	require.Len(t, addrs, 1)

	// nothing on chain 1
	require.NotEqual(t, registrySel, env.AllChainSelectors()[1])
	oaddrs, _ := resp.AddressBook.AddressesForChain(env.AllChainSelectors()[1])
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
		te := test.SetupTestEnv(t, test.TestConfig{
			WFDonConfig:     test.DonConfig{N: nWfNodes},
			AssetDonConfig:  test.DonConfig{N: 4},
			WriterDonConfig: test.DonConfig{N: 4},
			NumChains:       1,
		})

		var wfNodes []string
		for id := range te.WFNodes {
			wfNodes = append(wfNodes, id)
		}

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
		assert.Nil(t, csOut.Proposals)
	})

	t.Run("success multiple OCR3 contracts", func(t *testing.T) {
		te := test.SetupTestEnv(t, test.TestConfig{
			WFDonConfig:     test.DonConfig{N: nWfNodes},
			AssetDonConfig:  test.DonConfig{N: 4},
			WriterDonConfig: test.DonConfig{N: 4},
			NumChains:       1,
		})

		registrySel := te.Env.AllChainSelectors()[0]

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

		var wfNodes []string
		for id := range te.WFNodes {
			wfNodes = append(wfNodes, id)
		}

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
		assert.Nil(t, csOut.Proposals)
	})

	t.Run("fails multiple OCR3 contracts but unspecified address", func(t *testing.T) {
		te := test.SetupTestEnv(t, test.TestConfig{
			WFDonConfig:     test.DonConfig{N: nWfNodes},
			AssetDonConfig:  test.DonConfig{N: 4},
			WriterDonConfig: test.DonConfig{N: 4},
			NumChains:       1,
		})

		registrySel := te.Env.AllChainSelectors()[0]

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

		var wfNodes []string
		for id := range te.WFNodes {
			wfNodes = append(wfNodes, id)
		}

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
		te := test.SetupTestEnv(t, test.TestConfig{
			WFDonConfig:     test.DonConfig{N: nWfNodes},
			AssetDonConfig:  test.DonConfig{N: 4},
			WriterDonConfig: test.DonConfig{N: 4},
			NumChains:       1,
		})

		registrySel := te.Env.AllChainSelectors()[0]

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

		var wfNodes []string
		for id := range te.WFNodes {
			wfNodes = append(wfNodes, id)
		}

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
		te := test.SetupTestEnv(t, test.TestConfig{
			WFDonConfig:     test.DonConfig{N: nWfNodes},
			AssetDonConfig:  test.DonConfig{N: 4},
			WriterDonConfig: test.DonConfig{N: 4},
			NumChains:       1,
			UseMCMS:         true,
		})

		var wfNodes []string
		for id := range te.WFNodes {
			wfNodes = append(wfNodes, id)
		}

		w := &bytes.Buffer{}
		cfg := changeset.ConfigureOCR3Config{
			ChainSel:             te.RegistrySelector,
			NodeIDs:              wfNodes,
			OCR3Config:           &c,
			WriteGeneratedConfig: w,
			MCMSConfig:           &changeset.MCMSConfig{MinDuration: 0},
		}

		csOut, err := changeset.ConfigureOCR3Contract(te.Env, cfg)
		require.NoError(t, err)
		var got internal.OCR2OracleConfig
		err = json.Unmarshal(w.Bytes(), &got)
		require.NoError(t, err)
		assert.Len(t, got.Signers, 4)
		assert.Len(t, got.Transmitters, 4)
		assert.NotNil(t, csOut.Proposals)
		t.Logf("got: %v", csOut.Proposals[0])

		contracts := te.ContractSets()[te.RegistrySelector]
		require.NoError(t, err)
		var timelockContracts = map[uint64]*proposalutils.TimelockExecutionContracts{
			te.RegistrySelector: {
				Timelock:  contracts.Timelock,
				CallProxy: contracts.CallProxy,
			},
		}

		// now apply the changeset such that the proposal is signed and execed
		w2 := &bytes.Buffer{}
		cfg.WriteGeneratedConfig = w2
		_, err = commonchangeset.ApplyChangesets(t, te.Env, timelockContracts, []commonchangeset.ChangesetApplication{
			{
				Changeset: commonchangeset.WrapChangeSet(changeset.ConfigureOCR3Contract),
				Config:    cfg,
			},
		})
		require.NoError(t, err)
	})
}
