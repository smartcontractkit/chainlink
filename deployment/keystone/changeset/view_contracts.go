package changeset

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3confighelper"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	capocr3types "github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/ocr3/types"

	"github.com/smartcontractkit/chainlink/deployment"
	forwarder "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/forwarder_1_0_0"
	ocr3_capability "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/ocr3_capability_1_0_0"

	"github.com/smartcontractkit/chainlink/deployment/common/view"
	common_v1_0 "github.com/smartcontractkit/chainlink/deployment/common/view/v1_0"
)

type KeystoneChainView struct {
	CapabilityRegistry map[string]common_v1_0.CapabilityRegistryView `json:"capabilityRegistry,omitempty"`
	// OCRContracts is a map of OCR3 contract addresses to their configuration view
	OCRContracts     map[string]OCR3ConfigView                   `json:"ocrContracts,omitempty"`
	WorkflowRegistry map[string]common_v1_0.WorkflowRegistryView `json:"workflowRegistry,omitempty"`
	Forwarders       map[string]ForwarderView                    `json:"forwarders,omitempty"`
}

type OCR3ConfigView struct {
	Signers               []string            `json:"signers"`
	Transmitters          []ocr2types.Account `json:"transmitters"`
	F                     uint8               `json:"f"`
	OnchainConfig         []byte              `json:"onchainConfig"`
	OffchainConfigVersion uint64              `json:"offchainConfigVersion"`
	OffchainConfig        OracleConfig        `json:"offchainConfig"`
}

type ForwarderView struct {
	DonID         uint32   `json:"donId"`
	ConfigVersion uint32   `json:"configVersion"`
	F             uint8    `json:"f"`
	Signers       []string `json:"signers"`
}

var (
	ErrOCR3NotConfigured      = errors.New("OCR3 not configured")
	ErrForwarderNotConfigured = errors.New("forwarder not configured")
)

func GenerateOCR3ConfigView(ocr3Cap ocr3_capability.OCR3Capability) (OCR3ConfigView, error) {
	details, err := ocr3Cap.LatestConfigDetails(nil)
	if err != nil {
		return OCR3ConfigView{}, err
	}

	blockNumber := uint64(details.BlockNumber)
	configIterator, err := ocr3Cap.FilterConfigSet(&bind.FilterOpts{
		Start:   blockNumber,
		End:     &blockNumber,
		Context: nil,
	})
	if err != nil {
		return OCR3ConfigView{}, err
	}
	var config *ocr3_capability.OCR3CapabilityConfigSet
	for configIterator.Next() {
		// We wait for the iterator to receive an event
		if configIterator.Event == nil {
			return OCR3ConfigView{}, ErrOCR3NotConfigured
		}

		config = configIterator.Event
	}
	if config == nil {
		return OCR3ConfigView{}, ErrOCR3NotConfigured
	}

	var signers []ocr2types.OnchainPublicKey
	var readableSigners []string
	for _, s := range config.Signers {
		signers = append(signers, s)
		readableSigners = append(readableSigners, hex.EncodeToString(s))
	}
	var transmitters []ocr2types.Account
	for _, t := range config.Transmitters {
		transmitters = append(transmitters, ocr2types.Account(t.String()))
	}
	// `PublicConfigFromContractConfig` returns the `ocr2types.PublicConfig` that contains all the `OracleConfig` fields we need, including the
	// report plugin config.
	publicConfig, err := ocr3confighelper.PublicConfigFromContractConfig(true, ocr2types.ContractConfig{
		ConfigDigest:          config.ConfigDigest,
		ConfigCount:           config.ConfigCount,
		Signers:               signers,
		Transmitters:          transmitters,
		F:                     config.F,
		OnchainConfig:         nil, // empty onChain config, currently we always use a nil onchain config when calling SetConfig
		OffchainConfigVersion: config.OffchainConfigVersion,
		OffchainConfig:        config.OffchainConfig,
	})
	if err != nil {
		return OCR3ConfigView{}, err
	}
	var cfg capocr3types.ReportingPluginConfig
	if err = proto.Unmarshal(publicConfig.ReportingPluginConfig, &cfg); err != nil {
		return OCR3ConfigView{}, err
	}
	oracleConfig := OracleConfig{
		MaxQueryLengthBytes:       cfg.MaxQueryLengthBytes,
		MaxObservationLengthBytes: cfg.MaxObservationLengthBytes,
		MaxReportLengthBytes:      cfg.MaxReportLengthBytes,
		MaxOutcomeLengthBytes:     cfg.MaxOutcomeLengthBytes,
		MaxReportCount:            cfg.MaxReportCount,
		MaxBatchSize:              cfg.MaxBatchSize,
		OutcomePruningThreshold:   cfg.OutcomePruningThreshold,
		RequestTimeout:            cfg.RequestTimeout.AsDuration(),
		UniqueReports:             true, // This is hardcoded to true in the OCR3 contract

		DeltaProgressMillis:               millisecondsToUint32(publicConfig.DeltaProgress),
		DeltaResendMillis:                 millisecondsToUint32(publicConfig.DeltaResend),
		DeltaInitialMillis:                millisecondsToUint32(publicConfig.DeltaInitial),
		DeltaRoundMillis:                  millisecondsToUint32(publicConfig.DeltaRound),
		DeltaGraceMillis:                  millisecondsToUint32(publicConfig.DeltaGrace),
		DeltaCertifiedCommitRequestMillis: millisecondsToUint32(publicConfig.DeltaCertifiedCommitRequest),
		DeltaStageMillis:                  millisecondsToUint32(publicConfig.DeltaStage),
		MaxRoundsPerEpoch:                 publicConfig.RMax,
		TransmissionSchedule:              publicConfig.S,

		MaxDurationQueryMillis:          millisecondsToUint32(publicConfig.MaxDurationQuery),
		MaxDurationObservationMillis:    millisecondsToUint32(publicConfig.MaxDurationObservation),
		MaxDurationShouldAcceptMillis:   millisecondsToUint32(publicConfig.MaxDurationShouldAcceptAttestedReport),
		MaxDurationShouldTransmitMillis: millisecondsToUint32(publicConfig.MaxDurationShouldTransmitAcceptedReport),

		MaxFaultyOracles: publicConfig.F,
	}

	return OCR3ConfigView{
		Signers:               readableSigners,
		Transmitters:          transmitters,
		F:                     config.F,
		OnchainConfig:         nil, // empty onChain config
		OffchainConfigVersion: config.OffchainConfigVersion,
		OffchainConfig:        oracleConfig,
	}, nil
}

