package actions

import (
	"context"
	crypto_rand "crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"

	"dario.cat/mergo"
	"github.com/AlekSi/pointer"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/atomic"
	"golang.org/x/exp/rand"
	"golang.org/x/sync/errgroup"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/blockchain"
	ctfClient "github.com/smartcontractkit/chainlink-testing-framework/lib/client"
	ctftestenv "github.com/smartcontractkit/chainlink-testing-framework/lib/docker/test_env"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/environment"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/pkg/helm/foundry"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/pkg/helm/mockserver"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/pkg/helm/reorg"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/networks"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/ptr"
	"github.com/smartcontractkit/chainlink/integration-tests/ccip-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/ccip-tests/contracts/laneconfig"
	"github.com/smartcontractkit/chainlink/integration-tests/ccip-tests/testconfig"
	"github.com/smartcontractkit/chainlink/integration-tests/ccip-tests/testreporters"
	testutils "github.com/smartcontractkit/chainlink/integration-tests/ccip-tests/utils"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/docker/test_env"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/commit_store"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/evm_2_evm_offramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/evm_2_evm_onramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/evm_2_evm_onramp_1_2_0"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/fee_quoter"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/mock_rmn_contract"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/rmn_contract"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/token_pool"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/burn_mint_erc677"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/ccipexec"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/testhelpers"
	integrationtesthelpers "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/testhelpers/integration"
	bigmath "github.com/smartcontractkit/chainlink/v2/core/utils/big_math"
)

const (
	ChaosGroupExecution           = "ExecutionNodesAll"      // all execution nodes
	ChaosGroupCommit              = "CommitNodesAll"         // all commit nodes
	ChaosGroupCommitFaultyPlus    = "CommitMajority"         // >f number of nodes
	ChaosGroupCommitFaulty        = "CommitMinority"         //  f number of nodes
	ChaosGroupExecutionFaultyPlus = "ExecutionNodesMajority" // > f number of nodes
	ChaosGroupExecutionFaulty     = "ExecutionNodesMinority" //  f number of nodes

	ChaosGroupCommitAndExecFaulty     = "CommitAndExecutionNodesMinority" //  f number of nodes
	ChaosGroupCommitAndExecFaultyPlus = "CommitAndExecutionNodesMajority" // >f number of nodes
	ChaosGroupCCIPGeth                = "CCIPGeth"                        // both source and destination simulated geth networks
	ChaosGroupNetworkACCIPGeth        = "CCIPNetworkAGeth"
	ChaosGroupNetworkBCCIPGeth        = "CCIPNetworkBGeth"

	defaultUSDCDestBytesOverhead = 640
	defaultUSDCDestGasOverhead   = 150_000
	DefaultDestinationGasLimit   = 600_000
	// DefaultResubscriptionTimeout denotes the max backoff duration for resubscription for various watch events
	// if the subscription keeps failing even after this duration, the test will fail
	DefaultResubscriptionTimeout = 2 * time.Hour
	LBTCValidDestPoolData        = "9b11457aa29d65e4940b67b7da16bd370d29bf6a3247a28066f93ac407b8b811"
)

// TODO: These should be refactored along with the default CCIP test setup to use optional config functions
var (
	// DefaultPermissionlessExecThreshold denotes how long the DON will retry a transaction before giving up,
	// otherwise known as the "Smart Execution Time Window". If a transaction fails to execute within this time window,
	// the DON will give up and the transaction will need Manual Execution as detailed here: https://docs.chain.link/ccip/concepts/manual-execution#manual-execution
	// For performance tests: the higher the load/throughput, the higher value we might need here to guarantee that nonces are not blocked
	// 1 day should be enough for most of the cases
	DefaultPermissionlessExecThreshold        = time.Hour * 8
	DefaultMaxNoOfTokensInMsg          uint16 = 50
)

type CCIPTOMLEnv struct {
	Networks []blockchain.EVMNetwork
}

var (
	NetworkChart = reorg.TXNodesAppLabel
	NetworkName  = func(name string) string {
		return strings.ReplaceAll(strings.ToLower(name), " ", "-")
	}
	InflightExpiryExec   = 3 * time.Minute
	InflightExpiryCommit = 3 * time.Minute
	BatchGasLimit        = uint32(7_000_000)

	MaxDataBytes = uint32(50_000)

	RootSnoozeTime = 3 * time.Minute
	GethLabel      = func(name string) string {
		name = NetworkName(name)
		switch NetworkChart {
		case reorg.TXNodesAppLabel:
			return fmt.Sprintf("%s-ethereum-geth", name)
		case foundry.ChartName:
			return name
		}
		return ""
	}
	// ApprovedAmountToRouter is the default amount which gets approved for router so that it can transfer token and use the fee token for fee payment
	ApprovedAmountToRouter           = new(big.Int).Mul(big.NewInt(1e18), big.NewInt(1))
	ApprovedFeeAmountToRouter        = new(big.Int).Mul(big.NewInt(int64(GasFeeMultiplier)), big.NewInt(1e5))
	GasFeeMultiplier          uint64 = 12e17
	LinkToUSD                        = new(big.Int).Mul(big.NewInt(1e18), big.NewInt(20))
	WrappedNativeToUSD               = new(big.Int).Mul(big.NewInt(1e18), big.NewInt(1.7e3))
)

func GetUSDCDomain(networkName string, simulated bool) (uint32, error) {
	if simulated {
		// generate a random domain for simulated networks
		return rand.Uint32(), nil
	}
	lookup := map[string]uint32{
		networks.AvalancheFuji.Name:  1,
		networks.OptimismGoerli.Name: 2,
		networks.ArbitrumGoerli.Name: 3,
		networks.BaseGoerli.Name:     6,
		networks.PolygonMumbai.Name:  7,
	}
	if val, ok := lookup[networkName]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("USDC domain not found for chain %s", networkName)
}

type CCIPCommon struct {
	Logger      *zerolog.Logger
	ChainClient blockchain.EVMClient
	// Deployer deploys all CCIP contracts
	Deployer *contracts.CCIPContractsDeployer
	// tokenDeployer is used exclusively for deploying self-serve tokens and their pools
	tokenDeployer                 *contracts.CCIPContractsDeployer
	FeeToken                      *contracts.LinkToken
	BridgeTokens                  []*contracts.ERC20Token
	PriceAggregators              map[common.Address]*contracts.MockAggregator
	NoOfTokensNeedingDynamicPrice int
	BridgeTokenPools              []*contracts.TokenPool
	RateLimiterConfig             contracts.RateLimiterConfig
	RMNContract                   *common.Address
	ARM                           *contracts.ARM // populate only if the ARM contracts is not a mock and can be used to verify various ARM events; keep this nil for mock ARM
	Router                        *contracts.Router
	PriceRegistry                 *contracts.PriceRegistry
	TokenAdminRegistry            *contracts.TokenAdminRegistry
	WrappedNative                 common.Address
	MulticallEnabled              bool
	MulticallContract             common.Address
	ExistingDeployment            bool
	USDCMockDeployment            *bool
	LBTCMockDeployment            *bool
	LBTCDestPoolDataAs32Bytes     *bool
	TokenMessenger                *common.Address
	TokenTransmitter              *contracts.TokenTransmitter
	IsConnectionRestoredRecently  *atomic.Bool

	poolFunds                 *big.Int
	tokenPriceUpdateWatcherMu *sync.Mutex
	tokenPriceUpdateWatcher   map[common.Address]*big.Int // key - token; value - timestamp of update
	gasUpdateWatcherMu        *sync.Mutex
	gasUpdateWatcher          map[uint64]*big.Int // key - destchain id; value - timestamp of update
	GasUpdateEvents           []contracts.GasUpdateEvent
	AllowOutOfOrder           bool
}

// FreeUpUnusedSpace sets nil to various elements of ccipModule which are only used
// during lane set up and not used for rest of the test duration
// this is called mainly by load test to keep the memory usage minimum for high number of lanes
func (ccipModule *CCIPCommon) FreeUpUnusedSpace() {
	ccipModule.PriceAggregators = nil
	ccipModule.BridgeTokenPools = nil
	ccipModule.TokenMessenger = nil
	ccipModule.TokenTransmitter = nil
	runtime.GC()
}

func (ccipModule *CCIPCommon) UnvoteToCurseARM() error {
	if ccipModule.ARM != nil {
		return fmt.Errorf("real ARM deployed. cannot curse through test")
	}
	if ccipModule.RMNContract == nil {
		return fmt.Errorf("no ARM contract is set")
	}
	arm, err := mock_rmn_contract.NewMockRMNContract(*ccipModule.RMNContract, ccipModule.ChainClient.Backend())
	if err != nil {
		return fmt.Errorf("error instantiating arm %w", err)
	}
	opts, err := ccipModule.ChainClient.TransactionOpts(ccipModule.ChainClient.GetDefaultWallet())
	if err != nil {
		return fmt.Errorf("error getting owners for ARM OwnerUnvoteToCurse %w", err)
	}
	tx, err := arm.OwnerUnvoteToCurse0(opts, []mock_rmn_contract.RMNUnvoteToCurseRecord{})
	if err != nil {
		return fmt.Errorf("error in calling OwnerUnvoteToCurse %w", err)
	}
	err = ccipModule.ChainClient.ProcessTransaction(tx)
	if err != nil {
		return err
	}
	log.Info().
		Str("ARM", arm.Address().Hex()).
		Msg("ARM is uncursed")
	return ccipModule.ChainClient.WaitForEvents()
}

func (ccipModule *CCIPCommon) IsCursed() (bool, error) {
	if ccipModule.ARM != nil {
		return false, fmt.Errorf("real ARM deployed. cannot validate cursing")
	}
	if ccipModule.RMNContract == nil {
		return false, fmt.Errorf("no ARM contract is set")
	}
	arm, err := mock_rmn_contract.NewMockRMNContract(*ccipModule.RMNContract, ccipModule.ChainClient.Backend())
	if err != nil {
		return false, fmt.Errorf("error instantiating arm %w", err)
	}
	return arm.IsCursed0(nil)
}

func (ccipModule *CCIPCommon) CurseARM() (*types.Transaction, error) {
	if ccipModule.ARM != nil {
		return nil, fmt.Errorf("real ARM deployed. cannot curse through test")
	}
	if ccipModule.RMNContract == nil {
		return nil, fmt.Errorf("no ARM contract is set")
	}
	arm, err := mock_rmn_contract.NewMockRMNContract(*ccipModule.RMNContract, ccipModule.ChainClient.Backend())
	if err != nil {
		return nil, fmt.Errorf("error instantiating arm %w", err)
	}
	opts, err := ccipModule.ChainClient.TransactionOpts(ccipModule.ChainClient.GetDefaultWallet())
	if err != nil {
		return nil, fmt.Errorf("error getting owners for ARM VoteToCurse %w", err)
	}
	tx, err := arm.VoteToCurse(opts, [32]byte{})
	if err != nil {
		return nil, fmt.Errorf("error in calling VoteToCurse %w", err)
	}
	err = ccipModule.ChainClient.ProcessTransaction(tx)
	if err != nil {
		return tx, err
	}
	log.Info().
		Str("ARM", arm.Address().Hex()).
		Str("Network", ccipModule.ChainClient.GetNetworkName()).
		Msg("ARM is cursed")
	return tx, ccipModule.ChainClient.WaitForEvents()
}

func (ccipModule *CCIPCommon) LoadContractAddresses(conf *laneconfig.LaneConfig, noOfTokens *int) {
	if conf != nil {
		if conf.AllowOutOfOrder {
			ccipModule.AllowOutOfOrder = true
		}
		if common.IsHexAddress(conf.FeeToken) {
			ccipModule.FeeToken = &contracts.LinkToken{
				EthAddress: common.HexToAddress(conf.FeeToken),
			}
		}
		if conf.IsNativeFeeToken {
			ccipModule.FeeToken = &contracts.LinkToken{
				EthAddress: common.HexToAddress("0x0"),
			}
		}

		if common.IsHexAddress(conf.Router) {
			ccipModule.Router = &contracts.Router{
				EthAddress: common.HexToAddress(conf.Router),
			}
		}
		if common.IsHexAddress(conf.ARM) {
			addr := common.HexToAddress(conf.ARM)
			ccipModule.RMNContract = &addr
			if !conf.IsMockARM {
				ccipModule.ARM = &contracts.ARM{
					EthAddress: addr,
				}
			}
		}
		if common.IsHexAddress(conf.PriceRegistry) {
			ccipModule.PriceRegistry = &contracts.PriceRegistry{
				EthAddress: common.HexToAddress(conf.PriceRegistry),
			}
		}
		if common.IsHexAddress(conf.WrappedNative) {
			ccipModule.WrappedNative = common.HexToAddress(conf.WrappedNative)
		}
		if common.IsHexAddress(conf.Multicall) {
			ccipModule.MulticallContract = common.HexToAddress(conf.Multicall)
		}
		if common.IsHexAddress(conf.TokenMessenger) {
			addr := common.HexToAddress(conf.TokenMessenger)
			ccipModule.TokenMessenger = &addr
		}
		if common.IsHexAddress(conf.TokenTransmitter) {
			ccipModule.TokenTransmitter = &contracts.TokenTransmitter{
				ContractAddress: common.HexToAddress(conf.TokenTransmitter),
			}
		}
		if len(conf.BridgeTokens) > 0 {
			// if noOfTokens is set, then only take that many tokens from the list
			// the lane config can have more tokens than required for the test
			if noOfTokens != nil {
				if len(conf.BridgeTokens) > *noOfTokens {
					conf.BridgeTokens = conf.BridgeTokens[:*noOfTokens]
				}
			}
			var tokens []*contracts.ERC20Token
			for _, token := range conf.BridgeTokens {
				if common.IsHexAddress(token) {
					tokens = append(tokens, &contracts.ERC20Token{
						ContractAddress: common.HexToAddress(token),
					})
				}
			}
			ccipModule.BridgeTokens = tokens
		}
		if len(conf.BridgeTokenPools) > 0 {
			// if noOfTokens is set, then only take that many tokenpools from the list
			// the lane config can have more tokenpools than required for the test
			if noOfTokens != nil {
				if len(conf.BridgeTokenPools) > *noOfTokens {
					conf.BridgeTokenPools = conf.BridgeTokenPools[:*noOfTokens]
				}
			}
			var pools []*contracts.TokenPool
			for _, pool := range conf.BridgeTokenPools {
				if common.IsHexAddress(pool) {
					pools = append(pools, &contracts.TokenPool{
						EthAddress: common.HexToAddress(pool),
					})
				}
			}
			ccipModule.BridgeTokenPools = pools
		}
		if len(conf.PriceAggregators) > 0 {
			priceAggrs := make(map[common.Address]*contracts.MockAggregator)
			for token, aggr := range conf.PriceAggregators {
				if common.IsHexAddress(aggr) {
					priceAggrs[common.HexToAddress(token)] = &contracts.MockAggregator{
						ContractAddress: common.HexToAddress(aggr),
					}
				}
			}
			ccipModule.PriceAggregators = priceAggrs
		}
		if common.IsHexAddress(conf.TokenAdminRegistry) {
			ccipModule.TokenAdminRegistry = &contracts.TokenAdminRegistry{
				EthAddress: common.HexToAddress(conf.TokenAdminRegistry),
			}
		}
	}
}

// ApproveTokens approves tokens for the router to send usually a massive amount of tokens enough to cover all the ccip transfers
// to be triggered by the test.
// Also, if the test is using self-serve tokens and pools deployed by a separate `tokenDeployer` address, this sends some of those tokens
// to the default `ccipOwner` address to be used for the test.
func (ccipModule *CCIPCommon) ApproveTokens() error {
	isApproved := false
	for _, token := range ccipModule.BridgeTokens {
		// TODO: We send half of token funds back to the CCIP Deployer account, which isn't particularly realistic.
		// See CCIP-2477
		if token.OwnerWallet.Address() != ccipModule.ChainClient.GetDefaultWallet().Address() &&
			!ccipModule.ExistingDeployment {
			tokenBalance, err := token.BalanceOf(context.Background(), token.OwnerWallet.Address())
			if err != nil {
				return fmt.Errorf("failed to get balance of token %s: %w", token.ContractAddress.Hex(), err)
			}
			tokenBalance.Div(tokenBalance, big.NewInt(2)) // Send half of the balance to the default wallet
			err = token.Transfer(token.OwnerWallet, ccipModule.ChainClient.GetDefaultWallet().Address(), tokenBalance)
			if err != nil {
				return fmt.Errorf("failed to transfer token from '%s' to '%s' %s: %w",
					token.ContractAddress.Hex(), token.OwnerAddress.Hex(), ccipModule.ChainClient.GetDefaultWallet().Address(), err,
				)
			}
		}

		err := token.Approve(ccipModule.ChainClient.GetDefaultWallet(), ccipModule.Router.Address(), ApprovedAmountToRouter)
		if err != nil {
			return fmt.Errorf("failed to approve token %s: %w", token.ContractAddress.Hex(), err)
		}
		if token.ContractAddress == ccipModule.FeeToken.EthAddress {
			isApproved = true
		}
	}
	if ccipModule.FeeToken.EthAddress != common.HexToAddress("0x0") {
		amount := ApprovedFeeAmountToRouter
		if isApproved {
			amount = new(big.Int).Add(ApprovedAmountToRouter, ApprovedFeeAmountToRouter)
		}
		allowance, err := ccipModule.FeeToken.Allowance(ccipModule.ChainClient.GetDefaultWallet().Address(), ccipModule.Router.Address())
		if err != nil {
			return fmt.Errorf("failed to get allowance for token %s: %w", ccipModule.FeeToken.Address(), err)
		}
		if allowance.Cmp(amount) < 0 {
			err := ccipModule.FeeToken.Approve(ccipModule.Router.Address(), amount)
			if err != nil {
				return fmt.Errorf("failed to approve fee token %s: %w", ccipModule.FeeToken.EthAddress.String(), err)
			}
		}
	}
	ccipModule.Logger.Info().Msg("Tokens approved")

	return nil
}

func (ccipModule *CCIPCommon) CleanUp() error {
	if !ccipModule.ExistingDeployment {
		for i, pool := range ccipModule.BridgeTokenPools {
			if !pool.IsLockRelease() {
				continue
			}
			bal, err := ccipModule.BridgeTokens[i].BalanceOf(context.Background(), pool.Address())
			if err != nil {
				return fmt.Errorf("error in getting pool balance %w", err)
			}
			if bal.Cmp(big.NewInt(0)) == 0 {
				continue
			}
			err = pool.RemoveLiquidity(bal)
			if err != nil {
				return fmt.Errorf("error in removing liquidity %w", err)
			}
		}
		err := ccipModule.ChainClient.WaitForEvents()
		if err != nil {
			return fmt.Errorf("error in waiting for events %wfmt.Sprintf(\"Setting mockserver response\")", err)
		}
	}
	return nil
}

func (ccipModule *CCIPCommon) WaitForPriceUpdates(
	ctx context.Context,
	lggr *zerolog.Logger,
	timeout time.Duration,
	destChainId uint64,
	allTokens []common.Address,
) error {
	destChainSelector, err := chainselectors.SelectorFromChainId(destChainId)
	if err != nil {
		return err
	}
	// check if price is already updated
	price, err := ccipModule.PriceRegistry.Instance.GetDestinationChainGasPrice(nil, destChainSelector)
	if err != nil {
		return err
	}

	if price.Timestamp > 0 && price.Value.Cmp(big.NewInt(0)) > 0 {
		lggr.Info().
			Str("Price Registry", ccipModule.PriceRegistry.Address()).
			Uint64("dest chain", destChainId).
			Str("source chain", ccipModule.ChainClient.GetNetworkName()).
			Msg("Price already updated")
		return nil
	}
	// if not, wait for price update
	lggr.Info().Msgf("Waiting for UsdPerUnitGas and UsdPerTokenUpdated for dest chain %d Price Registry %s", destChainId, ccipModule.PriceRegistry.Address())
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	localCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	var tokensMissingForUpdate common.Address
	for {
		select {
		case <-ticker.C:
			ccipModule.gasUpdateWatcherMu.Lock()
			timestampOfUpdate, ok := ccipModule.gasUpdateWatcher[destChainId]
			ccipModule.gasUpdateWatcherMu.Unlock()
			tokenPricesUpdated := false
			if len(allTokens) > 0 {
				ccipModule.tokenPriceUpdateWatcherMu.Lock()
				for _, token := range allTokens {
					timestampOfTokenUpdate, okToken := ccipModule.tokenPriceUpdateWatcher[token]
					// we consider token prices updated only if all tokens have been updated
					// if any token is missing, we retry
					if !okToken || timestampOfTokenUpdate.Cmp(big.NewInt(0)) < 1 {
						tokenPricesUpdated = false
						tokensMissingForUpdate = token
						break
					}
					tokenPricesUpdated = true
				}
				ccipModule.tokenPriceUpdateWatcherMu.Unlock()
			}

			if tokenPricesUpdated && ok && timestampOfUpdate.Cmp(big.NewInt(0)) == 1 {
				lggr.Info().
					Str("Price Registry", ccipModule.PriceRegistry.Address()).
					Uint64("dest chain", destChainId).
					Str("source chain", ccipModule.ChainClient.GetNetworkName()).
					Msg("Price updated")
				return nil
			}
		case <-localCtx.Done():
			if tokensMissingForUpdate != (common.Address{}) {
				return fmt.Errorf("price Updates not found for token %s", tokensMissingForUpdate.Hex())
			}
			return fmt.Errorf("price Updates not found for chain %d", destChainId)
		}
	}
}

// WatchForPriceUpdates helps to ensure the price updates are happening in price registry by subscribing to a couple
// of price update events and add the event details to watchers. It subscribes to 'UsdPerUnitGasUpdated'
// and 'UsdPerTokenUpdated' event.
func (ccipModule *CCIPCommon) WatchForPriceUpdates(ctx context.Context, lggr *zerolog.Logger) error {
	gasUpdateEventLatest := make(chan *fee_quoter.FeeQuoterUsdPerUnitGasUpdated)
	tokenUpdateEvent := make(chan *fee_quoter.FeeQuoterUsdPerTokenUpdated)
	sub := event.Resubscribe(DefaultResubscriptionTimeout, func(_ context.Context) (event.Subscription, error) {
		lggr.Info().Msg("Subscribing to UsdPerUnitGasUpdated event")
		eventSub, err := ccipModule.PriceRegistry.WatchUsdPerUnitGasUpdated(nil, gasUpdateEventLatest, nil)
		if err != nil {
			log.Error().Err(err).Msg("error in subscribing to UsdPerUnitGasUpdated event")
		}
		return eventSub, err
	})
	if sub == nil {
		return fmt.Errorf("no event subscription found")
	}
	tokenUpdateSub := event.Resubscribe(DefaultResubscriptionTimeout, func(_ context.Context) (event.Subscription, error) {
		lggr.Info().Msg("Subscribing to UsdPerTokenUpdated event")
		eventSub, err := ccipModule.PriceRegistry.WatchUsdPerTokenUpdated(nil, tokenUpdateEvent)
		if err != nil {
			log.Error().Err(err).Msg("error in subscribing to UsdPerTokenUpdated event")
		}
		return eventSub, err
	})
	if tokenUpdateSub == nil {
		return fmt.Errorf("no event subscription found")
	}
	processEvent := func(value, timestamp *big.Int, destChainSelector uint64, raw types.Log) error {
		destChain, err := chainselectors.ChainIdFromSelector(destChainSelector)
		if err != nil {
			return err
		}
		ccipModule.gasUpdateWatcherMu.Lock()
		ccipModule.gasUpdateWatcher[destChain] = timestamp

		ccipModule.GasUpdateEvents = append(ccipModule.GasUpdateEvents, contracts.GasUpdateEvent{
			Sender:    raw.Address.Hex(),
			Tx:        raw.TxHash.Hex(),
			Value:     value,
			DestChain: destChain,
			Source:    ccipModule.ChainClient.GetNetworkName(),
		})
		ccipModule.gasUpdateWatcherMu.Unlock()
		lggr.Info().
			Uint64("chainSelector", destChainSelector).
			Uint64("dest_chain", destChain).
			Str("price_registry", ccipModule.PriceRegistry.Address()).
			Str("tx hash", raw.TxHash.Hex()).
			Msgf("UsdPerUnitGasUpdated event received for dest chain: %d, source chain: %s",
				destChain, ccipModule.ChainClient.GetNetworkName())
		return nil
	}
	go func() {
		defer func() {
			sub.Unsubscribe()
			tokenUpdateSub.Unsubscribe()
			ccipModule.gasUpdateWatcher = nil
			ccipModule.gasUpdateWatcherMu = nil
			ccipModule.GasUpdateEvents = nil
			ccipModule.tokenPriceUpdateWatcher = nil
			ccipModule.tokenPriceUpdateWatcherMu = nil
		}()
		for {
			select {
			case e := <-gasUpdateEventLatest:
				err := processEvent(e.Value, e.Timestamp, e.DestChain, e.Raw)
				if err != nil {
					continue
				}
			case tk := <-tokenUpdateEvent:
				ccipModule.tokenPriceUpdateWatcherMu.Lock()
				ccipModule.tokenPriceUpdateWatcher[tk.Token] = tk.Timestamp
				ccipModule.tokenPriceUpdateWatcherMu.Unlock()
				lggr.Info().
					Str("token", tk.Token.Hex()).
					Str("chain", ccipModule.ChainClient.GetNetworkName()).
					Str("price_registry", ccipModule.PriceRegistry.Address()).
					Msg("UsdPerTokenUpdated event received")
			case <-ctx.Done():
				return
			}
		}
	}()

	return nil
}

