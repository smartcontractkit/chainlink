package changeset

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink/deployment"
	ccipdeployment "github.com/smartcontractkit/chainlink/deployment/ccip"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestInitializePrerequisitesExisting(t *testing.T) {
	t.Parallel()
	lggr := logger.TestLogger(t)
	cfg := PrerequisiteConfig{
		ExistingContracts: map[deployment.ContractType]ContractConfig{
			ccipdeployment.LinkToken: {
				Address:        common.BigToAddress(big.NewInt(1)),
				TypeAndVersion: deployment.NewTypeAndVersion(ccipdeployment.LinkToken, deployment.Version1_0_0),
			},
			ccipdeployment.WETH9: {
				Address:        common.BigToAddress(big.NewInt(2)),
				TypeAndVersion: deployment.NewTypeAndVersion(ccipdeployment.WETH9, deployment.Version1_0_0),
			},
			ccipdeployment.TokenAdminRegistry: {
				Address:        common.BigToAddress(big.NewInt(3)),
				TypeAndVersion: deployment.NewTypeAndVersion(ccipdeployment.TokenAdminRegistry, deployment.Version1_5_0),
			},
		},
	}
	e := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		Bootstraps: 1,
		Chains:     2,
		Nodes:      4,
	})
	newChain := e.AllChainSelectors()[0]
	cfg.ChainSelectors = []uint64{newChain}
	_, err := InitializePrerequisites(e, cfg)
	require.Error(t, err)
	cfg.ExistingContracts[ccipdeployment.RegistryModule] = ContractConfig{
		Address:        common.BigToAddress(big.NewInt(4)),
		TypeAndVersion: deployment.NewTypeAndVersion(ccipdeployment.RegistryModule, deployment.Version1_5_0),
	}
	cfg.ExistingContracts[ccipdeployment.Router] = ContractConfig{
		Address:        common.BigToAddress(big.NewInt(5)),
		TypeAndVersion: deployment.NewTypeAndVersion(ccipdeployment.Router, deployment.Version1_2_0),
	}
	output, err := InitializePrerequisites(e, cfg)
	require.NoError(t, err)
	err = e.ExistingAddresses.Merge(output.AddressBook)
	require.NoError(t, err)
	state, err := ccipdeployment.LoadOnchainState(e)
	require.NoError(t, err)
	require.Equal(t, state.Chains[newChain].LinkToken.Address(), common.BigToAddress(big.NewInt(1)))
	require.Equal(t, state.Chains[newChain].Weth9.Address(), common.BigToAddress(big.NewInt(2)))
	require.Equal(t, state.Chains[newChain].TokenAdminRegistry.Address(), common.BigToAddress(big.NewInt(3)))
	require.Equal(t, state.Chains[newChain].RegistryModule.Address(), common.BigToAddress(big.NewInt(4)))
	require.Equal(t, state.Chains[newChain].Router.Address(), common.BigToAddress(big.NewInt(5)))
}

func TestInitializePrerequisitesNewDeploy(t *testing.T) {
	t.Parallel()
	lggr := logger.TestLogger(t)
	cfg := PrerequisiteConfig{
		Deploy: true,
		ExistingContracts: map[deployment.ContractType]ContractConfig{
			ccipdeployment.LinkToken: {
				Address:        common.BigToAddress(big.NewInt(1)),
				TypeAndVersion: deployment.NewTypeAndVersion(ccipdeployment.LinkToken, deployment.Version1_0_0),
			},
			ccipdeployment.WETH9: {
				Address:        common.BigToAddress(big.NewInt(2)),
				TypeAndVersion: deployment.NewTypeAndVersion(ccipdeployment.WETH9, deployment.Version1_0_0),
			},
			ccipdeployment.TokenAdminRegistry: {
				Address:        common.BigToAddress(big.NewInt(3)),
				TypeAndVersion: deployment.NewTypeAndVersion(ccipdeployment.TokenAdminRegistry, deployment.Version1_5_0),
			},
			ccipdeployment.RegistryModule: {
				Address:        common.BigToAddress(big.NewInt(4)),
				TypeAndVersion: deployment.NewTypeAndVersion(ccipdeployment.RegistryModule, deployment.Version1_5_0),
			},
			ccipdeployment.Router: {
				Address:        common.BigToAddress(big.NewInt(5)),
				TypeAndVersion: deployment.NewTypeAndVersion(ccipdeployment.Router, deployment.Version1_2_0),
			},
		},
	}
	e := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		Bootstraps: 1,
		Chains:     2,
		Nodes:      4,
	})
	newChain := e.AllChainSelectors()[0]
	cfg.ChainSelectors = []uint64{newChain}
	output, err := InitializePrerequisites(e, cfg)
	require.NoError(t, err)
	err = e.ExistingAddresses.Merge(output.AddressBook)
	require.NoError(t, err)
	state, err := ccipdeployment.LoadOnchainState(e)
	require.NoError(t, err)
	require.NotNil(t, state.Chains[newChain].LinkToken)
	require.NotNil(t, state.Chains[newChain].Weth9)
	require.NotNil(t, state.Chains[newChain].TokenAdminRegistry)
	require.NotNil(t, state.Chains[newChain].RegistryModule)
	// if new deployment we deploy router as part of DeployChainContracts due to the dependency on RMNProxy
	require.Nil(t, state.Chains[newChain].Router)
	// contracts should be newly deployed
	require.NotEqual(t, state.Chains[newChain].LinkToken.Address(), common.BigToAddress(big.NewInt(1)))
	require.NotEqual(t, state.Chains[newChain].Weth9.Address(), common.BigToAddress(big.NewInt(2)))
	require.NotEqual(t, state.Chains[newChain].TokenAdminRegistry.Address(), common.BigToAddress(big.NewInt(3)))
	require.NotEqual(t, state.Chains[newChain].RegistryModule.Address(), common.BigToAddress(big.NewInt(4)))
}
