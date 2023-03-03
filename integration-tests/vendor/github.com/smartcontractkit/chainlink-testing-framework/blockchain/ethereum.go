package blockchain

// Contans implementations for multi and single node ethereum clients
import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"math"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink-env/environment"
	"golang.org/x/sync/errgroup"

	"github.com/smartcontractkit/chainlink-testing-framework/utils"
)

// EthereumClient wraps the client and the BlockChain network to interact with an EVM based Blockchain
type EthereumClient struct {
	ID                  int
	Client              *ethclient.Client
	NetworkConfig       EVMNetwork
	Wallets             []*EthereumWallet
	DefaultWallet       *EthereumWallet
	NonceSettings       *NonceSettings
	headerSubscriptions map[string]HeaderEventSubscription
	subscriptionMutex   *sync.Mutex
	queueTransactions   bool
	gasStats            *GasStats
	doneChan            chan struct{}
}

// newEVMClient creates an EVM client for a single node/URL
func newEVMClient(networkSettings EVMNetwork) (EVMClient, error) {
	log.Info().
		Str("Name", networkSettings.Name).
		Str("URL", networkSettings.URL).
		Interface("Settings", networkSettings).
		Msg("Connecting client")
	cl, err := ethclient.Dial(networkSettings.URL)
	if err != nil {
		return nil, err
	}

	ec := &EthereumClient{
		NetworkConfig:       networkSettings,
		Client:              cl,
		Wallets:             make([]*EthereumWallet, 0),
		headerSubscriptions: map[string]HeaderEventSubscription{},
		subscriptionMutex:   &sync.Mutex{},
		queueTransactions:   false,
		doneChan:            make(chan struct{}),
	}

	if ec.NetworkConfig.Simulated {
		ec.NonceSettings = newNonceSettings()
	} else { // un-simulated chain means potentially running tests in parallel, need to share nonces
		ec.NonceSettings = useGlobalNonceManager(ec.GetChainID())
	}

	if err := ec.LoadWallets(networkSettings); err != nil {
		return nil, err
	}
	ec.gasStats = NewGasStats(ec.ID)
	go ec.newHeadersLoop()

	return wrapSingleClient(networkSettings, ec), nil
}

// SyncNonce sets the NonceMu and Nonces value based on existing EVMClient
// it ensures the instance of EthereumClient is synced with passed EVMClient's nonce updates.
func (e *EthereumClient) SyncNonce(c EVMClient) {
	n := c.GetNonceSetting()
	e.NonceSettings.NonceMu = n.NonceMu
	e.NonceSettings.Nonces = n.Nonces
}

// newHeadersLoop Logs when new headers come in
func (e *EthereumClient) newHeadersLoop() {
	for {
		if err := e.subscribeToNewHeaders(); err != nil {
			log.Error().
				Err(err).
				Str("NetworkName", e.NetworkConfig.Name).
				Msg("Error while subscribing to headers")
			time.Sleep(time.Second)
			continue
		}
		break
	}
	log.Debug().Str("NetworkName", e.NetworkConfig.Name).Msg("Stopped subscribing to new headers")
}

// Get returns the underlying client type to be used generically across the framework for switching
// network types
func (e *EthereumClient) Get() interface{} {
	return e
}

// GetNetworkName retrieves the ID of the network that the client interacts with
func (e *EthereumClient) GetNetworkName() string {
	return e.NetworkConfig.Name
}

// NetworkSimulated returns true if the network is a simulated geth instance, false otherwise
func (e *EthereumClient) NetworkSimulated() bool {
	return e.NetworkConfig.Simulated
}

func (e *EthereumClient) GetNonceSetting() NonceSettings {
	e.NonceSettings.NonceMu.Lock()
	defer e.NonceSettings.NonceMu.Unlock()
	return NonceSettings{
		NonceMu: e.NonceSettings.NonceMu,
		Nonces:  e.NonceSettings.Nonces,
	}
}

// GetChainID retrieves the ChainID of the network that the client interacts with
func (e *EthereumClient) GetChainID() *big.Int {
	return big.NewInt(e.NetworkConfig.ChainID)
}

// GetClients not used, only applicable to EthereumMultinodeClient
func (e *EthereumClient) GetClients() []EVMClient {
	return []EVMClient{e}
}

// DefaultWallet returns the default wallet for the network
func (e *EthereumClient) GetDefaultWallet() *EthereumWallet {
	return e.DefaultWallet
}

// DefaultWallet returns the default wallet for the network
func (e *EthereumClient) GetWallets() []*EthereumWallet {
	return e.Wallets
}

// DefaultWallet returns the default wallet for the network
func (e *EthereumClient) GetNetworkConfig() *EVMNetwork {
	return &e.NetworkConfig
}

// SetID sets client id, only used for multi-node networks
func (e *EthereumClient) SetID(id int) {
	e.ID = id
}

// SetDefaultWallet sets default wallet
func (e *EthereumClient) SetDefaultWallet(num int) error {
	if num >= len(e.GetWallets()) {
		return fmt.Errorf("no wallet #%d found for default client", num)
	}
	e.DefaultWallet = e.Wallets[num]
	return nil
}

