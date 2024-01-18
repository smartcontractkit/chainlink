package contracts

import (
	"context"
	"fmt"
	"math/big"
	"strconv"

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
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/price_registry"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/link_token_interface"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/shared/generated/erc20"
)

var HundredCoins = new(big.Int).Mul(big.NewInt(1e18), big.NewInt(100))

type RateLimiterConfig struct {
	Rate     *big.Int
	Capacity *big.Int
}

type ARMConfig struct {
	ARMWeightsByParticipants map[string]*big.Int // mapping : ARM participant address => weight
	ThresholdForBlessing     *big.Int
	ThresholdForBadSignal    *big.Int
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
		return nil, err
	}
	return balance, nil
}

func (token *ERC20Token) Approve(to string, amount *big.Int) error {
	opts, err := token.client.TransactionOpts(token.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	log.Info().
		Str("From", token.client.GetDefaultWallet().Address()).
		Str("To", to).
		Str("Token", token.Address()).
		Str("Amount", amount.String()).
		Uint64("Nonce", opts.Nonce.Uint64()).
		Str("Network Name", token.client.GetNetworkConfig().Name).
		Msg("Approving ERC20 Transfer")
	tx, err := token.instance.Approve(opts, common.HexToAddress(to), amount)
	if err != nil {
		return err
	}
	return token.client.ProcessTransaction(tx)
}

func (token *ERC20Token) Transfer(to string, amount *big.Int) error {
	opts, err := token.client.TransactionOpts(token.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	log.Info().
		Str("From", token.client.GetDefaultWallet().Address()).
		Str("To", to).
		Str("Amount", amount.String()).
		Uint64("Nonce", opts.Nonce.Uint64()).
		Str("Network Name", token.client.GetNetworkConfig().Name).
		Msg("Transferring ERC20")
	tx, err := token.instance.Transfer(opts, common.HexToAddress(to), amount)
	if err != nil {
		return err
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
		return nil, err
	}
	return balance, nil
}

func (l *LinkToken) Approve(to string, amount *big.Int) error {
	opts, err := l.client.TransactionOpts(l.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	log.Info().
		Str("From", l.client.GetDefaultWallet().Address()).
		Str("To", to).
		Str("Token", l.Address()).
		Str("Amount", amount.String()).
		Uint64("Nonce", opts.Nonce.Uint64()).
		Str("Network Name", l.client.GetNetworkConfig().Name).
		Msg("Approving LINK Transfer")
	tx, err := l.instance.Approve(opts, common.HexToAddress(to), amount)
	if err != nil {
		return err
	}
	return l.client.ProcessTransaction(tx)
}

func (l *LinkToken) Transfer(to string, amount *big.Int) error {
	opts, err := l.client.TransactionOpts(l.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	log.Info().
		Str("From", l.client.GetDefaultWallet().Address()).
		Str("To", to).
		Str("Amount", amount.String()).
		Uint64("Nonce", opts.Nonce.Uint64()).
		Str("Network Name", l.client.GetNetworkConfig().Name).
		Msg("Transferring LINK")
	tx, err := l.instance.Transfer(opts, common.HexToAddress(to), amount)
	if err != nil {
		return err
	}
	return l.client.ProcessTransaction(tx)
}

// LockReleaseTokenPool represents a LockReleaseTokenPool address
type LockReleaseTokenPool struct {
	client     blockchain.EVMClient
	Instance   *lock_release_token_pool.LockReleaseTokenPool
	EthAddress common.Address
}

func (pool *LockReleaseTokenPool) Address() string {
	return pool.EthAddress.Hex()
}

func (pool *LockReleaseTokenPool) RemoveLiquidity(amount *big.Int) error {
	opts, err := pool.client.TransactionOpts(pool.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	log.Info().
		Str("Token Pool", pool.Address()).
		Str("Amount", amount.String()).
		Msg("Initiating removing funds from pool")
	tx, err := pool.Instance.WithdrawLiquidity(opts, amount)
	if err != nil {
		return err
	}
	log.Info().
		Str("Token Pool", pool.Address()).
		Str("Amount", amount.String()).
		Str("Network Name", pool.client.GetNetworkConfig().Name).
		Msg("Liquidity removed")
	return pool.client.ProcessTransaction(tx)
}

type tokenApproveFn func(string, *big.Int) error

func (pool *LockReleaseTokenPool) AddLiquidity(approveFn tokenApproveFn, tokenAddr string, amount *big.Int) error {
	log.Info().
		Str("Link Token", tokenAddr).
		Str("Token Pool", pool.Address()).
		Msg("Initiating transferring of token to token pool")
	err := approveFn(pool.Address(), amount)
	if err != nil {
		return err
	}
	err = pool.client.WaitForEvents()
	if err != nil {
		return err
	}
	opts, err := pool.client.TransactionOpts(pool.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	_, err = pool.Instance.SetRebalancer(opts, opts.From)
	if err != nil {
		return err
	}
	opts, err = pool.client.TransactionOpts(pool.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	log.Info().
		Str("Token Pool", pool.Address()).
		Msg("Initiating adding Tokens in pool")
	tx, err := pool.Instance.ProvideLiquidity(opts, amount)
	if err != nil {
		return err
	}
	log.Info().
		Str("Token Pool", pool.Address()).
		Str("Link Token", tokenAddr).
		Str("Network Name", pool.client.GetNetworkConfig().Name).
		Msg("Liquidity added")
	return pool.client.ProcessTransaction(tx)
}

func (pool *LockReleaseTokenPool) SetOnRamp(onRamp common.Address) error {
	opts, err := pool.client.TransactionOpts(pool.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	log.Info().
		Str("Token Pool", pool.Address()).
		Msg("Setting on ramp for onramp router")
	tx, err := pool.Instance.ApplyRampUpdates(opts, []lock_release_token_pool.TokenPoolRampUpdate{
		{Ramp: onRamp, Allowed: true,
			RateLimiterConfig: lock_release_token_pool.RateLimiterConfig{
				IsEnabled: true,
				Capacity:  new(big.Int).Mul(big.NewInt(1e18), big.NewInt(1e9)),
				Rate:      new(big.Int).Mul(big.NewInt(1e18), big.NewInt(1e5)),
			}}}, []lock_release_token_pool.TokenPoolRampUpdate{})

	if err != nil {
		return err
	}
	log.Info().
		Str("Token Pool", pool.Address()).
		Str("OnRamp", onRamp.Hex()).
		Str("Network Name", pool.client.GetNetworkConfig().Name).
		Msg("OnRamp is set")
	return pool.client.ProcessTransaction(tx)
}

func (pool *LockReleaseTokenPool) SetOnRampRateLimit(onRamp common.Address, rl lock_release_token_pool.RateLimiterConfig) error {
	opts, err := pool.client.TransactionOpts(pool.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	log.Info().
		Str("Token Pool", pool.Address()).
		Str("OnRamp", onRamp.Hex()).
		Interface("RateLimiterConfig", rl).
		Msg("Setting Rate Limit on token pool")
	tx, err := pool.Instance.SetOnRampRateLimiterConfig(opts, onRamp, rl)

	if err != nil {
		return err
	}
	log.Info().
		Str("Token Pool", pool.Address()).
		Str("OnRamp", onRamp.Hex()).
		Interface("RateLimiterConfig", rl).
		Msg("Rate Limit on ramp is set")
	return pool.client.ProcessTransaction(tx)
}

func (pool *LockReleaseTokenPool) SetOffRampRateLimit(offRamp common.Address, rl lock_release_token_pool.RateLimiterConfig) error {
	opts, err := pool.client.TransactionOpts(pool.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	log.Info().
		Str("Token Pool", pool.Address()).
		Str("OffRamp", offRamp.Hex()).
		Interface("RateLimiterConfig", rl).
		Msg("Setting Rate Limit offramp")
	tx, err := pool.Instance.SetOffRampRateLimiterConfig(opts, offRamp, rl)

	if err != nil {
		return err
	}
	log.Info().
		Str("Token Pool", pool.Address()).
		Str("OffRamp", offRamp.Hex()).
		Interface("RateLimiterConfig", rl).
		Msg("Rate Limit offRamp is set")
	return pool.client.ProcessTransaction(tx)
}

func (pool *LockReleaseTokenPool) SetOffRamp(offRamp common.Address) error {
	opts, err := pool.client.TransactionOpts(pool.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	log.Info().
		Str("Token Pool", pool.Address()).
		Msg("Setting off ramp for Token Pool")

	tx, err := pool.Instance.ApplyRampUpdates(opts, []lock_release_token_pool.TokenPoolRampUpdate{}, []lock_release_token_pool.TokenPoolRampUpdate{
		{Ramp: offRamp, Allowed: true, RateLimiterConfig: lock_release_token_pool.RateLimiterConfig{
			IsEnabled: true,
			Capacity:  new(big.Int).Mul(big.NewInt(1e18), big.NewInt(1e9)),
			Rate:      new(big.Int).Mul(big.NewInt(1e18), big.NewInt(1e5)),
		}}})
	if err != nil {
		return err
	}
	log.Info().
		Str("Token Pool", pool.Address()).
		Str("OffRamp", offRamp.Hex()).
		Str("Network Name", pool.client.GetNetworkConfig().Name).
		Msg("OffRamp is set")
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
		return err
	}

	log.Info().
		Interface("signerAddresses", signers).
		Interface("transmitterAddresses", transmitters).
		Str("Network Name", b.client.GetNetworkConfig().Name).
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
		return err
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
		return err
	}
	tx, err := rDapp.instance.SetRevert(opts, revert)
	if err != nil {
		return err
	}
	log.Info().
		Bool("revert", revert).
		Str("tx", tx.Hash().String()).
		Str("ReceiverDapp", rDapp.Address()).
		Str("Network Name", rDapp.client.GetNetworkConfig().Name).
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
		return err
	}
	tx, err := c.Instance.ApplyPriceUpdatersUpdates(opts, []common.Address{addr}, []common.Address{})
	if err != nil {
		return err
	}
	log.Info().
		Str("updaters", addr.Hex()).
		Str("Network Name", c.client.GetNetworkConfig().Name).
		Msg("PriceRegistry updater added")
	return c.client.ProcessTransaction(tx)
}

func (c *PriceRegistry) AddFeeToken(addr common.Address) error {
	opts, err := c.client.TransactionOpts(c.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := c.Instance.ApplyFeeTokensUpdates(opts, []common.Address{addr}, []common.Address{})
	if err != nil {
		return err
	}
	log.Info().
		Str("feeTokens", addr.Hex()).
		Str("Network Name", c.client.GetNetworkConfig().Name).
		Msg("PriceRegistry feeToken set")
	return c.client.ProcessTransaction(tx)
}

func (c *PriceRegistry) UpdatePrices(priceUpdates price_registry.InternalPriceUpdates) error {
	opts, err := c.client.TransactionOpts(c.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := c.Instance.UpdatePrices(opts, priceUpdates)
	if err != nil {
		return err
	}
	log.Info().
		Str("Network Name", c.client.GetNetworkConfig().Name).
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
		return err
	}
	log.Info().
		Str("Router", r.Address()).
		Str("ChainSelector", strconv.FormatUint(chainSelector, 10)).
		Msg("Setting on ramp for r")

	tx, err := r.Instance.ApplyRampUpdates(opts, []router.RouterOnRamp{{DestChainSelector: chainSelector, OnRamp: onRamp}}, nil, nil)
	if err != nil {
		return err
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
		return nil, err
	}
	if valueForNative != nil {
		opts.Value = valueForNative
	}

	log.Info().
		Str("Network", r.client.GetNetworkName()).
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
		return nil, err
	}
	log.Info().
		Str("router", r.Address()).
		Str("txHash", tx.Hash().Hex()).
		Str("Network Name", r.client.GetNetworkConfig().Name).
		Str("chain selector", strconv.FormatUint(destChainSelector, 10)).
		Msg("msg is sent")
	return tx, r.client.ProcessTransaction(tx)
}

func (r *Router) AddOffRamp(offRamp common.Address, sourceChainId uint64) (*types.Transaction, error) {
	opts, err := r.client.TransactionOpts(r.client.GetDefaultWallet())
	if err != nil {
		return nil, err
	}
	tx, err := r.Instance.ApplyRampUpdates(opts, nil, nil, []router.RouterOffRamp{{SourceChainSelector: sourceChainId, OffRamp: offRamp}})
	if err != nil {
		return nil, err
	}
	log.Info().
		Str("offRamp", offRamp.Hex()).
		Str("Network Name", r.client.GetNetworkConfig().Name).
		Msg("offRamp is added to Router")
	return tx, r.client.ProcessTransaction(tx)
}

func (r *Router) SetWrappedNative(wNative common.Address) (*types.Transaction, error) {
	opts, err := r.client.TransactionOpts(r.client.GetDefaultWallet())
	if err != nil {
		return nil, err
	}
	tx, err := r.Instance.SetWrappedNative(opts, wNative)
	if err != nil {
		return nil, err
	}
	log.Info().
		Str("wrapped native", wNative.Hex()).
		Str("router", r.Address()).
		Str("Network Name", r.client.GetNetworkConfig().Name).
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
		return err
	}
	owner := common.HexToAddress(onRamp.client.GetDefaultWallet().Address())
	// set the payee to the default wallet
	tx, err := onRamp.Instance.SetNops(opts, []evm_2_evm_onramp.EVM2EVMOnRampNopAndWeight{{
		Nop:    owner,
		Weight: 1,
	}})
	if err != nil {
		return err
	}
	return onRamp.client.ProcessTransaction(tx)
}

func (onRamp *OnRamp) SetTokenTransferFeeConfig(tokenTransferFeeConfig []evm_2_evm_onramp.EVM2EVMOnRampTokenTransferFeeConfigArgs) error {
	opts, err := onRamp.client.TransactionOpts(onRamp.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := onRamp.Instance.SetTokenTransferFeeConfig(opts, tokenTransferFeeConfig)
	if err != nil {
		return err
	}
	log.Info().
		Interface("tokenTransferFeeConfig", tokenTransferFeeConfig).
		Str("onRamp", onRamp.Address()).
		Str("Network Name", onRamp.client.GetNetworkConfig().Name).
		Msg("TokenTransferFeeConfig set in OnRamp")
	return onRamp.client.ProcessTransaction(tx)
}

func (onRamp *OnRamp) ApplyPoolUpdates(poolUpdates []evm_2_evm_onramp.InternalPoolUpdate) error {
	opts, err := onRamp.client.TransactionOpts(onRamp.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := onRamp.Instance.ApplyPoolUpdates(opts, []evm_2_evm_onramp.InternalPoolUpdate{}, poolUpdates)
	if err != nil {
		return err
	}
	log.Info().
		Interface("poolUpdates", poolUpdates).
		Str("onRamp", onRamp.Address()).
		Str("Network Name", onRamp.client.GetNetworkConfig().Name).
		Msg("poolUpdates set in OnRamp")
	return onRamp.client.ProcessTransaction(tx)
}

func (onRamp *OnRamp) PayNops() error {
	opts, err := onRamp.client.TransactionOpts(onRamp.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := onRamp.Instance.PayNops(opts)
	if err != nil {
		return err
	}
	return onRamp.client.ProcessTransaction(tx)
}

func (onRamp *OnRamp) WithdrawNonLinkFees(wrappedNative common.Address) error {
	opts, err := onRamp.client.TransactionOpts(onRamp.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	owner := common.HexToAddress(onRamp.client.GetDefaultWallet().Address())
	tx, err := onRamp.Instance.WithdrawNonLinkFees(opts, wrappedNative, owner)
	if err != nil {
		return err
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
		return err
	}
	log.Info().
		Bool("Enabled", rlConfig.IsEnabled).
		Str("capacity", rlConfig.Capacity.String()).
		Str("rate", rlConfig.Rate.String()).
		Str("onRamp", onRamp.Address()).
		Str("Network Name", onRamp.client.GetNetworkConfig().Name).
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
		return err
	}
	log.Info().
		Interface("signerAddresses", signers).
		Interface("transmitterAddresses", transmitters).
		Str("Network Name", offRamp.client.GetNetworkConfig().Name).
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
		return err
	}
	return offRamp.client.ProcessTransaction(tx)
}

func (offRamp *OffRamp) SyncTokensAndPools(sourceTokens, pools []common.Address) error {
	opts, err := offRamp.client.TransactionOpts(offRamp.client.GetDefaultWallet())
	if err != nil {
		return err
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
		return err
	}
	log.Info().
		Interface("tokenUpdates", tokenUpdates).
		Str("offRamp", offRamp.Address()).
		Str("Network Name", offRamp.client.GetNetworkConfig().Name).
		Msg("tokenUpdates set in OffRamp")
	return offRamp.client.ProcessTransaction(tx)
}
