package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"math/big"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/smartcontractkit/chainlink/core/chains/evm"
	"github.com/smartcontractkit/chainlink/core/chains/evm/txmgr"
	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/transmission/generated/entry_point"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/services/pg"
)

var (
	entryPointABI = evmtypes.MustGetABI(entry_point.EntryPointABI)
)

type HttpHandler struct {
	lggr                logger.Logger
	chain               evm.Chain
	fromAddresses       []ethkey.EIP55Address
	q                   pg.Q
	userOpHashToEthTxID map[[32]byte]int64 // TODO: move to local state & add ORM.
	txmORM              txmgr.ORM
}

type UserOperationRequest struct {
	JsonRPCVersion string          `json:"jsonrpc"`
	Id             uint64          `json:"id"`
	Method         string          `json:"method"`
	Params         json.RawMessage `json:"params,omitempty"`
}

type UserOperationResponse struct {
	JsonRPCVersion string `json:"jsonrpc"`
	Id             uint64 `json:"id"`
	Result         string `json:"result"`
}

type GetUserOperationStatusResponse struct {
	JsonRPCVersion string `json:"jsonrpc"`
	Id             uint64 `json:"id"`
	Result         string `json:"result"`
}

type JsonUserOperation struct {
	Sender               string `json:"sender"`
	Nonce                string `json:"nonce"`
	InitCode             string `json:"initCode"`
	CallData             string `json:"callData"`
	CallGasLimit         string `json:"callGasLimit"`
	VerificationGasLimit string `json:"verificationGasLimit"`
	PreVerificationGas   string `json:"preVerificationGas"`
	MaxFeePerGas         string `json:"maxFeePerGas"`
	MaxPriorityFeePerGas string `json:"maxPriorityFeePerGas"`
	PaymasterAndData     string `json:"paymasterAndData"`
	Signature            string `json:"signature"`
}

func NewHandler(
	lggr logger.Logger,
	chain evm.Chain,
	fromAddresses []ethkey.EIP55Address,
	q pg.Q,
	txmORM txmgr.ORM,
) *HttpHandler {
	return &HttpHandler{
		lggr:                lggr,
		chain:               chain,
		fromAddresses:       fromAddresses,
		q:                   q,
		userOpHashToEthTxID: make(map[[32]byte]int64),
		txmORM:              txmORM,
	}
}

