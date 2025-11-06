package ccip

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/params"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_2_0/router"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_0/rmn_contract"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_1/token_pool"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_3/fee_quoter"
	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	evmChain "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/bindings"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/initial/burn_mint_erc677"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/initial/link_token"
	"github.com/smartcontractkit/chainlink-evm/pkg/utils"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/logging"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	v1_5testhelpers "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers/v1_5"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_5"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_5_1"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	usd_stablecoin "github.com/smartcontractkit/chainlink/deployment/ccip/shared/bindings/usd_stablecoin"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview/evm"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/deployment/environment/devenv"
	testsetups "github.com/smartcontractkit/chainlink/integration-tests/testsetups/ccip"
)

var (
	feeTokenLink   = "LINK"
	feeTokenNative = "NATIVE"
)

type balanceToken interface {
	BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error)
}

type balanceAssertion func(t *testing.T, sourceToken balanceToken, destinationToken balanceToken, address common.Address, description string)

type fastTransferE2ETestCase struct {
	name                                string
	enableFiller                        bool
	allowlistEnabled                    bool
	allowlistFiller                     bool
	tokenSymbol                         string
	preFastTransferFillerAssertions     []balanceAssertion
	postFastTransferFillerAssertions    []balanceAssertion
	postRegularTransferFillerAssertions []balanceAssertion
	preFastTransferUserAssertions       []balanceAssertion
	postFastTransferUserAssertions      []balanceAssertion
	postRegularTransferUserAssertions   []balanceAssertion
	preFastTransferPoolAssertions       []balanceAssertion
	postFastTransferPoolAssertions      []balanceAssertion
	postRegularTransferPoolAssertions   []balanceAssertion
	feeTokenType                        string // "LINK" or "NATIVE"
	fastTransferPoolFeeBps              uint16
	externalMinter                      bool
	isHybridPool                        bool // If true, the test case is for a hybrid pool
	hybridPool                          v1_5_1.Group
	settlementGasOverhead               uint32 // Used for fast transfer lane config
	expectNoExecutionError              bool
	customMaxFastTransferFee            *big.Int // Optional custom max fast transfer fee for negative testing
	expectRevert                        bool     // Whether the ccipSendToken call should revert
}

var (
	initialFillerTokenAmountOnDest = big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(1000))
	initialUserTokenAmountOnSource = big.NewInt(200000)
	defaultEthAmount               = big.NewInt(0).Mul(big.NewInt(params.Ether), big.NewInt(10000))
	transferAmount                 = big.NewInt(100000)
	expectedFastTransferFee        = big.NewInt(100)
	tokenDecimals                  = uint8(18)
	sourceChainID                  = uint64(1337)
	destinationChainID             = uint64(2337)
)

type fastTransferE2ETestCaseOption func(tc *fastTransferE2ETestCase) *fastTransferE2ETestCase

func ftfTc(name string, options ...fastTransferE2ETestCaseOption) *fastTransferE2ETestCase {
	tc := &fastTransferE2ETestCase{
		name:         name,
		enableFiller: true,
		preFastTransferFillerAssertions: []balanceAssertion{
			assertDestinationBalanceEqual(initialFillerTokenAmountOnDest),
		},
		preFastTransferUserAssertions: []balanceAssertion{
			assertSourceBalanceEqual(initialUserTokenAmountOnSource),
		},
		postFastTransferFillerAssertions:    []balanceAssertion{},
		postRegularTransferFillerAssertions: []balanceAssertion{},
		postFastTransferUserAssertions:      []balanceAssertion{},
		postRegularTransferUserAssertions:   []balanceAssertion{},
		preFastTransferPoolAssertions:       []balanceAssertion{},
		postFastTransferPoolAssertions:      []balanceAssertion{},
		postRegularTransferPoolAssertions:   []balanceAssertion{},
		feeTokenType:                        feeTokenLink,
		fastTransferPoolFeeBps:              0,
		settlementGasOverhead:               200000,
		expectNoExecutionError:              false,
	}

	for _, option := range options {
		tc = option(tc)
	}

	return tc
}

func withFillerDisabled() fastTransferE2ETestCaseOption {
	return func(tc *fastTransferE2ETestCase) *fastTransferE2ETestCase {
		tc.enableFiller = false
		return tc
	}
}

func withHybridPool(groupType v1_5_1.Group) fastTransferE2ETestCaseOption {
	return func(tc *fastTransferE2ETestCase) *fastTransferE2ETestCase {
		tc.isHybridPool = true
		tc.hybridPool = groupType
		return tc
	}
}

func withFastFillSuccessAmountAssertionsWithPoolAmount(poolAmount *big.Int, isLockRelease bool) fastTransferE2ETestCaseOption {
	transferAmountMinusFee := big.NewInt(0).Sub(transferAmount, expectedFastTransferFee)
	return func(tc *fastTransferE2ETestCase) *fastTransferE2ETestCase {
		// Calculate pool fee: (transferAmount * fastTransferPoolFeeBps) / 10000
		poolFee := big.NewInt(0).Mul(transferAmount, big.NewInt(int64(tc.fastTransferPoolFeeBps)))
		poolFee = big.NewInt(0).Div(poolFee, big.NewInt(10000))
		userReceivedAmount := big.NewInt(0).Sub(transferAmountMinusFee, poolFee)

		// Filler assertions
		tc.postRegularTransferFillerAssertions = append(tc.postRegularTransferFillerAssertions, assertDestinationBalanceEventuallyEqual(big.NewInt(0).Add(initialFillerTokenAmountOnDest, expectedFastTransferFee)))
		tc.postFastTransferFillerAssertions = append(tc.postFastTransferFillerAssertions, assertDestinationBalanceEventuallyEqual(big.NewInt(0).Sub(initialFillerTokenAmountOnDest, userReceivedAmount)))

		// User assertions
		tc.postFastTransferUserAssertions = append(tc.postFastTransferUserAssertions, assertDestinationBalanceEventuallyEqual(userReceivedAmount))
		tc.postRegularTransferUserAssertions = append(tc.postRegularTransferUserAssertions, assertDestinationBalanceEventuallyEqual(userReceivedAmount))

		// Pool assertions
		finalPoolAmount := big.NewInt(0).Add(poolAmount, poolFee)
		if isLockRelease {
			// In lock release mode we release the transfer amount from the pool this include the filler fee and the actual amount transferred to the user
			finalPoolAmount = big.NewInt(0).Sub(finalPoolAmount, transferAmount)
		}
		tc.preFastTransferPoolAssertions = append(tc.preFastTransferPoolAssertions, assertDestinationBalanceEqual(poolAmount))
		tc.postFastTransferPoolAssertions = append(tc.postFastTransferPoolAssertions, assertDestinationBalanceEventuallyEqual(poolAmount))
		tc.postRegularTransferPoolAssertions = append(tc.postRegularTransferPoolAssertions, assertDestinationBalanceEqual(finalPoolAmount))

		return tc
	}
}

func withFastFillSuccessAmountAssertions() fastTransferE2ETestCaseOption {
	return func(tc *fastTransferE2ETestCase) *fastTransferE2ETestCase {
		return withFastFillSuccessAmountAssertionsWithPoolAmount(big.NewInt(0), false)(tc)
	}
}

func withFastFillNoFillerSuccessAmountAssertionsWithPoolAmount(poolAmount *big.Int) fastTransferE2ETestCaseOption {
	return func(tc *fastTransferE2ETestCase) *fastTransferE2ETestCase {
		// Filler assertions
		tc.postFastTransferFillerAssertions = append(tc.postFastTransferFillerAssertions, assertDestinationBalanceEqual(initialFillerTokenAmountOnDest))
		tc.postRegularTransferFillerAssertions = append(tc.postRegularTransferFillerAssertions, assertDestinationBalanceEqual(initialFillerTokenAmountOnDest))

		// User assertions
		tc.postFastTransferUserAssertions = append(tc.postFastTransferUserAssertions, assertDestinationBalanceEventuallyEqual(big.NewInt(0)))
		tc.postRegularTransferUserAssertions = append(tc.postRegularTransferUserAssertions, assertDestinationBalanceEventuallyEqual(transferAmount))

		// Pool assertions
		tc.preFastTransferPoolAssertions = append(tc.preFastTransferPoolAssertions, assertDestinationBalanceEqual(poolAmount))
		tc.postFastTransferPoolAssertions = append(tc.postFastTransferPoolAssertions, assertDestinationBalanceEventuallyEqual(poolAmount))
		tc.postRegularTransferPoolAssertions = append(tc.postRegularTransferPoolAssertions, assertDestinationBalanceEqual(poolAmount))

		return tc
	}
}

func withFastFillNoFillerSuccessAmountAssertions() fastTransferE2ETestCaseOption {
	return func(tc *fastTransferE2ETestCase) *fastTransferE2ETestCase {
		return withFastFillNoFillerSuccessAmountAssertionsWithPoolAmount(big.NewInt(0))(tc)
	}
}

func withExpectNoExecutionError() fastTransferE2ETestCaseOption {
	return func(tc *fastTransferE2ETestCase) *fastTransferE2ETestCase {
		tc.expectNoExecutionError = true
		return tc
	}
}

func withSettlementGasOverhead(settlementGasOverhead uint32) fastTransferE2ETestCaseOption {
	return func(tc *fastTransferE2ETestCase) *fastTransferE2ETestCase {
		tc.settlementGasOverhead = settlementGasOverhead
		return tc
	}
}

