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

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/fee_quoter"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/message_hasher"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	ccipChangeSetSolana "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/solana"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_6"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview/evm"
	solanastateview "github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview/solana"

	commoncs "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ccipevm"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	chainsel "github.com/smartcontractkit/chain-selectors"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"

	"github.com/smartcontractkit/chainlink-ccip/pkg/reader"
	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/environment/devenv"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"

	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/base_token_pool"
	solCommon "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_common"
	solOffRamp "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_offramp"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_router"
	solRouter "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_router"
	solFeeQuoter "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/fee_quoter"
	solRmnRemote "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/rmn_remote"
	solTestReceiver "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/test_ccip_receiver"
	solTestTokenPool "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/test_token_pool"

	solconfig "github.com/smartcontractkit/chainlink-ccip/chains/solana/contracts/tests/config"
	solccip "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/ccip"
	solcommon "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/common"
	solstate "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/state"
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

	solbinary "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/token"
	"github.com/gagliardetto/solana-go/rpc"
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

func SendRequestEVM(
	e cldf.Environment,
	state stateview.CCIPOnChainState,
	cfg *CCIPSendReqConfig,
) (*onramp.OnRampCCIPMessageSent, error) {
	// Set default sender if not provided
	if cfg.Sender == nil {
		cfg.Sender = e.BlockChains.EVMChains()[cfg.SourceChain].DeployerKey
	}

	e.Logger.Infof("Sending CCIP request from chain selector %d to chain selector %d from sender %s",
		cfg.SourceChain, cfg.DestChain, cfg.Sender.From.String())

	tx, blockNum, err := CCIPSendRequest(e, state, cfg)
	if err != nil {
		return nil, err
	}

	it, err := state.MustGetEVMChainState(cfg.SourceChain).OnRamp.FilterCCIPMessageSent(&bind.FilterOpts{
		Start:   blockNum,
		End:     &blockNum,
		Context: context.Background(),
	}, []uint64{cfg.DestChain}, []uint64{})
	if err != nil {
		return nil, err
	}

	if !it.Next() {
		return nil, errors.New("no CCIP message sent event found")
	}

	e.Logger.Infof("CCIP message (id %s) sent from chain selector %d to chain selector %d tx %s seqNum %d nonce %d sender %s testRouterEnabled %t",
		common.Bytes2Hex(it.Event.Message.Header.MessageId[:]),
		cfg.SourceChain,
		cfg.DestChain,
		tx.Hash().String(),
		it.Event.SequenceNumber,
		it.Event.Message.Header.Nonce,
		it.Event.Message.Sender.String(),
		cfg.IsTestRouter,
	)
	return it.Event, nil
}

