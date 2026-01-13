package ccip

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/liteclient"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/ton"
	"github.com/xssnick/tonutils-go/ton/wallet"

	cldfchain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldfton "github.com/smartcontractkit/chainlink-deployments-framework/chain/ton"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	tonOps "github.com/smartcontractkit/chainlink-ton/deployment/ccip"
	tonstate "github.com/smartcontractkit/chainlink-ton/deployment/state"
	"github.com/smartcontractkit/chainlink-ton/pkg/ccip/bindings/onramp"
	tonrouter "github.com/smartcontractkit/chainlink-ton/pkg/ccip/bindings/router"
	"github.com/smartcontractkit/chainlink-ton/pkg/ton/chain"
	"github.com/smartcontractkit/chainlink-ton/pkg/ton/tvm"

	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/integration-tests/testconfig/ccip"
)

const (
	// Default gas limit for TON CCIP messages
	tonDefaultGasLimit = 200_000
)

// tonQueryCounter is used to generate unique query IDs for HighloadV3 wallet
// Initialized with a timestamp-based offset to avoid collisions with previous test runs
// HighloadV3 remembers used query IDs for the TTL period (120 seconds), so we must
// ensure each test run uses a different range of IDs
var tonQueryCounter = atomic.Uint32{}

func init() {
	// Initialize counter with time-based offset (seconds since epoch % max query ID)
	// This ensures different test runs use different ID ranges
	offset := uint32(time.Now().Unix() % (1 << 23))
	tonQueryCounter.Store(offset)
}

// loadTestHighloadV3Config is a custom HighloadV3 wallet configuration for load testing
// It uses an atomic counter to ensure unique query IDs for parallel message sending
// This fixes the issue where the default config generates the same query_id for messages
// sent within the same second, causing transaction conflicts
var loadTestHighloadV3Config = wallet.ConfigHighloadV3{
	MessageTTL: 120, // 2 minutes TTL
	MessageBuilder: func(ctx context.Context, subWalletId uint32) (id uint32, createdAt int64, err error) {
		tm := time.Now().Unix() - 30
		// Use atomic counter to ensure unique query IDs across parallel sends
		counter := tonQueryCounter.Add(1)
		return counter % (1 << 23), tm, nil
	},
}

// TonSourceManager manages TON source chains for load testing
// Uses HighloadWalletV3 which supports sending up to 254 messages per external message
// A mutex is used to prevent concurrent sends which can cause transaction conflicts
// in TON's seqno/query handling when multiple messages arrive in the same block.
//
// TODO(2026-01-12@jadepark-dev): For higher throughput, implement batch message sending
// using HighloadV3's native SendManyWaitTxHash to send all CCIP messages to different
// destinations in a single external transaction. This would require:
// 1. Collecting messages for all destinations before sending
// 2. Getting fees for each destination
// 3. Building wallet.Message objects for each destination
// 4. Sending all in one call to SendManyWaitTxHash
// This matches HighloadV3's design purpose and can send up to 254 messages per tx.
type TonSourceManager struct {
	l          logger.Logger
	client     *ton.APIClient
	wallet     *wallet.Wallet
	chainSel   uint64
	chainState tonstate.CCIPChainState // Stores Router, FeeQuoter, etc. addresses
	sendMu     sync.Mutex              // Prevents concurrent sends from the same wallet
}

