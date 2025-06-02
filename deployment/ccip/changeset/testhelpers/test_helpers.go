package testhelpers

import (
	"context"
	"encoding/json"
	"fmt"
	"maps"
	"math/big"
	"net/http"
	"net/http/httptest"
	"slices"
	"sort"
	"strings"
	"testing"
	"time"

	"golang.org/x/sync/errgroup"

	cldf_solana "github.com/smartcontractkit/chainlink-deployments-framework/chain/solana"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/message_hasher"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_1/fee_quoter"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_6"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview/evm"
	commoncs "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ccipevm"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	chainsel "github.com/smartcontractkit/chain-selectors"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink-ccip/pkg/reader"
	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/environment/devenv"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"

	solconfig "github.com/smartcontractkit/chainlink-ccip/chains/solana/contracts/tests/config"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_router"
	soltokens "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/tokens"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_2_0/router"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_1/burn_mint_token_pool"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_1/usdc_token_pool"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/onramp"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/mock_ethusd_aggregator_wrapper"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/aggregator_v3_interface"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/burn_mint_erc677"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/mock_v3_aggregator_contract"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/abihelpers"

	"github.com/gagliardetto/solana-go"
)

const (
	HomeChainIndex = 0
	FeedChainIndex = 1
)

var (
	routerABI = abihelpers.MustParseABI(router.RouterABI)

	DefaultLinkPrice = deployment.E18Mult(20)
	DefaultWethPrice = deployment.E18Mult(4000)
	DefaultGasPrice  = ToPackedFee(big.NewInt(8e14), big.NewInt(0))

	OneCoin     = new(big.Int).Mul(big.NewInt(1e18), big.NewInt(1))
	TinyOneCoin = new(big.Int).SetUint64(1)
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

// ReplayLogsOption represents an option for the ReplayLogs function
type ReplayLogsOption func(*replayLogsOptions)

type replayLogsOptions struct {
	assertOnError bool
}

// WithAssertOnError configures whether ReplayLogs should assert on errors
func WithAssertOnError(assert bool) ReplayLogsOption {
	return func(opts *replayLogsOptions) {
		opts.assertOnError = assert
	}
}

// SleepAndReplay sleeps for the specified duration and then replays logs for the given chain selectors.
func SleepAndReplay(t *testing.T, env cldf.Environment, duration time.Duration, chainSelectors ...uint64) {
	time.Sleep(duration)
	replayBlocks := make(map[uint64]uint64)
	for _, selector := range chainSelectors {
		replayBlocks[selector] = 1
	}
	ReplayLogs(t, env.Offchain, replayBlocks)
}

// ReplayLogs replays logs for the given blocks using the provided offchain client.
// By default, it will assert on errors. Use WithAssertOnError(false) to change this behavior.
func ReplayLogs(t *testing.T, oc cldf.OffchainClient, replayBlocks map[uint64]uint64, opts ...ReplayLogsOption) {
	options := &replayLogsOptions{
		assertOnError: true,
	}

	for _, opt := range opts {
		opt(options)
	}

	var err error

	switch oc := oc.(type) {
	case *memory.JobClient:
		err = oc.ReplayLogs(t.Context(), replayBlocks)
	case *devenv.JobDistributor:
		err = oc.ReplayLogs(replayBlocks)
	default:
		t.Fatalf("unsupported offchain client type %T", oc)
	}

	if err != nil {
		if options.assertOnError {
			require.NoError(t, err)
		} else {
			t.Logf("failed to replay logs: %v", err)
		}
	}
}

func DeployTestContracts(t *testing.T,
	lggr logger.Logger,
	ab cldf.AddressBook,
	homeChainSel,
	feedChainSel uint64,
	chains map[uint64]cldf_evm.Chain,
	linkPrice *big.Int,
	wethPrice *big.Int,
) deployment.CapabilityRegistryConfig {
	capReg, err := cldf.DeployContract(lggr, chains[homeChainSel], ab,
		func(chain cldf_evm.Chain) cldf.ContractDeploy[*capabilities_registry.CapabilitiesRegistry] {
			crAddr, tx, cr, err2 := capabilities_registry.DeployCapabilitiesRegistry(
				chain.DeployerKey,
				chain.Client,
			)
			return cldf.ContractDeploy[*capabilities_registry.CapabilitiesRegistry]{
				Address: crAddr, Contract: cr, Tv: cldf.NewTypeAndVersion(shared.CapabilitiesRegistry, deployment.Version1_0_0), Tx: tx, Err: err2,
			}
		})
	require.NoError(t, err)

	_, err = DeployFeeds(lggr, ab, chains[feedChainSel], linkPrice, wethPrice)
	require.NoError(t, err)

	evmChainID, err := chainsel.ChainIdFromSelector(homeChainSel)
	require.NoError(t, err)

	return deployment.CapabilityRegistryConfig{
		EVMChainID:  evmChainID,
		Contract:    capReg.Address,
		NetworkType: relay.NetworkEVM,
	}
}

func LatestBlock(ctx context.Context, env cldf.Environment, chainSelector uint64) (uint64, error) {
	family, err := chainsel.GetSelectorFamily(chainSelector)
	if err != nil {
		return 0, err
	}

	switch family {
	case chainsel.FamilyEVM:
		latesthdr, err := env.BlockChains.EVMChains()[chainSelector].Client.HeaderByNumber(ctx, nil)
		if err != nil {
			return 0, errors.Wrapf(err, "failed to get latest header for chain %d", chainSelector)
		}
		block := latesthdr.Number.Uint64()
		return block, nil
	case chainsel.FamilySolana:
		return env.BlockChains.SolanaChains()[chainSelector].Client.GetSlot(ctx, solconfig.DefaultCommitment)
	default:
		return 0, errors.New("unsupported chain family")
	}
}

func LatestBlocksByChain(ctx context.Context, env cldf.Environment) (map[uint64]uint64, error) {
	latestBlocks := make(map[uint64]uint64)

	chains := []uint64{}
	chains = slices.AppendSeq(chains, maps.Keys(env.BlockChains.EVMChains()))
	chains = slices.AppendSeq(chains, maps.Keys(env.BlockChains.SolanaChains()))
	for _, selector := range chains {
		block, err := LatestBlock(ctx, env, selector)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get latest block for chain %d", selector)
		}
		latestBlocks[selector] = block
	}
	return latestBlocks, nil
}

