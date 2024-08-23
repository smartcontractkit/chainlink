package chain

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers"
	relayevmtypes "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

const (
	MethodWrite = "write"
	MethodRead  = "read"

	contractName           = "my_contract"
	contractFunctionToCall = "increment"
)

type WriteRequestPayload struct {
	ChainFamily string `json:"chain_family`
	ChainID     string `json:"chain_id"`
	ToAddress   string `json:"to_address"`

	Abi          string `json:"abi"`
	FunctionName string `json:"function_name"`
}

type Response struct {
	Success      bool   `json:"success"`
	ErrorMessage string `json:"error_message,omitempty"`
}

type chainHandler struct {
	relayGetter handlers.RelayGetter
	mu          sync.Mutex
	lggr        logger.Logger
}

var _ handlers.Handler = (*chainHandler)(nil)

func NewChainHandler(relayGetter handlers.RelayGetter, lggr logger.Logger) (*chainHandler, error) {
	return &chainHandler{
		relayGetter: relayGetter,
		lggr:        lggr.Named("ChainHandler"),
	}, nil
}

func (h *chainHandler) HandleUserMessage(ctx context.Context, msg *api.Message, callbackCh chan<- handlers.UserCallbackPayload) error {
	if msg.Body.Method == MethodWrite {
		var req WriteRequestPayload
		err := json.Unmarshal(msg.Body.Payload, &req)
		msg.Body.Payload = nil
		if err != nil {
			return h.sendResponse(callbackCh, msg, Response{Success: false, ErrorMessage: "failed to parse payload: " + err.Error()})
		}

		h.mu.Lock()
		defer h.mu.Unlock()
		relayer, err := h.relayGetter.Get(types.RelayID{Network: req.ChainFamily, ChainID: req.ChainID})
		if err != nil {
			return h.sendResponse(callbackCh, msg, Response{Success: false, ErrorMessage: "chain not found: " + err.Error()})
		}

		// create a new chain writer
		chainWriterConfig := relayevmtypes.ChainWriterConfig{
			Contracts: map[string]*relayevmtypes.ContractConfig{
				contractName: {
					// TODO: extract function name and ABI from the request
					ContractABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"increment\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
					Configs: map[string]*relayevmtypes.ChainWriterDefinition{
						contractFunctionToCall: { // TODO: extract contract function name from the request
							ChainSpecificName: contractFunctionToCall,
							Checker:           "simulate",
							FromAddress:       common.HexToAddress("0xe82b221f8E6a4916D7FC057C18Fc36de4b714ee2"), // TODO: extract from somewhere?
							GasLimit:          100000,
						},
					},
				},
			},
		}
		chainWriterConfig.MaxGasPrice = assets.NewWei(big.NewInt(100 * 1000000000))

		encodedWriterConfig, err := json.Marshal(chainWriterConfig)
		if err != nil {
			return h.sendResponse(callbackCh, msg, Response{Success: false, ErrorMessage: "failed to encode chainwriter config: " + err.Error()})
		}

		h.lggr.Info("creating chain writer ...")
		cw, err := relayer.NewChainWriter(ctx, encodedWriterConfig)
		if err != nil {
			return h.sendResponse(callbackCh, msg, Response{Success: false, ErrorMessage: "failed to create chain writer: " + err.Error()})
		}
		h.lggr.Info("created chain writer")
		err = cw.SubmitTransaction(ctx, contractName, contractFunctionToCall, nil, msg.Body.MessageId, req.ToAddress, nil, nil)
		if err != nil {
			return h.sendResponse(callbackCh, msg, Response{Success: false, ErrorMessage: err.Error()})
		}
		return h.sendResponse(callbackCh, msg, Response{Success: true})
	} else {
		return h.sendResponse(callbackCh, msg, Response{Success: false, ErrorMessage: "unsupported method"})
	}
}

func (h *chainHandler) sendResponse(callbackCh chan<- handlers.UserCallbackPayload, msg *api.Message, response Response) error {
	defer close(callbackCh)
	rawResp, err := json.Marshal(response)
	if err != nil {
		return err
	}
	msg.Body.Payload = rawResp
	callbackCh <- handlers.UserCallbackPayload{Msg: msg, ErrCode: api.NoError, ErrMsg: ""}
	return nil
}

func (h *chainHandler) HandleNodeMessage(ctx context.Context, msg *api.Message, nodeAddr string) error {
	return fmt.Errorf("this handler does not expect any node messages")
}

func (h *chainHandler) Start(context.Context) error {
	return nil
}

func (h *chainHandler) Close() error {
	return nil
}
