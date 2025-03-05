package clientwrappers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/url"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/smartcontractkit/chainlink-integrations/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txm/types"
)

type DualBroadcastClientKeystore interface {
	SignMessage(ctx context.Context, address common.Address, message []byte) ([]byte, error)
}

type DualBroadcastClient struct {
	c         client.Client
	keystore  DualBroadcastClientKeystore
	customURL *url.URL
}

func NewDualBroadcastClient(c client.Client, keystore DualBroadcastClientKeystore, customURL *url.URL) *DualBroadcastClient {
	return &DualBroadcastClient{
		c:         c,
		keystore:  keystore,
		customURL: customURL,
	}
}

func (d *DualBroadcastClient) NonceAt(ctx context.Context, address common.Address, blockNumber *big.Int) (uint64, error) {
	return d.c.NonceAt(ctx, address, blockNumber)
}

func (d *DualBroadcastClient) PendingNonceAt(ctx context.Context, address common.Address) (uint64, error) {
	body := []byte(fmt.Sprintf(`{"jsonrpc":"2.0","method":"eth_getTransactionCount","params":["%s","pending"], "id":1}`, address.String()))
	response, err := d.signAndPostMessage(ctx, address, body, "")
	if err != nil {
		return 0, err
	}

	nonce, err := hexutil.DecodeUint64(response)
	if err != nil {
		return 0, fmt.Errorf("failed to decode response %v into uint64: %w", response, err)
	}
	return nonce, nil
}

func (d *DualBroadcastClient) SendTransaction(ctx context.Context, tx *types.Transaction, attempt *types.Attempt) error {
	meta, err := tx.GetMeta()
	if err != nil {
		return err
	}

	if meta != nil && meta.DualBroadcast != nil && *meta.DualBroadcast && !tx.IsPurgeable {
		data, err := attempt.SignedTransaction.MarshalBinary()
		if err != nil {
			return err
		}
		params := ""
		if meta.DualBroadcastParams != nil {
			params = *meta.DualBroadcastParams
		}
		body := []byte(fmt.Sprintf(`{"jsonrpc":"2.0","method":"eth_sendRawTransaction","params":["%s"], "id":1}`, hexutil.Encode(data)))
		_, err = d.signAndPostMessage(ctx, tx.FromAddress, body, params)
		return err
	}

	return d.c.SendTransaction(ctx, attempt.SignedTransaction)
}

func (d *DualBroadcastClient) signAndPostMessage(ctx context.Context, address common.Address, body []byte, urlParams string) (result string, err error) {
	bodyReader := bytes.NewReader(body)
	postReq, err := http.NewRequestWithContext(ctx, http.MethodPost, d.customURL.String()+"?"+urlParams, bodyReader)
	if err != nil {
		return
	}

	hashedBody := crypto.Keccak256Hash(body).Hex()
	signedMessage, err := d.keystore.SignMessage(ctx, address, []byte(hashedBody))
	if err != nil {
		return
	}

	postReq.Header.Add("X-Flashbots-signature", address.String()+":"+hexutil.Encode(signedMessage))
	postReq.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(postReq)
	if err != nil {
		return result, fmt.Errorf("request %v failed: %w", postReq, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return result, fmt.Errorf("request %v failed with status: %d", postReq, resp.StatusCode)
	}

	keyJSON, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	var response postResponse
	err = json.Unmarshal(keyJSON, &response)
	if err != nil {
		return result, fmt.Errorf("failed to unmarshal response into struct: %w: %s", err, string(keyJSON))
	}
	if response.Error.Message != "" {
		return result, errors.New(response.Error.Message)
	}
	return response.Result, nil
}

type postResponse struct {
	Result string `json:"result,omitempty"`
	Error  postError
}

type postError struct {
	Message string `json:"message,omitempty"`
}