func GenerateForwarderView(f *forwarder.KeystoneForwarder, chain deployment.Chain) (ForwarderView, error) {
	ctx := context.Background()
	// This could be effectively done with 2 other approaches:
	// 1. Fetching the transaction receipt of the contract deployment, getting the deployment block number,
	//    and extracting the config from the logs, but we don't have access to the transaction hash needed for this.
	// 2. Using `CodeAt()` to find the block number in which the contract was created, and use that.
	//    We would have to go from block number 0 to find it, which in the end is similar what's done here.
	configIterator, err := f.FilterConfigSet(&bind.FilterOpts{
		Start:   0,
		End:     nil,
		Context: ctx,
	}, nil, nil)
	if err != nil {
		return ForwarderView{}, fmt.Errorf("error filtering ConfigSet events: %w", err)
	}

	var configSet *forwarder.KeystoneForwarderConfigSet
	for configIterator.Next() {
		// We wait for the iterator to receive an event
		if configIterator.Event == nil {
			// Since we are going from the contract deployment block
			// to the latest block, we can't just return an error here
			// as we might not have reached the latest block yet
			// which may contain the config event.
			continue
		}
		configSet = configIterator.Event
	}
	if configSet == nil {
		return ForwarderView{}, ErrForwarderNotConfigured
	}

	var readableSigners []string
	for _, s := range configSet.Signers {
		readableSigners = append(readableSigners, s.String())
	}
	return ForwarderView{
		DonID:         configSet.DonId,
		ConfigVersion: configSet.ConfigVersion,
		F:             configSet.F,
		Signers:       readableSigners,
	}, nil
}

func millisecondsToUint32(dur time.Duration) uint32 {
	ms := dur.Milliseconds()
	if ms > int64(math.MaxUint32) {
		return math.MaxUint32
	}
	//nolint:gosec // disable G115 as it is practically impossible to overflow here
	return uint32(ms)
}

func NewKeystoneChainView() KeystoneChainView {
	return KeystoneChainView{
		CapabilityRegistry: make(map[string]common_v1_0.CapabilityRegistryView),
		OCRContracts:       make(map[string]OCR3ConfigView),
		WorkflowRegistry:   make(map[string]common_v1_0.WorkflowRegistryView),
		Forwarders:         make(map[string]ForwarderView),
	}
}

type KeystoneView struct {
	Chains map[string]KeystoneChainView `json:"chains,omitempty"`
	Nops   map[string]view.NopView      `json:"nops,omitempty"`
}

func (v KeystoneView) MarshalJSON() ([]byte, error) {
	// Alias to avoid recursive calls
	type Alias KeystoneView
	return json.MarshalIndent(&struct{ Alias }{Alias: Alias(v)}, "", " ")
}
