package deploymenttest

import (
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/fee_quoter"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/nonce_manager"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/offramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/onramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/token_admin_registry"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
)

// This tests ensures that all of the compiled contracts can be
// deployed to an actual blockchain (i.e no "max code size exceeded" errors).
// It does not attempt to correctly set up the contracts, so bogus inputs are used.
func TestDeployAllV1_6(t *testing.T) {
	owner := testutils.MustNewSimTransactor(t)
	chain := backends.NewSimulatedBackend(core.GenesisAlloc{
		owner.From: {Balance: assets.Ether(100).ToInt()},
	}, 30e6)

	// router
	_, _, _, err := router.DeployRouter(owner, chain, common.HexToAddress("0x1"), common.HexToAddress("0x2"))
	require.NoError(t, err)
	chain.Commit()

	// nonce manager
	_, _, _, err = nonce_manager.DeployNonceManager(owner, chain, []common.Address{common.HexToAddress("0x1")})
	require.NoError(t, err)
	chain.Commit()

	// offramp
	_, _, _, err = offramp.DeployOffRamp(owner, chain, offramp.OffRampStaticConfig{
		ChainSelector:      1,
		RmnRemote:          common.HexToAddress("0x1"),
		TokenAdminRegistry: common.HexToAddress("0x2"),
		NonceManager:       common.HexToAddress("0x3"),
	}, offramp.OffRampDynamicConfig{
		FeeQuoter:                               common.HexToAddress("0x4"),
		PermissionLessExecutionThresholdSeconds: uint32((8 * time.Hour).Seconds()),
		MessageInterceptor:                      common.HexToAddress("0x5"),
	}, nil)
	require.NoError(t, err)
	chain.Commit()

	// onramp
	_, _, _, err = onramp.DeployOnRamp(owner, chain, onramp.OnRampStaticConfig{
		ChainSelector:      1,
		RmnRemote:          common.HexToAddress("0x1"),
		NonceManager:       common.HexToAddress("0x2"),
		TokenAdminRegistry: common.HexToAddress("0x3"),
	}, onramp.OnRampDynamicConfig{
		FeeQuoter:          common.HexToAddress("0x4"),
		MessageInterceptor: common.HexToAddress("0x5"),
		FeeAggregator:      common.HexToAddress("0x6"),
		AllowListAdmin:     common.HexToAddress("0x7"),
	}, nil)
	require.NoError(t, err)
	chain.Commit()

	// fee quoter
	_, _, _, err = fee_quoter.DeployFeeQuoter(
		owner,
		chain,
		fee_quoter.FeeQuoterStaticConfig{
			MaxFeeJuelsPerMsg:            big.NewInt(1e18),
			LinkToken:                    common.HexToAddress("0x1"),
			TokenPriceStalenessThreshold: 10,
		},
		[]common.Address{common.HexToAddress("0x1")},
		[]common.Address{common.HexToAddress("0x2")},
		[]fee_quoter.FeeQuoterTokenPriceFeedUpdate{},
		[]fee_quoter.FeeQuoterTokenTransferFeeConfigArgs{},
		[]fee_quoter.FeeQuoterPremiumMultiplierWeiPerEthArgs{},
		[]fee_quoter.FeeQuoterDestChainConfigArgs{})
	require.NoError(t, err)
	chain.Commit()

	// token admin registry
	_, _, _, err = token_admin_registry.DeployTokenAdminRegistry(owner, chain)
	require.NoError(t, err)
	chain.Commit()

	// TODO: add rmn home and rmn remote
}
