package hyperliquid

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"

	"github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"

	"github.com/ugorji/go/codec"

	chain_selectors "github.com/smartcontractkit/chain-selectors"
)

var EnableBigBlockChangeset = cldf.CreateChangeSet(enableBigBlocksLogic, enableBigBlocksPreCondition)

type EnableBigBlocksConfig struct {
	APIURL   string
	ChainSel uint64
}

// Payload to be sent in the HTTP POST request
type EnableBigBlocksRequestPayload struct {
	Action    map[string]interface{} `json:"action"`    // Action details
	Nonce     int64                  `json:"nonce"`     // Unique nonce for the request
	Signature ECDSASignature         `json:"signature"` // ECDSA signature of the action
}

type ECDSASignature struct {
	R string `json:"r"`
	S string `json:"s"`
	V byte   `json:"v"`
}

type EnableBigBlocksDetailConfig struct {
	URL               string        // RPC URL
	ChainID           int64         // chain ID
	VerifyingContract string        // Verifying contract address
	RequestTimeout    time.Duration // HTTP request timeout
}

func enableBigBlocksPreCondition(env cldf.Environment, cfg EnableBigBlocksConfig) error {
	_, err := stateview.LoadOnchainState(env)
	if err != nil {
		return fmt.Errorf("failed to load onchain state: %w", err)
	}
	return nil
}

func enableBigBlocksLogic(env cldf.Environment, cfg EnableBigBlocksConfig) (cldf.ChangesetOutput, error) {
	chainIDStr, err := chain_selectors.GetChainIDFromSelector(cfg.ChainSel)

	out := cldf.ChangesetOutput{}
	if err != nil {
		return out, fmt.Errorf("invalid chain id: %w", err)
	}

	if err != nil {
		return out, fmt.Errorf("invalid private key: %w", err)
	}

	chainID, err := strconv.ParseUint(chainIDStr, 10, 64)
	const maxInt64 = int64(^uint64(0) >> 1)
	if chainID > uint64(maxInt64) {
		return out, fmt.Errorf("chain ID too large for int64: %d", chainID)
	}

	if err != nil {
		return out, fmt.Errorf("error converting string to uint64: %w", err)
	}

	action := map[string]interface{}{
		"type":           "evmUserModify",
		"usingBigBlocks": true,
	}

	chain, err := FindChainBySelector(env, cfg.ChainSel)

	if err != nil {
		return out, fmt.Errorf("error: %w finding chain by selector: %d", err, cfg.ChainSel)
	}

	// Verifying contract address for EIP-712 signing
	defaultVerifyingContract := "0x0000000000000000000000000000000000000000"

	// Default timeout for HTTP requests
	defaultRequestTimeout := 10 * time.Second

	nonce := time.Now().UnixMilli()
	config := EnableBigBlocksDetailConfig{
		URL:               cfg.APIURL,
		ChainID:           int64(chainID),
		VerifyingContract: defaultVerifyingContract,
		RequestTimeout:    defaultRequestTimeout,
	}

	sig, err := SignL1Action(action, nonce, true, config, chain)
	if err != nil {
		return out, fmt.Errorf("signing failed: %w", err)
	}

	err = sendRequest(EnableBigBlocksRequestPayload{
		Action:    action,
		Nonce:     nonce,
		Signature: sig,
	}, config)
	if err != nil {
		return out, fmt.Errorf("send failed: %w", err)
	}

	return out, nil
}

func SignL1Action(action map[string]interface{}, nonce int64, isMainnet bool, config EnableBigBlocksDetailConfig, chain chain.BlockChain) (ECDSASignature, error) {
	// Compute the action hash
	actionHash, err := ActionHash(action, nil, nonce)
	if err != nil {
		return ECDSASignature{}, err
	}

	// Construct the phantom agent for signing
	source := "a"
	if !isMainnet {
		source = "b"
	}
	phantomAgent := map[string]interface{}{
		"source":       source,
		"connectionId": actionHash,
	}

	// Define the EIP-712 domain
	domain := apitypes.TypedDataDomain{
		Name:              "Exchange",
		Version:           "1",
		ChainId:           (*math.HexOrDecimal256)(big.NewInt(config.ChainID)),
		VerifyingContract: config.VerifyingContract,
	}

	// Define the EIP-712 types
	types := apitypes.Types{
		"EIP712Domain": {
			{Name: "name", Type: "string"},
			{Name: "version", Type: "string"},
			{Name: "chainId", Type: "uint256"},
			{Name: "verifyingContract", Type: "address"},
		},
		"Agent": {
			{Name: "source", Type: "string"},
			{Name: "connectionId", Type: "bytes32"},
		},
	}

	// Construct the typed data for signing
	typedData := apitypes.TypedData{
		Domain:      domain,
		Types:       types,
		PrimaryType: "Agent",
		Message:     phantomAgent,
	}

	// Compute the hash of the typed data
	hash, _, err := apitypes.TypedDataAndHash(typedData)
	if err != nil {
		return ECDSASignature{}, err
	}

	// Sign the hash using the private key
	// signature, err := crypto.Sign(hash, privateKey)
	evmChain, ok := chain.(evm.Chain)
	if !ok {
		return ECDSASignature{}, errors.New("not an EVM chain")
	}

	signature, err := evmChain.SignHash(hash)
	if err != nil {
		return ECDSASignature{}, err
	}

	var r, s [32]byte
	copy(r[:], signature[:32])
	copy(s[:], signature[32:64])
	v := signature[64] + 27 // Adjust V to be 27 or 28

	return ECDSASignature{
		R: hex.EncodeToString(r[:]),
		S: hex.EncodeToString(s[:]),
		V: v,
	}, nil
}

// sendRequest sends the HTTP POST request with the signed payload
func sendRequest(payload EnableBigBlocksRequestPayload, config EnableBigBlocksDetailConfig) error {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error marshaling payload: %w", err)
	}

	// Create & send the HTTP POST request
	req, err := http.NewRequestWithContext(context.Background(), "POST", config.URL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: config.RequestTimeout}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %w", err)
	}

	fmt.Println("Response Status:", resp.Status)
	fmt.Println("Response Message:", string(resBody))
	return nil
}

// ActionHash computes the hash of the action, including the nonce and optional vault address
func ActionHash(action map[string]interface{}, vaultAddress *string, nonce int64) ([]byte, error) {
	var buf bytes.Buffer

	// Encode the action using MessagePack
	encoder := codec.NewEncoder(&buf, new(codec.MsgpackHandle))
	if err := encoder.Encode(action); err != nil {
		return nil, fmt.Errorf("error encoding action: %w", err)
	}

	// Append the nonce
	nonceBytes := make([]byte, 8)
	if nonce < 0 {
		return nil, fmt.Errorf("nonce must be non-negative, got %d", nonce)
	}
	binary.BigEndian.PutUint64(nonceBytes, uint64(nonce))
	buf.Write(nonceBytes)

	// Append vault address if provided
	if vaultAddress == nil {
		buf.WriteByte(0x00)
	} else {
		buf.WriteByte(0x01)
		addressBytes := hexutil.Bytes(*vaultAddress)
		buf.Write(addressBytes)
	}

	// Compute the Keccak-256 hash of the serialized data
	return crypto.Keccak256(buf.Bytes()), nil
}

func FindChainBySelector(e cldf.Environment, selector uint64) (chain.BlockChain, error) {
	evmChains := e.BlockChains.EVMChains()

	for _, chain := range evmChains {
		if chain.ChainSelector() == selector {
			return chain, nil
		}
	}
	return nil, fmt.Errorf("error finding chain with selector: %d", selector)
}