// SetWallets sets all wallets to be used by the client
func (e *EthereumClient) SetWallets(wallets []*EthereumWallet) {
	e.Wallets = wallets
}

// LoadWallets loads wallets from config
func (e *EthereumClient) LoadWallets(cfg EVMNetwork) error {
	pkStrings := cfg.PrivateKeys
	for _, pks := range pkStrings {
		w, err := NewEthereumWallet(pks)
		if err != nil {
			return err
		}
		e.Wallets = append(e.Wallets, w)
	}
	if len(e.Wallets) == 0 {
		return fmt.Errorf("no private keys found to load wallets")
	}
	e.DefaultWallet = e.Wallets[0]
	return nil
}

// BalanceAt returns the ETH balance of the specified address
func (e *EthereumClient) BalanceAt(ctx context.Context, address common.Address) (*big.Int, error) {
	return e.Client.BalanceAt(ctx, address, nil)
}

// SwitchNode not used, only applicable to EthereumMultinodeClient
func (e *EthereumClient) SwitchNode(_ int) error {
	return nil
}

// HeaderHashByNumber gets header hash by block number
func (e *EthereumClient) HeaderHashByNumber(ctx context.Context, bn *big.Int) (string, error) {
	h, err := e.Client.HeaderByNumber(ctx, bn)
	if err != nil {
		return "", err
	}
	return h.Hash().String(), nil
}

// HeaderTimestampByNumber gets header timestamp by number
func (e *EthereumClient) HeaderTimestampByNumber(ctx context.Context, bn *big.Int) (uint64, error) {
	h, err := e.Client.HeaderByNumber(ctx, bn)
	if err != nil {
		return 0, err
	}
	return h.Time, nil
}

// BlockNumber gets latest block number
func (e *EthereumClient) LatestBlockNumber(ctx context.Context) (uint64, error) {
	bn, err := e.Client.BlockNumber(ctx)
	if err != nil {
		return 0, err
	}
	return bn, nil
}

func (e *EthereumClient) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	// for instant txs, wait for permission to go
	if e.NetworkConfig.MinimumConfirmations <= 0 {
		fromAddr, err := types.Sender(types.LatestSignerForChainID(tx.ChainId()), tx)
		if err != nil {
			return err
		}
		<-e.NonceSettings.registerInstantTransaction(fromAddr.Hex(), tx.Nonce())
	}
	return e.Client.SendTransaction(ctx, tx)
}

// Fund sends some ETH to an address using the default wallet
func (e *EthereumClient) Fund(
	toAddress string,
	amount *big.Float,
) error {
	privateKey, err := crypto.HexToECDSA(e.DefaultWallet.PrivateKey())
	if err != nil {
		return fmt.Errorf("invalid private key: %v", err)
	}
	to := common.HexToAddress(toAddress)

	suggestedGasTipCap, err := e.Client.SuggestGasTipCap(context.Background())
	if err != nil {
		return err
	}

	// Bump Tip Cap
	gasPriceBuffer := big.NewInt(0).SetUint64(e.NetworkConfig.GasEstimationBuffer)
	suggestedGasTipCap.Add(suggestedGasTipCap, gasPriceBuffer)

	nonce, err := e.GetNonce(context.Background(), common.HexToAddress(e.DefaultWallet.Address()))
	if err != nil {
		return err
	}
	latestHeader, err := e.Client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return err
	}
	baseFeeMult := big.NewInt(1).Mul(latestHeader.BaseFee, big.NewInt(2))
	gasFeeCap := baseFeeMult.Add(baseFeeMult, suggestedGasTipCap)

	estimatedGas, err := e.Client.EstimateGas(context.Background(), ethereum.CallMsg{})
	if err != nil {
		return err
	}

	tx, err := types.SignNewTx(privateKey, types.LatestSignerForChainID(e.GetChainID()), &types.DynamicFeeTx{
		ChainID:   e.GetChainID(),
		Nonce:     nonce,
		To:        &to,
		Value:     utils.EtherToWei(amount),
		GasTipCap: suggestedGasTipCap,
		GasFeeCap: gasFeeCap,
		Gas:       estimatedGas,
	})
	if err != nil {
		return err
	}

	log.Info().
		Str("Token", "ETH").
		Str("From", e.DefaultWallet.Address()).
		Str("To", toAddress).
		Str("Hash", tx.Hash().Hex()).
		Uint64("Nonce", tx.Nonce()).
		Str("Network Name", e.GetNetworkName()).
		Str("Amount", amount.String()).
		Uint64("Estimated Gas Cost", new(big.Int).Mul(gasFeeCap, new(big.Int).SetUint64(estimatedGas)).Uint64()).
		Msg("Funding Address")
	if err := e.SendTransaction(context.Background(), tx); err != nil {
		if strings.Contains(err.Error(), "nonce") {
			err = errors.Wrap(err, fmt.Sprintf("using nonce %d", nonce))
		}
		return err
	}

	return e.ProcessTransaction(tx)
}

