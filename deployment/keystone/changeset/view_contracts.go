package changeset

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	forwarder "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/forwarder_1_0_0"

	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/common/view"
	common_v1_0 "github.com/smartcontractkit/chainlink/deployment/common/view/v1_0"
	creforwarder "github.com/smartcontractkit/chainlink/deployment/cre/forwarder"
	"github.com/smartcontractkit/chainlink/deployment/cre/ocr3"
	creview "github.com/smartcontractkit/chainlink/deployment/cre/view"
)

type KeystoneChainView struct {
	CapabilityRegistry map[string]common_v1_0.CapabilityRegistryView `json:"capabilityRegistry,omitempty"`
	// OCRContracts is a map of OCR3 contract addresses to their configuration view
	OCRContracts     map[string]creview.OCR3ConfigView           `json:"ocrContracts,omitempty"`
	WorkflowRegistry map[string]common_v1_0.WorkflowRegistryView `json:"workflowRegistry,omitempty"`
	Forwarders       map[string][]ForwarderView                  `json:"forwarders,omitempty"`
}

type ForwarderView struct {
	DonID                   uint32   `json:"donId"`
	ConfigVersion           uint32   `json:"configVersion"`
	F                       uint8    `json:"f"`
	Signers                 []string `json:"signers"`
	TxHash                  string   `json:"txHash,omitempty"`
	BlockNumber             uint64   `json:"blockNumber,omitempty"`
	LatestViewedBlockNumber uint64   `json:"latestViewedBlockNumber,omitempty"`
}

var (
	ErrForwarderNotConfigured = errors.New("forwarder not configured")
)

// GenerateKeystoneChainView is a view of the keystone chain
// It is best-effort, logs errors and generates the views in parallel.
func GenerateKeystoneChainView(
	ctx context.Context,
	lggr logger.Logger,
	prevView KeystoneChainView,
	contracts viewContracts,
	chain cldf_evm.Chain,
) (KeystoneChainView, error) {
	out := NewKeystoneChainView()
	var outMu sync.Mutex
	var allErrs error
	var wg sync.WaitGroup
	errCh := make(chan error, 4) // We are generating 4 views concurrently

	// Check if context is already done before starting work
	select {
	case <-ctx.Done():
		return out, ctx.Err()
	default:
		// Continue processing
	}

	if contracts.CapabilitiesRegistry != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for addr, capabilitiesRegistry := range contracts.CapabilitiesRegistry {
				select {
				case <-ctx.Done():
					errCh <- ctx.Err()
					return
				default:
					cr := capabilitiesRegistry
					addrCopy := addr
					capRegView, err := common_v1_0.GenerateCapabilityRegistryView(cr)
					if err != nil {
						lggr.Warnf("failed to generate capability registry view for address %s: %v", addrCopy, err)
						errCh <- err
					}
					outMu.Lock()
					out.CapabilityRegistry[addrCopy.String()] = capRegView
					outMu.Unlock()
				}
			}
		}()
	}

	if contracts.OCR3 != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for addr, ocr3Cap := range contracts.OCR3 {
				select {
				case <-ctx.Done():
					errCh <- ctx.Err()
					return
				default:
					oc := *ocr3Cap
					addrCopy := addr
					ocrView, err := creview.GenerateOCR3ConfigView(ctx, oc)
					if err != nil {
						// don't block view on single OCR3 not being configured
						if errors.Is(err, creview.ErrOCR3NotConfigured) {
							lggr.Warnf("ocr3 not configured for address %s", addrCopy)
						} else {
							lggr.Errorf("failed to generate OCR3 config view for address %s: %v", addrCopy, err)
							errCh <- err
						}
						continue
					}
					outMu.Lock()
					out.OCRContracts[addrCopy.String()] = ocrView
					outMu.Unlock()
				}
			}
		}()
	}

	// Process the workflow registry and print if WorkflowRegistryError errors.
	if contracts.WorkflowRegistry != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for addr, workflowRegistry := range contracts.WorkflowRegistry {
				select {
				case <-ctx.Done():
					errCh <- ctx.Err()
					return
				default:
					wr := workflowRegistry
					addrCopy := addr
					wrView, wrErrs := common_v1_0.GenerateWorkflowRegistryView(wr)
					for _, err := range wrErrs {
						lggr.Errorf("WorkflowRegistry error for address %s: %v", addrCopy, err)
						errCh <- err
					}
					outMu.Lock()
					out.WorkflowRegistry[addrCopy.String()] = wrView
					outMu.Unlock()
				}
			}
		}()
	}

	if contracts.Forwarder != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for _, fwr := range contracts.Forwarder {
				fwrCopy := fwr
				fwrAddr := fwrCopy.Address().String()
				var prevViews []ForwarderView
				if prevView.Forwarders != nil {
					pv, ok := prevView.Forwarders[fwrAddr]
					if !ok {
						prevViews = []ForwarderView{}
					} else {
						prevViews = pv
					}
				} else {
					prevViews = []ForwarderView{}
				}

				select {
				case <-ctx.Done():
					errCh <- ctx.Err()
					return
				default:
					fwrView, fwrErr := GenerateForwarderView(ctx, fwrCopy, prevViews, chain)
					if fwrErr != nil {
						// don't block view on single forwarder not being configured
						switch {
						case errors.Is(fwrErr, ErrForwarderNotConfigured):
							lggr.Warnf("forwarder not configured for address %s", fwrCopy.Address())
						case errors.Is(fwrErr, context.Canceled), errors.Is(fwrErr, context.DeadlineExceeded):
							lggr.Warnf("forwarder view generation cancelled for address %s", fwrCopy.Address())
							errCh <- fwrErr
						default:
							lggr.Errorf("failed to generate forwarder view for address %s: %v", fwrCopy.Address(), fwrErr)
							errCh <- fwrErr
						}
					} else {
						outMu.Lock()
						out.Forwarders[fwrAddr] = fwrView
						outMu.Unlock()
					}
				}
			}
		}()
	}

	wg.Wait()
	close(errCh)

	var errList []error
	// Collect all errors
	for err := range errCh {
		errList = append(errList, err)
	}
	allErrs = errors.Join(errList...)

	return out, allErrs
}

