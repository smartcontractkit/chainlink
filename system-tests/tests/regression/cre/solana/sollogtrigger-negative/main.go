//go:build wasip1

package main

import (
	"encoding/hex"
	"fmt"
	"log/slog"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/smartcontractkit/cre-sdk-go/capabilities/blockchain/solana"
	"github.com/smartcontractkit/cre-sdk-go/cre"
	"github.com/smartcontractkit/cre-sdk-go/cre/wasm"

	"github.com/smartcontractkit/chainlink/system-tests/tests/regression/cre/solana/sollogtrigger-negative/config"
)

// Function names, they must literally match the `functionToTest` values used by the
// regression test (system-tests/tests/regression/cre/solana_regression_test.go).
const (
	logTriggerInvalidAddress       = "LogTrigger - invalid address"
	logTriggerEmptyEventName       = "LogTrigger - empty event name"
	logTriggerEmptyFilterName      = "LogTrigger - empty filter name"
	logTriggerEmptyIDL             = "LogTrigger - empty IDL"
	logTriggerConflictingSubkeys   = "LogTrigger - conflicting subkey filters"
	logTriggerInvalidCPIDestAddres = "LogTrigger - invalid CPI destination address"
)

// anchorCPIMethodName is the method name Anchor uses for CPI-emitted events.
const anchorCPIMethodName = "anchor:event"

// testEventIDL is the Anchor IDL of the log_read_test program, copied from the generated bindings
// used by the Solana log trigger smoke test
// (system-tests/tests/smoke/cre/solana/sollogtrigger/contracts/solana/src/generated/log_read_test).
// The negative cases build their filters by hand, so they need the IDL as a plain value: every
// case keeps it valid except the one that asserts an empty IDL is rejected.
//
// It lives in this file, and not in a second file of the package, because workflows are compiled
// with `go build main.go` - a sibling file in package main would not be part of the build.
const testEventIDL = "{\"address\":\"J1zQwrBNBngz26jRPNWsUSZMHJwBwpkoDitXRV95LdK4\",\"metadata\":{\"name\":\"log_read_test\",\"version\":\"0.1.0\",\"spec\":\"0.1.0\",\"description\":\"Created with Anchor\"},\"instructions\":[{\"name\":\"create_log\",\"discriminator\":[215,95,248,114,153,204,208,48],\"accounts\":[{\"name\":\"authority\",\"signer\":true},{\"name\":\"system_program\"}],\"args\":[{\"name\":\"value\",\"type\":\"u64\"}]},{\"name\":\"create_log_cpi\",\"discriminator\":[101,207,236,250,214,98,173,146],\"accounts\":[{\"name\":\"authority\",\"signer\":true},{\"name\":\"system_program\"},{\"name\":\"event_authority\"},{\"name\":\"program\"}],\"args\":[{\"name\":\"value\",\"type\":\"u64\"}]},{\"name\":\"create_truncated_log\",\"discriminator\":[133,74,116,132,80,11,241,64],\"accounts\":[{\"name\":\"authority\",\"signer\":true},{\"name\":\"system_program\"}],\"args\":[{\"name\":\"value\",\"type\":\"u64\"}]}],\"events\":[{\"name\":\"TestEvent\",\"discriminator\":[28,52,39,105,8,210,91,9]}],\"types\":[{\"name\":\"TestEvent\",\"type\":{\"kind\":\"struct\",\"fields\":[{\"name\":\"str_val\",\"type\":\"string\"},{\"name\":\"u64_value\",\"type\":\"u64\"}]}}]}"

// testEventName is the event declared by the IDL above.
const testEventName = "TestEvent"

const validFilterName = "solana-negative-test-event-filter"

func main() {
	wasm.NewRunner(parseConfig).Run(RunSolanaLogTriggerNegativeWorkflow)
}

func parseConfig(b []byte) (config.Config, error) {
	var wfCfg config.Config
	if err := yaml.Unmarshal(b, &wfCfg); err != nil {
		return config.Config{}, fmt.Errorf("error unmarshalling config: %w", err)
	}
	return wfCfg, nil
}