// UpdateTokenPricesAtRegularInterval updates aggregator contract with updated answer at regular interval.
// At each iteration of ticker it chooses one of the aggregator contracts and updates its round answer.
func (ccipModule *CCIPCommon) UpdateTokenPricesAtRegularInterval(ctx context.Context, lggr *zerolog.Logger, interval time.Duration, conf *laneconfig.LaneConfig) error {
	if ccipModule.ExistingDeployment {
		return nil
	}
	var aggregators []*contracts.MockAggregator
	for _, aggregatorContract := range conf.PriceAggregators {
		contract, err := ccipModule.Deployer.NewMockAggregator(common.HexToAddress(aggregatorContract))
		if err != nil {
			return err
		}
		aggregators = append(aggregators, contract)
	}
	go func(aggregators []*contracts.MockAggregator) {
		rand.NewSource(uint64(time.Now().UnixNano()))
		ticker := time.NewTicker(interval)
		for {
			select {
			case <-ticker.C:
				// randomly choose an aggregator contract from slice of aggregators
				randomIndex := rand.Intn(len(aggregators))
				err := aggregators[randomIndex].UpdateRoundData(nil, ptr.Ptr(-5), ptr.Ptr(2))
				if err != nil {
					lggr.Error().Err(err).Msg("error in updating round data")
					continue
				}
			case <-ctx.Done():
				return
			}
		}
	}(aggregators)
	return nil
}

// SyncUSDCDomain makes domain updates to Source usdc pool domain with -
// 1. USDC domain from destination chain's token transmitter contract
// 2. Destination pool address as allowed caller
func (ccipModule *CCIPCommon) SyncUSDCDomain(destTransmitter *contracts.TokenTransmitter, destPools []*contracts.TokenPool, destChainID uint64) error {
	// if not USDC new deployment, return
	// if existing deployment, consider that no syncing is required and return
	if ccipModule.ExistingDeployment || !ccipModule.IsUSDCDeployment() {
		return nil
	}
	if destTransmitter == nil {
		return fmt.Errorf("invalid address")
	}
	destChainSelector, err := chainselectors.SelectorFromChainId(destChainID)
	if err != nil {
		return fmt.Errorf("invalid chain id %w", err)
	}

	// sync USDC domain
	for i, pool := range ccipModule.BridgeTokenPools {
		if !pool.IsUSDC() {
			continue
		}
		if destPools[i] == nil {
			return fmt.Errorf("invalid pool address")
		}
		if !destPools[i].IsUSDC() {
			return fmt.Errorf("corresponding dest pool is not USDC pool")
		}
		err = pool.SyncUSDCDomain(destTransmitter, destPools[i].EthAddress, destChainSelector)
		if err != nil {
			return err
		}
		err = destPools[i].MintUSDCToUSDCPool()
		if err != nil {
			return err
		}
	}

	return ccipModule.ChainClient.WaitForEvents()
}

func (ccipModule *CCIPCommon) PollRPCConnection(ctx context.Context, lggr *zerolog.Logger) {
	for {
		select {
		case reconnectTime := <-ccipModule.ChainClient.ConnectionRestored():
			if ccipModule.IsConnectionRestoredRecently == nil {
				ccipModule.IsConnectionRestoredRecently = atomic.NewBool(true)
			} else {
				ccipModule.IsConnectionRestoredRecently.Store(true)
			}
			lggr.Info().Time("Restored At", reconnectTime).Str("Network", ccipModule.ChainClient.GetNetworkName()).Msg("Connection Restored")
		case issueTime := <-ccipModule.ChainClient.ConnectionIssue():
			if ccipModule.IsConnectionRestoredRecently == nil {
				ccipModule.IsConnectionRestoredRecently = atomic.NewBool(false)
			} else {
				ccipModule.IsConnectionRestoredRecently.Store(false)
			}
			lggr.Info().Time("Started At", issueTime).Str("Network", ccipModule.ChainClient.GetNetworkName()).Msg("RPC Disconnected")
		case <-ctx.Done():
			return
		}
	}
}

func (ccipModule *CCIPCommon) IsUSDCDeployment() bool {
	return pointer.GetBool(ccipModule.USDCMockDeployment)
}

func (ccipModule *CCIPCommon) IsLBTCDeployment() bool {
	return pointer.GetBool(ccipModule.LBTCMockDeployment)
}

func (ccipModule *CCIPCommon) WriteLaneConfig(conf *laneconfig.LaneConfig) {
	var btAddresses, btpAddresses []string
	priceAggrs := make(map[string]string)
	for i, bt := range ccipModule.BridgeTokens {
		btAddresses = append(btAddresses, bt.Address())
		btpAddresses = append(btpAddresses, ccipModule.BridgeTokenPools[i].Address())
	}
	for k, v := range ccipModule.PriceAggregators {
		priceAggrs[k.Hex()] = v.ContractAddress.Hex()
	}
	conf.CommonContracts.FeeToken = ccipModule.FeeToken.Address()
	conf.CommonContracts.BridgeTokens = btAddresses
	conf.CommonContracts.BridgeTokenPools = btpAddresses
	conf.CommonContracts.ARM = ccipModule.RMNContract.Hex()
	conf.CommonContracts.Router = ccipModule.Router.Address()
	conf.CommonContracts.PriceRegistry = ccipModule.PriceRegistry.Address()
	conf.CommonContracts.PriceAggregators = priceAggrs
	conf.CommonContracts.WrappedNative = ccipModule.WrappedNative.Hex()
	conf.CommonContracts.Multicall = ccipModule.MulticallContract.Hex()

	if ccipModule.TokenAdminRegistry != nil {
		conf.CommonContracts.TokenAdminRegistry = ccipModule.TokenAdminRegistry.Address()
	}
	if ccipModule.TokenTransmitter != nil {
		conf.CommonContracts.TokenTransmitter = ccipModule.TokenTransmitter.ContractAddress.Hex()
	}
	if ccipModule.TokenMessenger != nil {
		conf.CommonContracts.TokenMessenger = ccipModule.TokenMessenger.Hex()
	}
	if ccipModule.ARM == nil {
		conf.CommonContracts.IsMockARM = true
	}
}

func (ccipModule *CCIPCommon) AddPriceAggregatorToken(token common.Address, initialAns *big.Int) error {
	// check if dynamic price update is enabled
	if ccipModule.NoOfTokensNeedingDynamicPrice <= 0 {
		return nil
	}
	var err error
	if aggregator, ok := ccipModule.PriceAggregators[token]; !ok {
		ccipModule.PriceAggregators[token], err = ccipModule.Deployer.DeployMockAggregator(18, initialAns)
		if err != nil {
			return fmt.Errorf("deploying mock aggregator contract shouldn't fail %w", err)
		}
	} else {
		ccipModule.PriceAggregators[token], err = ccipModule.Deployer.NewMockAggregator(aggregator.ContractAddress)
		if err != nil {
			return fmt.Errorf("error instantiating price aggregator for token %s", token.Hex())
		}
	}
	ccipModule.NoOfTokensNeedingDynamicPrice--
	return nil
}

// DeployContracts deploys the contracts which are necessary in both source and dest chain
// This reuses common contracts for bidirectional lanes
func (ccipModule *CCIPCommon) DeployContracts(
	noOfTokens int,
	tokenDeployerFns []blockchain.ContractDeployer,
	conf *laneconfig.LaneConfig,
) error {
	var err error
	cd := ccipModule.Deployer

	ccipModule.LoadContractAddresses(conf, &noOfTokens)
	if ccipModule.ARM != nil {
		arm, err := cd.NewRMNContract(ccipModule.ARM.EthAddress)
		if err != nil {
			return fmt.Errorf("getting new ARM contract shouldn't fail %w", err)
		}
		ccipModule.ARM = arm
	} else {
		// deploy a mock ARM contract
		if ccipModule.RMNContract == nil {
			if ccipModule.ExistingDeployment {
				return fmt.Errorf("ARM contract address is not provided in lane config")
			}
			ccipModule.RMNContract, err = cd.DeployMockRMNContract()
			if err != nil {
				return fmt.Errorf("deploying mock ARM contract shouldn't fail %w", err)
			}
			err = ccipModule.ChainClient.WaitForEvents()
			if err != nil {
				return fmt.Errorf("error in waiting for mock ARM deployment %w", err)
			}
		}
	}
	if ccipModule.WrappedNative == common.HexToAddress("0x0") {
		if ccipModule.ExistingDeployment {
			return fmt.Errorf("wrapped native contract address is not provided in lane config")
		}
		weth9addr, err := cd.DeployWrappedNative()
		if err != nil {
			return fmt.Errorf("deploying wrapped native shouldn't fail %w", err)
		}
		err = ccipModule.AddPriceAggregatorToken(*weth9addr, WrappedNativeToUSD)
		if err != nil {
			return fmt.Errorf("deploying mock aggregator contract shouldn't fail %w", err)
		}
		err = ccipModule.ChainClient.WaitForEvents()
		if err != nil {
			return fmt.Errorf("waiting for deploying wrapped native shouldn't fail %w", err)
		}
		ccipModule.WrappedNative = *weth9addr
	}

	if ccipModule.Router == nil {
		if ccipModule.ExistingDeployment {
			return fmt.Errorf("router contract address is not provided in lane config")
		}
		ccipModule.Router, err = cd.DeployRouter(ccipModule.WrappedNative, *ccipModule.RMNContract)
		if err != nil {
			return fmt.Errorf("deploying router shouldn't fail %w", err)
		}
		err = ccipModule.ChainClient.WaitForEvents()
		if err != nil {
			return fmt.Errorf("error in waiting for router deployment %w", err)
		}
	} else {
		r, err := cd.NewRouter(ccipModule.Router.EthAddress)
		if err != nil {
			return fmt.Errorf("getting new router contract shouldn't fail %w", err)
		}
		ccipModule.Router = r
	}
	if ccipModule.FeeToken == nil {
		if ccipModule.ExistingDeployment {
			return fmt.Errorf("FeeToken contract address is not provided in lane config")
		}
		// deploy link token
		token, err := cd.DeployLinkTokenContract()
		if err != nil {
			return fmt.Errorf("deploying fee token contract shouldn't fail %w", err)
		}

		ccipModule.FeeToken = token
		err = ccipModule.AddPriceAggregatorToken(ccipModule.FeeToken.EthAddress, LinkToUSD)
		if err != nil {
			return fmt.Errorf("deploying mock aggregator contract shouldn't fail %w", err)
		}
		err = ccipModule.ChainClient.WaitForEvents()
		if err != nil {
			return fmt.Errorf("error in waiting for feetoken deployment %w", err)
		}
	} else {
		token, err := cd.NewLinkTokenContract(common.HexToAddress(ccipModule.FeeToken.Address()))
		if err != nil {
			return fmt.Errorf("getting fee token contract shouldn't fail %w", err)
		}
		ccipModule.FeeToken = token
	}

	// If the number of deployed bridge tokens does not match noOfTokens, deploy rest of the tokens in case ExistingDeployment is false
	// In case of ExistingDeployment as true use whatever is provided in laneconfig
	if len(ccipModule.BridgeTokens) < noOfTokens && !ccipModule.ExistingDeployment {
		// deploy bridge token.
		for i := len(ccipModule.BridgeTokens); i < noOfTokens; i++ {
			var token *contracts.ERC20Token
			if len(tokenDeployerFns) != noOfTokens {
				if ccipModule.IsUSDCDeployment() && i == 0 {
					// if it's USDC deployment, we deploy the burn mint token 677 with decimal 6 and cast it to ERC20Token
					usdcToken, err := ccipModule.tokenDeployer.DeployBurnMintERC677(new(big.Int).Mul(big.NewInt(1e6), big.NewInt(1e18)))
					if err != nil {
						return fmt.Errorf("deploying bridge usdc token contract shouldn't fail %w", err)
					}
					token, err = ccipModule.tokenDeployer.NewERC20TokenContract(usdcToken.ContractAddress)
					if err != nil {
						return fmt.Errorf("getting new bridge usdc token contract shouldn't fail %w", err)
					}
					if ccipModule.TokenTransmitter == nil {
						domain, err := GetUSDCDomain(ccipModule.ChainClient.GetNetworkName(), ccipModule.ChainClient.NetworkSimulated())
						if err != nil {
							return fmt.Errorf("error in getting USDC domain %w", err)
						}

						ccipModule.TokenTransmitter, err = ccipModule.tokenDeployer.DeployTokenTransmitter(domain, usdcToken.ContractAddress)
						if err != nil {
							return fmt.Errorf("deploying token transmitter shouldn't fail %w", err)
						}
					}
					if ccipModule.TokenMessenger == nil {
						if ccipModule.TokenTransmitter == nil {
							return fmt.Errorf("TokenTransmitter contract address is not provided")
						}
						ccipModule.TokenMessenger, err = ccipModule.tokenDeployer.DeployTokenMessenger(ccipModule.TokenTransmitter.ContractAddress)
						if err != nil {
							return fmt.Errorf("deploying token messenger shouldn't fail %w", err)
						}
						err = ccipModule.ChainClient.WaitForEvents()
						if err != nil {
							return fmt.Errorf("error in waiting for mock TokenMessenger and Transmitter deployment %w", err)
						}
					}

					// grant minter role to token messenger
					err = usdcToken.GrantMintAndBurn(*ccipModule.TokenMessenger)
					if err != nil {
						return fmt.Errorf("granting minter role to token messenger shouldn't fail %w", err)
					}
					err = usdcToken.GrantMintAndBurn(ccipModule.TokenTransmitter.ContractAddress)
					if err != nil {
						return fmt.Errorf("granting minter role to token transmitter shouldn't fail %w", err)
					}
				} else if ccipModule.IsLBTCDeployment() && i == 0 {
					// if it's LBTC deployment, we deploy the burn mint token 677 with decimal 8 and cast it to ERC20Token
					lbtcToken, err := ccipModule.tokenDeployer.DeployCustomBurnMintERC677Token("Lombard LBTC", "LBTC", uint8(8), new(big.Int).Mul(big.NewInt(1e6), big.NewInt(1e18)))
					if err != nil {
						return fmt.Errorf("deploying bridge lbtc token contract shouldn't fail %w", err)
					}
					token, err = ccipModule.tokenDeployer.NewERC20TokenContract(lbtcToken.ContractAddress)
					if err != nil {
						return fmt.Errorf("getting new bridge lbtc token contract shouldn't fail %w", err)
					}
				} else {
					// otherwise we deploy link token and cast it to ERC20Token
					linkToken, err := ccipModule.tokenDeployer.DeployLinkTokenContract()
					if err != nil {
						return fmt.Errorf("deploying bridge token contract shouldn't fail %w", err)
					}
					token, err = ccipModule.tokenDeployer.NewERC20TokenContract(common.HexToAddress(linkToken.Address()))
					if err != nil {
						return fmt.Errorf("getting new bridge token contract shouldn't fail %w", err)
					}
					err = ccipModule.AddPriceAggregatorToken(linkToken.EthAddress, LinkToUSD)
					if err != nil {
						return fmt.Errorf("deploying mock aggregator contract shouldn't fail %w", err)
					}
				}
			} else {
				token, err = ccipModule.tokenDeployer.DeployERC20TokenContract(tokenDeployerFns[i])
				if err != nil {
					return fmt.Errorf("deploying bridge token contract shouldn't fail %w", err)
				}
				err = ccipModule.AddPriceAggregatorToken(token.ContractAddress, LinkToUSD)
				if err != nil {
					return fmt.Errorf("deploying mock aggregator contract shouldn't fail %w", err)
				}
			}
			ccipModule.BridgeTokens = append(ccipModule.BridgeTokens, token)

		}
		if err = ccipModule.ChainClient.WaitForEvents(); err != nil {
			return fmt.Errorf("error in waiting for bridge token deployment %w", err)
		}
	}

	var tokens []*contracts.ERC20Token
	for _, token := range ccipModule.BridgeTokens {
		newToken, err := ccipModule.tokenDeployer.NewERC20TokenContract(common.HexToAddress(token.Address()))
		if err != nil {
			return fmt.Errorf("getting new bridge token contract shouldn't fail %w", err)
		}
		tokens = append(tokens, newToken)
	}
	ccipModule.BridgeTokens = tokens
	if len(ccipModule.BridgeTokenPools) != len(ccipModule.BridgeTokens) {
		if ccipModule.ExistingDeployment {
			return fmt.Errorf("bridge token pool contract address is not provided in lane config")
		}
		// deploy native token pool
		for i := len(ccipModule.BridgeTokenPools); i < len(ccipModule.BridgeTokens); i++ {
			token := ccipModule.BridgeTokens[i]
			// usdc pool need to be the first one in the slice
			if ccipModule.IsUSDCDeployment() && i == 0 {
				// deploy usdc token pool in case of usdc deployment
				if ccipModule.TokenMessenger == nil {
					return fmt.Errorf("TokenMessenger contract address is not provided")
				}
				if ccipModule.TokenTransmitter == nil {
					return fmt.Errorf("TokenTransmitter contract address is not provided")
				}
				usdcPool, err := ccipModule.tokenDeployer.DeployUSDCTokenPoolContract(token.Address(), *ccipModule.TokenMessenger, *ccipModule.RMNContract, ccipModule.Router.Instance.Address())
				if err != nil {
					return fmt.Errorf("deploying bridge Token pool(usdc) shouldn't fail %w", err)
				}

				ccipModule.BridgeTokenPools = append(ccipModule.BridgeTokenPools, usdcPool)
			} else if ccipModule.IsLBTCDeployment() && i == 0 {
				if ccipModule.RMNContract == nil {
					return fmt.Errorf("RMNContract is not initialized")
				}
				rmnContract := *ccipModule.RMNContract

				destPoolData, err := hex.DecodeString(LBTCValidDestPoolData) // valid 32 bytes should call attestation api
				if err != nil {
					return errors.Wrapf(err, "decoding dest pool data shouldn't fail")
				}
				if !pointer.GetBool(ccipModule.LBTCDestPoolDataAs32Bytes) {
					// non 32 bytes data should not call attestation api and instead consider it as deposit payload.
					// lombard has both attested and non-attested flow
					destPoolData = []byte{0x12, 0x34, 0x56, 0x78}
				}
				lbtcPool, err := ccipModule.tokenDeployer.DeployMockLBTCTokenPoolContract(token.Address(), rmnContract, ccipModule.Router.Instance.Address(), destPoolData)
				if err != nil {
					return errors.Wrapf(err, "deploying mock lbtc bridge token pool shouldn't fail")
				}
				lbtcInstance, err := burn_mint_erc677.NewBurnMintERC677(token.ContractAddress, ccipModule.ChainClient.Backend())
				if err != nil {
					return errors.Wrapf(err, "failed to get dest usdc token instance")
				}
				opts, err := ccipModule.tokenDeployer.Client().TransactionOpts(token.OwnerWallet)
				if err != nil {
					return errors.Wrapf(err, "failed to get transaction opts")
				}
				tx, err := lbtcInstance.GrantMintAndBurnRoles(opts, common.HexToAddress(lbtcPool.Address()))
				if err != nil {
					return errors.Wrapf(err, "granting minter role to owner shouldn't fail")
				}
				err = ccipModule.tokenDeployer.Client().ProcessTransaction(tx)
				if err != nil {
					return errors.Wrapf(err, "failed to process grant mint role")
				}

				ccipModule.BridgeTokenPools = append(ccipModule.BridgeTokenPools, lbtcPool)
			} else {
				// deploy lock release token pool in case of non-usdc deployment
				btp, err := ccipModule.tokenDeployer.DeployLockReleaseTokenPoolContract(token.Address(), *ccipModule.RMNContract, ccipModule.Router.Instance.Address())
				if err != nil {
					return fmt.Errorf("deploying bridge Token pool(lock&release) shouldn't fail %w", err)
				}
				ccipModule.BridgeTokenPools = append(ccipModule.BridgeTokenPools, btp)

				err = btp.AddLiquidity(token, token.OwnerWallet, ccipModule.poolFunds)
				if err != nil {
					return fmt.Errorf("adding liquidity token to dest pool shouldn't fail %w", err)
				}
			}
		}
	} else {
		var pools []*contracts.TokenPool
		for _, pool := range ccipModule.BridgeTokenPools {
			newPool, err := ccipModule.tokenDeployer.NewLockReleaseTokenPoolContract(pool.EthAddress)
			if err != nil {
				return fmt.Errorf("getting new bridge token pool contract shouldn't fail %w", err)
			}
			pools = append(pools, newPool)
		}
		ccipModule.BridgeTokenPools = pools
	}

	// no need to have price registry for existing deployment, we consider that it's already deployed
	if !ccipModule.ExistingDeployment {
		if ccipModule.PriceRegistry == nil {
			// we will update the price updates later based on source and dest PriceUpdates
			ccipModule.PriceRegistry, err = cd.DeployPriceRegistry(
				[]common.Address{
					common.HexToAddress(ccipModule.FeeToken.Address()),
					common.HexToAddress(ccipModule.WrappedNative.Hex()),
				})
			if err != nil {
				return fmt.Errorf("deploying PriceRegistry shouldn't fail %w", err)
			}
			err = ccipModule.ChainClient.WaitForEvents()
			if err != nil {
				return fmt.Errorf("error in waiting for PriceRegistry deployment %w", err)
			}
		} else {
			ccipModule.PriceRegistry, err = cd.NewPriceRegistry(ccipModule.PriceRegistry.EthAddress)
			if err != nil {
				return fmt.Errorf("getting new PriceRegistry contract shouldn't fail %w", err)
			}
		}
	}
	if ccipModule.MulticallContract == (common.Address{}) && ccipModule.MulticallEnabled {
		ccipModule.MulticallContract, err = cd.DeployMultiCallContract()
		if err != nil {
			return fmt.Errorf("deploying multicall contract shouldn't fail %w", err)
		}
	}

	// if the version is after 1.4.0, we need to deploy TokenAdminRegistry
	// no need to have token admin registry for existing deployment, we consider that it's already deployed
	if contracts.NeedTokenAdminRegistry() && !ccipModule.ExistingDeployment {
		if ccipModule.TokenAdminRegistry == nil {
			// deploy token admin registry
			ccipModule.TokenAdminRegistry, err = cd.DeployTokenAdminRegistry()
			if err != nil {
				return fmt.Errorf("deploying token admin registry shouldn't fail %w", err)
			}
			err = ccipModule.ChainClient.WaitForEvents()
			if err != nil {
				return fmt.Errorf("error in waiting for token admin registry deployment %w", err)
			}

			if len(ccipModule.BridgeTokens) != len(ccipModule.BridgeTokenPools) {
				return fmt.Errorf("tokens number %d and pools number %d do not match", len(ccipModule.BridgeTokens), len(ccipModule.BridgeTokenPools))
			}
			// add all pools to registry
			for i, pool := range ccipModule.BridgeTokenPools {
				token := ccipModule.BridgeTokens[i]
				err := ccipModule.TokenAdminRegistry.SetAdminAndRegisterPool(token.ContractAddress, pool.EthAddress)
				if err != nil {
					return fmt.Errorf("error setting up token %s and pool %s on TokenAdminRegistry : %w", token.Address(), pool.Address(), err)
				}
			}
			err = ccipModule.ChainClient.WaitForEvents()
			if err != nil {
				return fmt.Errorf("error in waiting for token admin registry set up with tokens and pools %w", err)
			}
		} else {
			ccipModule.TokenAdminRegistry, err = cd.NewTokenAdminRegistry(ccipModule.TokenAdminRegistry.EthAddress)
			if err != nil {
				return fmt.Errorf("getting new token admin registry contract shouldn't fail %w", err)
			}
		}
	}
	ccipModule.Logger.Info().Msg("Finished deploying common contracts")
	// approve router to spend fee token
	return ccipModule.ApproveTokens()
}

func (ccipModule *CCIPCommon) AvgBlockTime(ctx context.Context) (time.Duration, error) {
	return ccipModule.ChainClient.AvgBlockTime(ctx)
}

// DynamicPriceGetterConfig specifies the configuration for the price getter in price pipeline.
// This should match pricegetter.DynamicPriceGetterConfig in core/services/ocr2/plugins/ccip/internal/pricegetter
type DynamicPriceGetterConfig struct {
	AggregatorPrices map[common.Address]AggregatorPriceConfig `json:"aggregatorPrices"`
	StaticPrices     map[common.Address]StaticPriceConfig     `json:"staticPrices"`
}

func (d *DynamicPriceGetterConfig) AddPriceConfig(
	tokenAddr string,
	aggregatorMap map[common.Address]*contracts.MockAggregator,
	price *big.Int,
	chainID uint64,
) error {
	aggregatorContract, ok := aggregatorMap[common.HexToAddress(tokenAddr)]
	if !ok || aggregatorContract == nil {
		return d.AddStaticPriceConfig(tokenAddr, chainID, price)
	}
	return d.AddAggregatorPriceConfig(tokenAddr, aggregatorMap, price)
}

func (d *DynamicPriceGetterConfig) AddAggregatorPriceConfig(
	tokenAddr string,
	aggregatorMap map[common.Address]*contracts.MockAggregator,
	price *big.Int,
) error {
	aggregatorContract, ok := aggregatorMap[common.HexToAddress(tokenAddr)]
	if !ok || aggregatorContract == nil {
		return fmt.Errorf("aggregator contract not found for token %s", tokenAddr)
	}
	// update round Data
	err := aggregatorContract.UpdateRoundData(price, nil, nil)
	if err != nil {
		return fmt.Errorf("error in updating round data %w", err)
	}

	d.AggregatorPrices[common.HexToAddress(tokenAddr)] = AggregatorPriceConfig{
		ChainID:                   aggregatorContract.ChainID(),
		AggregatorContractAddress: aggregatorContract.ContractAddress,
	}
	return nil
}

