package fee_manager_v0_5_0

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/deployment"
	commonChangesets "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/testutil"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/types"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/llo-feeds/generated/mock_fee_manager_v0_5_0"
	"github.com/stretchr/testify/require"
)

func TestPayLinkDeficit(t *testing.T) {
	testOptions := DataStreamsTestEnvOptions{
		DeployRewardManager: true,
		DeployLinkToken:     true,
		DeployMCMS:          true,
		DeployVerifier:      true,
		DeployVerifierProxy: true,
		DeployFeeManager:    false, // we override this below with a MockFeeManager
	}

	res, err := NewDataStreamsEnvironment(t, testOptions)
	require.NoError(t, err)

	e := res.Env
	ab := e.ExistingAddresses
	chain := e.Chains[testutil.TestChain.Selector]

	cc := DeployFeeManagerConfig{
		LinkTokenAddress:     res.LinkTokenAddress,
		NativeTokenAddress:   common.HexToAddress("0x55E8606BbE91513725f9e4d405B7956D9F9e236F"),
		ProxyAddress:         res.VerifierProxyAddress,
		RewardManagerAddress: res.RewardManagerAddress,
	}

	// This uses a MockFeeManager. The subject under test here is the PayLinkDeficit changeset.
	// This is modeled as a client/server test where the "client" is the PayLinkDeficit changeset
	// and the "server" is the MockFeeManager. The PayLinkDeficit changeset will call the MockFeeManager using
	// the real FeeManager interface. The MockFeeManager will then validate the call and return a response.
	_, err = changeset.DeployContract[*mock_fee_manager_v0_5_0.MockFeeManager](e, ab, chain, MockFeeManagerDeployFn(cc))
	require.NoError(t, err)

	feeManagerAddressHex, err := deployment.SearchAddressBook(e.ExistingAddresses, chain.Selector, types.FeeManager)
	require.NoError(t, err)
	feeManagerAddress := common.HexToAddress(feeManagerAddressHex)

	configDigest := [32]byte{1}

	e, err = commonChangesets.Apply(t, e, nil,
		commonChangesets.Configure(
			PayLinkDeficitChangeset,
			PayLinkDeficitConfig{
				ConfigPerChain: map[uint64][]PayLinkDeficit{
					testutil.TestChain.Selector: {
						{
							FeeManagerAddress: feeManagerAddress,
							ConfigDigest:      configDigest,
						},
					},
				},
			},
		))

	require.NoError(t, err)

	feeManager, err := LoadFeeManagerState(e, testutil.TestChain.Selector, feeManagerAddress.String())
	require.NoError(t, err)
	require.NotNil(t, feeManager)

	logIterator, err := feeManager.FilterLinkDeficitCleared(nil, [][32]byte{configDigest})
	require.NoError(t, err)

	foundExpected := false
	for logIterator.Next() {
		if logIterator.Event.ConfigDigest == configDigest {
			foundExpected = true
			break
		}
	}
	require.True(t, foundExpected)
	err = logIterator.Close()
	if err != nil {
		t.Errorf("Error closing log iterator: %v", err)
	}

}

func MockFeeManagerDeployFn(cfg DeployFeeManagerConfig) changeset.ContractDeployFn[*mock_fee_manager_v0_5_0.MockFeeManager] {
	return func(chain deployment.Chain) *changeset.ContractDeployment[*mock_fee_manager_v0_5_0.MockFeeManager] {
		ccsAddr, ccsTx, ccs, err := mock_fee_manager_v0_5_0.DeployMockFeeManager(
			chain.DeployerKey,
			chain.Client,
			cfg.LinkTokenAddress,
			cfg.NativeTokenAddress,
			cfg.ProxyAddress,
			cfg.RewardManagerAddress,
		)
		if err != nil {
			return &changeset.ContractDeployment[*mock_fee_manager_v0_5_0.MockFeeManager]{
				Err: err,
			}
		}
		return &changeset.ContractDeployment[*mock_fee_manager_v0_5_0.MockFeeManager]{
			Address:  ccsAddr,
			Contract: ccs,
			Tx:       ccsTx,
			Tv:       deployment.NewTypeAndVersion(types.FeeManager, deployment.Version0_5_0),
			Err:      nil,
		}
	}
}
