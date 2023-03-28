package contracts

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"
	ocrConfigHelper "github.com/smartcontractkit/libocr/offchainreporting/confighelper"
	ocrTypes "github.com/smartcontractkit/libocr/offchainreporting/types"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/authorized_forwarder"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/functions_billing_registry_events_mock"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/functions_oracle_events_mock"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/mock_aggregator_proxy"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/operator_factory"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/operator_wrapper"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/contracts/ethereum"

	"github.com/smartcontractkit/chainlink/integration-tests/client"
	eth_contracts "github.com/smartcontractkit/chainlink/integration-tests/contracts/ethereum"
	"github.com/smartcontractkit/chainlink/integration-tests/testreporters"
)

// EthereumOracle oracle for "directrequest" job tests
type EthereumOracle struct {
	address *common.Address
	client  blockchain.EVMClient
	oracle  *ethereum.Oracle
}

func (e *EthereumOracle) Address() string {
	return e.address.Hex()
}

func (e *EthereumOracle) Fund(ethAmount *big.Float) error {
	return e.client.Fund(e.address.Hex(), ethAmount)
}

// SetFulfillmentPermission sets fulfillment permission for particular address
func (e *EthereumOracle) SetFulfillmentPermission(address string, allowed bool) error {
	opts, err := e.client.TransactionOpts(e.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := e.oracle.SetFulfillmentPermission(opts, common.HexToAddress(address), allowed)
	if err != nil {
		return err
	}
	return e.client.ProcessTransaction(tx)
}

// EthereumAPIConsumer API consumer for job type "directrequest" tests
type EthereumAPIConsumer struct {
	address  *common.Address
	client   blockchain.EVMClient
	consumer *ethereum.APIConsumer
}

func (e *EthereumAPIConsumer) Address() string {
	return e.address.Hex()
}

func (e *EthereumAPIConsumer) RoundID(ctx context.Context) (*big.Int, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(e.client.GetDefaultWallet().Address()),
		Context: ctx,
	}
	return e.consumer.CurrentRoundID(opts)
}

func (e *EthereumAPIConsumer) Fund(ethAmount *big.Float) error {
	return e.client.Fund(e.address.Hex(), ethAmount)
}

func (e *EthereumAPIConsumer) WatchPerfEvents(ctx context.Context, eventChan chan<- *PerfEvent) error {
	ethEventChan := make(chan *ethereum.APIConsumerPerfMetricsEvent)
	sub, err := e.consumer.WatchPerfMetricsEvent(&bind.WatchOpts{}, ethEventChan)
	if err != nil {
		return err
	}
	defer sub.Unsubscribe()
	for {
		select {
		case event := <-ethEventChan:
			eventChan <- &PerfEvent{
				Contract:       e,
				RequestID:      event.RequestId,
				Round:          event.RoundID,
				BlockTimestamp: event.Timestamp,
			}
		case err := <-sub.Err():
			return err
		case <-ctx.Done():
			return nil
		}
	}
}

func (e *EthereumAPIConsumer) Data(ctx context.Context) (*big.Int, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(e.client.GetDefaultWallet().Address()),
		Context: ctx,
	}
	data, err := e.consumer.Data(opts)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// CreateRequestTo creates request to an oracle for particular jobID with params