func (d *DynamicPriceGetterConfig) AddStaticPriceConfig(tokenAddr string, chainID uint64, price *big.Int) error {
	d.StaticPrices[common.HexToAddress(tokenAddr)] = StaticPriceConfig{
		ChainID: chainID,
		Price:   price,
	}
	return nil
}

func (d *DynamicPriceGetterConfig) String() (string, error) {
	tokenPricesConfigBytes, err := json.MarshalIndent(d, "", " ")
	if err != nil {
		return "", fmt.Errorf("error in marshalling token prices config %w", err)
	}
	return string(tokenPricesConfigBytes), nil
}

// AggregatorPriceConfig specifies a price retrieved from an aggregator contract.
// This should match pricegetter.AggregatorPriceConfig in core/services/ocr2/plugins/ccip/internal/pricegetter
type AggregatorPriceConfig struct {
	ChainID                   uint64         `json:"chainID,string"`
	AggregatorContractAddress common.Address `json:"contractAddress"`
}

// StaticPriceConfig specifies a price defined statically.
// This should match pricegetter.StaticPriceConfig in core/services/ocr2/plugins/ccip/internal/pricegetter
type StaticPriceConfig struct {
	ChainID uint64   `json:"chainID,string"`
	Price   *big.Int `json:"price"`
}

func NewCCIPCommonFromConfig(
	logger *zerolog.Logger,
	testGroupConf *testconfig.CCIPTestGroupConfig,
	chainClient blockchain.EVMClient,
	laneConfig *laneconfig.LaneConfig,
) (*CCIPCommon, error) {
	newCCIPModule, err := DefaultCCIPModule(logger, testGroupConf, chainClient)
	if err != nil {
		return nil, err
	}
	newCD := newCCIPModule.Deployer
	newCCIPModule.LoadContractAddresses(laneConfig, testGroupConf.TokenConfig.NoOfTokensPerChain)
	if newCCIPModule.TokenAdminRegistry != nil {
		newCCIPModule.TokenAdminRegistry, err = newCD.NewTokenAdminRegistry(common.HexToAddress(newCCIPModule.TokenAdminRegistry.Address()))
		if err != nil {
			return nil, err
		}
	}
	var arm *contracts.ARM
	if newCCIPModule.ARM != nil {
		arm, err = newCD.NewRMNContract(*newCCIPModule.RMNContract)
		if err != nil {
			return nil, err
		}
		newCCIPModule.ARM = arm
	}
	var pools []*contracts.TokenPool
	for i := range newCCIPModule.BridgeTokenPools {
		// if there is usdc token, the corresponding pool will always be added as first one in the slice
		if newCCIPModule.IsUSDCDeployment() && i == 0 {
			pool, err := newCCIPModule.tokenDeployer.NewUSDCTokenPoolContract(common.HexToAddress(newCCIPModule.BridgeTokenPools[i].Address()))
			if err != nil {
				return nil, err
			}
			pools = append(pools, pool)
		} else {
			pool, err := newCCIPModule.tokenDeployer.NewLockReleaseTokenPoolContract(common.HexToAddress(newCCIPModule.BridgeTokenPools[i].Address()))
			if err != nil {
				return nil, err
			}
			pools = append(pools, pool)
		}
	}
	newCCIPModule.BridgeTokenPools = pools
	var tokens []*contracts.ERC20Token
	for i := range newCCIPModule.BridgeTokens {
		token, err := newCCIPModule.tokenDeployer.NewERC20TokenContract(common.HexToAddress(newCCIPModule.BridgeTokens[i].Address()))
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, token)
	}
	newCCIPModule.BridgeTokens = tokens
	priceAggregators := make(map[common.Address]*contracts.MockAggregator)
	for k, v := range newCCIPModule.PriceAggregators {
		aggregator, err := newCD.NewMockAggregator(v.ContractAddress)
		if err != nil {
			return nil, err
		}
		priceAggregators[k] = aggregator
	}
	newCCIPModule.PriceAggregators = priceAggregators
	newCCIPModule.FeeToken, err = newCCIPModule.Deployer.NewLinkTokenContract(common.HexToAddress(newCCIPModule.FeeToken.Address()))
	if err != nil {
		return nil, err
	}
	if newCCIPModule.PriceRegistry != nil {
		newCCIPModule.PriceRegistry, err = newCCIPModule.Deployer.NewPriceRegistry(common.HexToAddress(newCCIPModule.PriceRegistry.Address()))
		if err != nil {
			return nil, err
		}
	}
	newCCIPModule.Router, err = newCCIPModule.Deployer.NewRouter(common.HexToAddress(newCCIPModule.Router.Address()))
	if err != nil {
		return nil, err
	}
	if newCCIPModule.TokenTransmitter != nil {
		newCCIPModule.TokenTransmitter, err = newCCIPModule.Deployer.NewTokenTransmitter(newCCIPModule.TokenTransmitter.ContractAddress)
		if err != nil {
			return nil, err
		}
	}
	return newCCIPModule, nil
}

func DefaultCCIPModule(
	logger *zerolog.Logger,
	testGroupConf *testconfig.CCIPTestGroupConfig,
	chainClient blockchain.EVMClient,
) (*CCIPCommon, error) {
	networkCfg := chainClient.GetNetworkConfig()
	tokenDeployerChainClient, err := blockchain.ConcurrentEVMClient(*networkCfg, nil, chainClient, *logger)
	if err != nil {
		return nil, errors.WithStack(fmt.Errorf("failed to create token deployment chain client for %s: %w", networkCfg.Name, err))
	}
	// If we want to deploy tokens as a non CCIP owner, we need to set the default wallet to something other than the first one. The first wallet is used as default CCIP owner for all other ccip contract deployment.
	// This is not needed for existing deployment as the tokens and pools are already deployed.
	if contracts.NeedTokenAdminRegistry() &&
		!pointer.GetBool(testGroupConf.TokenConfig.CCIPOwnerTokens) &&
		!pointer.GetBool(testGroupConf.ExistingDeployment) &&
		len(tokenDeployerChainClient.GetWallets()) > 1 {
		if err = tokenDeployerChainClient.SetDefaultWallet(1); err != nil {
			return nil, errors.WithStack(fmt.Errorf("failed to set default wallet for token deployment client %s: %w", networkCfg.Name, err))
		}
	}
	cd, err := contracts.NewCCIPContractsDeployer(logger, chainClient)
	if err != nil {
		return nil, err
	}
	tokenCD, err := contracts.NewCCIPContractsDeployer(logger, tokenDeployerChainClient)
	if err != nil {
		return nil, err
	}
	return &CCIPCommon{
		Logger:        logger,
		ChainClient:   chainClient,
		Deployer:      cd,
		tokenDeployer: tokenCD,
		RateLimiterConfig: contracts.RateLimiterConfig{
			Rate:     contracts.FiftyCoins,
			Capacity: contracts.HundredCoins,
		},
		ExistingDeployment:            pointer.GetBool(testGroupConf.ExistingDeployment),
		MulticallEnabled:              pointer.GetBool(testGroupConf.MulticallInOneTx),
		USDCMockDeployment:            testGroupConf.USDCMockDeployment,
		LBTCMockDeployment:            testGroupConf.LBTCMockDeployment,
		LBTCDestPoolDataAs32Bytes:     testGroupConf.LBTCDestPoolDataAs32Bytes,
		NoOfTokensNeedingDynamicPrice: pointer.GetInt(testGroupConf.TokenConfig.NoOfTokensWithDynamicPrice),
		poolFunds:                     testhelpers.Link(5),
		gasUpdateWatcherMu:            &sync.Mutex{},
		gasUpdateWatcher:              make(map[uint64]*big.Int),
		tokenPriceUpdateWatcherMu:     &sync.Mutex{},
		tokenPriceUpdateWatcher:       make(map[common.Address]*big.Int),
		PriceAggregators:              make(map[common.Address]*contracts.MockAggregator),
	}, nil
}

type SourceCCIPModule struct {
	Common                     *CCIPCommon
	Sender                     common.Address
	TransferAmount             []*big.Int
	MsgDataLength              int64
	DestinationChainId         uint64
	DestChainSelector          uint64
	DestNetworkName            string
	OnRamp                     *contracts.OnRamp
	SrcStartBlock              uint64
	CCIPSendRequestedWatcher   *sync.Map // map[string]*evm_2_evm_onramp.EVM2EVMOnRampCCIPSendRequested
	NewFinalizedBlockNum       atomic.Uint64
	NewFinalizedBlockTimestamp atomic.Time
}

func (sourceCCIP *SourceCCIPModule) PayCCIPFeeToOwnerAddress() error {
	isNativeFee := sourceCCIP.Common.FeeToken.EthAddress == common.HexToAddress("0x0")
	if isNativeFee {
		err := sourceCCIP.OnRamp.WithdrawNonLinkFees(sourceCCIP.Common.WrappedNative)
		if err != nil {
			return err
		}
	} else {
		err := sourceCCIP.OnRamp.SetNops()
		if err != nil {
			return err
		}
		err = sourceCCIP.OnRamp.PayNops()
		if err != nil {
			return err
		}
	}
	return nil
}

func (sourceCCIP *SourceCCIPModule) LoadContracts(conf *laneconfig.LaneConfig) {
	if conf != nil {
		cfg, ok := conf.SrcContracts[sourceCCIP.DestNetworkName]
		if ok {
			if common.IsHexAddress(cfg.OnRamp) {
				sourceCCIP.OnRamp = &contracts.OnRamp{
					EthAddress: common.HexToAddress(cfg.OnRamp),
				}
			}
			if cfg.DeployedAt > 0 {
				sourceCCIP.SrcStartBlock = cfg.DeployedAt
			}
		}
	}
}

// SetAllTokenTransferFeeConfigs sets a default transfer fee config for all BridgeTokens on the CCIP source chain.
// enableAggregateRateLimit is used to enable/disable aggregate rate limit for all BridgeTokens.
func (sourceCCIP *SourceCCIPModule) SetAllTokenTransferFeeConfigs(enableAggregateRateLimit bool) error {
	var tokenTransferFeeConfig []evm_2_evm_onramp.EVM2EVMOnRampTokenTransferFeeConfigArgs
	var tokens, pools []common.Address
	if len(sourceCCIP.Common.BridgeTokens) != len(sourceCCIP.Common.BridgeTokenPools) {
		return fmt.Errorf("tokens number %d and pools number %d do not match", len(sourceCCIP.Common.BridgeTokens), len(sourceCCIP.Common.BridgeTokenPools))
	}
	for i, token := range sourceCCIP.Common.BridgeTokens {
		tokens = append(tokens, token.ContractAddress)
		pools = append(pools, sourceCCIP.Common.BridgeTokenPools[i].EthAddress)
		conf := evm_2_evm_onramp.EVM2EVMOnRampTokenTransferFeeConfigArgs{
			Token:                     token.ContractAddress,
			MinFeeUSDCents:            50,           // $0.5
			MaxFeeUSDCents:            1_000_000_00, // $ 1 million
			DeciBps:                   5_0,          // 5 bps
			AggregateRateLimitEnabled: enableAggregateRateLimit,
		}
		if sourceCCIP.Common.BridgeTokenPools[i].IsUSDC() {
			conf.DestBytesOverhead = defaultUSDCDestBytesOverhead
			conf.DestGasOverhead = defaultUSDCDestGasOverhead
		}
		tokenTransferFeeConfig = append(tokenTransferFeeConfig, conf)
	}
	err := sourceCCIP.OnRamp.SetTokenTransferFeeConfig(tokenTransferFeeConfig)
	if err != nil {
		return fmt.Errorf("setting token transfer fee config shouldn't fail %w", err)
	}
	// this is required for v1.2.0 ramps
	err = sourceCCIP.OnRamp.ApplyPoolUpdates(tokens, pools)
	if err != nil {
		return fmt.Errorf("applying pool updates shouldn't fail %w", err)
	}
	return nil
}

// DeployContracts deploys all CCIP contracts specific to the source chain
func (sourceCCIP *SourceCCIPModule) DeployContracts(lane *laneconfig.LaneConfig) error {
	var err error
	contractDeployer := sourceCCIP.Common.Deployer
	log.Info().Msg("Deploying source chain specific contracts")

	sourceCCIP.LoadContracts(lane)
	sourceChainSelector, err := chainselectors.SelectorFromChainId(sourceCCIP.Common.ChainClient.GetChainID().Uint64())
	if err != nil {
		return fmt.Errorf("getting chain selector shouldn't fail %w", err)
	}

	if sourceCCIP.OnRamp == nil {
		if sourceCCIP.Common.ExistingDeployment {
			return fmt.Errorf("existing deployment is set to true but no onramp address is provided")
		}
		var tokensAndPools []evm_2_evm_onramp_1_2_0.InternalPoolUpdate
		var tokenTransferFeeConfig []evm_2_evm_onramp.EVM2EVMOnRampTokenTransferFeeConfigArgs

		sourceCCIP.SrcStartBlock, err = sourceCCIP.Common.ChainClient.LatestBlockNumber(context.Background())
		if err != nil {
			return fmt.Errorf("getting latest block number shouldn't fail %w", err)
		}
		var tokenAdminReg common.Address
		if contracts.NeedTokenAdminRegistry() {
			if sourceCCIP.Common.TokenAdminRegistry == nil {
				return fmt.Errorf("token admin registry contract address is not provided in lane config")
			}
			tokenAdminReg = sourceCCIP.Common.TokenAdminRegistry.EthAddress
		}
		sourceCCIP.OnRamp, err = contractDeployer.DeployOnRamp(
			sourceChainSelector,
			sourceCCIP.DestChainSelector,
			tokensAndPools,
			*sourceCCIP.Common.RMNContract,
			sourceCCIP.Common.Router.EthAddress,
			sourceCCIP.Common.PriceRegistry.EthAddress,
			tokenAdminReg,
			sourceCCIP.Common.RateLimiterConfig,
			[]evm_2_evm_onramp.EVM2EVMOnRampFeeTokenConfigArgs{
				{
					Token:                      common.HexToAddress(sourceCCIP.Common.FeeToken.Address()),
					NetworkFeeUSDCents:         1_00,
					GasMultiplierWeiPerEth:     GasFeeMultiplier,
					PremiumMultiplierWeiPerEth: 1e18,
					Enabled:                    true,
				},
				{
					Token:                      sourceCCIP.Common.WrappedNative,
					NetworkFeeUSDCents:         1_00,
					GasMultiplierWeiPerEth:     GasFeeMultiplier,
					PremiumMultiplierWeiPerEth: 1e18,
					Enabled:                    true,
				},
			}, tokenTransferFeeConfig, sourceCCIP.Common.FeeToken.EthAddress)

		if err != nil {
			return fmt.Errorf("onRamp deployment shouldn't fail %w", err)
		}

		err = sourceCCIP.Common.ChainClient.WaitForEvents()
		if err != nil {
			return fmt.Errorf("waiting for onRamp deployment shouldn't fail %w", err)
		}

		// update source Router with OnRamp address
		err = sourceCCIP.Common.Router.SetOnRamp(sourceCCIP.DestChainSelector, sourceCCIP.OnRamp.EthAddress)
		if err != nil {
			return fmt.Errorf("setting onramp on the router shouldn't fail %w", err)
		}
		// now sync the pools and tokens
		err := sourceCCIP.SetAllTokenTransferFeeConfigs(true)
		if err != nil {
			return err
		}
	} else {
		sourceCCIP.OnRamp, err = contractDeployer.NewOnRamp(sourceCCIP.OnRamp.EthAddress)
		if err != nil {
			return fmt.Errorf("getting new onramp contractshouldn't fail %w", err)
		}
	}
	return nil
}

func (sourceCCIP *SourceCCIPModule) CollectBalanceRequirements() []testhelpers.BalanceReq {
	var balancesReq []testhelpers.BalanceReq
	for _, token := range sourceCCIP.Common.BridgeTokens {
		balancesReq = append(balancesReq, testhelpers.BalanceReq{
			Name:   fmt.Sprintf("BridgeToken-%s-Address-%s", token.Address(), sourceCCIP.Sender.Hex()),
			Addr:   sourceCCIP.Sender,
			Getter: GetterForLinkToken(token.BalanceOf, sourceCCIP.Sender.Hex()),
		})
	}
	for i, pool := range sourceCCIP.Common.BridgeTokenPools {
		balancesReq = append(balancesReq, testhelpers.BalanceReq{
			Name:   fmt.Sprintf("BridgeToken-%s-TokenPool-%s", sourceCCIP.Common.BridgeTokens[i].Address(), pool.Address()),
			Addr:   pool.EthAddress,
			Getter: GetterForLinkToken(sourceCCIP.Common.BridgeTokens[i].BalanceOf, pool.Address()),
		})
	}

	if sourceCCIP.Common.FeeToken.Address() != common.HexToAddress("0x0").String() {
		balancesReq = append(balancesReq, testhelpers.BalanceReq{
			Name:   fmt.Sprintf("FeeToken-%s-Address-%s", sourceCCIP.Common.FeeToken.Address(), sourceCCIP.Sender.Hex()),
			Addr:   sourceCCIP.Sender,
			Getter: GetterForLinkToken(sourceCCIP.Common.FeeToken.BalanceOf, sourceCCIP.Sender.Hex()),
		})
		balancesReq = append(balancesReq, testhelpers.BalanceReq{
			Name:   fmt.Sprintf("FeeToken-%s-Router-%s", sourceCCIP.Common.FeeToken.Address(), sourceCCIP.Common.Router.Address()),
			Addr:   sourceCCIP.Common.Router.EthAddress,
			Getter: GetterForLinkToken(sourceCCIP.Common.FeeToken.BalanceOf, sourceCCIP.Common.Router.Address()),
		})
		balancesReq = append(balancesReq, testhelpers.BalanceReq{
			Name:   fmt.Sprintf("FeeToken-%s-OnRamp-%s", sourceCCIP.Common.FeeToken.Address(), sourceCCIP.OnRamp.Address()),
			Addr:   sourceCCIP.OnRamp.EthAddress,
			Getter: GetterForLinkToken(sourceCCIP.Common.FeeToken.BalanceOf, sourceCCIP.OnRamp.Address()),
		})
		balancesReq = append(balancesReq, testhelpers.BalanceReq{
			Name:   fmt.Sprintf("FeeToken-%s-Prices-%s", sourceCCIP.Common.FeeToken.Address(), sourceCCIP.Common.PriceRegistry.Address()),
			Addr:   sourceCCIP.Common.PriceRegistry.EthAddress,
			Getter: GetterForLinkToken(sourceCCIP.Common.FeeToken.BalanceOf, sourceCCIP.Common.PriceRegistry.Address()),
		})
	}
	return balancesReq
}

func (sourceCCIP *SourceCCIPModule) UpdateBalance(
	noOfReq int64,
	totalFee *big.Int,
	balances *BalanceSheet,
) {
	if len(sourceCCIP.TransferAmount) > 0 {
		for i := range sourceCCIP.TransferAmount {
			if sourceCCIP.TransferAmount[i] == nil { // nil transfer amount means no transfer for this token
				continue
			}
			// if length of sourceCCIP.TransferAmount is more than available bridge token use first bridge token
			token := sourceCCIP.Common.BridgeTokens[0]
			if i < len(sourceCCIP.Common.BridgeTokens) {
				token = sourceCCIP.Common.BridgeTokens[i]
			}
			name := fmt.Sprintf("BridgeToken-%s-Address-%s", token.Address(), sourceCCIP.Sender.Hex())
			balances.Update(name, BalanceItem{
				Address:  sourceCCIP.Sender,
				Getter:   GetterForLinkToken(token.BalanceOf, sourceCCIP.Sender.Hex()),
				AmtToSub: bigmath.Mul(big.NewInt(noOfReq), sourceCCIP.TransferAmount[i]),
			})
		}
		for i := range sourceCCIP.TransferAmount {
			// if length of sourceCCIP.TransferAmount is more than available bridge token use first bridge token
			pool := sourceCCIP.Common.BridgeTokenPools[0]
			index := 0
			if i < len(sourceCCIP.Common.BridgeTokenPools) {
				pool = sourceCCIP.Common.BridgeTokenPools[i]
				index = i
			}

			name := fmt.Sprintf("BridgeToken-%s-TokenPool-%s", sourceCCIP.Common.BridgeTokens[index].Address(), pool.Address())
			balances.Update(name, BalanceItem{
				Address:  pool.EthAddress,
				Getter:   GetterForLinkToken(sourceCCIP.Common.BridgeTokens[index].BalanceOf, pool.Address()),
				AmtToAdd: bigmath.Mul(big.NewInt(noOfReq), sourceCCIP.TransferAmount[i]),
			})
		}
	}
	if sourceCCIP.Common.FeeToken.Address() != common.HexToAddress("0x0").String() {
		name := fmt.Sprintf("FeeToken-%s-Address-%s", sourceCCIP.Common.FeeToken.Address(), sourceCCIP.Sender.Hex())
		balances.Update(name, BalanceItem{
			Address:  sourceCCIP.Sender,
			Getter:   GetterForLinkToken(sourceCCIP.Common.FeeToken.BalanceOf, sourceCCIP.Sender.Hex()),
			AmtToSub: totalFee,
		})
		name = fmt.Sprintf("FeeToken-%s-Prices-%s", sourceCCIP.Common.FeeToken.Address(), sourceCCIP.Common.PriceRegistry.Address())
		balances.Update(name, BalanceItem{
			Address: sourceCCIP.Common.PriceRegistry.EthAddress,
			Getter:  GetterForLinkToken(sourceCCIP.Common.FeeToken.BalanceOf, sourceCCIP.Common.PriceRegistry.Address()),
		})
		name = fmt.Sprintf("FeeToken-%s-Router-%s", sourceCCIP.Common.FeeToken.Address(), sourceCCIP.Common.Router.Address())
		balances.Update(name, BalanceItem{
			Address: sourceCCIP.Common.Router.EthAddress,
			Getter:  GetterForLinkToken(sourceCCIP.Common.FeeToken.BalanceOf, sourceCCIP.Common.Router.Address()),
		})
		name = fmt.Sprintf("FeeToken-%s-OnRamp-%s", sourceCCIP.Common.FeeToken.Address(), sourceCCIP.OnRamp.Address())
		balances.Update(name, BalanceItem{
			Address:  sourceCCIP.OnRamp.EthAddress,
			Getter:   GetterForLinkToken(sourceCCIP.Common.FeeToken.BalanceOf, sourceCCIP.OnRamp.Address()),
			AmtToAdd: totalFee,
		})
	}
}

func (sourceCCIP *SourceCCIPModule) AssertSendRequestedLogFinalized(
	lggr *zerolog.Logger,
	txHash common.Hash,
	sendReqData []*contracts.SendReqEventData,
	prevEventAt time.Time,
	reqStats []*testreporters.RequestStat,
) (time.Time, uint64, error) {
	if len(sendReqData) != len(reqStats) {
		return time.Time{}, 0, fmt.Errorf("sendReqData and reqStats length mismatch")
	}
	var gasUsed uint64
	receipt, err := sourceCCIP.Common.ChainClient.GetTxReceipt(txHash)
	if err == nil {
		gasUsed = receipt.GasUsed
	}
	lggr.Info().Msg("Waiting for CCIPSendRequested event log to be finalized")
	finalizedBlockNum, finalizedAt, err := sourceCCIP.Common.ChainClient.WaitForFinalizedTx(txHash)
	if err != nil || finalizedBlockNum == nil {
		for i, stat := range reqStats {
			stat.UpdateState(lggr, stat.SeqNum, testreporters.SourceLogFinalized, time.Since(prevEventAt), testreporters.Failure, &testreporters.TransactionStats{
				MsgID:              fmt.Sprintf("0x%x", sendReqData[i].MessageId[:]),
				Fee:                sendReqData[i].Fee.String(),
				NoOfTokensSent:     sendReqData[i].NoOfTokens,
				MessageBytesLength: int64(sendReqData[i].DataLength),
				TxHash:             txHash.Hex(),
			})
		}
		return time.Time{}, 0, fmt.Errorf("error waiting for CCIPSendRequested event log to be finalized - %w", err)
	}
	for i, stat := range reqStats {
		stat.UpdateState(lggr, stat.SeqNum, testreporters.SourceLogFinalized, finalizedAt.Sub(prevEventAt), testreporters.Success,
			&testreporters.TransactionStats{
				MsgID:              fmt.Sprintf("0x%x", sendReqData[i].MessageId[:]),
				Fee:                sendReqData[i].Fee.String(),
				GasUsed:            gasUsed,
				NoOfTokensSent:     sendReqData[i].NoOfTokens,
				MessageBytesLength: int64(sendReqData[i].DataLength),
				TxHash:             txHash.Hex(),
				FinalizedByBlock:   finalizedBlockNum.String(),
				FinalizedAt:        finalizedAt.String(),
			})
	}
	return finalizedAt, finalizedBlockNum.Uint64(), nil
}

// IsRequestTriggeredWithinTimeframe monitors for live events occurring within the specified timeframe.
// Live events refer to those that are triggered after subscribing to the CCIP Send Requested event.
func (sourceCCIP *SourceCCIPModule) IsRequestTriggeredWithinTimeframe(timeframe *commonconfig.Duration) *time.Time {
	if timeframe == nil {
		return nil
	}
	var foundAt *time.Time
	lastSeenTimestamp := time.Now().UTC().Add(-timeframe.Duration())
	sourceCCIP.CCIPSendRequestedWatcher.Range(func(_, value any) bool {
		if sendRequestedEvents, exists := value.([]*evm_2_evm_onramp.EVM2EVMOnRampCCIPSendRequested); exists {
			for _, sendRequestedEvent := range sendRequestedEvents {
				raw := sendRequestedEvent.Raw
				hdr, err := sourceCCIP.Common.ChainClient.HeaderByNumber(context.Background(), new(big.Int).SetUint64(raw.BlockNumber))
				if err == nil {
					if hdr.Timestamp.After(lastSeenTimestamp) {
						foundAt = pointer.ToTime(hdr.Timestamp)
						return false
					}
				}
			}
		}
		return true
	})
	return foundAt
}

