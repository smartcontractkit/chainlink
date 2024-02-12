package contracts

import (
	"context"
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/arm_contract"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/commit_store"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/evm_2_evm_offramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/evm_2_evm_onramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/lock_release_token_pool"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/maybe_revert_message_receiver"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/mock_arm_contract"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/mock_usdc_token_transmitter"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/mock_v3_aggregator_contract"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/price_registry"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/token_pool"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/usdc_token_pool"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/link_token_interface"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/shared/generated/burn_mint_erc677"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/shared/generated/erc20"
)

var (
	FiftyCoins   = new(big.Int).Mul(big.NewInt(1e18), big.NewInt(50))
	HundredCoins = new(big.Int).Mul(big.NewInt(1e18), big.NewInt(100))
)

const (
	Network = "Network Name"
)

type RateLimiterConfig struct {
	Rate     *big.Int
	Capacity *big.Int
}

type ARMConfig struct {
	ARMWeightsByParticipants map[string]*big.Int // mapping : ARM participant address => weight
	ThresholdForBlessing     *big.Int
	ThresholdForBadSignal    *big.Int
}

type TokenTransmitter struct {
	client          blockchain.EVMClient
	instance        *mock_usdc_token_transmitter.MockE2EUSDCTransmitter
	ContractAddress common.Address
}

type ERC677Token struct {
	client          blockchain.EVMClient
	instance        *burn_mint_erc677.BurnMintERC677
	ContractAddress common.Address
}

func (token *ERC677Token) GrantMintAndBurn(burnAndMinter common.Address) error {
	opts, err := token.client.TransactionOpts(token.client.GetDefaultWallet())
	if err != nil {
		return fmt.Errorf("failed to get transaction opts: %w", err)
	}
	log.Info().
		Str(Network, token.client.GetNetworkName()).
		Str("BurnAndMinter", burnAndMinter.Hex()).
		Str("Token", token.ContractAddress.Hex()).
		Msg("Granting mint and burn roles")
	tx, err := token.instance.GrantMintAndBurnRoles(opts, burnAndMinter)
	if err != nil {
		return fmt.Errorf("failed to grant mint and burn roles: %w", err)
	}
	return token.client.ProcessTransaction(tx)
}

func (token *ERC677Token) GrantMintRole(minter common.Address) error {
	opts, err := token.client.TransactionOpts(token.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	log.Info().
		Str(Network, token.client.GetNetworkName()).
		Str("Minter", minter.Hex()).
		Str("Token", token.ContractAddress.Hex()).
		Msg("Granting mint roles")
	tx, err := token.instance.GrantMintRole(opts, minter)
	if err != nil {
		return fmt.Errorf("failed to grant mint role: %w", err)
	}
	return token.client.ProcessTransaction(tx)
}

func (token *ERC677Token) Mint(to common.Address, amount *big.Int) error {
	opts, err := token.client.TransactionOpts(token.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	log.Info().
		Str(Network, token.client.GetNetworkName()).
		Str("To", to.Hex()).
		Str("Token", token.ContractAddress.Hex()).
		Str("Amount", amount.String()).
		Msg("Minting tokens")
	tx, err := token.instance.Mint(opts, to, amount)
	if err != nil {
		return fmt.Errorf("failed to mint tokens: %w", err)
	}
	return token.client.ProcessTransaction(tx)
}

type ERC20Token struct {
	client          blockchain.EVMClient
	instance        *erc20.ERC20
	ContractAddress common.Address
}

func (token *ERC20Token) Address() string {
	return token.ContractAddress.Hex()
}

func (token *ERC20Token) BalanceOf(ctx context.Context, addr string) (*big.Int, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(token.client.GetDefaultWallet().Address()),
		Context: ctx,
	}
	balance, err := token.instance.BalanceOf(opts, common.HexToAddress(addr))
	if err != nil {
		return nil, fmt.Errorf("failed to get balance: %w", err)
	}
	return balance, nil
}

func (token *ERC20Token) Allowance(owner, spender string) (*big.Int, error) {
	allowance, err := token.instance.Allowance(nil, common.HexToAddress(owner), common.HexToAddress(spender))
	if err != nil {
		return nil, err
	}
	return allowance, nil
}

func (token *ERC20Token) Approve(to string, amount *big.Int) error {
	opts, err := token.client.TransactionOpts(token.client.GetDefaultWallet())
	if err != nil {
		return fmt.Errorf("failed to get transaction options: %w", err)
	}
	log.Info().
		Str("From", token.client.GetDefaultWallet().Address()).
		Str("To", to).
		Str("Token", token.Address()).
		Str("Amount", amount.String()).
		Uint64("Nonce", opts.Nonce.Uint64()).
		Str(Network, token.client.GetNetworkConfig().Name).
		Msg("Approving ERC20 Transfer")
	tx, err := token.instance.Approve(opts, common.HexToAddress(to), amount)
	if err != nil {
		return fmt.Errorf("failed to approve ERC20: %w", err)
	}
	return token.client.ProcessTransaction(tx)
}

func (token *ERC20Token) Transfer(to string, amount *big.Int) error {
	opts, err := token.client.TransactionOpts(token.client.GetDefaultWallet())
	if err != nil {
		return fmt.Errorf("failed to get transaction options: %w", err)
	}
	log.Info().
		Str("From", token.client.GetDefaultWallet().Address()).
		Str("To", to).
		Str("Amount", amount.String()).
		Uint64("Nonce", opts.Nonce.Uint64()).
		Str(Network, token.client.GetNetworkConfig().Name).
		Msg("Transferring ERC20")
	tx, err := token.instance.Transfer(opts, common.HexToAddress(to), amount)
	if err != nil {
		return fmt.Errorf("failed to transfer ERC20: %w", err)
	}
	return token.client.ProcessTransaction(tx)
}

type LinkToken struct {
	client     blockchain.EVMClient
	instance   *link_token_interface.LinkToken
	EthAddress common.Address
}

func (l *LinkToken) Address() string {
	return l.EthAddress.Hex()
}

func (l *LinkToken) BalanceOf(ctx context.Context, addr string) (*big.Int, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(l.client.GetDefaultWallet().Address()),
		Context: ctx,
	}
	balance, err := l.instance.BalanceOf(opts, common.HexToAddress(addr))
	if err != nil {
		return nil, fmt.Errorf("failed to get LINK balance: %w", err)
	}
	return balance, nil
}