func (e *EthereumClient) ReturnFunds(fromKey *ecdsa.PrivateKey) error {
	var tx *types.Transaction
	var err error
	for attempt := 1; attempt < 10; attempt++ {
		tx, err = attemptReturn(e, fromKey, attempt)
		if err == nil {
			return e.ProcessTransaction(tx)
		}
		log.Debug().Err(err).Int("Attempt", attempt+1).Msg("Error returning funds from Chainlink node, trying again")
	}
	return err
}

// a single fund return attempt, further attempts exponentially raise the error margin for fund returns
func attemptReturn(e *EthereumClient, fromKey *ecdsa.PrivateKey, attemptCount int) (*types.Transaction, error) {
	to := common.HexToAddress(e.DefaultWallet.Address())

	suggestedGasTipCap, err := e.Client.SuggestGasTipCap(context.Background())
	if err != nil {
		return nil, err
	}
	latestHeader, err := e.Client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	baseFeeMult := big.NewInt(1).Mul(latestHeader.BaseFee, big.NewInt(2))
	gasFeeCap := baseFeeMult.Add(baseFeeMult, suggestedGasTipCap)
	gasFeeCap.Add(gasFeeCap, big.NewInt(int64(math.Pow(float64(attemptCount), 2)*1000))) // exponentially increase error margin

	fromAddress, err := utils.PrivateKeyToAddress(fromKey)
	if err != nil {
		return nil, err
	}
	balance, err := e.Client.BalanceAt(context.Background(), fromAddress, nil)
	if err != nil {
		return nil, err
	}
	estGas, err := e.Client.EstimateGas(context.Background(), ethereum.CallMsg{})
	if err != nil {
		return nil, err
	}
	balance.Sub(balance, big.NewInt(1).Mul(gasFeeCap, big.NewInt(0).SetUint64(estGas)))

	nonce, err := e.GetNonce(context.Background(), fromAddress)
	if err != nil {
		return nil, err
	}
	estimatedGas, err := e.Client.EstimateGas(context.Background(), ethereum.CallMsg{})
	if err != nil {
		return nil, err
	}

	tx, err := types.SignNewTx(fromKey, types.LatestSignerForChainID(e.GetChainID()), &types.DynamicFeeTx{
		ChainID:   e.GetChainID(),
		Nonce:     nonce,
		To:        &to,
		Value:     balance,
		GasTipCap: suggestedGasTipCap,
		GasFeeCap: gasFeeCap,
		Gas:       estimatedGas,
	})
	if err != nil {
		return nil, err
	}
	log.Info().
		Str("Token", "ETH").
		Str("Amount", balance.String()).
		Str("From", fromAddress.Hex()).
		Msg("Returning Funds to Default Wallet")
	return tx, e.SendTransaction(context.Background(), tx)
}

// EstimateCostForChainlinkOperations calculates required amount of ETH for amountOfOperations Chainlink operations
// based on the network's suggested gas price and the chainlink gas limit. This is fairly imperfect and should be used
// as only a rough, upper-end estimate instead of an exact calculation.
// See https://ethereum.org/en/developers/docs/gas/#post-london for info on how gas calculation works
func (e *EthereumClient) EstimateCostForChainlinkOperations(amountOfOperations int) (*big.Float, error) {
	bigAmountOfOperations := big.NewInt(int64(amountOfOperations))
	gasPriceInWei, err := e.Client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}

	// https://ethereum.stackexchange.com/questions/19665/how-to-calculate-transaction-fee
	// total gas limit = chainlink gas limit + gas limit buffer
	gasLimit := e.NetworkConfig.GasEstimationBuffer + e.NetworkConfig.ChainlinkTransactionLimit
	// gas cost for TX = total gas limit * estimated gas price
	gasCostPerOperationWei := big.NewInt(1).Mul(big.NewInt(1).SetUint64(gasLimit), gasPriceInWei)
	gasCostPerOperationETH := utils.WeiToEther(gasCostPerOperationWei)
	// total Wei needed for all TXs = total value for TX * number of TXs
	totalWeiForAllOperations := big.NewInt(1).Mul(gasCostPerOperationWei, bigAmountOfOperations)
	totalEthForAllOperations := utils.WeiToEther(totalWeiForAllOperations)

	log.Debug().
		Int("Number of Operations", amountOfOperations).
		Uint64("Gas Limit per Operation", gasLimit).
		Str("Value per Operation (ETH)", gasCostPerOperationETH.String()).
		Str("Total (ETH)", totalEthForAllOperations.String()).
		Msg("Calculated ETH for Chainlink Operations")

	return totalEthForAllOperations, nil
}

