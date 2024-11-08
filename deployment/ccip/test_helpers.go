package ccipdeployment

import (
	"context"
	"fmt"
	"math/big"
	"sort"
	"testing"
	"time"

	mapset "github.com/deckarep/golang-set/v2"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/shared/generated/burn_mint_erc677"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"go.uber.org/multierr"
	"go.uber.org/zap/zapcore"

	chainsel "github.com/smartcontractkit/chain-selectors"

	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/mock_ethusd_aggregator_wrapper"

	"github.com/smartcontractkit/chainlink/deployment/environment/memory"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink/deployment/environment/devenv"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/burn_mint_token_pool"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/mock_v3_aggregator_contract"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/aggregator_v3_interface"
)

const (
	HomeChainIndex = 0
	FeedChainIndex = 1
)

// Context returns a context with the test's deadline, if available.
func Context(tb testing.TB) context.Context {
	ctx := context.Background()
	var cancel func()
	switch t := tb.(type) {
	case *testing.T:
		if d, ok := t.Deadline(); ok {
			ctx, cancel = context.WithDeadline(ctx, d)
		}
	}
	if cancel == nil {
		ctx, cancel = context.WithCancel(ctx)
	}
	tb.Cleanup(cancel)
	return ctx
}

type DeployedEnv struct {
	Env          deployment.Environment
	HomeChainSel uint64
	FeedChainSel uint64
	ReplayBlocks map[uint64]uint64
}

func (e *DeployedEnv) SetupJobs(t *testing.T) {
	ctx := testcontext.Get(t)
	jbs, err := NewCCIPJobSpecs(e.Env.NodeIDs, e.Env.Offchain)
	require.NoError(t, err)
	for nodeID, jobs := range jbs {
		for _, job := range jobs {
			// Note these auto-accept
			_, err := e.Env.Offchain.ProposeJob(ctx,
				&jobv1.ProposeJobRequest{
					NodeId: nodeID,
					Spec:   job,
				})
			require.NoError(t, err)
		}
	}
	// Wait for plugins to register filters?
	// TODO: Investigate how to avoid.
	time.Sleep(30 * time.Second)
	ReplayLogs(t, e.Env.Offchain, e.ReplayBlocks)
}

func ReplayLogs(t *testing.T, oc deployment.OffchainClient, replayBlocks map[uint64]uint64) {
	switch oc := oc.(type) {
	case *memory.JobClient:
		require.NoError(t, oc.ReplayLogs(replayBlocks))
	case *devenv.JobDistributor:
		require.NoError(t, oc.ReplayLogs(replayBlocks))
	default:
		t.Fatalf("unsupported offchain client type %T", oc)
	}
}

func DeployTestContracts(t *testing.T,
	lggr logger.Logger,
	ab deployment.AddressBook,
	homeChainSel,
	feedChainSel uint64,
	chains map[uint64]deployment.Chain,
) deployment.CapabilityRegistryConfig {
	capReg, err := DeployCapReg(lggr,
		// deploying cap reg for the first time on a blank chain state
		CCIPOnChainState{
			Chains: make(map[uint64]CCIPChainState),
		}, ab, chains[homeChainSel])
	require.NoError(t, err)
	_, err = DeployFeeds(lggr, ab, chains[feedChainSel])
	require.NoError(t, err)
	err = DeployFeeTokensToChains(lggr, ab, chains)
	require.NoError(t, err)
	evmChainID, err := chainsel.ChainIdFromSelector(homeChainSel)
	require.NoError(t, err)
	return deployment.CapabilityRegistryConfig{
		EVMChainID: evmChainID,
		Contract:   capReg.Address,
	}
}

func LatestBlocksByChain(ctx context.Context, chains map[uint64]deployment.Chain) (map[uint64]uint64, error) {
	latestBlocks := make(map[uint64]uint64)
	for _, chain := range chains {
		latesthdr, err := chain.Client.HeaderByNumber(ctx, nil)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get latest header for chain %d", chain.Selector)
		}
		block := latesthdr.Number.Uint64()
		latestBlocks[chain.Selector] = block
	}
	return latestBlocks, nil
}