func (l *LinkToken) Allowance(owner, spender string) (*big.Int, error) {
	allowance, err := l.instance.Allowance(nil, common.HexToAddress(owner), common.HexToAddress(spender))
	if err != nil {
		return nil, err
	}
	return allowance, nil
}

func (l *LinkToken) Approve(to string, amount *big.Int) error {
	opts, err := l.client.TransactionOpts(l.client.GetDefaultWallet())
	if err != nil {
		return fmt.Errorf("failed to get transaction opts: %w", err)
	}
	log.Info().
		Str("From", l.client.GetDefaultWallet().Address()).
		Str("To", to).
		Str("Token", l.Address()).
		Str("Amount", amount.String()).
		Uint64("Nonce", opts.Nonce.Uint64()).
		Str(Network, l.client.GetNetworkConfig().Name).
		Msg("Approving LINK Transfer")
	tx, err := l.instance.Approve(opts, common.HexToAddress(to), amount)
	if err != nil {
		return fmt.Errorf("failed to approve LINK transfer: %w", err)
	}
	return l.client.ProcessTransaction(tx)
}

func (l *LinkToken) Transfer(to string, amount *big.Int) error {
	opts, err := l.client.TransactionOpts(l.client.GetDefaultWallet())
	if err != nil {
		return fmt.Errorf("failed to get transaction opts: %w", err)
	}
	log.Info().
		Str("From", l.client.GetDefaultWallet().Address()).
		Str("To", to).
		Str("Amount", amount.String()).
		Uint64("Nonce", opts.Nonce.Uint64()).
		Str(Network, l.client.GetNetworkConfig().Name).
		Msg("Transferring LINK")
	tx, err := l.instance.Transfer(opts, common.HexToAddress(to), amount)
	if err != nil {
		return fmt.Errorf("failed to transfer LINK: %w", err)
	}
	return l.client.ProcessTransaction(tx)
}

// TokenPool represents a TokenPool address
type TokenPool struct {
	client          blockchain.EVMClient
	PoolInterface   *token_pool.TokenPool
	LockReleasePool *lock_release_token_pool.LockReleaseTokenPool
	USDCPool        *usdc_token_pool.USDCTokenPool
	EthAddress      common.Address
}

func (pool *TokenPool) Address() string {
	return pool.EthAddress.Hex()
}