func withFeeTokenType(feeTokenType string) fastTransferE2ETestCaseOption {
	return func(tc *fastTransferE2ETestCase) *fastTransferE2ETestCase {
		tc.feeTokenType = feeTokenType
		return tc
	}
}

func withFillerAllowlistEnabled() fastTransferE2ETestCaseOption {
	return func(tc *fastTransferE2ETestCase) *fastTransferE2ETestCase {
		tc.allowlistEnabled = true
		return tc
	}
}

func withAllowlistFiller() fastTransferE2ETestCaseOption {
	return func(tc *fastTransferE2ETestCase) *fastTransferE2ETestCase {
		tc.allowlistFiller = true
		return tc
	}
}

func withPoolFeeBps(poolFeeBps uint16) fastTransferE2ETestCaseOption {
	return func(tc *fastTransferE2ETestCase) *fastTransferE2ETestCase {
		tc.fastTransferPoolFeeBps = poolFeeBps
		return tc
	}
}

func withExternalMinter() fastTransferE2ETestCaseOption {
	return func(tc *fastTransferE2ETestCase) *fastTransferE2ETestCase {
		tc.externalMinter = true
		return tc
	}
}

func withCustomMaxFastTransferFee(fee *big.Int) fastTransferE2ETestCaseOption {
	return func(tc *fastTransferE2ETestCase) *fastTransferE2ETestCase {
		tc.customMaxFastTransferFee = fee
		return tc
	}
}

func withExpectRevert() fastTransferE2ETestCaseOption {
	return func(tc *fastTransferE2ETestCase) *fastTransferE2ETestCase {
		tc.expectRevert = true
		return tc
	}
}

var fastTransferTestCases = []*fastTransferE2ETestCase{
	ftfTc("fee token", withFeeTokenType(feeTokenLink), withFastFillSuccessAmountAssertions()),
	ftfTc("fee token and no filler", withFeeTokenType(feeTokenLink), withFastFillNoFillerSuccessAmountAssertions(), withFillerDisabled()),
	ftfTc("native fee token", withFeeTokenType(feeTokenNative), withFastFillSuccessAmountAssertions()),
	ftfTc("native fee token and no filler", withFeeTokenType(feeTokenNative), withFastFillNoFillerSuccessAmountAssertions(), withFillerDisabled()),
	ftfTc("allowlist enabled", withFillerAllowlistEnabled(), withAllowlistFiller(), withFastFillSuccessAmountAssertions()),
	ftfTc("allowlist enabled and filler not on allowlist", withFillerAllowlistEnabled(), withFastFillNoFillerSuccessAmountAssertions()),
	ftfTc("pool fee with filler", withPoolFeeBps(50), withFastFillSuccessAmountAssertions()),
	ftfTc("pool fee without filler", withPoolFeeBps(50), withFastFillNoFillerSuccessAmountAssertions(), withFillerDisabled()),
	ftfTc("external minter", withExternalMinter(), withFastFillSuccessAmountAssertions(), withFeeTokenType(feeTokenNative)),
	ftfTc("external minter feeToken", withExternalMinter(), withFastFillSuccessAmountAssertions(), withFeeTokenType(feeTokenLink)),
	ftfTc("settlement gas overhead too low", withSettlementGasOverhead(1), withExpectNoExecutionError(), withFeeTokenType(feeTokenNative)),
	ftfTc("max fast transfer fee too low", withCustomMaxFastTransferFee(big.NewInt(50)), withExpectRevert()),
	ftfTc("hybrid pool lock release", withHybridPool(v1_5_1.LockAndRelease), withFeeTokenType(feeTokenNative), withFastFillSuccessAmountAssertionsWithPoolAmount(big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(1000)), true)),
	ftfTc("hybrid pool", withHybridPool(v1_5_1.BurnAndMint), withFeeTokenType(feeTokenNative), withFastFillSuccessAmountAssertions()),
	ftfTc("hybrid pool with link fee", withHybridPool(v1_5_1.BurnAndMint), withFeeTokenType(feeTokenLink), withFastFillSuccessAmountAssertions()),
}

func assertDestinationBalanceEventuallyEqual(expectedBalance *big.Int) balanceAssertion {
	return func(t *testing.T, sourceToken balanceToken, destinationToken balanceToken, address common.Address, description string) {
		assert.EventuallyWithT(t, func(collect *assert.CollectT) {
			balance, err := destinationToken.BalanceOf(nil, address)
			assert.NoError(collect, err)
			assert.Equal(collect, expectedBalance.Int64(), balance.Int64(), "Balance should be equal to expected value")
		}, 30*time.Second, time.Second, description+" - Balance should eventually be equal to expected value")
	}
}

func assertSourceBalanceEqual(expectedBalance *big.Int) balanceAssertion {
	return func(t *testing.T, sourceToken balanceToken, destinationToken balanceToken, address common.Address, descriptuon string) {
		balance, err := sourceToken.BalanceOf(nil, address)
		require.NoError(t, err)
		require.Equal(t, expectedBalance.Int64(), balance.Int64(), descriptuon+" - Balance should be equal to expected value")
	}
}

func assertDestinationBalanceEqual(expectedBalance *big.Int) balanceAssertion {
	return func(t *testing.T, sourceToken balanceToken, destinationToken balanceToken, address common.Address, description string) {
		balance, err := destinationToken.BalanceOf(nil, address)
		require.NoError(t, err)
		require.Equal(t, expectedBalance.Int64(), balance.Int64(), description+" - Balance should be equal to expected value")
	}
}

func withdrawPoolFeesUsingChangeset(t *testing.T, env cldf.Environment, tokenSymbol string, contractType cldf.ContractType, contractVersion semver.Version, chainSelector uint64, recipientAddress common.Address, useMCMS bool) error {
	config := v1_5_1.FastTransferWithdrawPoolFeesConfig{
		TokenSymbol:     shared.TokenSymbol(tokenSymbol),
		ContractType:    contractType,
		ContractVersion: contractVersion,
		Withdrawals: map[uint64]common.Address{
			chainSelector: recipientAddress,
		},
	}

	if useMCMS {
		config.MCMS = &proposalutils.TimelockConfig{
			MinDelay:   0 * time.Second,
			MCMSAction: mcmstypes.TimelockActionSchedule,
		}
	}

	_, _, err := commonchangeset.ApplyChangesets(t, env,
		[]commonchangeset.ConfiguredChangeSet{commonchangeset.Configure(
			v1_5_1.FastTransferWithdrawPoolFeesChangeset,
			config,
		)}, commonchangeset.WithRealBackend())

	return err
}

func assertPoolFeeWithdrawal(expectedWithdrawnAmount *big.Int, env cldf.Environment, tokenSymbol string, contractType cldf.ContractType, contractVersion semver.Version, chainSelector uint64, destinationToken balanceToken, useMCMS bool, withdrawLock *sync.Mutex) balanceAssertion {
	return func(t *testing.T, sourceToken balanceToken, destinationTokenParam balanceToken, address common.Address, description string) {
		withdrawLock.Lock()
		defer withdrawLock.Unlock()

		recipientAddress, _, _ := createAccount(t, destinationChainID)

		pool, err := bindings.GetFastTransferTokenPoolContract(env, shared.TokenSymbol(tokenSymbol), contractType, contractVersion, chainSelector)
		require.NoError(t, err)

		accumulatedFees, err := pool.GetAccumulatedPoolFees(nil)
		require.NoError(t, err)
		require.Equal(t, expectedWithdrawnAmount.Int64(), accumulatedFees.Int64(), description+" - Accumulated pool fees should match expected amount")

		recipientBalanceBefore, err := destinationToken.BalanceOf(nil, recipientAddress)
		require.NoError(t, err)

		err = withdrawPoolFeesUsingChangeset(t, env, tokenSymbol, contractType, contractVersion, chainSelector, recipientAddress, useMCMS)
		require.NoError(t, err)

		recipientBalanceAfter, err := destinationToken.BalanceOf(nil, recipientAddress)
		require.NoError(t, err)

		expectedRecipientBalance := big.NewInt(0).Add(recipientBalanceBefore, expectedWithdrawnAmount)
		require.Equal(t, expectedRecipientBalance.Int64(), recipientBalanceAfter.Int64(), description+" - Recipient should receive the withdrawn pool fees")

		accumulatedFeesAfter, err := pool.GetAccumulatedPoolFees(nil)
		require.NoError(t, err)
		require.Equal(t, int64(0), accumulatedFeesAfter.Int64(), description+" - Pool accumulated fees should be zero after withdrawal")
	}
}

func createAccount(t *testing.T, chainID uint64) (common.Address, func() *bind.TransactOpts, *ecdsa.PrivateKey) {
	userPrivateKey, err := crypto.GenerateKey()
	require.NoError(t, err)
	userAddress := crypto.PubkeyToAddress(userPrivateKey.PublicKey)
	transactor := func() *bind.TransactOpts {
		userTransactor, err := bind.NewKeyedTransactorWithChainID(userPrivateKey, new(big.Int).SetUint64(chainID))
		require.NoError(t, err)
		return userTransactor
	}
	return userAddress, transactor, userPrivateKey
}