// IsPastRequestTriggeredWithinTimeframe determines the average block time and calculates the block numbers
// within the specified timeframe. It then uses FilterCCIPSendRequested to identify the past events.
func (sourceCCIP *SourceCCIPModule) IsPastRequestTriggeredWithinTimeframe(
	ctx context.Context,
	timeframe *commonconfig.Duration,
) (*types.Log, error) {
	if timeframe == nil {
		return nil, nil
	}
	//var foundAt *time.Time
	latestBlock, err := sourceCCIP.Common.ChainClient.LatestBlockNumber(ctx)
	if err != nil {
		return nil, fmt.Errorf("error while getting latest source block number. Error: %w", err)
	}
	avgBlockTime, err := sourceCCIP.Common.ChainClient.AvgBlockTime(ctx)
	if err != nil {
		return nil, fmt.Errorf("error while getting average source block time. Error: %w", err)
	}
	filterFromBlock := latestBlock - uint64(timeframe.Duration()/avgBlockTime)

	onRampContract, err := evm_2_evm_onramp.NewEVM2EVMOnRamp(common.HexToAddress(sourceCCIP.OnRamp.EthAddress.Hex()),
		sourceCCIP.Common.ChainClient.Backend())
	if err != nil {
		return nil, fmt.Errorf("error while on ramp contract. Error: %w", err)
	}
	iterator, err := onRampContract.FilterCCIPSendRequested(&bind.FilterOpts{
		Start: filterFromBlock,
	})
	if err != nil {
		return nil, fmt.Errorf("error while filtering CCIP send requested starting block number: %d. Error: %w", filterFromBlock, err)
	}
	defer func() {
		iterErr := iterator.Close()
		if iterErr != nil {
			sourceCCIP.Common.Logger.Error().Err(iterErr).Msg("Error closing iterator")
		}
	}()
	lastBlockNumber := uint64(0)
	var latestEvent *types.Log
	for iterator.Next() {
		blockNum := iterator.Event.Raw.BlockNumber
		if blockNum > lastBlockNumber {
			lastBlockNumber = blockNum
			latestEvent = &iterator.Event.Raw
		}
	}
	return latestEvent, nil
}

func (sourceCCIP *SourceCCIPModule) AssertEventCCIPSendRequested(
	lggr *zerolog.Logger,
	txHash string,
	timeout time.Duration,
	prevEventAt time.Time,
	reqStat []*testreporters.RequestStat,
) ([]*contracts.SendReqEventData, time.Time, error) {
	lggr.Info().Str("Timeout", timeout.String()).Msg("Waiting for CCIPSendRequested event")
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	timer := time.NewTimer(timeout)
	defer timer.Stop()
	resetTimer := 0
	for {
		select {
		case <-ticker.C:
			value, ok := sourceCCIP.CCIPSendRequestedWatcher.Load(txHash)
			if ok {
				// if sendrequested events are found, check if the number of events are same as the number of requests
				if sendRequestedEvents, exists := value.([]*contracts.SendReqEventData); exists && len(sendRequestedEvents) == len(reqStat) {
					// if the value is processed, delete it from the map
					sourceCCIP.CCIPSendRequestedWatcher.Delete(txHash)
					for i, sendRequestedEvent := range sendRequestedEvents {
						seqNum := sendRequestedEvent.SequenceNumber
						lggr = ptr.Ptr(lggr.With().
							Uint64("SequenceNumber", seqNum).
							Str("MsgID", fmt.Sprintf("0x%x", sendRequestedEvent.MessageId[:])).
							Logger())
						// prevEventAt is the time when the message was successful, this should be same as the time when the event was emitted
						reqStat[i].UpdateState(lggr, seqNum, testreporters.CCIPSendRe, 0, testreporters.Success, nil)
					}
					var err error
					if len(sendRequestedEvents) == 0 {
						err = fmt.Errorf("message logs not found, no CCIPSendRequested event found for tx %s", txHash)
					}
					return sendRequestedEvents, prevEventAt, err
				}
			}
		case <-timer.C:
			// if there is connection issue reset the timer :
			if sourceCCIP.Common.IsConnectionRestoredRecently != nil && !sourceCCIP.Common.IsConnectionRestoredRecently.Load() {
				if resetTimer > 2 {
					for _, stat := range reqStat {
						stat.UpdateState(lggr, 0, testreporters.CCIPSendRe, time.Since(prevEventAt), testreporters.Failure,
							&testreporters.TransactionStats{
								TxHash: txHash,
							})
					}
					return nil, time.Now(), fmt.Errorf("possible RPC issue - CCIPSendRequested event is not found for tx %s", txHash)
				}
				resetTimer++
				timer.Reset(timeout)
				lggr.Info().Int("count of reset", resetTimer).Msg("Resetting timer to validate CCIPSendRequested event")
				continue
			}
			for _, stat := range reqStat {
				stat.UpdateState(lggr, 0, testreporters.CCIPSendRe, time.Since(prevEventAt), testreporters.Failure,
					&testreporters.TransactionStats{
						TxHash: txHash,
					})
			}
			return nil, time.Now(), fmt.Errorf("CCIPSendRequested event is not found for tx %s", txHash)
		}
	}
}

// CCIPMsg constructs the message for a CCIP request
func (sourceCCIP *SourceCCIPModule) CCIPMsg(
	receiver common.Address,
	allowOutOfOrder bool,
	gasLimit *big.Int,
) (router.ClientEVM2AnyMessage, error) {
	length := sourceCCIP.MsgDataLength
	var data string
	if length > 0 {
		b := make([]byte, length)
		_, err := crypto_rand.Read(b)
		if err != nil {
			return router.ClientEVM2AnyMessage{}, fmt.Errorf("failed generating random string: %w", err)
		}
		randomString := base64.URLEncoding.EncodeToString(b)
		data = randomString[:length]
	}

	tokenAndAmounts := []router.ClientEVMTokenAmount{}
	for i, amount := range sourceCCIP.TransferAmount {
		if amount == nil { // make nil transfer amount 0 to avoid panics
			sourceCCIP.TransferAmount[i] = big.NewInt(0)
		}
		token := sourceCCIP.Common.BridgeTokens[0]
		// if length of sourceCCIP.TransferAmount is more than available bridge token use first bridge token
		if i < len(sourceCCIP.Common.BridgeTokens) {
			token = sourceCCIP.Common.BridgeTokens[i]
		}
		if amount == nil || amount.Cmp(big.NewInt(0)) == 0 {
			log.Warn().
				Str("Token Address", token.Address()).
				Int("Token Index", i).
				Msg("Not sending a request for token transfer as the amount is 0 or nil")
			continue
		}
		tokenAndAmounts = append(tokenAndAmounts, router.ClientEVMTokenAmount{
			Token: common.HexToAddress(token.Address()), Amount: amount,
		})
	}

	receiverAddr, err := utils.ABIEncode(`[{"type":"address"}]`, receiver)
	if err != nil {
		return router.ClientEVM2AnyMessage{}, fmt.Errorf("failed encoding the receiver address: %w", err)
	}

	var extraArgs []byte
	matchErr := contracts.MatchContractVersionsOrAbove(map[contracts.Name]contracts.Version{
		contracts.OnRampContract: contracts.V1_5_0,
	})
	if matchErr != nil {
		extraArgs, err = testhelpers.GetEVMExtraArgsV1(gasLimit, false)
	} else {
		extraArgs, err = testhelpers.GetEVMExtraArgsV2(gasLimit, allowOutOfOrder)
	}
	if err != nil {
		return router.ClientEVM2AnyMessage{}, fmt.Errorf("failed getting extra args: %w", err)
	}
	// form the message for transfer
	return router.ClientEVM2AnyMessage{
		Receiver:     receiverAddr,
		Data:         []byte(data),
		TokenAmounts: tokenAndAmounts,
		FeeToken:     common.HexToAddress(sourceCCIP.Common.FeeToken.Address()),
		ExtraArgs:    extraArgs,
	}, nil
}

// SendRequest sends a CCIP request to the source chain's router contract
func (sourceCCIP *SourceCCIPModule) SendRequest(
	receiver common.Address,
	gasLimit *big.Int,
) (common.Hash, time.Duration, *big.Int, error) {
	var d time.Duration
	destChainSelector, err := chainselectors.SelectorFromChainId(sourceCCIP.DestinationChainId)
	if err != nil {
		return common.Hash{}, d, nil, fmt.Errorf("failed getting the chain selector: %w", err)
	}
	// form the message for transfer
	msg, err := sourceCCIP.CCIPMsg(receiver, sourceCCIP.Common.AllowOutOfOrder, gasLimit)
	if err != nil {
		return common.Hash{}, d, nil, fmt.Errorf("failed forming the ccip msg: %w", err)
	}

	fee, err := sourceCCIP.Common.Router.GetFee(destChainSelector, msg)
	if err != nil {
		log.Info().Interface("Msg", msg).Msg("CCIP msg")
		reason, _ := blockchain.RPCErrorFromError(err)
		if reason != "" {
			return common.Hash{}, d, nil, fmt.Errorf("failed getting the fee: %s", reason)
		}
		return common.Hash{}, d, nil, fmt.Errorf("failed getting the fee: %w", err)
	}
	log.Info().Str("Fee", fee.String()).Msg("Calculated fee")

	var sendTx *types.Transaction
	timeNow := time.Now()
	feeToken := common.HexToAddress(sourceCCIP.Common.FeeToken.Address())
	// initiate the transfer
	// if the fee token address is 0x0 it will use Native as fee token and the fee amount should be mentioned in bind.TransactOpts's value
	if feeToken != (common.Address{}) {
		sendTx, err = sourceCCIP.Common.Router.CCIPSendAndProcessTx(destChainSelector, msg, nil)
		if err != nil {
			txHash := common.Hash{}
			if sendTx != nil {
				txHash = sendTx.Hash()
			}
			return txHash, time.Since(timeNow), nil, fmt.Errorf("failed initiating the transfer ccip-send: %w", err)
		}
	} else {
		sendTx, err = sourceCCIP.Common.Router.CCIPSendAndProcessTx(destChainSelector, msg, fee)
		if err != nil {
			txHash := common.Hash{}
			if sendTx != nil {
				txHash = sendTx.Hash()
			}
			return txHash, time.Since(timeNow), nil, fmt.Errorf("failed initiating the transfer ccip-send: %w", err)
		}
	}

	log.Info().
		Str("Network", sourceCCIP.Common.ChainClient.GetNetworkName()).
		Str("Send token transaction", sendTx.Hash().String()).
		Str("lane", fmt.Sprintf("%s-->%s", sourceCCIP.Common.ChainClient.GetNetworkName(), sourceCCIP.DestNetworkName)).
		Msg("Sending token")
	return sendTx.Hash(), time.Since(timeNow), fee, nil
}

func DefaultSourceCCIPModule(
	logger *zerolog.Logger,
	testConf *testconfig.CCIPTestGroupConfig,
	chainClient blockchain.EVMClient,
	destChainId uint64,
	destChain string,
	laneConf *laneconfig.LaneConfig,
) (*SourceCCIPModule, error) {
	cmn, err := NewCCIPCommonFromConfig(
		logger, testConf, chainClient, laneConf,
	)
	if err != nil {
		return nil, err
	}

	destChainSelector, err := chainselectors.SelectorFromChainId(destChainId)
	if err != nil {
		return nil, fmt.Errorf("failed getting the chain selector: %w", err)
	}
	source := &SourceCCIPModule{
		Common:                   cmn,
		TransferAmount:           testConf.MsgDetails.TransferAmounts(),
		MsgDataLength:            pointer.GetInt64(testConf.MsgDetails.DataLength),
		DestinationChainId:       destChainId,
		DestChainSelector:        destChainSelector,
		DestNetworkName:          destChain,
		Sender:                   common.HexToAddress(chainClient.GetDefaultWallet().Address()),
		CCIPSendRequestedWatcher: &sync.Map{},
	}

	return source, nil
}

type DestCCIPModule struct {
	Common                  *CCIPCommon
	SourceChainId           uint64
	SourceChainSelector     uint64
	SourceNetworkName       string
	CommitStore             *contracts.CommitStore
	ReceiverDapp            *contracts.ReceiverDapp
	OffRamp                 *contracts.OffRamp
	ReportAcceptedWatcher   *sync.Map
	ExecStateChangedWatcher *sync.Map
	ReportBlessedWatcher    *sync.Map
	ReportBlessedBySeqNum   *sync.Map
	NextSeqNumToCommit      *atomic.Uint64
	DestStartBlock          uint64
}

func (destCCIP *DestCCIPModule) LoadContracts(conf *laneconfig.LaneConfig) {
	if conf != nil {
		cfg, ok := conf.DestContracts[destCCIP.SourceNetworkName]
		if ok {
			if common.IsHexAddress(cfg.OffRamp) {
				destCCIP.OffRamp = &contracts.OffRamp{
					EthAddress: common.HexToAddress(cfg.OffRamp),
				}
			}
			if common.IsHexAddress(cfg.CommitStore) {
				destCCIP.CommitStore = &contracts.CommitStore{
					EthAddress: common.HexToAddress(cfg.CommitStore),
				}
			}
			if common.IsHexAddress(cfg.ReceiverDapp) {
				destCCIP.ReceiverDapp = &contracts.ReceiverDapp{
					EthAddress: common.HexToAddress(cfg.ReceiverDapp),
				}
			}
		}
	}
}

func (destCCIP *DestCCIPModule) SyncTokensAndPools(srcTokens []*contracts.ERC20Token) error {
	if destCCIP.OffRamp.Instance.V1_2_0 == nil {
		return nil
	}
	var sourceTokens, pools []common.Address

	for _, token := range srcTokens {
		sourceTokens = append(sourceTokens, common.HexToAddress(token.Address()))
	}

	for i := range destCCIP.Common.BridgeTokenPools {
		pools = append(pools, destCCIP.Common.BridgeTokenPools[i].EthAddress)
	}
	if len(sourceTokens) != len(pools) {
		return fmt.Errorf("source token and destination pool length mismatch")
	}
	// if number of tokens are more than 10, then we need to split the tokens in batch of 10 and call sync
	// otherwise the tx gets too large and we will get out of gas error
	if len(sourceTokens) > 10 {
		for i := 0; i < len(sourceTokens); i += 10 {
			end := i + 10
			if end > len(sourceTokens) {
				end = len(sourceTokens)
			}
			err := destCCIP.OffRamp.SyncTokensAndPools(sourceTokens[i:end], pools[i:end])
			if err != nil {
				return err
			}
		}
		return nil
	}
	return destCCIP.OffRamp.SyncTokensAndPools(sourceTokens, pools)
}

// AddRateLimitTokens adds token pairs to the OffRamp's rate limiting
func (destCCIP *DestCCIPModule) AddRateLimitTokens(srcTokens, destTokens []*contracts.ERC20Token) error {
	if destCCIP.OffRamp.Instance.Latest == nil {
		return nil
	}
	if srcTokens == nil || destTokens == nil {
		return fmt.Errorf("source or destination tokens are nil")
	}

	if len(srcTokens) != len(destTokens) {
		return fmt.Errorf("source and destination token length mismatch")
	}

	var sourceTokenAddresses, destTokenAddresses []common.Address

	for i, token := range srcTokens {
		sourceTokenAddresses = append(sourceTokenAddresses, common.HexToAddress(token.Address()))
		destTokenAddresses = append(destTokenAddresses, common.HexToAddress(destTokens[i].Address()))
	}

	// if number of tokens are more than 10, then we need to split the tokens in batch of 10 and update the rate limit
	// otherwise the tx gets too large and we will get out of gas error
	if len(sourceTokenAddresses) > 10 {
		for i := 0; i < len(sourceTokenAddresses); i += 10 {
			end := i + 10
			if end > len(sourceTokenAddresses) {
				end = len(sourceTokenAddresses)
			}
			err := destCCIP.OffRamp.AddRateLimitTokens(sourceTokenAddresses[i:end], destTokenAddresses[i:end])
			if err != nil {
				return err
			}
		}
		return nil
	}

	return destCCIP.OffRamp.AddRateLimitTokens(sourceTokenAddresses, destTokenAddresses)
}

// RemoveRateLimitTokens removes token pairs from the OffRamp's rate limiting.
// If you ask to remove a token pair that doesn't exist, it will return an error.
func (destCCIP *DestCCIPModule) RemoveRateLimitTokens(ctx context.Context, srcTokens, destTokens []*contracts.ERC20Token) error {
	if srcTokens == nil || destTokens == nil {
		return fmt.Errorf("source or destination tokens are nil")
	}

	if len(srcTokens) != len(destTokens) {
		return fmt.Errorf("source and destination token length mismatch")
	}

	var sourceTokenAddresses, destTokenAddresses []common.Address

	for i, token := range srcTokens {
		sourceTokenAddresses = append(sourceTokenAddresses, common.HexToAddress(token.Address()))
		destTokenAddresses = append(destTokenAddresses, common.HexToAddress(destTokens[i].Address()))
	}

	return destCCIP.OffRamp.RemoveRateLimitTokens(ctx, sourceTokenAddresses, destTokenAddresses)
}

// RemoveAllRateLimitTokens removes all token pairs from the OffRamp's rate limiting.
func (destCCIP *DestCCIPModule) RemoveAllRateLimitTokens(ctx context.Context) error {
	return destCCIP.OffRamp.RemoveAllRateLimitTokens(ctx)
}

// DeployContracts deploys all CCIP contracts specific to the destination chain
func (destCCIP *DestCCIPModule) DeployContracts(
	sourceCCIP SourceCCIPModule,
	lane *laneconfig.LaneConfig,
) error {
	var err error
	contractDeployer := destCCIP.Common.Deployer
	log.Info().Msg("Deploying destination chain specific contracts")
	destCCIP.LoadContracts(lane)
	destChainSelector, err := chainselectors.SelectorFromChainId(destCCIP.Common.ChainClient.GetChainID().Uint64())
	if err != nil {
		return fmt.Errorf("failed to get chain selector for destination chain id %d: %w", destCCIP.Common.ChainClient.GetChainID().Uint64(), err)
	}
	destCCIP.DestStartBlock, err = destCCIP.Common.ChainClient.LatestBlockNumber(context.Background())
	if err != nil {
		return fmt.Errorf("getting latest block number shouldn't fail %w", err)
	}
	if !destCCIP.Common.ExistingDeployment && len(sourceCCIP.Common.BridgeTokenPools) != len(destCCIP.Common.BridgeTokenPools) {
		return fmt.Errorf("source and destination token pool number does not match")
	}

	if destCCIP.CommitStore == nil {
		if destCCIP.Common.ExistingDeployment {
			return fmt.Errorf("commit store address not provided in lane config")
		}
		// commitStore responsible for validating the transfer message
		destCCIP.CommitStore, err = contractDeployer.DeployCommitStore(
			destCCIP.SourceChainSelector,
			destChainSelector,
			sourceCCIP.OnRamp.EthAddress,
			*destCCIP.Common.RMNContract,
		)
		if err != nil {
			return fmt.Errorf("deploying commitstore shouldn't fail %w", err)
		}
		err = destCCIP.Common.ChainClient.WaitForEvents()
		if err != nil {
			return fmt.Errorf("waiting for commitstore deployment shouldn't fail %w", err)
		}

		// CommitStore can update
		err = destCCIP.Common.PriceRegistry.AddPriceUpdater(destCCIP.CommitStore.EthAddress)
		if err != nil {
			return fmt.Errorf("setting commitstore as fee updater shouldn't fail %w", err)
		}
		err = destCCIP.Common.ChainClient.WaitForEvents()
		if err != nil {
			return fmt.Errorf("waiting for setting commitstore as fee updater shouldn't fail %w", err)
		}
	} else {
		destCCIP.CommitStore, err = contractDeployer.NewCommitStore(destCCIP.CommitStore.EthAddress)
		if err != nil {
			return fmt.Errorf("getting new commitstore shouldn't fail %w", err)
		}
	}

	if destCCIP.OffRamp == nil {
		if destCCIP.Common.ExistingDeployment {
			return fmt.Errorf("offramp address not provided in lane config")
		}
		var tokenAdminReg common.Address
		if contracts.NeedTokenAdminRegistry() {
			if destCCIP.Common.TokenAdminRegistry == nil {
				return fmt.Errorf("token admin registry contract address is not provided in lane config")
			}
			tokenAdminReg = destCCIP.Common.TokenAdminRegistry.EthAddress
		}
		destCCIP.OffRamp, err = contractDeployer.DeployOffRamp(
			destCCIP.SourceChainSelector,
			destChainSelector,
			destCCIP.CommitStore.EthAddress,
			sourceCCIP.OnRamp.EthAddress,
			destCCIP.Common.RateLimiterConfig,
			[]common.Address{},
			[]common.Address{},
			*destCCIP.Common.RMNContract,
			tokenAdminReg,
		)
		if err != nil {
			return fmt.Errorf("deploying offramp shouldn't fail %w", err)
		}
		err = destCCIP.Common.ChainClient.WaitForEvents()
		if err != nil {
			return fmt.Errorf("waiting for offramp deployment shouldn't fail %w", err)
		}

		// apply offramp updates
		_, err = destCCIP.Common.Router.AddOffRamp(destCCIP.OffRamp.EthAddress, destCCIP.SourceChainSelector)
		if err != nil {
			return fmt.Errorf("setting offramp as fee updater shouldn't fail %w", err)
		}

		err = destCCIP.AddRateLimitTokens(sourceCCIP.Common.BridgeTokens, destCCIP.Common.BridgeTokens)
		if err != nil {
			return fmt.Errorf("setting rate limited tokens shouldn't fail %w", err)
		}
		err = destCCIP.SyncTokensAndPools(sourceCCIP.Common.BridgeTokens)
		if err != nil {
			return fmt.Errorf("syncing tokens and pools shouldn't fail %w", err)
		}
		err = destCCIP.Common.ChainClient.WaitForEvents()
		if err != nil {
			return fmt.Errorf("waiting for events on destination contract shouldn't fail %w", err)
		}
	} else {
		destCCIP.OffRamp, err = contractDeployer.NewOffRamp(destCCIP.OffRamp.EthAddress)
		if err != nil {
			return fmt.Errorf("getting new offramp shouldn't fail %w", err)
		}
	}
	if destCCIP.ReceiverDapp == nil {
		// ReceiverDapp
		destCCIP.ReceiverDapp, err = contractDeployer.DeployReceiverDapp(false)
		if err != nil {
			return fmt.Errorf("receiverDapp contract should be deployed successfully %w", err)
		}
		err = destCCIP.Common.ChainClient.WaitForEvents()
		if err != nil {
			return fmt.Errorf("waiting for events on destination contract deployments %w", err)
		}
	} else {
		destCCIP.ReceiverDapp, err = contractDeployer.NewReceiverDapp(destCCIP.ReceiverDapp.EthAddress)
		if err != nil {
			return fmt.Errorf("getting new receiverDapp shouldn't fail %w", err)
		}
	}
	return nil
}

func (destCCIP *DestCCIPModule) CollectBalanceRequirements() []testhelpers.BalanceReq {
	var destBalancesReq []testhelpers.BalanceReq
	for _, token := range destCCIP.Common.BridgeTokens {
		destBalancesReq = append(destBalancesReq, testhelpers.BalanceReq{
			Name:   fmt.Sprintf("BridgeToken-%s-Address-%s", token.Address(), destCCIP.ReceiverDapp.Address()),
			Addr:   destCCIP.ReceiverDapp.EthAddress,
			Getter: GetterForLinkToken(token.BalanceOf, destCCIP.ReceiverDapp.Address()),
		})
	}
	for i, pool := range destCCIP.Common.BridgeTokenPools {
		destBalancesReq = append(destBalancesReq, testhelpers.BalanceReq{
			Name:   fmt.Sprintf("BridgeToken-%s-TokenPool-%s", destCCIP.Common.BridgeTokens[i].Address(), pool.Address()),
			Addr:   pool.EthAddress,
			Getter: GetterForLinkToken(destCCIP.Common.BridgeTokens[i].BalanceOf, pool.Address()),
		})
	}
	if destCCIP.Common.FeeToken.Address() != common.HexToAddress("0x0").String() {
		destBalancesReq = append(destBalancesReq, testhelpers.BalanceReq{
			Name:   fmt.Sprintf("FeeToken-%s-Address-%s", destCCIP.Common.FeeToken.Address(), destCCIP.ReceiverDapp.Address()),
			Addr:   destCCIP.ReceiverDapp.EthAddress,
			Getter: GetterForLinkToken(destCCIP.Common.FeeToken.BalanceOf, destCCIP.ReceiverDapp.Address()),
		})
		destBalancesReq = append(destBalancesReq, testhelpers.BalanceReq{
			Name:   fmt.Sprintf("FeeToken-%s-OffRamp-%s", destCCIP.Common.FeeToken.Address(), destCCIP.OffRamp.Address()),
			Addr:   destCCIP.OffRamp.EthAddress,
			Getter: GetterForLinkToken(destCCIP.Common.FeeToken.BalanceOf, destCCIP.OffRamp.Address()),
		})
	}
	return destBalancesReq
}