func (h *HttpHandler) GetNonce(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Unmarshal request.
	var req UserOperationRequest
	err = json.Unmarshal(b, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Route request.
	h.routeRequest(&req, w, r)
}

func (h *HttpHandler) routeRequest(
	req *UserOperationRequest,
	w http.ResponseWriter,
	r *http.Request,
) {
	if req.Method == "eth_sendUserOperation" {
		h.handleSendUserOperation(req, w, r)
	}

	if req.Method == "eth_getUserOperationStatus" {
		h.handleGetUserOperationByHash(req, w, r)
	}
}

func (h *HttpHandler) handleSendUserOperation(
	req *UserOperationRequest,
	w http.ResponseWriter,
	r *http.Request,
) {
	// Ensure params is an array.
	dec := json.NewDecoder(bytes.NewReader(req.Params))
	tok, err := dec.Token()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if tok != json.Delim('[') {
		http.Error(w, "params not correctly formatted", http.StatusBadRequest)
		return
	}

	// Parse arguments from params.
	var jsonUserOp JsonUserOperation
	var userOp *entry_point.UserOperation
	var entryPointAddress common.Address
	for i := 0; dec.More(); i++ {
		if i == 0 {
			if err := dec.Decode(&jsonUserOp); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
			}
			userOp = JsonUserOperationToUserOperation(&jsonUserOp)
		}
		if i == 1 {
			var addressString string
			if err := dec.Decode(&addressString); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
			}
			entryPointAddress = common.HexToAddress(addressString)
		}
	}
	_, err = dec.Token()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	// Pack user operation into transaction calldata.
	payload, err := entryPointABI.Pack(
		"handleOps",
		[]entry_point.UserOperation{*userOp},
		h.fromAddresses[0].Address(),
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Retrieve userOpHash.
	entryPoint, err := entry_point.NewEntryPoint(entryPointAddress, h.chain.Client())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	userOpHash, err := entryPoint.GetUserOpHash(nil, *userOp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Send transaction via txm.
	var ethTX txmgr.EthTx
	err = h.q.Transaction(func(tx pg.Queryer) error {
		t, err := h.chain.TxManager().CreateEthTransaction(txmgr.NewTx{
			FromAddress:    h.fromAddresses[0].Address(),
			ToAddress:      entryPointAddress,
			EncodedPayload: hexutil.MustDecode(hexutil.Encode(payload)),
			GasLimit:       10_000_000,
			Strategy:       txmgr.NewSendEveryStrategy(),
		}, pg.WithQueryer(tx), pg.WithParentCtx(context.Background()))
		ethTX = t
		return err
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	h.userOpHashToEthTxID[userOpHash] = ethTX.ID

	// Send response to user.
	serializedResp, err := json.Marshal(UserOperationResponse{
		Id:             req.Id,
		JsonRPCVersion: req.JsonRPCVersion,
		Result:         hexutil.Encode(userOpHash[:]),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(serializedResp)
}

func (h *HttpHandler) handleGetUserOperationByHash(
	req *UserOperationRequest,
	w http.ResponseWriter,
	r *http.Request,
) {
	// Ensure params is an array.
	dec := json.NewDecoder(bytes.NewReader(req.Params))
	tok, err := dec.Token()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if tok != json.Delim('[') {
		http.Error(w, "params not correctly formatted", http.StatusBadRequest)
		return
	}

	// Parse arguments from params.
	var userOpHash [32]byte
	for i := 0; dec.More(); i++ {
		var userOpHashString string
		if i == 0 {
			if err := dec.Decode(&userOpHashString); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
			}
			userOpHash = common.HexToHash(userOpHashString)
		}
	}
	_, err = dec.Token()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	// var txID uint64 = h.userOpHashToEthTxID[userOpHash]
	ethTxs, _, err := h.txmORM.EthTransactions(0, 100)
	for _, e := range ethTxs {
		if e.ID == h.userOpHashToEthTxID[userOpHash] {
			// Send response to user.
			serializedResp, err := json.Marshal(GetUserOperationStatusResponse{
				Id:             req.Id,
				JsonRPCVersion: req.JsonRPCVersion,
				Result:         string(e.State),
			})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Write(serializedResp)
			return
		}
	}
	http.Error(w, "could not find request for hash", http.StatusInternalServerError)
}

func JsonUserOperationToUserOperation(jsonUserOp *JsonUserOperation) *entry_point.UserOperation {
	nonce, _ := big.NewInt(0).SetString(jsonUserOp.Nonce, 10)
	callGasLimit, _ := big.NewInt(0).SetString(jsonUserOp.CallGasLimit, 10)
	verificationGasLimit, _ := big.NewInt(0).SetString(jsonUserOp.VerificationGasLimit, 10)
	preVerificationGas, _ := big.NewInt(0).SetString(jsonUserOp.PreVerificationGas, 10)
	maxFeePerGas, _ := big.NewInt(0).SetString(jsonUserOp.MaxFeePerGas, 10)
	maxPriorityFeePerGas, _ := big.NewInt(0).SetString(jsonUserOp.MaxPriorityFeePerGas, 10)

	return &entry_point.UserOperation{
		Sender:               common.HexToAddress(jsonUserOp.Sender),
		Nonce:                nonce,
		InitCode:             common.Hex2Bytes(jsonUserOp.InitCode),
		CallData:             common.Hex2Bytes(jsonUserOp.CallData),
		CallGasLimit:         callGasLimit,
		VerificationGasLimit: verificationGasLimit,
		PreVerificationGas:   preVerificationGas,
		MaxFeePerGas:         maxFeePerGas,
		MaxPriorityFeePerGas: maxPriorityFeePerGas,
		PaymasterAndData:     common.Hex2Bytes(jsonUserOp.PaymasterAndData),
		Signature:            common.Hex2Bytes(jsonUserOp.Signature),
	}
}
