package evm

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-evm/pkg/report/monitor"
	"github.com/smartcontractkit/chainlink-evm/pkg/writetarget"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/ocr3/types"
	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	forwarder "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/forwarder_1_0_0"
	"github.com/smartcontractkit/chainlink/v2/core/chains/legacyevm"
	relayevmtypes "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

func NewWriteTarget(ctx context.Context, relayer *Relayer, chain legacyevm.Chain, gasLimitDefault uint64, lggr logger.Logger) (capabilities.TargetCapability, error) {
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

	pollPeriod, err := commonconfig.NewDuration(2 * time.Second)
	if err != nil {
		return nil, fmt.Errorf("failed to create poll period: %w", err)
	}

	// TODO: Unsure what the timeout should be, I don't see an EVM config that corresponds to this.
	timeout, err := commonconfig.NewDuration(30 * time.Second)
	if err != nil {
		return nil, fmt.Errorf("failed to create timeout: %w", err)
	}

	chainInfo, err := getChainInfo(chain.ID().Uint64())
	if err != nil {
		return nil, fmt.Errorf("failed to get chain info: %w", err)
	}

	beholder, err := writetarget.NewMonitor(ctx, lggr)
	if err != nil {
		return nil, fmt.Errorf("failed to create Aptos WT monitor client: %+w", err)
	}

	opts := writetarget.WriteTargetOpts{
		ID:     id,
		Logger: lggr,
		Config: writetarget.Config{
			ConfirmerPollPeriod: pollPeriod,
			ConfirmerTimeout:    timeout,
		},
		ChainInfo:        chainInfo,
		Beholder:         beholder,
		ChainService:     chain,
		ContractReader:   cr,
		ChainWriter:      cw,
		ConfigValidateFn: evaluate,
		NodeAddress:      config.FromAddress().String(),
		ForwarderAddress: config.ForwarderAddress().String(),
		TargetStrategy:   NewEVMTargetStrategy(cr, cw, config.ForwarderAddress().String(), gasLimitDefault, lggr),
	}

	return writetarget.NewWriteTarget(opts), nil
}

type Inputs struct {
	SignedReport types.SignedReport
}

type Request struct {
	Metadata capabilities.RequestMetadata
	Config   Config
	Inputs   Inputs
}

func getEVMRequest(rawRequest capabilities.CapabilityRequest) (
	Request, error) {
	var r Request
	r.Metadata = rawRequest.Metadata

	if rawRequest.Config == nil {
		return Request{}, errors.New("missing config field")
	}

	if err := rawRequest.Config.UnwrapTo(&r.Config); err != nil {
		return Request{}, err
	}

	if !common.IsHexAddress(r.Config.Address) {
		return Request{}, fmt.Errorf("'%v' is not a valid address", r.Config.Address)
	}

	if rawRequest.Inputs == nil {
		return Request{}, errors.New("missing inputs field")
	}

	// required field of target's config in the workflow spec
	signedReport, ok := rawRequest.Inputs.Underlying[writetarget.KeySignedReport]
	if !ok {
		return Request{}, fmt.Errorf("missing required field %s", writetarget.KeySignedReport)
	}

	if err := signedReport.UnwrapTo(&r.Inputs.SignedReport); err != nil {
		return Request{}, err
	}
	return r, nil
}

