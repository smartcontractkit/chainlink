package changeset

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink/deployment"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	common_v1_0 "github.com/smartcontractkit/chainlink/deployment/common/view/v1_0"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/internal"

	capabilities_registry "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	forwarder "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/forwarder_1_0_0"
	ocr3_capability "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/ocr3_capability_1_0_0"
	workflow_registry "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/workflow_registry_wrapper"
)

type GetContractSetsRequest struct {
	Chains      map[uint64]deployment.Chain
	AddressBook deployment.AddressBook

	// Labels indicates the label set that a contract must include to be considered as a member
	// of the returned contract set.  By default, an empty label set implies that only contracts without
	// labels will be considered.  Otherwise, all labels must be on the contract (e.g., "label1" AND "label2").
	Labels []string
}

type GetContractSetsResponse struct {
	ContractSets map[uint64]ContractSet
}

type ContractSet struct {
	commonchangeset.MCMSWithTimelockState
	OCR3                 map[common.Address]*ocr3_capability.OCR3Capability
	Forwarder            *forwarder.KeystoneForwarder
	CapabilitiesRegistry *capabilities_registry.CapabilitiesRegistry
	WorkflowRegistry     *workflow_registry.WorkflowRegistry
}

func (cs ContractSet) Convert() internal.ContractSet {
	return internal.ContractSet{
		MCMSWithTimelockState: commonchangeset.MCMSWithTimelockState{
			MCMSWithTimelockContracts: cs.MCMSWithTimelockContracts,
		},
		Forwarder:            cs.Forwarder,
		WorkflowRegistry:     cs.WorkflowRegistry,
		OCR3:                 cs.OCR3,
		CapabilitiesRegistry: cs.CapabilitiesRegistry,
	}
}

func (cs ContractSet) TransferableContracts() []common.Address {
	var out []common.Address
	if cs.OCR3 != nil {
		for _, ocr := range cs.OCR3 {
			out = append(out, ocr.Address())
		}
	}
	if cs.Forwarder != nil {
		out = append(out, cs.Forwarder.Address())
	}
	if cs.CapabilitiesRegistry != nil {
		out = append(out, cs.CapabilitiesRegistry.Address())
	}
	if cs.WorkflowRegistry != nil {
		out = append(out, cs.WorkflowRegistry.Address())
	}
	return out
}