func SendRequestSol(
	e cldf.Environment,
	state stateview.CCIPOnChainState,
	cfg *CCIPSendReqConfig,
) (*onramp.OnRampCCIPMessageSent, error) { // TODO: chain independent return value
	ctx := e.GetContext()

	s := state.SolChains[cfg.SourceChain]
	c := e.BlockChains.SolanaChains()[cfg.SourceChain]

	destinationChainSelector := cfg.DestChain
	message := cfg.Message.(ccip_router.SVM2AnyMessage)
	feeToken := message.FeeToken
	client := c.Client

	// TODO: sender from cfg is EVM specific - need to revisit for Solana
	sender := c.DeployerKey

	e.Logger.Infof("Sending CCIP request from chain selector %d to chain selector %d from sender %s",
		cfg.SourceChain, cfg.DestChain, sender.PublicKey().String())

	feeTokenProgramID := solana.TokenProgramID
	feeTokenUserATA := solana.PublicKey{}
	if feeToken.IsZero() {
		// If the fee token is native SOL (i.e. message.FeeToken is the zero address), then we will
		// leave message.FeeToken as it is, but specify the WSOL mint account in the accounts list
		feeToken = solana.SolMint
	} else {
		feeTokenInfo, err := client.GetAccountInfo(ctx, feeToken)
		if err != nil {
			return nil, err
		}
		feeTokenProgramID = feeTokenInfo.Value.Owner

		_, err = GetSolanaTokenMintInfo(feeTokenInfo)
		if err != nil {
			return nil, fmt.Errorf("the provided fee token is not a valid token: (err = %w)", err)
		}

		ata, _, err := soltokens.FindAssociatedTokenAddress(feeTokenProgramID, feeToken, sender.PublicKey())
		if err != nil {
			return nil, err
		}

		feeTokenUserATA = ata
	}

	destinationChainStatePDA, err := solstate.FindDestChainStatePDA(destinationChainSelector, s.Router)
	if err != nil {
		return nil, err
	}

	noncePDA, err := solstate.FindNoncePDA(cfg.DestChain, sender.PublicKey(), s.Router)
	if err != nil {
		return nil, err
	}

	linkFqBillingConfigPDA, _, err := solstate.FindFqBillingTokenConfigPDA(s.LinkToken, s.FeeQuoter)
	if err != nil {
		return nil, err
	}

	feeTokenFqBillingConfigPDA, _, err := solstate.FindFqBillingTokenConfigPDA(feeToken, s.FeeQuoter)
	if err != nil {
		return nil, err
	}

	billingSignerPDA, _, err := solstate.FindFeeBillingSignerPDA(s.Router)
	if err != nil {
		return nil, err
	}

	feeTokenReceiverATA, _, err := soltokens.FindAssociatedTokenAddress(feeTokenProgramID, feeToken, billingSignerPDA)
	if err != nil {
		return nil, err
	}

	fqDestChainPDA, _, err := solstate.FindFqDestChainPDA(cfg.DestChain, s.FeeQuoter)
	if err != nil {
		return nil, err
	}

	rmnRemoteCursesPDA, _, err := solstate.FindRMNRemoteCursesPDA(s.RMNRemote)
	if err != nil {
		return nil, err
	}

	base := ccip_router.NewCcipSendInstruction(
		destinationChainSelector,
		message,
		[]byte{}, // starting indices for accounts, calculated later
		s.RouterConfigPDA,
		destinationChainStatePDA,
		noncePDA,
		sender.PublicKey(),
		solana.SystemProgramID,
		feeTokenProgramID,
		feeToken,
		feeTokenUserATA,
		feeTokenReceiverATA,
		billingSignerPDA,
		s.FeeQuoter,
		s.FeeQuoterConfigPDA,
		fqDestChainPDA,
		feeTokenFqBillingConfigPDA,
		linkFqBillingConfigPDA,
		s.RMNRemote,
		rmnRemoteCursesPDA,
		s.RMNRemoteConfigPDA,
	)

	// When paying with a non-native token (i.e. any SPL token), the user ATA must be writable so we
	// can debit the fees. If paying with native SOL, then the ATA passed in is just a zero-address
	// placeholder, and that can't be marked as writable.
	if !feeTokenUserATA.IsZero() {
		base.GetFeeTokenUserAssociatedAccountAccount().WRITE()
	}

	addressTables := map[solana.PublicKey]solana.PublicKeySlice{}

	requiredAccounts := len(base.AccountMetaSlice)
	tokenIndexes := []byte{}

	// set config.FeeQuoterProgram and CcipRouterProgram since they point to wrong addresses
	solconfig.FeeQuoterProgram = s.FeeQuoter
	solconfig.CcipRouterProgram = s.Router

	// Append token accounts to the account metas
	for _, tokenAmount := range message.TokenAmounts {
		tokenPubKey := tokenAmount.Token

		allTokenPools := solana.PublicKeySlice{}
		allTokenPools = slices.AppendSeq(allTokenPools, maps.Values(s.LockReleaseTokenPools))
		allTokenPools = slices.AppendSeq(allTokenPools, maps.Values(s.BurnMintTokenPools))

		e.Logger.Infof("Found %d token pools in state - searching for matching token pool", len(allTokenPools))
		tokenPoolPubKey, err := MatchTokenToTokenPool(ctx, client, tokenPubKey, allTokenPools)
		if err != nil {
			return nil, err
		}

		e.Logger.Infof("Token '%s' was matched to token pool '%s'",
			tokenPubKey.String(),
			tokenPoolPubKey.String(),
		)

		tokenProgramID, err := InferSolanaTokenProgramID(ctx, client, tokenPubKey)
		if err != nil {
			return nil, err
		}

		tokenPool, err := soltokens.NewTokenPool(tokenProgramID, tokenPoolPubKey, tokenPubKey)
		if err != nil {
			return nil, err
		}

		// Set the token pool's lookup table address
		var tokenAdminRegistry solCommon.TokenAdminRegistry
		err = solcommon.GetAccountDataBorshInto(ctx, client, tokenPool.AdminRegistryPDA, solconfig.DefaultCommitment, &tokenAdminRegistry)
		if err != nil {
			return nil, err
		}

		tokenPool.PoolLookupTable = tokenAdminRegistry.LookupTable

		// invalid config account, maybe this billing stuff isn't right

		chainPDA, _, err := soltokens.TokenPoolChainConfigPDA(cfg.DestChain, tokenPubKey, tokenPoolPubKey)
		if err != nil {
			return nil, err
		}

		tokenPool.Chain[cfg.DestChain] = chainPDA

		billingPDA, _, err := solstate.FindFqPerChainPerTokenConfigPDA(cfg.DestChain, tokenPubKey, s.FeeQuoter)
		if err != nil {
			return nil, err
		}

		tokenPool.Billing[cfg.DestChain] = billingPDA

		userTokenAccount, _, err := soltokens.FindAssociatedTokenAddress(tokenProgramID, tokenPubKey, sender.PublicKey())
		if err != nil {
			return nil, err
		}

		tokenMetas, tokenAddressTables, err := soltokens.ParseTokenLookupTableWithChain(ctx, client, tokenPool, userTokenAccount, cfg.DestChain)
		if err != nil {
			return nil, err
		}

		tokenIndexes = append(tokenIndexes, byte(len(base.AccountMetaSlice)-requiredAccounts))
		base.AccountMetaSlice = append(base.AccountMetaSlice, tokenMetas...)
		maps.Copy(addressTables, tokenAddressTables)
	}

	base.SetTokenIndexes(tokenIndexes)

	ix, err := base.ValidateAndBuild()
	if err != nil {
		return nil, err
	}

	// for some reason onchain doesn't see extraAccounts

	ixs := []solana.Instruction{ix}
	result, err := solcommon.SendAndConfirmWithLookupTables(ctx, client, ixs, *sender, solconfig.DefaultCommitment, addressTables, solcommon.AddComputeUnitLimit(400_000))
	if err != nil {
		return nil, err
	}

	// check CCIP event
	ccipMessageSentEvent := solccip.EventCCIPMessageSent{}
	printEvents := true
	err = solcommon.ParseEvent(result.Meta.LogMessages, "CCIPMessageSent", &ccipMessageSentEvent, printEvents)
	if err != nil {
		return nil, err
	}

	if len(message.TokenAmounts) != len(ccipMessageSentEvent.Message.TokenAmounts) {
		return nil, errors.New("token amounts mismatch")
	}

	// TODO: fee bumping?

	transactionID := "N/A"
	if tx, err := result.Transaction.GetTransaction(); err != nil {
		e.Logger.Warnf("could not obtain transaction details (err = %s)", err.Error())
	} else if len(tx.Signatures) == 0 {
		e.Logger.Warnf("transaction has no signatures: %v", tx)
	} else {
		transactionID = tx.Signatures[0].String()
	}

	e.Logger.Infof("CCIP message (id %s) sent from chain selector %d to chain selector %d tx %s seqNum %d nonce %d sender %s testRouterEnabled %t",
		common.Bytes2Hex(ccipMessageSentEvent.Message.Header.MessageId[:]),
		cfg.SourceChain,
		cfg.DestChain,
		transactionID,
		ccipMessageSentEvent.SequenceNumber,
		ccipMessageSentEvent.Message.Header.Nonce,
		ccipMessageSentEvent.Message.Sender.String(),
		cfg.IsTestRouter,
	)

	return &onramp.OnRampCCIPMessageSent{
		DestChainSelector: ccipMessageSentEvent.DestinationChainSelector,
		SequenceNumber:    ccipMessageSentEvent.SequenceNumber,
		Message: onramp.InternalEVM2AnyRampMessage{
			Header: onramp.InternalRampMessageHeader{
				SourceChainSelector: ccipMessageSentEvent.Message.Header.SourceChainSelector,
				DestChainSelector:   ccipMessageSentEvent.Message.Header.DestChainSelector,
				MessageId:           ccipMessageSentEvent.Message.Header.MessageId,
				SequenceNumber:      ccipMessageSentEvent.SequenceNumber,
				Nonce:               ccipMessageSentEvent.Message.Header.Nonce,
			},
			FeeTokenAmount: ConvertSolanaCrossChainAmountToBigInt(ccipMessageSentEvent.Message.FeeTokenAmount),
			FeeValueJuels:  ConvertSolanaCrossChainAmountToBigInt(ccipMessageSentEvent.Message.FeeValueJuels),
			ExtraArgs:      ccipMessageSentEvent.Message.ExtraArgs,
			Receiver:       ccipMessageSentEvent.Message.Receiver,
			Data:           ccipMessageSentEvent.Message.Data,

			// TODO: these fields are EVM specific - need to revisit for Solana
			FeeToken:     common.Address{}, // ccipMessageSentEvent.Message.FeeToken
			Sender:       common.Address{}, // ccipMessageSentEvent.Message.Sender
			TokenAmounts: []onramp.InternalEVM2AnyTokenTransfer{},
		},

		// TODO: EVM specific - need to revisit for Solana
		Raw: types.Log{},
	}, nil
}