func (pool *TokenPool) SyncUSDCDomain(destTokenTransmitter *TokenTransmitter, destPoolAddr common.Address, destChainSelector uint64) error {
	if pool.USDCPool == nil {
		return fmt.Errorf("USDCPool is nil")
	}

	var allowedCallerBytes [32]byte
	copy(allowedCallerBytes[12:], destPoolAddr.Bytes())
	destTokenTransmitterIns, err := mock_usdc_token_transmitter.NewMockE2EUSDCTransmitter(destTokenTransmitter.ContractAddress, destTokenTransmitter.client.Backend())
	if err != nil {
		return fmt.Errorf("failed to create mock USDC token transmitter: %w", err)
	}
	domain, err := destTokenTransmitterIns.LocalDomain(nil)
	if err != nil {
		return fmt.Errorf("failed to get local domain: %w", err)
	}
	opts, err := pool.client.TransactionOpts(pool.client.GetDefaultWallet())
	if err != nil {
		return fmt.Errorf("failed to get transaction opts: %w", err)
	}
	log.Info().
		Str("Token Pool", pool.Address()).
		Str(Network, pool.client.GetNetworkName()).
		Uint32("Domain", domain).
		Str("Allowed Caller", destPoolAddr.Hex()).
		Str("Dest Chain Selector", fmt.Sprintf("%d", destChainSelector)).
		Msg("Syncing USDC Domain")
	tx, err := pool.USDCPool.SetDomains(opts, []usdc_token_pool.USDCTokenPoolDomainUpdate{
		{
			AllowedCaller:     allowedCallerBytes,
			DomainIdentifier:  domain,
			DestChainSelector: destChainSelector,
			Enabled:           true,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to set domain: %w", err)
	}
	return pool.client.ProcessTransaction(tx)
}

func (pool *TokenPool) RemoveLiquidity(amount *big.Int) error {
	if pool.LockReleasePool == nil {
		return fmt.Errorf("LockReleasePool is nil")
	}
	opts, err := pool.client.TransactionOpts(pool.client.GetDefaultWallet())
	if err != nil {
		return fmt.Errorf("failed to get transaction opts: %w", err)
	}
	log.Info().
		Str("Token Pool", pool.Address()).
		Str("Amount", amount.String()).
		Msg("Initiating removing funds from pool")
	tx, err := pool.LockReleasePool.WithdrawLiquidity(opts, amount)
	if err != nil {
		return fmt.Errorf("failed to withdraw liquidity: %w", err)
	}
	log.Info().
		Str("Token Pool", pool.Address()).
		Str("Amount", amount.String()).
		Str(Network, pool.client.GetNetworkConfig().Name).
		Msg("Liquidity removed")
	return pool.client.ProcessTransaction(tx)
}

type tokenApproveFn func(string, *big.Int) error

func (pool *TokenPool) AddLiquidity(approveFn tokenApproveFn, tokenAddr string, amount *big.Int) error {
	if pool.LockReleasePool == nil {
		return fmt.Errorf("cannot add liquidity to pool")
	}
	log.Info().
		Str("Link Token", tokenAddr).
		Str("Token Pool", pool.Address()).
		Msg("Initiating transferring of token to token pool")
	err := approveFn(pool.Address(), amount)
	if err != nil {
		return fmt.Errorf("failed to approve token transfer: %w", err)
	}
	err = pool.client.WaitForEvents()
	if err != nil {
		return fmt.Errorf("failed to wait for events: %w", err)
	}
	opts, err := pool.client.TransactionOpts(pool.client.GetDefaultWallet())
	if err != nil {
		return fmt.Errorf("failed to get transaction opts: %w", err)
	}
	_, err = pool.LockReleasePool.SetRebalancer(opts, opts.From)
	if err != nil {
		return fmt.Errorf("failed to set rebalancer: %w", err)
	}
	opts, err = pool.client.TransactionOpts(pool.client.GetDefaultWallet())
	if err != nil {
		return fmt.Errorf("failed to get transaction opts: %w", err)
	}
	log.Info().
		Str("Token Pool", pool.Address()).
		Msg("Initiating adding Tokens in pool")
	tx, err := pool.LockReleasePool.ProvideLiquidity(opts, amount)
	if err != nil {
		return fmt.Errorf("failed to provide liquidity: %w", err)
	}
	log.Info().
		Str("Token Pool", pool.Address()).
		Str("Link Token", tokenAddr).
		Str(Network, pool.client.GetNetworkConfig().Name).
		Msg("Liquidity added")
	return pool.client.ProcessTransaction(tx)
}

func (pool *TokenPool) SetRemoteChainOnPool(remoteChainSelector uint64) error {
	log.Info().
		Str("Token Pool", pool.Address()).
		Msg("Setting remote chain on pool")
	isSupported, err := pool.PoolInterface.IsSupportedChain(nil, remoteChainSelector)
	if err != nil {
		return fmt.Errorf("failed to get if chain is supported: %w", err)
	}
	// Check if remote chain is already supported , if yes return
	if isSupported {
		log.Info().
			Str("Token Pool", pool.Address()).
			Str(Network, pool.client.GetNetworkName()).
			Uint64("Remote Chain Selector", remoteChainSelector).
			Msg("Remote chain is already supported")
		return nil
	}
	// If remote chain is not supported , add it
	opts, err := pool.client.TransactionOpts(pool.client.GetDefaultWallet())
	if err != nil {
		return fmt.Errorf("failed to get transaction opts: %w", err)
	}
	tx, err := pool.PoolInterface.ApplyChainUpdates(opts, []token_pool.TokenPoolChainUpdate{
		{
			RemoteChainSelector: remoteChainSelector,
			Allowed:             true,
			InboundRateLimiterConfig: token_pool.RateLimiterConfig{
				IsEnabled: true,
				Capacity:  new(big.Int).Mul(big.NewInt(1e18), big.NewInt(1e9)),
				Rate:      new(big.Int).Mul(big.NewInt(1e18), big.NewInt(1e5)),
			},
			OutboundRateLimiterConfig: token_pool.RateLimiterConfig{
				IsEnabled: true,
				Capacity:  new(big.Int).Mul(big.NewInt(1e18), big.NewInt(1e9)),
				Rate:      new(big.Int).Mul(big.NewInt(1e18), big.NewInt(1e5)),
			},
		},
	})

	if err != nil {
		return fmt.Errorf("failed to set chain updates on token pool: %w", err)
	}

	log.Info().
		Str("Token Pool", pool.Address()).
		Str("Chain selector", strconv.FormatUint(remoteChainSelector, 10)).
		Str(Network, pool.client.GetNetworkConfig().Name).
		Msg("Remote chain set on token pool")
	return pool.client.ProcessTransaction(tx)
}

func (pool *TokenPool) SetRemoteChainRateLimits(remoteChainSelector uint64, rl token_pool.RateLimiterConfig) error {
	opts, err := pool.client.TransactionOpts(pool.client.GetDefaultWallet())
	if err != nil {
		return fmt.Errorf("error getting transaction opts: %w", err)
	}
	log.Info().
		Str("Token Pool", pool.Address()).
		Str("Remote chain selector", strconv.FormatUint(remoteChainSelector, 10)).
		Interface("RateLimiterConfig", rl).
		Msg("Setting Rate Limit on token pool")
	tx, err := pool.PoolInterface.SetChainRateLimiterConfig(opts, remoteChainSelector, rl, rl)

	if err != nil {
		return fmt.Errorf("error setting rate limit token pool: %w", err)
	}

	log.Info().
		Str("Token Pool", pool.Address()).
		Str("Remote chain selector", strconv.FormatUint(remoteChainSelector, 10)).
		Interface("RateLimiterConfig", rl).
		Msg("Rate Limit on token pool is set")
	return pool.client.ProcessTransaction(tx)
}

type ARM struct {
	client     blockchain.EVMClient
	Instance   *arm_contract.ARMContract
	EthAddress common.Address
}

func (arm *ARM) Address() string {
	return arm.EthAddress.Hex()
}

type MockARM struct {
	client     blockchain.EVMClient
	Instance   *mock_arm_contract.MockARMContract
	EthAddress common.Address
}

func (arm *MockARM) SetClient(client blockchain.EVMClient) {
	arm.client = client
}
func (arm *MockARM) Address() string {
	return arm.EthAddress.Hex()
}

type CommitStore struct {
	client     blockchain.EVMClient
	Instance   *commit_store.CommitStore
	EthAddress common.Address
}

func (b *CommitStore) Address() string {
	return b.EthAddress.Hex()
}

// SetOCR2Config sets the offchain reporting protocol configuration
func (b *CommitStore) SetOCR2Config(
	signers []common.Address,
	transmitters []common.Address,
	f uint8,
	onchainConfig []byte,
	offchainConfigVersion uint64,
	offchainConfig []byte,
) error {
	log.Info().Str("Contract Address", b.Address()).Msg("Configuring OCR config for CommitStore Contract")
	// Set Config
	opts, err := b.client.TransactionOpts(b.client.GetDefaultWallet())
	if err != nil {
		return fmt.Errorf("error getting transaction opts: %w", err)
	}

	log.Info().
		Interface("signerAddresses", signers).
		Interface("transmitterAddresses", transmitters).
		Str(Network, b.client.GetNetworkConfig().Name).
		Msg("Configuring CommitStore")
	tx, err := b.Instance.SetOCR2Config(
		opts,
		signers,
		transmitters,
		f,
		onchainConfig,
		offchainConfigVersion,
		offchainConfig,
	)

	if err != nil {
		return fmt.Errorf("error setting OCR2 config: %w", err)
	}
	return b.client.ProcessTransaction(tx)
}

type ReceiverDapp struct {
	client     blockchain.EVMClient
	instance   *maybe_revert_message_receiver.MaybeRevertMessageReceiver
	EthAddress common.Address
}

func (rDapp *ReceiverDapp) Address() string {
	return rDapp.EthAddress.Hex()
}

func (rDapp *ReceiverDapp) ToggleRevert(revert bool) error {
	opts, err := rDapp.client.TransactionOpts(rDapp.client.GetDefaultWallet())
	if err != nil {
		return fmt.Errorf("error getting transaction opts: %w", err)
	}
	tx, err := rDapp.instance.SetRevert(opts, revert)
	if err != nil {
		return fmt.Errorf("error setting revert: %w", err)
	}
	log.Info().
		Bool("revert", revert).
		Str("tx", tx.Hash().String()).
		Str("ReceiverDapp", rDapp.Address()).
		Str(Network, rDapp.client.GetNetworkConfig().Name).
		Msg("ReceiverDapp revert set")
	return rDapp.client.ProcessTransaction(tx)
}

type PriceRegistry struct {
	client     blockchain.EVMClient
	Instance   *price_registry.PriceRegistry
	EthAddress common.Address
}

func (c *PriceRegistry) Address() string {
	return c.EthAddress.Hex()
}

func (c *PriceRegistry) AddPriceUpdater(addr common.Address) error {
	opts, err := c.client.TransactionOpts(c.client.GetDefaultWallet())
	if err != nil {
		return fmt.Errorf("error getting transaction opts: %w", err)
	}
	tx, err := c.Instance.ApplyPriceUpdatersUpdates(opts, []common.Address{addr}, []common.Address{})
	if err != nil {
		return fmt.Errorf("error adding price updater: %w", err)
	}
	log.Info().
		Str("updaters", addr.Hex()).
		Str(Network, c.client.GetNetworkConfig().Name).
		Msg("PriceRegistry updater added")
	return c.client.ProcessTransaction(tx)
}

func (c *PriceRegistry) AddFeeToken(addr common.Address) error {
	opts, err := c.client.TransactionOpts(c.client.GetDefaultWallet())
	if err != nil {
		return fmt.Errorf("error getting transaction opts: %w", err)
	}
	tx, err := c.Instance.ApplyFeeTokensUpdates(opts, []common.Address{addr}, []common.Address{})
	if err != nil {
		return fmt.Errorf("error adding fee token: %w", err)
	}
	log.Info().
		Str("feeTokens", addr.Hex()).
		Str(Network, c.client.GetNetworkConfig().Name).
		Msg("PriceRegistry feeToken set")
	return c.client.ProcessTransaction(tx)
}

func (c *PriceRegistry) UpdatePrices(priceUpdates price_registry.InternalPriceUpdates) error {
	opts, err := c.client.TransactionOpts(c.client.GetDefaultWallet())
	if err != nil {
		return fmt.Errorf("error getting transaction opts: %w", err)
	}
	tx, err := c.Instance.UpdatePrices(opts, priceUpdates)
	if err != nil {
		return fmt.Errorf("error updating prices: %w", err)
	}
	log.Info().
		Str(Network, c.client.GetNetworkConfig().Name).
		Interface("PriceUpdates", priceUpdates).
		Msg("Prices updated")
	return c.client.ProcessTransaction(tx)
}

type Router struct {
	client     blockchain.EVMClient
	Instance   *router.Router
	EthAddress common.Address
}

func (r *Router) Address() string {
	return r.EthAddress.Hex()
}

func (r *Router) SetOnRamp(chainSelector uint64, onRamp common.Address) error {
	opts, err := r.client.TransactionOpts(r.client.GetDefaultWallet())
	if err != nil {
		return fmt.Errorf("error getting transaction opts: %w", err)
	}
	log.Info().
		Str("Router", r.Address()).
		Str("OnRamp", onRamp.Hex()).
		Str(Network, r.client.GetNetworkName()).
		Str("ChainSelector", strconv.FormatUint(chainSelector, 10)).
		Msg("Setting on ramp for r")

	tx, err := r.Instance.ApplyRampUpdates(opts, []router.RouterOnRamp{{DestChainSelector: chainSelector, OnRamp: onRamp}}, nil, nil)
	if err != nil {
		return fmt.Errorf("error applying ramp updates: %w", err)
	}
	log.Info().
		Str("onRamp", onRamp.Hex()).
		Str("Network Name", r.client.GetNetworkConfig().Name).
		Msg("Router is configured")
	return r.client.ProcessTransaction(tx)
}

func (r *Router) CCIPSend(destChainSelector uint64, msg router.ClientEVM2AnyMessage, valueForNative *big.Int) (*types.Transaction, error) {
	opts, err := r.client.TransactionOpts(r.client.GetDefaultWallet())
	if err != nil {
		return nil, fmt.Errorf("error getting transaction opts: %w", err)
	}
	if valueForNative != nil {
		opts.Value = valueForNative
	}

	log.Info().
		Str(Network, r.client.GetNetworkName()).
		Str("Router", r.Address()).
		Interface("TokensAndAmounts", msg.TokenAmounts).
		Str("FeeToken", msg.FeeToken.Hex()).
		Str("ExtraArgs", fmt.Sprintf("0x%x", msg.ExtraArgs[:])).
		Str("Receiver", fmt.Sprintf("0x%x", msg.Receiver[:])).
		Msg("Sending msg")
	return r.Instance.CcipSend(opts, destChainSelector, msg)
}

func (r *Router) CCIPSendAndProcessTx(destChainSelector uint64, msg router.ClientEVM2AnyMessage, valueForNative *big.Int) (*types.Transaction, error) {
	tx, err := r.CCIPSend(destChainSelector, msg, valueForNative)
	if err != nil {
		return nil, fmt.Errorf("failed to send msg: %w", err)
	}
	log.Info().
		Str("router", r.Address()).
		Str("txHash", tx.Hash().Hex()).
		Str(Network, r.client.GetNetworkConfig().Name).
		Str("chain selector", strconv.FormatUint(destChainSelector, 10)).
		Msg("msg is sent")
	return tx, r.client.ProcessTransaction(tx)
}

func (r *Router) AddOffRamp(offRamp common.Address, sourceChainId uint64) (*types.Transaction, error) {
	opts, err := r.client.TransactionOpts(r.client.GetDefaultWallet())
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction opts: %w", err)
	}
	tx, err := r.Instance.ApplyRampUpdates(opts, nil, nil, []router.RouterOffRamp{{SourceChainSelector: sourceChainId, OffRamp: offRamp}})
	if err != nil {
		return nil, fmt.Errorf("failed to add offRamp: %w", err)
	}
	log.Info().
		Str("offRamp", offRamp.Hex()).
		Str(Network, r.client.GetNetworkConfig().Name).
		Msg("offRamp is added to Router")
	return tx, r.client.ProcessTransaction(tx)
}

func (r *Router) SetWrappedNative(wNative common.Address) (*types.Transaction, error) {
	opts, err := r.client.TransactionOpts(r.client.GetDefaultWallet())
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction opts: %w", err)
	}
	tx, err := r.Instance.SetWrappedNative(opts, wNative)
	if err != nil {
		return nil, fmt.Errorf("failed to set wrapped native: %w", err)
	}
	log.Info().
		Str("wrapped native", wNative.Hex()).
		Str("router", r.Address()).
		Str(Network, r.client.GetNetworkConfig().Name).
		Msg("wrapped native is added for Router")
	return tx, r.client.ProcessTransaction(tx)
}