func (destCCIP *DestCCIPModule) UpdateBalance(
	transferAmount []*big.Int,
	noOfReq int64,
	balance *BalanceSheet,
) {
	if len(transferAmount) > 0 {
		for i := range transferAmount {
			token := destCCIP.Common.BridgeTokens[0]
			if i < len(destCCIP.Common.BridgeTokens) {
				token = destCCIP.Common.BridgeTokens[i]
			}
			name := fmt.Sprintf("BridgeToken-%s-Address-%s", token.Address(), destCCIP.ReceiverDapp.Address())
			balance.Update(name, BalanceItem{
				Address:  destCCIP.ReceiverDapp.EthAddress,
				Getter:   GetterForLinkToken(token.BalanceOf, destCCIP.ReceiverDapp.Address()),
				AmtToAdd: bigmath.Mul(big.NewInt(noOfReq), transferAmount[i]),
			})
		}
		for i := range transferAmount {
			pool := destCCIP.Common.BridgeTokenPools[0]
			index := 0
			if i < len(destCCIP.Common.BridgeTokenPools) {
				pool = destCCIP.Common.BridgeTokenPools[i]
				index = i
			}
			name := fmt.Sprintf("BridgeToken-%s-TokenPool-%s", destCCIP.Common.BridgeTokens[index].Address(), pool.Address())
			balance.Update(name, BalanceItem{
				Address:  pool.EthAddress,
				Getter:   GetterForLinkToken(destCCIP.Common.BridgeTokens[index].BalanceOf, pool.Address()),
				AmtToSub: bigmath.Mul(big.NewInt(noOfReq), transferAmount[i]),
			})
		}
	}
	if destCCIP.Common.FeeToken.Address() != common.HexToAddress("0x0").String() {
		name := fmt.Sprintf("FeeToken-%s-OffRamp-%s", destCCIP.Common.FeeToken.Address(), destCCIP.OffRamp.Address())
		balance.Update(name, BalanceItem{
			Address: destCCIP.OffRamp.EthAddress,
			Getter:  GetterForLinkToken(destCCIP.Common.FeeToken.BalanceOf, destCCIP.OffRamp.Address()),
		})

		name = fmt.Sprintf("FeeToken-%s-Address-%s", destCCIP.Common.FeeToken.Address(), destCCIP.ReceiverDapp.Address())
		balance.Update(name, BalanceItem{
			Address: destCCIP.ReceiverDapp.EthAddress,
			Getter:  GetterForLinkToken(destCCIP.Common.FeeToken.BalanceOf, destCCIP.ReceiverDapp.Address()),
		})
	}
}

// AssertNoReportAcceptedEventReceived validates that no ExecutionStateChangedEvent is emitted for mentioned timeRange after lastSeenTimestamp
func (destCCIP *DestCCIPModule) AssertNoReportAcceptedEventReceived(lggr *zerolog.Logger, timeRange time.Duration, lastSeenTimestamp time.Time) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeRange)
	defer cancel()
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			var eventFoundAfterCursing *time.Time
			// verify if CommitReportAccepted is received, it's not generated after provided lastSeenTimestamp
			destCCIP.ReportAcceptedWatcher.Range(func(_, value any) bool {
				e, exists := value.(*evm_2_evm_offramp.EVM2EVMOffRampExecutionStateChanged)
				if exists {
					vLogs := e.Raw
					hdr, err := destCCIP.Common.ChainClient.HeaderByNumber(ctx, new(big.Int).SetUint64(vLogs.BlockNumber))
					if err != nil {
						return true
					}
					if hdr.Timestamp.After(lastSeenTimestamp) {
						eventFoundAfterCursing = pointer.ToTime(hdr.Timestamp)
						return false
					}
				}
				return true
			})
			if eventFoundAfterCursing != nil {
				return fmt.Errorf("CommitReportAccepted Event detected at %s after %s", lastSeenTimestamp, eventFoundAfterCursing.String())
			}
		case <-ctx.Done():
			lggr.Info().Msgf("successfully validated that no CommitReportAccepted detected after %s for %s", lastSeenTimestamp, timeRange)
			return nil
		}
	}
}

// AssertNoExecutionStateChangedEventReceived validates that no ExecutionStateChangedEvent is emitted for mentioned timeRange after lastSeenTimestamp
func (destCCIP *DestCCIPModule) AssertNoExecutionStateChangedEventReceived(
	lggr *zerolog.Logger,
	timeRange time.Duration,
	lastSeenTimestamp time.Time,
) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeRange)
	defer cancel()
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	lggr.Info().Str("Wait Time", timeRange.String()).Time("Since", lastSeenTimestamp).Msg("Waiting to ensure no ExecutionStateChanged event")
	for {
		select {
		case <-ticker.C:
			var eventFoundAfterCursing *time.Time
			// verify if ExecutionStateChanged is received, it's not generated after provided lastSeenTimestamp
			destCCIP.ExecStateChangedWatcher.Range(func(_, value any) bool {
				e, exists := value.(*contracts.EVM2EVMOffRampExecutionStateChanged)
				if exists {
					vLogs := e.LogInfo
					hdr, err := destCCIP.Common.ChainClient.HeaderByNumber(ctx, new(big.Int).SetUint64(vLogs.BlockNumber))
					if err != nil {
						return true
					}
					if hdr.Timestamp.After(lastSeenTimestamp) {
						eventFoundAfterCursing = pointer.ToTime(hdr.Timestamp)
						return false
					}
				}
				return true
			})
			if eventFoundAfterCursing != nil {
				return fmt.Errorf("ExecutionStateChanged Event detected at %s after %s", lastSeenTimestamp, eventFoundAfterCursing.String())
			}
		case <-ctx.Done():
			lggr.Info().Msgf("Successfully validated that no ExecutionStateChanged detected after %s for %s", lastSeenTimestamp, timeRange)
			return nil
		}
	}
}

func (destCCIP *DestCCIPModule) AssertEventExecutionStateChanged(
	lggr *zerolog.Logger,
	seqNum uint64,
	timeout time.Duration,
	timeNow time.Time,
	reqStat *testreporters.RequestStat,
	execState testhelpers.MessageExecutionState,
) (uint8, error) {
	lggr.Info().Uint64("seqNum", seqNum).Str("Timeout", timeout.String()).Msg("Waiting for ExecutionStateChanged event")
	timer := time.NewTimer(timeout)
	defer timer.Stop()
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	resetTimer := 0
	for {
		select {
		case <-ticker.C:
			value, ok := destCCIP.ExecStateChangedWatcher.Load(seqNum)
			if ok && value != nil {
				e, exists := value.(*contracts.EVM2EVMOffRampExecutionStateChanged)
				// find the type of the value
				if exists {
					// if the value is processed, delete it from the map
					destCCIP.ExecStateChangedWatcher.Delete(seqNum)
					vLogs := e.LogInfo
					receivedAt := time.Now().UTC()
					hdr, err := destCCIP.Common.ChainClient.HeaderByNumber(context.Background(), new(big.Int).SetUint64(vLogs.BlockNumber))
					if err == nil {
						receivedAt = hdr.Timestamp
					}
					receipt, err := destCCIP.Common.ChainClient.GetTxReceipt(vLogs.TxHash)
					if err != nil {
						lggr.Warn().Msg("Failed to get receipt for ExecStateChanged event")
					}
					var gasUsed uint64
					if receipt != nil {
						gasUsed = receipt.GasUsed
					}
					if testhelpers.MessageExecutionState(e.State) == execState {
						lggr.Info().Uint64("seqNum", seqNum).Uint8("ExecutionState", e.State).Msg("ExecutionStateChanged event received")
						reqStat.UpdateState(lggr, seqNum, testreporters.ExecStateChanged, receivedAt.Sub(timeNow),
							testreporters.Success,
							&testreporters.TransactionStats{
								TxHash:  vLogs.TxHash.Hex(),
								MsgID:   fmt.Sprintf("0x%x", e.MessageId[:]),
								GasUsed: gasUsed,
							},
						)
						return e.State, nil
					}
					reqStat.UpdateState(lggr, seqNum, testreporters.ExecStateChanged, time.Since(timeNow), testreporters.Failure, nil)
					return e.State, fmt.Errorf("ExecutionStateChanged event state - expected %d actual - %d with data %x for seq num %v for lane %d-->%d",
						execState, testhelpers.MessageExecutionState(e.State), e.ReturnData, seqNum, destCCIP.SourceChainId, destCCIP.Common.ChainClient.GetChainID())
				}
			}
		case <-timer.C:
			// if there is connection issue reset the context :
			if destCCIP.Common.IsConnectionRestoredRecently != nil && !destCCIP.Common.IsConnectionRestoredRecently.Load() {
				// if timer already has been reset 2 times we fail with warning
				if resetTimer > 2 {
					reqStat.UpdateState(lggr, seqNum, testreporters.ExecStateChanged, time.Since(timeNow), testreporters.Failure, nil)
					return 0, fmt.Errorf("possible RPC issues - ExecutionStateChanged event not found for seq num %d for lane %d-->%d",
						seqNum, destCCIP.SourceChainId, destCCIP.Common.ChainClient.GetChainID())
				}
				timer.Reset(timeout)
				resetTimer++
				lggr.Info().Int("count of reset", resetTimer).Msg("Resetting timer to validate ExecutionStateChanged event")
				continue
			}
			reqStat.UpdateState(lggr, seqNum, testreporters.ExecStateChanged, time.Since(timeNow), testreporters.Failure, nil)
			return 0, fmt.Errorf("ExecutionStateChanged event not found for seq num %d for lane %d-->%d",
				seqNum, destCCIP.SourceChainId, destCCIP.Common.ChainClient.GetChainID())
		}
	}
}

func (destCCIP *DestCCIPModule) AssertEventReportAccepted(
	lggr *zerolog.Logger,
	seqNum uint64,
	timeout time.Duration,
	prevEventAt time.Time,
	reqStat *testreporters.RequestStat,
) (*contracts.CommitStoreReportAccepted, time.Time, error) {
	lggr.Info().Uint64("seqNum", seqNum).Str("Timeout", timeout.String()).Msg("Waiting for ReportAccepted event")
	timer := time.NewTimer(timeout)
	defer timer.Stop()
	resetTimerCount := 0
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			value, ok := destCCIP.ReportAcceptedWatcher.Load(seqNum)
			if ok && value != nil {
				reportAccepted, exists := value.(*contracts.CommitStoreReportAccepted)
				if exists {
					// if the value is processed, delete it from the map
					destCCIP.ReportAcceptedWatcher.Delete(seqNum)
					receivedAt := time.Now().UTC()
					hdr, err := destCCIP.Common.ChainClient.HeaderByNumber(context.Background(), new(big.Int).SetUint64(reportAccepted.LogInfo.BlockNumber))
					if err == nil {
						receivedAt = hdr.Timestamp
					}

					totalTime := receivedAt.Sub(prevEventAt)
					// we cannot calculate the exact time at which block was finalized
					// as a result sometimes we get a time which is slightly after the block was marked as finalized
					// in such cases we get a negative time difference between finalized and report accepted if the commit
					// has happened almost immediately after block being finalized
					// in such cases we set the time difference to 1 second
					if totalTime < 0 {
						lggr.Warn().
							Uint64("seqNum", seqNum).
							Time("finalized at", prevEventAt).
							Time("ReportAccepted at", receivedAt).
							Msg("ReportAccepted event received before finalized timestamp")
						totalTime = time.Second
					}
					receipt, err := destCCIP.Common.ChainClient.GetTxReceipt(reportAccepted.LogInfo.TxHash)
					if err != nil {
						lggr.Warn().Msg("Failed to get receipt for ReportAccepted event")
					}
					var gasUsed uint64
					if receipt != nil {
						gasUsed = receipt.GasUsed
					}
					reqStat.UpdateState(lggr, seqNum, testreporters.Commit, totalTime, testreporters.Success,
						&testreporters.TransactionStats{
							GasUsed:    gasUsed,
							TxHash:     reportAccepted.LogInfo.TxHash.Hex(),
							CommitRoot: fmt.Sprintf("%x", reportAccepted.MerkleRoot),
						})
					return reportAccepted, receivedAt, nil
				}
			}
		case <-timer.C:
			// if there is connection issue reset the context :
			if destCCIP.Common.IsConnectionRestoredRecently != nil && !destCCIP.Common.IsConnectionRestoredRecently.Load() {
				if resetTimerCount > 2 {
					reqStat.UpdateState(lggr, seqNum, testreporters.Commit, time.Since(prevEventAt), testreporters.Failure, nil)
					return nil, time.Now().UTC(), fmt.Errorf("possible RPC issue - ReportAccepted is not found for seq num %d lane %d-->%d",
						seqNum, destCCIP.SourceChainId, destCCIP.Common.ChainClient.GetChainID())
				}
				timer.Reset(timeout)
				resetTimerCount++
				lggr.Info().Int("count of reset", resetTimerCount).Msg("Resetting timer to validate ReportAccepted event")
				continue
			}
			reqStat.UpdateState(lggr, seqNum, testreporters.Commit, time.Since(prevEventAt), testreporters.Failure, nil)
			return nil, time.Now().UTC(), fmt.Errorf("ReportAccepted is not found for seq num %d lane %d-->%d",
				seqNum, destCCIP.SourceChainId, destCCIP.Common.ChainClient.GetChainID())
		}
	}
}

func (destCCIP *DestCCIPModule) AssertReportBlessed(
	lggr *zerolog.Logger,
	seqNum uint64,
	timeout time.Duration,
	CommitReport contracts.CommitStoreReportAccepted,
	prevEventAt time.Time,
	reqStat *testreporters.RequestStat,
) (time.Time, error) {
	if destCCIP.Common.ARM == nil {
		lggr.Info().
			Uint64("commit store interval Min", CommitReport.Min).
			Uint64("commit store interval Max", CommitReport.Max).
			Hex("Root", CommitReport.MerkleRoot[:]).
			Msg("Skipping ReportBlessed check for mock ARM")
		return prevEventAt, nil
	}
	lggr.Info().
		Str("Timeout", timeout.String()).
		Uint64("commit store interval Min", CommitReport.Min).
		Uint64("commit store interval Max", CommitReport.Max).
		Msg("Waiting for Report To be blessed")
	timer := time.NewTimer(timeout)
	defer timer.Stop()
	resetTimerCount := 0
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			var value any
			var foundAsRoot, ok bool
			value, foundAsRoot = destCCIP.ReportBlessedWatcher.Load(CommitReport.MerkleRoot)
			receivedAt := time.Now().UTC()
			ok = foundAsRoot
			if !foundAsRoot {
				// if the value is not found as root, check if it is found as sequence number
				value, ok = destCCIP.ReportBlessedBySeqNum.Load(seqNum)
			}
			if ok && value != nil {
				vLogs, exists := value.(*contracts.LogInfo)
				if exists {
					// if the root is found, set the value for all the sequence numbers in the interval and delete the root from the map
					if foundAsRoot {
						// set the value for all the sequence numbers in the interval
						for i := CommitReport.Min; i <= CommitReport.Max; i++ {
							destCCIP.ReportBlessedBySeqNum.Store(i, vLogs)
						}
						// if the value is processed, delete it from the map
						destCCIP.ReportBlessedWatcher.Delete(CommitReport.MerkleRoot)
					} else {
						// if the value is processed, delete it from the map
						destCCIP.ReportBlessedBySeqNum.Delete(seqNum)
					}
					hdr, err := destCCIP.Common.ChainClient.HeaderByNumber(context.Background(), new(big.Int).SetUint64(vLogs.BlockNumber))
					if err == nil {
						receivedAt = hdr.Timestamp
					}
					receipt, err := destCCIP.Common.ChainClient.GetTxReceipt(vLogs.TxHash)
					if err != nil {
						lggr.Warn().Err(err).Msg("Failed to get receipt for ReportBlessed event")
					}
					var gasUsed uint64
					if receipt != nil {
						gasUsed = receipt.GasUsed
					}
					reqStat.UpdateState(lggr, seqNum, testreporters.ReportBlessed, receivedAt.Sub(prevEventAt), testreporters.Success,
						&testreporters.TransactionStats{
							GasUsed:    gasUsed,
							TxHash:     vLogs.TxHash.String(),
							CommitRoot: fmt.Sprintf("%x", CommitReport.MerkleRoot),
						})
					return receivedAt, nil
				}
			}
		case <-timer.C:
			// if there is connection issue reset the context :
			if destCCIP.Common.IsConnectionRestoredRecently != nil && !destCCIP.Common.IsConnectionRestoredRecently.Load() {
				if resetTimerCount > 2 {
					reqStat.UpdateState(lggr, seqNum, testreporters.ReportBlessed, time.Since(prevEventAt), testreporters.Failure, nil)
					return time.Now().UTC(), fmt.Errorf("possible RPC issue - ReportBlessed is not found for interval min - %d max - %d lane %d-->%d",
						CommitReport.Min, CommitReport.Max, destCCIP.SourceChainId, destCCIP.Common.ChainClient.GetChainID())
				}
				timer.Reset(timeout)
				resetTimerCount++
				lggr.Info().Int("count of reset", resetTimerCount).Msg("Resetting timer to validate ReportBlessed event")
				continue
			}
			reqStat.UpdateState(lggr, seqNum, testreporters.ReportBlessed, time.Since(prevEventAt), testreporters.Failure, nil)
			return time.Now().UTC(), fmt.Errorf("ReportBlessed is not found for interval min - %d max - %d lane %d-->%d",
				CommitReport.Min, CommitReport.Max, destCCIP.SourceChainId, destCCIP.Common.ChainClient.GetChainID())
		}
	}
}

func (destCCIP *DestCCIPModule) AssertSeqNumberExecuted(
	lggr *zerolog.Logger,
	seqNumberBefore uint64,
	timeout time.Duration,
	timeNow time.Time,
	reqStat *testreporters.RequestStat,
) error {
	lggr.Info().Uint64("seqNum", seqNumberBefore).Str("Timeout", timeout.String()).Msg("Waiting to be processed by commit store")
	timer := time.NewTimer(timeout)
	defer timer.Stop()
	resetTimerCount := 0
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if destCCIP.NextSeqNumToCommit.Load() > seqNumberBefore {
				return nil
			}
			seqNumberAfter, err := destCCIP.CommitStore.Instance.GetExpectedNextSequenceNumber(nil)
			if err != nil {
				// if we get error instead of returning error we continue, in case it's a temporary RPC failure .
				continue
			}
			if seqNumberAfter > seqNumberBefore {
				destCCIP.NextSeqNumToCommit.Store(seqNumberAfter)
				return nil
			}
		case <-timer.C:
			// if there is connection issue reset the context :
			if destCCIP.Common.IsConnectionRestoredRecently != nil && !destCCIP.Common.IsConnectionRestoredRecently.Load() {
				if resetTimerCount > 2 {
					reqStat.UpdateState(lggr, seqNumberBefore, testreporters.Commit, time.Since(timeNow), testreporters.Failure, nil)
					return fmt.Errorf("possible RPC issue - sequence number is not increased for seq num %d lane %d-->%d",
						seqNumberBefore, destCCIP.SourceChainId, destCCIP.Common.ChainClient.GetChainID())
				}
				timer.Reset(timeout)
				resetTimerCount++
				lggr.Info().Int("count of reset", resetTimerCount).Msg("Resetting timer to validate seqnumber increase in commit store")
				continue
			}
			reqStat.UpdateState(lggr, seqNumberBefore, testreporters.Commit, time.Since(timeNow), testreporters.Failure, nil)
			return fmt.Errorf("sequence number is not increased for seq num %d lane %d-->%d",
				seqNumberBefore, destCCIP.SourceChainId, destCCIP.Common.ChainClient.GetChainID())
		}
	}
}

func DefaultDestinationCCIPModule(
	logger *zerolog.Logger,
	testConf *testconfig.CCIPTestGroupConfig,
	chainClient blockchain.EVMClient,
	sourceChainId uint64,
	sourceChain string,
	laneConf *laneconfig.LaneConfig,
) (*DestCCIPModule, error) {
	cmn, err := NewCCIPCommonFromConfig(
		logger, testConf, chainClient, laneConf,
	)
	if err != nil {
		return nil, err
	}

	sourceChainSelector, err := chainselectors.SelectorFromChainId(sourceChainId)
	if err != nil {
		return nil, fmt.Errorf("failed to get chain selector for source chain id %d: %w", sourceChainId, err)
	}
	return &DestCCIPModule{
		Common:                  cmn,
		SourceChainId:           sourceChainId,
		SourceChainSelector:     sourceChainSelector,
		SourceNetworkName:       sourceChain,
		NextSeqNumToCommit:      atomic.NewUint64(1),
		ReportBlessedWatcher:    &sync.Map{},
		ReportBlessedBySeqNum:   &sync.Map{},
		ExecStateChangedWatcher: &sync.Map{},
		ReportAcceptedWatcher:   &sync.Map{},
	}, nil
}

type CCIPRequest struct {
	ReqNo                   int64
	txHash                  string
	txConfirmationTimestamp time.Time
	RequestStat             *testreporters.RequestStat
}

func CCIPRequestFromTxHash(txHash common.Hash, chainClient blockchain.EVMClient) (CCIPRequest, *types.Receipt, error) {
	rcpt, err := chainClient.GetTxReceipt(txHash)
	if err != nil {
		return CCIPRequest{}, nil, err
	}

	hdr, err := chainClient.HeaderByNumber(context.Background(), rcpt.BlockNumber)
	if err != nil {
		return CCIPRequest{}, nil, err
	}
	txConfirmationTimestamp := hdr.Timestamp

	return CCIPRequest{
		txHash:                  txHash.Hex(),
		txConfirmationTimestamp: txConfirmationTimestamp,
	}, rcpt, nil
}

type CCIPLane struct {
	Test                   *testing.T
	Logger                 *zerolog.Logger
	SourceNetworkName      string
	DestNetworkName        string
	SourceChain            blockchain.EVMClient
	DestChain              blockchain.EVMClient
	Source                 *SourceCCIPModule
	Dest                   *DestCCIPModule
	NumberOfReq            int
	Reports                *testreporters.CCIPLaneStats
	Balance                *BalanceSheet
	SentReqs               map[common.Hash][]CCIPRequest
	TotalFee               *big.Int // total fee for all the requests. Used for balance validation.
	ValidationTimeout      time.Duration
	Context                context.Context
	SrcNetworkLaneCfg      *laneconfig.LaneConfig
	DstNetworkLaneCfg      *laneconfig.LaneConfig
	PriceReportingDisabled bool
}

func (lane *CCIPLane) TokenPricesConfig() (string, error) {
	d := &DynamicPriceGetterConfig{
		AggregatorPrices: make(map[common.Address]AggregatorPriceConfig),
		StaticPrices:     make(map[common.Address]StaticPriceConfig),
	}
	// for each token if there is a price aggregator, add it to the aggregator prices
	// else add it to the static prices
	for _, token := range lane.Dest.Common.BridgeTokens {
		err := d.AddPriceConfig(token.Address(), lane.Dest.Common.PriceAggregators, LinkToUSD, lane.DestChain.GetChainID().Uint64())
		if err != nil {
			return "", fmt.Errorf("error in adding PriceConfig for source bridge token %s: %w", token.Address(), err)
		}
	}
	err := d.AddPriceConfig(lane.Dest.Common.FeeToken.Address(), lane.Dest.Common.PriceAggregators, LinkToUSD, lane.DestChain.GetChainID().Uint64())
	if err != nil {
		return "", fmt.Errorf("error adding PriceConfig for dest Fee token %s: %w", lane.Dest.Common.FeeToken.Address(), err)
	}
	err = d.AddPriceConfig(lane.Dest.Common.WrappedNative.Hex(), lane.Dest.Common.PriceAggregators, WrappedNativeToUSD, lane.DestChain.GetChainID().Uint64())
	if err != nil {
		return "", fmt.Errorf("error in adding PriceConfig for dest WrappedNative token %s: %w", lane.Dest.Common.WrappedNative.Hex(), err)
	}
	err = d.AddPriceConfig(lane.Source.Common.WrappedNative.Hex(), lane.Source.Common.PriceAggregators, WrappedNativeToUSD, lane.SourceChain.GetChainID().Uint64())
	if err != nil {
		return "", fmt.Errorf("error in adding PriceConfig for source WrappedNative token %s: %w", lane.Source.Common.WrappedNative.Hex(), err)
	}
	return d.String()
}

func (lane *CCIPLane) SetRemoteChainsOnPool() error {
	if lane.Source.Common.ExistingDeployment {
		return nil
	}
	if len(lane.Source.Common.BridgeTokenPools) != len(lane.Dest.Common.BridgeTokenPools) {
		return fmt.Errorf("source (%d) and dest (%d) bridge token pools length should be same",
			len(lane.Source.Common.BridgeTokenPools), len(lane.Dest.Common.BridgeTokenPools),
		)
	}
	for i, srcPool := range lane.Source.Common.BridgeTokenPools {
		sourceToken := lane.Source.Common.BridgeTokens[i]
		destToken := lane.Dest.Common.BridgeTokens[i]
		dstPool := lane.Dest.Common.BridgeTokenPools[i]

		err := srcPool.SetRemoteChainOnPool(lane.Source.DestChainSelector, dstPool.EthAddress, destToken.ContractAddress)
		if err != nil {
			return err
		}
		err = dstPool.SetRemoteChainOnPool(lane.Dest.SourceChainSelector, srcPool.EthAddress, sourceToken.ContractAddress)
		if err != nil {
			return err
		}
	}
	return nil
}