// DeployContract acts as a general contract deployment tool to an ethereum chain
func (e *EthereumClient) DeployContract(
	contractName string,
	deployer ContractDeployer,
) (*common.Address, *types.Transaction, interface{}, error) {
	suggestedTipCap, err := e.Client.SuggestGasTipCap(context.Background())
	if err != nil {
		return nil, nil, nil, err
	}
	gasPriceBuffer := big.NewInt(0).SetUint64(e.NetworkConfig.GasEstimationBuffer)

	opts, err := e.TransactionOpts(e.DefaultWallet)
	if err != nil {
		return nil, nil, nil, err
	}
	opts.GasTipCap = suggestedTipCap.Add(gasPriceBuffer, suggestedTipCap)

	if e.NetworkConfig.GasEstimationBuffer > 0 {
		log.Debug().
			Uint64("Suggested Gas Tip Cap", big.NewInt(0).Sub(suggestedTipCap, gasPriceBuffer).Uint64()).
			Uint64("Bumped Gas Tip Cap", suggestedTipCap.Uint64()).
			Str("Contract Name", contractName).
			Msg("Bumping Suggested Gas Price")
	}

	contractAddress, transaction, contractInstance, err := deployer(opts, e.Client)
	if err != nil {
		if strings.Contains(err.Error(), "nonce") {
			err = errors.Wrap(err, fmt.Sprintf("using nonce %d", opts.Nonce.Uint64()))
		}
		return nil, nil, nil, err
	}

	if err = e.ProcessTransaction(transaction); err != nil {
		return nil, nil, nil, err
	}

	log.Info().
		Str("Contract Address", contractAddress.Hex()).
		Str("Contract Name", contractName).
		Str("From", e.DefaultWallet.Address()).
		Str("Total Gas Cost (ETH)", utils.WeiToEther(transaction.Cost()).String()).
		Str("Network Name", e.NetworkConfig.Name).
		Msg("Deployed contract")
	return &contractAddress, transaction, contractInstance, err
}

// LoadContract load already deployed contract instance
func (e *EthereumClient) LoadContract(contractName string, contractAddress common.Address, loader ContractLoader) (interface{}, error) {
	contractInstance, err := loader(contractAddress, e.Client)
	if err != nil {
		return nil, err
	}
	log.Info().
		Str("Contract Address", contractAddress.Hex()).
		Str("Contract Name", contractName).
		Str("Network Name", e.NetworkConfig.Name).
		Msg("Loaded contract instance")
	return contractInstance, err
}

// TransactionOpts returns the base Tx options for 'transactions' that interact with a smart contract. Since most
// contract interactions in this framework are designed to happen through abigen calls, it's intentionally quite bare.
func (e *EthereumClient) TransactionOpts(from *EthereumWallet) (*bind.TransactOpts, error) {
	privateKey, err := crypto.HexToECDSA(from.PrivateKey())
	if err != nil {
		return nil, fmt.Errorf("invalid private key: %v", err)
	}
	opts, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(e.NetworkConfig.ChainID))
	if err != nil {
		return nil, err
	}
	opts.From = common.HexToAddress(from.Address())
	opts.Context = context.Background()

	nonce, err := e.GetNonce(context.Background(), common.HexToAddress(from.Address()))
	if err != nil {
		return nil, err
	}
	opts.Nonce = big.NewInt(int64(nonce))

	if e.NetworkConfig.MinimumConfirmations <= 0 { // Wait for your turn to send on an L2 chain
		<-e.NonceSettings.registerInstantTransaction(from.Address(), nonce)
	}
	return opts, nil
}

// ProcessTransaction will queue or wait on a transaction depending on whether parallel transactions are enabled
func (e *EthereumClient) ProcessTransaction(tx *types.Transaction) error {
	var txConfirmer HeaderEventSubscription
	if e.GetNetworkConfig().MinimumConfirmations <= 0 {
		fromAddr, err := types.Sender(types.LatestSignerForChainID(tx.ChainId()), tx)
		if err != nil {
			return err
		}
		e.NonceSettings.sentInstantTransaction(fromAddr.Hex()) // On an L2 chain, indicate the tx has been sent
		txConfirmer = NewInstantConfirmer(e, tx.Hash(), nil, nil)
	} else {
		txConfirmer = NewTransactionConfirmer(e, tx, e.GetNetworkConfig().MinimumConfirmations)
	}

	e.AddHeaderEventSubscription(tx.Hash().String(), txConfirmer)

	if !e.queueTransactions { // For sequential transactions
		log.Debug().Str("Hash", tx.Hash().String()).Msg("Waiting for TX to confirm before moving on")
		defer e.DeleteHeaderEventSubscription(tx.Hash().String())
		return txConfirmer.Wait()
	}
	return nil
}

// ProcessTransaction will queue or wait on a transaction depending on whether parallel transactions are enabled
func (e *EthereumClient) ProcessEvent(name string, event *types.Log, confirmedChan chan bool, errorChan chan error) error {
	var eventConfirmer HeaderEventSubscription
	if e.GetNetworkConfig().MinimumConfirmations <= 0 {
		eventConfirmer = NewInstantConfirmer(e, event.TxHash, confirmedChan, errorChan)
	} else {
		eventConfirmer = NewEventConfirmer(name, e, event, e.GetNetworkConfig().MinimumConfirmations, confirmedChan, errorChan)
	}

	subscriptionHash := fmt.Sprintf("%s-%s", event.TxHash.Hex(), name) // Many events can occupy the same tx hash
	e.AddHeaderEventSubscription(subscriptionHash, eventConfirmer)

	if !e.queueTransactions { // For sequential transactions
		log.Debug().Str("Hash", event.Address.Hex()).Msg("Waiting for Event to confirm before moving on")
		defer e.DeleteHeaderEventSubscription(subscriptionHash)
		return eventConfirmer.Wait()
	}
	return nil
}

