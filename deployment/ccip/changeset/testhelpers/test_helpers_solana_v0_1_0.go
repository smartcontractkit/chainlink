package testhelpers

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"math/big"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/aptos-labs/aptos-go-sdk"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	solbinary "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/token"
	"github.com/gagliardetto/solana-go/rpc"

	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"

	mcmstypes "github.com/smartcontractkit/mcms/types"

	tonOps "github.com/smartcontractkit/chainlink-ton/deployment/ccip"
	tonCfg "github.com/smartcontractkit/chainlink-ton/deployment/ccip/config"

	aptos_fee_quoter "github.com/smartcontractkit/chainlink-aptos/bindings/ccip/fee_quoter"
	"github.com/smartcontractkit/chainlink-aptos/bindings/helpers"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/latest/offramp"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_2_0/router"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_1/burn_mint_token_pool"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_1/usdc_token_pool"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_3/fee_quoter"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_3/message_hasher"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry"

	cldf_aptos "github.com/smartcontractkit/chainlink-deployments-framework/chain/aptos"
	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf_solana "github.com/smartcontractkit/chainlink-deployments-framework/chain/solana"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cldf_offchain "github.com/smartcontractkit/chainlink-deployments-framework/offchain"

	aptoscs "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/config"

	ccipChangeSetSolanaV0_1_0 "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/solana_v0_1_0"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_6"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview/evm"
	solanastateview "github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview/solana"

	commoncs "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ccipevm"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay"

	"github.com/smartcontractkit/chainlink-ccip/pkg/consts"
	"github.com/smartcontractkit/chainlink-ccip/pkg/reader"
	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/onramp"
	solconfig "github.com/smartcontractkit/chainlink-ccip/chains/solana/contracts/tests/config"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/v0_1_0/base_token_pool"
	solCommon "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/v0_1_0/ccip_common"
	solOffRamp "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/v0_1_0/ccip_offramp"
	solRouter "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/v0_1_0/ccip_router"
	solFeeQuoter "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/v0_1_0/fee_quoter"
	solRmnRemote "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/v0_1_0/rmn_remote"
	solTestReceiver "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/v0_1_0/test_ccip_receiver"
	solccip "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/ccip"
	solcommon "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/common"
	solstate "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/state"
	soltokens "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/tokens"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/mock_ethusd_aggregator_wrapper"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/initial/aggregator_v3_interface"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/initial/burn_mint_erc677"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/initial/mock_v3_aggregator_contract"
	"github.com/smartcontractkit/chainlink/deployment"

	ccipclient "github.com/smartcontractkit/chainlink/deployment/ccip/shared/client"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/environment/devenv"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/abihelpers"
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
func ReplayLogs(t *testing.T, oc cldf_offchain.Client, replayBlocks map[uint64]uint64, opts ...ReplayLogsOption) {
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

func WaitForEventFilterRegistration(t *testing.T, oc cldf_offchain.Client, chainSel uint64, eventName string, address []byte) error {
	family, err := chainsel.GetSelectorFamily(chainSel)
	if err != nil {
		return err
	}

	var eventID string
	switch family {
	case chainsel.FamilyEVM:
		evmOnRampABI, err := onramp.OnRampMetaData.GetAbi()
		require.NoError(t, err)
		if event, ok := evmOnRampABI.Events[eventName]; ok {
			eventID = event.ID.String()
			break
		}
		evmOffRampABI, err := offramp.OffRampMetaData.GetAbi()
		require.NoError(t, err)
		if event, ok := evmOffRampABI.Events[eventName]; ok {
			eventID = event.ID.String()
			break
		}
		return fmt.Errorf("failed to find event with name %s in onramp or offramp ABIs", eventName)
	case chainsel.FamilySolana:
		eventID = eventName
	case chainsel.FamilyAptos:
		// Aptos is not using LogPoller
		return nil
	case chainsel.FamilySui:
		// Sui is not using LogPoller
	case chainsel.FamilyTon:
		// TODO: TON is not using LogPoller
		return nil
	default:
		return fmt.Errorf("unsupported chain family; %v", family)
	}

	require.Eventually(t, func() bool {
		registered, err := isLogFilterRegistered(t, oc, chainSel, eventID, address)
		require.NoError(t, err)
		return registered
	}, 10*time.Minute, 5*time.Second)

	return nil
}

func isLogFilterRegistered(t *testing.T, oc cldf_offchain.Client, chainSel uint64, eventName string, address []byte) (bool, error) {
	var registered bool
	var err error
	switch oc := oc.(type) {
	case *memory.JobClient:
		registered, err = oc.IsLogFilterRegistered(t.Context(), chainSel, eventName, address)
	default:
		return false, fmt.Errorf("unsupported offchain client type %T", oc)
	}
	return registered, err
}

func WaitForEventFilterRegistrationOnLane(t *testing.T, onchainState stateview.CCIPOnChainState, onchainClient cldf_offchain.Client, sourceChainSel, destChainSel uint64) {
	onRampAddr, err := onchainState.GetOnRampAddressBytes(sourceChainSel)
	require.NoError(t, err)
	// Ensure CCIPMessageSent event filter is registered
	// Sending message too early could result in LogPoller missing the send event
	err = WaitForEventFilterRegistration(t, onchainClient, sourceChainSel, consts.EventNameCCIPMessageSent, onRampAddr)
	require.NoError(t, err)
	// Ensure CommitReportAccepted and ExecutionStateChanged event filters are registered for the offramp
	// The LogPoller could pick up the message sent event but miss the commit or execute event
	offRampAddr, err := onchainState.GetOffRampAddressBytes(destChainSel)
	require.NoError(t, err)
	err = WaitForEventFilterRegistration(t, onchainClient, destChainSel, consts.EventNameCommitReportAccepted, offRampAddr)
	require.NoError(t, err)
	err = WaitForEventFilterRegistration(t, onchainClient, destChainSel, consts.EventNameExecutionStateChanged, offRampAddr)
	require.NoError(t, err)

	t.Logf("%s, %s, and %s filters registered", consts.EventNameCCIPMessageSent, consts.EventNameCommitReportAccepted, consts.EventNameExecutionStateChanged)
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
			return 0, fmt.Errorf("failed to get latest header for chain %d: %w", chainSelector, err)
		}
		block := latesthdr.Number.Uint64()
		return block, nil
	case chainsel.FamilySolana:
		return env.BlockChains.SolanaChains()[chainSelector].Client.GetSlot(ctx, solconfig.DefaultCommitment)
	case chainsel.FamilySui:
		suiClient := env.BlockChains.SuiChains()[chainSelector].Client
		seqNum, err := suiClient.SuiGetLatestCheckpointSequenceNumber(ctx)
		if err != nil {
			return 0, fmt.Errorf("failed to get sui latest checkpoint: %w", err)
		}

		fmt.Println("LATEST BLOCK ON SUI: ", seqNum)
		return seqNum, nil
	case chainsel.FamilyAptos:
		chainInfo, err := env.BlockChains.AptosChains()[chainSelector].Client.Info()
		if err != nil {
			return 0, fmt.Errorf("failed to get chain info for chain %d: %w", chainSelector, err)
		}
		return chainInfo.LedgerVersion(), nil
	default:
		return 0, errors.New("unsupported chain family")
	}
}

