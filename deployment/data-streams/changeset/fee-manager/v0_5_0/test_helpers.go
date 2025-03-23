package fee_manager_v0_5_0

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/deployment"
	commonChangesets "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	reward_manager "github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/reward-manager"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/testutil"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/types"
	verifier "github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/verifier/v0_5_0"
	"github.com/stretchr/testify/require"
)

type DataStreamsTestEnvSetupOutput struct {
	Env                  deployment.Environment
	LinkTokenAddress     common.Address
	linkState            *commonChangesets.LinkTokenState
	FeeManagerAddress    common.Address
	RewardManagerAddress common.Address
	VerifierAddress      common.Address
	VerifierProxyAddress common.Address
}

type DataStreamsTestEnvOptions struct {
	DeployRewardManager bool
	DeployFeeManager    bool
	DeployLinkToken     bool
	DeployMCMS          bool
	DeployVerifier      bool
	DeployVerifierProxy bool
	initialEnv          deployment.Environment
}

func NewDefaultOptions() DataStreamsTestEnvOptions {
	return DataStreamsTestEnvOptions{
		DeployLinkToken:     true,
		DeployMCMS:          true,
		DeployRewardManager: false,
		DeployFeeManager:    true,
		DeployVerifier:      false,
		DeployVerifierProxy: false,
	}
}

func NewDataStreamsEnvironment(t *testing.T, opts DataStreamsTestEnvOptions) (DataStreamsTestEnvSetupOutput, error) {
	t.Helper()

	e := testutil.NewMemoryEnv(t, opts.DeployMCMS, 0)

	feeManagerAddress := common.HexToAddress("0x044304C47eD3B1C1357569960A537056AFE8c815")
	rewardManagerAddress := common.HexToAddress("0x0fd8b81e3d1143ec7f1ce474827ab93c43523ea2")
	verifierProxyAddr := common.Address{}
	verifierAddr := common.Address{}

	linkTokenAddress := common.Address{}
	if opts.DeployLinkToken {
		env, err := commonChangesets.Apply(t, e, nil,
			commonChangesets.Configure(
				deployment.CreateLegacyChangeSet(commonChangesets.DeployLinkToken),
				[]uint64{testutil.TestChain.Selector},
			),
		)
		require.NoError(t, err)

		addresses, err := env.ExistingAddresses.AddressesForChain(testutil.TestChain.Selector)
		require.NoError(t, err)

		chain := env.Chains[testutil.TestChain.Selector]
		linkState, err := commonChangesets.MaybeLoadLinkTokenChainState(chain, addresses)
		require.NoError(t, err)
		require.NotNil(t, linkState.LinkToken)
		linkTokenAddress = linkState.LinkToken.Address()
		e = env
	}

	if opts.DeployRewardManager {
		env, err := commonChangesets.Apply(t, e, nil,
			commonChangesets.Configure(
				reward_manager.DeployRewardManagerChangeset,
				reward_manager.DeployRewardManagerConfig{
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
			ProxyAddress:         common.HexToAddress("0x742d35Cc6634C0532925a3b844Bc454e4438f44e"),
			RewardManagerAddress: rewardManagerAddress,
		}

		env, err := commonChangesets.Apply(t, e, nil,
			commonChangesets.Configure(
				DeployFeeManagerChangeset,
				cc,
			),
		)
		require.NoError(t, err)

		fmAddressHex, err := deployment.SearchAddressBook(e.ExistingAddresses, testutil.TestChain.Selector, types.FeeManager)
		require.NoError(t, err)
		feeManagerAddress = common.HexToAddress(fmAddressHex)
		e = env
	}

	if opts.DeployVerifierProxy && opts.DeployVerifier {
		e, verifierProxyAddr, verifierAddr = verifier.DeployVerifierProxyAndVerifier(t, e)
	}

	return DataStreamsTestEnvSetupOutput{
		Env:                  e,
		LinkTokenAddress:     linkTokenAddress,
		FeeManagerAddress:    feeManagerAddress,
		RewardManagerAddress: rewardManagerAddress,
		VerifierAddress:      verifierAddr,
		VerifierProxyAddress: verifierProxyAddr,
	}, nil
}
