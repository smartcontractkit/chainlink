package changeset_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"go.uber.org/zap/zapcore"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	kslib "github.com/smartcontractkit/chainlink/deployment/keystone"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
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
	assert.Len(t, oaddrs, 0)
}

func TestConfigureOCR3(t *testing.T) {
	t.Parallel()

	c := kslib.OracleConfig{
		MaxFaultyOracles:    1,
		DeltaProgressMillis: 12345,
	}

	t.Run("no mcms", func(t *testing.T) {

		te := SetupTestEnv(t, TestConfig{
			WFDonConfig:     DonConfig{N: 4},
			AssetDonConfig:  DonConfig{N: 4},
			WriterDonConfig: DonConfig{N: 4},
			NumChains:       1,
		})

		var wfNodes []string
		for id, _ := range te.WFNodes {
			wfNodes = append(wfNodes, id)
		}

		w := &bytes.Buffer{}
		cfg := changeset.ConfigureOCR3Config{
			ChainSel:             te.RegistrySelector,
			NodeIDs:              wfNodes,
			OCR3Config:           &c,
			WriteGeneratedConfig: w,
			UseMCMS:              false,
		}

		csOut, err := changeset.ConfigureOCR3Contract(te.Env, cfg)
		require.NoError(t, err)
		var got kslib.OCR2OracleConfig
		err = json.Unmarshal(w.Bytes(), &got)
		require.NoError(t, err)
		assert.Len(t, got.Signers, 4)
		assert.Len(t, got.Transmitters, 4)
		assert.Nil(t, csOut.Proposals)
	})

	t.Run("mcms", func(t *testing.T) {
		te := SetupTestEnv(t, TestConfig{
			WFDonConfig:     DonConfig{N: 4},
			AssetDonConfig:  DonConfig{N: 4},
			WriterDonConfig: DonConfig{N: 4},
			NumChains:       1,
			UseMCMS:         true,
		})

		var wfNodes []string
		for id, _ := range te.WFNodes {
			wfNodes = append(wfNodes, id)
		}

		w := &bytes.Buffer{}
		cfg := changeset.ConfigureOCR3Config{
			ChainSel:             te.RegistrySelector,
			NodeIDs:              wfNodes,
			OCR3Config:           &c,
			WriteGeneratedConfig: w,
			UseMCMS:              true,
		}

		csOut, err := changeset.ConfigureOCR3Contract(te.Env, cfg)
		require.NoError(t, err)
		var got kslib.OCR2OracleConfig
		err = json.Unmarshal(w.Bytes(), &got)
		require.NoError(t, err)
		assert.Len(t, got.Signers, 4)
		assert.Len(t, got.Transmitters, 4)
		assert.NotNil(t, csOut.Proposals)
		t.Logf("got: %v", csOut.Proposals[0])

		contracts := te.ContractSets()[te.RegistrySelector]
		require.NoError(t, err)
		var timelockContracts = map[uint64]*commonchangeset.TimelockExecutionContracts{
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