func allocateCCIPChainSelectors(chains map[uint64]cldf_evm.Chain) (homeChainSel uint64, feeChainSel uint64) {
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

// mockAttestationResponse mocks the USDC attestation server, it returns random Attestation.
// We don't need to return exactly the same attestation, because our Mocked USDC contract doesn't rely on any specific
// value, but instead of that it just checks if the attestation is present. Therefore, it makes the test a bit simpler
// and doesn't require very detailed mocks. Please see tests in chainlink-ccip for detailed tests using real attestations
func mockAttestationResponse(isFaulty bool) *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := `{
			"status": "complete",
			"attestation": "0x9049623e91719ef2aa63c55f357be2529b0e7122ae552c18aff8db58b4633c4d3920ff03d3a6d1ddf11f06bf64d7fd60d45447ac81f527ba628877dc5ca759651b08ffae25a6d3b1411749765244f0a1c131cbfe04430d687a2e12fd9d2e6dc08e118ad95d94ad832332cf3c4f7a4f3da0baa803b7be024b02db81951c0f0714de1b"
		}`
		if isFaulty {
			response = `{
				"status": "pending",
				"error": "internal error"
			}`
		}
		_, err := w.Write([]byte(response))
		if err != nil {
			panic(err)
		}
	}))
	return server
}

func CCIPSendRequest(
	e cldf.Environment,
	state stateview.CCIPOnChainState,
	cfg *CCIPSendReqConfig,
) (*types.Transaction, uint64, error) {
	msg := cfg.Message.(router.ClientEVM2AnyMessage)
	r := state.MustGetEVMChainState(cfg.SourceChain).Router
	if cfg.IsTestRouter {
		r = state.MustGetEVMChainState(cfg.SourceChain).TestRouter
	}

	if msg.FeeToken == common.HexToAddress("0x0") { // fee is in native token
		return retryCcipSendUntilNativeFeeIsSufficient(e, r, cfg)
	}

	tx, err := r.CcipSend(cfg.Sender, cfg.DestChain, msg)
	blockNum, err := cldf.ConfirmIfNoErrorWithABI(e.BlockChains.EVMChains()[cfg.SourceChain], tx, router.RouterABI, err)
	if err != nil {
		return tx, 0, errors.Wrap(err, "failed to confirm CCIP message")
	}
	return tx, blockNum, nil
}

// retryCcipSendUntilNativeFeeIsSufficient sends a CCIP message with a native fee,
// and retries until the fee is sufficient. This is due to the fact that the fee is not known in advance,
// and the message will be rejected if the fee is insufficient.
// The function will retry based on the config's MaxRetries setting for errors other than insufficient fee.
func retryCcipSendUntilNativeFeeIsSufficient(
	e cldf.Environment,
	r *router.Router,
	cfg *CCIPSendReqConfig,
) (*types.Transaction, uint64, error) {
	const errCodeInsufficientFee = "0x07da6ee6"
	const cannotDecodeErrorReason = "could not decode error reason"
	const errMsgMissingTrieNode = "missing trie node"

	defer func() { cfg.Sender.Value = nil }()

	msg := cfg.Message.(router.ClientEVM2AnyMessage)
	var retryCount int
	for {
		fee, err := r.GetFee(&bind.CallOpts{Context: context.Background()}, cfg.DestChain, msg)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to get fee: %w", cldf.MaybeDataErr(err))
		}

		cfg.Sender.Value = fee

		tx, err := r.CcipSend(cfg.Sender, cfg.DestChain, msg)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to send CCIP message: %w", err)
		}

		blockNum, err := e.BlockChains.EVMChains()[cfg.SourceChain].Confirm(tx)
		if err != nil {
			if strings.Contains(err.Error(), errCodeInsufficientFee) {
				// Don't count insufficient fee as part of the retry count
				// because this is expected and we need to adjust the fee
				continue
			} else if strings.Contains(err.Error(), cannotDecodeErrorReason) ||
				strings.Contains(err.Error(), errMsgMissingTrieNode) {
				// If the error reason cannot be decoded, we retry to avoid transient issues. The retry behavior is disabled by default
				// It is configured in the CCIPSendReqConfig.
				// This retry was originally added to solve transient failure in end to end tests
				if retryCount >= cfg.MaxRetries {
					return nil, 0, fmt.Errorf("failed to confirm CCIP message after %d retries: %w", retryCount, cldf.MaybeDataErr(err))
				}
				retryCount++
				continue
			}

			return nil, 0, fmt.Errorf("failed to confirm CCIP message: %w", cldf.MaybeDataErr(err))
		}

		return tx, blockNum, nil
	}
}

// CCIPSendCalldata packs the calldata for the Router's ccipSend method.
// This is expected to be used in Multicall scenarios (i.e multiple ccipSend calls
// in a single transaction).
func CCIPSendCalldata(
	destChainSelector uint64,
	evm2AnyMessage router.ClientEVM2AnyMessage,
) ([]byte, error) {
	calldata, err := routerABI.Methods["ccipSend"].Inputs.Pack(
		destChainSelector,
		evm2AnyMessage,
	)
	if err != nil {
		return nil, fmt.Errorf("pack ccipSend calldata: %w", err)
	}

	calldata = append(routerABI.Methods["ccipSend"].ID, calldata...)
	return calldata, nil
}

// testhelpers.SendRequest(t, e, state, src, dest, msg, opts...)
// opts being testRouter, sender
// always return error
// note: there's also DoSendRequest vs SendRequest duplication, v1.6 vs v1.5

func TestSendRequest(
	t *testing.T,
	e cldf.Environment,
	state stateview.CCIPOnChainState,
	src, dest uint64,
	testRouter bool,
	msg any,
	opts ...SendReqOpts,
) (msgSentEvent *onramp.OnRampCCIPMessageSent) {
	baseOpts := []SendReqOpts{
		WithSourceChain(src),
		WithDestChain(dest),
		WithTestRouter(testRouter),
		WithMessage(msg),
	}
	baseOpts = append(baseOpts, opts...)

	msgSentEvent, err := SendRequest(e, state, baseOpts...)
	require.NoError(t, err)
	return msgSentEvent
}

type CCIPSendReqConfig struct {
	SourceChain  uint64
	DestChain    uint64
	IsTestRouter bool
	Sender       *bind.TransactOpts
	Message      any
	MaxRetries   int // Number of retries for errors (excluding insufficient fee errors)
	AssertError  bool
}

type SendReqOpts func(*CCIPSendReqConfig)

// WithMaxRetries sets the maximum number of retries for the CCIP send request.
func WithMaxRetries(maxRetries int) SendReqOpts {
	return func(c *CCIPSendReqConfig) {
		c.MaxRetries = maxRetries
	}
}