func ConvertSolanaCrossChainAmountToBigInt(amount ccip_router.CrossChainAmount) *big.Int {
	bytes := amount.LeBytes[:]
	slices.Reverse(bytes) // convert to big-endian
	return big.NewInt(0).SetBytes(bytes)
}

func InferSolanaTokenProgramID(ctx context.Context, client *rpc.Client, tokenPubKey solana.PublicKey) (solana.PublicKey, error) {
	tokenAcctInfo, err := client.GetAccountInfo(ctx, tokenPubKey)
	if errors.Is(err, rpc.ErrNotFound) {
		// NOTE: we use a fallback value of Token2022ProgramID to maintain backwards compatibility with the Solana tests
		return solana.Token2022ProgramID, nil
	}
	if err != nil {
		return solana.PublicKey{}, err
	}

	_, err = GetSolanaTokenMintInfo(tokenAcctInfo)
	if err != nil {
		return solana.PublicKey{}, fmt.Errorf("expected '%s' to be a token public key: (err = %w)", tokenPubKey, err)
	}

	return tokenAcctInfo.Value.Owner, nil
}

func GetSolanaTokenMintInfo(tokenAcctInfo *rpc.GetAccountInfoResult) (token.Mint, error) {
	var mint token.Mint

	err := solbinary.NewBinDecoder(tokenAcctInfo.Bytes()).Decode(&mint)
	if err != nil {
		return token.Mint{}, fmt.Errorf("failed to decode token mint data: (err = %w)", err)
	}

	return mint, nil
}

