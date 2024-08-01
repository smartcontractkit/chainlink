package rpc

import (
	"context"
	"errors"
	"net/http"
	"net/http/cookiejar"

	"github.com/NethermindEth/juno/core/felt"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
	"golang.org/x/net/publicsuffix"
)

// ErrNotFound is returned by API methods if the requested item does not exist.
var (
	errNotFound = errors.New("not found")
)

// Provider provides the provider for starknet.go/rpc implementation.
type Provider struct {
	c       callCloser
	chainID string
}

// NewProvider creates a new rpc Provider instance.
func NewProvider(url string, options ...ethrpc.ClientOption) (*Provider, error) {
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		return nil, err
	}
	client := &http.Client{Jar: jar}
	// prepend the custom client to allow users to override
	options = append([]ethrpc.ClientOption{ethrpc.WithHTTPClient(client)}, options...)
	c, err := ethrpc.DialOptions(context.Background(), url, options...)

	if err != nil {
		return nil, err
	}

	return &Provider{c: c}, nil
}

//go:generate mockgen -destination=../mocks/mock_rpc_provider.go -package=mocks -source=provider.go api
type RpcProvider interface {
	AddInvokeTransaction(ctx context.Context, invokeTxn BroadcastInvokeTxnType) (*AddInvokeTransactionResponse, error)
	AddDeclareTransaction(ctx context.Context, declareTransaction BroadcastDeclareTxnType) (*AddDeclareTransactionResponse, error)
	AddDeployAccountTransaction(ctx context.Context, deployAccountTransaction BroadcastAddDeployTxnType) (*AddDeployAccountTransactionResponse, error)
	BlockHashAndNumber(ctx context.Context) (*BlockHashAndNumberOutput, error)
	BlockNumber(ctx context.Context) (uint64, error)
	BlockTransactionCount(ctx context.Context, blockID BlockID) (uint64, error)
	BlockWithTxHashes(ctx context.Context, blockID BlockID) (interface{}, error)
	BlockWithTxs(ctx context.Context, blockID BlockID) (interface{}, error)
	Call(ctx context.Context, call FunctionCall, block BlockID) ([]*felt.Felt, error)
	ChainID(ctx context.Context) (string, error)
	Class(ctx context.Context, blockID BlockID, classHash *felt.Felt) (ClassOutput, error)
	ClassAt(ctx context.Context, blockID BlockID, contractAddress *felt.Felt) (ClassOutput, error)
	ClassHashAt(ctx context.Context, blockID BlockID, contractAddress *felt.Felt) (*felt.Felt, error)
	EstimateFee(ctx context.Context, requests []BroadcastTxn, simulationFlags []SimulationFlag, blockID BlockID) ([]FeeEstimate, error)
	EstimateMessageFee(ctx context.Context, msg MsgFromL1, blockID BlockID) (*FeeEstimate, error)
	Events(ctx context.Context, input EventsInput) (*EventChunk, error)
	BlockWithReceipts(ctx context.Context, blockID BlockID) (interface{}, error)
	GetTransactionStatus(ctx context.Context, transactionHash *felt.Felt) (*TxnStatusResp, error)
	Nonce(ctx context.Context, blockID BlockID, contractAddress *felt.Felt) (*felt.Felt, error)
	SimulateTransactions(ctx context.Context, blockID BlockID, txns []Transaction, simulationFlags []SimulationFlag) ([]SimulatedTransaction, error)
	StateUpdate(ctx context.Context, blockID BlockID) (*StateUpdateOutput, error)
	StorageAt(ctx context.Context, contractAddress *felt.Felt, key string, blockID BlockID) (string, error)
	SpecVersion(ctx context.Context) (string, error)
	Syncing(ctx context.Context) (*SyncStatus, error)
	TraceBlockTransactions(ctx context.Context, blockID BlockID) ([]Trace, error)
	TransactionByBlockIdAndIndex(ctx context.Context, blockID BlockID, index uint64) (Transaction, error)
	TransactionByHash(ctx context.Context, hash *felt.Felt) (Transaction, error)
	TransactionReceipt(ctx context.Context, transactionHash *felt.Felt) (*TransactionReceiptWithBlockInfo, error)
	TraceTransaction(ctx context.Context, transactionHash *felt.Felt) (TxnTrace, error)
}

var _ RpcProvider = &Provider{}