func WithSender(sender *bind.TransactOpts) SendReqOpts {
	return func(c *CCIPSendReqConfig) {
		c.Sender = sender
	}
}

// TODO: backwards compat, remove
func WithEvm2AnyMessage(msg router.ClientEVM2AnyMessage) SendReqOpts {
	return func(c *CCIPSendReqConfig) {
		c.Message = msg
	}
}

func WithMessage(msg any) SendReqOpts {
	return func(c *CCIPSendReqConfig) {
		c.Message = msg
	}
}

func WithTestRouter(isTestRouter bool) SendReqOpts {
	return func(c *CCIPSendReqConfig) {
		c.IsTestRouter = isTestRouter
	}
}

func WithSourceChain(sourceChain uint64) SendReqOpts {
	return func(c *CCIPSendReqConfig) {
		c.SourceChain = sourceChain
	}
}

func WithDestChain(destChain uint64) SendReqOpts {
	return func(c *CCIPSendReqConfig) {
		c.DestChain = destChain
	}
}

// SendRequest similar to TestSendRequest but returns an error.
func SendRequest(
	e cldf.Environment,
	state stateview.CCIPOnChainState,
	opts ...SendReqOpts,
) (*onramp.OnRampCCIPMessageSent, error) {
	cfg := &CCIPSendReqConfig{}
	for _, opt := range opts {
		opt(cfg)
	}
	family, err := chainsel.GetSelectorFamily(cfg.SourceChain)
	if err != nil {
		return nil, err
	}

	switch family {
	case chainsel.FamilyEVM:
		return SendRequestEVM(e, state, cfg)
	case chainsel.FamilySolana:
		return SendRequestSol(e, state, cfg)
	default:
		return nil, fmt.Errorf("send request: unsupported chain family: %v", family)
	}
}

// bytes4 public constant EVM_EXTRA_ARGS_V2_TAG = 0x181dcf10;
const GenericExtraArgsV2Tag = "0x181dcf10"
const SVMExtraArgsV1Tag = "0x1f3b3aba"

// MakeEVMExtraArgsV2 creates the extra args for the EVM2Any message that is destined
// for an EVM chain. The extra args contain the gas limit and allow out of order flag.
func MakeEVMExtraArgsV2(gasLimit uint64, allowOOO bool) []byte {
	extraArgs, err := ccipevm.SerializeClientGenericExtraArgsV2(message_hasher.ClientGenericExtraArgsV2{
		GasLimit:                 new(big.Int).SetUint64(gasLimit),
		AllowOutOfOrderExecution: allowOOO,
	})
	if err != nil {
		panic(err)
	}
	return extraArgs
}

func AddLane(
	t *testing.T,
	e *DeployedEnv,
	from, to uint64,
	isTestRouter bool,
	gasprice map[uint64]*big.Int,
	tokenPrices map[common.Address]*big.Int,
	fqCfg fee_quoter.FeeQuoterDestChainConfig,
) {
	var err error
	fromFamily, _ := chainsel.GetSelectorFamily(from)
	toFamily, _ := chainsel.GetSelectorFamily(to)
	changesets := []commoncs.ConfiguredChangeSet{}
	if fromFamily == chainsel.FamilyEVM {
		evmSrcChangesets := addEVMSrcChangesets(from, to, isTestRouter, gasprice, tokenPrices, fqCfg)
		changesets = append(changesets, evmSrcChangesets...)
	}
	if toFamily == chainsel.FamilyEVM {
		evmDstChangesets := addEVMDestChangesets(e, to, from, isTestRouter)
		changesets = append(changesets, evmDstChangesets...)
	}
	if fromFamily == chainsel.FamilySolana {
		changesets = append(changesets, addLaneSolanaChangesets(t, e, from, to, toFamily)...)
	}
	if toFamily == chainsel.FamilySolana {
		changesets = append(changesets, addLaneSolanaChangesets(t, e, to, from, fromFamily)...)
	}
	if fromFamily == chainsel.FamilyTon {
		fmt.Printf("Adding lane from %d to %d, fromFamily %s, toFamily %s\n", from, to, fromFamily, toFamily)
		changesets = append(changesets, addLaneTonChangesets(t, e, from, to, fromFamily, toFamily)...)
	}
	if toFamily == chainsel.FamilyTon {
		fmt.Printf("Adding lane from %d to %d, fromFamily %s, toFamily %s\n", to, from, toFamily, fromFamily)
		changesets = append(changesets, addLaneTonChangesets(t, e, to, from, toFamily, fromFamily)...)
	}
	e.Env, err = commoncs.ApplyChangesets(t, e.Env, e.TimelockContracts(t), changesets)
	require.NoError(t, err)
}

// RemoveLane removes a lane between the source and destination chains in the deployed environment.
func RemoveLane(t *testing.T, e *DeployedEnv, src, dest uint64, isTestRouter bool) {
	var err error
	apps := []commoncs.ConfiguredChangeSet{
		commoncs.Configure(
			cldf.CreateLegacyChangeSet(v1_6.UpdateRouterRampsChangeset),
			v1_6.UpdateRouterRampsConfig{
				UpdatesByChain: map[uint64]v1_6.RouterUpdates{
					// onRamp update on source chain
					src: {
						OnRampUpdates: map[uint64]bool{
							dest: false,
						},
					},
				},
			},
		),
		commoncs.Configure(
			cldf.CreateLegacyChangeSet(v1_6.UpdateFeeQuoterDestsChangeset),
			v1_6.UpdateFeeQuoterDestsConfig{
				UpdatesByChain: map[uint64]map[uint64]fee_quoter.FeeQuoterDestChainConfig{
					src: {
						dest: v1_6.DefaultFeeQuoterDestChainConfig(false),
					},
				},
			},
		),
		commoncs.Configure(
			cldf.CreateLegacyChangeSet(v1_6.UpdateOnRampsDestsChangeset),
			v1_6.UpdateOnRampDestsConfig{
				UpdatesByChain: map[uint64]map[uint64]v1_6.OnRampDestinationUpdate{
					src: {
						dest: {
							IsEnabled:        false,
							TestRouter:       isTestRouter,
							AllowListEnabled: false,
						},
					},
				},
			},
		),
	}
	e.Env, err = commoncs.ApplyChangesets(t, e.Env, e.TimelockContracts(t), apps)
	require.NoError(t, err)
}