// RunSolanaLogTriggerNegativeWorkflow registers a log trigger whose filter is invalid, so the
// workflow engine is expected to fail while initialising and onLogTrigger is never called.
func RunSolanaLogTriggerNegativeWorkflow(wfCfg config.Config, logger *slog.Logger, secretsProvider cre.SecretsProvider) (cre.Workflow[config.Config], error) {
	logger.Info("RunSolanaLogTriggerNegativeWorkflow called", "functionToTest", wfCfg.FunctionToTest)

	filter, err := buildFilter(wfCfg)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to build log trigger filter for '%s': %v", wfCfg.FunctionToTest, err))
		return nil, fmt.Errorf("failed to build log trigger filter for '%s': %w", wfCfg.FunctionToTest, err)
	}

	logger.Info(fmt.Sprintf(
		"FilterLogTriggerRequest content: Name=%q Address=%s EventName=%q ContractIdlJson=%d bytes Subkeys=%d CpiFilterConfig=%s",
		filter.Name,
		hex.EncodeToString(filter.Address),
		filter.EventName,
		len(filter.ContractIdlJson),
		len(filter.Subkeys),
		formatCPIFilterConfig(filter.CpiFilterConfig),
	))

	return cre.Workflow[config.Config]{
		cre.Handler(
			solana.LogTrigger(wfCfg.ChainSelector, filter),
			onLogTrigger,
		),
	}, nil
}

// buildFilter returns an intentionally invalid log trigger filter for the selected case. Every
// field other than the one under test stays valid, so that the capability rejects the filter for
// the reason the case is about.
func buildFilter(wfCfg config.Config) (*solana.FilterLogTriggerRequest, error) {
	filter := &solana.FilterLogTriggerRequest{
		Name:            validFilterName,
		Address:         wfCfg.ProgramID.Bytes(),
		EventName:       testEventName,
		ContractIdlJson: []byte(testEventIDL),
	}

	switch wfCfg.FunctionToTest {
	case logTriggerInvalidAddress:
		address, err := decodeHex(wfCfg.InvalidInput)
		if err != nil {
			return nil, fmt.Errorf("failed to decode address '%s': %w", wfCfg.InvalidInput, err)
		}
		filter.Address = address

	case logTriggerEmptyEventName:
		filter.EventName = ""

	case logTriggerEmptyFilterName:
		filter.Name = ""

	case logTriggerEmptyIDL:
		filter.ContractIdlJson = nil

	case logTriggerConflictingSubkeys:
		// Two equality filters with different values on the same subkey can never both hold.
		filter.Subkeys = []*solana.SubkeyConfig{
			{
				Path: []string{"u64_value"},
				Comparers: []*solana.ValueComparator{
					{Value: []byte{1}, Operator: solana.ComparisonOperator_COMPARISON_OPERATOR_EQ},
					{Value: []byte{2}, Operator: solana.ComparisonOperator_COMPARISON_OPERATOR_EQ},
				},
			},
		}

	case logTriggerInvalidCPIDestAddres:
		destAddress, err := decodeHex(wfCfg.InvalidInput)
		if err != nil {
			return nil, fmt.Errorf("failed to decode CPI destination address '%s': %w", wfCfg.InvalidInput, err)
		}
		filter.CpiFilterConfig = &solana.CPIFilterConfig{
			DestAddress: destAddress,
			MethodName:  []byte(anchorCPIMethodName),
		}

	default:
		return nil, fmt.Errorf("the provided name for function to test in regression Solana LogTrigger Workflow did not match any known functions: %s", wfCfg.FunctionToTest)
	}

	return filter, nil
}

// onLogTrigger should never be called: registering the trigger is expected to fail first.
func onLogTrigger(wfCfg config.Config, runtime cre.Runtime, payload *solana.Log) (string, error) {
	runtime.Logger().Error("onLogTrigger should not be called - trigger registration should have failed", "functionToTest", wfCfg.FunctionToTest)
	return "", fmt.Errorf("onLogTrigger should not be called - trigger registration should have failed for '%s'", wfCfg.FunctionToTest)
}

func formatCPIFilterConfig(cpiFilterConfig *solana.CPIFilterConfig) string {
	if cpiFilterConfig == nil {
		return "<nil>"
	}
	return fmt.Sprintf("{destAddress: %s, methodName: %s}", hex.EncodeToString(cpiFilterConfig.DestAddress), cpiFilterConfig.MethodName)
}

// decodeHex decodes a hex-encoded byte string of any length. An empty string decodes to no bytes,
// which reaches the capability as a nil field.
func decodeHex(s string) ([]byte, error) {
	return hex.DecodeString(strings.TrimPrefix(s, "0x"))
}