func deployTokenAndGrantAllRoles(t *testing.T, chain evmChain.Chain, tokenSymbol string, tokenDecimals uint8, lock *sync.Mutex, isExternalMinterToken bool) token {
	lock.Lock()
	defer lock.Unlock()

	if isExternalMinterToken {
		_, tx, token, err := usd_stablecoin.DeployStablecoin(
			chain.DeployerKey,
			chain.Client,
		)
		require.NoError(t, err)
		_, err = chain.Confirm(tx)
		require.NoError(t, err)

		tx, err = token.Initialize(chain.DeployerKey, tokenSymbol, tokenSymbol)
		require.NoError(t, err)
		_, err = chain.Confirm(tx)
		require.NoError(t, err)

		return token
	}

	_, tx, token, err := burn_mint_erc677.DeployBurnMintERC677(
		chain.DeployerKey,
		chain.Client,
		tokenSymbol,
		tokenSymbol,
		tokenDecimals,
		big.NewInt(0).Mul(big.NewInt(1e9), big.NewInt(1e18)),
	)
	require.NoError(t, err)
	_, err = chain.Confirm(tx)
	require.NoError(t, err)

	tx, err = token.GrantMintAndBurnRoles(chain.DeployerKey, chain.DeployerKey.From)
	require.NoError(t, err)
	_, err = chain.Confirm(tx)
	require.NoError(t, err)

	return token
}

func getLinkTokenAndGrantMintRole(t *testing.T, chain evmChain.Chain, state evm.CCIPChainState, sendLock *sync.Mutex) *link_token.LinkToken {
	sendLock.Lock()
	defer sendLock.Unlock()
	linkToken := state.LinkToken
	tx, err := linkToken.GrantMintRole(chain.DeployerKey, chain.DeployerKey.From)
	require.NoError(t, err)
	_, err = chain.Confirm(tx)
	require.NoError(t, err)

	return linkToken
}

type mintableToken interface {
	Mint(opts *bind.TransactOpts, account common.Address, amount *big.Int) (*types.Transaction, error)
	Address() common.Address
}

type token interface {
	balanceToken
	approvableToken
	Address() common.Address
}

type sequenceNumberRetriever func(opts *bind.CallOpts, destChainSelector uint64) (uint64, error)
type waitForExecutionFn func(t *testing.T, sequenceNumber uint64)

func fundAccountWithToken(t *testing.T, chain evmChain.Chain, receiver common.Address, token mintableToken, amount *big.Int, sendLock *sync.Mutex) {
	sendLock.Lock()
	defer sendLock.Unlock()
	tx, err := token.Mint(chain.DeployerKey, receiver, amount)
	require.NoError(t, err)
	_, err = chain.Confirm(tx)
	require.NoError(t, err)
}

type approvableToken interface {
	Approve(opts *bind.TransactOpts, spender common.Address, amount *big.Int) (*types.Transaction, error)
}

func approveToken(t *testing.T, chain evmChain.Chain, transactor *bind.TransactOpts, token approvableToken, spender common.Address, lock *sync.Mutex) {
	lock.Lock()
	defer lock.Unlock()
	tx, err := token.Approve(transactor, spender, big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(1e9))) // Approve a large amount
	require.NoError(t, err)
	_, err = chain.Confirm(tx)
	require.NoError(t, err)
}

type tokenPoolConfig struct {
	poolConfig        map[uint64]v1_5_1.DeployTokenPoolInput
	sourceMinter      mintableToken
	destinationMinter mintableToken
	postSetupAction   func(sourceTokenPool common.Address, destinationTokenPool common.Address)
	version           semver.Version
	poolType          cldf.ContractType
}

func configureExternalMinterTokenPool(t *testing.T, e cldf.Environment, sourceChainSelector, destinationChainSelector uint64, sourceTokenAddress, destinationTokenAddress common.Address, tokenDecimals uint8) tokenPoolConfig {
	sourceChain := e.BlockChains.EVMChains()[sourceChainSelector]
	destChain := e.BlockChains.EVMChains()[destinationChainSelector]

	_, sourceTokenGovernor := testhelpers.DeployTokenGovernor(t, e, sourceChainSelector, sourceTokenAddress)
	_, destinationTokenGovernor := testhelpers.DeployTokenGovernor(t, e, destinationChainSelector, destinationTokenAddress)

	bridgeBurnMintRole, err := sourceTokenGovernor.BRIDGEMINTERORBURNERROLE(nil)
	require.NoError(t, err)

	poolConfig := map[uint64]v1_5_1.DeployTokenPoolInput{
		sourceChainSelector: {
			Type:               shared.BurnMintWithExternalMinterFastTransferTokenPool,
			TokenAddress:       sourceTokenAddress,
			AllowList:          nil,
			LocalTokenDecimals: tokenDecimals,
			AcceptLiquidity:    nil,
			ExternalMinter:     sourceTokenGovernor.Address(),
		},
		destinationChainSelector: {
			Type:               shared.BurnMintWithExternalMinterFastTransferTokenPool,
			TokenAddress:       destinationTokenAddress,
			AllowList:          nil,
			LocalTokenDecimals: tokenDecimals,
			AcceptLiquidity:    nil,
			ExternalMinter:     destinationTokenGovernor.Address(),
		},
	}

	postSetupAction := func(sourceTokenPool common.Address, destinationTokenPool common.Address) {
		tx, err := sourceTokenGovernor.GrantRole(sourceChain.DeployerKey, bridgeBurnMintRole, sourceTokenPool)
		require.NoError(t, err)
		_, err = sourceChain.Confirm(tx)
		require.NoError(t, err)
		tx, err = destinationTokenGovernor.GrantRole(destChain.DeployerKey, bridgeBurnMintRole, destinationTokenPool)
		require.NoError(t, err)
		_, err = destChain.Confirm(tx)
		require.NoError(t, err)

		sourceToken, err := usd_stablecoin.NewStablecoin(sourceTokenAddress, sourceChain.Client)
		require.NoError(t, err)
		tx, err = sourceToken.TransferOwnership(sourceChain.DeployerKey, sourceTokenGovernor.Address())
		require.NoError(t, err)
		_, err = sourceChain.Confirm(tx)
		require.NoError(t, err)

		tx, err = sourceTokenGovernor.AcceptOwnership(sourceChain.DeployerKey)
		require.NoError(t, err)
		_, err = sourceChain.Confirm(tx)
		require.NoError(t, err)

		destinationToken, err := usd_stablecoin.NewStablecoin(destinationTokenAddress, destChain.Client)
		require.NoError(t, err)
		tx, err = destinationToken.TransferOwnership(destChain.DeployerKey, destinationTokenGovernor.Address())
		require.NoError(t, err)
		_, err = destChain.Confirm(tx)
		require.NoError(t, err)
		tx, err = destinationTokenGovernor.AcceptOwnership(destChain.DeployerKey)
		require.NoError(t, err)
		_, err = destChain.Confirm(tx)
		require.NoError(t, err)

		minterRole, err := sourceTokenGovernor.MINTERROLE(nil)
		require.NoError(t, err)
		tx, err = sourceTokenGovernor.GrantRole(sourceChain.DeployerKey, minterRole, sourceChain.DeployerKey.From)
		require.NoError(t, err)
		_, err = sourceChain.Confirm(tx)
		require.NoError(t, err)
		tx, err = destinationTokenGovernor.GrantRole(destChain.DeployerKey, minterRole, destChain.DeployerKey.From)
		require.NoError(t, err)
		_, err = destChain.Confirm(tx)
		require.NoError(t, err)
	}

	return tokenPoolConfig{
		poolConfig:        poolConfig,
		sourceMinter:      sourceTokenGovernor,
		destinationMinter: destinationTokenGovernor,
		postSetupAction:   postSetupAction,
		version:           shared.BurnMintWithExternalMinterFastTransferTokenPoolVersion,
		poolType:          shared.BurnMintWithExternalMinterFastTransferTokenPool,
	}
}

