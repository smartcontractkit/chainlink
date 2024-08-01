package client

import (
	"context"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"time"

	"github.com/cosmos/cosmos-sdk/types/query"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink-cosmos/pkg/cosmos/params"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	rpchttp "github.com/cometbft/cometbft/rpc/client/http"
	libclient "github.com/cometbft/cometbft/rpc/jsonrpc/client"
	cosmosclient "github.com/cosmos/cosmos-sdk/client"
	tmtypes "github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	txtypes "github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

//go:generate mockery --name ReaderWriter --output ./mocks/
type ReaderWriter interface {
	Writer
	Reader
}

// Reader provides methods for reading from a cosmos chain.
type Reader interface {
	Account(address sdk.AccAddress) (uint64, uint64, error)
	ContractState(contractAddress sdk.AccAddress, queryMsg []byte) ([]byte, error)
	TxsEvents(events []string, paginationParams *query.PageRequest) (*txtypes.GetTxsEventResponse, error)
	Tx(hash string) (*txtypes.GetTxResponse, error)
	LatestBlock() (*tmtypes.GetLatestBlockResponse, error)
	BlockByHeight(height int64) (*tmtypes.GetBlockByHeightResponse, error)
	Balance(addr sdk.AccAddress, denom string) (*sdk.Coin, error)
	// TODO: escape hatch for injective client
	Context() *cosmosclient.Context
}

// Writer provides methods for writing to a cosmos chain.
// Assumes all msgs are for the same from address.
// We may want to support multiple from addresses + signers if a use case arises.
type Writer interface {
	// TODO: SignAndBroadcast is only used for testing, remove it
	SignAndBroadcast(msgs []sdk.Msg, accountNum uint64, sequence uint64, gasPrice sdk.DecCoin, signer cryptotypes.PrivKey, mode txtypes.BroadcastMode) (*txtypes.BroadcastTxResponse, error)
	Broadcast(txBytes []byte, mode txtypes.BroadcastMode) (*txtypes.BroadcastTxResponse, error)
	Simulate(txBytes []byte) (*txtypes.SimulateResponse, error)
	BatchSimulateUnsigned(msgs SimMsgs, sequence uint64) (*BatchSimResults, error)
	SimulateUnsigned(msgs []sdk.Msg, sequence uint64) (*txtypes.SimulateResponse, error)
	CreateAndSign(msgs []sdk.Msg, account uint64, sequence uint64, gasLimit uint64, gasLimitMultiplier float64, gasPrice sdk.DecCoin, signer cryptotypes.PrivKey, timeoutHeight uint64) ([]byte, error)
}

var _ ReaderWriter = (*Client)(nil)

const (
	// DefaultTimeout is the default Cosmos client timeout.
	// Note that while the cosmos node is processing a heavy block,
	// requests can be delayed significantly (https://github.com/tendermint/tendermint/issues/6899),
	// however there's nothing we can do but wait until the block is processed.
	// So we set a fairly high timeout here.
	DefaultTimeout = 30 * time.Second
	// DefaultGasLimitMultiplier is the default gas limit multiplier.
	// It scales up the gas limit for 3 reasons:
	// 1. We simulate without a fee present (since we're simulating in order to determine the fee)
	// since we simulate unsigned. The fee is included in the signing data:
	// https://github.com/cosmos/cosmos-sdk/blob/master/x/auth/tx/direct.go#L40)
	// 2. Potential state changes between estimation and execution.
	// 3. The simulation doesn't include db writes in the tendermint node
	// (https://github.com/cosmos/cosmos-sdk/issues/4938)
	DefaultGasLimitMultiplier = 1.5
)

// Client is a cosmos client
type Client struct {
	chainID                 string
	clientCtx               cosmosclient.Context
	cosmosServiceClient     txtypes.ServiceClient
	authClient              authtypes.QueryClient
	wasmClient              wasmtypes.QueryClient
	bankClient              banktypes.QueryClient
	tendermintServiceClient tmtypes.ServiceClient
	log                     logger.Logger
}

// NewClient creates a new cosmos client
func NewClient(chainID string,
	tendermintURL string,
	requestTimeout time.Duration,
	lggr logger.Logger,
) (*Client, error) {
	if requestTimeout <= 0 {
		requestTimeout = DefaultTimeout
	}

	httpClient, err := libclient.DefaultHTTPClient(tendermintURL)
	if err != nil {
		return nil, err
	}
	httpClient.Timeout = requestTimeout
	tmClient, err := rpchttp.NewWithClient(tendermintURL, "/websocket", httpClient)
	if err != nil {
		return nil, err
	}

	// Note should cosmos nodes start exposing grpc, its preferable
	// to connect directly with grpc.Dial to avoid using clientCtx (according to tendermint team).
	// If so then we would start putting timeouts on the ctx we pass in to the generate grpc client calls.
	clientCtx := params.NewClientContext().
		WithAccountRetriever(authtypes.AccountRetriever{}).
		WithClient(tmClient).
		WithChainID(chainID)

	cosmosServiceClient := txtypes.NewServiceClient(clientCtx)
	authClient := authtypes.NewQueryClient(clientCtx)
	wasmClient := wasmtypes.NewQueryClient(clientCtx)
	tendermintServiceClient := tmtypes.NewServiceClient(clientCtx)
	bankClient := banktypes.NewQueryClient(clientCtx)

	return &Client{
		chainID:                 chainID,
		cosmosServiceClient:     cosmosServiceClient,
		authClient:              authClient,
		wasmClient:              wasmClient,
		tendermintServiceClient: tendermintServiceClient,
		bankClient:              bankClient,
		clientCtx:               clientCtx,
		log:                     lggr,
	}, nil
}

func (c *Client) Context() *cosmosclient.Context {
	return &c.clientCtx
}

// Account read the account address for the account number and sequence number.
// !!Note only one sequence number can be used per account per block!!
func (c *Client) Account(addr sdk.AccAddress) (uint64, uint64, error) {
	r, err := c.authClient.Account(context.Background(), &authtypes.QueryAccountRequest{Address: addr.String()})
	if err != nil {
		return 0, 0, err
	}
	var a authtypes.AccountI
	err = c.clientCtx.InterfaceRegistry.UnpackAny(r.Account, &a)
	if err != nil {
		return 0, 0, err
	}
	return a.GetAccountNumber(), a.GetSequence(), nil
}

// ContractState reads from a WASM contract store
func (c *Client) ContractState(contractAddress sdk.AccAddress, queryMsg []byte) ([]byte, error) {
	s, err := c.wasmClient.SmartContractState(context.Background(), &wasmtypes.QuerySmartContractStateRequest{
		Address:   contractAddress.String(),
		QueryData: queryMsg,
	})
	if err != nil {
		return nil, err
	}
	//  Note s will be nil on err
	return s.Data, err
}

// TxsEvents returns in tx events in descending order (latest txes first).
// Each event is ANDed together and follows the query language defined
// https://docs.cosmos.network/master/core/events.html
// Note one current issue https://github.com/cosmos/cosmos-sdk/issues/10448
func (c *Client) TxsEvents(events []string, paginationParams *query.PageRequest) (*txtypes.GetTxsEventResponse, error) {
	e, err := c.cosmosServiceClient.GetTxsEvent(context.Background(), &txtypes.GetTxsEventRequest{
		Events:     events,
		Pagination: paginationParams,
		OrderBy:    txtypes.OrderBy_ORDER_BY_DESC,
	})
	return e, err
}

// Tx gets a tx by hash
func (c *Client) Tx(hash string) (*txtypes.GetTxResponse, error) {
	e, err := c.cosmosServiceClient.GetTx(context.Background(), &txtypes.GetTxRequest{
		Hash: hash,
	})
	return e, err
}

// LatestBlock returns the latest block
func (c *Client) LatestBlock() (*tmtypes.GetLatestBlockResponse, error) {
	return c.tendermintServiceClient.GetLatestBlock(context.Background(), &tmtypes.GetLatestBlockRequest{})
}

// BlockByHeight gets a block by height
func (c *Client) BlockByHeight(height int64) (*tmtypes.GetBlockByHeightResponse, error) {
	return c.tendermintServiceClient.GetBlockByHeight(context.Background(), &tmtypes.GetBlockByHeightRequest{Height: height})
}

// CreateAndSign creates and signs a transaction
func (c *Client) CreateAndSign(msgs []sdk.Msg, account uint64, sequence uint64, gasLimit uint64, gasLimitMultiplier float64, gasPrice sdk.DecCoin, signer cryptotypes.PrivKey, timeoutHeight uint64) ([]byte, error) {
	// https://github.com/cosmos/cosmos-sdk/blob/a785bf5af602525cf7a5c5ea097056597e2eb7ef/client/tx/tx.go#L63-L117
	// https://docs.cosmos.network/main/run-node/txs#signing-a-transaction-1
	txConfig := params.ClientTxConfig()
	txBuilder := txConfig.NewTxBuilder()
	err := txBuilder.SetMsgs(msgs...)
	if err != nil {
		return nil, err
	}
	gasLimitBuffered := uint64(math.Ceil(float64(gasLimit) * gasLimitMultiplier))
	txBuilder.SetGasLimit(gasLimitBuffered)
	gasFee := sdk.NewCoin(gasPrice.Denom, gasPrice.Amount.MulInt64(int64(gasLimitBuffered)).Ceil().RoundInt())
	txBuilder.SetFeeAmount(sdk.NewCoins(gasFee))
	// 0 timeout height means unset.
	txBuilder.SetTimeoutHeight(timeoutHeight)

	// Sign
	// https://github.com/cosmos/cosmos-sdk/blob/a785bf5af602525cf7a5c5ea097056597e2eb7ef/client/tx/tx.go#L230-L337

	pubKey := signer.PubKey()

	// signMode := txConfig.SignModeHandler().DefaultMode()
	signMode := signing.SignMode_SIGN_MODE_DIRECT

	signerData := authsigning.SignerData{
		AccountNumber: account,
		ChainID:       c.chainID,
		Sequence:      sequence,
	}

	// For SIGN_MODE_DIRECT, calling SetSignatures calls setSignerInfos on
	// TxBuilder under the hood, and SignerInfos is needed to generated the
	// sign bytes. This is the reason for setting SetSignatures here, with a
	// nil signature.
	//
	// Note: this line is not needed for SIGN_MODE_LEGACY_AMINO, but putting it
	// also doesn't affect its generated sign bytes, so for code's simplicity
	// sake, we put it here.
	sigData := signing.SingleSignatureData{
		SignMode:  signMode,
		Signature: nil,
	}
	sig := signing.SignatureV2{
		PubKey:   pubKey,
		Data:     &sigData,
		Sequence: sequence,
	}
	if err = txBuilder.SetSignatures(sig); err != nil {
		return nil, err
	}

	// Sign those bytes
	signature, err := tx.SignWithPrivKey(
		signMode,
		signerData,
		txBuilder,
		signer,
		txConfig,
		sequence,
	)
	if err != nil {
		return nil, err
	}

	if err = txBuilder.SetSignatures(signature); err != nil {
		return nil, err
	}

	// TODO: return txBuilder.GetTx() for more flexibility

	return txConfig.TxEncoder()(txBuilder.GetTx())
}

// SimMsg binds an ID to a msg
type SimMsg struct {
	ID  int64
	Msg sdk.Msg
}

// SimMsgs is a slice of SimMsg
type SimMsgs []SimMsg

// GetMsgs extracts all msgs from SimMsgs
func (simMsgs SimMsgs) GetMsgs() []sdk.Msg {
	msgs := make([]sdk.Msg, len(simMsgs))
	for i := range simMsgs {
		msgs[i] = simMsgs[i].Msg
	}
	return msgs
}

// GetSimMsgsIDs extracts all IDs from SimMsgs
func (simMsgs SimMsgs) GetSimMsgsIDs() []int64 {
	ids := make([]int64, len(simMsgs))
	for i := range simMsgs {
		ids[i] = simMsgs[i].ID
	}
	return ids
}

// BatchSimResults indicates which msgs failed and which succeeded
type BatchSimResults struct {
	Failed    SimMsgs
	Succeeded SimMsgs
}

var failedMsgIndexRe = regexp.MustCompile(`^.*failed to execute message; message index: (?P<Index>\d+):.*$`)

func (c *Client) failedMsgIndex(err error) (bool, int) {
	if err == nil {
		return false, 0
	}

	m := failedMsgIndexRe.FindStringSubmatch(err.Error())
	if len(m) != 2 {
		return false, 0
	}
	index, err := strconv.ParseInt(m[1], 10, 32)
	if err != nil {
		return false, 0
	}
	return true, int(index)
}

// BatchSimulateUnsigned simulates a group of msgs.
// Assumes at least one msg is present.
// If we fail to simulate the batch, remove the offending tx
// and try again. Repeat until we have a successful batch.
// Keep track of failures so we can mark them as errored.
// Note that the error from simulating indicates the first
// msg in the slice which failed (it simply loops over the msgs
// and simulates them one by one, breaking at the first failure).
func (c *Client) BatchSimulateUnsigned(msgs SimMsgs, sequence uint64) (*BatchSimResults, error) {
	var succeeded []SimMsg
	var failed []SimMsg
	toSim := msgs
	for {
		_, err := c.SimulateUnsigned(toSim.GetMsgs(), sequence)
		containsFailure, failureIndex := c.failedMsgIndex(err)
		if err != nil && !containsFailure {
			return nil, err
		}
		if !containsFailure {
			// we're done they all succeeded
			succeeded = append(succeeded, toSim...)
			break
		}
		failed = append(failed, toSim[failureIndex])
		succeeded = append(succeeded, toSim[:failureIndex]...)
		// remove offending msg and retry
		if failureIndex == len(toSim)-1 {
			// we're done, last one failed
			c.log.Warnf("simulation error found in last msg, failure %v, index %v, err %v", toSim[failureIndex], failureIndex, err)
			break
		}
		// otherwise there may be more to sim
		c.log.Warnf("simulation error found in a msg, retrying with %v, failure %v, index %v, err %v", toSim[failureIndex+1:], toSim[failureIndex], failureIndex, err)
		toSim = toSim[failureIndex+1:]
	}
	return &BatchSimResults{
		Failed:    failed,
		Succeeded: succeeded,
	}, nil
}

// SimulateUnsigned simulates an unsigned msg
func (c *Client) SimulateUnsigned(msgs []sdk.Msg, sequence uint64) (*txtypes.SimulateResponse, error) {
	txConfig := params.ClientTxConfig()
	txBuilder := txConfig.NewTxBuilder()
	if err := txBuilder.SetMsgs(msgs...); err != nil {
		return nil, err
	}
	// Create an empty signature literal as the ante handler will populate with a
	// sentinel pubkey.
	// Note the simulation actually won't work without this
	sig := signing.SignatureV2{
		PubKey: &secp256k1.PubKey{},
		Data: &signing.SingleSignatureData{
			SignMode: signing.SignMode_SIGN_MODE_DIRECT,
		},
		Sequence: sequence,
	}
	if err := txBuilder.SetSignatures(sig); err != nil {
		return nil, err
	}
	txBytes, err := txConfig.TxEncoder()(txBuilder.GetTx())
	if err != nil {
		return nil, err
	}
	s, err := c.cosmosServiceClient.Simulate(context.Background(), &txtypes.SimulateRequest{
		TxBytes: txBytes,
	})
	return s, err
}

// Simulate simulates a signed transaction
func (c *Client) Simulate(txBytes []byte) (*txtypes.SimulateResponse, error) {
	s, err := c.cosmosServiceClient.Simulate(context.Background(), &txtypes.SimulateRequest{
		TxBytes: txBytes,
	})
	return s, err
}

// Broadcast broadcasts a tx
func (c *Client) Broadcast(txBytes []byte, mode txtypes.BroadcastMode) (*txtypes.BroadcastTxResponse, error) {
	res, err := c.cosmosServiceClient.BroadcastTx(context.Background(), &txtypes.BroadcastTxRequest{
		Mode:    mode,
		TxBytes: txBytes,
	})
	if err != nil {
		return nil, err
	}
	if res.TxResponse == nil {
		return nil, fmt.Errorf("got nil tx response")
	}
	if res.TxResponse.Code != 0 {
		return res, fmt.Errorf("tx failed with error code: %d, resp %v", res.TxResponse.Code, res.TxResponse)
	}
	return res, err
}

// SignAndBroadcast signs and broadcasts a group of msgs.
func (c *Client) SignAndBroadcast(msgs []sdk.Msg, account uint64, sequence uint64, gasPrice sdk.DecCoin, signer cryptotypes.PrivKey, mode txtypes.BroadcastMode) (*txtypes.BroadcastTxResponse, error) {
	sim, err := c.SimulateUnsigned(msgs, sequence)
	if err != nil {
		return nil, err
	}
	// TODO: replace with BroadcastTx()?
	txBytes, err := c.CreateAndSign(msgs, account, sequence, sim.GasInfo.GasUsed, DefaultGasLimitMultiplier, gasPrice, signer, 0)
	if err != nil {
		return nil, err
	}
	return c.Broadcast(txBytes, mode)
}

// Balance returns the balance of an address
func (c *Client) Balance(addr sdk.AccAddress, denom string) (*sdk.Coin, error) {
	b, err := c.bankClient.Balance(context.Background(), &banktypes.QueryBalanceRequest{Address: addr.String(), Denom: denom})
	if err != nil {
		return nil, err
	}
	return b.Balance, nil
}