func AddLaneWithDefaultPricesAndFeeQuoterConfig(t *testing.T, e *DeployedEnv, state stateview.CCIPOnChainState, from, to uint64, isTestRouter bool) {
	gasPrices := map[uint64]*big.Int{
		to: DefaultGasPrice,
	}
	fromFamily, _ := chainsel.GetSelectorFamily(from)
	tokenPrices := map[common.Address]*big.Int{}
	if fromFamily == chainsel.FamilyEVM {
		stateChainFrom := state.MustGetEVMChainState(from)
		tokenPrices = map[common.Address]*big.Int{
			stateChainFrom.LinkToken.Address(): DefaultLinkPrice,
			stateChainFrom.Weth9.Address():     DefaultWethPrice,
		}
	}
	fqCfg := v1_6.DefaultFeeQuoterDestChainConfig(true, to)
	AddLane(
		t,
		e,
		from, to,
		isTestRouter,
		gasPrices,
		tokenPrices,
		fqCfg,
	)
}

// AddLanesForAll adds densely connected lanes for all chains in the environment so that each chain
// is connected to every other chain except itself.
func AddLanesForAll(t *testing.T, e *DeployedEnv, state stateview.CCIPOnChainState) {
	chains := []uint64{}
	allEvmChainSelectors := maps.Keys(e.Env.BlockChains.EVMChains())
	allSolChainSelectors := maps.Keys(e.Env.BlockChains.SolanaChains())
	allTonChainSelectors := maps.Keys(e.Env.BlockChains.TonChains())
	chains = slices.AppendSeq(chains, allEvmChainSelectors)
	chains = slices.AppendSeq(chains, allSolChainSelectors)
	chains = slices.AppendSeq(chains, allTonChainSelectors)

	for _, source := range chains {
		for _, dest := range chains {
			if source != dest {
				AddLaneWithDefaultPricesAndFeeQuoterConfig(t, e, state, source, dest, false)
			}
		}
	}
}

func ToPackedFee(execFee, daFee *big.Int) *big.Int {
	daShifted := new(big.Int).Lsh(daFee, 112)
	return new(big.Int).Or(daShifted, execFee)
}

func DeployFeeds(
	lggr logger.Logger,
	ab cldf.AddressBook,
	chain cldf_evm.Chain,
	linkPrice *big.Int,
	wethPrice *big.Int,
) (map[string]common.Address, error) {
	linkTV := cldf.NewTypeAndVersion(shared.PriceFeed, deployment.Version1_0_0)
	mockLinkFeed := func(chain cldf_evm.Chain) cldf.ContractDeploy[*aggregator_v3_interface.AggregatorV3Interface] {
		linkFeed, tx, _, err1 := mock_v3_aggregator_contract.DeployMockV3Aggregator(
			chain.DeployerKey,
			chain.Client,
			shared.LinkDecimals, // decimals
			linkPrice,           // initialAnswer
		)
		aggregatorCr, err2 := aggregator_v3_interface.NewAggregatorV3Interface(linkFeed, chain.Client)

		return cldf.ContractDeploy[*aggregator_v3_interface.AggregatorV3Interface]{
			Address: linkFeed, Contract: aggregatorCr, Tv: linkTV, Tx: tx, Err: multierr.Append(err1, err2),
		}
	}

	mockWethFeed := func(chain cldf_evm.Chain) cldf.ContractDeploy[*aggregator_v3_interface.AggregatorV3Interface] {
		wethFeed, tx, _, err1 := mock_ethusd_aggregator_wrapper.DeployMockETHUSDAggregator(
			chain.DeployerKey,
			chain.Client,
			wethPrice, // initialAnswer
		)
		aggregatorCr, err2 := aggregator_v3_interface.NewAggregatorV3Interface(wethFeed, chain.Client)

		return cldf.ContractDeploy[*aggregator_v3_interface.AggregatorV3Interface]{
			Address: wethFeed, Contract: aggregatorCr, Tv: linkTV, Tx: tx, Err: multierr.Append(err1, err2),
		}
	}

	linkFeedAddress, linkFeedDescription, err := deploySingleFeed(lggr, ab, chain, mockLinkFeed, shared.LinkSymbol)
	if err != nil {
		return nil, err
	}

	wethFeedAddress, wethFeedDescription, err := deploySingleFeed(lggr, ab, chain, mockWethFeed, shared.WethSymbol)
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
	ab cldf.AddressBook,
	chain cldf_evm.Chain,
	deployFunc func(cldf_evm.Chain) cldf.ContractDeploy[*aggregator_v3_interface.AggregatorV3Interface],
	symbol shared.TokenSymbol,
) (common.Address, string, error) {
	// tokenTV := deployment.NewTypeAndVersion(PriceFeed, deployment.Version1_0_0)
	mockTokenFeed, err := cldf.DeployContract(lggr, chain, ab, deployFunc)
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

	if desc != shared.MockSymbolToDescription[symbol] {
		lggr.Errorw("Unexpected description for token", "symbol", symbol, "desc", desc)
		return common.Address{}, "", fmt.Errorf("unexpected description: %s", desc)
	}

	return mockTokenFeed.Address, desc, nil
}

func DeployTransferableToken(
	lggr logger.Logger,
	chains map[uint64]cldf_evm.Chain,
	src, dst uint64,
	srcActor, dstActor *bind.TransactOpts,
	state stateview.CCIPOnChainState,
	addresses cldf.AddressBook,
	token string,
) (*burn_mint_erc677.BurnMintERC677, *burn_mint_token_pool.BurnMintTokenPool, *burn_mint_erc677.BurnMintERC677, *burn_mint_token_pool.BurnMintTokenPool, error) {
	// Deploy token and pools
	srcToken, srcPool, dstToken, dstPool, err := deployTokenPoolsInParallel(lggr, chains, src, dst, srcActor, dstActor, state, addresses, token)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	// Configure pools in parallel
	configurePoolGrp := errgroup.Group{}
	configurePoolGrp.Go(func() error {
		err := setTokenPoolCounterPart(chains[src], srcPool, srcActor, dst, dstToken.Address().Bytes(), dstPool.Address().Bytes())
		if err != nil {
			return fmt.Errorf("failed to set token pool counter part chain %d: %w", src, err)
		}
		err = grantMintBurnPermissions(lggr, chains[src], srcToken, srcActor, srcPool.Address())
		if err != nil {
			return fmt.Errorf("failed to grant mint burn permissions chain %d: %w", src, err)
		}
		return nil
	})
	configurePoolGrp.Go(func() error {
		err := setTokenPoolCounterPart(chains[dst], dstPool, dstActor, src, srcToken.Address().Bytes(), srcPool.Address().Bytes())
		if err != nil {
			return fmt.Errorf("failed to set token pool counter part chain %d: %w", dst, err)
		}
		if err := grantMintBurnPermissions(lggr, chains[dst], dstToken, dstActor, dstPool.Address()); err != nil {
			return fmt.Errorf("failed to grant mint burn permissions chain %d: %w", dst, err)
		}
		return nil
	})
	if err := configurePoolGrp.Wait(); err != nil {
		return nil, nil, nil, nil, err
	}
	return srcToken, srcPool, dstToken, dstPool, nil
}

