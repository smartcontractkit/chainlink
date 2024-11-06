package smoke

import (
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/usdc_token_pool"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-ccip/pluginconfig"
	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"

	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"
	"github.com/smartcontractkit/chainlink/deployment"
	ccdeploy "github.com/smartcontractkit/chainlink/deployment/ccip"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/shared/generated/erc20"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestInitialDeployOnLocal(t *testing.T) {
	lggr := logger.TestLogger(t)
	ctx := ccdeploy.Context(t)
	tenv := ccdeploy.NewMemoryEnvironment(t, lggr, 3, 4)
	e := tenv.Env

	state, err := ccdeploy.LoadOnchainState(tenv.Env)
	require.NoError(t, err)

	feeds := state.Chains[tenv.FeedChainSel].USDFeeds
	tokenConfig := ccdeploy.NewTokenConfig()
	tokenConfig.UpsertTokenInfo(ccdeploy.LinkSymbol,
		pluginconfig.TokenInfo{
			AggregatorAddress: cciptypes.UnknownEncodedAddress(feeds[ccdeploy.LinkSymbol].Address().String()),
			Decimals:          ccdeploy.LinkDecimals,
			DeviationPPB:      cciptypes.NewBigIntFromInt64(1e9),
		},
	)
	// Apply migration
	output, err := changeset.InitialDeploy(tenv.Env, ccdeploy.DeployCCIPContractConfig{
		HomeChainSel:   tenv.HomeChainSel,
		FeedChainSel:   tenv.FeedChainSel,
		ChainsToDeploy: tenv.Env.AllChainSelectors(),
		TokenConfig:    tokenConfig,
		MCMSConfig:     ccdeploy.NewTestMCMSConfig(t, e),
		OCRSecrets:     deployment.XXXGenerateTestOCRSecrets(),
	})
	require.NoError(t, err)
	require.NoError(t, tenv.Env.ExistingAddresses.Merge(output.AddressBook))
	// Get new state after migration.
	state, err = ccdeploy.LoadOnchainState(e)
	require.NoError(t, err)

	// Ensure capreg logs are up to date.
	ccdeploy.ReplayLogs(t, e.Offchain, tenv.ReplayBlocks)

	// Apply the jobs.
	for nodeID, jobs := range output.JobSpecs {
		for _, job := range jobs {
			// Note these auto-accept
			_, err := e.Offchain.ProposeJob(ctx,
				&jobv1.ProposeJobRequest{
					NodeId: nodeID,
					Spec:   job,
				})
			require.NoError(t, err)
		}
	}

	// Add all lanes
	require.NoError(t, ccdeploy.AddLanesForAll(e, state))
	// Need to keep track of the block number for each chain so that event subscription can be done from that block.
	startBlocks := make(map[uint64]*uint64)
	// Send a message from each chain to every other chain.
	expectedSeqNum := make(map[uint64]uint64)
	for src := range e.Chains {
		for dest, destChain := range e.Chains {
			if src == dest {
				continue
			}
			latesthdr, err := destChain.Client.HeaderByNumber(testcontext.Get(t), nil)
			require.NoError(t, err)
			block := latesthdr.Number.Uint64()
			startBlocks[dest] = &block
			seqNum := ccdeploy.TestSendRequest(t, e, state, src, dest, false, nil)
			expectedSeqNum[dest] = seqNum
		}
	}

	// Wait for all commit reports to land.
	ccdeploy.ConfirmCommitForAllWithExpectedSeqNums(t, e, state, expectedSeqNum, startBlocks)

	// After commit is reported on all chains, token prices should be updated in FeeQuoter.
	for dest := range e.Chains {
		linkAddress := state.Chains[dest].LinkToken.Address()
		feeQuoter := state.Chains[dest].FeeQuoter
		timestampedPrice, err := feeQuoter.GetTokenPrice(nil, linkAddress)
		require.NoError(t, err)
		require.Equal(t, ccdeploy.MockLinkPrice, timestampedPrice.Value)
	}

	// Wait for all exec reports to land
	ccdeploy.ConfirmExecWithSeqNrForAll(t, e, state, expectedSeqNum, startBlocks)

	// TODO: Apply the proposal.
}

