package launcher

import (
	"errors"
	"fmt"
	"sync"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"go.uber.org/multierr"

	cctypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"
)

// MaxPlugins is the maximum number of plugins possible.
// A plugin represents a possible combination of (active/candidate) x (commit/exec)
// If we ever have more than 4 plugins in a prev or desired state, something went wrong
const MaxPlugins = 4

type pluginRegistry map[ocrtypes.ConfigDigest]cctypes.CCIPOracle

// StartAll will call Oracle.Start on an entire don
func (c pluginRegistry) StartAll() error {
	emptyPluginRegistry := make(pluginRegistry)
	return c.TransitionFrom(emptyPluginRegistry)
}

// CloseAll is used to shut down an entire don immediately
func (c pluginRegistry) CloseAll() error {
	emptyPluginRegistry := make(pluginRegistry)
	return emptyPluginRegistry.TransitionFrom(c)
}

// TransitionFrom manages starting and stopping ocr instances
// If there are any new config digests, we need to start those instances
// If any of the previous config digests are no longer present, we need to shut those down
// We don't care about if they're exec/commit or active/candidate, that all happens in the plugin
func (c pluginRegistry) TransitionFrom(prevPlugins pluginRegistry) error {
	var allErrs error

	if len(c) > MaxPlugins || len(prevPlugins) > MaxPlugins {
		return fmt.Errorf("current pluginRegistry or prevPlugins have more than 4 instances: len(prevPlugins): %d, len(currPlugins): %d", len(prevPlugins), len(c))
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	// This shuts down instances that were present previously, but are no longer needed
	for digest, oracle := range prevPlugins {
		if _, ok := c[digest]; !ok {
			wg.Add(1)
			go func(o cctypes.CCIPOracle) {
				defer wg.Done()
				if err := o.Close(); err != nil {
					mu.Lock()
					allErrs = multierr.Append(allErrs, err)
					mu.Unlock()
				}
			}(oracle)
		}
	}
	wg.Wait()

	// This will start the instances that were not previously present, but are in the new config
	for digest, oracle := range c {
		if digest == [32]byte{} {
			allErrs = multierr.Append(allErrs, errors.New("cannot start a plugin with an empty config digest"))
		} else if _, ok := prevPlugins[digest]; !ok {
			wg.Add(1)
			go func(o cctypes.CCIPOracle) {
				defer wg.Done()
				if err := o.Start(); err != nil {
					mu.Lock()
					allErrs = multierr.Append(allErrs, err)
					mu.Unlock()
				}
			}(oracle)
		}
	}
	wg.Wait()

	return allErrs
}
