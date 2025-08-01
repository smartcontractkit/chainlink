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
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
	"github.com/vmihailenco/msgpack/v5"

	"github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
)

var EnableBigBlockChangeset = cldf.CreateChangeSet(enableBigBlocksLogic, enableBigBlocksPreCondition)

type EnableBigBlocksConfig struct {
	APIURL    string
	ChainSel  uint64
	IsMainnet bool
}

// Payload to be sent in the HTTP POST request
type enableBigBlocksRequestPayload struct {
	Action       map[string]interface{} `json:"action"`       // Action details
	Nonce        int64                  `json:"nonce"`        // Unique nonce for the request
	Signature    ecdsaSignature         `json:"signature"`    // ECDSA signature of the action
	VaultAddress *string                `json:"vaultAddress"` // Vault address (null for most cases)
}

type ecdsaSignature struct {
	R string `json:"r"`
	S string `json:"s"`
	V int    `json:"v"`
}

type enableBigBlocksDetailConfig struct {
	URL               string        // RPC URL
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
	out := cldf.ChangesetOutput{}

	action := map[string]interface{}{
		"type":           "evmUserModify",
		"usingBigBlocks": true,
	}

	chain, err := findChainBySelector(env, cfg.ChainSel)

	if err != nil {
		return out, fmt.Errorf("error: %w finding chain by selector: %d", err, cfg.ChainSel)
	}

	// Verifying contract address for EIP-712 signing
	defaultVerifyingContract := "0x0000000000000000000000000000000000000000"

	// Default timeout for HTTP requests
	defaultRequestTimeout := 10 * time.Second

	nonce := time.Now().UnixMilli()
	config := enableBigBlocksDetailConfig{
		URL:               cfg.APIURL,
		VerifyingContract: defaultVerifyingContract,
		RequestTimeout:    defaultRequestTimeout,
	}

	sig, err := signL1Action(action, nonce, cfg.IsMainnet, config, chain)
	if err != nil {
		return out, fmt.Errorf("signing failed: %w", err)
	}

	err = sendRequest(enableBigBlocksRequestPayload{
		Action:       action,
		Nonce:        nonce,
		Signature:    sig,
		VaultAddress: nil,
	}, config)
	if err != nil {
		return out, fmt.Errorf("send failed: %w", err)
	}

	return out, nil
}

func signL1Action(action map[string]interface{}, nonce int64, isMainnet bool, config enableBigBlocksDetailConfig, chain chain.BlockChain) (ecdsaSignature, error) {
	// Compute the action hash
	actionHash, err := actionHash(action, "", nonce, nil)
	if err != nil {
		return ecdsaSignature{}, err
	}

	// Construct the phantom agent for signing
	source := "a"
	if !isMainnet {
		source = "b"
	}
	phantomAgent := map[string]interface{}{
		"source":       source,
		"connectionId": "0x" + hex.EncodeToString(actionHash),
	}

	// Define the EIP-712 domain
	domain := apitypes.TypedDataDomain{
		Name:              "Exchange",
		Version:           "1",
		ChainId:           (*math.HexOrDecimal256)(big.NewInt(1337)),
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
	domainSeparator, err := typedData.HashStruct("EIP712Domain", typedData.Domain.Map())
	if err != nil {
		return ecdsaSignature{}, err
	}
	typedDataHash, err := typedData.HashStruct(typedData.PrimaryType, typedData.Message)
	if err != nil {
		return ecdsaSignature{}, err
	}
	rawData := []byte{0x19, 0x01}
	rawData = append(rawData, domainSeparator...)
	rawData = append(rawData, typedDataHash...)
	msgHash := crypto.Keccak256Hash(rawData)

	// Sign the hash using the private key
	// signature, err := crypto.Sign(hash, privateKey)
	evmChain, ok := chain.(evm.Chain)
	if !ok {
		return ecdsaSignature{}, errors.New("not an EVM chain")
	}

	signature, err := evmChain.SignHash(msgHash.Bytes())
	if err != nil {
		return ecdsaSignature{}, err
	}

	// Extract r, s, v components
	r := new(big.Int).SetBytes(signature[:32])
	s := new(big.Int).SetBytes(signature[32:64])
	v := int(signature[64]) + 27

	return ecdsaSignature{
		R: hexutil.EncodeBig(r),
		S: hexutil.EncodeBig(s),
		V: v,
	}, nil
}

// sendRequest sends the HTTP POST request with the signed payload
func sendRequest(payload enableBigBlocksRequestPayload, config enableBigBlocksDetailConfig) error {
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
func actionHash(action any, vaultAddress string, nonce int64, expiresAfter *int64) ([]byte, error) {
	// Pack action using msgpack (like Python's msgpack.packb)
	var buf bytes.Buffer
	enc := msgpack.NewEncoder(&buf)
	enc.SetSortMapKeys(true)

	err := enc.Encode(action)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal action: %w", err)
	}
	data := buf.Bytes()

	// Add nonce as 8 bytes big endian
	if nonce < 0 {
		return nil, fmt.Errorf("nonce cannot be negative: %d", nonce)
	}
	nonceBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(nonceBytes, uint64(nonce))
	data = append(data, nonceBytes...)

	// Add vault address
	if vaultAddress == "" {
		data = append(data, 0x00)
	} else {
		data = append(data, 0x01)
		data = append(data, addressToBytes(vaultAddress)...)
	}

	// Add expires_after if provided
	if expiresAfter != nil {
		if *expiresAfter < 0 {
			panic(fmt.Sprintf("expiresAfter cannot be negative: %d", *expiresAfter))
		}
		data = append(data, 0x00)
		expiresAfterBytes := make([]byte, 8)
		binary.BigEndian.PutUint64(expiresAfterBytes, uint64(*expiresAfter))
		data = append(data, expiresAfterBytes...)
	}

	// Return keccak256 hash
	hash := crypto.Keccak256(data)
	fmt.Printf("go action hash: %s\n", hex.EncodeToString(hash))
	return hash, nil
}

// addressToBytes converts a hex address to bytes
func addressToBytes(address string) []byte {
	address = strings.TrimPrefix(address, "0x")
	bytes, _ := hex.DecodeString(address)
	return bytes
}

func findChainBySelector(e cldf.Environment, selector uint64) (chain.BlockChain, error) {
	evmChains := e.BlockChains.EVMChains()

	for _, chain := range evmChains {
		if chain.ChainSelector() == selector {
			return chain, nil
		}
	}
	return nil, fmt.Errorf("error finding chain with selector: %d", selector)
}