func TestTokenTransfer(t *testing.T) {
	lggr := logger.TestLogger(t)
	ctx := ccdeploy.Context(t)
	tenv := ccdeploy.NewMemoryEnvironment(t, lggr, 3, 4)

	e := tenv.Env
	state, err := ccdeploy.LoadOnchainState(e)
	require.NoError(t, err)

	feeds := state.Chains[tenv.FeedChainSel].USDFeeds
	tokenConfig := ccdeploy.NewTokenConfig()
	tokenConfig.UpsertTokenInfo(ccdeploy.LinkSymbol,
		pluginconfig.TokenInfo{
			AggregatorAddress: cciptypes.UnknownEncodedAddress(feeds[ccdeploy.LinkSymbol].Address().String()),
			Decimals:          ccdeploy.LinkDecimals,
			DeviationPPB:      cciptypes.NewBigIntFromInt64(1e9),
		},
	)

	// Apply migration
	output, err := changeset.InitialDeploy(e, ccdeploy.DeployCCIPContractConfig{
		HomeChainSel:   tenv.HomeChainSel,
		FeedChainSel:   tenv.FeedChainSel,
		ChainsToDeploy: e.AllChainSelectors(),
		TokenConfig:    tokenConfig,
		MCMSConfig:     ccdeploy.NewTestMCMSConfig(t, e),
		OCRSecrets:     deployment.XXXGenerateTestOCRSecrets(),
	})
	require.NoError(t, err)
	require.NoError(t, e.ExistingAddresses.Merge(output.AddressBook))
	// Get new state after migration and mock USDC token deployment.
	state, err = ccdeploy.LoadOnchainState(e)
	require.NoError(t, err)

	srcToken, _, dstToken, _, err := ccdeploy.DeployTransferableToken(
		lggr,
		tenv.Env.Chains,
		tenv.HomeChainSel,
		tenv.FeedChainSel,
		state,
		e.ExistingAddresses,
		"MY_TOKEN",
	)
	require.NoError(t, err)

	// Ensure capreg logs are up to date.
	ccdeploy.ReplayLogs(t, e.Offchain, tenv.ReplayBlocks)

	// Apply the jobs.
	for nodeID, jobs := range output.JobSpecs {
		for _, job := range jobs {
			// Note these auto-accept
			_, err := e.Offchain.ProposeJob(ctx,
				&jobv1.ProposeJobRequest{
					NodeId: nodeID,
					Spec:   job,
				})
			require.NoError(t, err)
		}
	}

	// Add all lanes
	require.NoError(t, ccdeploy.AddLanesForAll(e, state))
	// Need to keep track of the block number for each chain so that event subscription can be done from that block.
	startBlocks := make(map[uint64]*uint64)
	// Send a message from each chain to every other chain.
	expectedSeqNum := make(map[uint64]uint64)

	twoCoins := new(big.Int).Mul(big.NewInt(1e18), big.NewInt(2))
	tx, err := srcToken.Mint(
		e.Chains[tenv.HomeChainSel].DeployerKey,
		e.Chains[tenv.HomeChainSel].DeployerKey.From,
		new(big.Int).Mul(twoCoins, big.NewInt(10)),
	)
	require.NoError(t, err)
	_, err = e.Chains[tenv.HomeChainSel].Confirm(tx)
	require.NoError(t, err)

	tx, err = dstToken.Mint(
		e.Chains[tenv.FeedChainSel].DeployerKey,
		e.Chains[tenv.FeedChainSel].DeployerKey.From,
		new(big.Int).Mul(twoCoins, big.NewInt(10)),
	)
	require.NoError(t, err)
	_, err = e.Chains[tenv.FeedChainSel].Confirm(tx)
	require.NoError(t, err)

	tx, err = srcToken.Approve(e.Chains[tenv.HomeChainSel].DeployerKey, state.Chains[tenv.HomeChainSel].Router.Address(), twoCoins)
	require.NoError(t, err, "failed to approve USDC tokens in home chain")
	_, err = e.Chains[tenv.HomeChainSel].Confirm(tx)
	require.NoError(t, err, "failed to confirm USDC token approval in home chain")
	tx, err = dstToken.Approve(e.Chains[tenv.FeedChainSel].DeployerKey, state.Chains[tenv.FeedChainSel].Router.Address(), twoCoins)
	require.NoError(t, err, "failed to approve USDC tokens in feed chain")
	_, err = e.Chains[tenv.FeedChainSel].Confirm(tx)
	require.NoError(t, err, "failed to confirm USDC token approval in feed chain")

	tokens := map[uint64][]router.ClientEVMTokenAmount{
		tenv.HomeChainSel: {{
			Token:  srcToken.Address(),
			Amount: twoCoins,
		}},
		tenv.FeedChainSel: {{
			Token:  dstToken.Address(),
			Amount: twoCoins,
		}},
	}

	for src := range e.Chains {
		for dest, destChain := range e.Chains {
			if src == dest {
				continue
			}
			latesthdr, err := destChain.Client.HeaderByNumber(testcontext.Get(t), nil)
			require.NoError(t, err)
			block := latesthdr.Number.Uint64()
			startBlocks[dest] = &block

			if src == tenv.HomeChainSel && dest == tenv.FeedChainSel {
				seqNum := ccdeploy.TestSendRequest(t, e, state, src, dest, false, tokens[src])
				expectedSeqNum[dest] = seqNum
			} else {
				seqNum := ccdeploy.TestSendRequest(t, e, state, src, dest, false, nil)
				expectedSeqNum[dest] = seqNum
			}
		}
	}

	// Wait for all commit reports to land.
	ccdeploy.ConfirmCommitForAllWithExpectedSeqNums(t, e, state, expectedSeqNum, startBlocks)

	// After commit is reported on all chains, token prices should be updated in FeeQuoter.

	for dest := range e.Chains {
		linkAddress := state.Chains[dest].LinkToken.Address()
		feeQuoter := state.Chains[dest].FeeQuoter
		timestampedPrice, err := feeQuoter.GetTokenPrice(nil, linkAddress)
		require.NoError(t, err)
		require.Equal(t, ccdeploy.MockLinkPrice, timestampedPrice.Value)
	}

	// Wait for all exec reports to land
	ccdeploy.ConfirmExecWithSeqNrForAll(t, e, state, expectedSeqNum, startBlocks)

	balance, err := dstToken.BalanceOf(nil, state.Chains[tenv.FeedChainSel].Receiver.Address())
	require.NoError(t, err)
	require.Equal(t, twoCoins, balance)
}