// IsTxConfirmed checks if the transaction is confirmed on chain or not
func (e *EthereumClient) IsTxConfirmed(txHash common.Hash) (bool, error) {
	tx, isPending, err := e.Client.TransactionByHash(context.Background(), txHash)
	if err != nil {
		return !isPending, err
	}
	if !isPending {
		receipt, err := e.Client.TransactionReceipt(context.Background(), txHash)
		if err != nil {
			return !isPending, err
		}
		e.gasStats.AddClientTXData(TXGasData{
			TXHash:            txHash.String(),
			Value:             tx.Value().Uint64(),
			GasLimit:          tx.Gas(),
			GasUsed:           receipt.GasUsed,
			GasPrice:          tx.GasPrice().Uint64(),
			CumulativeGasUsed: receipt.CumulativeGasUsed,
		})
		if receipt.Status == 0 { // 0 indicates failure, 1 indicates success
			reason, err := e.errorReason(e.Client, tx, receipt)
			if err != nil {
				log.Warn().Str("TX Hash", txHash.Hex()).
					Str("To", tx.To().Hex()).
					Uint64("Nonce", tx.Nonce()).
					Msg("Transaction failed and was reverted! Unable to retrieve reason!")
				return false, err
			}
			log.Warn().Str("TX Hash", txHash.Hex()).
				Str("To", tx.To().Hex()).
				Str("Revert reason", reason).
				Msg("Transaction failed and was reverted!")
		}
	}
	return !isPending, err
}

// IsEventConfirmed returns if eth client can confirm that the event happened
func (e *EthereumClient) IsEventConfirmed(event *types.Log) (confirmed, removed bool, err error) {
	if event.Removed {
		return false, event.Removed, nil
	}
	eventTx, isPending, err := e.Client.TransactionByHash(context.Background(), event.TxHash)
	if err != nil {
		return false, event.Removed, err
	}
	if isPending {
		return false, event.Removed, nil
	}
	eventReceipt, err := e.Client.TransactionReceipt(context.Background(), eventTx.Hash())
	if err != nil {
		return false, event.Removed, err
	}
	if eventReceipt.Status == 0 { // Failed event tx
		reason, err := e.errorReason(e.Client, eventTx, eventReceipt)
		if err != nil {
			log.Warn().Str("TX Hash", eventTx.Hash().Hex()).Msg("Transaction failed and was reverted! Unable to retrieve reason!")
			return false, event.Removed, err
		}
		log.Warn().Str("TX Hash", eventTx.Hash().Hex()).
			Str("Revert reason", reason).
			Msg("Transaction failed and was reverted!")
		return false, event.Removed, err
	}
	headerByNumber, err := e.Client.HeaderByNumber(context.Background(), big.NewInt(0).SetUint64(event.BlockNumber))
	if err != nil || headerByNumber == nil {
		return false, event.Removed, err
	}
	if headerByNumber.Hash() != event.BlockHash {
		return false, event.Removed, nil
	}

	return true, event.Removed, nil
}

// GetTxReceipt returns the receipt of the transaction if available, error otherwise
func (e *EthereumClient) GetTxReceipt(txHash common.Hash) (*types.Receipt, error) {
	receipt, err := e.Client.TransactionReceipt(context.Background(), txHash)
	if err != nil {
		return nil, err
	}
	return receipt, nil
}

// ParallelTransactions when enabled, sends the transaction without waiting for transaction confirmations. The hashes
// are then stored within the client and confirmations can be waited on by calling WaitForEvents.
func (e *EthereumClient) ParallelTransactions(enabled bool) {
	e.queueTransactions = enabled
}

// Close tears down the current open Ethereum client
func (e *EthereumClient) Close() error {
	// close(e.NonceSettings.doneChan)
	close(e.doneChan)
	return nil
}

// EstimateTransactionGasCost estimates the current total gas cost for a simple transaction
func (e *EthereumClient) EstimateTransactionGasCost() (*big.Int, error) {
	gasPrice, err := e.Client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}
	return gasPrice.Mul(gasPrice, big.NewInt(21000)), err
}

// GasStats retrieves all information on gas spent by this client
func (e *EthereumClient) GasStats() *GasStats {
	return e.gasStats
}

func (e *EthereumClient) Backend() bind.ContractBackend {
	return e.Client
}

// AddHeaderEventSubscription adds a new header subscriber within the client to receive new headers
func (e *EthereumClient) AddHeaderEventSubscription(key string, subscriber HeaderEventSubscription) {
	e.subscriptionMutex.Lock()
	defer e.subscriptionMutex.Unlock()
	e.headerSubscriptions[key] = subscriber
}