func (r *Router) GetFee(destChainSelector uint64, message router.ClientEVM2AnyMessage) (*big.Int, error) {
	return r.Instance.GetFee(nil, destChainSelector, message)
}

type OnRamp struct {
	client     blockchain.EVMClient
	Instance   *evm_2_evm_onramp.EVM2EVMOnRamp
	EthAddress common.Address
}

func (onRamp *OnRamp) Address() string {
	return onRamp.EthAddress.Hex()
}

func (onRamp *OnRamp) SetNops() error {
	opts, err := onRamp.client.TransactionOpts(onRamp.client.GetDefaultWallet())
	if err != nil {
		return fmt.Errorf("failed to get transaction opts: %w", err)
	}
	owner := common.HexToAddress(onRamp.client.GetDefaultWallet().Address())
	// set the payee to the default wallet
	tx, err := onRamp.Instance.SetNops(opts, []evm_2_evm_onramp.EVM2EVMOnRampNopAndWeight{{
		Nop:    owner,
		Weight: 1,
	}})
	if err != nil {
		return fmt.Errorf("failed to set nops: %w", err)
	}
	return onRamp.client.ProcessTransaction(tx)
}

func (onRamp *OnRamp) SetTokenTransferFeeConfig(tokenTransferFeeConfig []evm_2_evm_onramp.EVM2EVMOnRampTokenTransferFeeConfigArgs) error {
	opts, err := onRamp.client.TransactionOpts(onRamp.client.GetDefaultWallet())
	if err != nil {
		return fmt.Errorf("failed to get transaction opts: %w", err)
	}
	tx, err := onRamp.Instance.SetTokenTransferFeeConfig(opts, tokenTransferFeeConfig)
	if err != nil {
		return fmt.Errorf("failed to set token transfer fee config: %w", err)
	}
	log.Info().
		Interface("tokenTransferFeeConfig", tokenTransferFeeConfig).
		Str("onRamp", onRamp.Address()).
		Str(Network, onRamp.client.GetNetworkConfig().Name).
		Msg("TokenTransferFeeConfig set in OnRamp")
	return onRamp.client.ProcessTransaction(tx)
}