// OptimizeStorage sets nil to various elements of CCIPLane which are only used
// during lane set up and not used for rest of the test duration
// this is called mainly by load test to keep the memory usage minimum for high number of lanes
func (lane *CCIPLane) OptimizeStorage() {
	lane.Source.Common.FreeUpUnusedSpace()
	lane.Dest.Common.FreeUpUnusedSpace()
	lane.DstNetworkLaneCfg = nil
	lane.SrcNetworkLaneCfg = nil
	// close all header subscriptions for dest chains
	queuedEvents := lane.Dest.Common.ChainClient.GetHeaderSubscriptions()
	for subName := range queuedEvents {
		lane.Dest.Common.ChainClient.DeleteHeaderEventSubscription(subName)
	}
	// close all header subscriptions for source chains except for finalized header
	queuedEvents = lane.Source.Common.ChainClient.GetHeaderSubscriptions()
	for subName := range queuedEvents {
		if subName == blockchain.FinalizedHeaderKey {
			continue
		}
		lane.Source.Common.ChainClient.DeleteHeaderEventSubscription(subName)
	}
}

func (lane *CCIPLane) UpdateLaneConfig() {
	lane.Source.Common.WriteLaneConfig(lane.SrcNetworkLaneCfg)
	lane.SrcNetworkLaneCfg.SrcContractsMu.Lock()
	lane.SrcNetworkLaneCfg.SrcContracts[lane.Source.DestNetworkName] = laneconfig.SourceContracts{
		OnRamp:     lane.Source.OnRamp.Address(),
		DeployedAt: lane.Source.SrcStartBlock,
	}
	lane.SrcNetworkLaneCfg.SrcContractsMu.Unlock()
	lane.Dest.Common.WriteLaneConfig(lane.DstNetworkLaneCfg)
	lane.DstNetworkLaneCfg.DestContractsMu.Lock()
	lane.DstNetworkLaneCfg.DestContracts[lane.Dest.SourceNetworkName] = laneconfig.DestContracts{
		OffRamp:      lane.Dest.OffRamp.Address(),
		CommitStore:  lane.Dest.CommitStore.Address(),
		ReceiverDapp: lane.Dest.ReceiverDapp.Address(),
	}
	lane.DstNetworkLaneCfg.DestContractsMu.Unlock()
}

func (lane *CCIPLane) RecordStateBeforeTransfer() {
	// collect the balance assert.ment to verify balances after transfer
	bal, err := testhelpers.GetBalances(lane.Test, lane.Source.CollectBalanceRequirements())
	require.NoError(lane.Test, err, "fetching source balance")
	lane.Balance.RecordBalance(bal)

	bal, err = testhelpers.GetBalances(lane.Test, lane.Dest.CollectBalanceRequirements())
	require.NoError(lane.Test, err, "fetching dest balance")
	lane.Balance.RecordBalance(bal)

	// save the current block numbers to use in various filter log requests
	lane.TotalFee = big.NewInt(0)
	lane.NumberOfReq = 0
	lane.SentReqs = make(map[common.Hash][]CCIPRequest)
}

func (lane *CCIPLane) AddToSentReqs(txHash common.Hash, reqStats []*testreporters.RequestStat) (*types.Receipt, error) {
	request, rcpt, err := CCIPRequestFromTxHash(txHash, lane.Source.Common.ChainClient)
	if err != nil {
		for _, stat := range reqStats {
			stat.UpdateState(lane.Logger, 0, testreporters.TX, 0, testreporters.Failure, nil)
		}
		return rcpt, fmt.Errorf("could not get request from tx hash %s: %w", txHash.Hex(), err)
	}
	var allRequests []CCIPRequest
	for _, stat := range reqStats {
		allRequests = append(allRequests, CCIPRequest{
			ReqNo:                   stat.ReqNo,
			txHash:                  rcpt.TxHash.Hex(),
			txConfirmationTimestamp: request.txConfirmationTimestamp,
			RequestStat:             stat,
		})
		lane.NumberOfReq++
	}
	lane.SentReqs[rcpt.TxHash] = allRequests
	return rcpt, nil
}

// Multicall sends multiple ccip-send requests in a single transaction
// It will create one transaction for all the requests and will wait for the confirmation
func (lane *CCIPLane) Multicall(noOfRequests int, multiSendAddr common.Address) error {
	var ccipMultipleMsg []contracts.CCIPMsgData
	feeToken := common.HexToAddress(lane.Source.Common.FeeToken.Address())
	genericMsg, err := lane.Source.CCIPMsg(
		lane.Dest.ReceiverDapp.EthAddress,
		lane.Source.Common.AllowOutOfOrder,
		big.NewInt(DefaultDestinationGasLimit),
	)
	if err != nil {
		return fmt.Errorf("failed to form the ccip message: %w", err)
	}
	destChainSelector, err := chainselectors.SelectorFromChainId(lane.Source.DestinationChainId)
	if err != nil {
		return fmt.Errorf("failed getting the chain selector: %w", err)
	}
	var reqStats []*testreporters.RequestStat
	var txstats []*testreporters.TransactionStats
	for i := 1; i <= noOfRequests; i++ {
		// form the message for transfer
		msg := genericMsg
		msg.Data = []byte(fmt.Sprintf("msg %d", i))
		sendData := contracts.CCIPMsgData{
			Msg:           msg,
			RouterAddr:    lane.Source.Common.Router.EthAddress,
			ChainSelector: destChainSelector,
		}

		fee, err := lane.Source.Common.Router.GetFee(destChainSelector, msg)
		if err != nil {
			reason, _ := blockchain.RPCErrorFromError(err)
			if reason != "" {
				return fmt.Errorf("failed getting the fee: %s", reason)
			}
			return fmt.Errorf("failed getting the fee: %w", err)
		}
		log.Info().Str("fee", fee.String()).Msg("calculated fee")
		sendData.Fee = fee
		lane.TotalFee = new(big.Int).Add(lane.TotalFee, fee)
		ccipMultipleMsg = append(ccipMultipleMsg, sendData)
		// if token transfer is required, transfer the token amount to multisend
		for j, amount := range lane.Source.TransferAmount {
			// if length of sourceCCIP.TransferAmount is more than available bridge token use first bridge token
			token := lane.Source.Common.BridgeTokens[0]
			if j < len(lane.Source.Common.BridgeTokens) {
				token = lane.Source.Common.BridgeTokens[j]
			}
			err = token.Transfer(lane.SourceChain.GetDefaultWallet(), multiSendAddr.Hex(), amount)
			if err != nil {
				return err
			}
		}
		stat := testreporters.NewCCIPRequestStats(int64(lane.NumberOfReq+i), lane.SourceNetworkName, lane.DestNetworkName)
		txstats = append(txstats, &testreporters.TransactionStats{
			Fee:                fee.String(),
			NoOfTokensSent:     len(msg.TokenAmounts),
			MessageBytesLength: int64(len(msg.Data)),
		})
		reqStats = append(reqStats, stat)
	}
	isNative := true
	// transfer the fee amount to multisend
	if feeToken != (common.Address{}) {
		isNative = false
		err = lane.Source.Common.FeeToken.Transfer(multiSendAddr.Hex(), lane.TotalFee)
		if err != nil {
			return err
		}
	}

	tx, err := contracts.MultiCallCCIP(lane.Source.Common.ChainClient, multiSendAddr.Hex(), ccipMultipleMsg, isNative)
	if err != nil {
		// update the stats as failure for all the requests in the multicall tx
		for _, stat := range reqStats {
			stat.UpdateState(lane.Logger, 0, testreporters.TX, 0, testreporters.Failure, nil)
		}
		return fmt.Errorf("failed to send the multicall: %w", err)
	}
	rcpt, err := lane.AddToSentReqs(tx.Hash(), reqStats)
	if err != nil {
		return err
	}
	var gasUsed uint64
	if rcpt != nil {
		gasUsed = rcpt.GasUsed
	}
	// update the stats for all the requests in the multicall tx
	for i, stat := range reqStats {
		txstats[i].GasUsed = gasUsed
		txstats[i].TxHash = tx.Hash().Hex()
		stat.UpdateState(lane.Logger, 0, testreporters.TX, 0, testreporters.Success, txstats[i])
	}
	return nil
}

// SendRequests sends individual ccip-send requests in different transactions
// It will create noOfRequests transactions
func (lane *CCIPLane) SendRequests(noOfRequests int, gasLimit *big.Int) error {
	for i := 1; i <= noOfRequests; i++ {
		stat := testreporters.NewCCIPRequestStats(int64(lane.NumberOfReq+i), lane.SourceNetworkName, lane.DestNetworkName)
		txHash, txConfirmationDur, fee, err := lane.Source.SendRequest(lane.Dest.ReceiverDapp.EthAddress, gasLimit)
		if err != nil {
			stat.UpdateState(lane.Logger, 0, testreporters.TX, txConfirmationDur, testreporters.Failure, nil)
			return fmt.Errorf("could not send request: %w", err)
		}
		err = lane.Source.Common.ChainClient.WaitForEvents()
		if err != nil {
			stat.UpdateState(lane.Logger, 0, testreporters.TX, txConfirmationDur, testreporters.Failure, nil)
			return fmt.Errorf("could not send request: %w", err)
		}

		noOfTokens := 0
		for _, tokenAmount := range lane.Source.TransferAmount { // Only count tokens that are actually sent
			if tokenAmount != nil && tokenAmount.Cmp(big.NewInt(0)) > 0 {
				noOfTokens++
			}
		}
		_, err = lane.AddToSentReqs(txHash, []*testreporters.RequestStat{stat})
		if err != nil {
			return err
		}
		stat.UpdateState(lane.Logger, 0, testreporters.TX, txConfirmationDur, testreporters.Success, nil)
		lane.TotalFee = bigmath.Add(lane.TotalFee, fee)
	}

	return nil
}

// manualExecutionOpts modify how ExecuteManually behaves
type manualExecutionOpts struct {
	timeout time.Duration
}

// ManualExecutionOption is a function that modifies ExecuteManually behavior
type ManualExecutionOption func(*manualExecutionOpts)

// WithConfirmationTimeout sets a custom timeout for waiting for the confirmation of the manual execution
func WithConfirmationTimeout(timeout time.Duration) ManualExecutionOption {
	return func(opts *manualExecutionOpts) {
		opts.timeout = timeout
	}
}

// ExecuteManually attempts to execute pending CCIP transactions manually.
// This is necessary in situations where Smart Execution window for that message is over and Offchain plugin
// will not attempt to execute the message.In such situation any further message from same sender will not be executed until
// the blocking message is executed by the OffRamp.
// More info: https://docs.chain.link/ccip/concepts/manual-execution#manual-execution
func (lane *CCIPLane) ExecuteManually(options ...ManualExecutionOption) error {
	var opts manualExecutionOpts
	for _, opt := range options {
		if opt != nil {
			opt(&opts)
		}
	}
	if opts.timeout == 0 {
		opts.timeout = lane.ValidationTimeout
	}

	onRampABI, err := abi.JSON(strings.NewReader(evm_2_evm_onramp.EVM2EVMOnRampABI))
	if err != nil {
		return err
	}
	sendReqTopic := onRampABI.Events["CCIPSendRequested"].ID
	for txHash, req := range lane.SentReqs {
		for _, ccipReq := range req {
			lane.Logger.Info().Str("ccip-send", txHash.Hex()).Msg("Executing request manually")
			seqNum := ccipReq.RequestStat.SeqNum
			sendReqReceipt, err := lane.Source.Common.ChainClient.GetTxReceipt(txHash)
			if err != nil {
				return err
			}
			if sendReqReceipt == nil {
				return fmt.Errorf("could not find the receipt for tx %s", txHash.Hex())
			}
			commitStat, ok := ccipReq.RequestStat.StatusByPhase[testreporters.Commit]
			if !ok {
				return fmt.Errorf("could not find the commit phase in the request stats, reqNo %d", ccipReq.RequestStat.ReqNo)
			}
			commitTx := commitStat.SendTransactionStats.TxHash
			commitReceipt, err := lane.DestChain.GetTxReceipt(common.HexToHash(commitTx))
			if err != nil {
				return err
			}
			var logIndex uint
			// find the send request log index sendReqReceipt
			for _, sendReqLog := range sendReqReceipt.Logs {
				if sendReqLog.Topics[0] == sendReqTopic {
					logSeqNum, err := lane.Source.OnRamp.Instance.ParseCCIPSendRequested(*sendReqLog)
					if err != nil {
						return err
					}
					if logSeqNum == seqNum {
						logIndex = sendReqLog.Index
					}
				}
			}
			destChainSelector, err := chainselectors.SelectorFromChainId(lane.DestChain.GetChainID().Uint64())
			if err != nil {
				return err
			}
			sourceChainSelector, err := chainselectors.SelectorFromChainId(lane.SourceChain.GetChainID().Uint64())
			if err != nil {
				return err
			}
			// Calling `TransactionOpts` will automatically increase the nonce, so if this fails, any other destination transactions will time out
			destUser, err := lane.DestChain.TransactionOpts(lane.DestChain.GetDefaultWallet())
			if err != nil {
				return err
			}
			args := testhelpers.ManualExecArgs{
				SourceChainID:    sourceChainSelector,
				DestChainID:      destChainSelector,
				DestUser:         destUser,
				SourceChain:      lane.SourceChain.Backend(),
				DestChain:        lane.DestChain.Backend(),
				SourceStartBlock: sendReqReceipt.BlockNumber,
				DestStartBlock:   commitReceipt.BlockNumber.Uint64(),
				SendReqTxHash:    txHash.Hex(),
				CommitStore:      lane.Dest.CommitStore.Address(),
				OnRamp:           lane.Source.OnRamp.Address(),
				OffRamp:          lane.Dest.OffRamp.Address(),
				SendReqLogIndex:  logIndex,
				GasLimit:         big.NewInt(DefaultDestinationGasLimit),
			}
			timeNow := time.Now().UTC()
			tx, err := args.ExecuteManually()
			if err != nil {
				return fmt.Errorf("could not execute manually: %w seqNum %d", err, seqNum)
			}

			ctx, cancel := context.WithTimeout(context.Background(), opts.timeout)
			rec, err := bind.WaitMined(ctx, lane.DestChain.DeployBackend(), tx)
			if err != nil {
				cancel()
				return fmt.Errorf("could not get receipt: %w seqNum %d", err, seqNum)
			}
			cancel()
			if rec.Status != 1 {
				return fmt.Errorf(
					"manual execution failed for seqNum %d with receipt status %d, use the revert-reason script on this transaction hash '%s' and this sender address '%s'",
					seqNum, rec.Status, tx.Hash().Hex(), destUser.From.Hex(),
				)
			}
			lane.Logger.Info().Uint64("seqNum", seqNum).Msg("Manual Execution completed")
			_, err = lane.Dest.AssertEventExecutionStateChanged(lane.Logger, seqNum, opts.timeout,
				timeNow, ccipReq.RequestStat, testhelpers.ExecutionStateSuccess,
			)
			if err != nil {
				return fmt.Errorf("could not validate ExecutionStateChanged event: %w", err)
			}
		}
	}
	return nil
}

// validationOptions are used in the ValidateRequests function to specify which phase is expected to fail and how
type validationOptions struct {
	expectAnyPhaseToFail bool
	phaseExpectedToFail  testreporters.Phase // the phase expected to fail
	expectedErrorMessage string              // if provided, we're looking for a specific error message
	timeout              time.Duration       // timeout for the validation
}

// ValidationOptionFunc is a function that can be passed to ValidateRequests to specify which phase is expected to fail
type ValidationOptionFunc func(opts *validationOptions)

// PhaseSpecificValidationOptionFunc can specify how exactly you want a phase to fail
type PhaseSpecificValidationOptionFunc func(*validationOptions)

// WithErrorMessage specifies the expected error message for the phase that is expected to fail.
func WithErrorMessage(expectedErrorMessage string) PhaseSpecificValidationOptionFunc {
	return func(opts *validationOptions) {
		opts.expectedErrorMessage = expectedErrorMessage
	}
}

// WithTimeout specifies a custom timeout for validating that the phase failed.
func WithTimeout(timeout time.Duration) PhaseSpecificValidationOptionFunc {
	return func(opts *validationOptions) {
		opts.timeout = timeout
	}
}

// ExpectPhaseToFail specifies that a specific phase is expected to fail.
// You can optionally provide an expected error message, if you don't have one in mind, just pass an empty string.
// shouldExist is used to specify whether the phase should exist or not, which is only applicable to the `ExecStateChanged` phase.
// If you expect the `ExecStateChanged` events to be there, but in a "failed" state, set this to true.
// It will otherwise be ignored.
func ExpectPhaseToFail(phase testreporters.Phase, phaseSpecificOptions ...PhaseSpecificValidationOptionFunc) ValidationOptionFunc {
	return func(opts *validationOptions) {
		opts.phaseExpectedToFail = phase
		for _, f := range phaseSpecificOptions {
			if f != nil {
				f(opts)
			}
		}
	}
}

// ExpectAnyPhaseToFail expects any phase in CCIP transaction to fail.
func ExpectAnyPhaseToFail(phaseSpecificOptions ...PhaseSpecificValidationOptionFunc) ValidationOptionFunc {
	return func(opts *validationOptions) {
		opts.expectAnyPhaseToFail = true
		for _, f := range phaseSpecificOptions {
			if f != nil {
				f(opts)
			}
		}
	}
}

// ValidateRequests validates all sent request events.
// If you expect a specific phase to fail, you can pass a validationOptionFunc to specify exactly which one.
// If not, just pass in nil.
func (lane *CCIPLane) ValidateRequests(validationOptionFuncs ...ValidationOptionFunc) {
	var opts validationOptions
	for _, f := range validationOptionFuncs {
		if f != nil {
			f(&opts)
		}
	}
	for txHash, ccipReqs := range lane.SentReqs {
		require.Greater(lane.Test, len(ccipReqs), 0, "no ccip requests found for tx hash")
		require.NoError(lane.Test, lane.ValidateRequestByTxHash(txHash, opts), "validating request events by tx hash")
	}
	if len(validationOptionFuncs) > 0 {
		return
	}
	// Asserting balances reliably work only for simulated private chains. The testnet contract balances might get updated by other transactions
	// verify the fee amount is deducted from sender, added to receiver token balances and
	if len(lane.Source.TransferAmount) > 0 && len(lane.Source.Common.BridgeTokens) > 0 {
		lane.Source.UpdateBalance(int64(lane.NumberOfReq), lane.TotalFee, lane.Balance)
		lane.Dest.UpdateBalance(lane.Source.TransferAmount, int64(lane.NumberOfReq), lane.Balance)
	}
}

// ValidateRequestByTxHash validates the request events by tx hash.
// If a phaseExpectedToFail is provided, it will return no error if that phase fails, but will error if it succeeds.
func (lane *CCIPLane) ValidateRequestByTxHash(txHash common.Hash, opts validationOptions) error {
	var (
		reqStats       []*testreporters.RequestStat
		timeout        = lane.ValidationTimeout
		ccipRequests   = lane.SentReqs[txHash]
		txConfirmation = ccipRequests[0].txConfirmationTimestamp
	)
	require.Greater(lane.Test, len(ccipRequests), 0, "no ccip requests found for tx hash")

	defer func() {
		for _, req := range ccipRequests {
			lane.Reports.UpdatePhaseStatsForReq(req.RequestStat)
		}
	}()
	for _, req := range ccipRequests {
		reqStats = append(reqStats, req.RequestStat)
	}

	if (opts.phaseExpectedToFail == testreporters.CCIPSendRe || opts.expectAnyPhaseToFail) && opts.timeout != 0 {
		timeout = opts.timeout
	}
	msgLogs, ccipSendReqGenAt, err := lane.Source.AssertEventCCIPSendRequested(
		lane.Logger, txHash.Hex(), timeout, txConfirmation, reqStats,
	)
	if shouldReturn, phaseErr := isPhaseValid(lane.Logger, testreporters.CCIPSendRe, opts, err); shouldReturn {
		return phaseErr
	}

	sourceLogFinalizedAt, _, err := lane.Source.AssertSendRequestedLogFinalized(lane.Logger, txHash, msgLogs, ccipSendReqGenAt, reqStats)
	if shouldReturn, phaseErr := isPhaseValid(lane.Logger, testreporters.SourceLogFinalized, opts, err); shouldReturn {
		return phaseErr
	}
	for _, msgLog := range msgLogs {
		seqNumber := msgLog.SequenceNumber
		lane.Logger = ptr.Ptr(lane.Logger.With().Str("msgId", fmt.Sprintf("0x%x", msgLog.MessageId[:])).Logger())
		var reqStat *testreporters.RequestStat
		for _, stat := range reqStats {
			if stat.SeqNum == seqNumber {
				reqStat = stat
				break
			}
		}
		if reqStat == nil {
			return fmt.Errorf("could not find request stat for seq number %d", seqNumber)
		}

		if opts.expectAnyPhaseToFail && opts.timeout != 0 {
			timeout = opts.timeout
		}

		if opts.phaseExpectedToFail == testreporters.Commit && opts.timeout != 0 {
			timeout = opts.timeout
		}
		err = lane.Dest.AssertSeqNumberExecuted(lane.Logger, seqNumber, timeout, sourceLogFinalizedAt, reqStat)
		if shouldReturn, phaseErr := isPhaseValid(lane.Logger, testreporters.Commit, opts, err); shouldReturn {
			return phaseErr
		}

		// Verify whether commitStore has accepted the report
		commitReport, reportAcceptedAt, err := lane.Dest.AssertEventReportAccepted(
			lane.Logger, seqNumber, timeout, sourceLogFinalizedAt, reqStat,
		)
		if shouldReturn, phaseErr := isPhaseValid(lane.Logger, testreporters.Commit, opts, err); shouldReturn {
			return phaseErr
		}

		if opts.phaseExpectedToFail == testreporters.ReportBlessed && opts.timeout != 0 {
			timeout = opts.timeout
		}
		reportBlessedAt, err := lane.Dest.AssertReportBlessed(lane.Logger, seqNumber, timeout, *commitReport, reportAcceptedAt, reqStat)
		if shouldReturn, phaseErr := isPhaseValid(lane.Logger, testreporters.ReportBlessed, opts, err); shouldReturn {
			return phaseErr
		}

		if opts.phaseExpectedToFail == testreporters.ExecStateChanged && opts.timeout != 0 {
			timeout = opts.timeout
		}
		// Verify whether the execution state is changed and the transfer is successful
		_, err = lane.Dest.AssertEventExecutionStateChanged(
			lane.Logger, seqNumber,
			timeout,
			reportBlessedAt,
			reqStat,
			testhelpers.ExecutionStateSuccess,
		)
		if shouldReturn, phaseErr := isPhaseValid(lane.Logger, testreporters.ExecStateChanged, opts, err); shouldReturn {
			return phaseErr
		}
	}
	if opts.expectAnyPhaseToFail {
		return fmt.Errorf("expected at least any one phase to fail but no phase got failed")
	}

	return nil
}

// isPhaseValid checks if the phase is in a valid state or not given expectations.
// If `shouldComplete` is true, it means that the phase validation is meant to end and we should return from the calling function.
func isPhaseValid(
	logger *zerolog.Logger,
	currentPhase testreporters.Phase,
	opts validationOptions,
	err error,
) (shouldComplete bool, validationError error) {
	if opts.expectAnyPhaseToFail && err != nil {
		logmsg := logger.Info().Str("Failed with Error", err.Error()).Str("Phase", string(currentPhase))
		logmsg.Msg("Phase failed, as expected")
		return true, nil
	}
	// If no phase is expected to fail or the current phase is not the one expected to fail, we just return what we were given
	if opts.phaseExpectedToFail == "" || currentPhase != opts.phaseExpectedToFail {
		return err != nil, err
	}
	if err == nil && currentPhase == opts.phaseExpectedToFail {
		return true, fmt.Errorf("expected phase '%s' to fail, but it passed", opts.phaseExpectedToFail)
	}
	logmsg := logger.Info().Str("Failed with Error", err.Error()).Str("Phase", string(currentPhase))

	if opts.expectedErrorMessage != "" {
		if !strings.Contains(err.Error(), opts.expectedErrorMessage) {
			return true, fmt.Errorf("expected phase '%s' to fail with error message '%s' but got error '%s'", currentPhase, opts.expectedErrorMessage, err.Error())
		}
		logmsg.Str("Expected Error Message", opts.expectedErrorMessage)
	}
	logmsg.Msg("Phase failed, as expected")
	return true, nil
}

// DisableAllRateLimiting disables all rate limiting for the lane, including ARL and token pool rate limits
func (lane *CCIPLane) DisableAllRateLimiting() error {
	src := lane.Source
	dest := lane.Dest

	// Tell OnRamp to not include any tokens in ARL
	err := src.SetAllTokenTransferFeeConfigs(false)
	if err != nil {
		return fmt.Errorf("error disabling token transfer fee config for OnRamp: %w", err)
	}
	err = dest.RemoveAllRateLimitTokens(context.Background())
	if err != nil {
		return fmt.Errorf("error removing rate limited tokens for OffRamp: %w", err)
	}
	// Disable ARL for OnRamp and OffRamp
	err = src.OnRamp.SetRateLimit(evm_2_evm_onramp.RateLimiterConfig{
		IsEnabled: false,
		Capacity:  big.NewInt(0),
		Rate:      big.NewInt(0),
	})
	if err != nil {
		return fmt.Errorf("error disabling rate limit for source onramp: %w", err)
	}
	err = dest.OffRamp.SetRateLimit(contracts.RateLimiterConfig{
		IsEnabled: false,
		Capacity:  big.NewInt(0),
		Rate:      big.NewInt(0),
	})
	if err != nil {
		return fmt.Errorf("error disabling rate limit for destination offramp: %w", err)
	}
	// Disable individual token pool rate limits
	for i, tokenPool := range src.Common.BridgeTokenPools {
		err = tokenPool.SetRemoteChainRateLimits(src.DestChainSelector, token_pool.RateLimiterConfig{
			IsEnabled: false,
			Capacity:  big.NewInt(0),
			Rate:      big.NewInt(0),
		})
		if err != nil {
			return fmt.Errorf("error disabling rate limit for token pool %d: %w", i, err)
		}
	}
	for i, tokenPool := range dest.Common.BridgeTokenPools {
		err = tokenPool.SetRemoteChainRateLimits(dest.SourceChainSelector, token_pool.RateLimiterConfig{
			IsEnabled: false,
			Capacity:  big.NewInt(0),
			Rate:      big.NewInt(0),
		})
		if err != nil {
			return fmt.Errorf("error disabling rate limit for token pool %d: %w", i, err)
		}
	}
	err = src.Common.ChainClient.WaitForEvents()
	if err != nil {
		return fmt.Errorf("error waiting for source chain events: %w", err)
	}
	err = dest.Common.ChainClient.WaitForEvents()
	if err != nil {
		return fmt.Errorf("error waiting for destination chain events: %w", err)
	}
	lane.Logger.Info().Msg("Disabled all rate limiting")
	return nil
}

