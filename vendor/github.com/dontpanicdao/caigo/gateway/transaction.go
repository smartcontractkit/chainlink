package gateway

import (
	"context"
	"fmt"
	"math/big"
	"net/http"
	"net/url"
	"time"

	"github.com/dontpanicdao/caigo/types"
	"github.com/google/go-querystring/query"
)

type StarknetTransaction struct {
	TransactionIndex int         `json:"transaction_index"`
	BlockNumber      int         `json:"block_number"`
	Transaction      Transaction `json:"transaction"`
	BlockHash        string      `json:"block_hash"`
	Status           string      `json:"status"`
}

// TODO: test this TX structure matches the case for InvokeV1 and DeployAccount

type Transaction struct {
	TransactionHash    string   `json:"transaction_hash,omitempty"`
	ClassHash          string   `json:"class_hash,omitempty"`
	ContractAddress    string   `json:"contract_address,omitempty"`
	SenderAddress      string   `json:"sender_address,omitempty"`
	EntryPointSelector string   `json:"entry_point_selector,omitempty"`
	Calldata           []string `json:"calldata"`
	Signature          []string `json:"signature"`
	EntryPointType     string   `json:"entry_point_type,omitempty"`
	MaxFee             string   `json:"max_fee,omitempty"`
	Nonce              string   `json:"nonce,omitempty"`
	Version            string   `json:"version,omitempty"`
	Type               string   `json:"type,omitempty"`
}

type TransactionReceipt struct {
	Status                types.TransactionState `json:"status"`
	BlockHash             string                 `json:"block_hash"`
	BlockNumber           int                    `json:"block_number"`
	TransactionIndex      int                    `json:"transaction_index"`
	TransactionHash       string                 `json:"transaction_hash"`
	L1ToL2ConsumedMessage struct {
		FromAddress string   `json:"from_address"`
		ToAddress   string   `json:"to_address"`
		Selector    string   `json:"selector"`
		Payload     []string `json:"payload"`
	} `json:"l1_to_l2_consumed_message"`
	L2ToL1Messages     []interface{}      `json:"l2_to_l1_messages"`
	Events             []interface{}      `json:"events"`
	ExecutionResources ExecutionResources `json:"execution_resources"`
}

type ExecutionResources struct {
	NSteps                 int `json:"n_steps"`
	BuiltinInstanceCounter struct {
		PedersenBuiltin   int `json:"pedersen_builtin"`
		RangeCheckBuiltin int `json:"range_check_builtin"`
		BitwiseBuiltin    int `json:"bitwise_builtin"`
		OutputBuiltin     int `json:"output_builtin"`
		EcdsaBuiltin      int `json:"ecdsa_builtin"`
		EcOpBuiltin       int `json:"ec_op_builtin"`
	} `json:"builtin_instance_counter"`
	NMemoryHoles int `json:"n_memory_holes"`
}

type TransactionReceiptType struct {
	TransactionHash string       `json:"txn_hash,omitempty"`
	Status          string       `json:"status,omitempty"`
	StatusData      string       `json:"status_data,omitempty"`
	MessagesSent    []*L1Message `json:"messages_sent,omitempty"`
	L1OriginMessage *L2Message   `json:"l1_origin_message,omitempty"`
	Events          []*Event     `json:"events,omitempty"`
}

type TransactionOptions struct {
	TransactionId   uint64 `url:"transactionId,omitempty"`
	TransactionHash string `url:"transactionHash,omitempty"`
}

func (gw *Gateway) TransactionByHash(ctx context.Context, hash string) (*Transaction, error) {
	t, err := gw.Transaction(ctx, TransactionOptions{TransactionHash: hash})
	if err != nil {
		return nil, err
	}
	return &t.Transaction, nil
}

// Gets the transaction information from a tx id.
//
// [Reference](https://github.com/starkware-libs/cairo-lang/blob/f464ec4797361b6be8989e36e02ec690e74ef285/src/starkware/starknet/services/api/feeder_gateway/feeder_gateway_client.py#L54-L58)
func (gw *Gateway) Transaction(ctx context.Context, opts TransactionOptions) (*StarknetTransaction, error) {
	req, err := gw.newRequest(ctx, http.MethodGet, "/get_transaction", nil)
	if err != nil {
		return nil, err
	}
	vs, err := query.Values(opts)
	if err != nil {
		return nil, err
	}
	appendQueryValues(req, vs)

	var resp StarknetTransaction
	return &resp, gw.do(req, &resp)
}

type TransactionStatusOptions struct {
	TransactionId   uint64 `url:"transactionId,omitempty"`
	TransactionHash string `url:"transactionHash,omitempty"`
}

type TransactionStatus struct {
	TxStatus        string `json:"tx_status"`
	BlockHash       string `json:"block_hash,omitempty"`
	TxFailureReason struct {
		ErrorMessage string `json:"error_message,omitempty"`
	} `json:"tx_failure_reason,omitempty"`
}