// View is a view of the keystone chain
// It is best-effort, logs errors and generates the views in parallel.
func (cs ContractSet) View(ctx context.Context, prevView KeystoneChainView, lggr logger.Logger) (KeystoneChainView, error) {
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

	if cs.CapabilitiesRegistry != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			select {
			case <-ctx.Done():
				errCh <- ctx.Err()
				return
			default:
				capRegView, err := common_v1_0.GenerateCapabilityRegistryView(cs.CapabilitiesRegistry)
				if err != nil {
					lggr.Warn("failed to generate capability registry view: %w", err)
					errCh <- err
				}
				outMu.Lock()
				out.CapabilityRegistry[cs.CapabilitiesRegistry.Address().String()] = capRegView
				outMu.Unlock()
			}
		}()
	}

	if cs.OCR3 != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for addr, ocr3Cap := range cs.OCR3 {
				select {
				case <-ctx.Done():
					errCh <- ctx.Err()
					return
				default:
					oc := *ocr3Cap
					addrCopy := addr
					ocrView, err := GenerateOCR3ConfigView(ctx, oc)
					if err != nil {
						// don't block view on single OCR3 not being configured
						if errors.Is(err, ErrOCR3NotConfigured) {
							lggr.Warnf("ocr3 not configured for address %s", addr)
						} else {
							lggr.Errorf("failed to generate OCR3 config view: %v", err)
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
	if cs.WorkflowRegistry != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			select {
			case <-ctx.Done():
				errCh <- ctx.Err()
				return
			default:
				wrView, wrErrs := common_v1_0.GenerateWorkflowRegistryView(cs.WorkflowRegistry)
				for _, err := range wrErrs {
					lggr.Errorf("WorkflowRegistry error: %v", err)
					errCh <- err
				}
				outMu.Lock()
				out.WorkflowRegistry[cs.WorkflowRegistry.Address().String()] = wrView
				outMu.Unlock()
			}
		}()
	}

	if cs.Forwarder != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fwrAddr := cs.Forwarder.Address().String()
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
				fwrView, fwrErr := GenerateForwarderView(ctx, cs.Forwarder, prevViews)
				if fwrErr != nil {
					// don't block view on single forwarder not being configured
					switch {
					case errors.Is(fwrErr, ErrForwarderNotConfigured):
						lggr.Warnf("forwarder not configured for address %s", cs.Forwarder.Address())
					case errors.Is(fwrErr, context.Canceled), errors.Is(fwrErr, context.DeadlineExceeded):
						lggr.Warnf("forwarder view generation cancelled for address %s", cs.Forwarder.Address())
						errCh <- fwrErr
					default:
						lggr.Errorf("failed to generate forwarder view: %v", fwrErr)
						errCh <- fwrErr
					}
				} else {
					outMu.Lock()
					out.Forwarders[fwrAddr] = fwrView
					outMu.Unlock()
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

func (cs ContractSet) GetOCR3Contract(addr *common.Address) (*ocr3_capability.OCR3Capability, error) {
	return getOCR3Contract(cs.OCR3, addr)
}

func GetContractSets(lggr logger.Logger, req *GetContractSetsRequest) (*GetContractSetsResponse, error) {
	resp := &GetContractSetsResponse{
		ContractSets: make(map[uint64]ContractSet),
	}
	for id, chain := range req.Chains {
		addrs, err := req.AddressBook.AddressesForChain(id)
		if err != nil {
			return nil, fmt.Errorf("failed to get addresses for chain %d: %w", id, err)
		}

		// Forwarder addresses now have informative labels, but we don't want them to be ignored if no labels are provided for filtering.
		// If labels are provided, just filter by those.
		forwarderAddrs := make(map[string]deployment.TypeAndVersion)
		if len(req.Labels) == 0 {
			for addr, tv := range addrs {
				if tv.Type == KeystoneForwarder {
					forwarderAddrs[addr] = tv
				}
			}
		}

		// TODO: we need to expand/refactor the way labeled addresses are filtered
		// see: https://smartcontract-it.atlassian.net/browse/CRE-363
		filtered := deployment.LabeledAddresses(addrs).And(req.Labels...)

		for addr, tv := range forwarderAddrs {
			filtered[addr] = tv
		}

		cs, err := loadContractSet(lggr, chain, filtered)
		if err != nil {
			return nil, fmt.Errorf("failed to load contract set for chain %d: %w", id, err)
		}
		resp.ContractSets[id] = *cs
	}
	return resp, nil
}

func loadContractSet(lggr logger.Logger, chain deployment.Chain, addresses map[string]deployment.TypeAndVersion) (*ContractSet, error) {
	var out ContractSet
	mcmsWithTimelock, err := commonchangeset.MaybeLoadMCMSWithTimelockChainState(chain, addresses)
	if err != nil {
		return nil, fmt.Errorf("failed to load mcms contract: %w", err)
	}
	out.MCMSWithTimelockState = *mcmsWithTimelock

	for addr, tv := range addresses {
		// todo handle versions
		switch tv.Type {
		case CapabilitiesRegistry:
			c, err := capabilities_registry.NewCapabilitiesRegistry(common.HexToAddress(addr), chain.Client)
			if err != nil {
				return nil, fmt.Errorf("failed to create capability registry contract from address %s: %w", addr, err)
			}
			out.CapabilitiesRegistry = c
		case KeystoneForwarder:
			c, err := forwarder.NewKeystoneForwarder(common.HexToAddress(addr), chain.Client)
			if err != nil {
				return nil, fmt.Errorf("failed to create forwarder contract from address %s: %w", addr, err)
			}
			out.Forwarder = c
		case OCR3Capability:
			c, err := ocr3_capability.NewOCR3Capability(common.HexToAddress(addr), chain.Client)
			if err != nil {
				return nil, fmt.Errorf("failed to create OCR3Capability contract from address %s: %w", addr, err)
			}
			if out.OCR3 == nil {
				out.OCR3 = make(map[common.Address]*ocr3_capability.OCR3Capability)
			}
			out.OCR3[common.HexToAddress(addr)] = c
		case WorkflowRegistry:
			c, err := workflow_registry.NewWorkflowRegistry(common.HexToAddress(addr), chain.Client)
			if err != nil {
				return nil, fmt.Errorf("failed to create OCR3Capability contract from address %s: %w", addr, err)
			}
			out.WorkflowRegistry = c
		default:
			lggr.Warnw("unknown contract type", "type", tv.Type)
			// ignore unknown contract types
		}
	}
	return &out, nil
}

// getOCR3Contract returns the OCR3 contract from the contract set.  By default, it returns the only
// contract in the set if there is no address specified.  If an address is specified, it returns the
// contract with that address.  If the address is specified but not found in the contract set, it returns
// an error.
func getOCR3Contract(contracts map[common.Address]*ocr3_capability.OCR3Capability, addr *common.Address) (*ocr3_capability.OCR3Capability, error) {
	// Fail if the OCR3 contract address is unspecified and there are multiple OCR3 contracts
	if addr == nil && len(contracts) > 1 {
		return nil, errors.New("OCR contract address is unspecified")
	}

	// Use the first OCR3 contract if the address is unspecified
	if addr == nil && len(contracts) == 1 {
		// use the first OCR3 contract
		for _, c := range contracts {
			return c, nil
		}
	}

	// Select the OCR3 contract by address
	if contract, ok := contracts[*addr]; ok {
		return contract, nil
	}

	addrSet := make([]string, 0, len(contracts))
	for a := range contracts {
		addrSet = append(addrSet, a.String())
	}

	// Fail if the OCR3 contract address is specified but not found in the contract set
	return nil, fmt.Errorf("OCR3 contract address %s not found in contract set %v", *addr, addrSet)
}