func evaluate(rawRequest capabilities.CapabilityRequest) (receiver string, err error) {
	r, err := getEVMRequest(rawRequest)
	if err != nil {
		return "", err
	}

	reportMetadata, err := decodeReportMetadata(r.Inputs.SignedReport.Report)
	if err != nil {
		return "", err
	}

	if reportMetadata.Version != 1 {
		return "", fmt.Errorf("unsupported report version: %d", reportMetadata.Version)
	}

	if hex.EncodeToString(reportMetadata.WorkflowExecutionID[:]) != rawRequest.Metadata.WorkflowExecutionID {
		return "", fmt.Errorf("WorkflowExecutionID in the report does not match WorkflowExecutionID in the request metadata. Report WorkflowExecutionID: %+v, request WorkflowExecutionID: %+v", hex.EncodeToString(reportMetadata.WorkflowExecutionID[:]), rawRequest.Metadata.WorkflowExecutionID)
	}

	// case-insensitive verification of the owner address (so that a check-summed address matches its non-checksummed version).
	if !strings.EqualFold(hex.EncodeToString(reportMetadata.WorkflowOwner[:]), rawRequest.Metadata.WorkflowOwner) {
		return "", fmt.Errorf("WorkflowOwner in the report does not match WorkflowOwner in the request metadata. Report WorkflowOwner: %+v, request WorkflowOwner: %+v", hex.EncodeToString(reportMetadata.WorkflowOwner[:]), rawRequest.Metadata.WorkflowOwner)
	}

	// workflowNames are padded to 10bytes
	decodedName, err := hex.DecodeString(rawRequest.Metadata.WorkflowName)
	if err != nil {
		return "", err
	}
	var workflowName [10]byte
	copy(workflowName[:], decodedName)
	if !bytes.Equal(reportMetadata.WorkflowName[:], workflowName[:]) {
		return "", fmt.Errorf("WorkflowName in the report does not match WorkflowName in the request metadata. Report WorkflowName: %+v, request WorkflowName: %+v", hex.EncodeToString(reportMetadata.WorkflowName[:]), hex.EncodeToString(workflowName[:]))
	}

	if hex.EncodeToString(reportMetadata.WorkflowCID[:]) != rawRequest.Metadata.WorkflowID {
		return "", fmt.Errorf("WorkflowID in the report does not match WorkflowID in the request metadata. Report WorkflowID: %+v, request WorkflowID: %+v", reportMetadata.WorkflowCID, rawRequest.Metadata.WorkflowID)
	}

	if !bytes.Equal(reportMetadata.ReportID[:], r.Inputs.SignedReport.ID) {
		return "", fmt.Errorf("ReportID in the report does not match ReportID in the inputs. reportMetadata.ReportID: %x, Inputs.SignedReport.ID: %x", reportMetadata.ReportID, r.Inputs.SignedReport.ID)
	}

	return r.Config.Address, nil
}

func decodeReportMetadata(data []byte) (metadata ReportV1Metadata, err error) {
	if len(data) < metadata.Length() {
		return metadata, fmt.Errorf("data too short: %d bytes", len(data))
	}
	return metadata, binary.Read(bytes.NewReader(data[:metadata.Length()]), binary.BigEndian, &metadata)
}

func (rm ReportV1Metadata) Encode() ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, rm)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (rm ReportV1Metadata) Length() int {
	bytes, err := rm.Encode()
	if err != nil {
		return 0
	}
	return len(bytes)
}

// Note: This should be a shared type that the OCR3 package validates as well
type ReportV1Metadata struct {
	Version             uint8
	WorkflowExecutionID [32]byte
	Timestamp           uint32
	DonID               uint32
	DonConfigVersion    uint32
	WorkflowCID         [32]byte
	WorkflowName        [10]byte
	WorkflowOwner       [20]byte
	ReportID            [2]byte
}

func getChainInfo(chainID uint64) (monitor.ChainInfo, error) {
	chainSelector := chainselectors.EvmChainIdToChainSelector()[chainID]
	chainFamily, err := chainselectors.GetSelectorFamily(chainSelector)
	if err != nil {
		return monitor.ChainInfo{}, fmt.Errorf("failed to get chain family for selector %d: %w", chainSelector, err)
	}
	chainDetails, err := chainselectors.GetChainDetailsByChainIDAndFamily(strconv.Itoa(int(chainID)), chainFamily)
	if err != nil {
		return monitor.ChainInfo{}, fmt.Errorf("failed to get chain details for chain %d and family %s: %w", chainID, chainFamily, err)
	}

	neworkName, err := ExtractNetwork(chainDetails.ChainName)
	if err != nil {
		return monitor.ChainInfo{}, fmt.Errorf("failed to get network name for chain %d: %w", chainID, err)
	}

	return monitor.ChainInfo{
		ChainFamilyName: chainFamily,
		ChainID:         strconv.Itoa(int(chainID)),
		NetworkName:     neworkName,
		NetworkNameFull: chainDetails.ChainName,
	}, nil
}

func ExtractNetwork(selector string) (string, error) {
	// Create a regexp pattern that matches any of the three.
	re := regexp.MustCompile(`(mainnet|testnet|devnet)`)
	name := re.FindString(selector)
	if name == "" {
		return "", fmt.Errorf("failed to extract network name from selector: %s", selector)
	}
	return name, nil
}

func GenerateWriteTargetName(chainID uint64) string {
	id := fmt.Sprintf("write_%v@1.0.0", chainID)
	chainName, err := chainselectors.NameFromChainId(chainID)
	if err == nil {
		id = fmt.Sprintf("write_%v@1.0.0", chainName)
	}

	return id
}
