package fee_manager

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment"
	commonchangesets "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	commonstate "github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/testutil"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/types"
)

type DataStreamsTestEnvSetupOutput struct {
	Env               deployment.Environment
	LinkTokenAddress  common.Address
	FeeManagerAddress common.Address
}

type DataStreamsTestEnvOptions struct {
	DeployFeeManager bool
	DeployLinkToken  bool
	DeployMCMS       bool
}

func NewDefaultOptions() DataStreamsTestEnvOptions {
	return DataStreamsTestEnvOptions{
		DeployLinkToken:  true,
		DeployMCMS:       true,
		DeployFeeManager: true,
	}
}

func NewDataStreamsEnvironment(t *testing.T, opts DataStreamsTestEnvOptions) (DataStreamsTestEnvSetupOutput, error) {
	t.Helper()

	e := testutil.NewMemoryEnv(t, opts.DeployMCMS, 0)

	feeManagerAddress := common.HexToAddress("0x044304C47eD3B1C1357569960A537056AFE8c815")

	linkTokenAddress := common.Address{}
	if opts.DeployLinkToken {
		env, err := commonchangesets.Apply(t, e, nil,
			commonchangesets.Configure(
				cldf.CreateLegacyChangeSet(commonchangesets.DeployLinkToken),
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

	if opts.DeployFeeManager {
		// FM checks LinkToken is ERC20 - but accepts any address for NativeTokenAddress
		cc := DeployFeeManagerConfig{
			ChainsToDeploy: map[uint64]DeployFeeManager{testutil.TestChain.Selector: {
				LinkTokenAddress:     linkTokenAddress,
				NativeTokenAddress:   common.HexToAddress("0x3e5e9111ae8eb78fe1cc3bb8915d5d461f3ef9a9"),
				VerifierProxyAddress: common.HexToAddress("0x742d35Cc6634C0532925a3b844Bc454e4438f44e"),
				RewardManagerAddress: common.HexToAddress("0x0fd8b81e3d1143ec7f1ce474827ab93c43523ea2"),
			}},
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
		Env:               e,
		LinkTokenAddress:  linkTokenAddress,
		FeeManagerAddress: feeManagerAddress,
	}, nil
}
