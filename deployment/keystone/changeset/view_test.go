package changeset_test

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	chain_selectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/internal"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/test"
)

var oracleConfig = changeset.OracleConfig{
	DeltaProgressMillis:               30000,
	DeltaResendMillis:                 5000,
	DeltaInitialMillis:                5000,
	DeltaRoundMillis:                  2000,
	DeltaGraceMillis:                  500,
	DeltaCertifiedCommitRequestMillis: 1000,
	DeltaStageMillis:                  30000,
	MaxRoundsPerEpoch:                 10,
	TransmissionSchedule:              []int{},
	MaxDurationQueryMillis:            1000,
	MaxDurationObservationMillis:      1000,
	MaxDurationShouldAcceptMillis:     1000,
	MaxDurationShouldTransmitMillis:   1000,
	MaxFaultyOracles:                  1,
	MaxQueryLengthBytes:               1000000,
	MaxObservationLengthBytes:         1000000,
	MaxReportLengthBytes:              1000000,
	MaxOutcomeLengthBytes:             1000000,
	MaxReportCount:                    20,
	MaxBatchSize:                      20,
	OutcomePruningThreshold:           3600,
	UniqueReports:                     true,
	RequestTimeout:                    30 * time.Second,
}

func TestKeystoneView(t *testing.T) {
	t.Parallel()
	env := test.SetupContractTestEnv(t, test.EnvWrapperConfig{
		WFDonConfig:     test.DonConfig{N: 4, Name: "wfDon"},
		AssetDonConfig:  test.DonConfig{N: 4, Name: "assetDon"},
		WriterDonConfig: test.DonConfig{N: 4, Name: "writerDon"},
		NumChains:       1,
	})
	registryChain := env.Env.AllChainSelectors()[0]
	oracleConfig.TransmissionSchedule = []int{len(env.Env.NodeIDs)}

	addrs, err := env.Env.ExistingAddresses.AddressesForChain(registryChain)
	require.NoError(t, err)

	var newOCR3Addr, newForwarderAddr string
	for addr, tv := range addrs {
		if newForwarderAddr != "" && newOCR3Addr != "" {
			break
		}
		switch tv.Type {
		case internal.KeystoneForwarder:
			newForwarderAddr = addr
			continue
		case internal.OCR3Capability:
			newOCR3Addr = addr
		default:
			continue
		}
	}

	t.Run("successfully generates a view of the keystone state", func(t *testing.T) {
		oracleConfigCopy := oracleConfig

		w := &bytes.Buffer{}
		na := common.HexToAddress(newOCR3Addr)
		cfg := changeset.ConfigureOCR3Config{
			ChainSel:             env.RegistrySelector,
			NodeIDs:              env.Env.NodeIDs,
			Address:              &na,
			OCR3Config:           &oracleConfigCopy,
			WriteGeneratedConfig: w,
		}
		_, err = changeset.ConfigureOCR3Contract(env.Env, cfg)
		require.NoError(t, err)

		var prevView json.RawMessage = []byte("{}")
		a, err := changeset.ViewKeystone(env.Env, prevView)
		require.NoError(t, err)
		b, err := a.MarshalJSON()
		require.NoError(t, err)
		require.NotEmpty(t, b)

		var outView changeset.KeystoneView
		require.NoError(t, json.Unmarshal(b, &outView))

		chainID, err := chain_selectors.ChainIdFromSelector(registryChain)
		require.NoError(t, err)
		chainName, err := chain_selectors.NameFromChainId(chainID)
		require.NoError(t, err)

		viewChain, ok := outView.Chains[chainName]
		require.True(t, ok)
		viewOCR3Config, ok := viewChain.OCRContracts[newOCR3Addr]
		require.True(t, ok)
		require.Len(t, viewChain.OCRContracts, 1)
		require.Equal(t, oracleConfig, viewOCR3Config.OffchainConfig)
		viewForwarders, ok := viewChain.Forwarders[newForwarderAddr]
		require.True(t, ok)
		require.Len(t, viewForwarders, 1)
		require.Equal(t, uint32(1), viewForwarders[0].DonID)
		require.Equal(t, uint8(1), viewForwarders[0].F)
		require.Equal(t, uint32(1), viewForwarders[0].ConfigVersion)
		require.Len(t, viewForwarders[0].Signers, 4)
	})

	t.Run("generates a partial view of the keystone state with OCR3 not configured", func(t *testing.T) {
		// Deploy a new OCR3 contract
		resp, err := changeset.DeployOCR3(env.Env, registryChain)
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.NoError(t, env.Env.ExistingAddresses.Merge(resp.AddressBook))

		var prevView json.RawMessage = []byte("{}")
		a, err := changeset.ViewKeystone(env.Env, prevView)
		require.NoError(t, err)
		b, err := a.MarshalJSON()
		require.NoError(t, err)
		require.NotEmpty(t, b)

		var outView changeset.KeystoneView
		require.NoError(t, json.Unmarshal(b, &outView))
		chainID, err := chain_selectors.ChainIdFromSelector(registryChain)
		require.NoError(t, err)
		chainName, err := chain_selectors.NameFromChainId(chainID)
		require.NoError(t, err)

		view, ok := outView.Chains[chainName]
		require.True(t, ok)
		assert.NotNil(t, view.Forwarders)
		assert.NotNil(t, view.OCRContracts)
		require.Len(t, view.OCRContracts, 1) // There already are OCR views available at this point
		assert.NotNil(t, view.WorkflowRegistry)
		assert.NotNil(t, view.CapabilityRegistry)
	})

	t.Run("fails to generate a view of the keystone state with a bad OracleConfig", func(t *testing.T) {
		oracleConfigCopy := oracleConfig
		oracleConfigCopy.DeltaRoundMillis = 0
		oracleConfigCopy.DeltaProgressMillis = 0

		w := &bytes.Buffer{}
		na := common.HexToAddress(newOCR3Addr)
		cfg := changeset.ConfigureOCR3Config{
			ChainSel:             env.RegistrySelector,
			NodeIDs:              env.Env.NodeIDs,
			Address:              &na,
			OCR3Config:           &oracleConfigCopy,
			WriteGeneratedConfig: w,
		}
		_, err = changeset.ConfigureOCR3Contract(env.Env, cfg)
		require.NoError(t, err)
		var prevView json.RawMessage = []byte("{}")
		_, err = changeset.ViewKeystone(env.Env, prevView)
		require.ErrorContains(t, err, "failed to view chain")
		require.ErrorContains(t, err, "DeltaRound (0s) must be less than DeltaProgress (0s)")
	})
}
