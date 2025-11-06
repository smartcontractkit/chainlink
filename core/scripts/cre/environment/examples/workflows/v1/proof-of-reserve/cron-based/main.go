package main

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"time"

	"gopkg.in/yaml.v3"

	readcontractcap "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/readcontract"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/ocr3/aggregators"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/ocr3/ocr3cap"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/targets/chainwriter"
	croncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities/triggers/cron"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/sdk"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm"
	types "github.com/smartcontractkit/chainlink/core/scripts/cre/environment/examples/workflows/v1/proof-of-reserve/cron-based/types"
)

func main() {
	runner := wasm.NewRunner()
	workflow := BuildWorkflow(runner)
	runner.Run(workflow)
}

func convertBigIntToFloat64(bi *big.Int) float64 {
	bigFloat := new(big.Float).SetInt(bi)
	f, _ := bigFloat.Float64()
	return f
}

func BuildWorkflow(runner *wasm.Runner) *sdk.WorkflowSpecFactory {
	workflow := sdk.NewWorkflowSpecFactory()

	cron := croncap.Config{
		Schedule: "*/30 * * * * *", // Every 30 seconds
	}.New(workflow)

	var workflowConfig types.WorkflowConfig
	err := yaml.Unmarshal(runner.Config(), &workflowConfig)
	if err != nil {
		runner.ExitWithError(errors.New("cannot unmarshal config : %w"))
	}

	// Configure Read Contract capability
	chainFamily := workflowConfig.ChainFamily
	chainID := workflowConfig.ChainID
	balanceReaderAddr := workflowConfig.BalanceReaderAddress
	addressesToRead := workflowConfig.BalanceReaderConfig.AddressesToRead
	readcontractCapID := fmt.Sprintf("read-contract-%s-%s@1.0.0", chainFamily, chainID)
	readcontractCapReaderConfig := `{"contracts":{"BalanceReader":{"contractABI":"[{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"addresses\",\"type\":\"address[]\"}],\"name\":\"getNativeBalances\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]","contractPollingFilter":{"genericEventNames":null,"pollingFilter":{"topic2":null,"topic3":null,"topic4":null,"retention":"0s","maxLogsKept":0,"logsPerBlock":0}},"configs":{"getNativeBalances":"{  \"chainSpecificName\": \"getNativeBalances\"}"}}}}`
	readcontractCapReadIdentifier := fmt.Sprintf("%s-%s-%s", balanceReaderAddr, "BalanceReader", "getNativeBalances")
	readcontractCapConfig := readcontractcap.Config{
		ContractReaderConfig: readcontractCapReaderConfig,
		ContractAddress:      balanceReaderAddr,
		ContractName:         "BalanceReader",
		ReadIdentifier:       readcontractCapReadIdentifier,
	}
	readcontractCapRef := "readSmokeTest"
	readcontractCapActionInput := readcontractcap.ActionInput{
		ConfidenceLevel: sdk.ConstantDefinition("unconfirmed"),
		Params: sdk.ConstantDefinition(readcontractcap.InputParams{
			"addresses": addressesToRead,
		}),
		StepDependency: sdk.ConstantDefinition(cron.Ref()),
	}

	chainRead := readcontractCapConfig.New(workflow, readcontractCapID, readcontractCapRef, readcontractCapActionInput)

	if workflowConfig.FeedID == "" {
		runner.ExitWithError(fmt.Errorf("feedID is empty in the config: %+v", workflowConfig))
	}

	computeConfig := types.ComputeConfig{
		FeedID:                workflowConfig.FeedID,
		URL:                   workflowConfig.URL,
		DataFeedsCacheAddress: workflowConfig.DataFeedsCacheAddress,
		WriteTargetName:       workflowConfig.WriteTargetName,
	}

	if workflowConfig.AuthKeySecretName != "" {
		// Secrets are only resolved by Custom Compute, if they are passed as config fields with type `sdk.SecretValue`
		// If we tried to call `sdk.Secret(workflowConfig.AuthKeySecretName)` directly inside the `Compute1WithConfig` function,
		// it would not be resolved and would be passed as a string to the compute function.
		computeConfig.AuthKey = sdk.Secret(workflowConfig.AuthKeySecretName)
	}

	compute := sdk.Compute1WithConfig(
		workflow,
		"compute",
		&sdk.ComputeConfig[types.ComputeConfig]{Config: computeConfig},
		sdk.Compute1Inputs[readcontractcap.Output]{Arg0: chainRead},
		func(runtime sdk.Runtime, config types.ComputeConfig, output readcontractcap.Output) (computeOutput, error) {
			feedID, err := convertFeedIDtoBytes(config.FeedID)
			if err != nil {
				return computeOutput{}, fmt.Errorf("cannot convert feedID to bytes : %w : %b", err, feedID)
			}

			// READ THE BALANCES
			balances, ok := output.LatestValue.([]any)
			if !ok {
				return computeOutput{}, fmt.Errorf("cannot convert latest value to []*big.Int, got type %T", output.LatestValue)
			}
			runtime.Emitter().With("feedID", config.FeedID).Emit(fmt.Sprintf("Balances read, %s", config.FeedID))

			totalBalance := &big.Int{}
			for _, balance := range balances {
				bi, ok := balance.(*big.Int)
				if !ok {
					return computeOutput{}, fmt.Errorf("cannot convert value to *big.Int, got %T", bi)
				}
				totalBalance = totalBalance.Add(totalBalance, bi)
			}
			runtime.Emitter().With("feedID", config.FeedID).Emit(fmt.Sprintf("Total Balances: %s", totalBalance.String()))

			fetchRequest := sdk.FetchRequest{
				URL:       config.URL + "?feedID=" + config.FeedID,
				Method:    "GET",
				TimeoutMs: 5000,
			}

			if string(config.AuthKey) != "" {
				fetchRequest.Headers = map[string]string{
					"Authorization": string(config.AuthKey),
				}
			}

			fresp, err := runtime.Fetch(fetchRequest)
			if err != nil {
				return computeOutput{}, err
			}

			var resp trueUSDResponse
			err = json.Unmarshal(fresp.Body, &resp)
			if err != nil {
				return computeOutput{}, err
			}
			runtime.Emitter().With("feedID", config.FeedID).Emit(fmt.Sprintf("TrueUSD price found: %.2f", resp.TotalTrust))

			if resp.Ripcord {
				runtime.Emitter().With("feedID", config.FeedID).Emit(fmt.Sprintf("ripcord flag set for feed ID %s", config.FeedID))
				return computeOutput{}, sdk.BreakErr
			}

			// COMPUTE THE TOTAL (by adding all the balances)
			runtime.Emitter().With("feedID", config.FeedID).Emit(fmt.Sprintf("Sum '%f' and '%.2f'", resp.TotalTrust, convertBigIntToFloat64(totalBalance)))
			roundedTotalTrust := resp.TotalTrust * 100
			total := roundedTotalTrust + convertBigIntToFloat64(totalBalance) // we multiply by 100 to convert the float to an integer and correctly sum
			runtime.Emitter().With("feedID", config.FeedID).Emit(fmt.Sprintf("Total computed for feed ID %s: %.2f", config.FeedID, total))

			return computeOutput{
				Price:     int(total),
				FeedID:    feedID, // TrueUSD
				Timestamp: resp.UpdatedAt.Unix(),
			}, nil
		},
	)

	consensusInput := ocr3cap.ReduceConsensusInput[computeOutput]{
		Observation: compute.Value(),
	}

	consensus := ocr3cap.ReduceConsensusConfig[computeOutput]{
		Encoder: ocr3cap.EncoderEVM,
		EncoderConfig: map[string]any{
			"abi": "(bytes32 FeedID, uint32 Timestamp, uint224 Price)[] Reports",
		},
		ReportID: "0001",
		KeyID:    "evm",
		AggregationConfig: aggregators.ReduceAggConfig{
			Fields: []aggregators.AggregationField{
				{
					InputKey:  "FeedID",
					OutputKey: "FeedID",
					Method:    "mode",
				},
				{
					InputKey:      "Price",
					OutputKey:     "Price",
					Method:        "median",
					DeviationType: "any",
				},
				{
					InputKey:        "Timestamp",
					OutputKey:       "Timestamp",
					Method:          "median",
					DeviationString: "30",
					DeviationType:   "absolute",
				},
			},
			ReportFormat: aggregators.REPORT_FORMAT_ARRAY,
		},
	}.New(workflow, "consensus", consensusInput)

	targetInput := chainwriter.TargetInput{
		SignedReport: consensus,
	}

	writeTargetName := "write_geth-testnet@1.0.0"
	if workflowConfig.WriteTargetName != "" {
		writeTargetName = workflowConfig.WriteTargetName
	}

	chainwriter.TargetConfig{
		Address:        workflowConfig.DataFeedsCacheAddress, // KeystoneConsumer contract address
		DeltaStage:     "15s",
		Schedule:       "oneAtATime",
		CreStepTimeout: 40,
	}.New(workflow, writeTargetName, targetInput)

	return workflow
}

type trueUSDResponse struct {
	AccountName string    `json:"accountName"`
	TotalTrust  float64   `json:"totalTrust"`
	Ripcord     bool      `json:"ripcord"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type computeOutput struct {
	Price     int
	FeedID    [32]byte
	Timestamp int64
}

func convertFeedIDtoBytes(feedIDStr string) ([32]byte, error) {
	if feedIDStr == "" {
		return [32]byte{}, fmt.Errorf("feedID string is empty")
	}

	if len(feedIDStr) < 2 {
		return [32]byte{}, fmt.Errorf("feedID string too short: %q", feedIDStr)
	}

	b, err := hex.DecodeString(feedIDStr[2:])
	if err != nil {
		return [32]byte{}, err
	}

	if len(b) < 32 {
		nb := [32]byte{}
		copy(nb[:], b[:])
		return nb, err
	}

	return [32]byte(b), nil
}