func configureHybridTokenPool(t *testing.T, e cldf.Environment, sourceChainSelector, destinationChainSelector uint64, sourceTokenAddress, destinationTokenAddress common.Address, tokenDecimals uint8, tokenSymbol shared.TokenSymbol, useMcms bool, groupType v1_5_1.Group) tokenPoolConfig {
	sourceChain := e.BlockChains.EVMChains()[sourceChainSelector]
	destChain := e.BlockChains.EVMChains()[destinationChainSelector]

	_, sourceTokenGovernor := testhelpers.DeployTokenGovernor(t, e, sourceChainSelector, sourceTokenAddress)
	_, destinationTokenGovernor := testhelpers.DeployTokenGovernor(t, e, destinationChainSelector, destinationTokenAddress)

	bridgeBurnMintRole, err := sourceTokenGovernor.BRIDGEMINTERORBURNERROLE(nil)
	require.NoError(t, err)

	poolConfig := map[uint64]v1_5_1.DeployTokenPoolInput{
		sourceChainSelector: {
			Type:               shared.HybridWithExternalMinterFastTransferTokenPool,
			TokenAddress:       sourceTokenAddress,
			AllowList:          nil,
			LocalTokenDecimals: tokenDecimals,
			AcceptLiquidity:    nil,
			ExternalMinter:     sourceTokenGovernor.Address(),
		},
		destinationChainSelector: {
			Type:               shared.HybridWithExternalMinterFastTransferTokenPool,
			TokenAddress:       destinationTokenAddress,
			AllowList:          nil,
			LocalTokenDecimals: tokenDecimals,
			AcceptLiquidity:    nil,
			ExternalMinter:     destinationTokenGovernor.Address(),
		},
	}

	postSetupAction := func(sourceTokenPool common.Address, destinationTokenPool common.Address) {
		tx, err := sourceTokenGovernor.GrantRole(sourceChain.DeployerKey, bridgeBurnMintRole, sourceTokenPool)
		require.NoError(t, err)
		_, err = sourceChain.Confirm(tx)
		require.NoError(t, err)
		tx, err = destinationTokenGovernor.GrantRole(destChain.DeployerKey, bridgeBurnMintRole, destinationTokenPool)
		require.NoError(t, err)
		_, err = destChain.Confirm(tx)
		require.NoError(t, err)

		sourceToken, err := usd_stablecoin.NewStablecoin(sourceTokenAddress, sourceChain.Client)
		require.NoError(t, err)
		tx, err = sourceToken.TransferOwnership(sourceChain.DeployerKey, sourceTokenGovernor.Address())
		require.NoError(t, err)
		_, err = sourceChain.Confirm(tx)
		require.NoError(t, err)

		tx, err = sourceTokenGovernor.AcceptOwnership(sourceChain.DeployerKey)
		require.NoError(t, err)
		_, err = sourceChain.Confirm(tx)
		require.NoError(t, err)

		destinationToken, err := usd_stablecoin.NewStablecoin(destinationTokenAddress, destChain.Client)
		require.NoError(t, err)
		tx, err = destinationToken.TransferOwnership(destChain.DeployerKey, destinationTokenGovernor.Address())
		require.NoError(t, err)
		_, err = destChain.Confirm(tx)
		require.NoError(t, err)
		tx, err = destinationTokenGovernor.AcceptOwnership(destChain.DeployerKey)
		require.NoError(t, err)
		_, err = destChain.Confirm(tx)
		require.NoError(t, err)

		minterRole, err := sourceTokenGovernor.MINTERROLE(nil)
		require.NoError(t, err)
		tx, err = sourceTokenGovernor.GrantRole(sourceChain.DeployerKey, minterRole, sourceChain.DeployerKey.From)
		require.NoError(t, err)
		_, err = sourceChain.Confirm(tx)
		require.NoError(t, err)
		tx, err = destinationTokenGovernor.GrantRole(destChain.DeployerKey, minterRole, destChain.DeployerKey.From)
		require.NoError(t, err)
		_, err = destChain.Confirm(tx)
		require.NoError(t, err)

		poolAmount := big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(1000))
		tx, err = sourceTokenGovernor.Mint(sourceChain.DeployerKey, sourceTokenPool, poolAmount)
		require.NoError(t, err)
		_, err = sourceChain.Confirm(tx)
		require.NoError(t, err)

		tx, err = destinationTokenGovernor.Mint(destChain.DeployerKey, destinationTokenPool, poolAmount)
		require.NoError(t, err)
		_, err = destChain.Confirm(tx)
		require.NoError(t, err)

		createGroupConfig := func(group v1_5_1.Group) v1_5_1.HybridTokenPoolUpdateGroupsConfig {
			config := v1_5_1.HybridTokenPoolUpdateGroupsConfig{
				TokenSymbol:     tokenSymbol,
				ContractType:    shared.HybridWithExternalMinterFastTransferTokenPool,
				ContractVersion: shared.HybridWithExternalMinterFastTransferTokenPoolVersion,
				Updates: map[uint64][]v1_5_1.GroupUpdateConfig{
					sourceChainSelector: {
						{
							RemoteChainSelector: destinationChainSelector,
							Group:               group,
							RemoteChainSupply:   poolAmount,
						},
					},
					destinationChainSelector: {
						{
							RemoteChainSelector: sourceChainSelector,
							Group:               group,
							RemoteChainSupply:   poolAmount,
						},
					},
				},
			}

			if useMcms {
				config.MCMS = &proposalutils.TimelockConfig{
					MinDelay:   0 * time.Second,
					MCMSAction: mcmstypes.TimelockActionSchedule,
				}
			}

			return config
		}

		// Apply BurnAndMint configuration
		_, _, err = commonchangeset.ApplyChangesets(t, e,
			[]commonchangeset.ConfiguredChangeSet{commonchangeset.Configure(
				v1_5_1.HybridTokenPoolUpdateGroupsChangeset,
				createGroupConfig(v1_5_1.BurnAndMint),
			)}, commonchangeset.WithRealBackend())
		require.NoError(t, err)

		// Apply LockAndRelease configuration if needed
		if groupType == v1_5_1.LockAndRelease {
			_, _, err = commonchangeset.ApplyChangesets(t, e,
				[]commonchangeset.ConfiguredChangeSet{commonchangeset.Configure(
					v1_5_1.HybridTokenPoolUpdateGroupsChangeset,
					createGroupConfig(v1_5_1.LockAndRelease),
				)}, commonchangeset.WithRealBackend())
			require.NoError(t, err)
		}
	}

	return tokenPoolConfig{
		poolConfig:        poolConfig,
		sourceMinter:      sourceTokenGovernor,
		destinationMinter: destinationTokenGovernor,
		postSetupAction:   postSetupAction,
		version:           shared.HybridWithExternalMinterFastTransferTokenPoolVersion,
		poolType:          shared.HybridWithExternalMinterFastTransferTokenPool,
	}
}

func configureBurnMintTokenPool(t *testing.T, e cldf.Environment, sourceChainSelector, destinationChainSelector uint64, sourceTokenAddress, destinationTokenAddress common.Address, tokenDecimals uint8) tokenPoolConfig {
	sourceChain := e.BlockChains.EVMChains()[sourceChainSelector]
	destChain := e.BlockChains.EVMChains()[destinationChainSelector]

	poolConfig := map[uint64]v1_5_1.DeployTokenPoolInput{
		sourceChainSelector: {
			Type:               shared.BurnMintFastTransferTokenPool,
			TokenAddress:       sourceTokenAddress,
			AllowList:          nil,
			LocalTokenDecimals: tokenDecimals,
			AcceptLiquidity:    nil,
		},
		destinationChainSelector: {
			Type:               shared.BurnMintFastTransferTokenPool,
			TokenAddress:       destinationTokenAddress,
			AllowList:          nil,
			LocalTokenDecimals: tokenDecimals,
			AcceptLiquidity:    nil,
		},
	}

	sourceToken, err := burn_mint_erc677.NewBurnMintERC677(sourceTokenAddress, sourceChain.Client)
	require.NoError(t, err)
	destToken, err := burn_mint_erc677.NewBurnMintERC677(destinationTokenAddress, destChain.Client)
	require.NoError(t, err)

	postSetupAction := func(sourceTokenPool common.Address, destinationTokenPool common.Address) {
		sourceTokenInstance, err := burn_mint_erc677.NewBurnMintERC677(sourceTokenAddress, sourceChain.Client)
		require.NoError(t, err)
		tx, err := sourceTokenInstance.GrantBurnRole(sourceChain.DeployerKey, sourceTokenPool)
		require.NoError(t, err)
		_, err = sourceChain.Confirm(tx)
		require.NoError(t, err)

		tx, err = destToken.GrantMintRole(destChain.DeployerKey, destinationTokenPool)
		require.NoError(t, err)
		_, err = destChain.Confirm(tx)
		require.NoError(t, err)
	}

	return tokenPoolConfig{
		poolConfig:        poolConfig,
		sourceMinter:      sourceToken,
		destinationMinter: destToken,
		postSetupAction:   postSetupAction,
		version:           shared.FastTransferTokenPoolVersion,
		poolType:          shared.BurnMintFastTransferTokenPool,
	}
}

func configureTokenPoolRateLimits(e cldf.Environment, tokenSymbol string, sourceChainSelector, destinationChainSelector uint64, poolType cldf.ContractType, version semver.Version) error {
	ratelimiterConfig := token_pool.RateLimiterConfig{
		IsEnabled: true,
		Capacity:  new(big.Int).Mul(big.NewInt(1e16), big.NewInt(2)),
		Rate:      big.NewInt(1),
	}
	tokenPoolConfig := map[uint64]v1_5_1.TokenPoolConfig{
		sourceChainSelector: {
			Type:    poolType,
			Version: version,
			ChainUpdates: v1_5_1.RateLimiterPerChain{
				destinationChainSelector: v1_5_1.RateLimiterConfig{
					Inbound:  ratelimiterConfig,
					Outbound: ratelimiterConfig,
				},
			},
		},
		destinationChainSelector: {
			Type:    poolType,
			Version: version,
			ChainUpdates: v1_5_1.RateLimiterPerChain{
				sourceChainSelector: v1_5_1.RateLimiterConfig{
					Inbound:  ratelimiterConfig,
					Outbound: ratelimiterConfig,
				},
			},
		},
	}
	_, err := v1_5_1.ConfigureTokenPoolContractsChangeset(e, v1_5_1.ConfigureTokenPoolContractsConfig{
		TokenSymbol: shared.TokenSymbol(tokenSymbol),
		PoolUpdates: tokenPoolConfig,
	})
	return err
}