// DeleteHeaderEventSubscription removes a header subscriber from the map
func (e *EthereumClient) DeleteHeaderEventSubscription(key string) {
	e.subscriptionMutex.Lock()
	defer e.subscriptionMutex.Unlock()
	delete(e.headerSubscriptions, key)
}

// WaitForEvents is a blocking function that waits for all event subscriptions that have been queued within the client.
func (e *EthereumClient) WaitForEvents() error {
	log.Debug().Msg("Waiting for blockchain events to finish before continuing")
	queuedEvents := e.GetHeaderSubscriptions()
	g := errgroup.Group{}

	for subName, sub := range queuedEvents {
		subName := subName
		sub := sub
		g.Go(func() error {
			defer e.DeleteHeaderEventSubscription(subName)
			return sub.Wait()
		})
	}
	return g.Wait()
}

// SubscribeFilterLogs subscribes to the results of a streaming filter query.
func (e *EthereumClient) SubscribeFilterLogs(ctx context.Context, q ethereum.FilterQuery, ch chan<- types.Log) (ethereum.Subscription, error) {
	return e.Client.SubscribeFilterLogs(ctx, q, ch)
}

// EthereumMultinodeClient wraps the client and the BlockChain network to interact with an EVM based Blockchain with multiple nodes
type EthereumMultinodeClient struct {
	DefaultClient EVMClient
	Clients       []EVMClient
}

func (e *EthereumMultinodeClient) Backend() bind.ContractBackend {
	return e.DefaultClient.Backend()
}

// LoadContract load already deployed contract instance
func (e *EthereumMultinodeClient) LoadContract(contractName string, address common.Address, loader ContractLoader) (interface{}, error) {
	return e.DefaultClient.LoadContract(contractName, address, loader)
}

// EstimateCostForChainlinkOperations calculates TXs cost as a dirty estimation based on transactionLimit for that network
func (e *EthereumMultinodeClient) EstimateCostForChainlinkOperations(amountOfOperations int) (*big.Float, error) {
	return e.DefaultClient.EstimateCostForChainlinkOperations(amountOfOperations)
}

// NewEVMClient returns a multi-node EVM client connected to the specified network
func NewEVMClient(networkSettings EVMNetwork, env *environment.Environment) (EVMClient, error) {
	ecl := &EthereumMultinodeClient{}
	if networkSettings.Simulated {
		if _, ok := env.URLs[networkSettings.Name]; !ok {
			return nil, fmt.Errorf("network %s not found in environment", networkSettings.Name)
		}
		if env == nil {
			log.Warn().Str("Network", networkSettings.Name).Msg("No test environment deployed")
		} else {
			networkSettings.URLs = env.URLs[networkSettings.Name]
		}
	} else {
		if len(networkSettings.URLs) == 0 {
			return nil, fmt.Errorf("no URL is provided to connect to network %s", networkSettings.Name)
		}
	}

	for idx, networkURL := range networkSettings.URLs {
		networkSettings.URL = networkURL
		ec, err := newEVMClient(networkSettings)

		if err != nil {
			return nil, err
		}
		ec.SetID(idx)
		ecl.Clients = append(ecl.Clients, ec)
	}
	ecl.DefaultClient = ecl.Clients[0]
	wrappedClient := wrapMultiClient(networkSettings, ecl)
	// required in Geth when you need to call "simulate" transactions from nodes
	if ecl.NetworkSimulated() {
		if err := ecl.Fund("0x0", big.NewFloat(1000)); err != nil {
			return nil, err
		}
	}

	return wrappedClient, nil
}

// ConcurrentEVMClient returns a multi-node EVM client connected to the specified network
// It is used for concurrent interactions from different threads with the same network and from same owner
// account. This ensures that correct nonce value is fetched when an instance of EVMClient is initiated using this method.
func ConcurrentEVMClient(networkSettings EVMNetwork, env *environment.Environment, existing EVMClient) (EVMClient, error) {
	ecl := &EthereumMultinodeClient{}
	if networkSettings.Simulated {
		if _, ok := env.URLs[networkSettings.Name]; !ok {
			return nil, fmt.Errorf("network %s not found in environment", networkSettings.Name)
		}
		if env == nil {
			log.Warn().Str("Network", networkSettings.Name).Msg("No test environment deployed")
		} else {
			networkSettings.URLs = env.URLs[networkSettings.Name]
		}
	} else {
		if len(networkSettings.URLs) == 0 {
			return nil, fmt.Errorf("no URL is provided to connect to network %s", networkSettings.Name)
		}
	}
	for idx, networkURL := range networkSettings.URLs {
		networkSettings.URL = networkURL
		ec, err := newEVMClient(networkSettings)
		ec.SyncNonce(existing)
		if err != nil {
			return nil, err
		}
		ec.SetID(idx)
		ecl.Clients = append(ecl.Clients, ec)
	}

	ecl.DefaultClient = ecl.Clients[0]
	wrappedClient := wrapMultiClient(networkSettings, ecl)
	// required in Geth when you need to call "simulate" transactions from nodes
	if ecl.NetworkSimulated() {
		if err := ecl.Fund("0x0", big.NewFloat(1000)); err != nil {
			return nil, err
		}
	}
	fundingBalance, err := wrappedClient.BalanceAt(context.Background(), wrappedClient.GetDefaultWallet().address)
	if err != nil {
		return nil, err
	}
	log.Debug().
		Str("Address", wrappedClient.GetDefaultWallet().address.Hex()).
		Uint64("Balance", fundingBalance.Uint64()).
		Msg("Funding Wallet")

	return wrappedClient, nil
}