// Gets the transaction status from a txn.
//
// [Reference](https://github.com/starkware-libs/cairo-lang/blob/fc97bdd8322a7df043c87c371634b26c15ed6cee/src/starkware/starknet/services/api/feeder_gateway/feeder_gateway_client.py#L87)
func (gw *Gateway) TransactionStatus(ctx context.Context, opts TransactionStatusOptions) (*TransactionStatus, error) {
	req, err := gw.newRequest(ctx, http.MethodGet, "/get_transaction_status", nil)
	if err != nil {
		return nil, err
	}
	vs, err := query.Values(opts)
	if err != nil {
		return nil, err
	}
	appendQueryValues(req, vs)

	var resp TransactionStatus
	return &resp, gw.do(req, &resp)
}

// Gets the transaction id from its hash.
//
// [Reference](https://github.com/starkware-libs/cairo-lang/blob/fc97bdd8322a7df043c87c371634b26c15ed6cee/src/starkware/starknet/services/api/feeder_gateway/feeder_gateway_client.py#L137)
func (gw *Gateway) TransactionID(ctx context.Context, hash string) (*big.Int, error) {
	req, err := gw.newRequest(ctx, http.MethodGet, "/get_transaction_id_by_hash", nil)
	if err != nil {
		return nil, err
	}

	appendQueryValues(req, url.Values{
		"transactionHash": []string{hash},
	})

	var resp big.Int
	return &resp, gw.do(req, &resp)
}

// Gets the transaction hash from its id.
//
// [Reference](https://github.com/starkware-libs/cairo-lang/blob/fc97bdd8322a7df043c87c371634b26c15ed6cee/src/starkware/starknet/services/api/feeder_gateway/feeder_gateway_client.py#L130)
func (gw *Gateway) TransactionHash(ctx context.Context, id *big.Int) (string, error) {
	req, err := gw.newRequest(ctx, http.MethodGet, "/get_transaction_hash_by_id", nil)
	if err != nil {
		return "", err
	}

	appendQueryValues(req, url.Values{
		"transactionId": []string{id.String()},
	})

	var resp string
	if err := gw.do(req, &resp); err != nil {
		return "", err
	}

	return resp, nil
}

// Get transaction receipt for specific tx
//
// [Reference](https://github.com/starkware-libs/cairo-lang/blob/fc97bdd8322a7df043c87c371634b26c15ed6cee/src/starkware/starknet/services/api/feeder_gateway/feeder_gateway_client.py#L104)
func (gw *Gateway) TransactionReceipt(ctx context.Context, txHash string) (*TransactionReceipt, error) {
	req, err := gw.newRequest(ctx, http.MethodGet, "/get_transaction_receipt", nil)
	if err != nil {
		return nil, err
	}

	appendQueryValues(req, url.Values{
		"transactionHash": []string{txHash},
	})

	var resp TransactionReceipt
	return &resp, gw.do(req, &resp)
}

func (gw *Gateway) TransactionTrace(ctx context.Context, txHash string) (*TransactionTrace, error) {
	req, err := gw.newRequest(ctx, http.MethodGet, "/get_transaction_trace", nil)
	if err != nil {
		return nil, err
	}

	appendQueryValues(req, url.Values{
		"transactionHash": []string{txHash},
	})

	var resp TransactionTrace
	return &resp, gw.do(req, &resp)
}

func (gw *Gateway) WaitForTransaction(ctx context.Context, txHash string, interval, maxPoll int) (n int, receipt *TransactionReceipt, err error) {
	errNotFound := fmt.Errorf("tx not finalized: %s", txHash)
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	count := 0
	for {
		select {
		case <-ticker.C:
			count++
			receipt, err = gw.TransactionReceipt(ctx, txHash)
			if err != nil || receipt.Status.IsTransactionFinal() {
				return count, receipt, err
			}
			if count >= maxPoll {
				return count, receipt, errNotFound
			}
		case <-ctx.Done():
			if count >= maxPoll {
				return count, nil, ctx.Err()
			}
		}
	}
}

type L1Message struct {
	ToAddress string        `json:"to_address,omitempty"`
	Payload   []*types.Felt `json:"payload,omitempty"`
}

type L2Message struct {
	FromAddress string        `json:"from_address,omitempty"`
	Payload     []*types.Felt `json:"payload,omitempty"`
}

type Event struct {
	Order       int           `json:"order,omitempty"`
	FromAddress string        `json:"from_address,omitempty"`
	Keys        []*types.Felt `json:"keys,omitempty"`
	Data        []*types.Felt `json:"data,omitempty"`
}

type TransactionTrace struct {
	FunctionInvocation FunctionInvocation `json:"function_invocation"`
	Signature          []*types.Felt      `json:"signature"`
}

type FunctionInvocation struct {
	CallerAddress      string               `json:"caller_address"`
	ContractAddress    string               `json:"contract_address"`
	Calldata           []string             `json:"calldata"`
	CallType           string               `json:"call_type"`
	ClassHash          string               `json:"class_hash"`
	Selector           string               `json:"selector"`
	EntryPointType     string               `json:"entry_point_type"`
	Result             []string             `json:"result"`
	ExecutionResources ExecutionResources   `json:"execution_resources"`
	InternalCalls      []FunctionInvocation `json:"internal_calls"`
	Events             []Event              `json:"events"`
	Messages           []interface{}        `json:"messages"`
}
