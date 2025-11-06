package state

import (
	"context"
	"encoding/json"
	"errors"
	"sync"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink/deployment/common/view"
	"github.com/smartcontractkit/chainlink/deployment/common/view/v2_0"
)

type CREChainView struct {
	CapabilityRegistry map[string]v2_0.CapabilityRegistryView `json:"capabilityRegistry,omitempty"`

	// TODO: add OCR, forwarders and workflow registry
}

func NewCREChainView() CREChainView {
	return CREChainView{
		CapabilityRegistry: make(map[string]v2_0.CapabilityRegistryView),
	}
}

// GenerateCREChainView is a view of the CRE contracts
// It is best-effort, logs errors and generates the views in parallel.
func GenerateCREChainView(
	ctx context.Context,
	lggr logger.Logger,
	prevView CREChainView,
	contracts viewContracts,
) (CREChainView, error) {
	out := NewCREChainView()
	var outMu sync.Mutex
	var allErrs error
	var wg sync.WaitGroup
	errCh := make(chan error, 1) // We are generating 1 view(s) concurrently

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
					cr := v2_0.ExtendedCapabilityRegistry{CapabilitiesRegistry: capabilitiesRegistry}
					addrCopy := addr
					capRegView, err := v2_0.GenerateCapabilityRegistryView(&cr)
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

type CREView struct {
	Chains map[string]CREChainView `json:"chains,omitempty"`
	Nops   map[string]view.NopView `json:"nops,omitempty"`
}

func (v CREView) MarshalJSON() ([]byte, error) {
	// Alias to avoid recursive calls
	type Alias CREView
	return json.MarshalIndent(&struct{ Alias }{Alias: Alias(v)}, "", " ")
}