func MatchTokenToTokenPool(ctx context.Context, client *rpc.Client, tokenPubKey solana.PublicKey, tokenPoolPubKeys solana.PublicKeySlice) (solana.PublicKey, error) {
	for _, tokenPoolPubKey := range tokenPoolPubKeys {
		tokenPoolConfigAddress, err := soltokens.TokenPoolConfigAddress(tokenPubKey, tokenPoolPubKey)
		if err != nil {
			return solana.PublicKey{}, err
		}

		var tokenPoolConfig base_token_pool.BaseConfig
		err = solcommon.GetAccountDataBorshInto(ctx, client, tokenPoolConfigAddress, solconfig.DefaultCommitment, &tokenPoolConfig)
		if errors.Is(err, rpc.ErrNotFound) {
			continue
		}
		if err != nil {
			return solana.PublicKey{}, err
		}

		return tokenPoolPubKey, nil
	}

	tokenPoolPubKeyStrs := make([]string, len(tokenPoolPubKeys))
	for i, tokenPoolPubKey := range tokenPoolPubKeys {
		tokenPoolPubKeyStrs[i] = "'" + tokenPoolPubKey.String() + "'"
	}

	msg := "token with public key '%s' is not associated with any of the following token pools: [ %s ]"
	return solana.PublicKey{}, fmt.Errorf(msg, tokenPubKey.String(), strings.Join(tokenPoolPubKeyStrs, ", "))
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
		evmSrcChangesets := AddEVMSrcChangesets(from, to, isTestRouter, gasprice, tokenPrices, fqCfg)
		changesets = append(changesets, evmSrcChangesets...)
	}
	if toFamily == chainsel.FamilyEVM {
		evmDstChangesets := AddEVMDestChangesets(e, to, from, isTestRouter)
		changesets = append(changesets, evmDstChangesets...)
	}
	if fromFamily == chainsel.FamilySolana {
		changesets = append(changesets, AddLaneSolanaChangesets(t, e, from, to, toFamily)...)
	}
	if toFamily == chainsel.FamilySolana {
		changesets = append(changesets, AddLaneSolanaChangesets(t, e, to, from, fromFamily)...)
	}

	e.Env, err = commoncs.ApplyChangesets(t, e.Env, e.TimelockContracts(t), changesets)
	require.NoError(t, err)
}

func AddLaneSolanaChangesets(t *testing.T, e *DeployedEnv, solChainSelector, remoteChainSelector uint64, remoteFamily string) []commoncs.ConfiguredChangeSet {
	chainFamilySelector := [4]uint8{}
	if remoteFamily == chainsel.FamilyEVM {
		// bytes4(keccak256("CCIP ChainFamilySelector EVM"))
		chainFamilySelector = [4]uint8{40, 18, 213, 44}
	} else if remoteFamily == chainsel.FamilySolana {
		// bytes4(keccak256("CCIP ChainFamilySelector SVM"));
		chainFamilySelector = [4]uint8{30, 16, 189, 196}
	}
	solanaChangesets := []commoncs.ConfiguredChangeSet{
		commoncs.Configure(
			cldf.CreateLegacyChangeSet(ccipChangeSetSolana.AddRemoteChainToRouter),
			ccipChangeSetSolana.AddRemoteChainToRouterConfig{
				ChainSelector: solChainSelector,
				UpdatesByChain: map[uint64]*ccipChangeSetSolana.RouterConfig{
					remoteChainSelector: {
						RouterDestinationConfig: solRouter.DestChainConfig{
							AllowListEnabled: true,
							AllowedSenders:   []solana.PublicKey{e.Env.BlockChains.SolanaChains()[solChainSelector].DeployerKey.PublicKey()},
						},
					},
				},
			},
		),
		commoncs.Configure(
			cldf.CreateLegacyChangeSet(ccipChangeSetSolana.AddRemoteChainToFeeQuoter),
			ccipChangeSetSolana.AddRemoteChainToFeeQuoterConfig{
				ChainSelector: solChainSelector,
				UpdatesByChain: map[uint64]*ccipChangeSetSolana.FeeQuoterConfig{
					remoteChainSelector: {
						FeeQuoterDestinationConfig: solFeeQuoter.DestChainConfig{
							IsEnabled:                   true,
							DefaultTxGasLimit:           200000,
							MaxPerMsgGasLimit:           3000000,
							MaxDataBytes:                30000,
							MaxNumberOfTokensPerMsg:     5,
							DefaultTokenDestGasOverhead: 90000,
							DestGasOverhead:             90000,
							ChainFamilySelector:         chainFamilySelector,
						},
					},
				},
			},
		),
		commoncs.Configure(
			cldf.CreateLegacyChangeSet(ccipChangeSetSolana.AddRemoteChainToOffRamp),
			ccipChangeSetSolana.AddRemoteChainToOffRampConfig{
				ChainSelector: solChainSelector,
				UpdatesByChain: map[uint64]*ccipChangeSetSolana.OffRampConfig{
					remoteChainSelector: {
						EnabledAsSource: true,
					},
				},
			},
		),
	}
	return solanaChangesets
}

