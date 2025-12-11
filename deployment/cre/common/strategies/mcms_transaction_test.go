package strategies_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	mcmstypes "github.com/smartcontractkit/mcms/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/deployment/cre/common/strategies"
	"github.com/smartcontractkit/chainlink/deployment/cre/contracts"
	"github.com/smartcontractkit/chainlink/deployment/cre/test"
)

func getMCMSTransaction(t *testing.T, env deployment.Environment) *strategies.MCMSTransaction {
	t.Helper()

	return &strategies.MCMSTransaction{
		Env:           env,
		ChainSel:      1,
		Description:   "test",
		Config:        &contracts.MCMSConfig{},
		Address:       common.HexToAddress("0x1"),
		MCMSContracts: &commonchangeset.MCMSWithTimelockState{},
	}
}

func TestMCMSTransaction_BuildProposal(t *testing.T) {
	t.Parallel()

	fixture := test.SetupEnvV2(t, true)

	t.Run("no config", func(t *testing.T) {
		m := getMCMSTransaction(t, *fixture.Env)
		m.Config = nil

		_, err := m.BuildProposal([]mcmstypes.BatchOperation{})
		require.Error(t, err)
		assert.Equal(t, "MCMS configuration or contracts are not provided", err.Error())
	})

	t.Run("no contracts", func(t *testing.T) {
		m := getMCMSTransaction(t, *fixture.Env)
		m.MCMSContracts = nil

		_, err := m.BuildProposal([]mcmstypes.BatchOperation{})
		require.Error(t, err)
		assert.Equal(t, "MCMS configuration or contracts are not provided", err.Error())
	})

	t.Run("no timelock", func(t *testing.T) {
		m := getMCMSTransaction(t, *fixture.Env)
		cfg := contracts.MCMSConfig{
			MinDelay: 0,
			TimelockQualifierPerChain: map[uint64]string{
				fixture.RegistrySelector: "",
			},
		}
		mcmsContracts, err := strategies.GetMCMSContracts(*fixture.Env, fixture.RegistrySelector, cfg)
		require.NoError(t, err)
		m.MCMSContracts = mcmsContracts
		m.MCMSContracts.Timelock = nil

		_, err = m.BuildProposal([]mcmstypes.BatchOperation{})
		require.Error(t, err)
		assert.Equal(t, "MCMS contracts are not properly initialized, missing Timelock or Proposer", err.Error())
	})

	t.Run("no proposer", func(t *testing.T) {
		m := getMCMSTransaction(t, *fixture.Env)
		cfg := contracts.MCMSConfig{
			MinDelay: 0,
			TimelockQualifierPerChain: map[uint64]string{
				fixture.RegistrySelector: "",
			},
		}
		mcmsContracts, err := strategies.GetMCMSContracts(*fixture.Env, fixture.RegistrySelector, cfg)
		require.NoError(t, err)
		m.MCMSContracts = mcmsContracts
		m.MCMSContracts.ProposerMcm = nil

		_, err = m.BuildProposal([]mcmstypes.BatchOperation{})
		require.Error(t, err)
		assert.Equal(t, "MCMS contracts are not properly initialized, missing Timelock or Proposer", err.Error())
	})

	t.Run("no operations", func(t *testing.T) {
		m := getMCMSTransaction(t, *fixture.Env)
		cfg := contracts.MCMSConfig{
			MinDelay: 0,
			TimelockQualifierPerChain: map[uint64]string{
				fixture.RegistrySelector: "",
			},
		}
		mcmsContracts, err := strategies.GetMCMSContracts(*fixture.Env, fixture.RegistrySelector, cfg)
		require.NoError(t, err)
		m.MCMSContracts = mcmsContracts

		_, err = m.BuildProposal([]mcmstypes.BatchOperation{})
		require.Error(t, err)
		assert.Equal(t, "no operations provided to build proposal", err.Error())
	})

	t.Run("MCMBasedOnAction fails", func(t *testing.T) {
		m := getMCMSTransaction(t, *fixture.Env)
		cfg := contracts.MCMSConfig{
			MinDelay: 0,
			TimelockQualifierPerChain: map[uint64]string{
				fixture.RegistrySelector: "",
			},
		}
		mcmsContracts, err := strategies.GetMCMSContracts(*fixture.Env, fixture.RegistrySelector, cfg)
		require.NoError(t, err)
		m.MCMSContracts = mcmsContracts
		m.Config.MCMSAction = "invalid"

		_, err = m.BuildProposal([]mcmstypes.BatchOperation{{}})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get mcms contract by action 'invalid'")
	})

	t.Run("uses Proposer when not specifying an action (defaults to `schedule`)", func(t *testing.T) {
		m := getMCMSTransaction(t, *fixture.Env)
		cfg := contracts.MCMSConfig{
			MinDelay: 0,
			TimelockQualifierPerChain: map[uint64]string{
				fixture.RegistrySelector: "",
			},
		}
		mcmsContracts, err := strategies.GetMCMSContracts(*fixture.Env, fixture.RegistrySelector, cfg)
		require.NoError(t, err)
		m.Config = &cfg
		m.MCMSContracts = mcmsContracts
		m.ChainSel = fixture.RegistrySelector

		op, err := proposalutils.BatchOperationForChain(m.ChainSel, m.Address.Hex(), []byte{0x01, 0x02, 0x03}, big.NewInt(0), "", nil)
		require.NoError(t, err)

		p, err := m.BuildProposal([]mcmstypes.BatchOperation{op})
		require.NoError(t, err)

		metadata, ok := p.ChainMetadata[mcmstypes.ChainSelector(fixture.RegistrySelector)]
		assert.True(t, ok)
		assert.Equal(t, mcmsContracts.ProposerMcm.Address().String(), metadata.MCMAddress)
	})

	t.Run("uses Bypasser when specified", func(t *testing.T) {
		m := getMCMSTransaction(t, *fixture.Env)
		cfg := contracts.MCMSConfig{
			MinDelay:   0,
			MCMSAction: mcmstypes.TimelockActionBypass,
			TimelockQualifierPerChain: map[uint64]string{
				fixture.RegistrySelector: "",
			},
		}
		mcmsContracts, err := strategies.GetMCMSContracts(*fixture.Env, fixture.RegistrySelector, cfg)
		require.NoError(t, err)
		m.Config = &cfg
		m.MCMSContracts = mcmsContracts
		m.ChainSel = fixture.RegistrySelector

		op, err := proposalutils.BatchOperationForChain(m.ChainSel, m.Address.Hex(), []byte{0x01, 0x02, 0x03}, big.NewInt(0), "", nil)
		require.NoError(t, err)

		p, err := m.BuildProposal([]mcmstypes.BatchOperation{op})
		require.NoError(t, err)

		metadata, ok := p.ChainMetadata[mcmstypes.ChainSelector(fixture.RegistrySelector)]
		assert.True(t, ok)
		assert.Equal(t, mcmsContracts.BypasserMcm.Address().String(), metadata.MCMAddress)
	})

	t.Run("uses Canceller when specified", func(t *testing.T) {
		m := getMCMSTransaction(t, *fixture.Env)
		cfg := contracts.MCMSConfig{
			MinDelay:   0,
			MCMSAction: mcmstypes.TimelockActionCancel,
			TimelockQualifierPerChain: map[uint64]string{
				fixture.RegistrySelector: "",
			},
		}
		mcmsContracts, err := strategies.GetMCMSContracts(*fixture.Env, fixture.RegistrySelector, cfg)
		require.NoError(t, err)
		m.Config = &cfg
		m.MCMSContracts = mcmsContracts
		m.ChainSel = fixture.RegistrySelector

		op, err := proposalutils.BatchOperationForChain(m.ChainSel, m.Address.Hex(), []byte{0x01, 0x02, 0x03}, big.NewInt(0), "", nil)
		require.NoError(t, err)

		p, err := m.BuildProposal([]mcmstypes.BatchOperation{op})
		require.NoError(t, err)

		metadata, ok := p.ChainMetadata[mcmstypes.ChainSelector(fixture.RegistrySelector)]
		assert.True(t, ok)
		assert.Equal(t, mcmsContracts.CancellerMcm.Address().String(), metadata.MCMAddress)
	})

	t.Run("uses custom ValidDuration value to set the proposal duration", func(t *testing.T) {
		m := getMCMSTransaction(t, *fixture.Env)
		validDuration, err := commonconfig.NewDuration(2 * time.Second)
		require.NoError(t, err)
		cfg := contracts.MCMSConfig{
			MinDelay: 0,
			TimelockQualifierPerChain: map[uint64]string{
				fixture.RegistrySelector: "",
			},
			ValidDuration: &validDuration,
		}
		mcmsContracts, err := strategies.GetMCMSContracts(*fixture.Env, fixture.RegistrySelector, cfg)
		require.NoError(t, err)
		m.Config = &cfg
		m.MCMSContracts = mcmsContracts
		m.ChainSel = fixture.RegistrySelector

		op, err := proposalutils.BatchOperationForChain(m.ChainSel, m.Address.Hex(), []byte{0x01, 0x02, 0x03}, big.NewInt(0), "", nil)
		require.NoError(t, err)

		p, err := m.BuildProposal([]mcmstypes.BatchOperation{op})
		require.NoError(t, err)

		expectedValidUntil := time.Now().Add(validDuration.Duration()).Unix()
		// Using InDelta to allow for slight timing differences during test execution
		assert.InDelta(t, uint32(expectedValidUntil), p.ValidUntil, 1, "ValidUntil should be within 1 second of expected value") //nolint:gosec // G115
	})
}