// NewTonSourceManager creates a new TON source manager with HighloadWalletV3
// HighloadV3 is optimized for high-throughput scenarios and doesn't have sequential nonce bottlenecks
// mnemonic should be a 24-word seed phrase
// chainState contains the contract addresses (Router, FeeQuoter, etc.)
func NewTonSourceManager(
	ctx context.Context,
	l logger.Logger,
	chainSel uint64,
	endpoint string,
	mnemonic string,
	chainState tonstate.CCIPChainState,
) (*TonSourceManager, error) {
	if mnemonic == "" {
		return nil, fmt.Errorf("TON mnemonic is required for source chain %d", chainSel)
	}

	client, err := connectTonClient(ctx, endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to TON client: %w", err)
	}

	// Use custom HighloadV3 config with atomic counter for unique query IDs
	// This allows parallel message sending without query_id conflicts
	w, err := wallet.FromSeed(client, strings.Fields(mnemonic), loadTestHighloadV3Config)
	if err != nil {
		return nil, fmt.Errorf("failed to create HighloadV3 wallet from mnemonic: %w", err)
	}

	// Log wallet balance for verification
	mc, _ := client.CurrentMasterchainInfo(ctx)
	balance, _ := w.GetBalance(ctx, mc)
	l.Infow("TON HighloadV3 wallet initialized for load testing",
		"chainSelector", chainSel,
		"address", w.Address().String(),
		"balance", balance.String(),
		"walletType", "HighloadV3",
		"router", chainState.Router.String(),
		"feeQuoter", chainState.FeeQuoter.String())

	return &TonSourceManager{
		l:          l,
		client:     client,
		wallet:     w,
		chainSel:   chainSel,
		chainState: chainState,
	}, nil
}

// TODO: (2026-01-12@jadepark-dev) import from chainlink-ton
// connectTonClient connects to TON network via liteserver or config URL
func connectTonClient(ctx context.Context, endpoint string) (*ton.APIClient, error) {
	if strings.HasPrefix(endpoint, "liteserver://") {
		pool, err := chain.CreateLiteserverConnectionPool(ctx, endpoint)
		if err != nil {
			return nil, fmt.Errorf("failed to create liteserver connection pool: %w", err)
		}
		return ton.NewAPIClient(pool, ton.ProofCheckPolicyFast), nil
	}

	// Connect via config URL (e.g., https://ton.org/testnet-global.config.json)
	cfg, err := liteclient.GetConfigFromUrl(ctx, endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to get TON config from URL: %w", err)
	}
	pool := liteclient.NewConnectionPool()
	err = pool.AddConnectionsFromConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to TON: %w", err)
	}
	return ton.NewAPIClient(pool, ton.ProofCheckPolicyFast), nil
}

// sendTONMessage sends a CCIP message from TON to a destination chain
// The mutex ensures only one message is sent at a time to prevent transaction conflicts
// when multiple destination guns try to send from the same TON wallet simultaneously.
// This is necessary because TON's HighloadV3 wallet can have issues when multiple
// external messages arrive in the same block with conflicting query IDs.
func (m *TonSourceManager) sendTONMessage(
	ctx context.Context,
	destChainSel uint64,
	receiver []byte,
	testConfig *ccip.LoadConfig,
) (uint64, string, error) {
	// Lock to prevent concurrent sends from the same wallet
	// This ensures messages are sent sequentially, avoiding query ID conflicts
	m.sendMu.Lock()
	defer m.sendMu.Unlock()

	// Use the chainState stored in manager (loaded from address book during initialization)
	// NOT from state.TonChains which has empty addresses due to loading issues

	// Build extra args
	extraArgs := onramp.GenericExtraArgsV2{
		GasLimit:                 big.NewInt(tonDefaultGasLimit),
		AllowOutOfOrderExecution: *testConfig.OOOExecution,
	}
	extraArgsCell, err := tlb.ToCell(extraArgs)
	if err != nil {
		return 0, "", fmt.Errorf("failed to serialize TON extra args: %w", err)
	}

	// Prepare receiver bytes (left-pad for EVM addresses)
	receiverBytes := receiver
	if len(receiverBytes) < 32 {
		paddedReceiver := make([]byte, 32)
		copy(paddedReceiver[32-len(receiverBytes):], receiverBytes)
		receiverBytes = paddedReceiver
	}

	// Generate random message data
	dataLength := 10 // default
	if testConfig.MessageDetails != nil && len(*testConfig.MessageDetails) > 0 {
		for _, msg := range *testConfig.MessageDetails {
			if msg.DataLengthBytes != nil && *msg.DataLengthBytes > 0 {
				dataLength = *msg.DataLengthBytes
				break
			}
		}
	}
	data := make([]byte, dataLength)
	_, err = rand.Read(data)
	if err != nil {
		return 0, "", fmt.Errorf("failed to generate random data: %w", err)
	}

	// Build CCIP message
	ccipSend := tonrouter.CCIPSend{
		QueryID:           0,
		DestChainSelector: destChainSel,
		Receiver:          receiverBytes,
		Data:              data,
		TokenAmounts:      nil, // No token transfers for messaging-only
		FeeToken:          tvm.TonTokenAddr,
		ExtraArgs:         extraArgsCell,
	}

	// Build minimal environment with initialized client/wallet (same pattern as staging-monitor)
	tonProvider := &cldfton.Chain{
		ChainMetadata: cldfton.ChainMetadata{
			Selector: m.chainSel,
		},
		Client:        m.client,
		Wallet:        m.wallet,
		WalletAddress: m.wallet.WalletAddress(),
	}
	blockchains := cldfchain.NewBlockChainsFromSlice([]cldfchain.BlockChain{tonProvider})
	env := cldf.Environment{
		GetContext:  func() context.Context { return ctx },
		Logger:      m.l,
		BlockChains: blockchains,
	}

	// Send using existing deployment ops - use m.chainState which has the correct addresses
	seqNum, event, err := tonOps.SendCCIPMessage(env, m.chainState, m.chainSel, ccipSend)
	if err != nil {
		return 0, "", fmt.Errorf("failed to send TON CCIP message: %w", err)
	}

	// Extract message ID from event
	ccipEvent, ok := event.(onramp.CCIPMessageSent)
	if !ok {
		return seqNum, "", fmt.Errorf("unexpected event type: %T", event)
	}

	messageID := hex.EncodeToString(ccipEvent.Message.Header.MessageID)
	m.l.Infow("CCIP message sent from TON",
		"sourceChain", m.chainSel,
		"destChain", destChainSel,
		"seqNum", seqNum,
		"messageID", messageID)

	return seqNum, messageID, nil
}