// Get gets default client as an interface{}
func (e *EthereumMultinodeClient) Get() interface{} {
	return e.DefaultClient
}
func (e *EthereumMultinodeClient) GetNonceSetting() NonceSettings {
	return e.DefaultClient.GetNonceSetting()
}

func (e *EthereumMultinodeClient) SyncNonce(c EVMClient) {
	e.DefaultClient.SyncNonce(c)
}

// GetNetworkName gets the ID of the chain that the clients are connected to
func (e *EthereumMultinodeClient) GetNetworkName() string {
	return e.DefaultClient.GetNetworkName()
}

// GetNetworkType retrieves the type of network this is running on
func (e *EthereumMultinodeClient) NetworkSimulated() bool {
	return e.DefaultClient.NetworkSimulated()
}

// GetChainID retrieves the ChainID of the network that the client interacts with
func (e *EthereumMultinodeClient) GetChainID() *big.Int {
	return e.DefaultClient.GetChainID()
}

// GetClients gets clients for all nodes connected
func (e *EthereumMultinodeClient) GetClients() []EVMClient {
	cl := make([]EVMClient, 0)
	cl = append(cl, e.Clients...)
	return cl
}

// GetDefaultWallet returns the default wallet for the network
func (e *EthereumMultinodeClient) GetDefaultWallet() *EthereumWallet {
	return e.DefaultClient.GetDefaultWallet()
}

// GetWallets returns the default wallet for the network
func (e *EthereumMultinodeClient) GetWallets() []*EthereumWallet {
	return e.DefaultClient.GetWallets()
}

// GetNetworkConfig return the network config
func (e *EthereumMultinodeClient) GetNetworkConfig() *EVMNetwork {
	return e.DefaultClient.GetNetworkConfig()
}

// SetID sets client ID in a multi-node environment
func (e *EthereumMultinodeClient) SetID(id int) {
	e.DefaultClient.SetID(id)
}

// SetDefaultWallet sets default wallet
func (e *EthereumMultinodeClient) SetDefaultWallet(num int) error {
	return e.DefaultClient.SetDefaultWallet(num)
}

// SetWallets sets the default client's wallets
func (e *EthereumMultinodeClient) SetWallets(wallets []*EthereumWallet) {
	e.DefaultClient.SetWallets(wallets)
}

// LoadWallets loads wallets using private keys provided in the config
func (e *EthereumMultinodeClient) LoadWallets(cfg EVMNetwork) error {
	pkStrings := cfg.PrivateKeys
	wallets := make([]*EthereumWallet, 0)
	for _, pks := range pkStrings {
		w, err := NewEthereumWallet(pks)
		if err != nil {
			return err
		}
		wallets = append(wallets, w)
	}
	for _, c := range e.Clients {
		c.SetWallets(wallets)
	}
	return nil
}

// BalanceAt returns the ETH balance of the specified address
func (e *EthereumMultinodeClient) BalanceAt(ctx context.Context, address common.Address) (*big.Int, error) {
	return e.DefaultClient.BalanceAt(ctx, address)
}

// SwitchNode sets default client to perform calls to the network
func (e *EthereumMultinodeClient) SwitchNode(clientID int) error {
	if clientID > len(e.Clients) {
		return fmt.Errorf("client for node %d not found", clientID)
	}
	e.DefaultClient = e.Clients[clientID]
	return nil
}

// HeaderHashByNumber gets header hash by block number
func (e *EthereumMultinodeClient) HeaderHashByNumber(ctx context.Context, bn *big.Int) (string, error) {
	return e.DefaultClient.HeaderHashByNumber(ctx, bn)
}

// HeaderTimestampByNumber gets header timestamp by number
func (e *EthereumMultinodeClient) HeaderTimestampByNumber(ctx context.Context, bn *big.Int) (uint64, error) {
	return e.DefaultClient.HeaderTimestampByNumber(ctx, bn)
}

// LatestBlockNumber gets the latest block number from the default client
func (e *EthereumMultinodeClient) LatestBlockNumber(ctx context.Context) (uint64, error) {
	return e.DefaultClient.LatestBlockNumber(ctx)
}

// SendTransaction wraps ethereum's SendTransaction to make it safe with instant transaction types
func (e *EthereumMultinodeClient) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	return e.DefaultClient.SendTransaction(ctx, tx)
}

// Fund funds a specified address with ETH from the given wallet
func (e *EthereumMultinodeClient) Fund(toAddress string, nativeAmount *big.Float) error {
	return e.DefaultClient.Fund(toAddress, nativeAmount)
}

func (e *EthereumMultinodeClient) ReturnFunds(fromKey *ecdsa.PrivateKey) error {
	return e.DefaultClient.ReturnFunds(fromKey)
}