func configureTokenAdminRegistry(e cldf.Environment, tokenSymbol string, sourceChainSelector, destinationChainSelector uint64, poolType cldf.ContractType, version semver.Version) error {
	registryConfig := map[uint64]map[shared.TokenSymbol]v1_5_1.TokenPoolInfo{
		sourceChainSelector: {
			shared.TokenSymbol(tokenSymbol): {
				Type:          poolType,
				Version:       version,
				ExternalAdmin: e.BlockChains.EVMChains()[sourceChainSelector].DeployerKey.From,
			},
		},
		destinationChainSelector: {
			shared.TokenSymbol(tokenSymbol): {
				Type:          poolType,
				Version:       version,
				ExternalAdmin: e.BlockChains.EVMChains()[destinationChainSelector].DeployerKey.From,
			},
		},
	}

	_, err := v1_5_1.ProposeAdminRoleChangeset(e, v1_5_1.TokenAdminRegistryChangesetConfig{
		Pools:                   registryConfig,
		SkipOwnershipValidation: true,
	})
	if err != nil {
		return err
	}

	_, err = v1_5_1.AcceptAdminRoleChangeset(e, v1_5_1.TokenAdminRegistryChangesetConfig{
		Pools:                   registryConfig,
		SkipOwnershipValidation: true,
	})
	if err != nil {
		return err
	}

	_, err = v1_5_1.SetPoolChangeset(e, v1_5_1.TokenAdminRegistryChangesetConfig{
		Pools:                   registryConfig,
		SkipOwnershipValidation: true,
	})
	return err
}

func getFirstAddressFromChain(t *testing.T, addressBook cldf.AddressBook, chainSelector uint64) common.Address {
	addresses, err := addressBook.AddressesForChain(chainSelector)
	require.NoError(t, err)

	for addr := range addresses {
		return common.HexToAddress(addr)
	}

	require.Failf(t, "No addresses found for chain", "ChainSelector: %d", chainSelector)
	return common.Address{}
}

func configureFastTransferSettingsWithMCMS(t *testing.T, e cldf.Environment, tokenSymbol string, sourceChainSelector, destinationChainSelector uint64, fillerAddress common.Address, tc *fastTransferE2ETestCase, poolType cldf.ContractType, version semver.Version, useMCMS bool) error {
	fillers := []common.Address{}
	if tc.allowlistEnabled && tc.allowlistFiller {
		fillers = append(fillers, fillerAddress)
	}

	// Configure filler allowlist
	if tc.allowlistFiller {
		config := v1_5_1.FastTransferFillerAllowlistConfig{
			TokenSymbol:     shared.TokenSymbol(tokenSymbol),
			ContractType:    poolType,
			ContractVersion: version,
			Updates: map[uint64]v1_5_1.FillerAllowlistConfig{
				sourceChainSelector: {
					AddFillers:    fillers,
					RemoveFillers: []common.Address{},
				},
				destinationChainSelector: {
					AddFillers:    fillers,
					RemoveFillers: []common.Address{},
				},
			},
		}

		// Add MCMS configuration if requested
		if useMCMS {
			config.MCMS = &proposalutils.TimelockConfig{
				MinDelay:   0 * time.Second,
				MCMSAction: mcmstypes.TimelockActionSchedule,
			}
		}

		_, _, err := commonchangeset.ApplyChangesets(t, e,
			[]commonchangeset.ConfiguredChangeSet{commonchangeset.Configure(
				v1_5_1.FastTransferFillerAllowlistChangeset,
				config,
			)}, commonchangeset.WithRealBackend())
		if err != nil {
			return err
		}
	}

	// Configure lane settings
	settlementGasOverhead := tc.settlementGasOverhead
	laneConfig := v1_5_1.FastTransferUpdateLaneConfigConfig{
		TokenSymbol:     shared.TokenSymbol(tokenSymbol),
		ContractType:    poolType,
		ContractVersion: version,
		Updates: map[uint64](map[uint64]v1_5_1.UpdateLaneConfig){
			sourceChainSelector: {
				destinationChainSelector: {
					FastTransferFillerFeeBps: 10,
					FastTransferPoolFeeBps:   tc.fastTransferPoolFeeBps,
					FillerAllowlistEnabled:   tc.allowlistEnabled,
					FillAmountMaxRequest:     big.NewInt(100000),
					SettlementOverheadGas:    &settlementGasOverhead,
					SkipAllowlistValidation:  true,
				},
			},
			destinationChainSelector: {
				sourceChainSelector: {
					FastTransferFillerFeeBps: 20,
					FastTransferPoolFeeBps:   tc.fastTransferPoolFeeBps,
					FillerAllowlistEnabled:   tc.allowlistEnabled,
					FillAmountMaxRequest:     big.NewInt(100000),
					SettlementOverheadGas:    &settlementGasOverhead,
					SkipAllowlistValidation:  true,
				},
			},
		},
	}

	// Add MCMS configuration if requested
	if useMCMS {
		laneConfig.MCMS = &proposalutils.TimelockConfig{
			MinDelay:   0 * time.Second,
			MCMSAction: mcmstypes.TimelockActionSchedule,
		}
	}

	_, _, err := commonchangeset.ApplyChangesets(t, e,
		[]commonchangeset.ConfiguredChangeSet{commonchangeset.Configure(
			v1_5_1.FastTransferUpdateLaneConfigChangeset,
			laneConfig,
		)}, commonchangeset.WithRealBackend())
	return err
}

func transferTokenPoolOwnershipToMCMS(t *testing.T, e cldf.Environment, poolAddresses map[uint64][]common.Address) {
	_, _, err := commonchangeset.ApplyChangesets(t, e,
		[]commonchangeset.ConfiguredChangeSet{commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(commonchangeset.TransferToMCMSWithTimelockV2),
			commonchangeset.TransferToMCMSWithTimelockConfig{
				ContractsByChain: poolAddresses,
				MCMSConfig: proposalutils.TimelockConfig{
					MinDelay: 0 * time.Second, // No delay for tests
				},
			},
		)}, commonchangeset.WithRealBackend(),
	)
	require.NoError(t, err)

	// Renounce timelock deployer for the chains
	for chainSelector := range poolAddresses {
		_, _, err := commonchangeset.ApplyChangesets(t, e,
			[]commonchangeset.ConfiguredChangeSet{commonchangeset.Configure(
				cldf.CreateLegacyChangeSet(commonchangeset.RenounceTimelockDeployer),
				commonchangeset.RenounceTimelockDeployerConfig{
					ChainSel: chainSelector,
				},
			)}, commonchangeset.WithRealBackend(),
		)
		require.NoError(t, err)
	}
}

func configureTokenPoolContractsWithMCMS(t *testing.T, e cldf.Environment, tokenSymbol string, sourceChainSelector, destinationChainSelector uint64, sourceTokenAddress, destinationTokenAddress common.Address, tokenDecimals uint8, fillerAddress common.Address, tc *fastTransferE2ETestCase, sourceLock *sync.Mutex, destinationLock *sync.Mutex, useMCMS bool) (sourcePoolAddr common.Address, destPoolAddr common.Address, version semver.Version, poolWrapper *bindings.FastTransferTokenPoolWrapper, sourceMinter mintableToken, destMinter mintableToken) {
	sourceLock.Lock()
	defer sourceLock.Unlock()
	destinationLock.Lock()
	defer destinationLock.Unlock()

	var config tokenPoolConfig
	if tc.externalMinter {
		config = configureExternalMinterTokenPool(t, e, sourceChainSelector, destinationChainSelector, sourceTokenAddress, destinationTokenAddress, tokenDecimals)
	} else if tc.isHybridPool {
		config = configureHybridTokenPool(t, e, sourceChainSelector, destinationChainSelector, sourceTokenAddress, destinationTokenAddress, tokenDecimals, shared.TokenSymbol(tokenSymbol), useMCMS, tc.hybridPool)
	} else {
		config = configureBurnMintTokenPool(t, e, sourceChainSelector, destinationChainSelector, sourceTokenAddress, destinationTokenAddress, tokenDecimals)
	}

	// Step 1: Deploy token pools without MCMS
	cs, err := v1_5_1.DeployTokenPoolContractsChangeset(e, v1_5_1.DeployTokenPoolContractsConfig{
		TokenSymbol: shared.TokenSymbol(tokenSymbol),
		NewPools:    config.poolConfig,
	})
	require.NoError(t, err)

	sourceTokenPoolAddress := getFirstAddressFromChain(t, cs.AddressBook, sourceChainSelector)           //nolint:staticcheck // AddressBook is deprecated but still required
	destinationTokenPoolAddress := getFirstAddressFromChain(t, cs.AddressBook, destinationChainSelector) //nolint:staticcheck // AddressBook is deprecated but still required

	err = e.ExistingAddresses.Merge(cs.AddressBook) //nolint:staticcheck // AddressBook is deprecated but still required
	require.NoError(t, err)

	// Step 2: Configure basic token pool settings without MCMS (rate limits, admin registry)
	err = configureTokenPoolRateLimits(e, tokenSymbol, sourceChainSelector, destinationChainSelector, config.poolType, config.version)
	require.NoError(t, err)

	err = configureTokenAdminRegistry(e, tokenSymbol, sourceChainSelector, destinationChainSelector, config.poolType, config.version)
	require.NoError(t, err)

	// Step 3: Transfer ownership to MCMS if requested
	if useMCMS {
		poolAddresses := map[uint64][]common.Address{
			sourceChainSelector:      {sourceTokenPoolAddress},
			destinationChainSelector: {destinationTokenPoolAddress},
		}
		transferTokenPoolOwnershipToMCMS(t, e, poolAddresses)
	}

	// Step 4: Configure fast transfer settings (with or without MCMS)
	err = configureFastTransferSettingsWithMCMS(t, e, tokenSymbol, sourceChainSelector, destinationChainSelector, fillerAddress, tc, config.poolType, config.version, useMCMS)
	require.NoError(t, err)

	sourceTokenPool, err := bindings.GetFastTransferTokenPoolContract(e, shared.TokenSymbol(tokenSymbol), config.poolType, config.version, sourceChainSelector)
	require.NoError(t, err)

	config.postSetupAction(sourceTokenPoolAddress, destinationTokenPoolAddress)

	sourcePoolAddr = sourceTokenPoolAddress
	destPoolAddr = destinationTokenPoolAddress
	version = config.version
	poolWrapper = sourceTokenPool
	sourceMinter = config.sourceMinter
	destMinter = config.destinationMinter
	return
}