// initializeTonSourceKeys initializes TON source managers for all TON chains
// It loads TON contract addresses directly from the address book
func initializeTonSourceKeys(
	ctx context.Context,
	l logger.Logger,
	env *cldf.Environment,
	tonMnemonic string,
) (map[uint64]*TonSourceManager, error) {
	tonSourceKeys := make(map[uint64]*TonSourceManager)

	// Get all TON chain selectors from environment
	for chainSel := range env.BlockChains.TonChains() {
		// Load addresses directly from address book
		addresses, err := env.ExistingAddresses.AddressesForChain(chainSel)
		if err != nil {
			l.Warnw("Failed to get addresses for TON chain from address book",
				"chainSelector", chainSel,
				"error", err)
			continue
		}

		// Build chainState from address book (same logic as ton/state.go loadChainState)
		chainState := tonstate.CCIPChainState{}
		for addrStr, tv := range addresses {
			addr, err := address.ParseAddr(addrStr)
			if err != nil {
				l.Warnw("Failed to parse TON address", "address", addrStr, "error", err)
				continue
			}
			switch tv.Type {
			case commontypes.LinkToken:
				chainState.LinkTokenAddress = *addr
			case shared.OffRamp:
				chainState.OffRamp = *addr
			case shared.Router:
				chainState.Router = *addr
			case shared.OnRamp:
				chainState.OnRamp = *addr
			case shared.FeeQuoter:
				chainState.FeeQuoter = *addr
			}
		}

		l.Infow("Loaded TON chain state from address book",
			"chainSelector", chainSel,
			"router", chainState.Router.String(),
			"feeQuoter", chainState.FeeQuoter.String(),
			"onRamp", chainState.OnRamp.String(),
			"offRamp", chainState.OffRamp.String())

		// Get endpoint - use liteserver from chains-details.json if available
		endpoint := "https://ton.org/testnet-global.config.json"

		manager, err := NewTonSourceManager(ctx, l, chainSel, endpoint, tonMnemonic, chainState)
		if err != nil {
			l.Warnw("Failed to initialize TON source manager",
				"chainSelector", chainSel,
				"error", err)
			continue
		}

		tonSourceKeys[chainSel] = manager
		l.Infow("Initialized TON source key",
			"chainSelector", chainSel,
			"wallet", manager.wallet.Address().String())
	}

	return tonSourceKeys, nil
}