// DeployContract deploys a specified contract
func (e *EthereumMultinodeClient) DeployContract(
	contractName string,
	deployer ContractDeployer,
) (*common.Address, *types.Transaction, interface{}, error) {
	return e.DefaultClient.DeployContract(contractName, deployer)
}

// TransactionOpts returns the base Tx options for 'transactions' that interact with a smart contract. Since most
// contract interactions in this framework are designed to happen through abigen calls, it's intentionally quite bare.
func (e *EthereumMultinodeClient) TransactionOpts(from *EthereumWallet) (*bind.TransactOpts, error) {
	return e.DefaultClient.TransactionOpts(from)
}

// ProcessTransaction returns the result of the default client's processed transaction
func (e *EthereumMultinodeClient) ProcessTransaction(tx *types.Transaction) error {
	return e.DefaultClient.ProcessTransaction(tx)
}

// ProcessEvent returns the result of the default client's processed event
func (e *EthereumMultinodeClient) ProcessEvent(name string, event *types.Log, confirmedChan chan bool, errorChan chan error) error {
	return e.DefaultClient.ProcessEvent(name, event, confirmedChan, errorChan)
}

// IsTxConfirmed returns the default client's transaction confirmations
func (e *EthereumMultinodeClient) IsTxConfirmed(txHash common.Hash) (bool, error) {
	return e.DefaultClient.IsTxConfirmed(txHash)
}

// IsEventConfirmed returns if the default client can confirm the event has happened
func (e *EthereumMultinodeClient) IsEventConfirmed(event *types.Log) (confirmed, removed bool, err error) {
	return e.DefaultClient.IsEventConfirmed(event)
}

// GetTxReceipt returns the receipt of the transaction if available, error otherwise
func (e *EthereumMultinodeClient) GetTxReceipt(txHash common.Hash) (*types.Receipt, error) {
	return e.DefaultClient.GetTxReceipt(txHash)
}

// ParallelTransactions when enabled, sends the transaction without waiting for transaction confirmations. The hashes
// are then stored within the client and confirmations can be waited on by calling WaitForEvents.
// When disabled, the minimum confirmations are waited on when the transaction is sent, so parallelisation is disabled.
func (e *EthereumMultinodeClient) ParallelTransactions(enabled bool) {
	for _, c := range e.Clients {
		c.ParallelTransactions(enabled)
	}
}

// Close tears down the all the clients
func (e *EthereumMultinodeClient) Close() error {
	for _, c := range e.Clients {
		if err := c.Close(); err != nil {
			return err
		}
	}
	return nil
}

func (e *EthereumMultinodeClient) EstimateTransactionGasCost() (*big.Int, error) {
	return e.DefaultClient.EstimateTransactionGasCost()
}

// GasStats gets gas stats instance
func (e *EthereumMultinodeClient) GasStats() *GasStats {
	return e.DefaultClient.GasStats()
}

// AddHeaderEventSubscription adds a new header subscriber within the client to receive new headers
func (e *EthereumMultinodeClient) AddHeaderEventSubscription(key string, subscriber HeaderEventSubscription) {
	e.DefaultClient.AddHeaderEventSubscription(key, subscriber)
}

// SubscribeFilterLogs subscribes to the results of a streaming filter query.
func (e *EthereumMultinodeClient) SubscribeFilterLogs(ctx context.Context, q ethereum.FilterQuery, logs chan<- types.Log) (ethereum.Subscription, error) {
	return e.DefaultClient.SubscribeFilterLogs(ctx, q, logs)
}

// DeleteHeaderEventSubscription removes a header subscriber from the map
func (e *EthereumMultinodeClient) DeleteHeaderEventSubscription(key string) {
	e.DefaultClient.DeleteHeaderEventSubscription(key)
}

// WaitForEvents is a blocking function that waits for all event subscriptions for all clients
func (e *EthereumMultinodeClient) WaitForEvents() error {
	g := errgroup.Group{}
	for _, c := range e.Clients {
		c := c
		g.Go(func() error {
			return c.WaitForEvents()
		})
	}
	return g.Wait()
}

// LogRevertReason prints the revert reason for the transaction error by parsing through abi defined error list
func LogRevertReason(err error, abiString string) error {
	var dataError rpc.DataError
	if errors.As(err, dataError) && dataError.ErrorData() != nil {
		data, err := hex.DecodeString(dataError.ErrorData().(string)[2:])
		if err != nil {
			return err
		}
		jsonABI, err := abi.JSON(strings.NewReader(abiString))
		if err != nil {
			return err
		}
		for k, abiError := range jsonABI.Errors {
			if bytes.Equal(data[:4], abiError.ID.Bytes()[:4]) {
				// Found a matching error
				v, err := abiError.Unpack(data)
				if err != nil {
					return err
				}
				log.Info().Interface("Error", k).Interface("args - ", v).Msg("Revert Reason")
				return nil
			}
		}
		return fmt.Errorf("Revert Reason could not be found with given abistring")
	}
	return nil
}