func (onRamp *OnRamp) ApplyPoolUpdates(poolUpdates []evm_2_evm_onramp.InternalPoolUpdate) error {
	opts, err := onRamp.client.TransactionOpts(onRamp.client.GetDefaultWallet())
	if err != nil {
		return fmt.Errorf("failed to get transaction opts: %w", err)
	}
	tx, err := onRamp.Instance.ApplyPoolUpdates(opts, []evm_2_evm_onramp.InternalPoolUpdate{}, poolUpdates)
	if err != nil {
		return fmt.Errorf("failed to apply pool updates: %w", err)
	}
	log.Info().
		Interface("poolUpdates", poolUpdates).
		Str("onRamp", onRamp.Address()).
		Str(Network, onRamp.client.GetNetworkConfig().Name).
		Msg("poolUpdates set in OnRamp")
	return onRamp.client.ProcessTransaction(tx)
}

func (onRamp *OnRamp) PayNops() error {
	opts, err := onRamp.client.TransactionOpts(onRamp.client.GetDefaultWallet())
	if err != nil {
		return fmt.Errorf("failed to get transaction opts: %w", err)
	}
	tx, err := onRamp.Instance.PayNops(opts)
	if err != nil {
		return fmt.Errorf("failed to pay nops: %w", err)
	}
	return onRamp.client.ProcessTransaction(tx)
}

