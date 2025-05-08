package changeset

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	capabilities_registry "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	forwarder "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/forwarder_1_0_0"
	ocr3_capability "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/ocr3_capability_1_0_0"
	workflow_registry "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/workflow_registry_wrapper"
	common_v1_0 "github.com/smartcontractkit/chainlink/deployment/common/view/v1_0"

	"github.com/smartcontractkit/chainlink/deployment"
	commonview "github.com/smartcontractkit/chainlink/deployment/common/view"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/internal"
)

var _ deployment.ViewStateV2 = ViewKeystone

type contractsPerChain map[uint64]viewContracts

type viewContracts struct {
	OCR3                 map[common.Address]*ocr3_capability.OCR3Capability
	Forwarder            *forwarder.KeystoneForwarder
	CapabilitiesRegistry *capabilities_registry.CapabilitiesRegistry
	WorkflowRegistry     *workflow_registry.WorkflowRegistry
}

func ViewKeystone(e deployment.Environment, previousView json.Marshaler) (json.Marshaler, error) {
	lggr := e.Logger
	contractsMap, err := getContractsPerChain(e)
	// This is an unrecoverable error
	if err != nil {
		return nil, fmt.Errorf("failed to get contract sets: %w", err)
	}

	prevViewBytes, err := previousView.MarshalJSON()
	if err != nil {
		// just log the error, we don't need to stop the execution since the previous view is optional
		lggr.Warnf("failed to marshal previous keystone view: %v", err)
	}
	var prevView KeystoneView
	if len(prevViewBytes) == 0 {
		prevView.Chains = make(map[string]KeystoneChainView)
	} else if err = json.Unmarshal(prevViewBytes, &prevView); err != nil {
		lggr.Warnf("failed to unmarshal previous keystone view: %v", err)
		prevView.Chains = make(map[string]KeystoneChainView)
	}

	var viewErrs error
	chainViews := make(map[string]KeystoneChainView)
	for chainSel, contracts := range contractsMap {
		chainid, err := chainsel.ChainIdFromSelector(chainSel)
		if err != nil {
			err2 := fmt.Errorf("failed to resolve chain id for selector %d: %w", chainSel, err)
			lggr.Error(err2)
			viewErrs = errors.Join(viewErrs, err2)
			continue
		}
		chainName, err := chainsel.NameFromChainId(chainid)
		if err != nil {
			err2 := fmt.Errorf("failed to resolve chain name for chain id %d: %w", chainid, err)
			lggr.Error(err2)
			viewErrs = errors.Join(viewErrs, err2)
			continue
		}
		v, err := GenerateKeystoneChainView(
			e.GetContext(),
			e.Logger,
			prevView.Chains[chainName],
			contracts.CapabilitiesRegistry,
			contracts.OCR3,
			contracts.WorkflowRegistry,
			contracts.Forwarder,
		)
		if err != nil {
			err2 := fmt.Errorf("failed to view chain %s: %w", chainName, err)
			lggr.Error(err2)
			viewErrs = errors.Join(viewErrs, err2)
			// don't continue; add the partial view
		}
		chainViews[chainName] = v
	}
	nopsView, err := commonview.GenerateNopsView(e.Logger, e.NodeIDs, e.Offchain)
	if err != nil {
		err2 := fmt.Errorf("failed to view nops: %w", err)
		lggr.Error(err2)
		viewErrs = errors.Join(viewErrs, err2)
	}
	return &KeystoneView{
		Chains: chainViews,
		Nops:   nopsView,
	}, viewErrs
}