func deployTokenPoolsInParallel(
	lggr logger.Logger,
	chains map[uint64]cldf_evm.Chain,
	src, dst uint64,
	srcActor, dstActor *bind.TransactOpts,
	state stateview.CCIPOnChainState,
	addresses cldf.AddressBook,
	token string,
) (
	*burn_mint_erc677.BurnMintERC677,
	*burn_mint_token_pool.BurnMintTokenPool,
	*burn_mint_erc677.BurnMintERC677,
	*burn_mint_token_pool.BurnMintTokenPool,
	error,
) {
	deployGrp := errgroup.Group{}
	// Deploy token and pools
	var srcToken *burn_mint_erc677.BurnMintERC677
	var srcPool *burn_mint_token_pool.BurnMintTokenPool
	var dstToken *burn_mint_erc677.BurnMintERC677
	var dstPool *burn_mint_token_pool.BurnMintTokenPool

	deployGrp.Go(func() error {
		var err error
		srcToken, srcPool, err = deployTransferTokenOneEnd(lggr, chains[src], srcActor, addresses, token)
		if err != nil {
			return err
		}
		err = attachTokenToTheRegistry(chains[src], state.MustGetEVMChainState(src), srcActor, srcToken.Address(), srcPool.Address())
		return err
	})
	deployGrp.Go(func() error {
		var err error
		dstToken, dstPool, err = deployTransferTokenOneEnd(lggr, chains[dst], dstActor, addresses, token)
		if err != nil {
			return err
		}
		err = attachTokenToTheRegistry(chains[dst], state.MustGetEVMChainState(dst), dstActor, dstToken.Address(), dstPool.Address())
		return err
	})
	if err := deployGrp.Wait(); err != nil {
		return nil, nil, nil, nil, err
	}
	if srcToken == nil || srcPool == nil || dstToken == nil || dstPool == nil {
		return nil, nil, nil, nil, errors.New("failed to deploy token and pool")
	}
	return srcToken, srcPool, dstToken, dstPool, nil
}

func grantMintBurnPermissions(lggr logger.Logger, chain cldf_evm.Chain, token *burn_mint_erc677.BurnMintERC677, actor *bind.TransactOpts, address common.Address) error {
	lggr.Infow("Granting burn/mint permissions", "token", token.Address(), "address", address)
	tx, err := token.GrantMintAndBurnRoles(actor, address)
	if err != nil {
		return err
	}
	_, err = chain.Confirm(tx)
	return err
}

func setUSDCTokenPoolCounterPart(
	chain cldf_evm.Chain,
	tokenPool *usdc_token_pool.USDCTokenPool,
	destChainSelector uint64,
	actor *bind.TransactOpts,
	destTokenAddress common.Address,
	destTokenPoolAddress common.Address,
) error {
	allowedCaller := common.LeftPadBytes(destTokenPoolAddress.Bytes(), 32)
	var fixedAddr [32]byte
	copy(fixedAddr[:], allowedCaller[:32])

	domain := reader.AllAvailableDomains()[destChainSelector]

	domains := []usdc_token_pool.USDCTokenPoolDomainUpdate{
		{
			AllowedCaller:     fixedAddr,
			DomainIdentifier:  domain,
			DestChainSelector: destChainSelector,
			Enabled:           true,
		},
	}
	tx, err := tokenPool.SetDomains(chain.DeployerKey, domains)
	if err != nil {
		return err
	}

	_, err = chain.Confirm(tx)
	if err != nil {
		return err
	}

	pool, err := burn_mint_token_pool.NewBurnMintTokenPool(tokenPool.Address(), chain.Client)
	if err != nil {
		return err
	}

	return setTokenPoolCounterPart(chain, pool, actor, destChainSelector, destTokenAddress.Bytes(), destTokenPoolAddress.Bytes())
}

