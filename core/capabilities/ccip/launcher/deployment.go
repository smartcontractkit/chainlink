package launcher

import (
	"golang.org/x/exp/maps"

	cctypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"go.uber.org/multierr"
)

type ccipPlugins map[ocrtypes.ConfigDigest]cctypes.CCIPOracle

// StartAll will call Oracle.Start on an entire don
func (c ccipPlugins) StartAll() error {
	nilPlugins := make(ccipPlugins)
	return c.Transition(nilPlugins)
}

// CloseAll is used to shut down an entire don immediately
func (c ccipPlugins) CloseAll() error {
	nilPlugins := make(ccipPlugins)
	return nilPlugins.Transition(c)
}

// Transition manages starting and stopping ocr instances
// If there are any new config digests, we need to start those instances
// If any of the previous config digests are no longer present, we need to shut those down
// We don't care about if they're exec/commit or active/candidate, that all happens in the plugin
func (c ccipPlugins) Transition(prevPlugins ccipPlugins) error {
	var err error

	// This shuts down instances that were present previously, but are no longer needed
	for _, digest := range maps.Keys(prevPlugins) {
		if c[digest] == nil {
			err = multierr.Append(err, prevPlugins[digest].Close())
		}
	}

	// This starts instances that were not previously present, but are in the new config
	for _, digest := range maps.Keys(c) {
		if prevPlugins[digest] == nil {
			err = multierr.Append(err, c[digest].Start())
		}
	}
	return err
}
