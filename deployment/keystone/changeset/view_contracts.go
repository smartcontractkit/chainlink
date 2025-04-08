package changeset

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3confighelper"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	capocr3types "github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/ocr3/types"

	forwarder "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/forwarder_1_0_0"
	ocr3_capability "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/ocr3_capability_1_0_0"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/view"
	common_v1_0 "github.com/smartcontractkit/chainlink/deployment/common/view/v1_0"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/internal"
)

type KeystoneChainView struct {
	CapabilityRegistry map[string]common_v1_0.CapabilityRegistryView `json:"capabilityRegistry,omitempty"`
	// OCRContracts is a map of OCR3 contract addresses to their configuration view
	OCRContracts     map[string]OCR3ConfigView                   `json:"ocrContracts,omitempty"`
	WorkflowRegistry map[string]common_v1_0.WorkflowRegistryView `json:"workflowRegistry,omitempty"`
	Forwarders       map[string][]ForwarderView                  `json:"forwarders,omitempty"`
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
	TxHash        string   `json:"txHash,omitempty"`
	BlockNumber   uint64   `json:"blockNumber,omitempty"`
}

var (
	ErrOCR3NotConfigured      = errors.New("OCR3 not configured")
	ErrForwarderNotConfigured = errors.New("forwarder not configured")
)

func GenerateOCR3ConfigView(ctx context.Context, ocr3Cap ocr3_capability.OCR3Capability) (OCR3ConfigView, error) {
	details, err := ocr3Cap.LatestConfigDetails(nil)
	if err != nil {
		return OCR3ConfigView{}, err
	}

	blockNumber := uint64(details.BlockNumber)
	configIterator, err := ocr3Cap.FilterConfigSet(&bind.FilterOpts{
		Start:   blockNumber,
		End:     &blockNumber,
		Context: ctx,
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

func GenerateForwarderView(ctx context.Context, f *forwarder.KeystoneForwarder, prevViews []ForwarderView) ([]ForwarderView, error) {
	startBlock := uint64(0)

	if len(prevViews) > 0 {
		// Sort `prevViews` by block number in ascending order, we make sure the last item has the highest block number
		sort.Slice(prevViews, func(i, j int) bool {
			return prevViews[i].BlockNumber < prevViews[j].BlockNumber
		})

		// If we have previous views, we will start from the last block number +1 of the previous views
		startBlock = prevViews[len(prevViews)-1].BlockNumber + 1
	} else {
		// If we don't have previous views, we will start from the deployment block number
		// which is stored in the forwarder's type and version labels.
		var deploymentBlock uint64
		lblPrefix := internal.DeploymentBlockLabel + ": "
		tvStr, err := f.TypeAndVersion(nil)
		if err != nil {
			return nil, fmt.Errorf("error getting TypeAndVersion for forwarder: %w", err)
		}
		tv, err := deployment.TypeAndVersionFromString(tvStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse type and version from %s: %w", tvStr, err)
		}
		for lbl := range tv.Labels {
			if strings.HasPrefix(lbl, lblPrefix) {
				// Extract the block number part after the prefix
				blockStr := strings.TrimPrefix(lbl, lblPrefix)
				blockNum, err := strconv.ParseUint(blockStr, 10, 64)
				if err == nil {
					deploymentBlock = blockNum
					break
				}
			}
		}

		if deploymentBlock > 0 {
			startBlock = deploymentBlock
		}
	}

	// Let's fetch the `SetConfig` events since the deployment block, since we don't have specific block numbers
	// for the `SetConfig` events.
	// If no deployment block is available, it will start from 0.
	configIterator, err := f.FilterConfigSet(&bind.FilterOpts{
		Start:   startBlock,
		End:     nil,
		Context: ctx,
	}, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("error filtering ConfigSet events: %w", err)
	}

	configSets := make([]*forwarder.KeystoneForwarderConfigSet, 0)
	for configIterator.Next() {
		// We wait for the iterator to receive an event
		if configIterator.Event == nil {
			// We cannot return an error, since we are capturing all `SetConfig` events, so if there's a nil event,
			// we ignore it.
			continue
		}
		configSets = append(configSets, configIterator.Event)
	}
	if len(configSets) == 0 {
		// Forwarder is not configured only if we don't have any previous configuration events.
		if len(prevViews) == 0 {
			return nil, ErrForwarderNotConfigured
		}

		// If we don't have any new config sets, we return the previous views as is.
		return prevViews, nil
	}

	// We now create a slice with all previous views and the new views, so they get all added to the final view.
	forwarderViews := append([]ForwarderView{}, prevViews...)
	for _, configSet := range configSets {
		var readableSigners []string
		for _, s := range configSet.Signers {
			readableSigners = append(readableSigners, s.String())
		}
		forwarderViews = append(forwarderViews, ForwarderView{
			DonID:         configSet.DonId,
			ConfigVersion: configSet.ConfigVersion,
			F:             configSet.F,
			Signers:       readableSigners,
			TxHash:        configSet.Raw.TxHash.String(),
			BlockNumber:   configSet.Raw.BlockNumber,
		})
	}

	return forwarderViews, nil
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
		Forwarders:         make(map[string][]ForwarderView),
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
