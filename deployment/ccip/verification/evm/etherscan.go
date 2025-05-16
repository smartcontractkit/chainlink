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
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/verification"
)

// etherscanAPIResponse defines a generic response from etherscan.
type etherscanAPIResponse[R any] struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Result  R      `json:"result"`
}

func (r etherscanAPIResponse[R]) String() string {
	return fmt.Sprintf("etherscanAPIResponse{Status: %s, Message: %s, Result: %+v}", r.Status, r.Message, r.Result)
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
	contractType              cldf.ContractType
	version                   semver.Version
	input                     solidityContractMetadata
	verificationCheckInterval time.Duration
	lggr                      logger.Logger
}

const (
	statusOK  = "1"
	messageOK = "OK"
)

// NewEtherscanContractVerifier creates a new EtherscanContractVerifier instance.
// TODO: Add logging support, use contract.getabi to check if the contract is verified
func NewEtherscanContractVerifier(
	apiURL, apiKey, address string, contractType cldf.ContractType, version semver.Version, verificationCheckInterval time.Duration, lggr logger.Logger) (verification.Verifiable, error) {
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
		lggr:                      lggr,
	}, nil
}

// String returns a string representation of the EtherscanContractVerifier.
func (v *EtherscanContractVerifier) String() string {
	return fmt.Sprintf("%s %s (%s)", v.contractType, v.version, v.address)
}

// Verify verifies the contract on Etherscan.
func (v *EtherscanContractVerifier) Verify(ctx context.Context) error {
	// Check if the contract is already verified
	verified, err := v.IsVerified(ctx)
	if err != nil {
		return fmt.Errorf("failed to check verification status: %w", err)
	}
	if verified {
		v.lggr.Infof("%s is already verified", v)
		return nil
	}

	constructorArgs, err := v.getConstructorArgs(ctx)
	if err != nil {
		return fmt.Errorf("failed to get constructor args: %w", err)
	}
	v.lggr.Infof("Got constructor args for %s: %s", v, constructorArgs)

	sourceCode, err := v.input.SourceCode()
	if err != nil {
		return fmt.Errorf("failed to get source code: %w", err)
	}

	resp, err := sendEtherscanRequest[string](ctx, "POST", v.apiURL, "contract", "verifysourcecode", v.apiKey, map[string]string{
		"contractaddress":       v.address,
		"sourceCode":            sourceCode,
		"codeformat":            "solidity-standard-json-input",
		"contractname":          v.input.Name,
		"compilerversion":       v.input.Version,
		"constructorArguements": constructorArgs, // "Arguements" is a typo in etherscan API
		"constructorArguments":  constructorArgs,
	})
	if err != nil {
		return fmt.Errorf("failed to verify contract: %w", err)
	}
	if resp.Status != statusOK || resp.Message != messageOK {
		return fmt.Errorf("etherscan error - %s", resp)
	}
	v.lggr.Infof("Verifiation request submitted for %s - %s", v, resp)

	for {
		// Check if the contract is verified until context is done OR the contract is verified
		verified, err := v.IsVerified(ctx)
		if err != nil {
			return fmt.Errorf("failed to check verification status: %w", err)
		}
		if verified {
			break
		}
		v.lggr.Infof("Verifiation still pending, checking again in %s", v.verificationCheckInterval)

		select {
		case <-time.After(v.verificationCheckInterval):
			// continue
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return nil
}

// IsVerified checks if the contract is verified on Etherscan.
func (v *EtherscanContractVerifier) IsVerified(ctx context.Context) (bool, error) {
	resp, err := sendEtherscanRequest[string](ctx, "GET", v.apiURL, "contract", "getabi", v.apiKey, map[string]string{
		"address": v.address,
	})
	if err != nil {
		return false, fmt.Errorf("failed to check verification status: %w", err)
	}

	// ABI should be a valid JSON string
	var js interface{}
	return json.Unmarshal([]byte(resp.Result), &js) == nil, nil
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
	form := url.Values{}
	form.Add("module", module)
	form.Add("action", action)
	form.Add("apikey", key)
	for key, value := range extraParams {
		form.Add(key, value)
	}

	httpReq, err := http.NewRequestWithContext(ctx, method, endpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return etherscanAPIResponse[R]{}, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return etherscanAPIResponse[R]{}, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return etherscanAPIResponse[R]{}, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return etherscanAPIResponse[R]{}, fmt.Errorf("http error - status=%d body=%s", resp.StatusCode, string(body))
	}

	var apiResp etherscanAPIResponse[R]
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return etherscanAPIResponse[R]{}, fmt.Errorf("failed to decode response body: %w", err)
	}

	return apiResp, nil
}