func (onRamp *OnRamp) WithdrawNonLinkFees(wrappedNative common.Address) error {
	opts, err := onRamp.client.TransactionOpts(onRamp.client.GetDefaultWallet())
	if err != nil {
		return fmt.Errorf("failed to get transaction opts: %w", err)
	}
	owner := common.HexToAddress(onRamp.client.GetDefaultWallet().Address())
	tx, err := onRamp.Instance.WithdrawNonLinkFees(opts, wrappedNative, owner)
	if err != nil {
		return fmt.Errorf("failed to withdraw non link fees: %w", err)
	}
	return onRamp.client.ProcessTransaction(tx)
}

func (onRamp *OnRamp) SetRateLimit(rlConfig evm_2_evm_onramp.RateLimiterConfig) error {
	opts, err := onRamp.client.TransactionOpts(onRamp.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := onRamp.Instance.SetRateLimiterConfig(opts, rlConfig)
	if err != nil {
		return fmt.Errorf("failed to set rate limit: %w", err)
	}
	log.Info().
		Bool("Enabled", rlConfig.IsEnabled).
		Str("capacity", rlConfig.Capacity.String()).
		Str("rate", rlConfig.Rate.String()).
		Str("onRamp", onRamp.Address()).
		Str(Network, onRamp.client.GetNetworkConfig().Name).
		Msg("Setting Rate limit in OnRamp")
	return onRamp.client.ProcessTransaction(tx)
}

type OffRamp struct {
	client     blockchain.EVMClient
	Instance   *evm_2_evm_offramp.EVM2EVMOffRamp
	EthAddress common.Address
}

func (offRamp *OffRamp) Address() string {
	return offRamp.EthAddress.Hex()
}

// SetOCR2Config sets the offchain reporting protocol configuration
func (offRamp *OffRamp) SetOCR2Config(
	signers []common.Address,
	transmitters []common.Address,
	f uint8,
	onchainConfig []byte,
	offchainConfigVersion uint64,
	offchainConfig []byte,
) error {
	log.Info().Str("Contract Address", offRamp.Address()).Msg("Configuring OffRamp Contract")
	// Set Config
	opts, err := offRamp.client.TransactionOpts(offRamp.client.GetDefaultWallet())
	if err != nil {
		return fmt.Errorf("failed to get transaction options: %w", err)
	}
	log.Info().
		Interface("signerAddresses", signers).
		Interface("transmitterAddresses", transmitters).
		Str(Network, offRamp.client.GetNetworkConfig().Name).
		Msg("Configuring OffRamp")
	tx, err := offRamp.Instance.SetOCR2Config(
		opts,
		signers,
		transmitters,
		f,
		onchainConfig,
		offchainConfigVersion,
		offchainConfig,
	)

	if err != nil {
		return fmt.Errorf("failed to set OCR2 config: %w", err)
	}
	return offRamp.client.ProcessTransaction(tx)
}

func (offRamp *OffRamp) SyncTokensAndPools(sourceTokens, pools []common.Address) error {
	opts, err := offRamp.client.TransactionOpts(offRamp.client.GetDefaultWallet())
	if err != nil {
		return fmt.Errorf("failed to get transaction opts: %w", err)
	}
	var tokenUpdates []evm_2_evm_offramp.InternalPoolUpdate
	for i, srcToken := range sourceTokens {
		tokenUpdates = append(tokenUpdates, evm_2_evm_offramp.InternalPoolUpdate{
			Token: srcToken,
			Pool:  pools[i],
		})
	}
	tx, err := offRamp.Instance.ApplyPoolUpdates(opts, []evm_2_evm_offramp.InternalPoolUpdate{}, tokenUpdates)
	if err != nil {
		return fmt.Errorf("failed to apply pool updates: %w", err)
	}
	log.Info().
		Interface("tokenUpdates", tokenUpdates).
		Str("offRamp", offRamp.Address()).
		Str(Network, offRamp.client.GetNetworkConfig().Name).
		Msg("tokenUpdates set in OffRamp")
	return offRamp.client.ProcessTransaction(tx)
}

type MockAggregator struct {
	client          blockchain.EVMClient
	Instance        *mock_v3_aggregator_contract.MockV3Aggregator
	ContractAddress common.Address
}

func (a *MockAggregator) ChainID() uint64 {
	return a.client.GetChainID().Uint64()
}

func (a *MockAggregator) UpdateRoundData(answer *big.Int) error {
	opts, err := a.client.TransactionOpts(a.client.GetDefaultWallet())
	if err != nil {
		return fmt.Errorf("unable to get transaction opts: %w", err)
	}
	log.Info().
		Str("Contract Address", a.ContractAddress.Hex()).
		Str("Network Name", a.client.GetNetworkConfig().Name).
		Msg("Updating Round Data")
	tx, err := a.Instance.UpdateRoundData(opts, big.NewInt(50), answer, big.NewInt(time.Now().UTC().UnixNano()), big.NewInt(time.Now().UTC().UnixNano()))
	if err != nil {
		return fmt.Errorf("unable to update round data: %w", err)
	}
	return a.client.ProcessTransaction(tx)
}

func (a *MockAggregator) WaitForTxConfirmations() error {
	return a.client.WaitForEvents()
}