func (lane *CCIPLane) StartEventWatchers() error {
	lane.Logger.Info().Msg("Starting event watchers")
	if lane.Source.Common.ChainClient.GetNetworkConfig().FinalityDepth == 0 {
		err := lane.Source.Common.ChainClient.PollFinality()
		if err != nil {
			return err
		}
	}

	go lane.Source.Common.PollRPCConnection(lane.Context, lane.Logger)
	go lane.Dest.Common.PollRPCConnection(lane.Context, lane.Logger)

	sendReqEventLatest := make(chan *evm_2_evm_onramp.EVM2EVMOnRampCCIPSendRequested)
	senReqSub := event.Resubscribe(DefaultResubscriptionTimeout, func(_ context.Context) (event.Subscription, error) {
		sub, err := lane.Source.OnRamp.WatchCCIPSendRequested(nil, sendReqEventLatest)
		if err != nil {
			log.Error().Err(err).Msg("error in subscribing to CCIPSendRequested event")
		}
		return sub, err
	})
	if senReqSub == nil {
		return fmt.Errorf("failed to subscribe to CCIPSendRequested event")
	}
	go func(sub event.Subscription) {
		defer sub.Unsubscribe()
		for {
			select {
			case e := <-sendReqEventLatest:
				lane.Logger.Info().Msgf("CCIPSendRequested event received for seq number %d", e.Message.SequenceNumber)
				eventsForTx, ok := lane.Source.CCIPSendRequestedWatcher.Load(e.Raw.TxHash.Hex())
				if ok {
					lane.Source.CCIPSendRequestedWatcher.Store(e.Raw.TxHash.Hex(), append(eventsForTx.([]*contracts.SendReqEventData),
						&contracts.SendReqEventData{
							MessageId:      e.Message.MessageId,
							SequenceNumber: e.Message.SequenceNumber,
							DataLength:     len(e.Message.Data),
							NoOfTokens:     len(e.Message.TokenAmounts),
							LogInfo: contracts.LogInfo{
								BlockNumber: e.Raw.BlockNumber,
								TxHash:      e.Raw.TxHash,
							},
							Fee: e.Message.FeeTokenAmount,
						}))
				} else {
					lane.Source.CCIPSendRequestedWatcher.Store(e.Raw.TxHash.Hex(), []*contracts.SendReqEventData{
						{
							MessageId:      e.Message.MessageId,
							SequenceNumber: e.Message.SequenceNumber,
							DataLength:     len(e.Message.Data),
							NoOfTokens:     len(e.Message.TokenAmounts),
							LogInfo: contracts.LogInfo{
								BlockNumber: e.Raw.BlockNumber,
								TxHash:      e.Raw.TxHash,
							},
							Fee: e.Message.FeeTokenAmount,
						},
					})
				}

				lane.Source.CCIPSendRequestedWatcher = testutils.DeleteNilEntriesFromMap(lane.Source.CCIPSendRequestedWatcher)
			case <-lane.Context.Done():
				return
			}
		}
	}(senReqSub)

	reportAcceptedEvent := make(chan *commit_store.CommitStoreReportAccepted)
	reportAccSub := event.Resubscribe(DefaultResubscriptionTimeout, func(_ context.Context) (event.Subscription, error) {
		sub, err := lane.Dest.CommitStore.WatchReportAccepted(nil, reportAcceptedEvent)
		if err != nil {
			log.Error().Err(err).Msg("error in subscribing to ReportAccepted event")
		}
		return sub, err
	})
	if reportAccSub == nil {
		return fmt.Errorf("failed to subscribe to ReportAccepted event")
	}
	go func(sub event.Subscription) {
		defer sub.Unsubscribe()
		for {
			select {
			case e := <-reportAcceptedEvent:
				lane.Logger.Info().Interface("Interval", e.Report.Interval).Msgf("ReportAccepted event received")
				for i := e.Report.Interval.Min; i <= e.Report.Interval.Max; i++ {
					lane.Dest.ReportAcceptedWatcher.Store(i, &contracts.CommitStoreReportAccepted{
						Min:        e.Report.Interval.Min,
						Max:        e.Report.Interval.Max,
						MerkleRoot: e.Report.MerkleRoot,
						LogInfo: contracts.LogInfo{
							BlockNumber: e.Raw.BlockNumber,
							TxHash:      e.Raw.TxHash,
						},
					})
				}
				lane.Dest.ReportAcceptedWatcher = testutils.DeleteNilEntriesFromMap(lane.Dest.ReportAcceptedWatcher)
			case <-lane.Context.Done():
				return
			}
		}
	}(reportAccSub)

	if lane.Dest.Common.ARM != nil {
		reportBlessedEvent := make(chan *rmn_contract.RMNContractTaggedRootBlessed)
		blessedSub := event.Resubscribe(DefaultResubscriptionTimeout, func(_ context.Context) (event.Subscription, error) {
			sub, err := lane.Dest.Common.ARM.Instance.WatchTaggedRootBlessed(nil, reportBlessedEvent, nil)
			if err != nil {
				log.Error().Err(err).Msg("error in subscribing to TaggedRootBlessed event")
			}
			return sub, err
		})
		if blessedSub == nil {
			return fmt.Errorf("failed to subscribe to TaggedRootBlessed event")
		}
		go func(sub event.Subscription) {
			defer sub.Unsubscribe()
			for {
				select {
				case e := <-reportBlessedEvent:
					lane.Logger.Info().Msgf("TaggedRootBlessed event received for root %x", e.TaggedRoot.Root)
					if e.TaggedRoot.CommitStore == lane.Dest.CommitStore.EthAddress {
						lane.Dest.ReportBlessedWatcher.Store(e.TaggedRoot.Root, &contracts.LogInfo{
							BlockNumber: e.Raw.BlockNumber,
							TxHash:      e.Raw.TxHash,
						})
					}
					lane.Dest.ReportBlessedWatcher = testutils.DeleteNilEntriesFromMap(lane.Dest.ReportBlessedWatcher)
				case <-lane.Context.Done():
					return
				}
			}
		}(blessedSub)
	}

	execStateChangedEventLatest := make(chan *evm_2_evm_offramp.EVM2EVMOffRampExecutionStateChanged)
	execSub := event.Resubscribe(DefaultResubscriptionTimeout, func(_ context.Context) (event.Subscription, error) {
		sub, err := lane.Dest.OffRamp.WatchExecutionStateChanged(nil, execStateChangedEventLatest, nil, nil)
		if err != nil {
			log.Error().Err(err).Msg("error in subscribing to ExecutionStateChanged event")
		}
		return sub, err
	})
	if execSub == nil {
		return fmt.Errorf("failed to subscribe to ExecutionStateChanged event")
	}
	go func(sub event.Subscription) {
		defer sub.Unsubscribe()
		for {
			select {
			case e := <-execStateChangedEventLatest:
				lane.Logger.Info().Msgf("Execution state changed event received for seq number %d", e.SequenceNumber)
				lane.Dest.ExecStateChangedWatcher.Store(e.SequenceNumber, &contracts.EVM2EVMOffRampExecutionStateChanged{
					SequenceNumber: e.SequenceNumber,
					MessageId:      e.MessageId,
					State:          e.State,
					ReturnData:     e.ReturnData,
					LogInfo: contracts.LogInfo{
						BlockNumber: e.Raw.BlockNumber,
						TxHash:      e.Raw.TxHash,
					},
				})
				lane.Dest.ExecStateChangedWatcher = testutils.DeleteNilEntriesFromMap(lane.Dest.ExecStateChangedWatcher)
			case <-lane.Context.Done():
				return
			}
		}
	}(execSub)
	return nil
}

func (lane *CCIPLane) CleanUp(clearFees bool) error {
	lane.Logger.Info().Msg("Cleaning up lane")
	if lane.Source.Common.ChainClient.GetNetworkConfig().FinalityDepth == 0 {
		lane.Source.Common.ChainClient.CancelFinalityPolling()
	}
	// recover fees from onRamp contract
	if clearFees && !lane.Source.Common.ChainClient.NetworkSimulated() {
		err := lane.Source.PayCCIPFeeToOwnerAddress()
		if err != nil {
			return err
		}
	}
	err := lane.Dest.Common.ChainClient.Close()
	if err != nil {
		return err
	}
	return lane.Source.Common.ChainClient.Close()
}

// DeployNewCCIPLane sets up a lane and initiates lane.Source and lane.Destination
// If configureCLNodes is true it sets up jobs and contract config for the lane
func (lane *CCIPLane) DeployNewCCIPLane(
	setUpCtx context.Context,
	env *CCIPTestEnv,
	testConf *testconfig.CCIPTestGroupConfig,
	bootstrapAdded *atomic.Bool,
	jobErrGroup *errgroup.Group,
) error {
	var (
		err                    error
		sourceChainClient      = lane.SourceChain
		destChainClient        = lane.DestChain
		srcConf                = lane.SrcNetworkLaneCfg
		destConf               = lane.DstNetworkLaneCfg
		commitAndExecOnSameDON = pointer.GetBool(testConf.CommitAndExecuteOnSameDON)
		allowOutOfOrder        = pointer.GetBool(testConf.AllowOutOfOrder)
		withPipeline           = pointer.GetBool(testConf.TokenConfig.WithPipeline)
		configureCLNodes       = !pointer.GetBool(testConf.ExistingDeployment)
	)

	lane.Source, err = DefaultSourceCCIPModule(
		lane.Logger,
		testConf,
		sourceChainClient, destChainClient.GetChainID().Uint64(),
		destChainClient.GetNetworkName(),
		srcConf,
	)
	if err != nil {
		return fmt.Errorf("failed to create source module: %w", err)
	}

	// If AllowOutOfOrder is set globally in test config, the assumption is to set it for every lane.
	// However, if this is set as false, and set as true for specific chain in lane_configs, then apply it
	// only for a lane where source network is of that chain.
	if allowOutOfOrder {
		lane.Source.Common.AllowOutOfOrder = true
	}
	lane.Dest, err = DefaultDestinationCCIPModule(
		lane.Logger, testConf,
		destChainClient, sourceChainClient.GetChainID().Uint64(),
		sourceChainClient.GetNetworkName(), destConf,
	)
	if err != nil {
		return fmt.Errorf("failed to create destination module: %w", err)
	}

	// deploy all source contracts
	err = lane.Source.DeployContracts(srcConf)
	if err != nil {
		return fmt.Errorf("failed to deploy source contracts: %w", err)
	}
	// deploy all destination contracts
	err = lane.Dest.DeployContracts(*lane.Source, destConf)
	if err != nil {
		return fmt.Errorf("failed to deploy destination contracts: %w", err)
	}
	// if it's a new USDC deployment, sync the USDC domain
	err = lane.Source.Common.SyncUSDCDomain(lane.Dest.Common.TokenTransmitter, lane.Dest.Common.BridgeTokenPools, lane.Source.DestinationChainId)
	if err != nil {
		return fmt.Errorf("failed to sync USDC domain: %w", err)
	}

	lane.UpdateLaneConfig()

	// if lane is being set up for already configured CL nodes and contracts
	// no further action is necessary
	if !configureCLNodes {
		return nil
	}
	err = lane.Source.Common.WatchForPriceUpdates(setUpCtx, lane.Logger)
	if err != nil {
		return fmt.Errorf("error in starting price update watch %w", err)
	}
	if env == nil {
		return fmt.Errorf("test environment not set")
	}
	// wait for the CL nodes to be ready before moving ahead with job creation
	err = env.CLNodeWithKeyReady.Wait()
	if err != nil {
		return fmt.Errorf("failed to wait for CL nodes to be ready: %w", err)
	}
	clNodesWithKeys := env.CLNodesWithKeys
	// set up ocr2 jobs
	clNodes, exists := clNodesWithKeys[lane.Dest.Common.ChainClient.GetChainID().String()]
	if !exists {
		return fmt.Errorf("could not find CL nodes for %s", lane.Dest.Common.ChainClient.GetChainID().String())
	}
	bootstrapCommit := clNodes[0]
	var bootstrapExec *client.CLNodesWithKeys
	commitNodes := clNodes[env.CommitNodeStartIndex : env.CommitNodeStartIndex+env.NumOfCommitNodes]
	execNodes := clNodes[env.ExecNodeStartIndex : env.ExecNodeStartIndex+env.NumOfExecNodes]
	if !commitAndExecOnSameDON {
		if len(clNodes) < 11 {
			return fmt.Errorf("not enough CL nodes for separate commit and execution nodes")
		}
		bootstrapExec = clNodes[1] // for a set-up of different commit and execution nodes second node is the bootstrapper for execution nodes
	}

	// save the current block numbers. If there is a delay between job start up and ocr config set up, the jobs will
	// replay the log polling from these mentioned block number. The dest block number should ideally be the block number on which
	// contract config is set and the source block number should be the one on which the ccip send request is performed.
	// Here for simplicity we are just taking the current block number just before the job is created.
	currentBlockOnDest, err := destChainClient.LatestBlockNumber(context.Background())
	if err != nil {
		return fmt.Errorf("getting current block should be successful in destination chain %w", err)
	}

	var killgrave *ctftestenv.Killgrave
	if env.LocalCluster != nil {
		killgrave = env.LocalCluster.MockAdapter
	}
	var tokenAddresses []string
	for _, token := range lane.Dest.Common.BridgeTokens {
		tokenAddresses = append(tokenAddresses, token.Address())
	}
	tokenAddresses = append(tokenAddresses, lane.Dest.Common.FeeToken.Address(), lane.Source.Common.WrappedNative.Hex(), lane.Dest.Common.WrappedNative.Hex())

	// Only one off pipeline or price getter to be set.
	tokenPricesUSDPipeline := ""
	tokenPricesConfigJson := ""
	if withPipeline {
		tokensUSDUrl := TokenPricePipelineURLs(tokenAddresses, killgrave, env.MockServer)
		tokenPricesUSDPipeline = TokenFeeForMultipleTokenAddr(tokensUSDUrl)
	} else {
		tokenPricesConfigJson, err = lane.TokenPricesConfig()
		if err != nil {
			return fmt.Errorf("error getting token prices config %w", err)
		}
		lane.Logger.Info().Str("tokenPricesConfigJson", tokenPricesConfigJson).Msg("Price getter config")
	}

	jobParams := integrationtesthelpers.CCIPJobSpecParams{
		OffRamp:                lane.Dest.OffRamp.EthAddress,
		CommitStore:            lane.Dest.CommitStore.EthAddress,
		SourceChainName:        sourceChainClient.GetNetworkName(),
		DestChainName:          destChainClient.GetNetworkName(),
		DestEvmChainId:         destChainClient.GetChainID().Uint64(),
		SourceStartBlock:       lane.Source.SrcStartBlock,
		TokenPricesUSDPipeline: tokenPricesUSDPipeline,
		PriceGetterConfig:      tokenPricesConfigJson,
		DestStartBlock:         currentBlockOnDest,
	}
	if !lane.Source.Common.ExistingDeployment && lane.Source.Common.IsUSDCDeployment() {
		api := ""
		if killgrave != nil {
			api = killgrave.InternalEndpoint
		}
		if env.MockServer != nil {
			api = env.MockServer.Config.ClusterURL
		}
		if lane.Source.Common.TokenTransmitter == nil {
			return fmt.Errorf("token transmitter address not set")
		}
		// Only one USDC allowed per chain
		jobParams.USDCConfig = &config.USDCConfig{
			SourceTokenAddress:              common.HexToAddress(lane.Source.Common.BridgeTokens[0].Address()),
			SourceMessageTransmitterAddress: lane.Source.Common.TokenTransmitter.ContractAddress,
			AttestationAPI:                  api,
			AttestationAPITimeoutSeconds:    5,
		}
	}
	if !lane.Source.Common.ExistingDeployment && lane.Source.Common.IsLBTCDeployment() {
		api := ""
		if killgrave != nil {
			api = killgrave.InternalEndpoint
		}
		if env.MockServer != nil {
			api = env.MockServer.Config.ClusterURL
		}
		// Only one LBTC allowed per chain
		jobParams.LBTCConfig = &config.LBTCConfig{
			SourceTokenAddress:           common.HexToAddress(lane.Source.Common.BridgeTokens[0].Address()),
			AttestationAPI:               api,
			AttestationAPITimeoutSeconds: 5,
		}
	}
	if !bootstrapAdded.Load() {
		bootstrapAdded.Store(true)
		err := CreateBootstrapJob(jobParams, bootstrapCommit, bootstrapExec)
		if err != nil {
			return fmt.Errorf("failed to create bootstrap job: %w", err)
		}
	}

	bootstrapCommitP2PId := bootstrapCommit.KeysBundle.P2PKeys.Data[0].Attributes.PeerID
	var p2pBootstrappersExec, p2pBootstrappersCommit *client.P2PData
	if bootstrapExec != nil {
		p2pBootstrappersExec = &client.P2PData{
			InternalIP: bootstrapExec.Node.InternalIP(),
			PeerID:     bootstrapExec.KeysBundle.P2PKeys.Data[0].Attributes.PeerID,
		}
	}

	p2pBootstrappersCommit = &client.P2PData{
		InternalIP: bootstrapCommit.Node.InternalIP(),
		PeerID:     bootstrapCommitP2PId,
	}

	jobParams.P2PV2Bootstrappers = []string{p2pBootstrappersCommit.P2PV2Bootstrapper()}

	err = SetOCR2Config(lane.Context, lane.Logger, *testConf, commitNodes, execNodes, *lane.Dest, lane.PriceReportingDisabled)
	if err != nil {
		return fmt.Errorf("failed to set ocr2 config: %w", err)
	}

	err = CreateOCR2CCIPCommitJobs(lane.Logger, jobParams, commitNodes, env.nodeMutexes, jobErrGroup)
	if err != nil {
		return fmt.Errorf("failed to create ocr2 commit jobs: %w", err)
	}
	if p2pBootstrappersExec != nil {
		jobParams.P2PV2Bootstrappers = []string{p2pBootstrappersExec.P2PV2Bootstrapper()}
	}

	err = CreateOCR2CCIPExecutionJobs(lane.Logger, jobParams, execNodes, env.nodeMutexes, jobErrGroup)
	if err != nil {
		return fmt.Errorf("failed to create ocr2 execution jobs: %w", err)
	}

	if err := lane.Source.Common.ChainClient.WaitForEvents(); err != nil {
		return fmt.Errorf("failed to wait for events: %w", err)
	}
	if err := lane.Dest.Common.ChainClient.WaitForEvents(); err != nil {
		return fmt.Errorf("failed to wait for events: %w", err)
	}
	lane.Dest.Common.ChainClient.ParallelTransactions(false)
	lane.Source.Common.ChainClient.ParallelTransactions(false)

	return nil
}

// SetOCR2Config sets the oracle config in ocr2 contracts. If execNodes is nil, commit and execution jobs are set up in same DON
func SetOCR2Config(
	ctx context.Context,
	lggr *zerolog.Logger,
	testConf testconfig.CCIPTestGroupConfig,
	commitNodes,
	execNodes []*client.CLNodesWithKeys,
	destCCIP DestCCIPModule,
	priceReportingDisabled bool,
) error {
	inflightExpiryExec := commonconfig.MustNewDuration(InflightExpiryExec)
	inflightExpiryCommit := commonconfig.MustNewDuration(InflightExpiryCommit)
	blockTime, err := destCCIP.Common.AvgBlockTime(ctx)
	if err != nil {
		return fmt.Errorf("failed to get avg block time: %w", err)
	}

	OCR2ParamsForCommit := contracts.OCR2ParamsForCommit(blockTime)
	OCR2ParamsForExec := contracts.OCR2ParamsForExec(blockTime)
	// if test config has custom ocr2 params, merge them with default params to replace default with custom ocr2 params provided in config
	// for commit and exec
	if testConf.CommitOCRParams != nil {
		err := mergo.Merge(&OCR2ParamsForCommit, testConf.CommitOCRParams, mergo.WithOverride)
		if err != nil {
			return err
		}
	}
	if testConf.ExecOCRParams != nil {
		err := mergo.Merge(&OCR2ParamsForExec, testConf.ExecOCRParams, mergo.WithOverride)
		if err != nil {
			return err
		}
	}
	lggr.Info().
		Dur("AvgBlockTimeOnDest", blockTime).
		Interface("OCRParmsForCommit", OCR2ParamsForCommit).
		Interface("OCRParmsForExec", OCR2ParamsForExec).
		Msg("Setting OCR2 config")
	commitOffchainCfg, err := contracts.NewCommitOffchainConfig(
		*commonconfig.MustNewDuration(5 * time.Second),
		1e6,
		1e6,
		*commonconfig.MustNewDuration(5 * time.Second),
		1e6,
		*inflightExpiryCommit,
		priceReportingDisabled,
	)
	if err != nil {
		return fmt.Errorf("failed to create commit offchain config: %w", err)
	}

	commitOnchainCfg, err := contracts.NewCommitOnchainConfig(
		destCCIP.Common.PriceRegistry.EthAddress,
	)
	if err != nil {
		return fmt.Errorf("failed to create commit onchain config: %w", err)
	}
	signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig, err := contracts.NewOffChainAggregatorV2ConfigForCCIPPlugin(
		commitNodes, commitOffchainCfg, commitOnchainCfg, OCR2ParamsForCommit, 3*time.Minute)
	if err != nil {
		return fmt.Errorf("failed to create ocr2 config params for commit: %w", err)
	}

	err = destCCIP.CommitStore.SetOCR2Config(signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig)
	if err != nil {
		return fmt.Errorf("failed to set ocr2 config for commit: %w", err)
	}

	nodes := commitNodes
	// if commit and exec job is set up in different DON
	if len(execNodes) > 0 {
		nodes = execNodes
	}

	// Use out of order batching strategy if we expect to be sending out of order messages
	batchingStrategyID := ccipexec.BestEffortBatchingStrategyID
	if pointer.GetBool(testConf.AllowOutOfOrder) {
		batchingStrategyID = ccipexec.ZKOverflowBatchingStrategyID
	}
	if destCCIP.OffRamp != nil {
		execOffchainCfg, err := contracts.NewExecOffchainConfig(
			1,
			BatchGasLimit,
			0.7,
			*inflightExpiryExec,
			*commonconfig.MustNewDuration(RootSnoozeTime),
			batchingStrategyID,
		)
		if err != nil {
			return fmt.Errorf("failed to create exec offchain config: %w", err)
		}
		execOnchainCfg, err := contracts.NewExecOnchainConfig(
			uint32(DefaultPermissionlessExecThreshold.Seconds()),
			destCCIP.Common.Router.EthAddress,
			destCCIP.Common.PriceRegistry.EthAddress,
			DefaultMaxNoOfTokensInMsg,
			MaxDataBytes,
			200_000,
		)
		if err != nil {
			return fmt.Errorf("failed to create exec onchain config: %w", err)
		}
		signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig, err = contracts.NewOffChainAggregatorV2ConfigForCCIPPlugin(
			nodes,
			execOffchainCfg,
			execOnchainCfg,
			OCR2ParamsForExec,
			3*time.Minute,
		)
		if err != nil {
			return fmt.Errorf("failed to create ocr2 config params for exec: %w", err)
		}
		err = destCCIP.OffRamp.SetOCR2Config(signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig)
		if err != nil {
			return fmt.Errorf("failed to set ocr2 config for exec: %w", err)
		}
	}
	return destCCIP.Common.ChainClient.WaitForEvents()
}

func CreateBootstrapJob(
	jobParams integrationtesthelpers.CCIPJobSpecParams,
	bootstrapCommit *client.CLNodesWithKeys,
	bootstrapExec *client.CLNodesWithKeys,
) error {
	_, err := bootstrapCommit.Node.MustCreateJob(jobParams.BootstrapJob(jobParams.CommitStore.Hex()))
	if err != nil {
		return fmt.Errorf("shouldn't fail creating bootstrap job on bootstrap node %w", err)
	}
	if bootstrapExec != nil {
		_, err := bootstrapExec.Node.MustCreateJob(jobParams.BootstrapJob(jobParams.OffRamp.Hex()))
		if err != nil {
			return fmt.Errorf("shouldn't fail creating bootstrap job on bootstrap node %w", err)
		}
	}
	return nil
}