func allocateCCIPChainSelectors(chains map[uint64]deployment.Chain) (homeChainSel uint64, feeChainSel uint64) {
	// Lower chainSel is home chain.
	var chainSels []uint64
	// Say first chain is home chain.
	for chainSel := range chains {
		chainSels = append(chainSels, chainSel)
	}
	sort.Slice(chainSels, func(i, j int) bool {
		return chainSels[i] < chainSels[j]
	})
	// Take lowest for determinism.
	return chainSels[HomeChainIndex], chainSels[FeedChainIndex]
}

// NewMemoryEnvironment creates a new CCIP environment
// with capreg, fee tokens, feeds and nodes set up.
func NewMemoryEnvironment(t *testing.T, lggr logger.Logger, numChains int, numNodes int) DeployedEnv {
	require.GreaterOrEqual(t, numChains, 2, "numChains must be at least 2 for home and feed chains")
	require.GreaterOrEqual(t, numNodes, 4, "numNodes must be at least 4")
	ctx := testcontext.Get(t)
	chains := memory.NewMemoryChains(t, numChains)
	homeChainSel, feedSel := allocateCCIPChainSelectors(chains)
	replayBlocks, err := LatestBlocksByChain(ctx, chains)
	require.NoError(t, err)

	ab := deployment.NewMemoryAddressBook()
	crConfig := DeployTestContracts(t, lggr, ab, homeChainSel, feedSel, chains)
	nodes := memory.NewNodes(t, zapcore.InfoLevel, chains, numNodes, 1, crConfig)
	for _, node := range nodes {
		require.NoError(t, node.App.Start(ctx))
		t.Cleanup(func() {
			require.NoError(t, node.App.Stop())
		})
	}
	e := memory.NewMemoryEnvironmentFromChainsNodes(t, lggr, chains, nodes)
	envNodes, err := deployment.NodeInfo(e.NodeIDs, e.Offchain)
	require.NoError(t, err)
	e.ExistingAddresses = ab
	_, err = DeployHomeChain(lggr, e, e.ExistingAddresses, chains[homeChainSel],
		NewTestRMNStaticConfig(),
		NewTestRMNDynamicConfig(),
		NewTestNodeOperator(chains[homeChainSel].DeployerKey.From),
		map[string][][32]byte{
			"NodeOperator": envNodes.NonBootstraps().PeerIDs(),
		},
	)
	require.NoError(t, err)

	return DeployedEnv{
		Env:          e,
		HomeChainSel: homeChainSel,
		FeedChainSel: feedSel,
		ReplayBlocks: replayBlocks,
	}
}

func NewMemoryEnvironmentWithJobs(t *testing.T, lggr logger.Logger, numChains int, numNodes int) DeployedEnv {
	e := NewMemoryEnvironment(t, lggr, numChains, numNodes)
	e.SetupJobs(t)
	return e
}

func CCIPSendRequest(
	e deployment.Environment,
	state CCIPOnChainState,
	src, dest uint64,
	data []byte,
	tokensAndAmounts []router.ClientEVMTokenAmount,
	feeToken common.Address,
	testRouter bool,
	extraArgs []byte,
) (*types.Transaction, uint64, error) {
	msg := router.ClientEVM2AnyMessage{
		Receiver:     common.LeftPadBytes(state.Chains[dest].Receiver.Address().Bytes(), 32),
		Data:         data,
		TokenAmounts: tokensAndAmounts,
		FeeToken:     feeToken,
		ExtraArgs:    extraArgs,
	}
	r := state.Chains[src].Router
	if testRouter {
		r = state.Chains[src].TestRouter
	}
	fee, err := r.GetFee(
		&bind.CallOpts{Context: context.Background()}, dest, msg)
	if err != nil {
		return nil, 0, errors.Wrap(deployment.MaybeDataErr(err), "failed to get fee")
	}
	if msg.FeeToken == common.HexToAddress("0x0") {
		e.Chains[src].DeployerKey.Value = fee
		defer func() { e.Chains[src].DeployerKey.Value = nil }()
	}
	tx, err := r.CcipSend(
		e.Chains[src].DeployerKey,
		dest,
		msg)
	if err != nil {
		return nil, 0, errors.Wrap(err, "failed to send CCIP message")
	}
	blockNum, err := e.Chains[src].Confirm(tx)
	if err != nil {
		return tx, 0, errors.Wrap(err, "failed to confirm CCIP message")
	}
	return tx, blockNum, nil
}

