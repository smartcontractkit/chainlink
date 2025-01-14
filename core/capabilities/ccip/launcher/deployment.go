package launcher

import (
	"fmt"

	mapset "github.com/deckarep/golang-set/v2"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"golang.org/x/exp/maps"
	"golang.org/x/sync/errgroup"

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
	if len(c) > MaxPlugins || len(prevPlugins) > MaxPlugins {
		return fmt.Errorf("current pluginRegistry or prevPlugins have more than 4 instances: len(prevPlugins): %d, len(currPlugins): %d", len(prevPlugins), len(c))
	}

	prevOracles := mapset.NewSet[ocrtypes.ConfigDigest](maps.Keys(prevPlugins)...)
	currOracles := mapset.NewSet[ocrtypes.ConfigDigest](maps.Keys(c)...)

	var ops = make([]syncAction, 0, 2*MaxPlugins)
	for digest := range prevOracles.Difference(currOracles).Iterator().C {
		ops = append(ops, syncAction{
			command: closeAction,
			oracle:  prevPlugins[digest],
		})
	}

	for digest := range currOracles.Difference(prevOracles).Iterator().C {
		ops = append(ops, syncAction{
			command: openAction,
			oracle:  c[digest],
		})
	}

	g := new(errgroup.Group)
	for _, op := range ops {
		op := op
		g.Go(func() error {
			if op.command == closeAction {
				if err := op.oracle.Close(); err != nil {
					return err
				}
			} else if op.command == openAction {
				if err := op.oracle.Start(); err != nil {
					return err
				}
			}
			return nil
		})
	}

	return g.Wait()
}

const (
	closeAction = iota
	openAction
)

type syncAction struct {
	command int
	oracle  cctypes.CCIPOracle
}