func AddEVMSrcChangesets(from, to uint64, isTestRouter bool, gasprice map[uint64]*big.Int, tokenPrices map[common.Address]*big.Int, fqCfg fee_quoter.FeeQuoterDestChainConfig) []commoncs.ConfiguredChangeSet {
	evmSrcChangesets := []commoncs.ConfiguredChangeSet{
		commoncs.Configure(
			cldf.CreateLegacyChangeSet(v1_6.UpdateOnRampsDestsChangeset),
			v1_6.UpdateOnRampDestsConfig{
				UpdatesByChain: map[uint64]map[uint64]v1_6.OnRampDestinationUpdate{
					from: {
						to: {
							IsEnabled:        true,
							TestRouter:       isTestRouter,
							AllowListEnabled: false,
						},
					},
				},
			},
		),
		commoncs.Configure(
			cldf.CreateLegacyChangeSet(v1_6.UpdateFeeQuoterPricesChangeset),
			v1_6.UpdateFeeQuoterPricesConfig{
				PricesByChain: map[uint64]v1_6.FeeQuoterPriceUpdatePerSource{
					from: {
						TokenPrices: tokenPrices,
						GasPrices:   gasprice,
					},
				},
			},
		),
		commoncs.Configure(
			cldf.CreateLegacyChangeSet(v1_6.UpdateFeeQuoterDestsChangeset),
			v1_6.UpdateFeeQuoterDestsConfig{
				UpdatesByChain: map[uint64]map[uint64]fee_quoter.FeeQuoterDestChainConfig{
					from: {
						to: fqCfg,
					},
				},
			},
		),
		commoncs.Configure(
			cldf.CreateLegacyChangeSet(v1_6.UpdateRouterRampsChangeset),
			v1_6.UpdateRouterRampsConfig{
				TestRouter: isTestRouter,
				UpdatesByChain: map[uint64]v1_6.RouterUpdates{
					// onRamp update on source chain
					from: {
						OnRampUpdates: map[uint64]bool{
							to: true,
						},
					},
				},
			},
		),
	}
	return evmSrcChangesets
}

