package fee_manager

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/deployment"
	commonchangesets "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	commonstate "github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	rewardmanager "github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/reward-manager"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/testutil"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/types"
)

type DataStreamsTestEnvSetupOutput struct {
	Env                  deployment.Environment
	LinkTokenAddress     common.Address
	linkState            *commonchangesets.LinkTokenState
	FeeManagerAddress    common.Address
	RewardManagerAddress common.Address
}

type DataStreamsTestEnvOptions struct {
	DeployRewardManager bool
	DeployFeeManager    bool
	DeployLinkToken     bool
	DeployMCMS          bool
	initialEnv          deployment.Environment
}

func NewDefaultOptions() DataStreamsTestEnvOptions {
	return DataStreamsTestEnvOptions{
		DeployLinkToken:     true,
		DeployMCMS:          true,
		DeployRewardManager: false,
		DeployFeeManager:    true,
	}
}

func NewDataStreamsEnvironment(t *testing.T, opts DataStreamsTestEnvOptions) (DataStreamsTestEnvSetupOutput, error) {
	t.Helper()

	e := testutil.NewMemoryEnv(t, opts.DeployMCMS, 0)

	feeManagerAddress := common.HexToAddress("0x044304C47eD3B1C1357569960A537056AFE8c815")
	rewardManagerAddress := common.HexToAddress("0x0fd8b81e3d1143ec7f1ce474827ab93c43523ea2")

	linkTokenAddress := common.Address{}
	if opts.DeployLinkToken {
		env, err := commonchangesets.Apply(t, e, nil,
			commonchangesets.Configure(
				deployment.CreateLegacyChangeSet(commonchangesets.DeployLinkToken),
				[]uint64{testutil.TestChain.Selector},
			),
		)
		require.NoError(t, err)

		addresses, err := env.ExistingAddresses.AddressesForChain(testutil.TestChain.Selector)
		require.NoError(t, err)

		chain := env.Chains[testutil.TestChain.Selector]
		linkState, err := commonstate.MaybeLoadLinkTokenChainState(chain, addresses)
		require.NoError(t, err)
		require.NotNil(t, linkState.LinkToken)
		linkTokenAddress = linkState.LinkToken.Address()
		e = env
	}

	if opts.DeployRewardManager {
		env, err := commonchangesets.Apply(t, e, nil,
			commonchangesets.Configure(
				rewardmanager.DeployRewardManagerChangeset,
				rewardmanager.DeployRewardManagerConfig{
					ChainsToDeploy:   []uint64{testutil.TestChain.Selector},
					LinkTokenAddress: linkTokenAddress,
				},
			),
		)
		require.NoError(t, err)
		e = env

		// Ensure the Contract was deployed
		rewardManagerAddressHex, err := deployment.SearchAddressBook(e.ExistingAddresses, testutil.TestChain.Selector, types.RewardManager)
		require.NoError(t, err)
		rewardManagerAddress = common.HexToAddress(rewardManagerAddressHex)
	}

	if opts.DeployFeeManager {
		// FM checks LinkToken is ERC20 - but accepts any address for NativeTokenAddress
		cc := DeployFeeManagerConfig{
			ChainsToDeploy:       []uint64{testutil.TestChain.Selector},
			LinkTokenAddress:     linkTokenAddress,
			NativeTokenAddress:   common.HexToAddress("0x3e5e9111ae8eb78fe1cc3bb8915d5d461f3ef9a9"),
			VerifierProxyAddress: common.HexToAddress("0x742d35Cc6634C0532925a3b844Bc454e4438f44e"),
			RewardManagerAddress: rewardManagerAddress,
		}

		env, err := commonchangesets.Apply(t, e, nil,
			commonchangesets.Configure(
				DeployFeeManagerChangeset,
				cc,
			),
		)
		require.NoError(t, err)

		fmAddressHex, err := deployment.SearchAddressBook(env.ExistingAddresses, testutil.TestChain.Selector, types.FeeManager)
		require.NoError(t, err)
		feeManagerAddress = common.HexToAddress(fmAddressHex)
		e = env
	}

	return DataStreamsTestEnvSetupOutput{
		Env:                  e,
		LinkTokenAddress:     linkTokenAddress,
		FeeManagerAddress:    feeManagerAddress,
		RewardManagerAddress: rewardManagerAddress,
	}, nil
}