func getFillerImage() (string, error) {
	envVersion := os.Getenv(devenv.E2eFastFillerVersion)
	envImage := os.Getenv(devenv.E2eFastFillerImage)

	if envVersion == "" || envImage == "" {
		return devenv.DefaultFastFillerImage, nil
	}

	return envImage + ":" + envVersion, nil
}

func runAssertions(t *testing.T, sourceToken balanceToken, destinationToken balanceToken, address common.Address, assertions []balanceAssertion, description string) {
	for _, assertion := range assertions {
		assertion(t, sourceToken, destinationToken, address, description)
	}
}

func startRelayer(t *testing.T, sourceChainSelector, destinationChainSelector uint64, sourceTokenPoolAddress common.Address, destinationTokenPoolAddress common.Address, deployedEnv testhelpers.TestEnvironment, fillerPrivateKey *ecdsa.PrivateKey) func() error {
	dockerEnv, ok := deployedEnv.(*testsetups.DeployedLocalDevEnvironment)
	require.True(t, ok, "deployedEnv is not of type *testsetups.DeployedLocalDevEnvironment")

	networks := dockerEnv.GetCLClusterTestEnv().EVMNetworks
	var sourceChainNetwork *blockchain.EVMNetwork
	for _, network := range networks {
		if network.ChainID >= 0 && uint64(network.ChainID) == sourceChainID {
			sourceChainNetwork = network
			break
		}
	}
	require.NotNil(t, sourceChainNetwork, "Source chain network not found in EVM networks")

	var destinationChainNetwork *blockchain.EVMNetwork
	for _, network := range networks {
		if network.ChainID >= 0 && uint64(network.ChainID) == destinationChainID {
			destinationChainNetwork = network
			break
		}
	}
	require.NotNil(t, destinationChainNetwork, "Destination chain network not found in EVM networks")

	marshalledKey := crypto.FromECDSA(fillerPrivateKey)

	hexString := hex.EncodeToString(marshalledKey)

	fastFillerConfig := devenv.CCIPFastFillerConfig{
		SignerProviders: []devenv.SignerProvider{
			{
				Name:       "filler",
				Type:       "raw",
				PrivateKey: hexString,
			},
		},
		Listeners: []devenv.ListenerConfig{
			{
				ChainSelector:        strconv.FormatUint(sourceChainSelector, 10),
				TokenPoolAddress:     sourceTokenPoolAddress.Hex(),
				RPCURL:               sourceChainNetwork.HTTPURLs[0],
				DestinationTokenPool: destinationTokenPoolAddress.Hex(),
			},
		},
		Fillers: []devenv.FillerConfig{
			{
				ChainSelector:    strconv.FormatUint(destinationChainSelector, 10),
				TokenPoolAddress: destinationTokenPoolAddress.Hex(),
				RPCURL:           destinationChainNetwork.HTTPURLs[0],
				SignerProvider:   "filler",
				SourceTokenPool:  sourceTokenPoolAddress.Hex(),
			},
		},
	}
	image, err := getFillerImage()
	require.NoError(t, err, "Failed to get filler image")
	l := logging.GetTestLogger(t)
	relayer := devenv.NewCCIPFastFiller(fastFillerConfig, l, []string{dockerEnv.GetCLClusterTestEnv().DockerNetwork.ID}, image)
	err = relayer.Start(t.Context(), t)
	require.NoError(t, err, "Failed to start the relayer")

	return func() error { return relayer.Stop(context.Background()) }
}

func setupFastTransfer1_5TestEnvironment(t *testing.T, useMCMS bool) *fastTransferTestContext {
	e, _, tEnv := testsetups.NewIntegrationEnvironment(
		t,
		testhelpers.WithPrerequisiteDeploymentOnly(
			&changeset.V1_5DeploymentConfig{
				PriceRegStalenessThreshold: 60 * 60 * 24 * 14, // two weeks
				RMNConfig: &rmn_contract.RMNConfig{
					BlessWeightThreshold: 2,
					CurseWeightThreshold: 2,
					// setting dummy voters, we will permabless this later
					Voters: []rmn_contract.RMNVoter{
						{
							BlessWeight:   2,
							CurseWeight:   2,
							BlessVoteAddr: utils.RandomAddress(),
							CurseVoteAddr: utils.RandomAddress(),
						},
					},
				},
			}),
	)
	state, err := stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)
	allChains := e.Env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chainselectors.FamilyEVM))
	src1, dest := allChains[0], allChains[1]
	pairs := []testhelpers.SourceDestPair{
		{SourceChainSelector: src1, DestChainSelector: dest},
		{SourceChainSelector: dest, DestChainSelector: src1},
	}
	// wire up all lanes
	// deploy onRamp, commit store, offramp , set ocr2config and send corresponding jobs
	e.Env = v1_5testhelpers.AddLanes(t, e.Env, state, pairs)

	// permabless the commit stores
	e.Env, err = commonchangeset.Apply(t, e.Env,
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(v1_5.PermaBlessCommitStoreChangeset),
			v1_5.PermaBlessCommitStoreConfig{
				Configs: map[uint64]v1_5.PermaBlessCommitStoreConfigPerDest{
					dest: {
						Sources: []v1_5.PermaBlessConfigPerSourceChain{
							{
								SourceChainSelector: src1,
								PermaBless:          true,
							},
						},
					},
					src1: {
						Sources: []v1_5.PermaBlessConfigPerSourceChain{
							{
								SourceChainSelector: dest,
								PermaBless:          true,
							},
						},
					},
				},
			},
		),
	)
	require.NoError(t, err)

	onChainState, err := stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)

	sourceChainSelector := e.Env.BlockChains.ListChainSelectors()[0]
	destinationChainSelector := e.Env.BlockChains.ListChainSelectors()[1]

	sourceChainState := onChainState.Chains[sourceChainSelector]
	destinationChain := e.Env.BlockChains.EVMChains()[destinationChainSelector]

	require.NoError(t, err)

	seqNumRetriever := func(opts *bind.CallOpts, destChainSelector uint64) (uint64, error) {
		onramp := onChainState.Chains[sourceChainSelector].EVM2EVMOnRamp[destChainSelector]
		seq, err := onramp.GetExpectedNextSequenceNumber(opts)
		if err != nil {
			return 0, fmt.Errorf("failed to get expected next sequence number: %w", err)
		}
		return seq, nil
	}

	offramp := onChainState.Chains[destinationChainSelector].EVM2EVMOffRamp[sourceChainSelector]
	commitStore := onChainState.Chains[destinationChainSelector].CommitStore[sourceChainSelector]

	waitForExecution := func(t *testing.T, sequenceNumber uint64) {
		sourceChain := e.Env.BlockChains.EVMChains()[sourceChainSelector]
		v1_5testhelpers.WaitForCommit(t, sourceChain, destinationChain, commitStore, sequenceNumber)
		e.Env.Logger.Infof("Commit confirmed, waiting for offramp execution for sequence number %d", sequenceNumber)
		v1_5testhelpers.WaitForExecute(t, sourceChain, destinationChain, offramp, []uint64{sequenceNumber}, uint64(0))
	}

	waitForExecutionError := func(t *testing.T, sequenceNumber uint64) {
		sourceChain := e.Env.BlockChains.EVMChains()[sourceChainSelector]
		v1_5testhelpers.WaitForNoExec(t, sourceChain, destinationChain, offramp, sequenceNumber)
	}

	// Create shared locks for coordination between parallel test cases
	sourceLock := &sync.Mutex{}
	destinationLock := &sync.Mutex{}
	sendLock := &sync.Mutex{}

	return newFastTransferTestContext(
		e.Env,
		sourceChainState,
		tEnv,
		seqNumRetriever,
		waitForExecution,
		waitForExecutionError,
		sourceLock,
		destinationLock,
		sendLock,
		useMCMS,
	)
}

func TestFastTransfer1_5Lanes(t *testing.T) {
	baseCtx := setupFastTransfer1_5TestEnvironment(t, false)

	for i, tc := range fastTransferTestCases {
		ctx := baseCtx.WithTestIndex(i)
		runFastTransferTestCase(t, ctx, tc)
	}
}