func CreateOCR2CCIPCommitJobs(
	lggr *zerolog.Logger,
	jobParams integrationtesthelpers.CCIPJobSpecParams,
	commitNodes []*client.CLNodesWithKeys,
	mutexes []*sync.Mutex,
	group *errgroup.Group,
) error {
	ocr2SpecCommit, err := jobParams.CommitJobSpec()
	if err != nil {
		return fmt.Errorf("failed to create ocr2 commit job spec: %w", err)
	}
	createJob := func(index int, node *client.CLNodesWithKeys, ocr2SpecCommit client.OCR2TaskJobSpec, mu *sync.Mutex) error {
		mu.Lock()
		defer mu.Unlock()
		ocr2SpecCommit.OCR2OracleSpec.OCRKeyBundleID.SetValid(node.KeysBundle.OCR2Key.Data.ID)
		ocr2SpecCommit.OCR2OracleSpec.TransmitterID.SetValid(node.KeysBundle.EthAddress)
		lggr.Info().Msgf("Creating CCIP-Commit job on OCR node %d job name %s", index+1, ocr2SpecCommit.Name)
		_, err = node.Node.MustCreateJob(&ocr2SpecCommit)
		if err != nil {
			return fmt.Errorf("shouldn't fail creating CCIP-Commit job on OCR node %d job name %s - %w", index+1, ocr2SpecCommit.Name, err)
		}
		return nil
	}

	testSpec := client.OCR2TaskJobSpec{
		Name:           ocr2SpecCommit.Name,
		JobType:        ocr2SpecCommit.JobType,
		OCR2OracleSpec: ocr2SpecCommit.OCR2OracleSpec,
	}
	for i, node := range commitNodes {
		node := node
		i := i
		group.Go(func() error {
			return createJob(i, node, testSpec, mutexes[i])
		})
	}
	return nil
}

func CreateOCR2CCIPExecutionJobs(
	lggr *zerolog.Logger,
	jobParams integrationtesthelpers.CCIPJobSpecParams,
	execNodes []*client.CLNodesWithKeys,
	mutexes []*sync.Mutex,
	group *errgroup.Group,
) error {
	ocr2SpecExec, err := jobParams.ExecutionJobSpec()
	if err != nil {
		return fmt.Errorf("failed to create ocr2 execution job spec: %w", err)
	}
	createJob := func(index int, node *client.CLNodesWithKeys, ocr2SpecExec client.OCR2TaskJobSpec, mu *sync.Mutex) error {
		mu.Lock()
		defer mu.Unlock()
		ocr2SpecExec.OCR2OracleSpec.OCRKeyBundleID.SetValid(node.KeysBundle.OCR2Key.Data.ID)
		ocr2SpecExec.OCR2OracleSpec.TransmitterID.SetValid(node.KeysBundle.EthAddress)
		lggr.Info().Msgf("Creating CCIP-Exec job on OCR node %d job name %s", index+1, ocr2SpecExec.Name)
		_, err = node.Node.MustCreateJob(&ocr2SpecExec)
		if err != nil {
			return fmt.Errorf("shouldn't fail creating CCIP-Exec job on OCR node %d job name %s - %w", index+1,
				ocr2SpecExec.Name, err)
		}
		return nil
	}
	if ocr2SpecExec != nil {
		for i, node := range execNodes {
			node := node
			i := i
			group.Go(func() error {
				return createJob(i, node, client.OCR2TaskJobSpec{
					Name:              ocr2SpecExec.Name,
					JobType:           ocr2SpecExec.JobType,
					MaxTaskDuration:   ocr2SpecExec.MaxTaskDuration,
					ForwardingAllowed: ocr2SpecExec.ForwardingAllowed,
					OCR2OracleSpec:    ocr2SpecExec.OCR2OracleSpec,
					ObservationSource: ocr2SpecExec.ObservationSource,
				}, mutexes[i])
			})
		}
	}
	return nil
}

func TokenFeeForMultipleTokenAddr(tokenAddrToURL map[string]string) string {
	source := ""
	right := ""
	i := 1
	for addr, url := range tokenAddrToURL {
		source = source + fmt.Sprintf(`
token%d [type=http method=GET url="%s"];
token%d_parse [type=jsonparse path="data,result"];
token%d->token%d_parse;`, i, url, i, i, i)
		right = right + fmt.Sprintf(` \\\"%s\\\":$(token%d_parse),`, addr, i)
		i++
	}
	right = right[:len(right)-1]
	source = fmt.Sprintf(`%s
merge [type=merge left="{}" right="{%s}"];`, source, right)

	return source
}

type CCIPTestEnv struct {
	MockServer               *ctfClient.MockserverClient
	LocalCluster             *test_env.CLClusterTestEnv
	CLNodesWithKeys          map[string][]*client.CLNodesWithKeys // key - network chain-id
	CLNodes                  []*client.ChainlinkK8sClient
	nodeMutexes              []*sync.Mutex
	ExecNodeStartIndex       int
	CommitNodeStartIndex     int
	NumOfAllowedFaultyCommit int
	NumOfAllowedFaultyExec   int
	NumOfCommitNodes         int
	NumOfExecNodes           int
	K8Env                    *environment.Environment
	CLNodeWithKeyReady       *errgroup.Group // denotes if keys are created in chainlink node and ready to be used for job creation
}

func (c *CCIPTestEnv) ChaosLabelForGeth(t *testing.T, srcChain, destChain string) {
	err := c.K8Env.Client.LabelChaosGroupByLabels(c.K8Env.Cfg.Namespace, map[string]string{
		"app": GethLabel(srcChain),
	}, ChaosGroupNetworkACCIPGeth)
	require.NoError(t, err)

	err = c.K8Env.Client.LabelChaosGroupByLabels(c.K8Env.Cfg.Namespace, map[string]string{
		"app": GethLabel(destChain),
	}, ChaosGroupNetworkBCCIPGeth)
	require.NoError(t, err)
	gethNetworksLabels := []string{GethLabel(srcChain), GethLabel(destChain)}
	c.ChaosLabelForAllGeth(t, gethNetworksLabels)

}

func (c *CCIPTestEnv) ChaosLabelForAllGeth(t *testing.T, gethNetworksLabels []string) {
	for _, gethNetworkLabel := range gethNetworksLabels {
		err := c.K8Env.Client.AddLabel(c.K8Env.Cfg.Namespace,
			fmt.Sprintf("app=%s", gethNetworkLabel),
			fmt.Sprintf("geth=%s", ChaosGroupCCIPGeth))
		require.NoError(t, err)
	}
}

func (c *CCIPTestEnv) ChaosLabelForCLNodes(t *testing.T) {
	allowedFaulty := c.NumOfAllowedFaultyCommit
	commitStartInstance := c.CommitNodeStartIndex + 1
	execStartInstance := c.ExecNodeStartIndex + 1
	for i := commitStartInstance; i < len(c.CLNodes); i++ {
		labelSelector := map[string]string{
			"app":      "chainlink-0",
			"instance": fmt.Sprintf("node-%d", i),
		}
		if i >= commitStartInstance && i < commitStartInstance+allowedFaulty+1 {
			err := c.K8Env.Client.LabelChaosGroupByLabels(c.K8Env.Cfg.Namespace, labelSelector, ChaosGroupCommitAndExecFaultyPlus)
			require.NoError(t, err)
		}
		if i >= commitStartInstance && i < commitStartInstance+allowedFaulty {
			err := c.K8Env.Client.LabelChaosGroupByLabels(c.K8Env.Cfg.Namespace, labelSelector, ChaosGroupCommitAndExecFaulty)
			require.NoError(t, err)
		}

		// commit node starts from index 2
		if i >= commitStartInstance && i < commitStartInstance+c.NumOfCommitNodes {
			err := c.K8Env.Client.LabelChaosGroupByLabels(c.K8Env.Cfg.Namespace, labelSelector, ChaosGroupCommit)
			require.NoError(t, err)
		}
		if i >= commitStartInstance && i < commitStartInstance+c.NumOfAllowedFaultyCommit+1 {
			err := c.K8Env.Client.LabelChaosGroupByLabels(c.K8Env.Cfg.Namespace, labelSelector, ChaosGroupCommitFaultyPlus)
			require.NoError(t, err)
		}
		if i >= commitStartInstance && i < commitStartInstance+c.NumOfAllowedFaultyCommit {
			err := c.K8Env.Client.LabelChaosGroupByLabels(c.K8Env.Cfg.Namespace, labelSelector, ChaosGroupCommitFaulty)
			require.NoError(t, err)
		}
		if i >= execStartInstance && i < execStartInstance+c.NumOfExecNodes {
			err := c.K8Env.Client.LabelChaosGroupByLabels(c.K8Env.Cfg.Namespace, labelSelector, ChaosGroupExecution)
			require.NoError(t, err)
		}
		if i >= execStartInstance && i < execStartInstance+c.NumOfAllowedFaultyExec+1 {
			err := c.K8Env.Client.LabelChaosGroupByLabels(c.K8Env.Cfg.Namespace, labelSelector, ChaosGroupExecutionFaultyPlus)
			require.NoError(t, err)
		}
		if i >= execStartInstance && i < execStartInstance+c.NumOfAllowedFaultyExec {
			err := c.K8Env.Client.LabelChaosGroupByLabels(c.K8Env.Cfg.Namespace, labelSelector, ChaosGroupExecutionFaulty)
			require.NoError(t, err)
		}
	}
}

func (c *CCIPTestEnv) ConnectToExistingNodes(envConfig *testconfig.Common) error {
	if envConfig.ExistingCLCluster == nil {
		return fmt.Errorf("existing cluster is nil")
	}
	noOfNodes := pointer.GetInt(envConfig.ExistingCLCluster.NoOfNodes)
	namespace := pointer.GetString(envConfig.ExistingCLCluster.Name)

	for i := 0; i < noOfNodes; i++ {
		cfg := envConfig.ExistingCLCluster.NodeConfigs[i]
		if cfg == nil {
			return fmt.Errorf("node %d config is nil", i+1)
		}
		clClient, err := client.NewChainlinkK8sClient(cfg, cfg.InternalIP, namespace)
		if err != nil {
			return fmt.Errorf("failed to create chainlink client: %w for node %d config %v", err, i+1, cfg)
		}
		clClient.ChainlinkClient.WithRetryCount(3)
		c.CLNodes = append(c.CLNodes, clClient)
		c.nodeMutexes = append(c.nodeMutexes, &sync.Mutex{})
	}

	return nil
}

func (c *CCIPTestEnv) ConnectToDeployedNodes() error {
	if c.LocalCluster != nil {
		// for local cluster, fetch the values from the local cluster
		for _, chainlinkNode := range c.LocalCluster.ClCluster.Nodes {
			c.nodeMutexes = append(c.nodeMutexes, &sync.Mutex{})
			c.CLNodes = append(c.CLNodes, &client.ChainlinkK8sClient{
				ChainlinkClient: chainlinkNode.API.WithRetryCount(3),
			})
		}
	} else {
		// in case of k8s, we need to connect to the chainlink nodes
		log.Info().Msg("Connecting to launched resources")
		chainlinkK8sNodes, err := client.ConnectChainlinkNodes(c.K8Env)
		if err != nil {
			return fmt.Errorf("failed to connect to chainlink nodes: %w", err)
		}
		if len(chainlinkK8sNodes) == 0 {
			return fmt.Errorf("no CL node found")
		}

		for range chainlinkK8sNodes {
			c.nodeMutexes = append(c.nodeMutexes, &sync.Mutex{})
		}
		c.CLNodes = chainlinkK8sNodes
		if _, exists := c.K8Env.URLs[mockserver.InternalURLsKey]; exists {
			c.MockServer = ctfClient.ConnectMockServer(c.K8Env)
		}
	}
	return nil
}

// SetUpNodeKeysAndFund creates node keys and funds the nodes
func (c *CCIPTestEnv) SetUpNodeKeysAndFund(
	logger *zerolog.Logger,
	nodeFund *big.Float,
	chains []blockchain.EVMClient,
) error {
	if c.CLNodes == nil || len(c.CLNodes) == 0 {
		return fmt.Errorf("no chainlink nodes to setup")
	}
	var chainlinkNodes []*client.ChainlinkClient
	for _, node := range c.CLNodes {
		chainlinkNodes = append(chainlinkNodes, node.ChainlinkClient)
	}
	nodesWithKeys := make(map[string][]*client.CLNodesWithKeys)

	populateKeys := func(chain blockchain.EVMClient) error {
		log.Info().Str("chain id", chain.GetChainID().String()).Msg("creating node keys for chain")
		_, clNodes, err := client.CreateNodeKeysBundle(chainlinkNodes, "evm", chain.GetChainID().String())
		if err != nil {
			return fmt.Errorf("failed to create node keys for chain %s: %w", chain.GetChainID().String(), err)
		}
		if len(clNodes) == 0 {
			return fmt.Errorf("no CL node with keys found for chain %s", chain.GetNetworkName())
		}

		nodesWithKeys[chain.GetChainID().String()] = clNodes
		return nil
	}

	fund := func(ec blockchain.EVMClient) error {
		cfg := ec.GetNetworkConfig()
		if cfg == nil {
			return fmt.Errorf("blank network config")
		}
		c1, err := blockchain.ConcurrentEVMClient(*cfg, c.K8Env, ec, *logger)
		if err != nil {
			return fmt.Errorf("getting concurrent evmclient chain %s %w", ec.GetNetworkName(), err)
		}
		defer func() {
			if c1 != nil {
				c1.Close()
			}
		}()
		log.Info().Str("chain id", c1.GetChainID().String()).Msg("Funding Chainlink nodes for chain")
		for i := 1; i < len(chainlinkNodes); i++ {
			cl := chainlinkNodes[i]
			m := c.nodeMutexes[i]
			toAddress, err := cl.EthAddressesForChain(c1.GetChainID().String())
			if err != nil {
				return err
			}
			for _, addr := range toAddress {
				toAddr := common.HexToAddress(addr)
				gasEstimates, err := c1.EstimateGas(ethereum.CallMsg{
					To: &toAddr,
				})
				if err != nil {
					return err
				}
				m.Lock()
				err = c1.Fund(addr, nodeFund, gasEstimates)
				m.Unlock()
				if err != nil {
					return err
				}
			}
		}
		return c1.WaitForEvents()
	}
	grp, _ := errgroup.WithContext(context.Background())
	for _, chain := range chains {
		err := populateKeys(chain)
		if err != nil {
			return err
		}
	}
	for _, chain := range chains {
		chain := chain
		grp.Go(func() error {
			return fund(chain)
		})
	}
	err := grp.Wait()
	if err != nil {
		return fmt.Errorf("error funding nodes %w", err)
	}
	c.CLNodesWithKeys = nodesWithKeys

	return nil
}

func AssertBalances(t *testing.T, bas []testhelpers.BalanceAssertion) {
	logEvent := log.Info()
	for _, b := range bas {
		actual := b.Getter(t, b.Address)
		assert.NotNil(t, actual, "%v getter return nil", b.Name)
		if b.Within == "" {
			assert.Equal(t, b.Expected, actual.String(), "wrong balance for %s got %s want %s", b.Name, actual, b.Expected)
			logEvent.Interface(b.Name, struct {
				Exp    string
				Actual string
			}{
				Exp:    b.Expected,
				Actual: actual.String(),
			})
		} else {
			bi, _ := big.NewInt(0).SetString(b.Expected, 10)
			withinI, _ := big.NewInt(0).SetString(b.Within, 10)
			high := big.NewInt(0).Add(bi, withinI)
			low := big.NewInt(0).Sub(bi, withinI)
			assert.Equal(t, -1, actual.Cmp(high),
				"wrong balance for %s got %s outside expected range [%s, %s]", b.Name, actual, low, high)
			assert.Equal(t, 1, actual.Cmp(low),
				"wrong balance for %s got %s outside expected range [%s, %s]", b.Name, actual, low, high)
			logEvent.Interface(b.Name, struct {
				ExpRange string
				Actual   string
			}{
				ExpRange: fmt.Sprintf("[%s, %s]", low, high),
				Actual:   actual.String(),
			})
		}
	}
	logEvent.Msg("balance assertions succeeded")
}

type BalFunc func(ctx context.Context, addr string) (*big.Int, error)

func GetterForLinkToken(getBalance BalFunc, addr string) func(t *testing.T, _ common.Address) *big.Int {
	return func(t *testing.T, _ common.Address) *big.Int {
		balance, err := getBalance(context.Background(), addr)
		assert.NoError(t, err)
		return balance
	}
}

type BalanceItem struct {
	Address         common.Address
	Getter          func(t *testing.T, addr common.Address) *big.Int
	PreviousBalance *big.Int
	AmtToAdd        *big.Int
	AmtToSub        *big.Int
}

type BalanceSheet struct {
	mu          *sync.Mutex
	Items       map[string]BalanceItem
	PrevBalance map[string]*big.Int
}

func (b *BalanceSheet) Update(key string, item BalanceItem) {
	b.mu.Lock()
	defer b.mu.Unlock()
	prev, ok := b.Items[key]
	if !ok {
		b.Items[key] = item
		return
	}
	amtToAdd, amtToSub := big.NewInt(0), big.NewInt(0)
	if prev.AmtToAdd != nil {
		amtToAdd = prev.AmtToAdd
	}
	if prev.AmtToSub != nil {
		amtToSub = prev.AmtToSub
	}
	if item.AmtToAdd != nil {
		amtToAdd = new(big.Int).Add(amtToAdd, item.AmtToAdd)
	}
	if item.AmtToSub != nil {
		amtToSub = new(big.Int).Add(amtToSub, item.AmtToSub)
	}

	b.Items[key] = BalanceItem{
		Address:  item.Address,
		Getter:   item.Getter,
		AmtToAdd: amtToAdd,
		AmtToSub: amtToSub,
	}
}

func (b *BalanceSheet) RecordBalance(bal map[string]*big.Int) {
	b.mu.Lock()
	defer b.mu.Unlock()
	for key, value := range bal {
		if _, ok := b.PrevBalance[key]; !ok {
			b.PrevBalance[key] = value
		}
	}
}

func (b *BalanceSheet) Verify(t *testing.T) {
	var balAssertions []testhelpers.BalanceAssertion
	for key, item := range b.Items {
		prevBalance, ok := b.PrevBalance[key]
		require.Truef(t, ok, "previous balance is not captured for %s", key)
		exp := prevBalance
		if item.AmtToAdd != nil {
			exp = new(big.Int).Add(exp, item.AmtToAdd)
		}
		if item.AmtToSub != nil {
			exp = new(big.Int).Sub(exp, item.AmtToSub)
		}
		balAssertions = append(balAssertions, testhelpers.BalanceAssertion{
			Name:     key,
			Address:  item.Address,
			Getter:   item.Getter,
			Expected: exp.String(),
		})
	}
	AssertBalances(t, balAssertions)
}

func NewBalanceSheet() *BalanceSheet {
	return &BalanceSheet{
		mu:          &sync.Mutex{},
		Items:       make(map[string]BalanceItem),
		PrevBalance: make(map[string]*big.Int),
	}
}

// SetMockServerWithUSDCAttestation responds with a mock attestation for any msgHash
// The path is set with regex to match any path that starts with /v1/attestations
func SetMockServerWithUSDCAttestation(
	killGrave *ctftestenv.Killgrave,
	mockserver *ctfClient.MockserverClient,
) error {
	path := "/v1/attestations"
	response := struct {
		Status      string `json:"status"`
		Attestation string `json:"attestation"`
		Error       string `json:"error"`
	}{
		Status:      "complete",
		Attestation: "0x9049623e91719ef2aa63c55f357be2529b0e7122ae552c18aff8db58b4633c4d3920ff03d3a6d1ddf11f06bf64d7fd60d45447ac81f527ba628877dc5ca759651b08ffae25a6d3b1411749765244f0a1c131cbfe04430d687a2e12fd9d2e6dc08e118ad95d94ad832332cf3c4f7a4f3da0baa803b7be024b02db81951c0f0714de1b",
	}
	if killGrave == nil && mockserver == nil {
		return fmt.Errorf("both killgrave and mockserver are nil")
	}
	log.Info().Str("path", path).Msg("setting attestation-api response for any msgHash")
	if killGrave != nil {
		err := killGrave.SetAnyValueResponse(fmt.Sprintf("%s/{_hash:.*}", path), []string{http.MethodGet}, response)
		if err != nil {
			return fmt.Errorf("failed to set killgrave server value: %w", err)
		}
	}
	if mockserver != nil {
		err := mockserver.SetAnyValueResponse(fmt.Sprintf("%s/.*", path), response)
		if err != nil {
			return fmt.Errorf("failed to set mockserver value: %w URL = %s", err, fmt.Sprintf("%s/%s/.*", mockserver.LocalURL(), path))
		}
	}
	return nil
}

// SetMockServerWithLBTCAttestation responds with a mock attestation for any msgHash
// The path is set with regex to match any path that starts with /v1/attestations
func SetMockServerWithLBTCAttestation(
	killGrave *ctftestenv.Killgrave,
	mockserver *ctfClient.MockserverClient,
) error {
	path := "/bridge/v1/deposits/getByHash"
	type attestation struct {
		Status      string `json:"status"`
		Attestation string `json:"attestation"`
		MessageHash string `json:"message_hash"`
	}
	response := struct {
		Attestations []attestation `json:"attestations"`
	}{
		Attestations: []attestation{
			{
				MessageHash: "0x" + LBTCValidDestPoolData, //sample hash
				Status:      "NOTARIZATION_STATUS_SESSION_APPROVED",
				Attestation: "0x0000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000016000000000000000000000000000000000000000000000000000000000000000e45c70a5050000000000000000000000000000000000000000000000000000000000aa36a7000000000000000000000000845f8e3c214d8d0e4d83fc094f302aa26a12a0bc0000000000000000000000000000000000000000000000000000000000014a34000000000000000000000000845f8e3c214d8d0e4d83fc094f302aa26a12a0bc00000000000000000000000062f10ce5b727edf787ea45776bd050308a61150800000000000000000000000000000000000000000000000000000000000003e60000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000003000000000000000000000000000000000000000000000000000000000000006000000000000000000000000000000000000000000000000000000000000000c000000000000000000000000000000000000000000000000000000000000001200000000000000000000000000000000000000000000000000000000000000040277eeafba008d767c2636d9428f2ebb13ab29ac70337f4fc34b0f5606767cae546f9be3f12160de6d142e5b3c1c3ebd0bf4298662b32b597d0cc5970c7742fc10000000000000000000000000000000000000000000000000000000000000040bbcd60ecc9e06f2effe7c94161219498a1eb435b419387adadb86ec9a52dfb066ce027532517df7216404049d193a25b85c35edfa3e7c5aa4757bfe84887a3980000000000000000000000000000000000000000000000000000000000000040da4a6dc619b5ca2349783cabecc4efdbc910090d3e234d7b8d0430165f8fae532f9a965ceb85c18bb92e059adefa7ce5835850a705761ab9e026d2db4a13ef9a",
			},
		},
	}
	if killGrave == nil && mockserver == nil {
		return fmt.Errorf("both killgrave and mockserver are nil")
	}
	log.Info().Str("path", path).Msg("setting attestation-api response for any msgHash")
	if killGrave != nil {
		err := killGrave.SetAnyValueResponse(path, []string{http.MethodPost}, response)
		if err != nil {
			return fmt.Errorf("failed to set killgrave server value: %w", err)
		}
	}
	if mockserver != nil {
		err := mockserver.SetAnyValueResponse(path, response)
		if err != nil {
			return fmt.Errorf("failed to set mockserver value: %w URL = %s", err, fmt.Sprintf("%s/%s/.*", mockserver.LocalURL(), path))
		}
	}
	return nil
}

// SetMockserverWithTokenPriceValue sets the mock responses in mockserver that are read by chainlink nodes
// to simulate different price feed value.
// it keeps updating the response every 15 seconds to simulate price feed updates
func SetMockserverWithTokenPriceValue(
	killGrave *ctftestenv.Killgrave,
	mockserver *ctfClient.MockserverClient,
) {
	wg := &sync.WaitGroup{}
	path := "token_contract_"
	wg.Add(1)
	go func() {
		set := true
		// keep updating token value every 15 second
		for {
			if killGrave == nil && mockserver == nil {
				log.Fatal().Msg("both killgrave and mockserver are nil")
				return
			}
			tokenValue := big.NewInt(time.Now().UnixNano()).String()
			if killGrave != nil {
				err := killGrave.SetAdapterBasedAnyValuePath(fmt.Sprintf("%s{.*}", path), []string{http.MethodGet}, tokenValue)
				if err != nil {
					log.Fatal().Err(err).Msg("failed to set killgrave server value")
					return
				}
			}
			if mockserver != nil {
				err := mockserver.SetAnyValuePath(fmt.Sprintf("/%s.*", path), tokenValue)
				if err != nil {
					log.Fatal().Err(err).Str("URL", fmt.Sprintf("%s/%s/.*", mockserver.LocalURL(), path)).Msg("failed to set mockserver value")
					return
				}
			}
			if set {
				set = false
				wg.Done()
			}
			time.Sleep(15 * time.Second)
		}
	}()
	// wait for the first value to be set
	wg.Wait()
}

// TokenPricePipelineURLs returns the mockserver urls for the token price pipeline
func TokenPricePipelineURLs(
	tokenAddresses []string,
	killGrave *ctftestenv.Killgrave,
	mockserver *ctfClient.MockserverClient,
) map[string]string {
	mapTokenURL := make(map[string]string)

	for _, tokenAddr := range tokenAddresses {
		path := fmt.Sprintf("token_contract_%s", tokenAddr[2:12])
		if mockserver != nil {
			mapTokenURL[tokenAddr] = fmt.Sprintf("%s/%s", mockserver.Config.ClusterURL, path)
		}
		if killGrave != nil {
			mapTokenURL[tokenAddr] = fmt.Sprintf("%s/%s", killGrave.InternalEndpoint, path)
		}
	}

	return mapTokenURL
}