func getContractsPerChain(e deployment.Environment) (contractsPerChain, error) {
	// Cannot do a single `Filter` call because it appears to work as an AND filter.
	ocr3CapabilityContracts := e.DataStore.Addresses().Filter(
		datastore.AddressRefByType(datastore.ContractType(internal.OCR3Capability)),
	)
	workflowRegistryContracts := e.DataStore.Addresses().Filter(
		datastore.AddressRefByType(datastore.ContractType(internal.WorkflowRegistry)),
	)
	keystoneForwarderContracts := e.DataStore.Addresses().Filter(
		datastore.AddressRefByType(datastore.ContractType(internal.KeystoneForwarder)),
	)
	capabilitiesRegistryContracts := e.DataStore.Addresses().Filter(
		datastore.AddressRefByType(datastore.ContractType(internal.CapabilitiesRegistry)),
	)
	contractAddresses := make([]datastore.AddressRef, 0, len(ocr3CapabilityContracts)+len(workflowRegistryContracts)+
		len(keystoneForwarderContracts)+len(capabilitiesRegistryContracts))
	contractAddresses = append(contractAddresses, ocr3CapabilityContracts...)
	contractAddresses = append(contractAddresses, workflowRegistryContracts...)
	contractAddresses = append(contractAddresses, keystoneForwarderContracts...)
	contractAddresses = append(contractAddresses, capabilitiesRegistryContracts...)

	contracts := make(contractsPerChain)
	var errs error

	// Initialize all contract sets first
	for _, addr := range contractAddresses {
		if _, ok := contracts[addr.ChainSelector]; !ok {
			contracts[addr.ChainSelector] = viewContracts{
				OCR3: make(map[common.Address]*ocr3_capability.OCR3Capability),
			}
		}
	}

	// TODO: what happens when there are multiple contracts of the same type deployed on the same chain?
	// As of now, it assigns the last one found to the contract set.
	for _, contractAddress := range contractAddresses {
		chain, ok := e.Chains[contractAddress.ChainSelector]
		if !ok {
			errs = errors.Join(errs, fmt.Errorf("chain with selector %d not found", contractAddress.ChainSelector))
			continue
		}

		// Get a mutable copy of the ContractSet
		set := contracts[contractAddress.ChainSelector]

		switch contractAddress.Type {
		case datastore.ContractType(internal.CapabilitiesRegistry):
			ownedContract, err := GetOwnedContractV2[*capabilities_registry.CapabilitiesRegistry](
				e.DataStore.Addresses(), chain, contractAddress.Address,
			)
			if err != nil {
				errs = errors.Join(errs, fmt.Errorf("failed to retrieve capabilities registry contract at %s: %w",
					contractAddress.Address, err))
				continue
			}
			set.CapabilitiesRegistry = ownedContract.Contract

		case datastore.ContractType(internal.OCR3Capability):
			ownedContract, err := GetOwnedContractV2[*ocr3_capability.OCR3Capability](
				e.DataStore.Addresses(), chain, contractAddress.Address,
			)
			if err != nil {
				errs = errors.Join(errs, fmt.Errorf("failed to retrieve OCR3 capability contract at %s: %w",
					contractAddress.Address, err))
				continue
			}
			set.OCR3[common.HexToAddress(contractAddress.Address)] = ownedContract.Contract

		case datastore.ContractType(internal.KeystoneForwarder):
			ownedContract, err := GetOwnedContractV2[*forwarder.KeystoneForwarder](
				e.DataStore.Addresses(), chain, contractAddress.Address,
			)
			if err != nil {
				errs = errors.Join(errs, fmt.Errorf("failed to retrieve forwarder contract at %s: %w",
					contractAddress.Address, err))
				continue
			}
			set.Forwarder = ownedContract.Contract

		case datastore.ContractType(internal.WorkflowRegistry):
			ownedContract, err := GetOwnedContractV2[*workflow_registry.WorkflowRegistry](
				e.DataStore.Addresses(), chain, contractAddress.Address,
			)
			if err != nil {
				errs = errors.Join(errs, fmt.Errorf("failed to retrieve workflow registry contract at %s: %w",
					contractAddress.Address, err))
				continue
			}
			set.WorkflowRegistry = ownedContract.Contract
		}

		// Store the updated `contractSet` back in the map
		contracts[contractAddress.ChainSelector] = set
	}

	return contracts, errs
}

// GenerateKeystoneChainView is a view of the keystone chain
// It is best-effort, logs errors and generates the views in parallel.
func GenerateKeystoneChainView(
	ctx context.Context,
	lggr logger.Logger,
	prevView KeystoneChainView,
	capabilitiesRegistry *capabilities_registry.CapabilitiesRegistry,
	ocr3Contracts map[common.Address]*ocr3_capability.OCR3Capability,
	workflowRegistry *workflow_registry.WorkflowRegistry,
	forwarder *forwarder.KeystoneForwarder,
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

	if capabilitiesRegistry != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			select {
			case <-ctx.Done():
				errCh <- ctx.Err()
				return
			default:
				capRegView, err := common_v1_0.GenerateCapabilityRegistryView(capabilitiesRegistry)
				if err != nil {
					lggr.Warn("failed to generate capability registry view: %w", err)
					errCh <- err
				}
				outMu.Lock()
				out.CapabilityRegistry[capabilitiesRegistry.Address().String()] = capRegView
				outMu.Unlock()
			}
		}()
	}

	if ocr3Contracts != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for addr, ocr3Cap := range ocr3Contracts {
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
	if workflowRegistry != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			select {
			case <-ctx.Done():
				errCh <- ctx.Err()
				return
			default:
				wrView, wrErrs := common_v1_0.GenerateWorkflowRegistryView(workflowRegistry)
				for _, err := range wrErrs {
					lggr.Errorf("WorkflowRegistry error: %v", err)
					errCh <- err
				}
				outMu.Lock()
				out.WorkflowRegistry[workflowRegistry.Address().String()] = wrView
				outMu.Unlock()
			}
		}()
	}

	if forwarder != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fwrAddr := forwarder.Address().String()
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
				fwrView, fwrErr := GenerateForwarderView(ctx, forwarder, prevViews)
				if fwrErr != nil {
					// don't block view on single forwarder not being configured
					switch {
					case errors.Is(fwrErr, ErrForwarderNotConfigured):
						lggr.Warnf("forwarder not configured for address %s", forwarder.Address())
					case errors.Is(fwrErr, context.Canceled), errors.Is(fwrErr, context.DeadlineExceeded):
						lggr.Warnf("forwarder view generation cancelled for address %s", forwarder.Address())
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