func setTokenPoolCounterPart(
	chain cldf_evm.Chain,
	tokenPool *burn_mint_token_pool.BurnMintTokenPool,
	actor *bind.TransactOpts,
	destChainSelector uint64,
	destTokenAddress []byte,
	destTokenPoolAddress []byte,
) error {
	tx, err := tokenPool.ApplyChainUpdates(
		actor,
		[]uint64{},
		[]burn_mint_token_pool.TokenPoolChainUpdate{
			{
				RemoteChainSelector: destChainSelector,
				RemotePoolAddresses: [][]byte{common.LeftPadBytes(destTokenPoolAddress, 32)},
				RemoteTokenAddress:  common.LeftPadBytes(destTokenAddress, 32),
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
		return fmt.Errorf("failed to apply chain updates on token pool %s: %w", tokenPool.Address(), err)
	}

	_, err = chain.Confirm(tx)
	return err
}

func attachTokenToTheRegistry(
	chain cldf_evm.Chain,
	state evm.CCIPChainState,
	owner *bind.TransactOpts,
	token common.Address,
	tokenPool common.Address,
) error {
	pool, err := state.TokenAdminRegistry.GetPool(nil, token)
	if err != nil {
		return err
	}
	// Pool is already registered, don't reattach it, because it would cause revert
	if pool != (common.Address{}) {
		return nil
	}

	for _, reg := range state.RegistryModules1_6 {
		tx, err := reg.RegisterAdminViaOwner(owner, token)
		if err != nil {
			return err
		}
		_, err = chain.Confirm(tx)
		if err != nil {
			return err
		}
	}

	tx, err := state.TokenAdminRegistry.AcceptAdminRole(owner, token)
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
	chain cldf_evm.Chain,
	deployer *bind.TransactOpts,
	addressBook cldf.AddressBook,
	tokenSymbol string,
) (*burn_mint_erc677.BurnMintERC677, *burn_mint_token_pool.BurnMintTokenPool, error) {
	var rmnAddress, routerAddress string
	chainAddresses, err := addressBook.AddressesForChain(chain.Selector)
	if err != nil {
		return nil, nil, err
	}
	for address, v := range chainAddresses {
		if cldf.NewTypeAndVersion(shared.ARMProxy, deployment.Version1_0_0).Equal(v) {
			rmnAddress = address
		}
		if cldf.NewTypeAndVersion(shared.Router, deployment.Version1_2_0).Equal(v) {
			routerAddress = address
		}
		if rmnAddress != "" && routerAddress != "" {
			break
		}
	}

	tokenDecimals := uint8(18)

	tokenContract, err := cldf.DeployContract(lggr, chain, addressBook,
		func(chain cldf_evm.Chain) cldf.ContractDeploy[*burn_mint_erc677.BurnMintERC677] {
			tokenAddress, tx, token, err2 := burn_mint_erc677.DeployBurnMintERC677(
				deployer,
				chain.Client,
				tokenSymbol,
				tokenSymbol,
				tokenDecimals,
				big.NewInt(0).Mul(big.NewInt(1e9), big.NewInt(1e18)),
			)
			return cldf.ContractDeploy[*burn_mint_erc677.BurnMintERC677]{
				Address: tokenAddress, Contract: token, Tx: tx, Tv: cldf.NewTypeAndVersion(shared.BurnMintToken, deployment.Version1_0_0), Err: err2,
			}
		})
	if err != nil {
		lggr.Errorw("Failed to deploy Token ERC677", "err", err)
		return nil, nil, err
	}

	tx, err := tokenContract.Contract.GrantMintRole(deployer, deployer.From)
	if err != nil {
		return nil, nil, err
	}
	_, err = chain.Confirm(tx)
	if err != nil {
		return nil, nil, err
	}

	tokenPool, err := cldf.DeployContract(lggr, chain, addressBook,
		func(chain cldf_evm.Chain) cldf.ContractDeploy[*burn_mint_token_pool.BurnMintTokenPool] {
			tokenPoolAddress, tx, tokenPoolContract, err2 := burn_mint_token_pool.DeployBurnMintTokenPool(
				deployer,
				chain.Client,
				tokenContract.Address,
				tokenDecimals,
				[]common.Address{},
				common.HexToAddress(rmnAddress),
				common.HexToAddress(routerAddress),
			)
			return cldf.ContractDeploy[*burn_mint_token_pool.BurnMintTokenPool]{
				Address: tokenPoolAddress, Contract: tokenPoolContract, Tx: tx, Tv: cldf.NewTypeAndVersion(shared.BurnMintTokenPool, deployment.Version1_5_1), Err: err2,
			}
		})
	if err != nil {
		lggr.Errorw("Failed to deploy token pool", "err", err)
		return nil, nil, err
	}

	return tokenContract.Contract, tokenPool.Contract, nil
}

type MintTokenInfo struct {
	auth   *bind.TransactOpts
	sender *bind.TransactOpts
	tokens []*burn_mint_erc677.BurnMintERC677
}

func NewMintTokenInfo(auth *bind.TransactOpts, tokens ...*burn_mint_erc677.BurnMintERC677) MintTokenInfo {
	return MintTokenInfo{auth: auth, tokens: tokens}
}

func NewMintTokenWithCustomSender(auth *bind.TransactOpts, sender *bind.TransactOpts, tokens ...*burn_mint_erc677.BurnMintERC677) MintTokenInfo {
	return MintTokenInfo{auth: auth, sender: sender, tokens: tokens}
}

// ApproveToken approves the router to spend the given amount of tokens
// Keeping this proxy method in order to not break compatibility
func ApproveToken(env cldf.Environment, src uint64, tokenAddress common.Address, routerAddress common.Address, amount *big.Int) error {
	return commoncs.ApproveToken(env, src, tokenAddress, routerAddress, amount)
}

// MintAndAllow mints tokens for deployers and allow router to spend them
func MintAndAllow(
	t *testing.T,
	e cldf.Environment,
	state stateview.CCIPOnChainState,
	tokenMap map[uint64][]MintTokenInfo,
) {
	configurePoolGrp := errgroup.Group{}
	tenCoins := new(big.Int).Mul(big.NewInt(1e18), big.NewInt(10))

	for chain, mintTokenInfos := range tokenMap {
		mintTokenInfos := mintTokenInfos

		configurePoolGrp.Go(func() error {
			for _, mintTokenInfo := range mintTokenInfos {
				sender := mintTokenInfo.sender
				if sender == nil {
					sender = e.BlockChains.EVMChains()[chain].DeployerKey
				}

				for _, token := range mintTokenInfo.tokens {
					tx, err := token.Mint(
						mintTokenInfo.auth,
						sender.From,
						new(big.Int).Mul(tenCoins, big.NewInt(10)),
					)
					require.NoError(t, err)
					_, err = e.BlockChains.EVMChains()[chain].Confirm(tx)
					require.NoError(t, err)

					tx, err = token.Approve(sender, state.MustGetEVMChainState(chain).Router.Address(), tenCoins)
					require.NoError(t, err)
					_, err = e.BlockChains.EVMChains()[chain].Confirm(tx)
					require.NoError(t, err)
				}
			}
			return nil
		})
	}

	require.NoError(t, configurePoolGrp.Wait())
}

func Transfer(
	ctx context.Context,
	t *testing.T,
	env cldf.Environment,
	state stateview.CCIPOnChainState,
	sourceChain, destChain uint64,
	tokens any,
	receiver []byte,
	useTestRouter bool,
	data, extraArgs []byte,
	feeToken string,
) (*onramp.OnRampCCIPMessageSent, map[uint64]*uint64) {
	startBlocks := make(map[uint64]*uint64)

	block, err := LatestBlock(ctx, env, destChain)
	require.NoError(t, err)
	startBlocks[destChain] = &block
	family, err := chainsel.GetSelectorFamily(sourceChain)
	require.NoError(t, err)

	var msg any
	switch family {
	case chainsel.FamilyEVM:
		feeTokenAddr := common.HexToAddress("0x0")
		if len(feeToken) > 0 {
			feeTokenAddr = common.HexToAddress(feeToken)
		}

		msg = router.ClientEVM2AnyMessage{
			Receiver:     common.LeftPadBytes(receiver, 32),
			Data:         data,
			TokenAmounts: tokens.([]router.ClientEVMTokenAmount),
			FeeToken:     feeTokenAddr,
			ExtraArgs:    extraArgs,
		}
	case chainsel.FamilySolana:
		feeTokenAddr := solana.PublicKey{}
		if len(feeToken) > 0 {
			feeTokenAddr, err = solana.PublicKeyFromBase58(feeToken)
			require.NoError(t, err)
		}

		msg = ccip_router.SVM2AnyMessage{
			Receiver:     common.LeftPadBytes(receiver, 32),
			Data:         data,
			TokenAmounts: tokens.([]ccip_router.SVMTokenAmount),
			FeeToken:     feeTokenAddr,
			ExtraArgs:    extraArgs,
		}

	default:
		t.Errorf("unsupported source chain: %v", family)
	}

	msgSentEvent := TestSendRequest(t, env, state, sourceChain, destChain, useTestRouter, msg)
	return msgSentEvent, startBlocks
}

type TestTransferRequest struct {
	Name                   string
	SourceChain, DestChain uint64
	Receiver               []byte
	TokenReceiver          []byte
	ExpectedStatus         int
	// optional
	Tokens                []router.ClientEVMTokenAmount
	SolTokens             []ccip_router.SVMTokenAmount
	Data                  []byte
	ExtraArgs             []byte
	ExpectedTokenBalances []ExpectedBalance
	RouterAddress         common.Address // Expected for long-living environments
	UseTestRouter         bool
	FeeToken              string
}

// TransferMultiple sends multiple CCIPMessages (represented as TestTransferRequest) sequentially.
// It verifies whether message is not reverted on the source and proper event is emitted by OnRamp.
// However, it doesn't wait for message to be committed or executed. Therefore, you can send multiple messages very fast,
// but you need to make sure they are committed/executed on your own (if that's the intention).
// It saves some time during test execution, because we let plugins batch instead of executing one by one
// If you want to wait for execution in a "batch" manner you will need to pass maps returned by TransferMultiple to
// either ConfirmMultipleCommits (for commit) or ConfirmExecWithSeqNrsForAll (for exec). Check example usage in the tests.
func TransferMultiple(
	ctx context.Context,
	t *testing.T,
	env cldf.Environment,
	state stateview.CCIPOnChainState,
	requests []TestTransferRequest,
) (
	map[uint64]*uint64,
	map[SourceDestPair]cciptypes.SeqNumRange,
	map[SourceDestPair]map[uint64]int,
	map[uint64][]ExpectedTokenBalance,
) {
	startBlocks := make(map[uint64]*uint64)
	expectedSeqNums := make(map[SourceDestPair]cciptypes.SeqNumRange)
	expectedExecutionStates := make(map[SourceDestPair]map[uint64]int)
	expectedTokenBalances := make(TokenBalanceAccumulator)

	for _, tt := range requests {
		t.Run(tt.Name, func(t *testing.T) {
			pairId := SourceDestPair{
				SourceChainSelector: tt.SourceChain,
				DestChainSelector:   tt.DestChain,
			}

			// TODO: inline this in Transfer
			family, err := chainsel.GetSelectorFamily(tt.SourceChain)
			require.NoError(t, err)
			var tokens any
			switch family {
			case chainsel.FamilyEVM:
				destFamily, err := chainsel.GetSelectorFamily(tt.DestChain)
				require.NoError(t, err)
				if destFamily == chainsel.FamilySolana {
					// for EVM2Solana token transfer we need to use tokenReceiver instead logical receiver
					expectedTokenBalances.add(tt.DestChain, tt.TokenReceiver, tt.ExpectedTokenBalances)
				} else {
					expectedTokenBalances.add(tt.DestChain, tt.Receiver, tt.ExpectedTokenBalances)
				}

				tokens = tt.Tokens

				// TODO: handle this for all chains

				// Approve router to spend tokens
				if tt.RouterAddress != (common.Address{}) {
					for _, ta := range tt.Tokens {
						err := commoncs.ApproveToken(env, tt.SourceChain, ta.Token, tt.RouterAddress, new(big.Int).Mul(ta.Amount, big.NewInt(10)))
						require.NoError(t, err)
					}
				}
			case chainsel.FamilySolana:
				tokens = tt.SolTokens
				expectedTokenBalances.add(tt.DestChain, tt.Receiver, tt.ExpectedTokenBalances)
			default:
				t.Errorf("unsupported source chain: %v", family)
			}

			msg, blocks := Transfer(
				ctx, t, env, state, tt.SourceChain, tt.DestChain, tokens, tt.Receiver, tt.UseTestRouter, tt.Data, tt.ExtraArgs, tt.FeeToken)
			if _, ok := expectedExecutionStates[pairId]; !ok {
				expectedExecutionStates[pairId] = make(map[uint64]int)
			}
			expectedExecutionStates[pairId][msg.SequenceNumber] = tt.ExpectedStatus

			if _, ok := startBlocks[tt.DestChain]; !ok {
				startBlocks[tt.DestChain] = blocks[tt.DestChain]
			}

			seqNr, ok := expectedSeqNums[pairId]
			if ok {
				expectedSeqNums[pairId] = cciptypes.NewSeqNumRange(
					seqNr.Start(), cciptypes.SeqNum(msg.SequenceNumber),
				)
			} else {
				expectedSeqNums[pairId] = cciptypes.NewSeqNumRange(
					cciptypes.SeqNum(msg.SequenceNumber), cciptypes.SeqNum(msg.SequenceNumber),
				)
			}
		})
	}

	return startBlocks, expectedSeqNums, expectedExecutionStates, expectedTokenBalances
}

// TokenBalanceAccumulator is a convenient accumulator to aggregate expected balances of different tokens
// used across the tests. You can iterate over your test cases and build the final "expected" balances for tokens (per chain, per sender)
// For instance, if your test runs multiple transfers for the same token, and you want to verify the balance of tokens at
// the end of the execution, you can simply use that struct for aggregating expected tokens
// Please also see WaitForTokenBalances to better understand how you can assert token balances
type TokenBalanceAccumulator map[uint64][]ExpectedTokenBalance

func (t TokenBalanceAccumulator) add(
	destChain uint64,
	receiver []byte,
	expectedBalances []ExpectedBalance) {
	for _, expected := range expectedBalances {
		token := expected.Token
		balance := expected.Amount
		tkIdentifier := TokenReceiverIdentifier{token, receiver}

		idx := slices.IndexFunc(t[destChain], func(b ExpectedTokenBalance) bool {
			return slices.Equal(b.Receiver.receiver, tkIdentifier.receiver) && slices.Equal(b.Receiver.token, tkIdentifier.token)
		})

		if idx < 0 {
			t[destChain] = append(t[destChain], ExpectedTokenBalance{
				Receiver: tkIdentifier,
				Amount:   balance,
			})
		} else {
			t[destChain][idx].Amount = new(big.Int).Add(t[destChain][idx].Amount, balance)
		}
	}
}

type ExpectedBalance struct {
	Token  []byte
	Amount *big.Int
}

type ExpectedTokenBalance struct {
	Receiver TokenReceiverIdentifier
	Amount   *big.Int
}
type TokenReceiverIdentifier struct {
	token    []byte
	receiver []byte
}

// WaitForTokenBalances waits for multiple ERC20 tokens to reach a particular balance
// It works in a batch manner, so you can pass and exhaustive list of different tokens (per senders and chains)
// and it would work concurrently for the balance to be met. Check WaitForTheTokenBalance to see how balance
// checking is made for a token/receiver pair
func WaitForTokenBalances(
	ctx context.Context,
	t *testing.T,
	env cldf.Environment,
	expectedBalances map[uint64][]ExpectedTokenBalance,
) {
	errGrp := &errgroup.Group{}
	for chainSelector, tokens := range expectedBalances {
		for _, expected := range tokens {
			id := expected.Receiver
			balance := expected.Amount
			errGrp.Go(func() error {
				family, err := chainsel.GetSelectorFamily(chainSelector)
				if err != nil {
					return err
				}

				switch family {
				case chainsel.FamilyEVM:
					token := common.BytesToAddress(id.token)
					receiver := common.BytesToAddress(id.receiver)
					WaitForTheTokenBalance(ctx, t, token, receiver, env.BlockChains.EVMChains()[chainSelector], balance)
				case chainsel.FamilySolana:
					expectedBalance := balance.Uint64()
					// TODO: need to pass env rather than chains
					token := solana.PublicKeyFromBytes(id.token)
					receiver := solana.PublicKeyFromBytes(id.receiver)
					// TODO: could be spl instead of spl2022
					// TODO: receiver is actually the receiver's ATA
					tokenReceiver, _, err := soltokens.FindAssociatedTokenAddress(solana.Token2022ProgramID, token, receiver)
					if err != nil {
						return err
					}
					WaitForTheTokenBalanceSol(ctx, t, token, tokenReceiver, env.BlockChains.SolanaChains()[chainSelector], expectedBalance)
				default:
				}
				return nil
			})
		}
	}
	require.NoError(t, errGrp.Wait())
}

func WaitForTheTokenBalance(
	ctx context.Context,
	t *testing.T,
	token common.Address,
	receiver common.Address,
	chain cldf_evm.Chain,
	expected *big.Int,
) {
	tokenContract, err := burn_mint_erc677.NewBurnMintERC677(token, chain.Client)
	require.NoError(t, err)

	require.Eventually(t, func() bool {
		actualBalance, err := tokenContract.BalanceOf(&bind.CallOpts{Context: ctx}, receiver)
		require.NoError(t, err)

		t.Log("Waiting for the token balance",
			"expected", expected,
			"actual", actualBalance,
			"token", token,
			"receiver", receiver,
		)

		return actualBalance.Cmp(expected) == 0
	}, tests.WaitTimeout(t), 100*time.Millisecond)
}

func WaitForTheTokenBalanceSol(
	ctx context.Context,
	t *testing.T,
	token solana.PublicKey,
	receiver solana.PublicKey,
	chain cldf_solana.Chain,
	expected uint64,
) {
	require.Eventually(t, func() bool {
		_, balance, berr := soltokens.TokenBalance(ctx, chain.Client, receiver, solconfig.DefaultCommitment)
		require.NoError(t, berr)
		// TODO: validate receiver's token mint == token

		t.Log("Waiting for the token balance",
			"expected", expected,
			"actual", balance,
			"token", token,
			"receiver", receiver,
		)
		return uint64(balance) == expected //nolint:gosec // value is always unsigned
	}, tests.WaitTimeout(t), 100*time.Millisecond)
}

func DefaultRouterMessage(receiverAddress common.Address) router.ClientEVM2AnyMessage {
	return router.ClientEVM2AnyMessage{
		Receiver:     common.LeftPadBytes(receiverAddress.Bytes(), 32),
		Data:         []byte("hello world"),
		TokenAmounts: nil,
		FeeToken:     common.HexToAddress("0x0"),
		ExtraArgs:    nil,
	}
}

func DeployCCIPContractsTest(t *testing.T, solChains int) {
	e, _ := NewMemoryEnvironment(t, WithSolChains(solChains))
	// Deploy all the CCIP contracts.
	state, err := stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)
	evmChainSelectors := e.Env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chainsel.FamilyEVM))
	solChainSelectors := e.Env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chainsel.FamilySolana))
	var allChains []uint64
	allChains = append(allChains, evmChainSelectors...)
	allChains = append(allChains, solChainSelectors...)
	snap, solana, err := state.View(&e.Env, allChains)
	require.NoError(t, err)
	if solChains > 0 {
		DeploySolanaCcipReceiver(t, e.Env)
	}

	// Assert expect every deployed address to be in the address book.
	// TODO (CCIP-3047): Add the rest of CCIPv2 representation
	b, err := json.MarshalIndent(snap, "", "	")
	require.NoError(t, err)
	fmt.Println(string(b))
	b, err = json.MarshalIndent(solana, "", "	")
	require.NoError(t, err)
	fmt.Println(string(b))
}