func (e *EthereumAPIConsumer) CreateRequestTo(
	oracleAddr string,
	jobID [32]byte,
	payment *big.Int,
	url string,
	path string,
	times *big.Int,
) error {
	opts, err := e.client.TransactionOpts(e.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := e.consumer.CreateRequestTo(opts, common.HexToAddress(oracleAddr), jobID, payment, url, path, times)
	if err != nil {
		return err
	}
	return e.client.ProcessTransaction(tx)
}

// EthereumStaking
type EthereumStaking struct {
	client  blockchain.EVMClient
	staking *eth_contracts.Staking
	address *common.Address
}

func (f *EthereumStaking) Address() string {
	return f.address.Hex()
}

// Fund sends specified currencies to the contract
func (f *EthereumStaking) Fund(ethAmount *big.Float) error {
	return f.client.Fund(f.address.Hex(), ethAmount)
}

func (f *EthereumStaking) AddOperators(operators []common.Address) error {
	opts, err := f.client.TransactionOpts(f.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := f.staking.AddOperators(opts, operators)
	if err != nil {
		return err
	}
	return f.client.ProcessTransaction(tx)
}

func (f *EthereumStaking) RemoveOperators(operators []common.Address) error {
	opts, err := f.client.TransactionOpts(f.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := f.staking.RemoveOperators(opts, operators)
	if err != nil {
		return err
	}
	return f.client.ProcessTransaction(tx)
}

func (f *EthereumStaking) SetFeedOperators(operators []common.Address) error {
	opts, err := f.client.TransactionOpts(f.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := f.staking.SetFeedOperators(opts, operators)
	if err != nil {
		return err
	}
	return f.client.ProcessTransaction(tx)
}

func (f *EthereumStaking) RaiseAlert() error {
	opts, err := f.client.TransactionOpts(f.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := f.staking.RaiseAlert(opts)
	if err != nil {
		return err
	}
	return f.client.ProcessTransaction(tx)
}

func (f *EthereumStaking) Start(amount *big.Int, initialRewardRate *big.Int) error {
	opts, err := f.client.TransactionOpts(f.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := f.staking.Start(opts, amount, initialRewardRate)
	if err != nil {
		return err
	}
	return f.client.ProcessTransaction(tx)
}

func (f *EthereumStaking) SetMerkleRoot(newMerkleRoot [32]byte) error {
	opts, err := f.client.TransactionOpts(f.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := f.staking.SetMerkleRoot(opts, newMerkleRoot)
	if err != nil {
		return err
	}
	return f.client.ProcessTransaction(tx)
}

// EthereumFunctionsOracleEventsMock represents the basic events mock contract
type EthereumFunctionsOracleEventsMock struct {
	client     blockchain.EVMClient
	eventsMock *functions_oracle_events_mock.FunctionsOracleEventsMock
	address    *common.Address
}

func (f *EthereumFunctionsOracleEventsMock) Address() string {
	return f.address.Hex()
}

func (f *EthereumFunctionsOracleEventsMock) OracleResponse(requestId [32]byte) error {
	opts, err := f.client.TransactionOpts(f.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := f.eventsMock.EmitOracleResponse(opts, requestId)
	if err != nil {
		return err
	}
	return f.client.ProcessTransaction(tx)
}

func (f *EthereumFunctionsOracleEventsMock) OracleRequest(requestId [32]byte, requestingContract common.Address, requestInitiator common.Address, subscriptionId uint64, subscriptionOwner common.Address, data []byte) error {
	opts, err := f.client.TransactionOpts(f.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := f.eventsMock.EmitOracleRequest(opts, requestId, requestingContract, requestInitiator, subscriptionId, subscriptionOwner, data)
	if err != nil {
		return err
	}
	return f.client.ProcessTransaction(tx)
}

func (f *EthereumFunctionsOracleEventsMock) UserCallbackError(requestId [32]byte, reason string) error {
	opts, err := f.client.TransactionOpts(f.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := f.eventsMock.EmitUserCallbackError(opts, requestId, reason)
	if err != nil {
		return err
	}
	return f.client.ProcessTransaction(tx)
}

func (f *EthereumFunctionsOracleEventsMock) UserCallbackRawError(requestId [32]byte, lowLevelData []byte) error {
	opts, err := f.client.TransactionOpts(f.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := f.eventsMock.EmitUserCallbackRawError(opts, requestId, lowLevelData)
	if err != nil {
		return err
	}
	return f.client.ProcessTransaction(tx)
}

// EthereumFunctionsBillingRegistryEventsMock represents the basic events mock contract
type EthereumFunctionsBillingRegistryEventsMock struct {
	client     blockchain.EVMClient
	eventsMock *functions_billing_registry_events_mock.FunctionsBillingRegistryEventsMock
	address    *common.Address
}

func (f *EthereumFunctionsBillingRegistryEventsMock) Address() string {
	return f.address.Hex()
}

func (f *EthereumFunctionsBillingRegistryEventsMock) SubscriptionFunded(subscriptionId uint64, oldBalance *big.Int, newBalance *big.Int) error {
	opts, err := f.client.TransactionOpts(f.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := f.eventsMock.EmitSubscriptionFunded(opts, subscriptionId, oldBalance, newBalance)
	if err != nil {
		return err
	}
	return f.client.ProcessTransaction(tx)
}

func (f *EthereumFunctionsBillingRegistryEventsMock) BillingStart(requestId [32]byte, commitment functions_billing_registry_events_mock.FunctionsBillingRegistryEventsMockCommitment) error {
	opts, err := f.client.TransactionOpts(f.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := f.eventsMock.EmitBillingStart(opts, requestId, commitment)
	if err != nil {
		return err
	}
	return f.client.ProcessTransaction(tx)
}

func (f *EthereumFunctionsBillingRegistryEventsMock) BillingEnd(requestId [32]byte, subscriptionId uint64, signerPayment *big.Int, transmitterPayment *big.Int, totalCost *big.Int, success bool) error {
	opts, err := f.client.TransactionOpts(f.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := f.eventsMock.EmitBillingEnd(opts, requestId, subscriptionId, signerPayment, transmitterPayment, totalCost, success)
	if err != nil {
		return err
	}
	return f.client.ProcessTransaction(tx)
}

// EthereumFluxAggregator represents the basic flux aggregation contract
type EthereumFluxAggregator struct {
	client         blockchain.EVMClient
	fluxAggregator *ethereum.FluxAggregator
	address        *common.Address
}

func (f *EthereumFluxAggregator) Address() string {
	return f.address.Hex()
}

// Fund sends specified currencies to the contract
func (f *EthereumFluxAggregator) Fund(ethAmount *big.Float) error {
	return f.client.Fund(f.address.Hex(), ethAmount)
}

func (f *EthereumFluxAggregator) UpdateAvailableFunds() error {
	opts, err := f.client.TransactionOpts(f.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := f.fluxAggregator.UpdateAvailableFunds(opts)
	if err != nil {
		return err
	}
	return f.client.ProcessTransaction(tx)
}

func (f *EthereumFluxAggregator) PaymentAmount(ctx context.Context) (*big.Int, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(f.client.GetDefaultWallet().Address()),
		Context: ctx,
	}
	payment, err := f.fluxAggregator.PaymentAmount(opts)
	if err != nil {
		return nil, err
	}
	return payment, nil
}

func (f *EthereumFluxAggregator) RequestNewRound(ctx context.Context) error {
	opts, err := f.client.TransactionOpts(f.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := f.fluxAggregator.RequestNewRound(opts)
	if err != nil {
		return err
	}
	return f.client.ProcessTransaction(tx)
}

// WatchSubmissionReceived subscribes to any submissions on a flux feed
func (f *EthereumFluxAggregator) WatchSubmissionReceived(ctx context.Context, eventChan chan<- *SubmissionEvent) error {
	ethEventChan := make(chan *ethereum.FluxAggregatorSubmissionReceived)
	sub, err := f.fluxAggregator.WatchSubmissionReceived(&bind.WatchOpts{}, ethEventChan, nil, nil, nil)
	if err != nil {
		return err
	}
	defer sub.Unsubscribe()

	for {
		select {
		case event := <-ethEventChan:
			eventChan <- &SubmissionEvent{
				Contract:    event.Raw.Address,
				Submission:  event.Submission,
				Round:       event.Round,
				BlockNumber: event.Raw.BlockNumber,
				Oracle:      event.Oracle,
			}
		case err := <-sub.Err():
			return err
		case <-ctx.Done():
			return nil
		}
	}
}

func (f *EthereumFluxAggregator) SetRequesterPermissions(ctx context.Context, addr common.Address, authorized bool, roundsDelay uint32) error {
	opts, err := f.client.TransactionOpts(f.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := f.fluxAggregator.SetRequesterPermissions(opts, addr, authorized, roundsDelay)
	if err != nil {
		return err
	}
	return f.client.ProcessTransaction(tx)
}

func (f *EthereumFluxAggregator) GetOracles(ctx context.Context) ([]string, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(f.client.GetDefaultWallet().Address()),
		Context: ctx,
	}
	addresses, err := f.fluxAggregator.GetOracles(opts)
	if err != nil {
		return nil, err
	}
	var oracleAddrs []string
	for _, o := range addresses {
		oracleAddrs = append(oracleAddrs, o.Hex())
	}
	return oracleAddrs, nil
}

func (f *EthereumFluxAggregator) LatestRoundID(ctx context.Context) (*big.Int, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(f.client.GetDefaultWallet().Address()),
		Context: ctx,
	}
	rID, err := f.fluxAggregator.LatestRound(opts)
	if err != nil {
		return nil, err
	}
	return rID, nil
}

func (f *EthereumFluxAggregator) WithdrawPayment(
	ctx context.Context,
	from common.Address,
	to common.Address,
	amount *big.Int) error {
	opts, err := f.client.TransactionOpts(f.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := f.fluxAggregator.WithdrawPayment(opts, from, to, amount)
	if err != nil {
		return err
	}
	return f.client.ProcessTransaction(tx)
}

func (f *EthereumFluxAggregator) WithdrawablePayment(ctx context.Context, addr common.Address) (*big.Int, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(f.client.GetDefaultWallet().Address()),
		Context: ctx,
	}
	balance, err := f.fluxAggregator.WithdrawablePayment(opts, addr)
	if err != nil {
		return nil, err
	}
	return balance, nil
}

func (f *EthereumFluxAggregator) LatestRoundData(ctx context.Context) (RoundData, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(f.client.GetDefaultWallet().Address()),
		Context: ctx,
	}
	lr, err := f.fluxAggregator.LatestRoundData(opts)
	if err != nil {
		return RoundData{}, err
	}
	return lr, nil
}

// GetContractData retrieves basic data for the flux aggregator contract
func (f *EthereumFluxAggregator) GetContractData(ctx context.Context) (*FluxAggregatorData, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(f.client.GetDefaultWallet().Address()),
		Context: ctx,
	}

	allocated, err := f.fluxAggregator.AllocatedFunds(opts)
	if err != nil {
		return &FluxAggregatorData{}, err
	}

	available, err := f.fluxAggregator.AvailableFunds(opts)
	if err != nil {
		return &FluxAggregatorData{}, err
	}

	lr, err := f.fluxAggregator.LatestRoundData(opts)
	if err != nil {
		return &FluxAggregatorData{}, err
	}
	latestRound := RoundData(lr)

	oracles, err := f.fluxAggregator.GetOracles(opts)
	if err != nil {
		return &FluxAggregatorData{}, err
	}

	return &FluxAggregatorData{
		AllocatedFunds:  allocated,
		AvailableFunds:  available,
		LatestRoundData: latestRound,
		Oracles:         oracles,
	}, nil
}

// SetOracles allows the ability to add and/or remove oracles from the contract, and to set admins
func (f *EthereumFluxAggregator) SetOracles(o FluxAggregatorSetOraclesOptions) error {
	opts, err := f.client.TransactionOpts(f.client.GetDefaultWallet())
	if err != nil {
		return err
	}

	tx, err := f.fluxAggregator.ChangeOracles(opts, o.RemoveList, o.AddList, o.AdminList, o.MinSubmissions, o.MaxSubmissions, o.RestartDelayRounds)
	if err != nil {
		return err
	}
	return f.client.ProcessTransaction(tx)
}

// Description returns the description of the flux aggregator contract
func (f *EthereumFluxAggregator) Description(ctxt context.Context) (string, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(f.client.GetDefaultWallet().Address()),
		Context: ctxt,
	}
	return f.fluxAggregator.Description(opts)
}

// FluxAggregatorRoundConfirmer is a header subscription that awaits for a certain flux round to be completed
type FluxAggregatorRoundConfirmer struct {
	fluxInstance FluxAggregator
	roundID      *big.Int
	doneChan     chan struct{}
	context      context.Context
	cancel       context.CancelFunc
	complete     bool
}

// NewFluxAggregatorRoundConfirmer provides a new instance of a FluxAggregatorRoundConfirmer
func NewFluxAggregatorRoundConfirmer(
	contract FluxAggregator,
	roundID *big.Int,
	timeout time.Duration,
) *FluxAggregatorRoundConfirmer {
	ctx, ctxCancel := context.WithTimeout(context.Background(), timeout)
	return &FluxAggregatorRoundConfirmer{
		fluxInstance: contract,
		roundID:      roundID,
		doneChan:     make(chan struct{}),
		context:      ctx,
		cancel:       ctxCancel,
	}
}

// ReceiveHeader will query the latest FluxAggregator round and check to see whether the round has confirmed
func (f *FluxAggregatorRoundConfirmer) ReceiveHeader(header blockchain.NodeHeader) error {
	if f.complete {
		return nil
	}
	lr, err := f.fluxInstance.LatestRoundID(context.Background())
	if err != nil {
		return err
	}
	logFields := map[string]any{
		"Contract Address":  f.fluxInstance.Address(),
		"Current Round":     lr.Int64(),
		"Waiting for Round": f.roundID.Int64(),
		"Header Number":     header.Number.Uint64(),
	}
	if lr.Cmp(f.roundID) >= 0 {
		log.Info().Fields(logFields).Msg("FluxAggregator round completed")
		f.complete = true
		f.doneChan <- struct{}{}
	} else {
		log.Debug().Fields(logFields).Msg("Waiting for FluxAggregator round")
	}
	return nil
}

// Wait is a blocking function that will wait until the round has confirmed, and timeout if the deadline has passed
func (f *FluxAggregatorRoundConfirmer) Wait() error {
	defer func() { f.complete = true }()
	for {
		select {
		case <-f.doneChan:
			f.cancel()
			return nil
		case <-f.context.Done():
			return fmt.Errorf("timeout waiting for flux round to confirm: %d", f.roundID)
		}
	}
}

func (f *FluxAggregatorRoundConfirmer) Complete() bool {
	return f.complete
}

// EthereumLinkToken represents a LinkToken address
type EthereumLinkToken struct {
	client   blockchain.EVMClient
	instance *ethereum.LinkToken
	address  common.Address
}

// Fund the LINK Token contract with ETH to distribute the token
func (l *EthereumLinkToken) Fund(ethAmount *big.Float) error {
	return l.client.Fund(l.address.Hex(), ethAmount)
}

func (l *EthereumLinkToken) BalanceOf(ctx context.Context, addr string) (*big.Int, error) {
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

// Name returns the name of the link token
func (l *EthereumLinkToken) Name(ctxt context.Context) (string, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(l.client.GetDefaultWallet().Address()),
		Context: ctxt,
	}
	return l.instance.Name(opts)
}

func (l *EthereumLinkToken) Address() string {
	return l.address.Hex()
}

func (l *EthereumLinkToken) Approve(to string, amount *big.Int) error {
	opts, err := l.client.TransactionOpts(l.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	log.Info().
		Str("From", l.client.GetDefaultWallet().Address()).
		Str("To", to).
		Str("Amount", amount.String()).
		Uint64("Nonce", opts.Nonce.Uint64()).
		Msg("Approving LINK Transfer")
	tx, err := l.instance.Approve(opts, common.HexToAddress(to), amount)
	if err != nil {
		return err
	}
	return l.client.ProcessTransaction(tx)
}

func (l *EthereumLinkToken) Transfer(to string, amount *big.Int) error {
	opts, err := l.client.TransactionOpts(l.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	log.Info().
		Str("From", l.client.GetDefaultWallet().Address()).
		Str("To", to).
		Str("Amount", amount.String()).
		Uint64("Nonce", opts.Nonce.Uint64()).
		Msg("Transferring LINK")
	tx, err := l.instance.Transfer(opts, common.HexToAddress(to), amount)
	if err != nil {
		return err
	}
	return l.client.ProcessTransaction(tx)
}

func (l *EthereumLinkToken) TransferAndCall(to string, amount *big.Int, data []byte) (*types.Transaction, error) {
	opts, err := l.client.TransactionOpts(l.client.GetDefaultWallet())
	if err != nil {
		return nil, err
	}
	tx, err := l.instance.TransferAndCall(opts, common.HexToAddress(to), amount, data)
	if err != nil {
		return nil, err
	}
	log.Info().
		Str("From", l.client.GetDefaultWallet().Address()).
		Str("To", to).
		Str("Amount", amount.String()).
		Uint64("Nonce", opts.Nonce.Uint64()).
		Str("TxHash", tx.Hash().String()).
		Msg("Transferring and Calling LINK")
	return tx, l.client.ProcessTransaction(tx)
}

// LoadExistingLinkToken loads an EthereumLinkToken with a specific address
func (l *EthereumLinkToken) LoadExistingLinkToken(address string, client blockchain.EVMClient) error {
	l.address = common.HexToAddress(address)
	instance, err := ethereum.NewLinkToken(l.address, client.(*blockchain.EthereumClient).Client)
	if err != nil {
		return err
	}
	l.client = client
	l.instance = instance
	return nil
}

// EthereumOffchainAggregator represents the offchain aggregation contract
type EthereumOffchainAggregator struct {
	client  blockchain.EVMClient
	ocr     *ethereum.OffchainAggregator
	address *common.Address
}

// Fund sends specified currencies to the contract
func (o *EthereumOffchainAggregator) Fund(ethAmount *big.Float) error {
	return o.client.Fund(o.address.Hex(), ethAmount)
}

// GetContractData retrieves basic data for the offchain aggregator contract
func (o *EthereumOffchainAggregator) GetContractData(ctxt context.Context) (*OffchainAggregatorData, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(o.client.GetDefaultWallet().Address()),
		Context: ctxt,
	}

	lr, err := o.ocr.LatestRoundData(opts)
	if err != nil {
		return &OffchainAggregatorData{}, err
	}
	latestRound := RoundData(lr)

	return &OffchainAggregatorData{
		LatestRoundData: latestRound,
	}, nil
}

// SetPayees sets wallets for the contract to pay out to?
func (o *EthereumOffchainAggregator) SetPayees(
	transmitters, payees []string,
) error {
	opts, err := o.client.TransactionOpts(o.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	var transmittersAddr, payeesAddr []common.Address
	for _, tr := range transmitters {
		transmittersAddr = append(transmittersAddr, common.HexToAddress(tr))
	}
	for _, p := range payees {
		payeesAddr = append(payeesAddr, common.HexToAddress(p))
	}

	log.Info().
		Str("Transmitters", fmt.Sprintf("%v", transmitters)).
		Str("Payees", fmt.Sprintf("%v", payees)).
		Str("OCR Address", o.Address()).
		Msg("Setting OCR Payees")

	tx, err := o.ocr.SetPayees(opts, transmittersAddr, payeesAddr)
	if err != nil {
		return err
	}
	return o.client.ProcessTransaction(tx)
}

// SetConfig sets the payees and the offchain reporting protocol configuration
func (o *EthereumOffchainAggregator) SetConfig(
	chainlinkNodes []*client.Chainlink,
	ocrConfig OffChainAggregatorConfig,
	transmitters []common.Address,
) error {
	// Gather necessary addresses and keys from our chainlink nodes to properly configure the OCR contract
	log.Info().Str("Contract Address", o.address.Hex()).Msg("Configuring OCR Contract")
	for i, node := range chainlinkNodes {
		ocrKeys, err := node.MustReadOCRKeys()
		if err != nil {
			return err
		}
		primaryOCRKey := ocrKeys.Data[0]
		if err != nil {
			return err
		}
		p2pKeys, err := node.MustReadP2PKeys()
		if err != nil {
			return err
		}
		primaryP2PKey := p2pKeys.Data[0]

		// Need to convert the key representations
		var onChainSigningAddress [20]byte
		var configPublicKey [32]byte
		offchainSigningAddress, err := hex.DecodeString(primaryOCRKey.Attributes.OffChainPublicKey)
		if err != nil {
			return err
		}
		decodeConfigKey, err := hex.DecodeString(primaryOCRKey.Attributes.ConfigPublicKey)
		if err != nil {
			return err
		}

		// https://stackoverflow.com/questions/8032170/how-to-assign-string-to-bytes-array
		copy(onChainSigningAddress[:], common.HexToAddress(primaryOCRKey.Attributes.OnChainSigningAddress).Bytes())
		copy(configPublicKey[:], decodeConfigKey)

		oracleIdentity := ocrConfigHelper.OracleIdentity{
			TransmitAddress:       transmitters[i],
			OnChainSigningAddress: onChainSigningAddress,
			PeerID:                primaryP2PKey.Attributes.PeerID,
			OffchainPublicKey:     offchainSigningAddress,
		}
		oracleIdentityExtra := ocrConfigHelper.OracleIdentityExtra{
			OracleIdentity:                  oracleIdentity,
			SharedSecretEncryptionPublicKey: ocrTypes.SharedSecretEncryptionPublicKey(configPublicKey),
		}

		ocrConfig.OracleIdentities = append(ocrConfig.OracleIdentities, oracleIdentityExtra)
	}

	signers, transmitters, threshold, encodedConfigVersion, encodedConfig, err := ocrConfigHelper.ContractSetConfigArgs(
		ocrConfig.DeltaProgress,
		ocrConfig.DeltaResend,
		ocrConfig.DeltaRound,
		ocrConfig.DeltaGrace,
		ocrConfig.DeltaC,
		ocrConfig.AlphaPPB,
		ocrConfig.DeltaStage,
		ocrConfig.RMax,
		ocrConfig.S,
		ocrConfig.OracleIdentities,
		ocrConfig.F,
	)
	if err != nil {
		return err
	}

	// Set Config
	opts, err := o.client.TransactionOpts(o.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := o.ocr.SetConfig(opts, signers, transmitters, threshold, encodedConfigVersion, encodedConfig)
	if err != nil {
		return err
	}
	return o.client.ProcessTransaction(tx)
}

// RequestNewRound requests the OCR contract to create a new round
func (o *EthereumOffchainAggregator) RequestNewRound() error {
	opts, err := o.client.TransactionOpts(o.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := o.ocr.RequestNewRound(opts)
	if err != nil {
		return err
	}
	log.Info().Str("Contract Address", o.address.Hex()).Msg("New OCR round requested")

	return o.client.ProcessTransaction(tx)
}

// GetLatestAnswer returns the latest answer from the OCR contract
func (o *EthereumOffchainAggregator) GetLatestAnswer(ctxt context.Context) (*big.Int, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(o.client.GetDefaultWallet().Address()),
		Context: ctxt,
	}
	return o.ocr.LatestAnswer(opts)
}

func (o *EthereumOffchainAggregator) Address() string {
	return o.address.Hex()
}

// GetLatestRound returns data from the latest round
func (o *EthereumOffchainAggregator) GetLatestRound(ctx context.Context) (*RoundData, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(o.client.GetDefaultWallet().Address()),
		Context: ctx,
	}

	roundData, err := o.ocr.LatestRoundData(opts)
	if err != nil {
		return nil, err
	}

	return &RoundData{
		RoundId:         roundData.RoundId,
		Answer:          roundData.Answer,
		AnsweredInRound: roundData.AnsweredInRound,
		StartedAt:       roundData.StartedAt,
		UpdatedAt:       roundData.UpdatedAt,
	}, err
}

func (o *EthereumOffchainAggregator) LatestRoundDataUpdatedAt() (*big.Int, error) {
	data, err := o.ocr.LatestRoundData(&bind.CallOpts{
		From:    common.HexToAddress(o.client.GetDefaultWallet().Address()),
		Context: context.Background(),
	})
	if err != nil {
		return nil, err
	}
	return data.UpdatedAt, nil
}

// GetRound retrieves an OCR round by the round ID
func (o *EthereumOffchainAggregator) GetRound(ctx context.Context, roundID *big.Int) (*RoundData, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(o.client.GetDefaultWallet().Address()),
		Context: ctx,
	}
	roundData, err := o.ocr.GetRoundData(opts, roundID)
	if err != nil {
		return nil, err
	}

	return &RoundData{
		RoundId:         roundData.RoundId,
		Answer:          roundData.Answer,
		AnsweredInRound: roundData.AnsweredInRound,
		StartedAt:       roundData.StartedAt,
		UpdatedAt:       roundData.UpdatedAt,
	}, nil
}

// ParseEventAnswerUpdated parses the log for event AnswerUpdated
func (o *EthereumOffchainAggregator) ParseEventAnswerUpdated(eventLog types.Log) (*ethereum.OffchainAggregatorAnswerUpdated, error) {
	return o.ocr.ParseAnswerUpdated(eventLog)
}

// RunlogRoundConfirmer is a header subscription that awaits for a certain Runlog round to be completed
type RunlogRoundConfirmer struct {
	consumer APIConsumer
	roundID  *big.Int
	doneChan chan struct{}
	context  context.Context
	cancel   context.CancelFunc
}

// NewRunlogRoundConfirmer provides a new instance of a RunlogRoundConfirmer
func NewRunlogRoundConfirmer(
	contract APIConsumer,
	roundID *big.Int,
	timeout time.Duration,
) *RunlogRoundConfirmer {
	ctx, ctxCancel := context.WithTimeout(context.Background(), timeout)
	return &RunlogRoundConfirmer{
		consumer: contract,
		roundID:  roundID,
		doneChan: make(chan struct{}),
		context:  ctx,
		cancel:   ctxCancel,
	}
}

// ReceiveHeader will query the latest Runlog round and check to see whether the round has confirmed
func (o *RunlogRoundConfirmer) ReceiveHeader(_ blockchain.NodeHeader) error {
	currentRoundID, err := o.consumer.RoundID(context.Background())
	if err != nil {
		return err
	}
	logFields := map[string]any{
		"Contract Address":  o.consumer.Address(),
		"Current Round":     currentRoundID.Int64(),
		"Waiting for Round": o.roundID.Int64(),
	}
	if currentRoundID.Cmp(o.roundID) >= 0 {
		log.Info().Fields(logFields).Msg("Runlog round completed")
		o.doneChan <- struct{}{}
	} else {
		log.Debug().Fields(logFields).Msg("Waiting for Runlog round")
	}
	return nil
}

// Wait is a blocking function that will wait until the round has confirmed, and timeout if the deadline has passed
func (o *RunlogRoundConfirmer) Wait() error {
	for {
		select {
		case <-o.doneChan:
			o.cancel()
			return nil
		case <-o.context.Done():
			return fmt.Errorf("timeout waiting for OCR round to confirm: %d", o.roundID)
		}
	}
}

// OffchainAggregatorRoundConfirmer is a header subscription that awaits for a certain OCR round to be completed
type OffchainAggregatorRoundConfirmer struct {
	ocrInstance        OffchainAggregator
	roundID            *big.Int
	doneChan           chan struct{}
	context            context.Context
	cancel             context.CancelFunc
	optionalTestReport *testreporters.OCRSoakTestReport
	blocksSinceAnswer  uint
	complete           bool
}

// NewOffchainAggregatorRoundConfirmer provides a new instance of a OffchainAggregatorRoundConfirmer
func NewOffchainAggregatorRoundConfirmer(
	contract OffchainAggregator,
	roundID *big.Int,
	timeout time.Duration,
	optionalTestReport *testreporters.OCRSoakTestReport,
) *OffchainAggregatorRoundConfirmer {
	ctx, ctxCancel := context.WithTimeout(context.Background(), timeout)
	return &OffchainAggregatorRoundConfirmer{
		ocrInstance:        contract,
		roundID:            roundID,
		doneChan:           make(chan struct{}),
		context:            ctx,
		cancel:             ctxCancel,
		optionalTestReport: optionalTestReport,
		complete:           false,
	}
}

// ReceiveHeader will query the latest OffchainAggregator round and check to see whether the round has confirmed
func (o *OffchainAggregatorRoundConfirmer) ReceiveHeader(_ blockchain.NodeHeader) error {
	if channelClosed(o.doneChan) {
		return nil
	}

	lr, err := o.ocrInstance.GetLatestRound(context.Background())
	if err != nil {
		return err
	}
	o.blocksSinceAnswer++
	currRound := lr.RoundId
	logFields := map[string]any{
		"Contract Address":  o.ocrInstance.Address(),
		"Current Round":     currRound.Int64(),
		"Waiting for Round": o.roundID.Int64(),
	}
	if currRound.Cmp(o.roundID) >= 0 {
		log.Info().Fields(logFields).Msg("OCR round completed")
		o.doneChan <- struct{}{}
		o.complete = true
	} else {
		log.Debug().Fields(logFields).Msg("Waiting on OCR Round")
	}
	return nil
}

// Wait is a blocking function that will wait until the round has confirmed, and timeout if the deadline has passed
func (o *OffchainAggregatorRoundConfirmer) Wait() error {
	defer func() { o.complete = true }()
	for {
		select {
		case <-o.doneChan:
			o.cancel()
			close(o.doneChan)
			return nil
		case <-o.context.Done():
			return fmt.Errorf("timeout waiting for OCR round to confirm: %d", o.roundID)
		}
	}
}

func (o *OffchainAggregatorRoundConfirmer) Complete() bool {
	return o.complete
}

// EthereumStorage acts as a conduit for the ethereum version of the storage contract
type EthereumStorage struct {
	client blockchain.EVMClient
	store  *ethereum.Store
}

// Set sets a value in the storage contract
func (e *EthereumStorage) Set(value *big.Int) error {
	opts, err := e.client.TransactionOpts(e.client.GetDefaultWallet())
	if err != nil {
		return err
	}

	tx, err := e.store.Set(opts, value)
	if err != nil {
		return err
	}
	return e.client.ProcessTransaction(tx)
}

// Get retrieves a set value from the storage contract
func (e *EthereumStorage) Get(ctxt context.Context) (*big.Int, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(e.client.GetDefaultWallet().Address()),
		Context: ctxt,
	}
	return e.store.Get(opts)
}

// EthereumMockETHLINKFeed represents mocked ETH/LINK feed contract
type EthereumMockETHLINKFeed struct {
	client  blockchain.EVMClient
	feed    *ethereum.MockETHLINKAggregator
	address *common.Address
}

func (v *EthereumMockETHLINKFeed) Address() string {
	return v.address.Hex()
}

func (v *EthereumMockETHLINKFeed) LatestRoundData() (*big.Int, error) {
	data, err := v.feed.LatestRoundData(&bind.CallOpts{
		From:    common.HexToAddress(v.client.GetDefaultWallet().Address()),
		Context: context.Background(),
	})
	if err != nil {
		return nil, err
	}
	return data.Ans, nil
}

func (v *EthereumMockETHLINKFeed) LatestRoundDataUpdatedAt() (*big.Int, error) {
	data, err := v.feed.LatestRoundData(&bind.CallOpts{
		From:    common.HexToAddress(v.client.GetDefaultWallet().Address()),
		Context: context.Background(),
	})
	if err != nil {
		return nil, err
	}
	return data.UpdatedAt, nil
}

// EthereumMockGASFeed represents mocked Gas feed contract
type EthereumMockGASFeed struct {
	client  blockchain.EVMClient
	feed    *ethereum.MockGASAggregator
	address *common.Address
}

func (v *EthereumMockGASFeed) Address() string {
	return v.address.Hex()
}

// EthereumReadAccessController represents read access controller contract
type EthereumReadAccessController struct {
	client  blockchain.EVMClient
	rac     *ethereum.SimpleReadAccessController
	address *common.Address
}

// AddAccess grants access to particular address to raise a flag
func (e *EthereumReadAccessController) AddAccess(addr string) error {
	opts, err := e.client.TransactionOpts(e.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	log.Debug().Str("Address", addr).Msg("Adding access for address")
	tx, err := e.rac.AddAccess(opts, common.HexToAddress(addr))
	if err != nil {
		return err
	}
	return e.client.ProcessTransaction(tx)
}

// DisableAccessCheck disables all access checks
func (e *EthereumReadAccessController) DisableAccessCheck() error {
	opts, err := e.client.TransactionOpts(e.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := e.rac.DisableAccessCheck(opts)
	if err != nil {
		return err
	}
	return e.client.ProcessTransaction(tx)
}

func (e *EthereumReadAccessController) Address() string {
	return e.address.Hex()
}

// EthereumFlags represents flags contract
type EthereumFlags struct {
	client  blockchain.EVMClient
	flags   *ethereum.Flags
	address *common.Address
}

func (e *EthereumFlags) Address() string {
	return e.address.Hex()
}

// GetFlag returns boolean if a flag was set for particular address
func (e *EthereumFlags) GetFlag(ctx context.Context, addr string) (bool, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(e.client.GetDefaultWallet().Address()),
		Context: ctx,
	}
	flag, err := e.flags.GetFlag(opts, common.HexToAddress(addr))
	if err != nil {
		return false, err
	}
	return flag, nil
}

// EthereumDeviationFlaggingValidator represents deviation flagging validator contract
type EthereumDeviationFlaggingValidator struct {
	client  blockchain.EVMClient
	dfv     *ethereum.DeviationFlaggingValidator
	address *common.Address
}

func (e *EthereumDeviationFlaggingValidator) Address() string {
	return e.address.Hex()
}

// EthereumOperatorFactory represents operator factory contract
type EthereumOperatorFactory struct {
	address         *common.Address
	client          blockchain.EVMClient
	operatorFactory *operator_factory.OperatorFactory
}

func (e *EthereumOperatorFactory) ParseAuthorizedForwarderCreated(eventLog types.Log) (*operator_factory.OperatorFactoryAuthorizedForwarderCreated, error) {
	return e.operatorFactory.ParseAuthorizedForwarderCreated(eventLog)
}

func (e *EthereumOperatorFactory) ParseOperatorCreated(eventLog types.Log) (*operator_factory.OperatorFactoryOperatorCreated, error) {
	return e.operatorFactory.ParseOperatorCreated(eventLog)
}

func (e *EthereumOperatorFactory) Address() string {
	return e.address.Hex()
}

func (e *EthereumOperatorFactory) DeployNewOperatorAndForwarder() (*types.Transaction, error) {
	opts, err := e.client.TransactionOpts(e.client.GetDefaultWallet())
	if err != nil {
		return nil, err
	}
	tx, err := e.operatorFactory.DeployNewOperatorAndForwarder(opts)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

// EthereumOperator represents operator contract
type EthereumOperator struct {
	address  common.Address
	client   blockchain.EVMClient
	operator *operator_wrapper.Operator
}

func (e *EthereumOperator) Address() string {
	return e.address.Hex()
}

func (e *EthereumOperator) AcceptAuthorizedReceivers(forwarders []common.Address, eoa []common.Address) error {
	opts, err := e.client.TransactionOpts(e.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	log.Info().
		Str("ForwardersAddresses", fmt.Sprint(forwarders)).
		Str("EoaAddresses", fmt.Sprint(eoa)).
		Msg("Accepting Authorized Receivers")
	tx, err := e.operator.AcceptAuthorizedReceivers(opts, forwarders, eoa)
	if err != nil {
		return err
	}
	return e.client.ProcessTransaction(tx)
}

// EthereumAuthorizedForwarder represents authorized forwarder contract
type EthereumAuthorizedForwarder struct {
	address             common.Address
	client              blockchain.EVMClient
	authorizedForwarder *authorized_forwarder.AuthorizedForwarder
}

// Owner return authorized forwarder owner address
func (e *EthereumAuthorizedForwarder) Owner(ctx context.Context) (string, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(e.client.GetDefaultWallet().Address()),
		Context: ctx,
	}
	owner, err := e.authorizedForwarder.Owner(opts)

	return owner.Hex(), err
}

func (e *EthereumAuthorizedForwarder) GetAuthorizedSenders(ctx context.Context) ([]string, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(e.client.GetDefaultWallet().Address()),
		Context: ctx,
	}
	authorizedSenders, err := e.authorizedForwarder.GetAuthorizedSenders(opts)
	if err != nil {
		return nil, err
	}
	var sendersAddrs []string
	for _, o := range authorizedSenders {
		sendersAddrs = append(sendersAddrs, o.Hex())
	}
	return sendersAddrs, nil
}

func (e *EthereumAuthorizedForwarder) Address() string {
	return e.address.Hex()
}

// EthereumMockAggregatorProxy represents mock aggregator proxy contract
type EthereumMockAggregatorProxy struct {
	address             *common.Address
	client              blockchain.EVMClient
	mockAggregatorProxy *mock_aggregator_proxy.MockAggregatorProxy
}

func (e *EthereumMockAggregatorProxy) Address() string {
	return e.address.Hex()
}

func (e *EthereumMockAggregatorProxy) UpdateAggregator(aggregator common.Address) error {
	opts, err := e.client.TransactionOpts(e.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := e.mockAggregatorProxy.UpdateAggregator(opts, aggregator)
	if err != nil {
		return err
	}
	return e.client.ProcessTransaction(tx)
}

func (e *EthereumMockAggregatorProxy) Aggregator() (common.Address, error) {
	addr, err := e.mockAggregatorProxy.Aggregator(&bind.CallOpts{
		From:    common.HexToAddress(e.client.GetDefaultWallet().Address()),
		Context: context.Background(),
	})
	if err != nil {
		return common.Address{}, err
	}
	return addr, nil
}

func channelClosed(ch <-chan struct{}) bool {
	select {
	case <-ch:
		return true
	default:
	}

	return false
}
