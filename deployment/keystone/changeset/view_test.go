package changeset_test

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/internal"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/test"
)

func TestKeystoneView(t *testing.T) {
	t.Parallel()
	env := test.SetupTestEnv(t, test.TestConfig{
		WFDonConfig:     test.DonConfig{N: 4},
		AssetDonConfig:  test.DonConfig{N: 4},
		WriterDonConfig: test.DonConfig{N: 4},
		NumChains:       1,
	})
	registryChain := env.Env.AllChainSelectors()[0]

	resp, err := changeset.DeployCapabilityRegistry(env.Env, registryChain)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NoError(t, env.Env.ExistingAddresses.Merge(resp.AddressBook))

	addrs, err := env.Env.ExistingAddresses.AddressesForChain(registryChain)
	require.NoError(t, err)

	var newOCR3Addr string
	for addr, tv := range addrs {
		if tv.Type == internal.OCR3Capability {
			newOCR3Addr = addr
			break
		}
	}

	var wfNodes []string
	for id := range env.WFNodes {
		wfNodes = append(wfNodes, id)
	}

	oracleConfig := changeset.OracleConfig{
		DeltaProgressMillis:               30000,
		DeltaResendMillis:                 5000,
		DeltaInitialMillis:                5000,
		DeltaRoundMillis:                  2000,
		DeltaGraceMillis:                  500,
		DeltaCertifiedCommitRequestMillis: 1000,
		DeltaStageMillis:                  30000,
		MaxRoundsPerEpoch:                 10,
		TransmissionSchedule:              []int{len(wfNodes)},
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
		MaxBatchSize:                      1000,
		OutcomePruningThreshold:           3600,
		UniqueReports:                     true,
		RequestTimeout:                    30 * time.Second,
	}

	w := &bytes.Buffer{}
	na := common.HexToAddress(newOCR3Addr)
	cfg := changeset.ConfigureOCR3Config{
		ChainSel:             env.RegistrySelector,
		NodeIDs:              wfNodes,
		Address:              &na,
		OCR3Config:           &oracleConfig,
		WriteGeneratedConfig: w,
	}
	_, err = changeset.ConfigureOCR3Contract(env.Env, cfg)
	require.NoError(t, err)

	resp, err = changeset.DeployForwarder(env.Env, changeset.DeployForwarderRequest{})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NoError(t, env.Env.ExistingAddresses.Merge(resp.AddressBook))

	a, err := changeset.ViewKeystone(env.Env)
	require.NoError(t, err)
	b, err := a.MarshalJSON()
	require.NoError(t, err)
	require.NotEmpty(t, b)
	t.Log(string(b))

	var outView changeset.KeystoneView
	require.NoError(t, json.Unmarshal(b, &outView))

	chainID, err := chain_selectors.ChainIdFromSelector(registryChain)
	require.NoError(t, err)
	chainName, err := chain_selectors.NameFromChainId(chainID)
	require.NoError(t, err)

	viewChain, ok := outView.Chains[chainName]
	require.True(t, ok)
	viewOCR3Config, ok := viewChain.OCR3ConfigView[newOCR3Addr]
	require.True(t, ok)
	require.Equal(t, oracleConfig, viewOCR3Config.OffchainConfig)
}