func TransferToTimelock(
	t *testing.T,
	tenv DeployedEnv,
	state stateview.CCIPOnChainState,
	chains []uint64,
	withTestRouterTransfer bool,
) {
	timelockContracts := make(map[uint64]*proposalutils.TimelockExecutionContracts, len(chains)+1)
	for _, chain := range chains {
		timelockContracts[chain] = &proposalutils.TimelockExecutionContracts{
			Timelock:  state.MustGetEVMChainState(chain).Timelock,
			CallProxy: state.MustGetEVMChainState(chain).CallProxy,
		}
	}
	// Add the home chain to the timelock contracts.
	timelockContracts[tenv.HomeChainSel] = &proposalutils.TimelockExecutionContracts{
		Timelock:  state.MustGetEVMChainState(tenv.HomeChainSel).Timelock,
		CallProxy: state.MustGetEVMChainState(tenv.HomeChainSel).CallProxy,
	}
	// Transfer ownership to timelock so that we can promote the zero digest later down the line.
	_, err := commoncs.Apply(t, tenv.Env,
		timelockContracts,
		commoncs.Configure(
			cldf.CreateLegacyChangeSet(commoncs.TransferToMCMSWithTimelockV2),
			GenTestTransferOwnershipConfig(tenv, chains, state, withTestRouterTransfer),
		),
	)
	require.NoError(t, err)
	AssertTimelockOwnership(t, tenv, chains, state, withTestRouterTransfer)
}
