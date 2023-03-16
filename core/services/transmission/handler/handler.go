package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/core/logger"
)

type GetNonceRequest struct {
	SenderAddress common.Address `json:"sender_address"`
	SourceChainID uint64         `json:"source_chain_id"`
	TokenAddress  common.Address `json:"token_address"`
}

type GetNonceResponse struct {
	Nonce uint64 `json:"nonce"`
}

type HttpHandler struct {
	lggr logger.Logger
}

func NewHandler(lggr logger.Logger) *HttpHandler {
	return &HttpHandler{
		lggr: lggr,
	}
}

func (h *HttpHandler) GetNonce(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var req GetNonceRequest
	err = json.Unmarshal(b, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// stub response for now
	resp := GetNonceResponse{
		Nonce: 5,
	}
	serializedResp, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(serializedResp)
}