func setupFastTransfer1_6TestEnvironment(t *testing.T, useMCMS bool) *fastTransferTestContext {
	e, _, deployedEnv := testsetups.NewIntegrationEnvironment(t)

	onChainState, err := stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)
	testhelpers.AddLanesForAll(t, &e, onChainState)

	sourceChainSelector := e.Env.BlockChains.ListChainSelectors()[0]
	destinationChainSelector := e.Env.BlockChains.ListChainSelectors()[1]

	sourceChainState := onChainState.Chains[sourceChainSelector]
	destinationChain := e.Env.BlockChains.EVMChains()[destinationChainSelector]

	seqNumRetriever := func(opts *bind.CallOpts, destChainSelector uint64) (uint64, error) {
		onramp := onChainState.Chains[sourceChainSelector].OnRamp
		seq, err := onramp.GetExpectedNextSequenceNumber(opts, destChainSelector)
		if err != nil {
			return 0, fmt.Errorf("failed to get expected next sequence number: %w", err)
		}
		return seq, nil
	}

	offramp := onChainState.Chains[destinationChainSelector].OffRamp
	waitForExecution := func(t *testing.T, sequenceNumber uint64) {
		zero := uint64(0)
		_, _ = testhelpers.ConfirmExecWithSeqNrs(t, sourceChainSelector, destinationChain, offramp, &zero, []uint64{sequenceNumber})
	}

	waitForExecutionError := func(t *testing.T, sequenceNumber uint64) {
		testhelpers.ConfirmNoExecSuccessConsistentlyWithSeqNr(t, sourceChainSelector, destinationChain, offramp, sequenceNumber, 30*time.Second)
	}

	// Create shared locks for coordination between parallel test cases
	sourceLock := &sync.Mutex{}
	destinationLock := &sync.Mutex{}
	sendLock := &sync.Mutex{}

	return newFastTransferTestContext(
		e.Env,
		sourceChainState,
		deployedEnv,
		seqNumRetriever,
		waitForExecution,
		waitForExecutionError,
		sourceLock,
		destinationLock,
		sendLock,
		useMCMS,
	)
}

func TestFastTransfer1_6Lanes(t *testing.T) {
	baseCtx := setupFastTransfer1_6TestEnvironment(t, false)

	for i, tc := range fastTransferTestCases {
		ctx := baseCtx.WithTestIndex(i)
		runFastTransferTestCase(t, ctx, tc)
	}
}

func TestFastTransfer1_5LanesWithMCMS(t *testing.T) {
	baseCtx := setupFastTransfer1_5TestEnvironment(t, true)

	for i, tc := range fastTransferTestCases {
		ctx := baseCtx.WithTestIndex(i)
		runFastTransferTestCase(t, ctx, tc)
	}
}

func TestFastTransfer1_6LanesWithMCMS(t *testing.T) {
	baseCtx := setupFastTransfer1_6TestEnvironment(t, true)

	for i, tc := range fastTransferTestCases {
		ctx := baseCtx.WithTestIndex(i)
		runFastTransferTestCase(t, ctx, tc)
	}
}

type fastTransferTestContext struct {
	env                     cldf.Environment
	testIndex               int
	sourceLock              *sync.Mutex
	destinationLock         *sync.Mutex
	sendLock                *sync.Mutex
	sourceChainState        evm.CCIPChainState
	deployedEnv             testhelpers.TestEnvironment
	sequenceNumberRetriever sequenceNumberRetriever
	waitForExecution        waitForExecutionFn
	waitForExecutionError   waitForExecutionFn
	useMCMS                 bool
}

func (ctx *fastTransferTestContext) SourceChainSelector() uint64 {
	return ctx.env.BlockChains.ListChainSelectors()[0]
}

func (ctx *fastTransferTestContext) DestinationChainSelector() uint64 {
	return ctx.env.BlockChains.ListChainSelectors()[1]
}

func (ctx *fastTransferTestContext) SourceChain() evmChain.Chain {
	return ctx.env.BlockChains.EVMChains()[ctx.SourceChainSelector()]
}

func (ctx *fastTransferTestContext) DestinationChain() evmChain.Chain {
	return ctx.env.BlockChains.EVMChains()[ctx.DestinationChainSelector()]
}

func (ctx *fastTransferTestContext) SourceLock() *sync.Mutex {
	return ctx.sourceLock
}

func (ctx *fastTransferTestContext) DestinationLock() *sync.Mutex {
	return ctx.destinationLock
}

func (ctx *fastTransferTestContext) SendLock() *sync.Mutex {
	return ctx.sendLock
}

func (ctx *fastTransferTestContext) WithTestIndex(testIndex int) *fastTransferTestContext {
	clone := *ctx
	clone.testIndex = testIndex
	return &clone
}

func newFastTransferTestContext(
	env cldf.Environment,
	sourceChainState evm.CCIPChainState,
	deployedEnv testhelpers.TestEnvironment,
	sequenceNumberRetriever sequenceNumberRetriever,
	waitForExecution waitForExecutionFn,
	waitForExecutionError waitForExecutionFn,
	sourceLock *sync.Mutex,
	destinationLock *sync.Mutex,
	sendLock *sync.Mutex,
	useMCMS bool,
) *fastTransferTestContext {
	return &fastTransferTestContext{
		env:                     env,
		testIndex:               0,
		sourceLock:              sourceLock,
		destinationLock:         destinationLock,
		sendLock:                sendLock,
		sourceChainState:        sourceChainState,
		deployedEnv:             deployedEnv,
		sequenceNumberRetriever: sequenceNumberRetriever,
		waitForExecution:        waitForExecution,
		waitForExecutionError:   waitForExecutionError,
		useMCMS:                 useMCMS,
	}
}