func GenerateForwarderView(ctx context.Context, f *forwarder.KeystoneForwarder, prevViews []ForwarderView, chain cldf_evm.Chain) ([]ForwarderView, error) {
	startBlock := uint64(0)

	if len(prevViews) > 0 {
		// Sort `prevViews` by block number in ascending order, we make sure the last item has the highest block number
		sort.Slice(prevViews, func(i, j int) bool {
			// having `LatestViewedBlockNumber` means all prev views have it, so we use it for sorting
			if prevViews[i].LatestViewedBlockNumber > 0 {
				return prevViews[i].LatestViewedBlockNumber < prevViews[j].LatestViewedBlockNumber
			}
			return prevViews[i].BlockNumber < prevViews[j].BlockNumber
		})

		// If we have previous views, we will start from the last block number +1 of the previous views
		startBlock = prevViews[len(prevViews)-1].BlockNumber + 1
		if prevViews[len(prevViews)-1].LatestViewedBlockNumber > 0 {
			startBlock = prevViews[len(prevViews)-1].LatestViewedBlockNumber + 1
		}
	} else {
		// If we don't have previous views, we will start from the deployment block number
		// which is stored in the forwarder's type and version labels.
		var deploymentBlock uint64
		lblPrefix := creforwarder.DeploymentBlockLabel + ": "
		tvStr, err := f.TypeAndVersion(nil)
		if err != nil {
			return nil, fmt.Errorf("error getting TypeAndVersion for forwarder: %w", err)
		}
		tv, err := cldf.TypeAndVersionFromString(tvStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse type and version from %s: %w", tvStr, err)
		}
		for lbl := range tv.Labels {
			if after, ok := strings.CutPrefix(lbl, lblPrefix); ok {
				// Extract the block number part after the prefix
				blockStr := after
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

	// We'll iterate in chunks to avoid fetching too many logs at once.
	// We need the latest block number to know when to stop.
	l, err := chain.Client.HeaderByNumber(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest block number: %w", err)
	}
	latestBlock := l.Number.Uint64()

	configSets := make([]*forwarder.KeystoneForwarderConfigSet, 0)
	// We use a batch size of 50k blocks to avoid timeouts or limits
	batchSize := uint64(50000)

	// Let's fetch the `SetConfig` events since the deployment block, since we don't have specific block numbers
	// for the `SetConfig` events.
	// If no deployment block is available, it will start from 0.
	for start := startBlock; start <= latestBlock; start += batchSize {
		end := start + batchSize - 1
		if end > latestBlock {
			end = latestBlock
		}

		configIterator, configSetErr := f.FilterConfigSet(&bind.FilterOpts{
			Start:   start,
			End:     &end,
			Context: ctx,
		}, nil, nil)
		if configSetErr != nil {
			return nil, fmt.Errorf("error filtering ConfigSet events: %w", configSetErr)
		}

		for configIterator.Next() {
			// We wait for the iterator to receive an event
			if configIterator.Event == nil {
				// We cannot return an error, since we are capturing all `SetConfig` events, so if there's a nil event,
				// we ignore it.
				continue
			}
			configSets = append(configSets, configIterator.Event)
		}
	}
	updatedPrevViews := make([]ForwarderView, 0, len(prevViews)+len(configSets))
	for _, prevView := range prevViews {
		// Let's set the LatestViewedBlockNumber to the previous views
		prevView.LatestViewedBlockNumber = latestBlock
		updatedPrevViews = append(updatedPrevViews, prevView)
	}

	if len(configSets) == 0 {
		// Forwarder is not configured only if we don't have any previous configuration events.
		if len(updatedPrevViews) == 0 {
			return nil, ErrForwarderNotConfigured
		}

		// If we don't have any new config sets, we return the previous views as is.
		return updatedPrevViews, nil
	}

	// We now create a slice with all previous views and the new views, so they get all added to the final view.
	forwarderViews := updatedPrevViews
	for _, configSet := range configSets {
		var readableSigners []string
		for _, s := range configSet.Signers {
			readableSigners = append(readableSigners, s.String())
		}
		forwarderViews = append(forwarderViews, ForwarderView{
			DonID:                   configSet.DonId,
			ConfigVersion:           configSet.ConfigVersion,
			F:                       configSet.F,
			Signers:                 readableSigners,
			TxHash:                  configSet.Raw.TxHash.String(),
			BlockNumber:             configSet.Raw.BlockNumber,
			LatestViewedBlockNumber: latestBlock,
		})
	}

	return forwarderViews, nil
}

func NewKeystoneChainView() KeystoneChainView {
	return KeystoneChainView{
		CapabilityRegistry: make(map[string]common_v1_0.CapabilityRegistryView),
		OCRContracts:       make(map[string]creview.OCR3ConfigView),
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

type KeystoneViewV2 struct {
	Chains map[string]KeystoneChainView `json:"chains,omitempty"`

	// Nops will be DEPRECATED, but needs to stay here so we don't break downstream consumers, as tracking all the
	// consumers is not straightforward. We would need to list them for the rdd monster Kafka topic and BQ table, but
	// also list all the consumers who read the state.json file directly (are there any? we don't know).
	// The best way to go about it is to add a `nops_v2` entry and announce depreciation of the `nops` entry.
	Nops   map[string]view.NopView   `json:"nops,omitempty"`
	NopsV2 map[string]view.NopViewV2 `json:"nops_v2,omitempty"`
}

func (v KeystoneViewV2) MarshalJSON() ([]byte, error) {
	// Alias to avoid recursive calls
	type Alias KeystoneViewV2
	return json.MarshalIndent(&struct{ Alias }{Alias: Alias(v)}, "", " ")
}

// KeystoneChainViewLegacy is the legacy version of KeystoneChainView, which contains the legacy version of the OCR3 config view.
// This is used for auto-migration from the legacy view to the new view.
type KeystoneChainViewLegacy struct {
	CapabilityRegistry map[string]common_v1_0.CapabilityRegistryView `json:"capabilityRegistry,omitempty"`
	// OCRContracts is a map of OCR3 contract addresses to their configuration view
	OCRContracts     map[string]creview.OCR3ConfigViewLegacy     `json:"ocrContracts,omitempty"`
	WorkflowRegistry map[string]common_v1_0.WorkflowRegistryView `json:"workflowRegistry,omitempty"`
	Forwarders       map[string][]ForwarderView                  `json:"forwarders,omitempty"`
}

// Migrate migrates the legacy KeystoneChainView to the new KeystoneChainView.
// It converts the legacy OCR3 config view to the new OCR3 config view.
func (v KeystoneChainViewLegacy) Migrate() (KeystoneChainView, error) {
	newChainView := KeystoneChainView{
		CapabilityRegistry: v.CapabilityRegistry,
		OCRContracts:       make(map[string]creview.OCR3ConfigView),
		WorkflowRegistry:   v.WorkflowRegistry,
		Forwarders:         v.Forwarders,
	}

	for addr, legacyOCRView := range v.OCRContracts {
		addrCopy := addr
		newChainView.OCRContracts[addrCopy] = creview.OCR3ConfigView{
			Signers:               legacyOCRView.Signers,
			Transmitters:          legacyOCRView.Transmitters,
			F:                     legacyOCRView.F,
			OnchainConfig:         legacyOCRView.OnchainConfig,
			OffchainConfigVersion: legacyOCRView.OffchainConfigVersion,
			OffchainConfig: ocr3.OracleConfig{
				UniqueReports:                     legacyOCRView.OffchainConfig.UniqueReports,
				DeltaProgressMillis:               legacyOCRView.OffchainConfig.DeltaProgressMillis,
				DeltaResendMillis:                 legacyOCRView.OffchainConfig.DeltaResendMillis,
				DeltaInitialMillis:                legacyOCRView.OffchainConfig.DeltaInitialMillis,
				DeltaRoundMillis:                  legacyOCRView.OffchainConfig.DeltaRoundMillis,
				DeltaGraceMillis:                  legacyOCRView.OffchainConfig.DeltaGraceMillis,
				DeltaCertifiedCommitRequestMillis: legacyOCRView.OffchainConfig.DeltaCertifiedCommitRequestMillis,
				DeltaStageMillis:                  legacyOCRView.OffchainConfig.DeltaStageMillis,
				MaxRoundsPerEpoch:                 legacyOCRView.OffchainConfig.MaxRoundsPerEpoch,
				TransmissionSchedule:              legacyOCRView.OffchainConfig.TransmissionSchedule,
				MaxDurationQueryMillis:            legacyOCRView.OffchainConfig.MaxDurationQueryMillis,
				MaxDurationObservationMillis:      legacyOCRView.OffchainConfig.MaxDurationObservationMillis,
				MaxDurationShouldAcceptMillis:     legacyOCRView.OffchainConfig.MaxDurationShouldAcceptMillis,
				MaxDurationShouldTransmitMillis:   legacyOCRView.OffchainConfig.MaxDurationShouldTransmitMillis,
				MaxFaultyOracles:                  legacyOCRView.OffchainConfig.MaxFaultyOracles,
				ConsensusCapOffchainConfig: &ocr3.ConsensusCapOffchainConfig{
					MaxQueryLengthBytes:       legacyOCRView.OffchainConfig.MaxQueryLengthBytes,
					MaxObservationLengthBytes: legacyOCRView.OffchainConfig.MaxObservationLengthBytes,
					MaxReportLengthBytes:      legacyOCRView.OffchainConfig.MaxReportLengthBytes,
					MaxOutcomeLengthBytes:     legacyOCRView.OffchainConfig.MaxOutcomeLengthBytes,
					MaxReportCount:            legacyOCRView.OffchainConfig.MaxReportCount,
					MaxBatchSize:              legacyOCRView.OffchainConfig.MaxBatchSize,
					OutcomePruningThreshold:   legacyOCRView.OffchainConfig.OutcomePruningThreshold,
					RequestTimeout:            legacyOCRView.OffchainConfig.RequestTimeout,
				},
			},
		}
	}

	return newChainView, nil
}

type KeystoneViewLegacy struct {
	Chains map[string]KeystoneChainViewLegacy `json:"chains,omitempty"`
	Nops   map[string]view.NopView            `json:"nops,omitempty"`
}

func (v KeystoneViewLegacy) MarshalJSON() ([]byte, error) {
	// Alias to avoid recursive calls
	type Alias KeystoneViewLegacy
	return json.MarshalIndent(&struct{ Alias }{Alias: Alias(v)}, "", " ")
}
