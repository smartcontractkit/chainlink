package evm

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	dfprocessor "github.com/smartcontractkit/chainlink-evm/pkg/report/datafeeds/processor"
	porprocessor "github.com/smartcontractkit/chainlink-evm/pkg/report/por/processor"
	df "github.com/smartcontractkit/chainlink-framework/capabilities/writetarget/monitoring/pb/data-feeds/on-chain/registry"
	processor "github.com/smartcontractkit/chainlink-framework/capabilities/writetarget/report/platform/processor"

	"github.com/smartcontractkit/chainlink-framework/capabilities/writetarget"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	ocr3types "github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/ocr3/types"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	forwarder "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/forwarder_1_0_0"
	"github.com/smartcontractkit/chainlink-evm/pkg/chains/legacyevm"
	relayevmtypes "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

func NewWriteTarget(ctx context.Context, relayer *Relayer, chain legacyevm.Chain, gasLimitDefault uint64, lggr logger.Logger) (capabilities.ExecutableCapability, error) {
	// generate ID based on chain selector
	id := GenerateWriteTargetName(chain.ID().Uint64())

	// EVM-specific init
	config := chain.Config().EVM().Workflow()

	// Initialize a reader to check whether a value was already transmitted on chain
	contractReaderConfigEncoded, err := json.Marshal(relayevmtypes.ChainReaderConfig{
		Contracts: map[string]relayevmtypes.ChainContractReader{
			"forwarder": {
				ContractABI: forwarder.KeystoneForwarderABI,
				Configs: map[string]*relayevmtypes.ChainReaderDefinition{
					"getTransmissionInfo": {
						ChainSpecificName: "getTransmissionInfo",
					},
				},
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal contract reader config %w", err)
	}
	cr, err := relayer.NewContractReader(ctx, contractReaderConfigEncoded)
	if err != nil {
		return nil, err
	}

	chainWriterConfig := relayevmtypes.ChainWriterConfig{
		Contracts: map[string]*relayevmtypes.ContractConfig{
			"forwarder": {
				ContractABI: forwarder.KeystoneForwarderABI,
				Configs: map[string]*relayevmtypes.ChainWriterDefinition{
					"report": {
						ChainSpecificName: "report",
						FromAddress:       config.FromAddress().Address(),
						GasLimit:          gasLimitDefault,
					},
				},
			},
		},
	}
	chainWriterConfig.MaxGasPrice = chain.Config().EVM().GasEstimator().PriceMax()

	encodedWriterConfig, err := json.Marshal(chainWriterConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal chainwriter config: %w", err)
	}

	cw, err := relayer.NewContractWriter(ctx, encodedWriterConfig)
	if err != nil {
		return nil, err
	}

	chainInfo, err := relayer.GetChainInfo(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get chain info: %w", err)
	}

	registryMetrics, err := df.NewMetrics()
	if err != nil {
		return nil, fmt.Errorf("failed to create new registry metrics: %w", err)
	}

	emitter := writetarget.NewMonitorEmitter(lggr)

	dfProcessor := dfprocessor.NewDataFeedsProcessor(registryMetrics, emitter)
	ccipDfProcessor := dfprocessor.NewCCIPDataFeedsProcessor(registryMetrics, emitter)
	porProcessor := porprocessor.NewPORProcessor(registryMetrics, emitter)

	processors, err := processor.NewPlatformProcessors(emitter)
	if err != nil {
		return nil, fmt.Errorf("failed to create EVM platform processors: %w", err)
	}

	processors["evm-data-feeds"] = dfProcessor
	processors["evm-data-feeds-ccip"] = ccipDfProcessor
	processors["evm-por-feeds"] = porProcessor

	beholder, err := writetarget.NewMonitor(writetarget.MonitorOpts{
		Lggr:              lggr,
		Processors:        processors,
		EnabledProcessors: processor.PlatformDefaultProcessors,
		Emitter:           emitter,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Aptos WT monitor client: %+w", err)
	}
	ts, err := NewEVMTargetStrategy(cr, cw, relayer.chain.TxManager(), config.ForwarderAddress().String(), gasLimitDefault, lggr)
	if err != nil {
		return nil, fmt.Errorf("failed to create target strategy: %w", err)
	}
	opts := writetarget.WriteTargetOpts{
		ID:     id,
		Logger: lggr,
		Config: writetarget.Config{
			PollPeriod:        config.PollPeriod(),
			AcceptanceTimeout: config.AcceptanceTimeout(),
		},
		ChainInfo:            chainInfo,
		Beholder:             beholder,
		ChainService:         chain,
		ConfigValidateFn:     evaluate,
		NodeAddress:          config.FromAddress().String(),
		ForwarderAddress:     config.ForwarderAddress().String(),
		TargetStrategy:       ts,
		WriteAcceptanceState: *config.TxAcceptanceState(),
	}

	return writetarget.NewWriteTarget(opts), nil
}

type Inputs struct {
	SignedReport ocr3types.SignedReport
}

type TargetRequest struct {
	Metadata capabilities.RequestMetadata
	Config   Config
	Inputs   Inputs
}

func getEVMRequest(rawRequest capabilities.CapabilityRequest) (
	TargetRequest, error) {
	var r TargetRequest
	r.Metadata = rawRequest.Metadata

	if rawRequest.Config == nil {
		return TargetRequest{}, errors.New("missing config field")
	}

	if err := rawRequest.Config.UnwrapTo(&r.Config); err != nil {
		return TargetRequest{}, err
	}

	if !common.IsHexAddress(r.Config.Address) {
		return TargetRequest{}, fmt.Errorf("'%v' is not a valid address", r.Config.Address)
	}

	if rawRequest.Inputs == nil {
		return TargetRequest{}, errors.New("missing inputs field")
	}

	// required field of target's config in the workflow spec
	signedReport, ok := rawRequest.Inputs.Underlying[writetarget.KeySignedReport]
	if !ok {
		return TargetRequest{}, fmt.Errorf("missing required field %s", writetarget.KeySignedReport)
	}

	if err := signedReport.UnwrapTo(&r.Inputs.SignedReport); err != nil {
		return TargetRequest{}, err
	}
	return r, nil
}

func evaluate(rawRequest capabilities.CapabilityRequest) (receiver string, err error) {
	r, err := getEVMRequest(rawRequest)
	if err != nil {
		return "", err
	}

	// don't need tail in this case
	reportMetadata, _, err := ocr3types.Decode(r.Inputs.SignedReport.Report)
	if err != nil {
		return "", fmt.Errorf("failed to decode report metadata: %w", err)
	}

	if reportMetadata.Version != 1 {
		return "", fmt.Errorf("unsupported report version: %d", reportMetadata.Version)
	}

	if reportMetadata.ExecutionID != rawRequest.Metadata.WorkflowExecutionID {
		return "", fmt.Errorf("WorkflowExecutionID in the report does not match WorkflowExecutionID in the request metadata. Report WorkflowExecutionID: %+v, request WorkflowExecutionID: %+v", hex.EncodeToString([]byte(reportMetadata.ExecutionID)), rawRequest.Metadata.WorkflowExecutionID)
	}

	// case-insensitive verification of the owner address (so that a check-summed address matches its non-checksummed version).
	if !strings.EqualFold(reportMetadata.WorkflowOwner, rawRequest.Metadata.WorkflowOwner) {
		return "", fmt.Errorf("WorkflowOwner in the report does not match WorkflowOwner in the request metadata. Report WorkflowOwner: %+v, request WorkflowOwner: %+v", reportMetadata.WorkflowOwner, rawRequest.Metadata.WorkflowOwner)
	}

	// pad workflow name to match the report which is padded to 20 characters
	if len(rawRequest.Metadata.WorkflowName) < 20 {
		suffix := strings.Repeat("0", 20-len(rawRequest.Metadata.WorkflowName))
		rawRequest.Metadata.WorkflowName += suffix
	}

	if !strings.EqualFold(reportMetadata.WorkflowName, rawRequest.Metadata.WorkflowName) {
		return "", fmt.Errorf("WorkflowName in the report does not match WorkflowName in the request metadata. Report WorkflowName: %+v, request WorkflowName: %+v", reportMetadata.WorkflowName, rawRequest.Metadata.WorkflowName)
	}

	if reportMetadata.WorkflowID != rawRequest.Metadata.WorkflowID {
		return "", fmt.Errorf("WorkflowID in the report does not match WorkflowID in the request metadata. Report WorkflowID: %+v, request WorkflowID: %+v", reportMetadata.WorkflowID, rawRequest.Metadata.WorkflowID)
	}

	byteID, err := hex.DecodeString(reportMetadata.ReportID)
	if err != nil {
		return "", fmt.Errorf("failed to decode report ID: %w", err)
	}

	if !bytes.Equal(byteID, r.Inputs.SignedReport.ID) {
		return "", fmt.Errorf("ReportID in the report does not match ReportID in the inputs. reportMetadata.ReportID: %x, Inputs.SignedReport.ID: %x", reportMetadata.ReportID, r.Inputs.SignedReport.ID)
	}

	return r.Config.Address, nil
}

func GenerateWriteTargetName(chainID uint64) string {
	id := fmt.Sprintf("write_%v@1.0.0", chainID)

	chainName, err := chainselectors.NameFromChainId(chainID)
	if err == nil {
		wtID, err := writetarget.NewWriteTargetID("", chainName, strconv.FormatUint(chainID, 10), "1.0.0")
		if err == nil {
			id = wtID
		}
	}

	return id
}