func runFastTransferTestCase(t *testing.T, ctx *fastTransferTestContext, tc *fastTransferE2ETestCase) {
	tc.tokenSymbol = fmt.Sprintf("FTF_TEST_%d", ctx.testIndex+1)
	t.Run(tc.name, func(t *testing.T) {
		t.Parallel()
		userAddress, userTransactor, _ := createAccount(t, sourceChainID)
		fillerAddress, fillerTransactor, fillerPrivateKey := createAccount(t, destinationChainID)
		sourceToken := deployTokenAndGrantAllRoles(t, ctx.SourceChain(), tc.tokenSymbol, tokenDecimals, ctx.sourceLock, tc.externalMinter || tc.isHybridPool)
		destinationToken := deployTokenAndGrantAllRoles(t, ctx.DestinationChain(), tc.tokenSymbol, tokenDecimals, ctx.destinationLock, tc.externalMinter || tc.isHybridPool)

		sourceTokenPoolAddress, destinationTokenPoolAddress, contractVersion, _, sourceMinter, destinationMinter := configureTokenPoolContractsWithMCMS(t, ctx.env, tc.tokenSymbol, ctx.SourceChainSelector(), ctx.DestinationChainSelector(), sourceToken.Address(), destinationToken.Address(), tokenDecimals, fillerAddress, tc, ctx.sourceLock, ctx.destinationLock, ctx.useMCMS)
		var contractType cldf.ContractType
		switch {
		case tc.isHybridPool:
			contractType = shared.HybridWithExternalMinterFastTransferTokenPool
		case tc.externalMinter:
			contractType = shared.BurnMintWithExternalMinterFastTransferTokenPool
		default:
			contractType = shared.BurnMintFastTransferTokenPool
		}
		pool, err := bindings.NewFastTransferTokenPoolWrapper(sourceTokenPoolAddress, ctx.SourceChain().Client, contractType)
		require.NoError(t, err)

		onChainState, err := stateview.LoadOnchainState(ctx.env)
		require.NoError(t, err)

		userEncodedAddress := common.LeftPadBytes(userAddress.Bytes(), 32)

		var feeTokenAddress common.Address
		switch tc.feeTokenType {
		case feeTokenLink:
			feeTokenAddress = onChainState.Chains[ctx.SourceChainSelector()].LinkToken.Address()
		case feeTokenNative:
			feeTokenAddress = common.HexToAddress("0x0")
		default:
			t.Fatalf("Unknown fee token type: %s", tc.feeTokenType)
		}

		fees, err := pool.GetCcipSendTokenFee(nil, ctx.DestinationChainSelector(), transferAmount, userEncodedAddress, feeTokenAddress, []byte{})
		require.NoError(t, err)

		// Setup source chain funding and approvals
		fundAccount(t, ctx.SourceChain(), userAddress, defaultEthAmount, ctx.sourceLock)
		fundAccountWithToken(t, ctx.SourceChain(), userAddress, sourceMinter, initialUserTokenAmountOnSource, ctx.sourceLock)
		approveToken(t, ctx.SourceChain(), userTransactor(), sourceToken, sourceTokenPoolAddress, ctx.sourceLock)

		if tc.feeTokenType == feeTokenLink {
			sourceLinkToken := getLinkTokenAndGrantMintRole(t, ctx.SourceChain(), ctx.sourceChainState, ctx.sourceLock)
			fundAccountWithToken(t, ctx.SourceChain(), userAddress, sourceLinkToken, fees.CcipSettlementFee, ctx.sourceLock)
			approveToken(t, ctx.SourceChain(), userTransactor(), sourceLinkToken, sourceTokenPoolAddress, ctx.sourceLock)
		}

		// Setup destination chain funding and approvals
		fundAccount(t, ctx.DestinationChain(), fillerAddress, defaultEthAmount, ctx.destinationLock)
		fundAccountWithToken(t, ctx.DestinationChain(), fillerAddress, destinationMinter, initialFillerTokenAmountOnDest, ctx.destinationLock)
		approveToken(t, ctx.DestinationChain(), fillerTransactor(), destinationToken, destinationTokenPoolAddress, ctx.destinationLock)

		if tc.enableFiller {
			stop := startRelayer(t, ctx.SourceChainSelector(), ctx.DestinationChainSelector(), sourceTokenPoolAddress, destinationTokenPoolAddress, ctx.deployedEnv, fillerPrivateKey)
			ctx.env.Logger.Infof("Started relayer for source chain %d and destination chain %d", ctx.SourceChainSelector(), ctx.DestinationChainSelector())

			defer func() {
				ctx.env.Logger.Infof("Stopping relayer for source chain %d and destination chain %d", ctx.SourceChainSelector(), ctx.DestinationChainSelector())
				_ = stop()
			}()
		}

		runAssertions(t, sourceToken, destinationToken, fillerAddress, tc.preFastTransferFillerAssertions, "Pre Fast Transfer Filler Assertions")
		runAssertions(t, sourceToken, destinationToken, userAddress, tc.preFastTransferUserAssertions, "Pre Fast Transfer User Assertions")
		runAssertions(t, sourceToken, destinationToken, destinationTokenPoolAddress, tc.preFastTransferPoolAssertions, "Pre Fast Transfer Pool Assertions")

		userTransac := userTransactor()
		if tc.feeTokenType == feeTokenNative {
			userTransac.Value = fees.CcipSettlementFee
		}

		var seqNum uint64
		func() {
			ctx.sendLock.Lock()
			defer ctx.sendLock.Unlock()
			require.NoError(t, err)
			seqNum, err = ctx.sequenceNumberRetriever(nil, ctx.DestinationChainSelector())
			require.NoError(t, err)
			ctx.env.Logger.Infof("Sending transaction from user address: %s", userTransac.From.Hex())

			// Determine max fast transfer fee - use custom fee if set, otherwise use calculated fee
			maxFastTransferFee := fees.FastTransferFee
			if tc.customMaxFastTransferFee != nil {
				maxFastTransferFee = tc.customMaxFastTransferFee
			}

			tx, err := pool.CcipSendToken(userTransac, ctx.DestinationChainSelector(), transferAmount, maxFastTransferFee, userEncodedAddress, feeTokenAddress, []byte{})

			if tc.expectRevert {
				// Expect the transaction to fail
				require.Error(t, err, "Expected ccipSendToken to revert when maxFastTransferFee is too low")
				ctx.env.Logger.Infof("Transaction correctly reverted as expected: %v", err)
				return
			}

			ctx.env.Logger.Infof("Sending transaction: %s", tx.Hash().Hex())
			require.NoError(t, err)
			_, err = ctx.SourceChain().Confirm(tx)
			require.NoError(t, err)

			filter, err := pool.FilterFastTransferRequested(nil, nil, nil, nil)
			require.NoError(t, err)
			for filter.Next() {
				event := filter.Event()
				ctx.env.Logger.Infof("FastTransferRequested event: %s, fillId: %s, settlementId: %s", event.Raw.TxHash.Hex(), hex.EncodeToString(event.FillID[:]), hex.EncodeToString(event.SettlementID[:]))
			}
		}()

		// Skip post-transaction logic if we expect the transaction to revert
		if !tc.expectRevert {
			runAssertions(t, sourceToken, destinationToken, fillerAddress, tc.postFastTransferFillerAssertions, "Post Fast Transfer Filler Assertions")
			runAssertions(t, sourceToken, destinationToken, userAddress, tc.postFastTransferUserAssertions, "Post Fast Transfer User Assertions")
			runAssertions(t, sourceToken, destinationToken, destinationTokenPoolAddress, tc.postFastTransferPoolAssertions, "Post Fast Transfer Pool Assertions")

			if tc.expectNoExecutionError {
				ctx.waitForExecutionError(t, seqNum)
			} else {
				ctx.waitForExecution(t, seqNum)
			}

			runAssertions(t, sourceToken, destinationToken, fillerAddress, tc.postRegularTransferFillerAssertions, "Post Regular Transfer Filler Assertions")
			runAssertions(t, sourceToken, destinationToken, userAddress, tc.postRegularTransferUserAssertions, "Post Regular Transfer User Assertions")
			runAssertions(t, sourceToken, destinationToken, destinationTokenPoolAddress, tc.postRegularTransferPoolAssertions, "Post Regular Transfer Pool Assertions")

			if tc.enableFiller {
				expectedPoolFee := big.NewInt(0).Mul(transferAmount, big.NewInt(int64(tc.fastTransferPoolFeeBps)))
				expectedPoolFee = big.NewInt(0).Div(expectedPoolFee, big.NewInt(10000))
				// Pool fees are only accumulated during fast fills
				poolFeeAssertion := assertPoolFeeWithdrawal(expectedPoolFee, ctx.env, tc.tokenSymbol, contractType, contractVersion, ctx.DestinationChainSelector(), destinationToken, ctx.useMCMS, ctx.destinationLock)
				poolFeeAssertion(t, destinationToken, destinationToken, destinationTokenPoolAddress, "Pool Fee Withdrawal Test")
			} else {
				// When no filler was used, pool fees should be 0
				poolFeeAssertion := assertPoolFeeWithdrawal(big.NewInt(0), ctx.env, tc.tokenSymbol, contractType, contractVersion, ctx.DestinationChainSelector(), destinationToken, ctx.useMCMS, ctx.destinationLock)
				poolFeeAssertion(t, destinationToken, destinationToken, destinationTokenPoolAddress, "Pool Fee Withdrawal Test (No Filler)")
			}
		}

		if !tc.expectNoExecutionError && !tc.expectRevert {
			ctx.env.Logger.Info("Sanity check regular token transfer (slow-path)")

			// Apply transfer fee config updates only to hybrid pools on 1.6 lanes
			if tc.isHybridPool && contractVersion.Compare(&deployment.Version1_6_0) >= 0 {
				state, _ := onChainState.EVMChainState(ctx.SourceChainSelector())
				if state.FeeQuoter != nil {
					configs := []fee_quoter.FeeQuoterTokenTransferFeeConfigArgs{
						{
							DestChainSelector: ctx.DestinationChainSelector(),
							TokenTransferFeeConfigs: []fee_quoter.FeeQuoterTokenTransferFeeConfigSingleTokenArgs{
								{
									Token: sourceToken.Address(),
									TokenTransferFeeConfig: fee_quoter.FeeQuoterTokenTransferFeeConfig{
										MinFeeUSDCents:    150,
										MaxFeeUSDCents:    4294967295,
										DeciBps:           0,
										DestGasOverhead:   200_000,
										DestBytesOverhead: 640,
										IsEnabled:         true,
									},
								},
							},
						},
					}
					tx, err := state.FeeQuoter.ApplyTokenTransferFeeConfigUpdates(ctx.SourceChain().DeployerKey, configs, nil)
					require.NoError(t, err, "Failed to apply token transfer fee config updates")
					ctx.env.Logger.Infof("Applied token transfer fee config updates transaction: %s", tx.Hash().Hex())
					_, err = ctx.SourceChain().Confirm(tx)
					require.NoError(t, err, "Failed to confirm token transfer fee config updates transaction")
				} else {
					ctx.env.Logger.Infof("FeeQuoter not available on chain %d, skipping token transfer fee config updates", ctx.SourceChainSelector())
				}
			}
			// We want to ensure regular transfer works as expected
			message := router.ClientEVM2AnyMessage{
				Receiver: common.LeftPadBytes(userAddress.Bytes(), 32),
				Data:     []byte{},
				TokenAmounts: []router.ClientEVMTokenAmount{
					{
						Token:  sourceToken.Address(),
						Amount: initialUserTokenAmountOnSource,
					},
				},
				FeeToken:  common.HexToAddress("0x0"),
				ExtraArgs: nil,
			}
			userBalance, err := destinationToken.BalanceOf(nil, userAddress)
			require.NoError(t, err)
			// Top-up user account on source chain
			fundAccountWithToken(t, ctx.SourceChain(), userAddress, sourceMinter, initialUserTokenAmountOnSource, ctx.sourceLock)
			approveToken(t, ctx.SourceChain(), userTransactor(), sourceToken, ctx.sourceChainState.Router.Address(), ctx.sourceLock)
			func() {
				ctx.sendLock.Lock()
				defer ctx.sendLock.Unlock()
				seqNum, err = ctx.sequenceNumberRetriever(nil, ctx.DestinationChainSelector())
				require.NoError(t, err)
				router := onChainState.Chains[ctx.SourceChainSelector()].Router
				fee, err := router.GetFee(&bind.CallOpts{Context: context.Background()}, ctx.DestinationChainSelector(), message)
				require.NoError(t, err)
				userTransac := userTransactor()
				userTransac.Value = fee
				tx, err := router.CcipSend(userTransac, ctx.DestinationChainSelector(), message)
				require.NoError(t, err)
				ctx.env.Logger.Infof("Sending regular transfer transaction: %s", tx.Hash().Hex())
				_, err = ctx.SourceChain().Confirm(tx)
				require.NoError(t, err)
			}()

			ctx.waitForExecution(t, seqNum)
			finalBalance, err := destinationToken.BalanceOf(nil, userAddress)
			require.NoError(t, err)
			expectedBalance := new(big.Int).Add(userBalance, initialUserTokenAmountOnSource)
			require.Equal(t, expectedBalance.String(), finalBalance.String(), "Final balance after regular transfer does not match expected value")
		}
	})
}

func fundAccount(
	t *testing.T,
	chain evmChain.Chain,
	receiver common.Address,
	amount *big.Int,
	sendLock *sync.Mutex,
) {
	sendLock.Lock()
	defer sendLock.Unlock()
	client := chain.Client
	sender := chain.DeployerKey

	nonce, err := client.NonceAt(t.Context(), sender.From, nil)
	require.NoError(t, err)

	gasPrice, err := client.SuggestGasPrice(t.Context())
	require.NoError(t, err)
	gasLimit := uint64(21000)

	tx := types.NewTransaction(nonce, receiver, amount, gasLimit, gasPrice, nil)

	signedTx, err := sender.Signer(sender.From, tx)
	require.NoError(t, err)

	err = client.SendTransaction(t.Context(), signedTx)
	require.NoError(t, err)

	_, err = chain.Confirm(signedTx)
	require.NoError(t, err)
}