func TestUSDCTokenTransfer(t *testing.T) {
	lggr := logger.TestLogger(t)
	ctx := ccdeploy.Context(t)
	tenv := ccdeploy.NewMemoryEnvironment(t, lggr, 3, 4)

	e := tenv.Env
	// if it's a new USDC deployment, set up mock server for attestation,
	// we need to set it only once for all the lanes as the attestation path uses regex to match the path for
	// all messages across all lanes
	//err := actions.SetMockServerWithUSDCAttestation(e.MockAdapter, nil)
	//require.NoError(t, err, "failed to set up mock server for attestation")
	state, err := ccdeploy.LoadOnchainState(e)
	require.NoError(t, err)

	feeds := state.Chains[tenv.FeedChainSel].USDFeeds
	tokenConfig := ccdeploy.NewTokenConfig()
	tokenConfig.UpsertTokenInfo(ccdeploy.LinkSymbol,
		pluginconfig.TokenInfo{
			AggregatorAddress: cciptypes.UnknownEncodedAddress(feeds[ccdeploy.LinkSymbol].Address().String()),
			Decimals:          ccdeploy.LinkDecimals,
			DeviationPPB:      cciptypes.NewBigIntFromInt64(1e9),
		},
	)

	// Apply migration
	output, err := changeset.InitialDeploy(e, ccdeploy.DeployCCIPContractConfig{
		HomeChainSel:   tenv.HomeChainSel,
		FeedChainSel:   tenv.FeedChainSel,
		ChainsToDeploy: e.AllChainSelectors(),
		TokenConfig:    tokenConfig,
		MCMSConfig:     ccdeploy.NewTestMCMSConfig(t, e),
		OCRSecrets:     deployment.XXXGenerateTestOCRSecrets(),
		USDCConfig: ccdeploy.USDCConfig{
			Enabled: true,
			USDCAttestationConfig: ccdeploy.USDCAttestationConfig{
				API:         "http://elo.com",
				APITimeout:  commonconfig.MustNewDuration(time.Second),
				APIInterval: commonconfig.MustNewDuration(500 * time.Millisecond),
			},
		},
	})
	require.NoError(t, err)
	require.NoError(t, e.ExistingAddresses.Merge(output.AddressBook))
	// Get new state after migration and mock USDC token deployment.
	state, err = ccdeploy.LoadOnchainState(e)
	require.NoError(t, err)

	err = ccdeploy.SyncUSDCDomains(lggr, e.Chains, tenv.HomeChainSel, tenv.FeedChainSel, state)
	require.NoError(t, err, "error while syncing USDC domains")

	// Ensure capreg logs are up to date.
	ccdeploy.ReplayLogs(t, e.Offchain, tenv.ReplayBlocks)

	// Apply the jobs.
	for nodeID, jobs := range output.JobSpecs {
		for _, job := range jobs {
			// Note these auto-accept
			_, err := e.Offchain.ProposeJob(ctx,
				&jobv1.ProposeJobRequest{
					NodeId: nodeID,
					Spec:   job,
				})
			require.NoError(t, err)
		}
	}

	// Add all lanes
	require.NoError(t, ccdeploy.AddLanesForAll(e, state))
	// Need to keep track of the block number for each chain so that event subscription can be done from that block.
	startBlocks := make(map[uint64]*uint64)
	// Send a message from each chain to every other chain.
	expectedSeqNum := make(map[uint64]uint64)

	oneCoin := new(big.Int).Mul(big.NewInt(1e18), big.NewInt(1))
	// cast 677 token into ERC20Token
	homeChainUSDCtoken, err := erc20.NewERC20(state.Chains[tenv.HomeChainSel].BurnMintTokens677[ccdeploy.USDCSymbol].Address(),
		e.Chains[tenv.HomeChainSel].Client)
	require.NoError(t, err, "failed to cast USDC ERC677 token into ERC20 token in home chain")

	feedChainUSDCtoken, err := erc20.NewERC20(state.Chains[tenv.FeedChainSel].BurnMintTokens677[ccdeploy.USDCSymbol].Address(),
		e.Chains[tenv.FeedChainSel].Client)
	require.NoError(t, err, "failed to cast USDC ERC677 token into ERC20 token in feed chain")
	// approve tokens

	tx, err := state.Chains[tenv.HomeChainSel].BurnMintTokens677[ccdeploy.USDCSymbol].Mint(
		e.Chains[tenv.HomeChainSel].DeployerKey,
		e.Chains[tenv.HomeChainSel].DeployerKey.From,
		new(big.Int).Mul(oneCoin, big.NewInt(10)),
	)
	require.NoError(t, err)
	_, err = e.Chains[tenv.HomeChainSel].Confirm(tx)
	require.NoError(t, err)

	tx, err = state.Chains[tenv.FeedChainSel].BurnMintTokens677[ccdeploy.USDCSymbol].Mint(
		e.Chains[tenv.FeedChainSel].DeployerKey,
		e.Chains[tenv.FeedChainSel].DeployerKey.From,
		new(big.Int).Mul(oneCoin, big.NewInt(10)),
	)
	require.NoError(t, err)
	_, err = e.Chains[tenv.FeedChainSel].Confirm(tx)
	require.NoError(t, err)

	tx, err = homeChainUSDCtoken.Approve(e.Chains[tenv.HomeChainSel].DeployerKey, state.Chains[tenv.HomeChainSel].Router.Address(), oneCoin)
	require.NoError(t, err, "failed to approve USDC tokens in home chain")
	_, err = e.Chains[tenv.HomeChainSel].Confirm(tx)
	require.NoError(t, err, "failed to confirm USDC token approval in home chain")
	tx, err = feedChainUSDCtoken.Approve(e.Chains[tenv.FeedChainSel].DeployerKey, state.Chains[tenv.FeedChainSel].Router.Address(), oneCoin)
	require.NoError(t, err, "failed to approve USDC tokens in feed chain")
	_, err = e.Chains[tenv.FeedChainSel].Confirm(tx)
	require.NoError(t, err, "failed to confirm USDC token approval in feed chain")
	tokens := map[uint64][]router.ClientEVMTokenAmount{
		tenv.HomeChainSel: {{
			Token:  homeChainUSDCtoken.Address(),
			Amount: oneCoin,
		}},
		tenv.FeedChainSel: {{
			Token:  feedChainUSDCtoken.Address(),
			Amount: oneCoin,
		}},
	}

	AttachTokenAndPoolToRegistry(
		t,
		e.Chains[tenv.HomeChainSel],
		state.Chains[tenv.HomeChainSel],
		e.Chains[tenv.HomeChainSel].DeployerKey,
		homeChainUSDCtoken.Address(),
		state.Chains[tenv.HomeChainSel].USDCTokenPool.Address(),
	)

	AttachTokenAndPoolToRegistry(
		t,
		e.Chains[tenv.FeedChainSel],
		state.Chains[tenv.FeedChainSel],
		e.Chains[tenv.FeedChainSel].DeployerKey,
		feedChainUSDCtoken.Address(),
		state.Chains[tenv.FeedChainSel].USDCTokenPool.Address(),
	)

	tx, err = state.Chains[tenv.HomeChainSel].USDCTokenPool.ApplyChainUpdates(
		e.Chains[tenv.HomeChainSel].DeployerKey,
		[]usdc_token_pool.TokenPoolChainUpdate{
			{
				RemoteChainSelector: tenv.FeedChainSel,
				Allowed:             true,
				RemotePoolAddress:   state.Chains[tenv.FeedChainSel].USDCTokenPool.Address().Bytes(),
				RemoteTokenAddress:  feedChainUSDCtoken.Address().Bytes(),
				OutboundRateLimiterConfig: usdc_token_pool.RateLimiterConfig{
					IsEnabled: false,
					Capacity:  big.NewInt(0),
					Rate:      big.NewInt(0),
				},
				InboundRateLimiterConfig: usdc_token_pool.RateLimiterConfig{
					IsEnabled: false,
					Capacity:  big.NewInt(0),
					Rate:      big.NewInt(0),
				},
			},
		},
	)
	require.NoError(t, err)
	_, err = e.Chains[tenv.HomeChainSel].Confirm(tx)
	require.NoError(t, err)

	tx, err = state.Chains[tenv.FeedChainSel].USDCTokenPool.ApplyChainUpdates(
		e.Chains[tenv.FeedChainSel].DeployerKey,
		[]usdc_token_pool.TokenPoolChainUpdate{
			{
				RemoteChainSelector: tenv.HomeChainSel,
				Allowed:             true,
				RemotePoolAddress:   state.Chains[tenv.HomeChainSel].USDCTokenPool.Address().Bytes(),
				RemoteTokenAddress:  homeChainUSDCtoken.Address().Bytes(),
				OutboundRateLimiterConfig: usdc_token_pool.RateLimiterConfig{
					IsEnabled: false,
					Capacity:  big.NewInt(0),
					Rate:      big.NewInt(0),
				},
				InboundRateLimiterConfig: usdc_token_pool.RateLimiterConfig{
					IsEnabled: false,
					Capacity:  big.NewInt(0),
					Rate:      big.NewInt(0),
				},
			},
		},
	)
	require.NoError(t, err)
	_, err = e.Chains[tenv.FeedChainSel].Confirm(tx)
	require.NoError(t, err)

	tx, err = state.Chains[tenv.HomeChainSel].USDCTokenPool.SetRemotePool(
		e.Chains[tenv.HomeChainSel].DeployerKey,
		tenv.FeedChainSel,
		state.Chains[tenv.FeedChainSel].USDCTokenPool.Address().Bytes(),
	)
	require.NoError(t, err)
	_, err = e.Chains[tenv.HomeChainSel].Confirm(tx)
	require.NoError(t, err)

	tx, err = state.Chains[tenv.FeedChainSel].USDCTokenPool.SetRemotePool(
		e.Chains[tenv.FeedChainSel].DeployerKey,
		tenv.HomeChainSel,
		state.Chains[tenv.HomeChainSel].USDCTokenPool.Address().Bytes(),
	)
	require.NoError(t, err)
	_, err = e.Chains[tenv.FeedChainSel].Confirm(tx)
	require.NoError(t, err)

	for src := range e.Chains {
		for dest, destChain := range e.Chains {
			if src == dest {
				continue
			}
			latesthdr, err := destChain.Client.HeaderByNumber(testcontext.Get(t), nil)
			require.NoError(t, err)
			block := latesthdr.Number.Uint64()
			startBlocks[dest] = &block
			_ = tokens[src]

			seqNum := ccdeploy.TestSendRequest(t, e, state, src, dest, false, tokens[src])
			expectedSeqNum[dest] = seqNum
		}
	}

	// Wait for all commit reports to land.
	ccdeploy.ConfirmCommitForAllWithExpectedSeqNums(t, e, state, expectedSeqNum, startBlocks)

	// After commit is reported on all chains, token prices should be updated in FeeQuoter.
	for dest := range e.Chains {
		linkAddress := state.Chains[dest].LinkToken.Address()
		feeQuoter := state.Chains[dest].FeeQuoter
		timestampedPrice, err := feeQuoter.GetTokenPrice(nil, linkAddress)
		require.NoError(t, err)
		require.Equal(t, ccdeploy.MockLinkPrice, timestampedPrice.Value)
	}

	// Wait for all exec reports to land
	ccdeploy.ConfirmExecWithSeqNrForAll(t, e, state, expectedSeqNum, startBlocks)

	// TODO: Apply the proposal.
}

func AttachTokenAndPoolToRegistry(
	t *testing.T,
	chain deployment.Chain,
	state ccdeploy.CCIPChainState,
	owner *bind.TransactOpts,
	token common.Address,
	tokenPool common.Address,
) {
	tx, err := state.RegistryModule.RegisterAdminViaOwner(owner, token)
	require.NoError(t, err)
	_, err = chain.Confirm(tx)
	require.NoError(t, err)

	tx, err = state.TokenAdminRegistry.AcceptAdminRole(owner, token)
	require.NoError(t, err)
	_, err = chain.Confirm(tx)
	require.NoError(t, err)

	tx, err = state.TokenAdminRegistry.SetPool(owner, token, tokenPool)
	require.NoError(t, err)
	_, err = chain.Confirm(tx)
	require.NoError(t, err)

	homePool, err := state.TokenAdminRegistry.GetPool(nil, token)
	fmt.Print(homePool)
}
