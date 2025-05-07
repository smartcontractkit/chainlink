package evm

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/verification"
)

// etherscanAPIResponse defines a generic response from etherscan.
type etherscanAPIResponse[R any] struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Result  R      `json:"result"`
}

// transactionInfo details information about a transaction.
type transactionInfo struct {
	Input string `json:"input"`
}

// EtherscanContractVerifier verifies a contract on Etherscan.
type EtherscanContractVerifier struct {
	apiURL                    string
	apiKey                    string
	address                   string
	contractType              deployment.ContractType
	version                   semver.Version
	input                     solidityContractMetadata
	verificationCheckInterval time.Duration
}

const (
	statusOK       = "1"
	messageOK      = "OK"
	resultVerified = "Pass - Verified"
)

// NewEtherscanContractVerifier creates a new EtherscanContractVerifier instance.
func NewEtherscanContractVerifier(
	apiURL, apiKey, address string, contractType deployment.ContractType, version semver.Version, verificationCheckInterval time.Duration) (verification.Verifiable, error) {
	input, err := loadSolidityContractMetadata(contractType, version)
	if err != nil {
		return nil, fmt.Errorf("failed to load solidity standard JSON input: %w", err)
	}

	return &EtherscanContractVerifier{
		apiURL:                    apiURL,
		apiKey:                    apiKey,
		address:                   address,
		contractType:              contractType,
		version:                   version,
		input:                     input,
		verificationCheckInterval: verificationCheckInterval,
	}, nil
}

// String returns a string representation of the EtherscanContractVerifier.
func (v *EtherscanContractVerifier) String() string {
	return fmt.Sprintf("%s %s (%s)", v.contractType, v.version, v.address)
}

// Verify verifies the contract on Etherscan.
func (v *EtherscanContractVerifier) Verify(ctx context.Context) (bool, error) {
	constructorArgs, err := v.getConstructorArgs(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to get constructor args: %w", err)
	}
	fmt.Println("Constructor Args: ", constructorArgs)

	sourceCode, err := v.input.SourceCode()
	if err != nil {
		return false, fmt.Errorf("failed to get source code: %w", err)
	}

	version := v.input.Version
	if !strings.HasPrefix(version, "v") {
		version = "v" + version
	}

	resp, err := sendEtherscanRequest[string](ctx, "POST", v.apiURL, "contract", "verifysourcecode", v.apiKey, map[string]string{
		"contractaddress":       v.address,
		"sourceCode":            sourceCode,
		"codeformat":            "solidity-standard-json-input",
		"contractname":          v.input.Name,
		"compilerversion":       version,
		"constructorArguements": constructorArgs, // "Arguements" is a typo in etherscan API
	})
	if err != nil {
		return false, fmt.Errorf("failed to verify contract: %w", err)
	}
	guid := resp.Result

	for {
		// Check if the contract is verified until context is done OR the contract is verified
		verified, err := v.isVerified(ctx, guid)
		if err != nil {
			return false, fmt.Errorf("failed to check verification status: %w", err)
		}
		if verified {
			break
		}
		// Wait for configured time before checking again
		time.Sleep(v.verificationCheckInterval)

		select {
		case <-ctx.Done():
			return false, ctx.Err()
		}
	}

	return true, nil
}

// isVerified checks if the contract is verified on Etherscan.
func (v *EtherscanContractVerifier) isVerified(ctx context.Context, guid string) (bool, error) {
	resp, err := sendEtherscanRequest[string](ctx, "GET", v.apiURL, "contract", "checkverifystatus", v.apiKey, map[string]string{
		"guid": guid,
	})
	if err != nil {
		return false, fmt.Errorf("failed to check verification status: %w", err)
	}

	return resp.Result == resultVerified, nil
}

// getConstructorArgs returns the construction arguments used to deploy the contract.
func (v *EtherscanContractVerifier) getConstructorArgs(ctx context.Context) (string, error) {
	// We request the first page with an offset of 1 because we only care about the contract creation transaction.
	resp, err := sendEtherscanRequest[[]transactionInfo](ctx, "GET", v.apiURL, "account", "txlist", v.apiKey, map[string]string{
		"address": v.address,
		"page":    "1",
		"offset":  "1",
		"sort":    "asc",
	})
	if err != nil {
		return "", fmt.Errorf("failed to get contract creation info: %w", err)
	}

	l := len(resp.Result)
	if l != 1 {
		return "", fmt.Errorf("expected 1 result, got %d", l)
	}
	tx := resp.Result[0]

	// Trim to remove "0x" as a variable prefix
	bytecode := strings.TrimPrefix(v.input.Bytecode, "0x")
	txInput := strings.TrimPrefix(tx.Input, "0x")

	// Check if transaction input includes the contract bytecode
	if !strings.HasPrefix(txInput, bytecode) {
		return "", errors.New("contract creation tx input does not contain contract bytecode")
	}

	return txInput[len(bytecode):], nil
}

// sendEtherscanRequest sends a request to the Etherscan API and returns the response.
func sendEtherscanRequest[R any](ctx context.Context, method, endpoint, module, action, key string, extraParams map[string]string) (etherscanAPIResponse[R], error) {
	baseURL, err := url.Parse(endpoint)
	if err != nil {
		return etherscanAPIResponse[R]{}, fmt.Errorf("invalid endpoint URL: %w", err)
	}
	query := url.Values{}
	query.Set("module", module)
	query.Set("action", action)
	query.Set("apikey", key)
	for key, value := range extraParams {
		query.Set(key, value)
	}
	baseURL.RawQuery = query.Encode()

	httpReq, err := http.NewRequestWithContext(ctx, method, baseURL.String(), nil)
	if err != nil {
		return etherscanAPIResponse[R]{}, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return etherscanAPIResponse[R]{}, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return etherscanAPIResponse[R]{}, fmt.Errorf("received non-200 response from etherscan: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return etherscanAPIResponse[R]{}, fmt.Errorf("failed to read response body: %w", err)
	}

	var apiResp etherscanAPIResponse[R]
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return etherscanAPIResponse[R]{}, fmt.Errorf("failed to decode response body: %w", err)
	}

	if apiResp.Status != statusOK || apiResp.Message != messageOK {
		return etherscanAPIResponse[R]{}, fmt.Errorf("etherscan error - status=%s message=%s result=%+v", apiResp.Status, apiResp.Message, apiResp.Result)
	}

	return apiResp, nil
}