func TestSendRequest(t *testing.T, e deployment.Environment, state CCIPOnChainState, src, dest uint64, testRouter bool, tokensAndAmounts []router.ClientEVMTokenAmount) uint64 {
	t.Logf("Sending CCIP request from chain selector %d to chain selector %d",
		src, dest)
	tx, blockNum, err := CCIPSendRequest(e, state, src, dest, []byte("hello"), tokensAndAmounts, common.HexToAddress("0x0"), testRouter, nil)
	require.NoError(t, err)
	it, err := state.Chains[src].OnRamp.FilterCCIPMessageSent(&bind.FilterOpts{
		Start:   blockNum,
		End:     &blockNum,
		Context: context.Background(),
	}, []uint64{dest}, []uint64{})
	require.NoError(t, err)
	require.True(t, it.Next())
	seqNum := it.Event.Message.Header.SequenceNumber
	t.Logf("CCIP message sent from chain selector %d to chain selector %d tx %s seqNum %d", src, dest, tx.Hash().String(), seqNum)
	return seqNum
}

// AddLanesForAll adds densely connected lanes for all chains in the environment so that each chain
// is connected to every other chain except itself.
func AddLanesForAll(e deployment.Environment, state CCIPOnChainState) error {
	for source := range e.Chains {
		for dest := range e.Chains {
			if source != dest {
				err := AddLane(e, state, source, dest)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

const (
	// MockLinkAggregatorDescription This is the description of the MockV3Aggregator.sol contract
	// nolint:lll
	// https://github.com/smartcontractkit/chainlink/blob/a348b98e90527520049c580000a86fb8ceff7fa7/contracts/src/v0.8/tests/MockV3Aggregator.sol#L76-L76
	MockLinkAggregatorDescription = "v0.8/tests/MockV3Aggregator.sol"
	// MockWETHAggregatorDescription WETH use description from MockETHUSDAggregator.sol
	// nolint:lll
	// https://github.com/smartcontractkit/chainlink/blob/a348b98e90527520049c580000a86fb8ceff7fa7/contracts/src/v0.8/automation/testhelpers/MockETHUSDAggregator.sol#L19-L19
	MockWETHAggregatorDescription = "MockETHUSDAggregator"
)

var (
	MockLinkPrice = big.NewInt(5e18)
	MockWethPrice = big.NewInt(9e8)
	// MockDescriptionToTokenSymbol maps a mock feed description to token descriptor
	MockDescriptionToTokenSymbol = map[string]TokenSymbol{
		MockLinkAggregatorDescription: LinkSymbol,
		MockWETHAggregatorDescription: WethSymbol,
	}
	MockSymbolToDescription = map[TokenSymbol]string{
		LinkSymbol: MockLinkAggregatorDescription,
		WethSymbol: MockWETHAggregatorDescription,
	}
	MockSymbolToDecimals = map[TokenSymbol]uint8{
		LinkSymbol: LinkDecimals,
		WethSymbol: WethDecimals,
	}
)

func DeployFeeds(lggr logger.Logger, ab deployment.AddressBook, chain deployment.Chain) (map[string]common.Address, error) {
	linkTV := deployment.NewTypeAndVersion(PriceFeed, deployment.Version1_0_0)
	mockLinkFeed := func(chain deployment.Chain) ContractDeploy[*aggregator_v3_interface.AggregatorV3Interface] {
		linkFeed, tx, _, err1 := mock_v3_aggregator_contract.DeployMockV3Aggregator(
			chain.DeployerKey,
			chain.Client,
			LinkDecimals,  // decimals
			MockLinkPrice, // initialAnswer
		)
		aggregatorCr, err2 := aggregator_v3_interface.NewAggregatorV3Interface(linkFeed, chain.Client)

		return ContractDeploy[*aggregator_v3_interface.AggregatorV3Interface]{
			Address: linkFeed, Contract: aggregatorCr, Tv: linkTV, Tx: tx, Err: multierr.Append(err1, err2),
		}
	}

	mockWethFeed := func(chain deployment.Chain) ContractDeploy[*aggregator_v3_interface.AggregatorV3Interface] {
		wethFeed, tx, _, err1 := mock_ethusd_aggregator_wrapper.DeployMockETHUSDAggregator(
			chain.DeployerKey,
			chain.Client,
			MockWethPrice, // initialAnswer
		)
		aggregatorCr, err2 := aggregator_v3_interface.NewAggregatorV3Interface(wethFeed, chain.Client)

		return ContractDeploy[*aggregator_v3_interface.AggregatorV3Interface]{
			Address: wethFeed, Contract: aggregatorCr, Tv: linkTV, Tx: tx, Err: multierr.Append(err1, err2),
		}
	}

	linkFeedAddress, linkFeedDescription, err := deploySingleFeed(lggr, ab, chain, mockLinkFeed, LinkSymbol)
	if err != nil {
		return nil, err
	}

	wethFeedAddress, wethFeedDescription, err := deploySingleFeed(lggr, ab, chain, mockWethFeed, WethSymbol)
	if err != nil {
		return nil, err
	}

	descriptionToAddress := map[string]common.Address{
		linkFeedDescription: linkFeedAddress,
		wethFeedDescription: wethFeedAddress,
	}

	return descriptionToAddress, nil
}

func deploySingleFeed(
	lggr logger.Logger,
	ab deployment.AddressBook,
	chain deployment.Chain,
	deployFunc func(deployment.Chain) ContractDeploy[*aggregator_v3_interface.AggregatorV3Interface],
	symbol TokenSymbol,
) (common.Address, string, error) {
	//tokenTV := deployment.NewTypeAndVersion(PriceFeed, deployment.Version1_0_0)
	mockTokenFeed, err := deployContract(lggr, chain, ab, deployFunc)
	if err != nil {
		lggr.Errorw("Failed to deploy token feed", "err", err, "symbol", symbol)
		return common.Address{}, "", err
	}

	lggr.Infow("deployed mockTokenFeed", "addr", mockTokenFeed.Address)

	desc, err := mockTokenFeed.Contract.Description(&bind.CallOpts{})
	if err != nil {
		lggr.Errorw("Failed to get description", "err", err, "symbol", symbol)
		return common.Address{}, "", err
	}

	if desc != MockSymbolToDescription[symbol] {
		lggr.Errorw("Unexpected description for token", "symbol", symbol, "desc", desc)
		return common.Address{}, "", fmt.Errorf("unexpected description: %s", desc)
	}

	return mockTokenFeed.Address, desc, nil
}

func ConfirmRequestOnSourceAndDest(t *testing.T, env deployment.Environment, state CCIPOnChainState, sourceCS, destCS, expectedSeqNr uint64) error {
	latesthdr, err := env.Chains[destCS].Client.HeaderByNumber(testcontext.Get(t), nil)
	require.NoError(t, err)
	startBlock := latesthdr.Number.Uint64()
	fmt.Printf("startblock %d", startBlock)
	seqNum := TestSendRequest(t, env, state, sourceCS, destCS, false, nil)
	require.Equal(t, expectedSeqNr, seqNum)

	fmt.Printf("Request sent for seqnr %d", seqNum)
	require.NoError(t,
		ConfirmCommitWithExpectedSeqNumRange(t, env.Chains[sourceCS], env.Chains[destCS], state.Chains[destCS].OffRamp, &startBlock, cciptypes.SeqNumRange{
			cciptypes.SeqNum(seqNum),
			cciptypes.SeqNum(seqNum),
		}))

	fmt.Printf("Commit confirmed for seqnr %d", seqNum)
	require.NoError(t,
		ConfirmExecWithSeqNr(t, env.Chains[sourceCS], env.Chains[destCS], state.Chains[destCS].OffRamp, &startBlock, seqNum))

	return nil
}

func ProcessChangeset(t *testing.T, e deployment.Environment, c deployment.ChangesetOutput) {

	// TODO: Add support for jobspecs as well

	// sign and execute all proposals provided
	if len(c.Proposals) != 0 {
		state, err := LoadOnchainState(e)
		require.NoError(t, err)
		for _, prop := range c.Proposals {
			chains := mapset.NewSet[uint64]()
			for _, op := range prop.Transactions {
				chains.Add(uint64(op.ChainIdentifier))
			}

			signed := SignProposal(t, e, &prop)
			for _, sel := range chains.ToSlice() {
				ExecuteProposal(t, e, signed, state, sel)
			}
		}
	}

	// merge address books
	if c.AddressBook != nil {
		err := e.ExistingAddresses.Merge(c.AddressBook)
		require.NoError(t, err)
	}
}

func DeployTransferableToken(
	lggr logger.Logger,
	chains map[uint64]deployment.Chain,
	src, dst uint64,
	state CCIPOnChainState,
	addresses deployment.AddressBook,
	token string,
) (*burn_mint_erc677.BurnMintERC677, *burn_mint_token_pool.BurnMintTokenPool, *burn_mint_erc677.BurnMintERC677, *burn_mint_token_pool.BurnMintTokenPool, error) {
	// Deploy token and pools
	srcToken, srcPool, err := deployTransferTokenOneEnd(lggr, chains[src], addresses, token)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	dstToken, dstPool, err := deployTransferTokenOneEnd(lggr, chains[dst], addresses, token)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	// Attach token pools to registry
	if err := attachTokenToTheRegistry(chains[src], state.Chains[src], chains[src].DeployerKey, srcToken.Address(), srcPool.Address()); err != nil {
		return nil, nil, nil, nil, err
	}

	if err := attachTokenToTheRegistry(chains[dst], state.Chains[dst], chains[dst].DeployerKey, dstToken.Address(), dstPool.Address()); err != nil {
		return nil, nil, nil, nil, err
	}

	// Connect pool to each other
	if err := setTokenPoolCounterPart(chains[src], srcPool, dst, dstToken.Address(), dstPool.Address()); err != nil {
		return nil, nil, nil, nil, err
	}

	if err := setTokenPoolCounterPart(chains[dst], dstPool, src, srcToken.Address(), srcPool.Address()); err != nil {
		return nil, nil, nil, nil, err
	}

	// Add burn/mint permissions
	if err := grantMintBurnPermissions(chains[src], srcToken, srcPool.Address()); err != nil {
		return nil, nil, nil, nil, err
	}

	if err := grantMintBurnPermissions(chains[dst], dstToken, dstPool.Address()); err != nil {
		return nil, nil, nil, nil, err
	}

	return srcToken, srcPool, dstToken, dstPool, nil
}

func grantMintBurnPermissions(chain deployment.Chain, token *burn_mint_erc677.BurnMintERC677, address common.Address) error {
	tx, err := token.GrantBurnRole(chain.DeployerKey, address)
	if err != nil {
		return err
	}
	_, err = chain.Confirm(tx)
	if err != nil {
		return err
	}

	tx, err = token.GrantMintRole(chain.DeployerKey, address)
	if err != nil {
		return err
	}
	_, err = chain.Confirm(tx)
	return err
}

func setTokenPoolCounterPart(
	chain deployment.Chain,
	tokenPool *burn_mint_token_pool.BurnMintTokenPool,
	destChainSelector uint64,
	destTokenAddress common.Address,
	destTokenPoolAddress common.Address,
) error {
	tx, err := tokenPool.ApplyChainUpdates(
		chain.DeployerKey,
		[]burn_mint_token_pool.TokenPoolChainUpdate{
			{
				RemoteChainSelector: destChainSelector,
				Allowed:             true,
				RemotePoolAddress:   common.LeftPadBytes(destTokenPoolAddress.Bytes(), 32),
				RemoteTokenAddress:  common.LeftPadBytes(destTokenAddress.Bytes(), 32),
				OutboundRateLimiterConfig: burn_mint_token_pool.RateLimiterConfig{
					IsEnabled: false,
					Capacity:  big.NewInt(0),
					Rate:      big.NewInt(0),
				},
				InboundRateLimiterConfig: burn_mint_token_pool.RateLimiterConfig{
					IsEnabled: false,
					Capacity:  big.NewInt(0),
					Rate:      big.NewInt(0),
				},
			},
		},
	)
	if err != nil {
		return err
	}

	_, err = chain.Confirm(tx)
	if err != nil {
		return err
	}

	tx, err = tokenPool.SetRemotePool(
		chain.DeployerKey,
		destChainSelector,
		destTokenPoolAddress.Bytes(),
	)
	return err
}

func attachTokenToTheRegistry(
	chain deployment.Chain,
	state CCIPChainState,
	owner *bind.TransactOpts,
	token common.Address,
	tokenPool common.Address,
) error {
	tx, err := state.RegistryModule.RegisterAdminViaOwner(owner, token)
	if err != nil {
		return err
	}
	_, err = chain.Confirm(tx)
	if err != nil {
		return err
	}

	tx, err = state.TokenAdminRegistry.AcceptAdminRole(owner, token)
	if err != nil {
		return err
	}
	_, err = chain.Confirm(tx)
	if err != nil {
		return err
	}

	tx, err = state.TokenAdminRegistry.SetPool(owner, token, tokenPool)
	if err != nil {
		return err
	}
	_, err = chain.Confirm(tx)
	if err != nil {
		return err
	}
	return nil
}

func deployTransferTokenOneEnd(
	lggr logger.Logger,
	chain deployment.Chain,
	addressBook deployment.AddressBook,
	tokenSymbol string,
) (*burn_mint_erc677.BurnMintERC677, *burn_mint_token_pool.BurnMintTokenPool, error) {
	var rmnAddress, routerAddress string
	chainAddresses, err := addressBook.AddressesForChain(chain.Selector)
	if err != nil {
		return nil, nil, err
	}
	for address, v := range chainAddresses {
		if deployment.NewTypeAndVersion(ARMProxy, deployment.Version1_0_0) == v {
			rmnAddress = address
		}
		if deployment.NewTypeAndVersion(Router, deployment.Version1_2_0) == v {
			routerAddress = address
		}
		if rmnAddress != "" && routerAddress != "" {
			break
		}
	}

	tokenContract, err := deployContract(lggr, chain, addressBook,
		func(chain deployment.Chain) ContractDeploy[*burn_mint_erc677.BurnMintERC677] {
			USDCTokenAddr, tx, token, err2 := burn_mint_erc677.DeployBurnMintERC677(
				chain.DeployerKey,
				chain.Client,
				tokenSymbol,
				tokenSymbol,
				uint8(18),
				big.NewInt(0).Mul(big.NewInt(1e9), big.NewInt(1e18)),
			)
			return ContractDeploy[*burn_mint_erc677.BurnMintERC677]{
				USDCTokenAddr, token, tx, deployment.NewTypeAndVersion(BurnMintToken, deployment.Version1_0_0), err2,
			}
		})
	if err != nil {
		lggr.Errorw("Failed to deploy Token ERC677", "err", err)
		return nil, nil, err
	}

	tx, err := tokenContract.Contract.GrantMintRole(chain.DeployerKey, chain.DeployerKey.From)
	if err != nil {
		return nil, nil, err
	}
	_, err = chain.Confirm(tx)
	if err != nil {
		return nil, nil, err
	}

	tokenPool, err := deployContract(lggr, chain, addressBook,
		func(chain deployment.Chain) ContractDeploy[*burn_mint_token_pool.BurnMintTokenPool] {
			tokenPoolAddress, tx, tokenPoolContract, err2 := burn_mint_token_pool.DeployBurnMintTokenPool(
				chain.DeployerKey,
				chain.Client,
				tokenContract.Address,
				[]common.Address{},
				common.HexToAddress(rmnAddress),
				common.HexToAddress(routerAddress),
			)
			return ContractDeploy[*burn_mint_token_pool.BurnMintTokenPool]{
				tokenPoolAddress, tokenPoolContract, tx, deployment.NewTypeAndVersion(BurnMintTokenPool, deployment.Version1_0_0), err2,
			}
		})
	if err != nil {
		lggr.Errorw("Failed to deploy token pool", "err", err)
		return nil, nil, err
	}

	return tokenContract.Contract, tokenPool.Contract, nil
}