func AddEVMDestChangesets(e *DeployedEnv, to, from uint64, isTestRouter bool) []commoncs.ConfiguredChangeSet {
	evmDstChangesets := []commoncs.ConfiguredChangeSet{
		commoncs.Configure(
			cldf.CreateLegacyChangeSet(v1_6.UpdateOffRampSourcesChangeset),
			v1_6.UpdateOffRampSourcesConfig{
				UpdatesByChain: map[uint64]map[uint64]v1_6.OffRampSourceUpdate{
					to: {
						from: {
							IsEnabled:                 true,
							TestRouter:                isTestRouter,
							IsRMNVerificationDisabled: !e.RmnEnabledSourceChains[from],
						},
					},
				},
			},
		),
		commoncs.Configure(
			cldf.CreateLegacyChangeSet(v1_6.UpdateRouterRampsChangeset),
			v1_6.UpdateRouterRampsConfig{
				TestRouter: isTestRouter,
				UpdatesByChain: map[uint64]v1_6.RouterUpdates{
					// offramp update on dest chain
					to: {
						OffRampUpdates: map[uint64]bool{
							from: true,
						},
					},
				},
			},
		),
	}
	return evmDstChangesets
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
	chains = slices.AppendSeq(chains, allEvmChainSelectors)
	chains = slices.AppendSeq(chains, allSolChainSelectors)

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

// assuming one out of the src and dst is solana and the other is evm
func DeployTransferableTokenSolana(
	lggr logger.Logger,
	e cldf.Environment,
	evmChainSel, solChainSel uint64,
	evmDeployer *bind.TransactOpts,
	evmTokenName string,
) (*burn_mint_erc677.BurnMintERC677,
	*burn_mint_token_pool.BurnMintTokenPool, solana.PublicKey, error) {
	selectorFamily, err := chainsel.GetSelectorFamily(evmChainSel)
	if err != nil {
		return nil, nil, solana.PublicKey{}, err
	}
	if selectorFamily != chainsel.FamilyEVM {
		return nil, nil, solana.PublicKey{}, fmt.Errorf("evmChainSel %d is not an evm chain", evmChainSel)
	}
	selectorFamily, err = chainsel.GetSelectorFamily(solChainSel)
	if err != nil {
		return nil, nil, solana.PublicKey{}, err
	}
	if selectorFamily != chainsel.FamilySolana {
		return nil, nil, solana.PublicKey{}, fmt.Errorf("solChainSel %d is not a solana chain", solChainSel)
	}
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return nil, nil, solana.PublicKey{}, err
	}

	addresses := e.ExistingAddresses //nolint:staticcheck // addressbook still valid
	// deploy evm token and pool
	evmToken, evmPool, err := deployTransferTokenOneEnd(lggr, e.BlockChains.EVMChains()[evmChainSel], evmDeployer, addresses, evmTokenName)
	if err != nil {
		return nil, nil, solana.PublicKey{}, err
	}
	// attach token and pool to the registry
	if err := attachTokenToTheRegistry(e.BlockChains.EVMChains()[evmChainSel], state.MustGetEVMChainState(evmChainSel), evmDeployer, evmToken.Address(), evmPool.Address()); err != nil {
		return nil, nil, solana.PublicKey{}, err
	}
	solDeployerKey := e.BlockChains.SolanaChains()[solChainSel].DeployerKey.PublicKey()

	// deploy solana token
	solTokenName := evmTokenName
	e, err = commoncs.Apply(nil, e, nil,
		commoncs.Configure(
			// this makes the deployer the mint authority by default
			cldf.CreateLegacyChangeSet(ccipChangeSetSolana.DeploySolanaToken),
			ccipChangeSetSolana.DeploySolanaTokenConfig{
				ChainSelector:    solChainSel,
				TokenProgramName: shared.SPL2022Tokens,
				TokenDecimals:    9,
				TokenSymbol:      solTokenName,
				ATAList:          []string{solDeployerKey.String()},
				MintAmountToAddress: map[string]uint64{
					solDeployerKey.String(): uint64(1000e9),
				},
			},
		),
	)
	if err != nil {
		return nil, nil, solana.PublicKey{}, err
	}
	// find solana token address
	solAddresses, err := e.ExistingAddresses.AddressesForChain(solChainSel)
	if err != nil {
		return nil, nil, solana.PublicKey{}, err
	}
	solTokenAddress := solanastateview.FindSolanaAddress(
		cldf.TypeAndVersion{
			Type:    shared.SPL2022Tokens,
			Version: deployment.Version1_0_0,
			Labels:  cldf.NewLabelSet(solTokenName),
		},
		solAddresses,
	)
	bnm := solTestTokenPool.BurnAndMint_PoolType

	// deploy and configure solana token pool
	e, err = commoncs.Apply(nil, e, nil,
		commoncs.Configure(
			// deploy token pool and set the burn/mint authority to the tokenPool
			cldf.CreateLegacyChangeSet(ccipChangeSetSolana.E2ETokenPool),
			ccipChangeSetSolana.E2ETokenPoolConfig{
				AddTokenPoolAndLookupTable: []ccipChangeSetSolana.TokenPoolConfig{
					{
						ChainSelector: solChainSel,
						TokenPubKey:   solTokenAddress,
						PoolType:      &bnm,
						Metadata:      shared.CLLMetadata,
					},
				},
				RegisterTokenAdminRegistry: []ccipChangeSetSolana.RegisterTokenAdminRegistryConfig{
					{
						ChainSelector:           solChainSel,
						TokenPubKey:             solTokenAddress,
						TokenAdminRegistryAdmin: solDeployerKey.String(),
						RegisterType:            ccipChangeSetSolana.ViaGetCcipAdminInstruction,
					},
				},
				AcceptAdminRoleTokenAdminRegistry: []ccipChangeSetSolana.AcceptAdminRoleTokenAdminRegistryConfig{
					{
						ChainSelector: solChainSel,
						TokenPubKey:   solTokenAddress,
					},
				},
				SetPool: []ccipChangeSetSolana.SetPoolConfig{
					{
						ChainSelector:   solChainSel,
						TokenPubKey:     solTokenAddress,
						PoolType:        &bnm,
						Metadata:        shared.CLLMetadata,
						WritableIndexes: []uint8{3, 4, 7},
					},
				},
				RemoteChainTokenPool: []ccipChangeSetSolana.RemoteChainTokenPoolConfig{
					{
						SolChainSelector: solChainSel,
						SolTokenPubKey:   solTokenAddress,
						SolPoolType:      &bnm,
						Metadata:         shared.CLLMetadata,
						EVMRemoteConfigs: map[uint64]ccipChangeSetSolana.EVMRemoteConfig{
							evmChainSel: {
								TokenSymbol: shared.TokenSymbol(evmTokenName),
								PoolType:    shared.BurnMintTokenPool,
								PoolVersion: shared.CurrentTokenPoolVersion,
								RateLimiterConfig: ccipChangeSetSolana.RateLimiterConfig{
									Inbound: solTestTokenPool.RateLimitConfig{
										Enabled:  false,
										Capacity: 0,
										Rate:     0,
									},
									Outbound: solTestTokenPool.RateLimitConfig{
										Enabled:  false,
										Capacity: 0,
										Rate:     0,
									},
								},
							},
						},
					},
				},
			},
		),
	)
	if err != nil {
		return nil, nil, solana.PublicKey{}, err
	}

	// configure evm
	poolConfigPDA, err := soltokens.TokenPoolConfigAddress(solTokenAddress, state.SolChains[solChainSel].BurnMintTokenPools[shared.CLLMetadata])
	if err != nil {
		return nil, nil, solana.PublicKey{}, err
	}
	err = setTokenPoolCounterPart(e.BlockChains.EVMChains()[evmChainSel], evmPool, evmDeployer, solChainSel, solTokenAddress.Bytes(), poolConfigPDA.Bytes())
	if err != nil {
		return nil, nil, solana.PublicKey{}, err
	}

	err = grantMintBurnPermissions(lggr, e.BlockChains.EVMChains()[evmChainSel], evmToken, evmDeployer, evmPool.Address())
	if err != nil {
		return nil, nil, solana.PublicKey{}, err
	}

	return evmToken, evmPool, solTokenAddress, nil
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

// TODO: this should be linked to the solChain function
func SavePreloadedSolAddresses(e cldf.Environment, solChainSelector uint64) error {
	tv := cldf.NewTypeAndVersion(shared.Router, deployment.Version1_0_0)
	err := e.ExistingAddresses.Save(solChainSelector, memory.SolanaProgramIDs["ccip_router"], tv)
	if err != nil {
		return err
	}
	tv = cldf.NewTypeAndVersion(shared.Receiver, deployment.Version1_0_0)
	err = e.ExistingAddresses.Save(solChainSelector, memory.SolanaProgramIDs["test_ccip_receiver"], tv)
	if err != nil {
		return err
	}
	tv = cldf.NewTypeAndVersion(shared.FeeQuoter, deployment.Version1_0_0)
	err = e.ExistingAddresses.Save(solChainSelector, memory.SolanaProgramIDs["fee_quoter"], tv)
	if err != nil {
		return err
	}
	tv = cldf.NewTypeAndVersion(shared.OffRamp, deployment.Version1_0_0)
	err = e.ExistingAddresses.Save(solChainSelector, memory.SolanaProgramIDs["ccip_offramp"], tv)
	if err != nil {
		return err
	}
	tv = cldf.NewTypeAndVersion(shared.BurnMintTokenPool, deployment.Version1_0_0)
	err = e.ExistingAddresses.Save(solChainSelector, memory.SolanaProgramIDs["burnmint_token_pool"], tv)
	if err != nil {
		return err
	}
	tv = cldf.NewTypeAndVersion(shared.LockReleaseTokenPool, deployment.Version1_0_0)
	err = e.ExistingAddresses.Save(solChainSelector, memory.SolanaProgramIDs["lockrelease_token_pool"], tv)
	if err != nil {
		return err
	}
	tv = cldf.NewTypeAndVersion(commontypes.ManyChainMultisigProgram, deployment.Version1_0_0)
	err = e.ExistingAddresses.Save(solChainSelector, memory.SolanaProgramIDs["mcm"], tv)
	if err != nil {
		return err
	}
	tv = cldf.NewTypeAndVersion(commontypes.AccessControllerProgram, deployment.Version1_0_0)
	err = e.ExistingAddresses.Save(solChainSelector, memory.SolanaProgramIDs["access_controller"], tv)
	if err != nil {
		return err
	}
	tv = cldf.NewTypeAndVersion(commontypes.RBACTimelockProgram, deployment.Version1_0_0)
	err = e.ExistingAddresses.Save(solChainSelector, memory.SolanaProgramIDs["timelock"], tv)
	if err != nil {
		return err
	}
	tv = cldf.NewTypeAndVersion(shared.RMNRemote, deployment.Version1_0_0)
	err = e.ExistingAddresses.Save(solChainSelector, memory.SolanaProgramIDs["rmn_remote"], tv)
	if err != nil {
		return err
	}
	return nil
}

func ValidateSolanaState(t *testing.T, e cldf.Environment, solChainSelectors []uint64) {
	state, err := stateview.LoadOnchainStateSolana(e)
	require.NoError(t, err, "Failed to load Solana state")

	for _, sel := range solChainSelectors {
		// Validate chain exists in state
		chainState, exists := state.SolChains[sel]
		require.True(t, exists, "Chain selector %d not found in Solana state", sel)

		// Validate addresses
		require.False(t, chainState.Router.IsZero(), "Router address is zero for chain %d", sel)
		require.False(t, chainState.OffRamp.IsZero(), "OffRamp address is zero for chain %d", sel)
		require.False(t, chainState.FeeQuoter.IsZero(), "FeeQuoter address is zero for chain %d", sel)
		require.False(t, chainState.LinkToken.IsZero(), "Link token address is zero for chain %d", sel)
		require.False(t, chainState.RMNRemote.IsZero(), "RMNRemote address is zero for chain %d", sel)

		// Get router config
		var routerConfigAccount solRouter.Config
		err = e.BlockChains.SolanaChains()[sel].GetAccountDataBorshInto(testcontext.Get(t), chainState.RouterConfigPDA, &routerConfigAccount)
		require.NoError(t, err, "Failed to deserialize router config for chain %d", sel)

		// Get fee quoter config
		var feeQuoterConfigAccount solFeeQuoter.Config
		err = e.BlockChains.SolanaChains()[sel].GetAccountDataBorshInto(testcontext.Get(t), chainState.FeeQuoterConfigPDA, &feeQuoterConfigAccount)
		require.NoError(t, err, "Failed to deserialize fee quoter config for chain %d", sel)

		// Get offramp config
		var offRampConfigAccount solOffRamp.Config
		err = e.BlockChains.SolanaChains()[sel].GetAccountDataBorshInto(testcontext.Get(t), chainState.OffRampConfigPDA, &offRampConfigAccount)
		require.NoError(t, err, "Failed to deserialize offramp config for chain %d", sel)

		// Get rmn remote config
		var rmnRemoteConfigAccount solRmnRemote.Config
		err = e.BlockChains.SolanaChains()[sel].GetAccountDataBorshInto(testcontext.Get(t), chainState.RMNRemoteConfigPDA, &rmnRemoteConfigAccount)
		require.NoError(t, err, "Failed to deserialize rmn remote config for chain %d", sel)

		addressLookupTable, err := solanastateview.FetchOfframpLookupTable(e.GetContext(), e.BlockChains.SolanaChains()[sel], chainState.OffRamp)
		require.NoError(t, err, "Failed to get offramp lookup table for chain %d", sel)

		addresses, err := solcommon.GetAddressLookupTable(
			e.GetContext(),
			e.BlockChains.SolanaChains()[sel].Client,
			addressLookupTable)
		require.NoError(t, err, "Failed to get address lookup table for chain %d", sel)
		require.GreaterOrEqual(t, len(addresses), 22, "Not enough addresses found in lookup table for chain %d", sel)
	}
}

func DeploySolanaCcipReceiver(t *testing.T, e cldf.Environment) {
	state, err := stateview.LoadOnchainStateSolana(e)
	require.NoError(t, err)
	for solSelector, chainState := range state.SolChains {
		solTestReceiver.SetProgramID(chainState.Receiver)
		externalExecutionConfigPDA, _, _ := solana.FindProgramAddress([][]byte{[]byte("external_execution_config")}, chainState.Receiver)
		instruction, ixErr := solTestReceiver.NewInitializeInstruction(
			chainState.Router,
			solanastateview.FindReceiverTargetAccount(chainState.Receiver),
			externalExecutionConfigPDA,
			e.BlockChains.SolanaChains()[solSelector].DeployerKey.PublicKey(),
			solana.SystemProgramID,
		).ValidateAndBuild()
		require.NoError(t, ixErr)
		err = e.BlockChains.SolanaChains()[solSelector].Confirm([]solana.Instruction{instruction})
		require.NoError(t, err)
	}
}

func TransferOwnershipSolana(
	t *testing.T,
	e *cldf.Environment,
	solChain uint64,
	needTimelockDeployed bool,
	contractsToTransfer ccipChangeSetSolana.CCIPContractsToTransfer,
) (timelockSignerPDA solana.PublicKey, mcmSignerPDA solana.PublicKey) {
	var err error
	if needTimelockDeployed {
		*e, _, err = commoncs.ApplyChangesetsV2(t, *e, []commoncs.ConfiguredChangeSet{
			commoncs.Configure(
				cldf.CreateLegacyChangeSet(commoncs.DeployMCMSWithTimelockV2),
				map[uint64]commontypes.MCMSWithTimelockConfigV2{
					solChain: {
						Canceller:        proposalutils.SingleGroupMCMSV2(t),
						Proposer:         proposalutils.SingleGroupMCMSV2(t),
						Bypasser:         proposalutils.SingleGroupMCMSV2(t),
						TimelockMinDelay: big.NewInt(0),
					},
				},
			),
		})
		require.NoError(t, err)
	}

	addresses, err := e.ExistingAddresses.AddressesForChain(solChain)
	require.NoError(t, err)
	mcmState, err := state.MaybeLoadMCMSWithTimelockChainStateSolana(e.BlockChains.SolanaChains()[solChain], addresses)
	require.NoError(t, err)

	// Fund signer PDAs for timelock and mcm
	// If we don't fund, execute() calls will fail with "no funds" errors.
	timelockSignerPDA = state.GetTimelockSignerPDA(mcmState.TimelockProgram, mcmState.TimelockSeed)
	mcmSignerPDA = state.GetMCMSignerPDA(mcmState.McmProgram, mcmState.ProposerMcmSeed)
	err = memory.FundSolanaAccounts(e.GetContext(), []solana.PublicKey{timelockSignerPDA, mcmSignerPDA},
		100, e.BlockChains.SolanaChains()[solChain].Client)
	require.NoError(t, err)
	t.Logf("funded timelock signer PDA: %s", timelockSignerPDA.String())
	t.Logf("funded mcm signer PDA: %s", mcmSignerPDA.String())
	// Apply transfer ownership changeset
	*e, _, err = commoncs.ApplyChangesetsV2(t, *e, []commoncs.ConfiguredChangeSet{
		commoncs.Configure(
			cldf.CreateLegacyChangeSet(ccipChangeSetSolana.TransferCCIPToMCMSWithTimelockSolana),
			ccipChangeSetSolana.TransferCCIPToMCMSWithTimelockSolanaConfig{
				MCMSCfg: proposalutils.TimelockConfig{MinDelay: 1 * time.Second},
				ContractsByChain: map[uint64]ccipChangeSetSolana.CCIPContractsToTransfer{
					solChain: contractsToTransfer,
				},
			},
		),
	})
	require.NoError(t, err)
	return timelockSignerPDA, mcmSignerPDA
}

func GenTestTransferOwnershipConfig(
	e DeployedEnv,
	chains []uint64,
	state stateview.CCIPOnChainState,
	withTestRouterTransfer bool,
) commoncs.TransferToMCMSWithTimelockConfig {
	var (
		timelocksPerChain = make(map[uint64]common.Address)
		contracts         = make(map[uint64][]common.Address)
	)

	// chain contracts
	for _, chain := range chains {
		timelocksPerChain[chain] = state.MustGetEVMChainState(chain).Timelock.Address()
		contracts[chain] = []common.Address{
			state.MustGetEVMChainState(chain).OnRamp.Address(),
			state.MustGetEVMChainState(chain).OffRamp.Address(),
			state.MustGetEVMChainState(chain).FeeQuoter.Address(),
			state.MustGetEVMChainState(chain).NonceManager.Address(),
			state.MustGetEVMChainState(chain).RMNRemote.Address(),
			state.MustGetEVMChainState(chain).Router.Address(),
			state.MustGetEVMChainState(chain).TokenAdminRegistry.Address(),
			state.MustGetEVMChainState(chain).RMNProxy.Address(),
		}
		if withTestRouterTransfer {
			contracts[chain] = append(contracts[chain], state.MustGetEVMChainState(chain).TestRouter.Address())
		}
	}

	// home chain
	homeChainTimelockAddress := state.MustGetEVMChainState(e.HomeChainSel).Timelock.Address()
	timelocksPerChain[e.HomeChainSel] = homeChainTimelockAddress
	contracts[e.HomeChainSel] = append(contracts[e.HomeChainSel],
		state.MustGetEVMChainState(e.HomeChainSel).CapabilityRegistry.Address(),
		state.MustGetEVMChainState(e.HomeChainSel).CCIPHome.Address(),
		state.MustGetEVMChainState(e.HomeChainSel).RMNHome.Address(),
	)

	return commoncs.TransferToMCMSWithTimelockConfig{
		ContractsByChain: contracts,
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