func LatestBlocksByChain(ctx context.Context, env cldf.Environment) (map[uint64]uint64, error) {
	latestBlocks := make(map[uint64]uint64)

	chains := []uint64{}
	chains = slices.AppendSeq(chains, maps.Keys(env.BlockChains.EVMChains()))
	chains = slices.AppendSeq(chains, maps.Keys(env.BlockChains.SolanaChains()))
	suiChains := env.BlockChains.SuiChains()
	chains = slices.AppendSeq(chains, maps.Keys(suiChains))

	chains = slices.AppendSeq(chains, maps.Keys(env.BlockChains.AptosChains()))
	for _, selector := range chains {
		block, err := LatestBlock(ctx, env, selector)
		if err != nil {
			return nil, fmt.Errorf("failed to get latest block for chain %d: %w", selector, err)
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
	slices.Sort(chainSels)
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
	cfg *ccipclient.CCIPSendReqConfig,
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
		return tx, 0, fmt.Errorf("failed to confirm CCIP message: %w", err)
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
	cfg *ccipclient.CCIPSendReqConfig,
) (*types.Transaction, uint64, error) {
	const errCodeInsufficientFee = "0x07da6ee6"
	const cannotDecodeErrorReason = "could not decode error reason"
	const errMsgMissingTrieNode = "missing trie node"

	defer func() { cfg.Sender.Value = nil }()

	msg := cfg.Message.(router.ClientEVM2AnyMessage)
	var retryCount int
	for {
		fmt.Println("ABOUT TO SEND THIS MSG: ", msg, cfg.DestChain)
		fee, err := r.GetFee(&bind.CallOpts{Context: context.Background()}, cfg.DestChain, msg)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to get EVM fee: %w", cldf.MaybeDataErr(err))
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
	opts ...ccipclient.SendReqOpts,
) (msgSentEvent *ccipclient.AnyMsgSentEvent) {
	baseOpts := []ccipclient.SendReqOpts{
		ccipclient.WithSourceChain(src),
		ccipclient.WithDestChain(dest),
		ccipclient.WithTestRouter(testRouter),
		ccipclient.WithMessage(msg),
	}
	baseOpts = append(baseOpts, opts...)

	msgSentEvent, err := SendRequest(e, state, baseOpts...)
	require.NoError(t, err)
	return msgSentEvent
}

// SendRequest similar to TestSendRequest but returns an error.
func SendRequest(
	e cldf.Environment,
	state stateview.CCIPOnChainState,
	opts ...ccipclient.SendReqOpts,
) (*ccipclient.AnyMsgSentEvent, error) {
	cfg := &ccipclient.CCIPSendReqConfig{}
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
	case chainsel.FamilySui:
		return SendRequestSui(e, state, cfg)
	case chainsel.FamilyAptos:
		return SendRequestAptos(e, state, cfg)
	case chainsel.FamilyTon:
		seq, raw, err := tonOps.SendTonRequest(e, state.TonChains[cfg.SourceChain], cfg.SourceChain, cfg.DestChain, cfg.Message.(tonOps.TonSendRequest))
		if err != nil {
			return nil, err
		}

		return &ccipclient.AnyMsgSentEvent{
			SequenceNumber: seq,
			RawEvent:       raw,
		}, nil
	default:
		return nil, fmt.Errorf("send request: unsupported chain family: %v", family)
	}
}

func SendRequestEVM(
	e cldf.Environment,
	state stateview.CCIPOnChainState,
	cfg *ccipclient.CCIPSendReqConfig,
) (*ccipclient.AnyMsgSentEvent, error) {
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
	return &ccipclient.AnyMsgSentEvent{
		SequenceNumber: it.Event.SequenceNumber,
		RawEvent:       it.Event,
	}, nil
}

func SendRequestSui(
	e cldf.Environment,
	state stateview.CCIPOnChainState,
	cfg *ccipclient.CCIPSendReqConfig,
) (*ccipclient.AnyMsgSentEvent, error) {
	return SendSuiCCIPRequest(e, cfg)
}

func SendRequestSol(
	e cldf.Environment,
	state stateview.CCIPOnChainState,
	cfg *ccipclient.CCIPSendReqConfig,
) (*ccipclient.AnyMsgSentEvent, error) { // TODO: chain independent return value
	ctx := e.GetContext()

	s := state.SolChains[cfg.SourceChain]
	c := e.BlockChains.SolanaChains()[cfg.SourceChain]

	destinationChainSelector := cfg.DestChain
	message := cfg.Message.(solRouter.SVM2AnyMessage)
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

	base := solRouter.NewCcipSendInstruction(
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
		allTokenPools = append(allTokenPools, s.CCTPTokenPool)

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

	tempIx, err := base.ValidateAndBuild()
	if err != nil {
		return nil, err
	}
	ixData, err := tempIx.Data()
	if err != nil {
		return nil, fmt.Errorf("failed to extract data payload from router ccip send instruction: %w", err)
	}
	ix := solana.NewInstruction(s.Router, tempIx.Accounts(), ixData)

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

	return &ccipclient.AnyMsgSentEvent{
		SequenceNumber: ccipMessageSentEvent.SequenceNumber,
		RawEvent: &onramp.OnRampCCIPMessageSent{
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
				FeeTokenAmount: ConvertSolanaCrossChainAmountToBigInt(ccipMessageSentEvent.Message.FeeTokenAmount.LeBytes),
				FeeValueJuels:  ConvertSolanaCrossChainAmountToBigInt(ccipMessageSentEvent.Message.FeeValueJuels.LeBytes),
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
		},
	}, nil
}

func ConvertSolanaCrossChainAmountToBigInt(amountLeBytes [32]uint8) *big.Int {
	bytes := amountLeBytes[:]
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
	state stateview.CCIPOnChainState,
	from, to uint64,
	isTestRouter bool,
	gasPrices map[uint64]*big.Int,
	tokenPrices map[string]*big.Int,
	fqCfg fee_quoter.FeeQuoterDestChainConfig,
) error {
	var err error
	fromFamily, err := chainsel.GetSelectorFamily(from)
	require.NoError(t, err)
	toFamily, err := chainsel.GetSelectorFamily(to)
	require.NoError(t, err)
	changesets := []commoncs.ConfiguredChangeSet{}

	switch fromFamily {
	case chainsel.FamilyEVM:
		evmTokenPrices := make(map[common.Address]*big.Int, len(tokenPrices))
		for address, price := range tokenPrices {
			evmTokenPrices[common.HexToAddress(address)] = price
		}
		changesets = append(changesets, AddEVMSrcChangesets(from, to, isTestRouter, gasPrices, evmTokenPrices, fqCfg)...)
	case chainsel.FamilySolana:
		changesets = append(changesets, AddLaneSolanaChangesetsV0_1_0(e, from, to, toFamily)...)
	case chainsel.FamilyAptos:
		aptosTokenPrices := make(map[aptos.AccountAddress]*big.Int, len(tokenPrices))
		for address, price := range tokenPrices {
			aptosTokenPrices[aptoscs.MustParseAddress(t, address)] = price
		}
		changesets = append(changesets, AddLaneAptosChangesets(t, from, to, gasPrices, aptosTokenPrices)...)
	case chainsel.FamilyTon:
		onRamp, err := state.GetOnRampAddressBytes(to)
		if err != nil {
			return err
		}
		addLaneConfig := tonOps.AddLaneTONConfig(&e.Env, onRamp, from, to, fromFamily, toFamily, gasPrices)
		changesets = append(changesets, commoncs.Configure(tonOps.AddTonLanes{},
			tonCfg.UpdateTonLanesConfig{
				Lanes:      []tonCfg.LaneConfig{addLaneConfig},
				TestRouter: false,
			}))
	}

	// changesets = append(changesets, AddEVMDestChangesets(e, 909606746561742123, 18395503381733958356, false)...)

	switch toFamily {
	case chainsel.FamilyEVM:
		changesets = append(changesets, AddEVMDestChangesets(e, to, from, isTestRouter)...)
	case chainsel.FamilySolana:
		changesets = append(changesets, AddLaneSolanaChangesetsV0_1_0(e, to, from, fromFamily)...)
	case chainsel.FamilyAptos:
		changesets = append(changesets, AddLaneAptosChangesets(t, from, to, gasPrices, nil)...)
	case chainsel.FamilyTon:
		onRamp, err := state.GetOnRampAddressBytes(from)
		if err != nil {
			return err
		}
		addLaneConfig := tonOps.AddLaneTONConfig(&e.Env, onRamp, from, to, fromFamily, toFamily, gasPrices)
		changesets = append(changesets, commoncs.Configure(tonOps.AddTonLanes{},
			tonCfg.UpdateTonLanesConfig{
				Lanes:      []tonCfg.LaneConfig{addLaneConfig},
				TestRouter: false,
			}))
	}

	fmt.Println("ADDLANE CHANGESETS: ", changesets)

	e.Env, _, err = commoncs.ApplyChangesets(t, e.Env, changesets)
	if err != nil {
		fmt.Println("ERROR APPLYING CHANGESET", err)
		return err
	}
	return nil
}

func AddLaneSolanaChangesetsV0_1_0(e *DeployedEnv, solChainSelector, remoteChainSelector uint64, remoteFamily string) []commoncs.ConfiguredChangeSet {
	var chainFamilySelector [4]uint8
	switch remoteFamily {
	case chainsel.FamilyEVM:
		// bytes4(keccak256("CCIP ChainFamilySelector EVM"))
		chainFamilySelector = [4]uint8{40, 18, 213, 44}
	case chainsel.FamilySolana:
		// bytes4(keccak256("CCIP ChainFamilySelector SVM"));
		chainFamilySelector = [4]uint8{30, 16, 189, 196}
	case chainsel.FamilyAptos:
		// bytes4(keccak256("CCIP ChainFamilySelector APTOS"));
		chainFamilySelector = [4]uint8{0xac, 0x77, 0xff, 0xec}
	default:
		panic("unsupported remote family")
	}
	solanaChangesets := []commoncs.ConfiguredChangeSet{
		commoncs.Configure(
			cldf.CreateLegacyChangeSet(ccipChangeSetSolanaV0_1_0.AddRemoteChainToRouter),
			ccipChangeSetSolanaV0_1_0.AddRemoteChainToRouterConfig{
				ChainSelector: solChainSelector,
				UpdatesByChain: map[uint64]*ccipChangeSetSolanaV0_1_0.RouterConfig{
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
			cldf.CreateLegacyChangeSet(ccipChangeSetSolanaV0_1_0.AddRemoteChainToFeeQuoter),
			ccipChangeSetSolanaV0_1_0.AddRemoteChainToFeeQuoterConfig{
				ChainSelector: solChainSelector,
				UpdatesByChain: map[uint64]*ccipChangeSetSolanaV0_1_0.FeeQuoterConfig{
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
			cldf.CreateLegacyChangeSet(ccipChangeSetSolanaV0_1_0.AddRemoteChainToOffRamp),
			ccipChangeSetSolanaV0_1_0.AddRemoteChainToOffRampConfig{
				ChainSelector: solChainSelector,
				UpdatesByChain: map[uint64]*ccipChangeSetSolanaV0_1_0.OffRampConfig{
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

func AddSuiDestChangeset(e *DeployedEnv, to, from uint64, isTestRouter bool) []commoncs.ConfiguredChangeSet {
	suiDstChangesets := []commoncs.ConfiguredChangeSet{
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
	}

	return suiDstChangesets
}

func AddLaneAptosChangesets(t *testing.T, srcChainSelector, destChainSelector uint64, gasPrices map[uint64]*big.Int, tokenPrices map[aptos.AccountAddress]*big.Int) []commoncs.ConfiguredChangeSet {
	srcFamily, err := chainsel.GetSelectorFamily(srcChainSelector)
	require.NoError(t, err)
	destFamily, err := chainsel.GetSelectorFamily(destChainSelector)
	require.NoError(t, err)

	if srcFamily != chainsel.FamilyAptos &&
		destFamily != chainsel.FamilyAptos {
		t.Fatalf("At least one of the provided source/destination chains has to be Aptos. srcFamily: %v destFamily: %v", srcFamily, destFamily)
	}

	var src, dest config.ChainDefinition

	switch srcFamily {
	case chainsel.FamilyEVM:
		src = config.EVMChainDefinition{
			ChainDefinition: v1_6.ChainDefinition{
				ConnectionConfig: v1_6.ConnectionConfig{
					RMNVerificationDisabled: true,
				},
				Selector: srcChainSelector,
			},
		}
	case chainsel.FamilyAptos:
		src = config.AptosChainDefinition{
			TokenPrices: tokenPrices,
			ConnectionConfig: v1_6.ConnectionConfig{
				RMNVerificationDisabled: true,
			},
			Selector:                      srcChainSelector,
			AddTokenTransferFeeConfigs:    nil,
			RemoveTokenTransferFeeConfigs: nil,
		}
	default:
		t.Fatalf("Unsupported source chain family: %v", srcFamily)
	}

	switch destFamily {
	case chainsel.FamilyEVM:
		dest = config.EVMChainDefinition{
			ChainDefinition: v1_6.ChainDefinition{
				ConnectionConfig: v1_6.ConnectionConfig{
					AllowListEnabled: false,
				},
				Selector: destChainSelector,
				GasPrice: gasPrices[destChainSelector],
				FeeQuoterDestChainConfig: fee_quoter.FeeQuoterDestChainConfig{
					IsEnabled:                         true,
					MaxNumberOfTokensPerMsg:           10,
					MaxDataBytes:                      30_000,
					MaxPerMsgGasLimit:                 3_000_000,
					DestGasOverhead:                   ccipevm.DestGasOverhead,
					DestGasPerPayloadByteBase:         ccipevm.CalldataGasPerByteBase,
					DestGasPerPayloadByteHigh:         ccipevm.CalldataGasPerByteHigh,
					DestGasPerPayloadByteThreshold:    ccipevm.CalldataGasPerByteThreshold,
					DestDataAvailabilityOverheadGas:   100,
					DestGasPerDataAvailabilityByte:    16,
					DestDataAvailabilityMultiplierBps: 1,
					ChainFamilySelector:               [4]byte{0x28, 0x12, 0xd5, 0x2c},
					EnforceOutOfOrder:                 false,
					DefaultTokenFeeUSDCents:           25,
					DefaultTokenDestGasOverhead:       90_000,
					DefaultTxGasLimit:                 200_000,
					GasMultiplierWeiPerEth:            11e8, // TODO what's the scale here ?
					GasPriceStalenessThreshold:        0,
					NetworkFeeUSDCents:                10,
				},
			},
			OnRampVersion: []byte{1, 6, 0},
		}
	case chainsel.FamilyAptos:
		dest = config.AptosChainDefinition{
			ConnectionConfig: v1_6.ConnectionConfig{
				AllowListEnabled: false,
			},
			Selector: destChainSelector,
			GasPrice: gasPrices[destChainSelector],
			FeeQuoterDestChainConfig: aptos_fee_quoter.DestChainConfig{
				IsEnabled:                         true,
				MaxNumberOfTokensPerMsg:           10,
				MaxDataBytes:                      30_000,
				MaxPerMsgGasLimit:                 3_000_000,
				DestGasOverhead:                   ccipevm.DestGasOverhead,
				DestGasPerPayloadByteBase:         ccipevm.CalldataGasPerByteBase,
				DestGasPerPayloadByteHigh:         ccipevm.CalldataGasPerByteHigh,
				DestGasPerPayloadByteThreshold:    ccipevm.CalldataGasPerByteThreshold,
				DestDataAvailabilityOverheadGas:   100,
				DestGasPerDataAvailabilityByte:    16,
				DestDataAvailabilityMultiplierBps: 1,
				ChainFamilySelector:               []byte{0xac, 0x77, 0xff, 0xec},
				EnforceOutOfOrder:                 true,
				DefaultTokenFeeUsdCents:           25,
				DefaultTokenDestGasOverhead:       90_000,
				DefaultTxGasLimit:                 200_000,
				GasMultiplierWeiPerEth:            11e17,
				GasPriceStalenessThreshold:        0,
				NetworkFeeUsdCents:                10,
			},
		}
	default:
		t.Fatalf("Unsupported dstination chain family: %v", srcFamily)
	}

	return []commoncs.ConfiguredChangeSet{
		commoncs.Configure(
			aptoscs.AddAptosLanes{},
			config.UpdateAptosLanesConfig{
				AptosMCMSConfig: &proposalutils.TimelockConfig{
					MinDelay:     time.Second,
					MCMSAction:   mcmstypes.TimelockActionSchedule,
					OverrideRoot: false,
				},
				Lanes: []config.LaneConfig{
					{
						Source:     src,
						Dest:       dest,
						IsDisabled: false,
					},
				},
				TestRouter: false,
			},
		),
	}
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
	e.Env, _, err = commoncs.ApplyChangesets(t, e.Env, apps)
	require.NoError(t, err)
}

func AddLaneWithDefaultPricesAndFeeQuoterConfig(t *testing.T, e *DeployedEnv, state stateview.CCIPOnChainState, from, to uint64, isTestRouter bool) error {
	gasPrices := map[uint64]*big.Int{
		to: DefaultGasPrice,
	}
	fromFamily, err := chainsel.GetSelectorFamily(from)
	require.NoError(t, err)

	// Maps token address => price
	// Uses string to be re-usable across chains
	tokenPrices := make(map[string]*big.Int)
	switch fromFamily {
	case chainsel.FamilyEVM:
		stateChainFrom := state.MustGetEVMChainState(from)
		tokenPrices[stateChainFrom.LinkToken.Address().String()] = DefaultLinkPrice
		tokenPrices[stateChainFrom.Weth9.Address().String()] = DefaultWethPrice
	case chainsel.FamilyAptos:
		aptosState := state.AptosChains[from]
		tokenPrices[aptosState.LinkTokenAddress.StringLong()] = deployment.EDecMult(20, 28)
		tokenPrices[shared.AptosAPTAddress] = deployment.EDecMult(5, 28)
	case chainsel.FamilyTon:
		// TODO Need to double check this, LINK will have 9 decimals on TON like on Solana (not 18)
		tonState := state.TonChains[from]
		gasPrices[from] = big.NewInt(1e17)
		gasPrices[to] = big.NewInt(1e17)
		tokenPrices[tonState.LinkTokenAddress.String()] = deployment.EDecMult(20, 28)
	case chainsel.FamilySui:
		suiState := state.SuiChains[from]
		gasPrices[from] = big.NewInt(1e17)
		gasPrices[to] = big.NewInt(1e17)
		tokenPrices[suiState.LinkTokenCoinMetadataId] = deployment.EDecMult(20, 28)
	}
	fqCfg := v1_6.DefaultFeeQuoterDestChainConfig(true, to)

	err = AddLane(
		t,
		e,
		state,
		from, to,
		isTestRouter,
		gasPrices,
		tokenPrices,
		fqCfg,
	)
	if err != nil {
		return err
	}
	return nil
}

func AddLaneWithEnforceOutOfOrder(t *testing.T, e *DeployedEnv, state stateview.CCIPOnChainState, from, to uint64, isTestRouter bool) {
	gasPrices := map[uint64]*big.Int{
		to: DefaultGasPrice,
	}
	fromFamily, err := chainsel.GetSelectorFamily(from)
	require.NoError(t, err)

	// Maps token address => price
	// Uses string to be re-usable across chains
	tokenPrices := make(map[string]*big.Int)
	switch fromFamily {
	case chainsel.FamilyEVM:
		stateChainFrom := state.MustGetEVMChainState(from)
		tokenPrices[stateChainFrom.LinkToken.Address().String()] = DefaultLinkPrice
		tokenPrices[stateChainFrom.Weth9.Address().String()] = DefaultWethPrice
	case chainsel.FamilyAptos:
		aptosState := state.AptosChains[from]
		tokenPrices[aptosState.LinkTokenAddress.StringLong()] = deployment.EDecMult(20, 28)
		tokenPrices[shared.AptosAPTAddress] = deployment.EDecMult(5, 28)
	}
	fqCfg := v1_6.DefaultFeeQuoterDestChainConfig(true, to)
	fqCfg.EnforceOutOfOrder = true
	AddLane(
		t,
		e,
		state,
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
			Address: linkFeed, Contract: aggregatorCr, Tv: linkTV, Tx: tx, Err: errors.Join(err1, err2),
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
			Address: wethFeed, Contract: aggregatorCr, Tv: linkTV, Tx: tx, Err: errors.Join(err1, err2),
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
	allowance := new(big.Int).Mul(big.NewInt(1e18), big.NewInt(100))

	for chain, mintTokenInfos := range tokenMap {

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
						new(big.Int).Mul(allowance, big.NewInt(10)),
					)
					require.NoError(t, err)
					_, err = e.BlockChains.EVMChains()[chain].Confirm(tx)
					require.NoError(t, err)

					tx, err = token.Approve(sender, state.MustGetEVMChainState(chain).Router.Address(), allowance)
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
) (*ccipclient.AnyMsgSentEvent, map[uint64]*uint64) {
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

		msg = solRouter.SVM2AnyMessage{
			Receiver:     common.LeftPadBytes(receiver, 32),
			Data:         data,
			TokenAmounts: tokens.([]solRouter.SVMTokenAmount),
			FeeToken:     feeTokenAddr,
			ExtraArgs:    extraArgs,
		}
	case chainsel.FamilyAptos:
		feeTokenAddr := aptos.AccountAddress{}
		if len(feeToken) > 0 {
			feeTokenAddr = aptoscs.MustParseAddress(t, feeToken)
		}
		msg = AptosSendRequest{
			Data:         data,
			Receiver:     common.LeftPadBytes(receiver, 32),
			ExtraArgs:    extraArgs,
			FeeToken:     feeTokenAddr,
			TokenAmounts: tokens.([]AptosTokenAmount),
		}
	case chainsel.FamilySui:
		msg = SuiSendRequest{
			Data:         data,
			Receiver:     common.LeftPadBytes(receiver, 32),
			ExtraArgs:    extraArgs,
			FeeToken:     feeToken,
			TokenAmounts: tokens.([]SuiTokenAmount),
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
	TokenReceiverATA       []byte
	ExpectedStatus         int
	// optional
	Tokens                []router.ClientEVMTokenAmount
	SolTokens             []solRouter.SVMTokenAmount
	AptosTokens           []AptosTokenAmount
	SuiTokens             []SuiTokenAmount
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
				if destFamily == chainsel.FamilySolana || destFamily == chainsel.FamilySui {
					// for EVM2Solana token transfer we need to use tokenReceiver instead logical receiver
					expectedTokenBalances.add(tt.DestChain, tt.TokenReceiverATA, tt.ExpectedTokenBalances)
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
			case chainsel.FamilyAptos:
				tokens = tt.AptosTokens
				expectedTokenBalances.add(tt.DestChain, tt.Receiver, tt.ExpectedTokenBalances)
			case chainsel.FamilySui:
				tokens = tt.SuiTokens
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

			if prev, ok := startBlocks[tt.DestChain]; !ok || *blocks[tt.DestChain] < *prev {
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
					WaitForTheTokenBalanceSol(ctx, t, token, receiver, env.BlockChains.SolanaChains()[chainSelector], expectedBalance)
				case chainsel.FamilyAptos:
					expectedBalance := balance.Uint64()
					fungibleAssetMetadata := aptos.AccountAddress{}
					copy(fungibleAssetMetadata[32-len(id.token):], id.token)
					receiver := aptos.AccountAddress{}
					copy(receiver[32-len(id.receiver):], id.receiver)
					WaitForTokenBalanceAptos(ctx, t, fungibleAssetMetadata, receiver, env.BlockChains.AptosChains()[chainSelector], expectedBalance)
				case chainsel.FamilySui:
					tokenHex := "0x" + hex.EncodeToString(id.token)
					tokenReceiverHex := "0x" + hex.EncodeToString(id.receiver)
					fmt.Println("Waiting for TokenBalance sui: ", tokenHex, tokenReceiverHex)
					WaitForTokenBalanceSui(ctx, t, tokenHex, tokenReceiverHex, env.BlockChains.SuiChains()[chainSelector], balance)
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

func WaitForTokenBalanceAptos(
	ctx context.Context,
	t *testing.T,
	fungibleAsset aptos.AccountAddress,
	account aptos.AccountAddress,
	chain cldf_aptos.Chain,
	expected uint64,
) {
	require.Eventually(t, func() bool {
		balance, err := helpers.GetFungibleAssetBalance(chain.Client, account, fungibleAsset)
		require.NoError(t, err)

		t.Log("(Aptos) Waiting for the token balance",
			"expected", expected,
			"actual", balance,
			"fungibleAsset", fungibleAsset.StringLong(),
			"receiver", account.StringLong(),
		)

		return balance == expected
	}, tests.WaitTimeout(t), 500*time.Millisecond)
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

// GetSolanaPreloadedAddressBook returns an address book with the preloaded Solana addresses for
// the given selector.
//
// This is used because Solana programs have already been predeployed, and we need to seed the
// address book with the preloaded addresses.
func GetSolanaPreloadedAddressBook(t *testing.T, selector uint64) *cldf.AddressBookMap {
	t.Helper()

	ab := cldf.NewMemoryAddressBook()

	tv := cldf.NewTypeAndVersion(shared.Router, deployment.Version1_0_0)
	err := ab.Save(selector, memory.SolanaProgramIDs["ccip_router"], tv)
	require.NoError(t, err)

	tv = cldf.NewTypeAndVersion(shared.Receiver, deployment.Version1_0_0)
	err = ab.Save(selector, memory.SolanaProgramIDs["test_ccip_receiver"], tv)
	require.NoError(t, err)

	tv = cldf.NewTypeAndVersion(shared.FeeQuoter, deployment.Version1_0_0)
	err = ab.Save(selector, memory.SolanaProgramIDs["fee_quoter"], tv)
	require.NoError(t, err)

	tv = cldf.NewTypeAndVersion(shared.OffRamp, deployment.Version1_0_0)
	err = ab.Save(selector, memory.SolanaProgramIDs["ccip_offramp"], tv)
	require.NoError(t, err)

	tv = cldf.NewTypeAndVersion(shared.BurnMintTokenPool, deployment.Version1_0_0)
	err = ab.Save(selector, memory.SolanaProgramIDs["burnmint_token_pool"], tv)
	require.NoError(t, err)

	tv = cldf.NewTypeAndVersion(shared.LockReleaseTokenPool, deployment.Version1_0_0)
	err = ab.Save(selector, memory.SolanaProgramIDs["lockrelease_token_pool"], tv)
	require.NoError(t, err)

	tv = cldf.NewTypeAndVersion(shared.CCTPTokenPool, deployment.Version1_0_0)
	err = ab.Save(selector, memory.SolanaProgramIDs["cctp_token_pool"], tv)
	require.NoError(t, err)

	tv = cldf.NewTypeAndVersion(commontypes.ManyChainMultisigProgram, deployment.Version1_0_0)
	err = ab.Save(selector, memory.SolanaProgramIDs["mcm"], tv)
	require.NoError(t, err)

	tv = cldf.NewTypeAndVersion(commontypes.AccessControllerProgram, deployment.Version1_0_0)
	err = ab.Save(selector, memory.SolanaProgramIDs["access_controller"], tv)
	require.NoError(t, err)

	tv = cldf.NewTypeAndVersion(commontypes.RBACTimelockProgram, deployment.Version1_0_0)
	err = ab.Save(selector, memory.SolanaProgramIDs["timelock"], tv)
	require.NoError(t, err)

	tv = cldf.NewTypeAndVersion(shared.RMNRemote, deployment.Version1_0_0)
	err = ab.Save(selector, memory.SolanaProgramIDs["rmn_remote"], tv)
	require.NoError(t, err)

	return ab
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
	tv = cldf.NewTypeAndVersion(shared.CCTPTokenPool, deployment.Version1_0_0)
	err = e.ExistingAddresses.Save(solChainSelector, memory.SolanaProgramIDs["cctp_token_pool"], tv)
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

func ValidateSolanaState(e cldf.Environment, solChainSelectors []uint64) error {
	state, err := stateview.LoadOnchainStateSolana(e)
	if err != nil {
		return fmt.Errorf("failed to load Solana state: %w", err)
	}

	for _, sel := range solChainSelectors {
		// Validate chain exists in state
		chainState, exists := state.SolChains[sel]
		if !exists {
			return fmt.Errorf("chain selector %d not found in Solana state", sel)
		}

		// Validate addresses
		if chainState.Router.IsZero() {
			return fmt.Errorf("router address is zero for chain %d", sel)
		}
		if chainState.OffRamp.IsZero() {
			return fmt.Errorf("offRamp address is zero for chain %d", sel)
		}
		if chainState.FeeQuoter.IsZero() {
			return fmt.Errorf("feeQuoter address is zero for chain %d", sel)
		}
		if chainState.LinkToken.IsZero() {
			return fmt.Errorf("link token address is zero for chain %d", sel)
		}
		if chainState.RMNRemote.IsZero() {
			return fmt.Errorf("RMNRemote address is zero for chain %d", sel)
		}

		// Get router config
		var routerConfigAccount solRouter.Config
		err = e.BlockChains.SolanaChains()[sel].GetAccountDataBorshInto(context.Background(), chainState.RouterConfigPDA, &routerConfigAccount)
		if err != nil {
			return fmt.Errorf("failed to deserialize router config for chain %d: %w", sel, err)
		}

		// Get fee quoter config
		var feeQuoterConfigAccount solFeeQuoter.Config
		err = e.BlockChains.SolanaChains()[sel].GetAccountDataBorshInto(e.GetContext(), chainState.FeeQuoterConfigPDA, &feeQuoterConfigAccount)
		if err != nil {
			return fmt.Errorf("failed to deserialize fee quoter config for chain %d: %w", sel, err)
		}

		// Get offramp config
		var offRampConfigAccount solOffRamp.Config
		err = e.BlockChains.SolanaChains()[sel].GetAccountDataBorshInto(
			context.Background(),
			chainState.OffRampConfigPDA,
			&offRampConfigAccount,
		)
		if err != nil {
			return fmt.Errorf("failed to deserialize off-ramp config for chain %d: %w", sel, err)
		}

		// Get rmn remote config
		var rmnRemoteConfigAccount solRmnRemote.Config
		err = e.BlockChains.SolanaChains()[sel].GetAccountDataBorshInto(e.GetContext(), chainState.RMNRemoteConfigPDA, &rmnRemoteConfigAccount)
		if err != nil {
			return fmt.Errorf("failed to deserialize rmn remote config for chain %d: %w", sel, err)
		}

		addressLookupTable, err := solanastateview.FetchOfframpLookupTable(e.GetContext(), e.BlockChains.SolanaChains()[sel], chainState.OffRamp)
		if err != nil {
			return fmt.Errorf("failed to get offramp lookup table for chain %d: %w", sel, err)
		}

		addresses, err := solcommon.GetAddressLookupTable(
			e.GetContext(),
			e.BlockChains.SolanaChains()[sel].Client,
			addressLookupTable,
		)
		if err != nil {
			return fmt.Errorf("failed to get address lookup table for chain %d: %w", sel, err)
		}
		if len(addresses) < 22 {
			return fmt.Errorf("not enough addresses found in lookup table for chain %d: got %d, expected at least 22", sel, len(addresses))
		}
	}
	return nil
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

func TransferOwnershipSolanaV0_1_0(
	t *testing.T,
	e *cldf.Environment,
	solChain uint64,
	needTimelockDeployed bool,
	contractsToTransfer ccipChangeSetSolanaV0_1_0.CCIPContractsToTransfer,
) (timelockSignerPDA solana.PublicKey, mcmSignerPDA solana.PublicKey) {
	var err error
	if needTimelockDeployed {
		*e, _, err = commoncs.ApplyChangesets(t, *e, []commoncs.ConfiguredChangeSet{
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
	*e, _, err = commoncs.ApplyChangesets(t, *e, []commoncs.ConfiguredChangeSet{
		commoncs.Configure(
			cldf.CreateLegacyChangeSet(ccipChangeSetSolanaV0_1_0.TransferCCIPToMCMSWithTimelockSolana),
			ccipChangeSetSolanaV0_1_0.TransferCCIPToMCMSWithTimelockSolanaConfig{
				MCMSCfg: proposalutils.TimelockConfig{MinDelay: 1 * time.Second},
				ContractsByChain: map[uint64]ccipChangeSetSolanaV0_1_0.CCIPContractsToTransfer{
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
		contracts = make(map[uint64][]common.Address)
	)

	// chain contracts
	for _, chain := range chains {
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
	contracts[e.HomeChainSel] = append(contracts[e.HomeChainSel],
		state.MustGetEVMChainState(e.HomeChainSel).CapabilityRegistry.Address(),
		state.MustGetEVMChainState(e.HomeChainSel).CCIPHome.Address(),
		state.MustGetEVMChainState(e.HomeChainSel).RMNHome.Address(),
	)

	return commoncs.TransferToMCMSWithTimelockConfig{
		ContractsByChain: contracts,
	}
}

func DeployCCIPContractsTest(t *testing.T, solChains int, tonChains int) {
	e, _ := NewMemoryEnvironment(t, WithSolChains(solChains), WithTonChains(tonChains))
	// Deploy all the CCIP contracts.
	state, err := stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)
	var allChains []uint64
	allChains = append(allChains, e.Env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chainsel.FamilyEVM))...)
	allChains = append(allChains, e.Env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chainsel.FamilySolana))...)
	allChains = append(allChains, e.Env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chainsel.FamilyAptos))...)
	allChains = append(allChains, e.Env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chainsel.FamilyTon))...)
	stateView, err := state.View(&e.Env, allChains)
	require.NoError(t, err)
	if solChains > 0 {
		DeploySolanaCcipReceiver(t, e.Env)
	}

	// Assert expect every deployed address to be in the address book.
	// TODO (CCIP-3047): Add the rest of CCIPv2 representation
	b, err := json.MarshalIndent(stateView.Chains, "", "	")
	require.NoError(t, err)
	fmt.Println(string(b))
	b, err = json.MarshalIndent(stateView.SolChains, "", "	")
	require.NoError(t, err)
	fmt.Println(string(b))
	b, err = json.MarshalIndent(stateView.AptosChains, "", "	")
	require.NoError(t, err)
	fmt.Println(string(b))
	b, err = json.MarshalIndent(stateView.TONChains, "", "	")
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
	// Transfer ownership to timelock so that we can promote the zero digest later down the line.
	_, err := commoncs.Apply(t, tenv.Env,
		commoncs.Configure(
			cldf.CreateLegacyChangeSet(commoncs.TransferToMCMSWithTimelockV2),
			GenTestTransferOwnershipConfig(tenv, chains, state, withTestRouterTransfer),
		),
	)
	require.NoError(t, err)
	AssertTimelockOwnership(t, tenv, chains, state, withTestRouterTransfer)
}

func UpdateFeeQuoterForToken(
	t *testing.T,
	e cldf.Environment,
	lggr logger.Logger,
	chain cldf_evm.Chain,
	dstChain uint64,
	tokenSymbol shared.TokenSymbol,
) error {
	config := fee_quoter.FeeQuoterTokenTransferFeeConfig{
		MinFeeUSDCents:    50,
		MaxFeeUSDCents:    50_000,
		DeciBps:           0,
		DestGasOverhead:   180_000,
		DestBytesOverhead: 640,
		IsEnabled:         true,
	}
	_, err := commoncs.Apply(t, e,
		commoncs.Configure(
			cldf.CreateLegacyChangeSet(v1_6.ApplyTokenTransferFeeConfigUpdatesFeeQuoterChangeset),
			v1_6.ApplyTokenTransferFeeConfigUpdatesConfig{
				UpdatesByChain: map[uint64]v1_6.ApplyTokenTransferFeeConfigUpdatesConfigPerChain{
					chain.Selector: {
						TokenTransferFeeConfigArgs: []v1_6.TokenTransferFeeConfigArg{
							{
								DestChain: dstChain,
								TokenTransferFeeConfigPerToken: map[shared.TokenSymbol]fee_quoter.FeeQuoterTokenTransferFeeConfig{
									tokenSymbol: config,
								},
							},
						},
					},
				},
			}),
	)

	if err != nil {
		lggr.Errorw("Failed to apply token transfer fee config updates", "err", err, "config", config)
		return err
	}
	return nil
}
